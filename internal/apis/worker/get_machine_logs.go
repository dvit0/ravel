package api

import (
	"os"

	"github.com/labstack/echo/v4"
	apigen "github.com/valyentdev/ravel/pkg/api/worker"
)

func (s *WorkerServer) GetMachineLogs(ctx echo.Context) error {
	machineId := ctx.Param("id")
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
