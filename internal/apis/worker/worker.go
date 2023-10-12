package api

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/charmbracelet/log"
	"github.com/labstack/echo/v4"
	"github.com/valyentdev/ravel/internal/worker"
	apigen "github.com/valyentdev/ravel/pkg/api/worker"
)

type WorkerServer struct {
	worker *worker.Worker
	server *http.Server
}

func StartWorkerApi(worker *worker.Worker) {
	workerServer := &WorkerServer{
		worker: worker,
	}

	wrapper := apigen.ServerInterfaceWrapper{
		Handler: workerServer,
	}
	e := echo.New()

	e.GET("/api/v1/machines/:id/logs", wrapper.GetMachineLogs)

	e.GET("/api/v1/machines", wrapper.ListMachines)
	e.POST("/api/v1/machines", wrapper.CreateMachine)

	e.GET("/api/v1/machines/:id", wrapper.GetMachine)
	e.DELETE("/api/v1/machines/:id", wrapper.DeleteMachine)
	e.POST("/api/v1/machines/:id/start", wrapper.StartMachine)
	e.POST("/api/v1/machines/:id/stop", wrapper.StopMachine)
	e.POST("/api/v1/exit", wrapper.ExitWorker)

	e.Server.Addr = ":3000"
	workerServer.server = e.Server

	log.Info("Ravel worker api start listening on", "port", e.Server.Addr)
	e.Server.ListenAndServe()
}

func (s *WorkerServer) ExitWorker(ctx echo.Context) error {
	defer func() {
		time.Sleep(1 * time.Second)
		s.server.Shutdown(context.Background())
	}()

	s.worker.Cleanup()

	return nil
}

type LogsResponseWriter struct {
	ctx echo.Context
}

func (r *LogsResponseWriter) Write(p []byte) (n int, err error) {

	if string(p) == "\n" {
		return 0, nil
	}
	n, err = r.ctx.Response().Write(p)
	r.ctx.Response().Flush()
	return n, err
}

func (s *WorkerServer) GetMachineLogs(ctx echo.Context, machineId string) error {
	machine, found, err := s.worker.GetMachine(machineId)

	if err != nil {
		return echo.ErrInternalServerError
	}

	if !found {
		return echo.ErrNotFound
	}

	if machine.Status != apigen.Running {
		return echo.ErrBadRequest
	}
	ctx.Response().Header().Set("Content-Type", "text/event-stream")
	ctx.Response().Header().Set("Cache-Control", "no-cache")
	ctx.Response().Header().Set("Connection", "keep-alive")

	logFile, err := os.ReadFile("/var/log/ravel/machines/" + machineId + "/machine.log")
	if err != nil {
		return echo.ErrInternalServerError
	}

	responseWriter := &LogsResponseWriter{
		ctx: ctx,
	}

	responseWriter.Write(logFile)

	logBroadcaster := s.worker.LogsManager.GetLogBroadcaster(machineId)
	if logBroadcaster == nil {
		return nil
	}

	logBroadcaster.Subscribe(ctx.Request().Context(), responseWriter)

	return nil
}

func (s *WorkerServer) ListMachines(ctx echo.Context) error {
	machines, err := s.worker.ListMachines()
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, map[string][]apigen.RavelMachine{
		"machines": machines,
	})
}

func (s *WorkerServer) CreateMachine(ctx echo.Context) error {
	var body apigen.CreateMachineJSONRequestBody

	err := ctx.Bind(&body)
	if err != nil {
		return err
	}

	id, err := s.worker.CreateMachine(context.Background(), body)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusCreated, map[string]string{
		"machineId": id,
	})

}

func (s *WorkerServer) DeleteMachine(ctx echo.Context, machineId string) error {
	machine, found, err := s.worker.GetMachine(machineId)
	if err != nil {
		return err
	}
	if !found {
		return echo.ErrNotFound
	}

	err = s.worker.DeleteMachine(machine.Id)
	if err != nil {
		return err
	}

	return ctx.NoContent(http.StatusNoContent)
}

func (s *WorkerServer) GetMachine(ctx echo.Context, machineId string) error {
	machine, found, err := s.worker.GetMachine(machineId)
	if err != nil {
		return err
	}

	if !found {
		return echo.ErrNotFound
	}

	return ctx.JSON(http.StatusOK, machine)

}

func (s *WorkerServer) StartMachine(ctx echo.Context, machineId string) error {
	machine, found, err := s.worker.GetMachine(machineId)
	if err != nil {
		return err
	}

	if !found {
		return echo.ErrNotFound
	}

	if machine.Status == apigen.Running {
		return echo.NewHTTPError(http.StatusBadRequest, "Machine already started")
	}

	err = s.worker.StartMachine(machine.Id)
	if err != nil {
		return err
	}

	return ctx.NoContent(http.StatusNoContent)

}

func (s *WorkerServer) StopMachine(ctx echo.Context, machineId string) error {
	machine, found, err := s.worker.GetMachine(machineId)
	if err != nil {
		return err
	}

	if !found {
		return echo.ErrNotFound
	}

	if machine.Status == apigen.Stopped {
		return echo.NewHTTPError(http.StatusBadRequest, "Machine already stopped")
	}

	err = s.worker.StopMachine(machine.Id)
	if err != nil {
		return err
	}

	return ctx.NoContent(http.StatusNoContent)
}
