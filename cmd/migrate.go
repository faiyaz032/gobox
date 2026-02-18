package cmd

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"

	"github.com/faiyaz032/gobox/internal/config"
	"github.com/faiyaz032/gobox/internal/infra/logger"
)

func RunMigrate(cfg *config.Config) error {
	// Initialize logger
	log, err := logger.New(cfg.Environment)
	if err != nil {
		return fmt.Errorf("failed to initialize logger: %w", err)
	}
	defer log.Sync()

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.DBName,
		cfg.Database.SSLMode,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("failed to open db: %w", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping db: %w", err)
	}

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("failed to set goose dialect: %w", err)
	}

	migrationsDir := "./migrations"

	log.Info("Running migrations...")
	if err := goose.Up(db, migrationsDir); err != nil {
		return fmt.Errorf("goose up failed: %w", err)
	}

	log.Info("Migrations completed successfully")
	return nil
}
