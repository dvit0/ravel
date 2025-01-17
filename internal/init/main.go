package init

import (
	"io/fs"
	"os"
	"strings"
	"syscall"

	"github.com/charmbracelet/log"
	"github.com/opencontainers/runc/libcontainer/system"
	"github.com/opencontainers/runc/libcontainer/user"
	"golang.org/x/sys/unix"
)

func Init() {
	log.Info("started init")

	defer func() {
		if err := reboot(); err != nil {
			log.Fatal(err)
		}
	}()

	pid, err := unix.Setsid()
	if err != nil {
		log.Fatal(err)
	}

	pgid, err := unix.Getpgid(1)
	if err != nil {
		log.Fatal(err)
	}

	log.Infof("pid: %d", pid)
	log.Infof("pgid: %d", pgid)

	log.Debug("decoding run.json file")
	config, err := DecodeMachine("/ravel/run.json")
	if err != nil {
		log.Fatalf("could not parse run.json file %v", err)
	}

	if err := mkdir("/dev", perm0755); err != nil {
		log.Fatal(err)
	}

	initialMnts := MakeInitialMounts(config.RootDevice)
	if err := initialMnts.Mount(); err != nil {
		log.Fatal(err)
	}

	log.Debug("Switching root")
	if err := os.Chdir("/newroot"); err != nil {
		log.Fatal(err)
	}

	if err := mount(".", "/", "", unix.MS_MOVE, ""); err != nil {
		log.Fatal(err)
	}

	if err := unix.Chroot("."); err != nil {
		log.Fatal(err)
	}

	if err := os.Chdir("/"); err != nil {
		log.Fatal(err)
	}

	mnts := MakeMounts()
	if err := mnts.Mount(); err != nil {
		log.Fatal(err)
	}

	if err := mkdir("/run/lock", fs.FileMode(^uint32(0))); err != nil {
		log.Fatalf("could not create /run/lock directory: %v", err)
	}

	if err := unix.Symlinkat("/proc/self/fd", 0, "/dev/fd"); err != nil {
		log.Fatal(err)
	}

	if err := unix.Symlinkat("/proc/self/fd/0", 0, "/dev/stdin"); err != nil {
		log.Fatal(err)
	}

	if err := unix.Symlinkat("/proc/self/fd/1", 0, "/dev/stdout"); err != nil {
		log.Fatal(err)
	}

	if err := unix.Symlinkat("/proc/self/fd/2", 0, "/dev/stderr"); err != nil {
		log.Fatal(err)
	}

	if err := mkdir("/root", unix.S_IRWXU); err != nil {
		log.Fatalf("could not create /root dir: %v", err)
	}

	cgroupMnt := MakeCgroupMounts()
	if err := cgroupMnt.Mount(); err != nil {
		log.Fatal(err)
	}

	if err := unix.Setrlimit(0, &unix.Rlimit{Cur: 10240, Max: 10240}); err != nil {
		log.Fatal(err)
	}

	// parse user and  group names
	username := config.ImageConfig.User
	if username == "" {
		username = "root"
	}
	usrSplit := strings.Split(username, ":")

	var group user.Group
	if len(usrSplit) < 1 {
		log.Fatal("no username set, something is terribly wrong!")
	} else if len(usrSplit) >= 2 {
		group, err = user.LookupGroup(usrSplit[1])
		if err != nil {
			log.Fatalf("group %s not found: %v", usrSplit[1], err)
		}
	}
	_ = group

	nixUser, err := user.LookupUser(usrSplit[0])
	if err != nil {
		log.Fatalf("user %s not found: %v", username, err)
	}

	if err := system.Setgid(nixUser.Gid); err != nil {
		log.Fatalf("unable to set group id: %v", err)
	}

	if err := system.Setuid(nixUser.Uid); err != nil {
		log.Fatalf("unable to set group id: %v", err)
	}

	if err := PopulateProcessEnv(config.ImageConfig.Env); err != nil {
		log.Fatal(err)
	}

	// set the home dir if not already set
	if envHome := os.Getenv("HOME"); envHome == "" {
		if err := os.Setenv("HOME", nixUser.Home); err != nil {
			log.Fatal("unable to set user home directory")
		}
	}

	if err := MountAdditionalDrives(config.Mounts, nixUser.Uid, nixUser.Gid); err != nil {
		log.Fatalf("error mounting drives: %v", err)
	}

	if err := unix.Sethostname([]byte(config.Hostname)); err != nil {
		log.Fatalf("error setting hostname: %v", err)
	}

	if err := mkdir("/etc", perm0755); err != nil {
		log.Fatalf("could not create /etc dir: %v", err)
	}

	if err := os.WriteFile("/etc/hostname", []byte(config.Hostname+"\n"), perm0755); err != nil {
		log.Fatalf("error writing /etc/hostname: %v", err)
	}

	if err := WriteEtcResolv(config.EtcResolv); err != nil {
		log.Fatal(err)
	}

	if err := WriteEtcHost(config.EtcHost); err != nil {
		log.Fatal(err)
	}

	if err := NetworkSetup(); err != nil {
		log.Error(err)
		err = nil
	}

	p, err := NewProcess(config)
	if err != nil {
		log.Fatal(err)
	}

	if err = p.Run(); err != nil {
		log.Printf("error running process: %v", err)
	}
}

func reboot() error {
	return syscall.Reboot(syscall.LINUX_REBOOT_CMD_RESTART)
}
