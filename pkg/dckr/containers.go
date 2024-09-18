package dckr

import (
	"context"
	"io"

	"github.com/docker/docker/api/types"
)

func (dc *DockerClient) ListContainers() ([]types.Container, error) {

	containers, err := dc.cli.ContainerList(context.Background(), dc.containerListArgs)

	if err != nil {
		return nil, err
	}

	return containers, nil
}

func (dc *DockerClient) CopyToContainer(ctx context.Context, containerID, dstPath string, content io.Reader, options types.CopyToContainerOptions) error {
	return dc.cli.CopyToContainer(ctx, containerID, dstPath, content, options)
}
