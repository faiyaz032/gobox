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
	LastActive  time.Time `db:"last_active" json:"last_active"`
}

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, sc *SessionContainer) error {
	const query = `
		INSERT INTO session_containers (session_id, container_id, last_active)
		VALUES ($1, $2, $3)
		RETURNING id
	`
	err := r.db.QueryRowContext(ctx, query, sc.SessionID, sc.ContainerID, sc.LastActive).Scan(&sc.ID)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetOne(ctx context.Context, sessionID string) (*SessionContainer, error) {
	const query = `
		SELECT id, session_id, container_id, last_active
		FROM session_containers
		WHERE session_id = $1
	`
	var sc SessionContainer
	if err := r.db.GetContext(ctx, &sc, query, sessionID); err != nil {
		return nil, err
	}
	return &sc, nil
}

func (r *Repository) GetAll(ctx context.Context) ([]SessionContainer, error) {
	const query = `SELECT * FROM session_container`
	var containers []SessionContainer
	if err := r.db.SelectContext(ctx, &containers, query); err != nil {
		return nil, err
	}
	return containers, nil
}

func (r *Repository) Update(ctx context.Context, sc *SessionContainer) error {
	const query = `
		UPDATE session_containers
		SET session_id = $1, container_id = $2, last_active = $3
		WHERE id = $4
	`
	result, err := r.db.ExecContext(ctx, query, sc.SessionID, sc.ContainerID, sc.LastActive, sc.ID)
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
