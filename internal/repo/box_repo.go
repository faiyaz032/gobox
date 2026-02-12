package repo

import (
	"context"
	"time"

	"github.com/faiyaz032/gobox/internal/domain"
	db "github.com/faiyaz032/gobox/internal/infra/db/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
)

type BoxRepo struct {
	queries *db.Queries
}

func NewBoxRepo(queries *db.Queries) *BoxRepo {
	return &BoxRepo{
		queries: queries,
	}
}

func (r *BoxRepo) Create(ctx context.Context, box domain.Box) (*domain.Box, error) {
	params := db.CreateBoxParams{
		FingerprintID: box.FingerprintID,
		ContainerID:   box.ContainerID,
		Status:        box.Status,
		LastActive: pgtype.Timestamp{
			Time:  box.LastActive,
			Valid: true,
		},
	}

	dbBox, err := r.queries.CreateBox(ctx, params)
	if err != nil {
		return nil, err
	}

	return r.toDomain(dbBox), nil
}

func (r *BoxRepo) GetByFingerprint(ctx context.Context, fingerprintID string) (*domain.Box, error) {
	dbBox, err := r.queries.GetBoxByFingerprint(ctx, fingerprintID)
	if err != nil {
		return nil, err
	}

	return r.toDomain(dbBox), nil
}

func (r *BoxRepo) GetByContainerID(ctx context.Context, containerID string) (*domain.Box, error) {
	dbBox, err := r.queries.GetBoxByContainerID(ctx, containerID)
	if err != nil {
		return nil, err
	}

	return r.toDomain(dbBox), nil
}

func (r *BoxRepo) Touch(ctx context.Context, fingerprintID string) (*domain.Box, error) {
	params := db.TouchBoxParams{
		FingerprintID: fingerprintID,
		LastActive: pgtype.Timestamp{
			Time:  time.Now(),
			Valid: true,
		},
	}

	err := r.queries.TouchBox(ctx, params)
	if err != nil {
		return nil, err
	}

	return r.GetByFingerprint(ctx, fingerprintID)
}

func (r *BoxRepo) UpdateStatus(ctx context.Context, fingerprintID string, status string) (*domain.Box, error) {
	params := db.UpdateBoxStatusParams{
		FingerprintID: fingerprintID,
		Status:        status,
	}

	err := r.queries.UpdateBoxStatus(ctx, params)
	if err != nil {
		return nil, err
	}

	return r.GetByFingerprint(ctx, fingerprintID)
}

func (r *BoxRepo) toDomain(dbBox db.Box) *domain.Box {
	var lastActive time.Time
	if dbBox.LastActive.Valid {
		lastActive = dbBox.LastActive.Time
	}

	status := domain.BoxStatus(dbBox.Status.(string))

	return &domain.Box{
		ID:            dbBox.ID,
		FingerprintID: dbBox.FingerprintID,
		ContainerID:   dbBox.ContainerID,
		Status:        status,
		LastActive:    lastActive,
	}
}
