package dckr

import (
	"context"
	"testing"

	"github.com/docker/docker/api/types/container"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListContainers(t *testing.T) {
	dc, err := NewDockerClient()
	require.NoError(t, err, "Failed to create Docker client")
	defer dc.Close()

	ctx := context.Background()

	err = dc.cli.ContainerStart(ctx, "test-container", container.StartOptions{})
	require.NoError(t, err, "Failed to start test container")

	// Call the function
	containers, err := dc.ListContainers()

	// Assert the results
	assert.NoError(t, err, "ListContainers should not return an error")
	assert.NotEmpty(t, containers, "ListContainers should return at least one container")

}
