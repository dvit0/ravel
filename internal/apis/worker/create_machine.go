package api

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/valyentdev/ravel/pkg/types"
)

type CreateMachineBody = types.RavelMachineSpec

type CreateMachineResponse struct {
	MachineID string `json:"machine_id"` // Id of the created machine
}

func (s *WorkerServer) CreateMachine(ctx echo.Context) error {
	var body CreateMachineBody

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
