package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (s *WorkerServer) GetMachine(ctx echo.Context) error {
	machineId := ctx.Param("id")
	machine, found, err := s.worker.GetMachine(machineId)
	if err != nil {
		return err
	}

	if !found {
		return echo.ErrNotFound
	}

	return ctx.JSON(http.StatusOK, machine)

}
