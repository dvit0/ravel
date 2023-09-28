package drives

import (
	"os"
	"os/exec"

	"github.com/valyentdev/ravel/internal/worker/config"
	"github.com/valyentdev/ravel/pkg/types"
	"golang.org/x/sys/unix"
)

type Drive struct {
	*types.RavelDrive
}

func (drive *Drive) Mount() error {
	err := os.Mkdir(drive.GetMountPath(), os.FileMode(0666))
	if err != nil {
		return err
	}

	err = exec.Command("mount", "-o", "loop", drive.GetDrivePath(), drive.GetMountPath()).Run()
	if err != nil {
		return err
	}

	return nil
}

func (drive *Drive) Unmount() error {
	unix.Unmount(drive.GetMountPath(), unix.MNT_FORCE)
	err := exec.Command("umount", drive.GetMountPath()).Run()
	if err != nil {
		return err
	}

	err = os.RemoveAll(drive.GetMountPath())
	if err != nil {
		return err
	}

	return nil
}

func getMountPath(id string) string {
	return config.RAVEL_DRIVES_PATH + id + "/_mount"
}

func getDrivePath(id string) string {
	return config.RAVEL_DRIVES_PATH + id + "/drive.img"
}

func (drive *Drive) GetMountPath() string {
	return getMountPath(drive.Id)
}

func (drive *Drive) GetDrivePath() string {
	return getDrivePath(drive.Id)
}
