# GoBox 📦

**GoBox** is an isolated, browser-based Linux playground that provides users with disposable, high-performance containers on demand. Think of it as your own personal Linux sandbox that's ready in seconds and fully destructible.

## 🎯 Purpose
The core objective of this project is to explore and master the management of **distributed systems at scale**. GoBox serves as a "hands-on" experiment in handling hundreds of isolated containers, resource orchestration, and secure networking. 

This project is a stepping stone toward mastering cloud-native infrastructure, with a specific focus on preparing for complex **Kubernetes (K8s) environments** by understanding the low-level container mechanics that power them.

---

## 🏗️ Architecture Overview
GoBox is built with a high-performance backend and a modern, responsive frontend, unified in a monorepo structure.

### 1. Backend (Go)
The Go backend acts as the brain of the operation, utilizing the **Docker SDK** to orchestrate container lifecycles. 
*   **Orchestration**: Dynamically handles container creation, startup, and automatic cleanup.
*   **WebSocket Proxy**: Manages bi-directional terminal I/O between the user's browser and the container's PTY.
*   **Persistence**: Uses **PostgreSQL** to track container states, user fingerprints, and session metadata.

### 2. Frontend (React)
A sleek, **Ubuntu-inspired** UI that provides a premium terminal experience.
*   **xterm.js**: A full-featured terminal emulator running in the browser.
*   **Hyper.js Aesthetic**: A refined, centered, and bordered terminal window for a professional feel.
*   **Responsive Design**: Fully optimized for both desktop and mobile exploration.

---

## 🌐 Networking & Subnetting
One of the key technical highlights of GoBox is how it handles container networking at scale.

### ⚡ The /16 Subnet Choice
GoBox uses a custom Docker bridge network (`gobox-c-network`) configured with a **`172.25.0.0/16`** subnet.
*   **Why /16?**: This mask provides a massive pool of **65,534 usable IP addresses**.
*   **Scale**: This ensures the system can support a vast number of concurrent, isolated containers without ever running into IP exhaustion or collisions.

### 🤖 Dynamic IP Assignment
The system leverages Docker's internal **IPAM (IP Address Management)**. 
*   When a user requests a box, the backend creates a container on the custom bridge. 
*   Docker automatically assigns a unique IP from the `/16` pool to that container.
*   This approach abstracts away the networking complexity, allowing the server to focus strictly on orchestration.

---

## 🚀 Key Features
- **Instant Boot**: Linux containers ready in <3 seconds.
- **Full Root Access**: Install packages (`apt`), manage services, and hack freely.
- **Deep Isolation**: Every user gets a physically isolated environment.
- **Persistence**: Fingerprint-based session recovery — come back to where you left off.
- **Disposable**: One-click "Destroy Box" to wipe everything and start fresh.

## 🛠️ Tech Stack
- **Languages**: Go (Backend), JavaScript (Frontend)
- **Engine**: Docker SDK
- **Database**: PostgreSQL (SQLC + Goose migrations)
- **Frontend**: React, xterm.js, Lucide Icons, CSS3
- **DevOps**: Docker Compose, Makefile, Multi-stage Dockerfiles

---

## 🔧 Local Development
1. Clone the repository.
2. Ensure Docker and Make are installed.
3. Run `make dev` to spin up the local environment.
4. Access the app at `http://localhost:8010`.
