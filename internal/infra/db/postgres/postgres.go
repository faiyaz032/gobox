package postgres

import (
	"context"
	"fmt"
	"log"
	"time"

	db "github.com/faiyaz032/gobox/internal/infra/db/sqlc"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Connect connects to PostgreSQL and returns a pgxpool.Pool
func Connect(dsn string) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to create pgx pool: %w", err)
	}

	// Ping to verify connection
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Successfully connected to PostgreSQL database")
	return pool, nil
}

// NewQueries creates a new Queries instance
func NewQueries(pool *pgxpool.Pool) *db.Queries {
	return db.New(pool)
}

// Close closes the pool
func Close(pool *pgxpool.Pool) {
	pool.Close()
}
