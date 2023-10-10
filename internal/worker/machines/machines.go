package machines

import (
	"github.com/valyentdev/ravel/internal/worker/docker"
	"github.com/valyentdev/ravel/internal/worker/driversmanager"
	"github.com/valyentdev/ravel/internal/worker/drives"
	"github.com/valyentdev/ravel/internal/worker/images"
	"github.com/valyentdev/ravel/internal/worker/store"
)

type MachineManager struct {
	driversManager *driversmanager.DriversManager
	drives         *drives.DrivesManager
	store          *store.Store
	images         *images.ImagesManager
}

func NewMachineManager(store *store.Store, docker *docker.Docker) *MachineManager {
	return &MachineManager{
		driversManager: driversmanager.NewDriversManager(store),
		drives:         drives.NewDrivesManager(store),
		store:          store,
		images:         images.NewImagesManager(),
	}
}
