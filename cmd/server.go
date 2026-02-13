package cmd

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/faiyaz032/gobox/internal/box"
	"github.com/faiyaz032/gobox/internal/config"
	"github.com/faiyaz032/gobox/internal/docker"
	"github.com/faiyaz032/gobox/internal/infra/db/postgres"
	"github.com/faiyaz032/gobox/internal/repo"
	boxhandler "github.com/faiyaz032/gobox/internal/rest/handler/box"
)

func RunServer(cfg *config.Config) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.DBName,
		cfg.Database.SSLMode,
	)

	db, err := postgres.Connect(dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	queries := postgres.NewQueries(db)

	boxRepo := repo.NewBoxRepo(queries)
	boxSvc := box.NewSvc(boxRepo)
	boxHandler := boxhandler.NewHandler(boxSvc)

	dockerSvc, err := docker.NewSvc()
	if err != nil {
		log.Fatalf("Failed to initialize docker client: %v", err)
	}
	defer dockerSvc.Close()

	ctx := context.Background()

	imageName := "gobox-base:latest"
	dockerfilePath := "./base-image"
	networkName := "gobox-c-network"
	subnet := "172.25.0.0/16"

	// Ensure network
	_, err = dockerSvc.EnsureNetwork(ctx, networkName, subnet)
	if err != nil {
		log.Fatalf("Failed to ensure network: %v", err)
	}

	if err := dockerSvc.EnsureImage(ctx, imageName, dockerfilePath); err != nil {
		log.Fatalf("Failed to ensure docker image: %v", err)
	}

	log.Println("Docker base image ensured ✅")

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)

	// Routes
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Server is running 🚀"))
	})

	boxhandler.RegisterRoutes(r, boxHandler)

	addr := fmt.Sprintf(":%s", cfg.Server.Port)

	log.Printf("Starting server on %s\n", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
