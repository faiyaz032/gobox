package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/docker/docker/client"
	"github.com/faiyaz032/gobox/internal/handler"
	"github.com/faiyaz032/gobox/internal/infra/database"
	"github.com/faiyaz032/gobox/internal/service"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {

	apiClient, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return err
	}
	defer apiClient.Close()

	repository := database.NewRepository()
	svc := service.NewService(repository, apiClient)
	h := handler.NewHandler(svc, apiClient)

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		if err := h.HandleWS(w, r); err != nil {
			fmt.Fprintf(os.Stderr, "WebSocket error: %v\n", err)
		}
	})

	fmt.Println("Server running on :8080")
	return http.ListenAndServe(":8080", nil)
}
