package service

import (
	"bufio"
	"context"
	"net"
	"time"

	"github.com/docker/docker/client"
	"github.com/faiyaz032/gobox/internal/docker"
	"github.com/faiyaz032/gobox/internal/errors"
	"github.com/faiyaz032/gobox/internal/repository"
)

type Service struct {
	repository *repository.Repository
	apiClient  *client.Client
}

func NewService(repository *repository.Repository, apiClient *client.Client) *Service {
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
	item, err := s.repository.GetOne(ctx, sessionID)
	if err != nil {
		containerID, err := docker.CreateContainer(ctx, s.apiClient)
		if err != nil {
			return nil, errors.Wrap(err, 500, "failed to create container")
		}

		mapItem := &repository.SessionContainer{
			SessionID:   sessionID,
			ContainerID: containerID,
		}

		if err := s.repository.Create(ctx, mapItem); err != nil {
			return nil, errors.Wrap(err, 500, "failed to create session container mapping")
		}

		item = mapItem
	} else {
		paused, err := docker.IsContainerPaused(ctx, s.apiClient, item.ContainerID)
		if err == nil && paused {
			_ = docker.UnpauseContainer(ctx, s.apiClient, item.ContainerID)
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

	item, err := s.repository.GetOneByContainerID(ctx, containerID)
	if err != nil {
		return nil
	}

	item.IsPaused = true
	item.PausedAt = time.Now()

	if err := s.repository.Update(ctx, item); err != nil {
		return err
	}

	return nil
}

func (s *Service) CleanContainers(ctx context.Context) {
	containers, err := s.repository.GetAll(ctx)
	if err != nil || len(containers) == 0 {
		return
	}

	now := time.Now()
	threshold := now.Add(-1 * time.Minute)

	for _, container := range containers {
		if container.IsPaused && container.PausedAt.Before(threshold) {
			if err := s.repository.Delete(ctx, container.ID); err != nil {
				continue
			}
			_ = docker.RemoveContainer(ctx, s.apiClient, container.ContainerID)
		}
	}
}
