package api

import (
	"net/http"

	"github.com/charmbracelet/log"
	"github.com/labstack/echo/v4"
	"github.com/valyentdev/ravel/internal/worker"
)

type WorkerServer struct {
	worker *worker.Worker
	server *http.Server
}

func StartWorkerApi(worker *worker.Worker) {
	workerServer := &WorkerServer{
		worker: worker,
	}

	e := echo.New()

	e.GET("/api/v1/machines/:id/logs", workerServer.GetMachineLogs)

	e.GET("/api/v1/machines", workerServer.ListMachines)
	e.POST("/api/v1/machines", workerServer.CreateMachine)

	e.GET("/api/v1/machines/:id", workerServer.GetMachine)
	e.DELETE("/api/v1/machines/:id", workerServer.DeleteMachine)
	e.POST("/api/v1/machines/:id/start", workerServer.StartMachine)
	e.POST("/api/v1/machines/:id/stop", workerServer.StopMachine)
	e.POST("/api/v1/exit", workerServer.ExitWorker)

	e.Server.Addr = ":3000"
	workerServer.server = e.Server

	log.Info("Ravel worker api start listening on", "port", e.Server.Addr)
	e.Server.ListenAndServe()
}
