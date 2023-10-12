package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	api "github.com/valyentdev/ravel/pkg/api/worker"
)

func newErrorFromResponse(body []byte) error {

	var errorResponse api.ErrorResponse
	err := json.Unmarshal(body, &errorResponse)
	if err != nil || errorResponse.Message == nil {
		return errors.New("unknown error")
	}

	return errors.New(*errorResponse.Message)
}

func unexpectedError() error {
	return errors.New("unexpected error")
}

func authError() error {
	return errors.New("authentication error, please check your API key")
}

type WorkerClient struct {
	client *api.ClientWithResponses
}

func NewWorkerClient(server string, apiKey string) (*WorkerClient, error) {
	provideApiKey := func(ctx context.Context, req *http.Request) error {
		req.Header.Set("x-api-key", apiKey)
		return nil
	}
	c, err := api.NewClientWithResponses(server, api.WithRequestEditorFn(provideApiKey))
	if err != nil {
		return nil, err
	}
	return &WorkerClient{client: c}, nil
}

func (c *WorkerClient) CreateMachine(ctx context.Context, machineSpec api.RavelMachineSpec) (machineId string, err error) {
	res, err := c.client.CreateMachineWithResponse(context.Background(), machineSpec)
	if err != nil {
		return "", err
	}

	if res.StatusCode() == 401 {
		return "", authError()
	}

	if res.JSON201 == nil {
		return "", newErrorFromResponse(res.Body)
	}

	return *res.JSON201.MachineId, nil
}

func (c *WorkerClient) ListMachines() ([]api.RavelMachine, error) {
	res, err := c.client.ListMachinesWithResponse(context.Background())
	if err != nil {
		return nil, err
	}

	if res.StatusCode() == 401 {
		return nil, authError()
	}

	if res.JSON200 == nil {
		return nil, unexpectedError()
	}

	return *res.JSON200.Machines, nil
}

func (c *WorkerClient) DeleteMachine(ctx context.Context, machineId string) error {
	res, err := c.client.DeleteMachineWithResponse(context.Background(), machineId)
	if err != nil {
		return err
	}

	if res.StatusCode() == 401 {
		return authError()
	}

	if res.StatusCode() != 204 {
		return newErrorFromResponse(res.Body)
	}

	return nil
}

func (c *WorkerClient) GetMachine(ctx context.Context, machineId string) (*api.RavelMachine, error) {
	res, err := c.client.GetMachineWithResponse(context.Background(), machineId)
	if err != nil {
		return nil, err
	}

	if res.StatusCode() == 401 {
		return nil, authError()
	}

	if res.JSON200 == nil {
		return nil, newErrorFromResponse(res.Body)
	}

	return res.JSON200, nil
}

func (c *WorkerClient) StartMachine(ctx context.Context, machineId string) error {
	res, err := c.client.StartMachineWithResponse(context.Background(), machineId)
	if err != nil {
		return err
	}

	if res.StatusCode() == 401 {
		return authError()
	}

	if res.StatusCode() != 204 {
		return newErrorFromResponse(res.Body)
	}

	return nil
}

func (c *WorkerClient) StopMachine(ctx context.Context, machineId string) error {
	res, err := c.client.StopMachineWithResponse(context.Background(), machineId)
	if err != nil {
		return err
	}

	if res.StatusCode() == 401 {
		return authError()
	}

	if res.StatusCode() != 204 {
		return newErrorFromResponse(res.Body)
	}

	return nil
}

func (c *WorkerClient) GetStreamedLogs(ctx context.Context, machineId string) {
	res, err := c.client.GetMachineLogs(context.Background(), machineId)
	if err != nil {
		return
	}

	reader := res.Body

	for {
		buf := make([]byte, 1024)
		n, err := reader.Read(buf)
		if err != nil {
			return
		}
		fmt.Println(string(buf[:n]))
	}

}
