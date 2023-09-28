package machines

import (
	"github.com/valyentdev/ravel/internal/worker/docker"
	"github.com/valyentdev/ravel/internal/worker/drives"
	vmm_manager "github.com/valyentdev/ravel/internal/worker/firecracker_manager"
	"github.com/valyentdev/ravel/internal/worker/images"
	"github.com/valyentdev/ravel/internal/worker/store"
)

type MachineManager struct {
	drives   *drives.DrivesManager
	machines *vmm_manager.VMMManager
	store    *store.Store
	images   *images.ImagesManager
}

func NewMachineManager(store *store.Store, docker *docker.Docker) *MachineManager {
	return &MachineManager{
		drives:   drives.NewDrivesManager(store),
		machines: vmm_manager.NewMachinesManager(store),
		store:    store,
		images:   images.NewImagesManager(),
	}
}
