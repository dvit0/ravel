package api

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/charmbracelet/log"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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

	handler := apigen.NewStrictHandler(workerServer, nil)
	e := echo.New()

	e.Use(middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		KeyLookup: "header:x-api-key",
		Validator: func(key string, c echo.Context) (bool, error) {
			return key == os.Getenv("API_KEY"), nil
		},
	}))

	e.Server.Addr = ":3000"

	apigen.RegisterHandlers(e, handler)

	workerServer.server = e.Server

	log.Info("Ravel worker api start listening on", "port", e.Server.Addr)
	e.Server.ListenAndServe()
}

func strptr(s string) *string {
	return &s
}

func (s *WorkerServer) ExitWorker(ctx context.Context, request apigen.ExitWorkerRequestObject) (apigen.ExitWorkerResponseObject, error) {
	defer func() {
		time.Sleep(1 * time.Second)
		s.server.Shutdown(context.Background())
	}()

	s.worker.Cleanup()

	return nil, nil
}

func (s *WorkerServer) ListMachines(ctx context.Context, request apigen.ListMachinesRequestObject) (apigen.ListMachinesResponseObject, error) {
	machines, err := s.worker.ListMachines()
	if err != nil {
		return nil, err
	}

	return apigen.ListMachines200JSONResponse{
		Machines: &machines,
	}, nil

}

func (s *WorkerServer) CreateMachine(ctx context.Context, request apigen.CreateMachineRequestObject) (apigen.CreateMachineResponseObject, error) {
	id, err := s.worker.CreateMachine(context.Background(), *request.Body)
	if err != nil {
		return nil, err
	}

	return apigen.CreateMachine201JSONResponse{
		MachineId: &id,
	}, nil
}

func (s *WorkerServer) DeleteMachine(ctx context.Context, request apigen.DeleteMachineRequestObject) (apigen.DeleteMachineResponseObject, error) {
	machine, found, err := s.worker.GetMachine(request.Id)
	if err != nil {
		return nil, err
	}

	if !found {
		return apigen.DeleteMachine404JSONResponse{
			Message: strptr("Machine not found"),
		}, nil
	}

	err = s.worker.DeleteMachine(machine.Id)
	if err != nil {
		return nil, err
	}

	return apigen.DeleteMachine204Response{}, nil
}

func (s *WorkerServer) GetMachine(ctx context.Context, request apigen.GetMachineRequestObject) (apigen.GetMachineResponseObject, error) {
	machine, found, err := s.worker.GetMachine(request.Id)
	if err != nil {
		return nil, err
	}

	if !found {
		return apigen.GetMachine404JSONResponse{
			Message: strptr("Machine not found"),
		}, nil
	}

	return apigen.GetMachine200JSONResponse(*machine), nil

}

func (s *WorkerServer) StartMachine(ctx context.Context, request apigen.StartMachineRequestObject) (apigen.StartMachineResponseObject, error) {
	machine, found, err := s.worker.GetMachine(request.Id)
	if err != nil {
		return nil, err
	}

	if !found {
		return apigen.StartMachine404JSONResponse{
			Message: strptr("Machine not found"),
		}, nil
	}

	if machine.Status == apigen.Running {
		return apigen.StartMachine200Response{}, nil
	}

	err = s.worker.StartMachine(machine.Id)
	if err != nil {
		return nil, err
	}

	return apigen.StartMachine200Response{}, nil

}

func (s *WorkerServer) StopMachine(ctx context.Context, request apigen.StopMachineRequestObject) (apigen.StopMachineResponseObject, error) {
	machine, found, err := s.worker.GetMachine(request.Id)
	if err != nil {
		return nil, err
	}

	if !found {
		return apigen.StopMachine404JSONResponse{
			Message: strptr("Machine not found"),
		}, nil
	}

	if machine.Status == apigen.Stopped {
		return apigen.StopMachine400JSONResponse{
			Message: strptr("Machine already stopped"),
		}, nil
	}

	err = s.worker.StopMachine(machine.Id)
	if err != nil {
		return nil, err
	}

	return apigen.StopMachine200Response{}, nil
}
