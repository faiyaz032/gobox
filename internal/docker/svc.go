package docker

import (
	"archive/tar"
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/containerd/errdefs"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/google/uuid"
)

type Svc struct {
	client *client.Client
}

func NewSvc() (*Svc, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	return &Svc{client: cli}, nil
}

func (s *Svc) Close() error {
	if s.client != nil {
		return s.client.Close()
	}
	return nil
}

func (s *Svc) EnsureImage(ctx context.Context, imageName, dockerfilePath string) error {
	_, err := s.client.ImageInspect(ctx, imageName)
	if err != nil {
		if errdefs.IsNotFound(err) {
			if err := s.BuildBaseImage(ctx, dockerfilePath, imageName); err != nil {
				return err
			}
		} else {
			return err
		}
	}

	return nil
}

func (s *Svc) EnsureNetwork(ctx context.Context, networkName, subnet string) (string, error) {
	nwList, err := s.client.NetworkList(ctx, network.ListOptions{})
	if err != nil {
		return "", fmt.Errorf("failed to list networks: %w", err)
	}

	for _, nw := range nwList {
		if nw.Name == networkName {
			return nw.ID, nil
		}
	}

	resp, err := s.client.NetworkCreate(ctx, networkName, network.CreateOptions{
		Driver: "bridge",
		IPAM: &network.IPAM{
			Driver: "default",
			Config: []network.IPAMConfig{
				{
					Subnet: subnet,
				},
			},
		},
	})
	if err != nil {
		return "", fmt.Errorf("failed to create network: %w", err)
	}
	return resp.ID, nil
}

func (s *Svc) BuildBaseImage(ctx context.Context, contextDir string, imageName string) error {
	buf := new(bytes.Buffer)
	tw := tar.NewWriter(buf)
	defer tw.Close()

	err := filepath.Walk(contextDir, func(file string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		header, err := tar.FileInfoHeader(fi, file)
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(contextDir, file)
		if err != nil {
			return err
		}

		header.Name = relPath

		if err := tw.WriteHeader(header); err != nil {
			return err
		}

		if fi.Mode().IsRegular() {
			f, err := os.Open(file)
			if err != nil {
				return err
			}
			defer f.Close()

			if _, err := io.Copy(tw, f); err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return err
	}

	res, err := s.client.ImageBuild(ctx, buf, types.ImageBuildOptions{
		Tags:       []string{imageName},
		Dockerfile: "Dockerfile",
		Remove:     true,
	})
	if err != nil {
		return err
	}
	defer res.Body.Close()

	_, err = io.Copy(os.Stdout, res.Body)
	return err
}

func (s *Svc) CreateContainer(ctx context.Context) (string, error) {
	// generate container name
	containerName := fmt.Sprintf("box-%s", uuid.New().String())
	// create container with resource limits
	resp, err := s.client.ContainerCreate(ctx, &container.Config{
		Image:     "gobox-base:latest",
		Cmd:       []string{"bash"},
		Tty:       true,
		OpenStdin: true,
		Hostname:  "box",
	}, &container.HostConfig{
		AutoRemove: false,
		Resources: container.Resources{
			Memory:            256 * 1024 * 1024, // 256MB
			MemoryReservation: 128 * 1024 * 1024,
			MemorySwap:        256 * 1024 * 1024,
			NanoCPUs:          500000000, // 0.5 CPU cores
			CPUShares:         512,       // Half of default priority
			BlkioWeight:       300,
		},
		StorageOpt: map[string]string{
			"size": "512MB", // 512MB
		},
	}, &network.NetworkingConfig{
		EndpointsConfig: map[string]*network.EndpointSettings{
			"gobox-c-network": {},
		},
	}, nil, containerName)
	if err != nil {
		return "", fmt.Errorf("failed to create container: %w", err)
	}
	return resp.ID, nil
}

func (s *Svc) AttachContainer(ctx context.Context, containerID string) (types.HijackedResponse, error) {
	attachResp, err := s.client.ContainerAttach(ctx, containerID, container.AttachOptions{
		Stream: true,
		Stdin:  true,
		Stdout: true,
		Stderr: true,
		Logs:   true,
	})
	if err != nil {
		return types.HijackedResponse{}, fmt.Errorf("failed to attach to container: %w", err)
	}

	return attachResp, nil
}

func (s *Svc) StartIfNotRunning(ctx context.Context, containerID string) error {
	inspect, err := s.client.ContainerInspect(ctx, containerID)
	if err != nil {
		return fmt.Errorf("failed to inspect container: %w", err)
	}

	if !inspect.State.Running {
		if err := s.client.ContainerStart(ctx, containerID, container.StartOptions{}); err != nil {
			return fmt.Errorf("failed to start container: %w", err)
		}
	}

	return nil
}

func (s *Svc) StopContainer(ctx context.Context, containerID string) error {
	timeout := 10
	stopOptions := container.StopOptions{
		Timeout: &timeout,
	}

	if err := s.client.ContainerStop(ctx, containerID, stopOptions); err != nil {
		return fmt.Errorf("failed to stop container: %w", err)
	}

	return nil
}
