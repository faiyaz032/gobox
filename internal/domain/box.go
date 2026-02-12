package domain

import (
	"time"

	"github.com/google/uuid"
)

type Box struct {
	ID            uuid.UUID `json:"id"`
	FingerprintID uuid.UUID `json:"fingerprint_id"`
	ContainerID   string    `json:"container_id"`
	Status        string    `json:"status"`
	LastActive    time.Time `json:"last_active"`
}
