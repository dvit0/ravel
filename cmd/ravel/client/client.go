package workerclient

import (
	"github.com/valyentdev/ravel/pkg/client"
)

var Client *client.WorkerClient

func init() {
	var err error
	Client, err = client.NewWorkerClient("http://localhost:3000", "very secret password")
	if err != nil {
		panic(err)
	}
}
