package docker

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/faiyaz032/gobox/util"
)

func createNetwork(apiClient *client.Client, ctx context.Context, containerName string) (string, error) {
	networkName := fmt.Sprintf("net-%s", containerName)
	networkResp, err := apiClient.NetworkCreate(ctx, networkName, network.CreateOptions{
		Driver: "bridge",
	})
	if err != nil {
		return "", fmt.Errorf("create network failed: %w", err)
	}

	return networkResp.ID, nil
}

func CreateContainer(apiClient *client.Client, ctx context.Context) (string, error) {

	containerName := util.GenerateContainerName("ubuntu")

	networkID, err := createNetwork(apiClient, ctx, containerName)
	if err != nil {
		return "", err
	}

	containerResp, err := apiClient.ContainerCreate(ctx,
		&container.Config{
			Image: "gobox-base:latest",
			Cmd:   []string{"sleep", "3600"},
			Tty:   true,
		},
		&container.HostConfig{
			CapDrop:        []string{"ALL"},
			CapAdd:         []string{"CAP_NET_RAW"},
			ReadonlyRootfs: true,
			NetworkMode:    container.NetworkMode(networkID),
			Resources: container.Resources{
				Memory:   256 * 1024 * 1024,
				NanoCPUs: 500000000,
			},
		},
		nil, nil, containerName,
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

func RemoveContainer(apiClient *client.Client, ctx context.Context, containerID string) error {
	if err := apiClient.ContainerRemove(ctx, containerID, container.RemoveOptions{
		RemoveVolumes: true,
		Force:         true,
	}); err != nil {
		return fmt.Errorf("remove container failed: %w", err)
	}
	return nil
}

func PauseContainer(apiClient *client.Client, ctx context.Context, containerID string) error {
	if err := apiClient.ContainerPause(ctx, containerID); err != nil {
		return fmt.Errorf("pause container failed: %w", err)
	}
	return nil
}

func UnpauseContainer(ctx context.Context, apiClient *client.Client, containerID string) error {
	if err := apiClient.ContainerUnpause(ctx, containerID); err != nil {
		return fmt.Errorf("unpause container failed: %w", err)
	}
	return nil
}
