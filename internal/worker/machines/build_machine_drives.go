package machines

import (
	"context"
	"encoding/json"
	"os"

	"github.com/charmbracelet/log"
	initPkg "github.com/valyentdev/ravel/internal/init"
	"github.com/valyentdev/ravel/internal/utils"
	"github.com/valyentdev/ravel/internal/worker/drives"
	"github.com/valyentdev/ravel/pkg/types"
	"github.com/valyentdev/ravel/pkg/units"
	"golang.org/x/sync/errgroup"
)

func (machineManager *MachineManager) buildMachineDrives(machineSpec types.RavelMachineSpec) (initDriveId string, rootDriveId string, err error) {
	log.Info("Pulling image", "image", machineSpec.Image)
	err = machineManager.images.PullImage(machineSpec.Image)
	if err != nil {
		log.Error("Error pulling image", "error", err)
		return "", "", err
	}

	log.Info("Inspecting image", "image", machineSpec.Image)
	image, err := machineManager.images.GetImage(machineSpec.Image)
	if err != nil {
		log.Error("Error getting image", "error", err)
		return "", "", err
	}

	log.Info("Getting image config", "image", machineSpec.Image)
	imageInitConfig := image.GetInitImageConfig()

	g := new(errgroup.Group)
	var initDrive *drives.Drive
	var mainDrive *drives.Drive

	g.Go(func() error {
		log.Info("Building init drive")
		initDrive, err = machineManager.buildInitDrive(machineSpec, imageInitConfig)
		if err != nil {
			log.Error("Error building init drive", "error", err)
			return err
		}
		return nil
	})

	g.Go(func() error {
		log.Info("Building main drive")
		mainDrive, err = machineManager.buildMainDrive(machineSpec)
		if err != nil {
			log.Error("Error building main drive", "error", err)
			return err
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		log.Error("Error building drives", "error", err)
		if initDrive != nil {
			machineManager.drives.DeleteDrive(initDrive.Id)
		}
		if mainDrive != nil {
			machineManager.drives.DeleteDrive(mainDrive.Id)
		}
		return "", "", err
	}

	return initDrive.Id, mainDrive.Id, nil
}

func (machineManager *MachineManager) buildInitDrive(machineSpec types.RavelMachineSpec, imageConfig initPkg.ImageConfig) (*drives.Drive, error) {
	drive, err := machineManager.drives.CreateDrive(types.RavelDriveSpec{
		Name: "init",
		Size: 50 * units.MB,
	}, true)
	if err != nil {
		log.Error("Error creating init drive", "error", err)
		return nil, err
	}

	err = drive.Mount()
	if err != nil {
		log.Error("Failed to create valyent directory", "err", err)
		return nil, err
	}

	defer drive.Unmount()

	mountPath := drive.GetMountPath()

	err = os.Mkdir(mountPath+"/ravel", os.FileMode(0777))
	if err != nil {
		log.Error("Failed to create ravel directory", "err", err)
		return nil, err
	}

	_, err = utils.Copy(os.Getenv("INIT_BINARY"), mountPath+"/ravel/init")
	if err != nil {
		log.Error("Failed to copy init binary", "err", err)
		return nil, err
	}

	err = os.Chmod(mountPath+"/ravel/init", 0777)
	if err != nil {
		log.Error("Failed to chmod init binary", "err", err)
		return nil, err
	}

	initConfig := initPkg.NewInitConfig(imageConfig)

	initConfigJson, _ := json.Marshal(initConfig)

	err = os.WriteFile(mountPath+"/ravel/run.json", initConfigJson, os.FileMode(0777))
	if err != nil {
		log.Error("Failed to write init config", "err", err)
		return nil, err
	}

	return drive, nil
}

func (machineManager *MachineManager) buildMainDrive(machineSpec types.RavelMachineSpec) (*drives.Drive, error) {
	drive, err := machineManager.drives.CreateDrive(types.RavelDriveSpec{
		Name: "init",
		Size: 5 * units.GB,
	}, true)
	if err != nil {
		log.Error("Error creating main drive", "error", err)
		return nil, err
	}
	defer func() {
		if err != nil {
			machineManager.drives.DeleteDrive(drive.Id)
		}
	}()

	err = drive.Mount()
	if err != nil {
		log.Error("Failed to mount drive", "err", err)
		return nil, err
	}
	defer drive.Unmount()

	image, err := machineManager.images.GetImage(machineSpec.Image)
	if err != nil {
		log.Error("Failed to get image", "err", err)
		return nil, err
	}

	err = image.Unpack(context.Background(), drive.GetMountPath())
	if err != nil {
		log.Error("Failed to unpack image", "err", err)
		return nil, err
	}

	return drive, nil
}
