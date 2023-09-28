package docker

import (
	"log"

	"github.com/docker/docker/client"
)

type Docker struct {
	*client.Client
	RegistryAuth string
}

func NewDockerClient() *Docker {
	client, err := client.NewClientWithOpts(
		client.FromEnv,
		client.WithAPIVersionNegotiation(),
	)

	if err != nil {
		log.Fatal(err)
	}

	return &Docker{
		Client: client,
	}
}
