package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/docker/docker/client"
	"github.com/faiyaz032/gobox/internal/docker"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func wsHandler(apiClient *client.Client, ctx context.Context) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			fmt.Println("WebSocket upgrade error:", err)
			return
		}
		defer conn.Close()

		containerId, err := docker.CreateContainer(apiClient, ctx)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating container: %v\n", err)
			return
		}
		defer docker.CleanUP(apiClient, ctx, containerId)

		if err := docker.StartContainer(apiClient, ctx, containerId); err != nil {
			fmt.Fprintf(os.Stderr, "Error starting container: %v\n", err)
			return
		}

		hijackResp, err := docker.AttachShell(apiClient, ctx, containerId)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error attaching shell: %v\n", err)
			return
		}

		go func() {
			buf := make([]byte, 1024)
			for {
				n, err := hijackResp.Reader.Read(buf)
				if err != nil {
					break
				}
				conn.WriteMessage(websocket.TextMessage, buf[:n])
			}
		}()

		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				break
			}
			hijackResp.Conn.Write(msg)
		}
	}
}

func main() {

	ctx := context.Background()
	apiClient, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating Docker client: %v\n", err)
		return
	}
	defer apiClient.Close()

	http.HandleFunc("/ws", wsHandler(apiClient, ctx))
	fmt.Println("Server running on :8080")
	http.ListenAndServe(":8080", nil)
}
