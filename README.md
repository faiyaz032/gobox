# Gobox

Gobox is a web service that provides an instant Linux playground in your browser. Users can visit the site and immediately get an interactive shell in a secure, isolated Ubuntu container. Behind the scenes, Gobox creates a dedicated, isolated Docker container and network for each user session, ensuring a safe and private environment.

## Features

*   **Instant Linux Playground:** Get an interactive `bash` shell in an Ubuntu container directly in your browser.
*   **On-demand Container Creation:** A fresh container is created for each user session.
*   **Isolated Environments:** Each container runs in its own isolated Docker network for enhanced security.
*   **Interactive Shell:** Provides an interactive `bash` shell within the container, accessible via a WebSocket connection.
*   **Resource Limitation:** Containers are created with limited CPU and memory resources to prevent abuse.
*   **Automatic Cleanup:** A scheduler periodically cleans up and removes inactive containers.

## Architecture

Gobox follows a layered architecture, which separates the application into two main components:

*   **Handler:** The handler layer is responsible for handling incoming HTTP and WebSocket requests. It decodes requests and calls the appropriate service methods.
*   **Service:** The service layer contains the core business logic of the application. It orchestrates the interaction with the Docker client to manage the lifecycle of the containers.

```
+-------------------+      +-------------------+
|      Handler      |----->|      Service      |
+-------------------+      +-------------------+
        |                          |
        |                          |
        v                          v
+-------------------+      +-------------------+
|    WebSocket      |      |   Docker Client   |
+-------------------+      +-------------------+
```

## Technologies Used

*   **Go:** The application is written in Go, a fast and concurrent language that is well-suited for building web services.
*   **Docker:** Docker is used to create and manage the sandboxed container environments.
*   **WebSocket:** WebSockets are used for real-time, bidirectional communication between the client and the container's shell.
*   **gocron:** The `gocron` library is used to schedule the cleanup of old containers.

## How to Run

To run the application, you will need to have Docker and Docker Compose installed.

1.  **Build the base Docker image:**

    ```bash
    docker build -t gobox-base:latest -f Dockerfile.base .
    ```

2.  **Run the application using Docker Compose:**

    ```bash
    docker-compose up --build
    ```

The application will be running on `http://localhost:8080`. You can connect to the WebSocket endpoint at `ws://localhost:8080/ws?session_id=<your_session_id>`.

## How it Works

1.  When a user connects to the WebSocket endpoint with a unique `session_id`, the application creates a new Docker container using the `gobox-base:latest` image.
2.  The container is started with resource limitations and a `sleep` command to keep it running.
3.  An interactive `bash` shell is started within the container, and the input and output streams are attached to the WebSocket connection.
4.  When the user disconnects from the WebSocket, the container is paused.
5.  A background job runs periodically to check for inactive containers. These containers are removed to free up resources.