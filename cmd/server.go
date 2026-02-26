package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

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
	defer func() {
		_ = log.Sync()
	}()

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
		_, err := w.Write([]byte("Server is running 🚀"))
		if err != nil {
			log.Error("Failed to write health check response", zap.Error(err))
		}
	})

	boxhandler.RegisterRoutes(r, boxHandler)

	// Serve static files from the frontend build directory
	staticPath := "./frontend/build"
	fileServer := http.FileServer(http.Dir(staticPath))

	r.Handle("/*", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if the requested file exists in the static path
		path := filepath.Join(staticPath, r.URL.Path)
		_, err := os.Stat(path)
		if os.IsNotExist(err) || r.URL.Path == "/" {
			// If not found, serve index.html (SPA routing)
			http.ServeFile(w, r, filepath.Join(staticPath, "index.html"))
			return
		}
		// Otherwise, serve the static file
		fileServer.ServeHTTP(w, r)
	}))

	addr := fmt.Sprintf(":%s", cfg.Server.Port)

	log.Info("Starting server", zap.String("address", addr))
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatal("Server failed", zap.Error(err))
	}
}
