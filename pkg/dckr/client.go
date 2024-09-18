package dckr

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/docker/cli/cli/config"
	"github.com/docker/cli/cli/connhelper"
	"github.com/docker/cli/cli/context/store"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	// Add this new import
)

type DockerClient struct {
	cli               client.CommonAPIClient
	containerListArgs container.ListOptions
}

func (dc *DockerClient) Close() error {
	return dc.cli.Close()
}

func NewDockerClient() (*DockerClient, error) {
	socket, err := GetActiveDockerContext()
	if err != nil {
		return nil, fmt.Errorf("failed to get active Docker context: %w", err)
	}

	cli, err := client.NewClientWithOpts(
		client.WithHost(socket),
		client.WithAPIVersionNegotiation(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create Docker client: %w", err)
	}

	return &DockerClient{
		cli: cli,
		containerListArgs: container.ListOptions{
			Size:   true,
			All:    false,
			Latest: false,
		},
	}, nil
}

func NewDockerClientWithSSH(ssh string) (*DockerClient, error) {
	helper, err := connhelper.GetConnectionHelper(ssh)

	if err != nil {
		return nil, fmt.Errorf("failed to create SSH Docker client: %w", err)
	}

	httpClient := &http.Client{
		// No tls
		// No proxy
		Transport: &http.Transport{
			DialContext: helper.Dialer,
		},
	}

	cli, err := client.NewClientWithOpts(
		client.WithHTTPClient(httpClient),
		client.WithHost(helper.Host),
		client.WithDialContext(helper.Dialer),
		client.WithAPIVersionNegotiation(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create Docker client: %w", err)
	}

	return &DockerClient{
		cli: cli,
		containerListArgs: container.ListOptions{
			Size:   true,
			All:    false,
			Latest: false,
		},
	}, nil
}

type DockerContext struct {
	Description      string
	AdditionalFields map[string]any
}

// Add this new function
func ListDockerContexts() ([]store.Metadata, error) {
	contextStore := store.New(config.ContextStoreDir(), store.NewConfig(func() any { return &DockerContext{} }))
	if contextStore == nil {
		return nil, fmt.Errorf("failed to create context store")
	}

	return contextStore.List()

}

func GetActiveDockerContext() (string, error) {
	contexts, err := ListDockerContexts()
	if err != nil {
		return "", fmt.Errorf("failed to list Docker contexts: %w", err)
	}
	for _, ctx := range contexts {
		data := ctx.Endpoints
		ep := data["docker"].(map[string]any)
		host := GetHostFromEndpoint(ep)

		if strings.HasPrefix(host, "unix://") {
			if _, err := os.Stat(host[7:]); err == nil {
				return host, nil
			}
		}
	}

	return "", fmt.Errorf("failed to find active Docker context")
}

func GetHostFromEndpoint(endpoint map[string]any) string {
	host, ok := endpoint["Host"].(string)
	if !ok {
		return ""
	}
	return host
}
