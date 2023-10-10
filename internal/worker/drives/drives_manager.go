package drives

import (
	"os"
	"os/exec"

	"github.com/charmbracelet/log"
	"github.com/valyentdev/ravel/internal/utils"
	"github.com/valyentdev/ravel/internal/worker/config"
	"github.com/valyentdev/ravel/internal/worker/store"
	"github.com/valyentdev/ravel/pkg/types"
)

type DrivesManager struct {
	store *store.Store
}

func NewDrivesManager(store *store.Store) *DrivesManager {
	return &DrivesManager{
		store: store,
	}
}

func (drivesManager *DrivesManager) CreateDrive(driveConfig types.RavelDriveSpec, internal bool) (*Drive, error) {
	driveId := utils.NewId()
	err := os.Mkdir(config.RAVEL_DRIVES_PATH+driveId, os.FileMode(0666))
	if err != nil {
		log.Error("Error creating drive directory", "error", err)
		return nil, err
	}

	defer func() {
		if err != nil {
			os.RemoveAll(config.RAVEL_DRIVES_PATH + driveId)
		}
	}()

	err = utils.Fallocate(getDrivePath(driveId), driveConfig.Size)
	if err != nil {
		log.Error("Error allocating drive space", "error", err)
		return nil, err
	}

	err = exec.Command("mkfs.ext4", getDrivePath(driveId)).Run()
	if err != nil {
		log.Error("Error formatting drive", "error", err)
		return nil, err
	}

	drive, err := drivesManager.store.StoreRavelDrive(driveId, driveConfig, internal)
	if err != nil {
		log.Error("Error storing drive", "error", err)
		return nil, err
	}

	return &Drive{
		RavelDrive: drive,
	}, nil
}

func (drivesManager *DrivesManager) DeleteDrive(id string) error {
	drive, err := drivesManager.store.GetRavelDrive(id)
	if err != nil {
		return err
	}

	err = os.RemoveAll(config.RAVEL_DRIVES_PATH + drive.Id)
	if err != nil {
		return err
	}

	return nil
}

func (drivesManager *DrivesManager) GetDrive(id string) *Drive {
	drive, err := drivesManager.store.GetRavelDrive(id)
	if err != nil {
		return nil
	}

	return &Drive{
		RavelDrive: &drive,
	}
}
