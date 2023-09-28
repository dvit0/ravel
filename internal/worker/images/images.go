package images

import (
	"context"
	"io"
	"os"
	"os/exec"

	"github.com/charmbracelet/log"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
	initPkg "github.com/valyentdev/ravel/internal/init"
	"github.com/valyentdev/ravel/internal/utils"
	"github.com/valyentdev/ravel/internal/worker/config"
	"github.com/valyentdev/ravel/internal/worker/docker"
)

type ImagesManager struct {
	docker *docker.Docker
}

type Image struct {
	id           string
	docker       *docker.Docker
	imageInspect types.ImageInspect
	Name         string `json:"name"`
}

func (imageManager ImagesManager) GetImage(name string) (*Image, error) {
	imageInspect, _, err := imageManager.docker.Client.ImageInspectWithRaw(context.Background(), name)
	if err != nil {
		log.Error("Image not found", "err", err)
		return nil, err
	}

	return &Image{
		id:           utils.NewId(),
		docker:       imageManager.docker,
		Name:         name,
		imageInspect: imageInspect,
	}, nil
}

func (imageManager *ImagesManager) PullImage(name string) error {
	return imageManager.docker.PullImage(name)
}

func (image *Image) Unpack(ctx context.Context, dest string) error {
	log.Info("Unpacking", "image", image.Name)

	containerId := image.id

	// Create a container from the image
	_, err := image.docker.ContainerCreate(ctx, &container.Config{
		Image: image.Name,
	}, &container.HostConfig{}, &network.NetworkingConfig{}, &v1.Platform{}, containerId)
	if err != nil {
		return err
	}

	// Export the container to a tar archive
	readClose, err := image.docker.ContainerExport(ctx, containerId)
	if err != nil {
		return err
	}
	defer readClose.Close()

	archivePath := config.RAVEL_TEMP_ARCHIVES + "/" + image.id + ".tar"
	archive, err := os.Create(archivePath)
	if err != nil {
		return err
	}
	defer archive.Close()

	_, err = io.Copy(archive, readClose)

	if err != nil {
		return err
	}

	// Extract the archive to the destination
	err = exec.Command("tar", "-xvf", archivePath, "-C", dest).Run()
	if err != nil {
		log.Error("Failed to extract the archive", "err", err)
		return err
	}

	return nil
}

func (image *Image) GetInitImageConfig() initPkg.ImageConfig {

	return initPkg.ImageConfig{
		Cmd:        image.imageInspect.Config.Cmd,
		Entrypoint: image.imageInspect.Config.Entrypoint,
		Env:        image.imageInspect.Config.Env,
		WorkingDir: image.imageInspect.Config.WorkingDir,
		User:       image.imageInspect.Config.User,
	}
}

func NewImagesManager() *ImagesManager {

	return &ImagesManager{
		docker: docker.NewDockerClient(),
	}
}
