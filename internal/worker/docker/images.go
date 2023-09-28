package docker

import (
	"context"
	"io"
	"os"

	"github.com/docker/docker/api/types"
)

func (docker *Docker) PullImage(image string) error {
	reader, err := docker.ImagePull(context.Background(), image, types.ImagePullOptions{})
	if err != nil {
		return err
	}
	defer reader.Close()

	_, err = io.Copy(os.Stdout, reader)
	if err != nil {
		return err
	}

	return nil
}

func (docker *Docker) ImageExists(imageName string) bool {
	_, _, err := docker.ImageInspectWithRaw(context.Background(), imageName)

	return err == nil
}
