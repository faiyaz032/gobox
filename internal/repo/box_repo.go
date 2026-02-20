package repo

import (
	"context"
	"errors"
	"time"

	"github.com/faiyaz032/gobox/internal/domain"
	db "github.com/faiyaz032/gobox/internal/infra/db/sqlc"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
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
		Status:        string(box.Status), // Convert BoxStatus to string
		LastActive: pgtype.Timestamp{
			Time:  box.LastActive,
			Valid: true,
		},
	}

	dbBox, err := r.queries.CreateBox(ctx, params)
	if err != nil {
		return nil, r.mapError(err, "create box")
	}

	return r.toDomain(dbBox), nil
}

func (r *BoxRepo) GetByFingerprint(ctx context.Context, fingerprintID string) (*domain.Box, error) {
	dbBox, err := r.queries.GetBoxByFingerprint(ctx, fingerprintID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.NewNotFoundError("box", fingerprintID)
		}
		return nil, r.mapError(err, "get box by fingerprint")
	}

	return r.toDomain(dbBox), nil
}

func (r *BoxRepo) GetByContainerID(ctx context.Context, containerID string) (*domain.Box, error) {
	dbBox, err := r.queries.GetBoxByContainerID(ctx, containerID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.NewNotFoundError("box", containerID)
		}
		return nil, r.mapError(err, "get box by container ID")
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
		return nil, r.mapError(err, "touch box")
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
		return nil, r.mapError(err, "update box status")
	}

	return r.GetByFingerprint(ctx, fingerprintID)
}

func (r *BoxRepo) GetExpiredBoxes(ctx context.Context, lastActive time.Time) ([]domain.Box, error) {
	dbBoxes, err := r.queries.GetExpiredBoxes(ctx, pgtype.Timestamp{
		Time:  lastActive,
		Valid: true,
	})
	if err != nil {
		return nil, r.mapError(err, "get expired boxes")
	}

	boxes := make([]domain.Box, len(dbBoxes))
	for i, dbBox := range dbBoxes {
		boxes[i] = *r.toDomain(dbBox)
	}

	return boxes, nil
}

func (r *BoxRepo) Delete(ctx context.Context, fingerprintID string) error {
	err := r.queries.DeleteBox(ctx, fingerprintID)
	if err != nil {
		return r.mapError(err, "delete box")
	}
	return nil
}

func (r *BoxRepo) toDomain(dbBox db.Box) *domain.Box {
	var lastActive time.Time
	if dbBox.LastActive.Valid {
		lastActive = dbBox.LastActive.Time
	}

	return &domain.Box{
		ID:            dbBox.ID,
		FingerprintID: dbBox.FingerprintID,
		ContainerID:   dbBox.ContainerID,
		Status:        domain.BoxStatus(dbBox.Status),
		LastActive:    lastActive,
	}
}

// converts database errors to AppError
func (r *BoxRepo) mapError(err error, operation string) error {
	if err == nil {
		return nil
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505": // unique_violation
			return domain.NewConflictError("record already exists")
		case "23503": // foreign_key_violation
			return domain.NewValidationError("referenced record does not exist")
		case "23514": // check_violation
			return domain.NewValidationError("constraint violation: " + pgErr.Message)
		default:
			return domain.NewDatabaseError(operation, err)
		}
	}

	// Handle connection errors
	if errors.Is(err, context.DeadlineExceeded) {
		return domain.NewDatabaseError(operation+" (timeout)", err)
	}

	if errors.Is(err, context.Canceled) {
		return domain.NewDatabaseError(operation+" (canceled)", err)
	}

	return domain.NewDatabaseError(operation, err)
}
