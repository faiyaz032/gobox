package domain

import (
	"time"

	"github.com/google/uuid"
)

type BoxStatus string

const (
	StatusRunning BoxStatus = "running"
	StatusPaused  BoxStatus = "paused"
)

type Box struct {
	ID            uuid.UUID `json:"id"`
	FingerprintID string    `json:"fingerprint_id"`
	ContainerID   string    `json:"container_id"`
	Status        BoxStatus `json:"status"`
	LastActive    time.Time `json:"last_active"`
}
