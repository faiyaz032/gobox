package docker

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/faiyaz032/gobox/internal/errors"
	"github.com/faiyaz032/gobox/internal/util"
)

func CreateContainer(ctx context.Context, apiClient *client.Client) (string, error) {
	containerName := util.GenerateContainerName("ubuntu")

	containerResp, err := apiClient.ContainerCreate(ctx,
		&container.Config{
			Image: "gobox-base:latest",
			Cmd:   []string{"sleep", "3600"},
			Tty:   true,
		},
		&container.HostConfig{
			CapDrop:     []string{"ALL"},
			CapAdd:      []string{"CAP_NET_RAW"},
			NetworkMode: "bridge",
			Resources: container.Resources{
				Memory:   256 * 1024 * 1024,
				NanoCPUs: 500_000_000,
			},
		},
		nil, nil, containerName,
	)
	if err != nil {
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

func IsContainerPaused(ctx context.Context, apiClient *client.Client, containerID string) (bool, error) {
	containerJSON, err := apiClient.ContainerInspect(ctx, containerID)
	if err != nil {
		return false, errors.Wrap(err, 500, "failed to inspect container")
	}

	return containerJSON.State.Paused, nil
}
