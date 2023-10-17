package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
	apigen "github.com/valyentdev/ravel/pkg/api/worker"
)

func (s *WorkerServer) StartMachine(ctx echo.Context) error {
	machineId := ctx.Param("id")
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
