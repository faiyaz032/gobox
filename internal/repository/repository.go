package repository

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
)

type SessionContainer struct {
	ID          int64     `db:"id" json:"id"`
	SessionID   string    `db:"session_id" json:"session_id"`
	ContainerID string    `db:"container_id" json:"container_id"`
	PausedAt    time.Time `db:"paused_at" json:"paused_at"`
	IsPaused    bool      `db:"is_paused" json:"is_paused"`
}

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, sc *SessionContainer) error {
	const query = `
		INSERT INTO session_containers (session_id, container_id, paused_at, is_paused)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`
	// Force UTC before saving
	pausedAtUTC := sc.PausedAt.UTC()
	err := r.db.QueryRowContext(ctx, query, sc.SessionID, sc.ContainerID, pausedAtUTC, sc.IsPaused).Scan(&sc.ID)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetOne(ctx context.Context, sessionID string) (*SessionContainer, error) {
	const query = `
		SELECT id, session_id, container_id, paused_at, is_paused
		FROM session_containers
		WHERE session_id = $1
	`
	var sc SessionContainer
	if err := r.db.GetContext(ctx, &sc, query, sessionID); err != nil {
		return nil, err
	}
	// Ensure UTC on retrieval
	sc.PausedAt = sc.PausedAt.UTC()
	return &sc, nil
}

func (r *Repository) GetOneByContainerID(ctx context.Context, containerID string) (*SessionContainer, error) {
	const query = `
		SELECT id, session_id, container_id, paused_at, is_paused
		FROM session_containers
		WHERE container_id = $1
	`
	var sc SessionContainer
	if err := r.db.GetContext(ctx, &sc, query, containerID); err != nil {
		return nil, err
	}
	// Ensure UTC on retrieval
	sc.PausedAt = sc.PausedAt.UTC()
	return &sc, nil
}

func (r *Repository) GetAll(ctx context.Context) ([]SessionContainer, error) {
	const query = `SELECT id, session_id, container_id, paused_at, is_paused FROM session_containers`
	var containers []SessionContainer
	if err := r.db.SelectContext(ctx, &containers, query); err != nil {
		return nil, err
	}
	// Ensure UTC for all containers
	for i := range containers {
		containers[i].PausedAt = containers[i].PausedAt.UTC()
	}
	return containers, nil
}

func (r *Repository) Update(ctx context.Context, sc *SessionContainer) error {
	const query = `
		UPDATE session_containers
		SET session_id = $1, container_id = $2, paused_at = $3, is_paused = $4
		WHERE id = $5
	`
	// Force UTC before saving
	pausedAtUTC := sc.PausedAt.UTC()
	result, err := r.db.ExecContext(ctx, query, sc.SessionID, sc.ContainerID, pausedAtUTC, sc.IsPaused, sc.ID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return err
	}
	return nil
}

func (r *Repository) Delete(ctx context.Context, id int64) error {
	const query = `
		DELETE FROM session_containers
		WHERE id = $1
	`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return err
	}
	return nil
}
