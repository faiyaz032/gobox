package cmd

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"

	"github.com/faiyaz032/gobox/internal/box"
	"github.com/faiyaz032/gobox/internal/config"
	"github.com/faiyaz032/gobox/internal/docker"
	"github.com/faiyaz032/gobox/internal/infra/db/postgres"
	"github.com/faiyaz032/gobox/internal/infra/logger"
	"github.com/faiyaz032/gobox/internal/repo"
	boxhandler "github.com/faiyaz032/gobox/internal/rest/handler/box"
)

func RunServer(cfg *config.Config) {
	// Initialize logger
	log, err := logger.New(cfg.Environment)
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize logger: %v", err))
	}
	defer log.Sync()

	log.Info("Starting GoBox server", zap.String("environment", cfg.Environment))

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
		log.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer db.Close()

	queries := postgres.NewQueries(db)

	dockerSvc, err := docker.NewSvc()
	if err != nil {
		log.Fatal("Failed to initialize docker client", zap.Error(err))
	}
	defer dockerSvc.Close()

	boxRepo := repo.NewBoxRepo(queries)
	boxSvc := box.NewSvc(boxRepo, dockerSvc, log)
	boxHandler := boxhandler.NewHandler(boxSvc, log)

	ctx := context.Background()

	imageName := "gobox-base:latest"
	dockerfilePath := "./base-image"
	networkName := "gobox-c-network"
	subnet := "172.25.0.0/16"

	// Ensure network
	_, err = dockerSvc.EnsureNetwork(ctx, networkName, subnet)
	if err != nil {
		log.Fatal("Failed to ensure network", zap.Error(err))
	}

	if err := dockerSvc.EnsureImage(ctx, imageName, dockerfilePath); err != nil {
		log.Fatal("Failed to ensure docker image", zap.Error(err))
	}

	log.Info("Docker base image ensured")

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)

	// Routes
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		log.Info("Server is healthy")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Server is running 🚀"))
	})

	boxhandler.RegisterRoutes(r, boxHandler)

	addr := fmt.Sprintf(":%s", cfg.Server.Port)

	log.Info("Starting server", zap.String("address", addr))
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatal("Server failed", zap.Error(err))
	}
}
