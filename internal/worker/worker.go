package worker

import (
	"github.com/charmbracelet/log"
	"github.com/valyentdev/ravel/internal/worker/config"
	"github.com/valyentdev/ravel/internal/worker/docker"
	"github.com/valyentdev/ravel/internal/worker/machines"
	"github.com/valyentdev/ravel/internal/worker/store"
)

type Worker struct {
	store *store.Store
	*machines.MachineManager
}

func NewWorker() *Worker {
	log.SetPrefix("[Worker]")
	docker := docker.NewDockerClient()

	log.Info("Initializing store...")
	store := store.NewStore()

	log.Info("Initializing ravel directories...")
	config.InitRavelDirectories()

	log.Info("Worker is ready !")

	return &Worker{
		store:          &store,
		MachineManager: machines.NewMachineManager(&store, docker),
	}
}

func (worker *Worker) Cleanup() {
	err := worker.store.CloseStore()

	if err != nil {
		log.Error("Error closing store", "error", err)
	}
}
