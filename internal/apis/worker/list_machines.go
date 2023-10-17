package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/valyentdev/ravel/pkg/types"
)

func (s *WorkerServer) ListMachines(ctx echo.Context) error {
	machines, err := s.worker.ListMachines()
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, map[string][]types.RavelMachine{
		"machines": machines,
	})
}
