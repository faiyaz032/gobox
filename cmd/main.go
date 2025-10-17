package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/docker/docker/client"
	"github.com/faiyaz032/gobox/internal/docker"
	"github.com/gorilla/websocket"
)

type DatabaseItem struct {
	ContainerID string    `json:"container_id"`
	LastActive  time.Time `json:"last_active"`
}
type Database struct {
	store map[string]DatabaseItem
	mu    sync.RWMutex
}

func (d *Database) Get(ctx context.Context, key string) (DatabaseItem, bool) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	val, ok := d.store[key]
	return val, ok
}

func (d *Database) Set(ctx context.Context, key string, value DatabaseItem) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.store[key] = value
	return nil
}

func (d *Database) Delete(ctx context.Context, key string) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	delete(d.store, key)
	return nil
}

var database = &Database{
	store: make(map[string]DatabaseItem),
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func wsHandler(ctx context.Context, apiClient *client.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionID := r.URL.Query().Get("session_id")
		if sessionID == "" {
			http.Error(w, "session_id missing", http.StatusBadRequest)
			return
		}
		fmt.Println("Session ID:", sessionID)

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			fmt.Println("WebSocket upgrade error:", err)
			return
		}
		defer conn.Close()

		item, ok := database.Get(ctx, sessionID)
		if !ok {

			containerId, err := docker.CreateContainer(apiClient, ctx)
			if err != nil {
				fmt.Println("failed  to create container:", err)
				return
			}
			mapItem := DatabaseItem{
				ContainerID: containerId,
				LastActive:  time.Now(),
			}

			database.Set(ctx, sessionID, mapItem)
			item = mapItem
		} else {
			err := docker.UnpauseContainer(ctx, apiClient, item.ContainerID)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error unpausing container: %v\n", err)
				return
			}
		}

		// Log database state
		fmt.Println("---- DATABASE STATE ----")
		dbJSON, _ := json.MarshalIndent(database.store, "", "  ")
		fmt.Println(string(dbJSON))
		fmt.Println("-------------------------")

		//defer docker.RemoveContainer(apiClient, ctx, containerID)

		if err := docker.StartContainer(apiClient, ctx, item.ContainerID); err != nil {
			fmt.Fprintf(os.Stderr, "Error starting container: %v\n", err)
			return
		}

		hijackResp, err := docker.AttachShell(apiClient, ctx, item.ContainerID)
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
				HandleWsDisconnect(ctx, apiClient, sessionID)
				break
			}
			hijackResp.Conn.Write(msg)
		}
	}
}

func HandleWsDisconnect(ctx context.Context, apiClient *client.Client, sessionID string) {
	item, ok := database.Get(ctx, sessionID)
	if !ok {
		return
	}
	err := docker.PauseContainer(apiClient, ctx, item.ContainerID)
	if err != nil {
		return
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

	http.HandleFunc("/ws", wsHandler(ctx, apiClient))
	fmt.Println("Server running on :8080")
	http.ListenAndServe(":8080", nil)
}
