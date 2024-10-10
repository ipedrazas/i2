package dckr

import (
	"context"
	"fmt"

	"github.com/compose-spec/compose-go/cli"
	"github.com/compose-spec/compose-go/types"
	"github.com/docker/docker/client"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
)

func StartComposeServices(ctx context.Context, name, composeFile string) error {
	project, err := loadComposeFile(composeFile)
	if err != nil {
		return err
	}

	// Create Docker client
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}

	// Run containers in detached mode
	for _, svc := range project.Services {
		fmt.Printf("Starting container: %s in detached mode\n", svc.Name)
		resp, err := cli.ContainerCreate(context.Background(), &container.Config{
			Image: svc.Image,
		}, nil, nil, nil, svc.Name)
		if err != nil {
			return err
		}

		err = cli.ContainerStart(context.Background(), resp.ID, container.StartOptions{})
		if err != nil {
			return err
		}
	}

	return nil
}

// Load the compose file
func loadComposeFile(filename string) (*types.Project, error) {
	fmt.Println(filename)
	var opts *cli.ProjectOptions
	opts, err := cli.NewProjectOptions([]string{filename})
	if err != nil {
		fmt.Println("Load Option Errors", err)
		return nil, err
	}
	fmt.Println("Load Options", opts)

	project, err := cli.ProjectFromOptions(opts)
	if err != nil {
		return nil, err
	}
	fmt.Println("Load Project", project)
	return project, nil
}

func PullComposeImages(ctx context.Context, name, composeFile string) error {
	project, err := loadComposeFile(composeFile)
	if err != nil {
		return err
	}

	// Create Docker client
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}

	// Pull images
	options := image.PullOptions{}
	for _, svc := range project.Services {
		fmt.Printf("Pulling image: %s\n", svc.Image)
		_, err := cli.ImagePull(context.Background(), svc.Image, options)
		if err != nil {
			return err
		}
	}
	return nil
}
