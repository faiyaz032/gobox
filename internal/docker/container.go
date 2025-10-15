package docker

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/faiyaz032/gobox/util"
)

func CreateContainer(apiClient *client.Client, ctx context.Context) (string, error) {
	containerName := util.GenerateContainerName("ubuntu")
	containerResp, err := apiClient.ContainerCreate(ctx,
		&container.Config{
			Image: "ubuntu:latest",
			Cmd:   []string{"sleep", "3600"},
			Tty:   true,
		},
		nil, nil, nil, containerName,
	)
	if err != nil {
		return "", fmt.Errorf("create container failed: %w", err)
	}
	return containerResp.ID, nil
}

func StartContainer(apiClient *client.Client, ctx context.Context, containerId string) error {
	if err := apiClient.ContainerStart(ctx, containerId, container.StartOptions{}); err != nil {
		return fmt.Errorf("start container failed: %w", err)
	}
	return nil
}

func AttachShell(apiClient *client.Client, ctx context.Context, containerID string) (types.HijackedResponse, error) {
	execResp, err := apiClient.ContainerExecCreate(ctx, containerID, container.ExecOptions{
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		Tty:          true,
		Cmd:          []string{"bash"},
	})
	if err != nil {
		return types.HijackedResponse{}, fmt.Errorf("create exec failed: %w", err)
	}

	hijackResp, err := apiClient.ContainerExecAttach(ctx, execResp.ID, container.ExecAttachOptions{Tty: true})
	if err != nil {
		return types.HijackedResponse{}, fmt.Errorf("attach exec failed: %w", err)
	}

	return hijackResp, nil
}

func CleanUP(apiClient *client.Client, ctx context.Context, containerID string) error {
	if err := apiClient.ContainerRemove(ctx, containerID, container.RemoveOptions{
		RemoveVolumes: true,
		Force:         true,
	}); err != nil {
		return fmt.Errorf("remove container failed: %w", err)
	}
	return nil
}
