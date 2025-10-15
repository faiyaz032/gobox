// internal/docker/image.go
package docker

import (
	"context"
	"fmt"
	"io"

	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
)

func PullImage(apiClient *client.Client, ctx context.Context) error {
	out, err := apiClient.ImagePull(ctx, "ubuntu:latest", image.PullOptions{})
	if err != nil {
		return fmt.Errorf("pull image failed: %w", err)
	}
	defer out.Close()
	_, _ = io.Copy(io.Discard, out)
	return nil
}
