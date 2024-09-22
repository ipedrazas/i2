package dckr

import (
	"context"
	"fmt"
	"io"
	"strings"

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

func PortsAsString(ports []types.Port) string {
	var portsString []string
	for _, port := range ports {
		portsString = append(portsString, fmt.Sprintf("%s:%d->%d/%s", port.IP, port.PublicPort, port.PrivatePort, port.Type))
	}
	return strings.Join(portsString, ", ")
}
