package docker

import (
	"archive/tar"
	"bytes"
	"context"
	"io"
	"os"
	"path/filepath"

	"github.com/containerd/errdefs"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
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
