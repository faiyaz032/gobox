package box

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/faiyaz032/gobox/internal/domain"
)

type Repo interface {
	Create(context.Context, domain.Box) (*domain.Box, error)
	GetByFingerprint(context.Context, string) (*domain.Box, error)
	GetByContainerID(context.Context, string) (*domain.Box, error)
	Touch(context.Context, string) (*domain.Box, error)
	UpdateStatus(context.Context, string, string) (*domain.Box, error)
}

type DockerSvc interface {
	CreateContainer(ctx context.Context) (string, error)
	AttachContainer(ctx context.Context, containerID string) (types.HijackedResponse, error)
	StartIfNotRunning(ctx context.Context, containerID string) error
}
