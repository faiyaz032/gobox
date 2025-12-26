package database

import (
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // PostgreSQL driver
)

func GetDB() (*sqlx.DB, error) {
	dbUser := os.Getenv("POSTGRES_USER")
	dbPassword := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB")
	dbHost := os.Getenv("POSTGRES_HOST")

	connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s sslmode=disable", dbUser, dbPassword, dbName, dbHost)

	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		log.Printf("❌ Database connection failed: %v", err)
		return nil, fmt.Errorf("database connection failed: %w", err)
	}

	// Test connection
	if err := db.Ping(); err != nil {
		log.Printf("❌ Database ping failed: %v", err)
		db.Close()
		return nil, fmt.Errorf("database ping failed: %w", err)
	}

	log.Println("✅ Database connected")
	return db, nil
}
