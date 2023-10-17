package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/valyentdev/ravel/pkg/types"
)

func (s *WorkerServer) StopMachine(ctx echo.Context) error {
	machineId := ctx.Param("id")
	machine, found, err := s.worker.GetMachine(machineId)
	if err != nil {
		return err
	}

	if !found {
		return echo.ErrNotFound
	}

	if machine.Status == types.RavelMachineStatusStopped {
		return echo.NewHTTPError(http.StatusBadRequest, "Machine already stopped")
	}

	err = s.worker.StopMachine(machine.Id)
	if err != nil {
		return err
	}

	return ctx.NoContent(http.StatusNoContent)
}
