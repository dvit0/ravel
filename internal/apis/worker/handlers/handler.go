package handlers

import "github.com/valyentdev/ravel/internal/worker"

type Handler struct {
	worker *worker.Worker
}

func NewHandler(worker *worker.Worker) *Handler {
	return &Handler{
		worker: worker,
	}
}
