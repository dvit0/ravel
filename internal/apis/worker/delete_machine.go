package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (s *WorkerServer) DeleteMachine(ctx echo.Context) error {
	machineId := ctx.Param("id")

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
