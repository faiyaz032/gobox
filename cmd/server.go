package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/docker/docker/client"
	"github.com/faiyaz032/gobox/internal/handler"
	"github.com/faiyaz032/gobox/internal/infra/database"
	"github.com/faiyaz032/gobox/internal/repository"
	"github.com/faiyaz032/gobox/internal/service"
	"github.com/go-co-op/gocron/v2"
)

func Serve() error {
	ctx := context.Background()

	apiClient, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return err
	}
	defer apiClient.Close()

	db, _ := database.GetDB()
	repo := repository.NewRepository(db)

	svc := service.NewService(repo, apiClient)
	h := handler.NewHandler(ctx, svc, apiClient)

	scheduler, err := gocron.NewScheduler()
	if err != nil {
		return err
	}
	defer scheduler.Shutdown()

	_, err = scheduler.NewJob(
		gocron.DurationJob(30*time.Second),
		gocron.NewTask(func() {
			log.Print("Running every 30 seconds")
			svc.CleanContainers(ctx)
		}),
	)
	if err != nil {
		return err
	}

	scheduler.Start()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		if err := h.HandleWS(w, r); err != nil {
			fmt.Printf("WebSocket error: %v\n", err)
		}
	})

	fmt.Println("Server running on :8080")
	return http.ListenAndServe(":8080", nil)
}
