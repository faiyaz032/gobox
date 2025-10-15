package main

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/docker/docker/client"
	"github.com/faiyaz032/gobox/internal/docker"
	"github.com/moby/term"
)

func main() {
	ctx := context.Background()
	apiClient, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating Docker client: %v\n", err)
		return
	}
	defer apiClient.Close()

	if err := docker.PullImage(apiClient, ctx); err != nil {
		fmt.Fprintf(os.Stderr, "Error pulling image: %v\n", err)
		return
	}

	containerId, err := docker.CreateContainer(apiClient, ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating container: %v\n", err)
		return
	}

	if err := docker.StartContainer(apiClient, ctx, containerId); err != nil {
		fmt.Fprintf(os.Stderr, "Error starting container: %v\n", err)
		return
	}

	fd := os.Stdin.Fd()
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error setting raw mode: %v\n", err)
		return
	}
	defer term.RestoreTerminal(fd, oldState)

	hijackResp, err := docker.AttachShell(apiClient, ctx, containerId)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error attaching shell: %v\n", err)
		return
	}

	go io.Copy(hijackResp.Conn, os.Stdin)
	io.Copy(os.Stdout, hijackResp.Reader)

	defer docker.CleanUP(apiClient, ctx, containerId)
}
