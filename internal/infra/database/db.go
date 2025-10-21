package database

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
)

func GetDB() (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", "user=gobox password=gobox123 sslmode=disable")
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
