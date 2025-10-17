package database

import (
	"context"
	"sync"
	"time"
)

type SessionContainer struct {
	ContainerID string    `json:"container_id"`
	LastActive  time.Time `json:"last_active"`
}

type Repository struct {
	store map[string]SessionContainer
	mu    sync.RWMutex
}

func NewRepository() *Repository {
	return &Repository{
		store: make(map[string]SessionContainer),
	}
}

func (s *Repository) Get(ctx context.Context, sessionID string) (SessionContainer, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	val, ok := s.store[sessionID]
	return val, ok
}

func (s *Repository) Set(ctx context.Context, sessionID string, value SessionContainer) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.store[sessionID] = value
	return nil
}

func (s *Repository) Delete(ctx context.Context, sessionID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.store, sessionID)
	return nil
}
