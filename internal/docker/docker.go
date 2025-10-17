package docker

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/faiyaz032/gobox/util"
)

func createNetwork(ctx context.Context, apiClient *client.Client, containerName string) (string, error) {
	networkName := util.GenerateContainerName(containerName)
	networkResp, err := apiClient.NetworkCreate(ctx, networkName, network.CreateOptions{
		Driver: "bridge",
	})
	if err != nil {
		return "", err
	}
	return networkResp.ID, nil
}

func removeNetwork(ctx context.Context, apiClient *client.Client, networkID string) error {
	return apiClient.NetworkRemove(ctx, networkID)
}

func CreateContainer(ctx context.Context, apiClient *client.Client) (string, error) {
	containerName := util.GenerateContainerName("ubuntu")

	networkID, err := createNetwork(ctx, apiClient, containerName)
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
			CapDrop:     []string{"ALL"},
			CapAdd:      []string{"CAP_NET_RAW"},
			NetworkMode: container.NetworkMode(networkID),
			Resources: container.Resources{
				Memory:   256 * 1024 * 1024,
				NanoCPUs: 500000000,
			},
		},
		nil, nil, containerName,
	)
	if err != nil {
		removeNetwork(ctx, apiClient, networkID)
		return "", err
	}

	return containerResp.ID, nil
}

func StartContainer(ctx context.Context, apiClient *client.Client, containerID string) error {
	return apiClient.ContainerStart(ctx, containerID, container.StartOptions{})
}

func AttachShell(ctx context.Context, apiClient *client.Client, containerID string) (types.HijackedResponse, error) {
	execResp, err := apiClient.ContainerExecCreate(ctx, containerID, container.ExecOptions{
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		Tty:          true,
		Cmd:          []string{"bash"},
	})
	if err != nil {
		return types.HijackedResponse{}, err
	}

	return apiClient.ContainerExecAttach(ctx, execResp.ID, container.ExecAttachOptions{Tty: true})
}

func RemoveContainer(ctx context.Context, apiClient *client.Client, containerID string) error {
	return apiClient.ContainerRemove(ctx, containerID, container.RemoveOptions{
		RemoveVolumes: true,
		Force:         true,
	})
}

func PauseContainer(ctx context.Context, apiClient *client.Client, containerID string) error {
	return apiClient.ContainerPause(ctx, containerID)
}

func UnpauseContainer(ctx context.Context, apiClient *client.Client, containerID string) error {
	return apiClient.ContainerUnpause(ctx, containerID)
}
