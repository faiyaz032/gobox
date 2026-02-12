package box

import (
	"context"

	"github.com/faiyaz032/gobox/internal/domain"
)

type Repo interface {
	Create(context.Context, domain.Box) (*domain.Box, error)
	GetByFingerprint(context.Context, string) (*domain.Box, error)
	GetByContainerID(context.Context, string) (*domain.Box, error)
	Touch(context.Context, string) (*domain.Box, error)
	UpdateStatus(context.Context, string, string) (*domain.Box, error)
}
