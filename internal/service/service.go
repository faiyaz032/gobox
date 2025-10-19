package service

import (
	"bufio"
	"context"
	"net"
	"time"

	"github.com/docker/docker/client"
	"github.com/faiyaz032/gobox/internal/docker"
	"github.com/faiyaz032/gobox/internal/errors"
	"github.com/faiyaz032/gobox/internal/infra/database"
)

type Service struct {
	repository *database.Repository
	apiClient  *client.Client
}

func NewService(repository *database.Repository, apiClient *client.Client) *Service {
	return &Service{
		repository: repository,
		apiClient:  apiClient,
	}
}

type StartResponse struct {
	ContainerID string
	Conn        net.Conn
	Reader      *bufio.Reader
}

func (s *Service) Start(ctx context.Context, sessionID string) (*StartResponse, error) {
	item, ok := s.repository.Get(ctx, sessionID)
	if !ok {
		containerID, err := docker.CreateContainer(ctx, s.apiClient)
		if err != nil {
			return nil, errors.Wrap(err, 500, "failed to create container")
		}

		mapItem := database.SessionContainer{
			ContainerID: containerID,
			LastActive:  time.Now(),
		}
		s.repository.Set(ctx, sessionID, mapItem)
		item = mapItem
	} else {
		if err := docker.UnpauseContainer(ctx, s.apiClient, item.ContainerID); err != nil {
			return nil, errors.Wrap(err, 500, "failed to unpause container")
		}
	}

	if err := docker.StartContainer(ctx, s.apiClient, item.ContainerID); err != nil {
		return nil, errors.Wrap(err, 500, "failed to start container")
	}

	hijackResp, err := docker.AttachShell(ctx, s.apiClient, item.ContainerID)
	if err != nil {
		return nil, errors.Wrap(err, 500, "failed to attach shell")
	}

	return &StartResponse{
		ContainerID: item.ContainerID,
		Conn:        hijackResp.Conn,
		Reader:      hijackResp.Reader,
	}, nil
}

func (s *Service) Pause(ctx context.Context, containerID string) error {
	if err := docker.PauseContainer(ctx, s.apiClient, containerID); err != nil {
		return errors.Wrap(err, 500, "failed to pause container")
	}
	return nil
}
