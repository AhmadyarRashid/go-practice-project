# Docker Complete Guide for Beginners

This guide will teach you everything about Docker from zero to deploying your Go application.

## Table of Contents

1. [What is Docker?](#what-is-docker)
2. [Why Use Docker?](#why-use-docker)
3. [Installing Docker](#installing-docker)
4. [Docker Concepts](#docker-concepts)
5. [Basic Docker Commands](#basic-docker-commands)
6. [Understanding Dockerfile](#understanding-dockerfile)
7. [Understanding docker-compose.yml](#understanding-docker-composeyml)
8. [Running This Project with Docker](#running-this-project-with-docker)
9. [Common Docker Operations](#common-docker-operations)
10. [Troubleshooting](#troubleshooting)
11. [Best Practices](#best-practices)

---

## What is Docker?

### Simple Explanation

Imagine you're moving to a new house. Instead of packing items individually and hoping they work in the new place, you put everything in a **shipping container**. This container has everything needed - furniture, appliances, and even the electricity setup.

**Docker does the same thing for software!**

```
Traditional Way:
┌─────────────────────────────────────────────────────┐
│  Your Computer                                       │
│  ┌─────────┐ ┌─────────┐ ┌─────────┐               │
│  │ Go 1.21 │ │ Node 18 │ │ Python 3│               │
│  └─────────┘ └─────────┘ └─────────┘               │
│  Different versions, conflicts, "works on my machine"│
└─────────────────────────────────────────────────────┘

Docker Way:
┌─────────────────────────────────────────────────────┐
│  Your Computer (Host)                                │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐   │
│  │ Container 1 │ │ Container 2 │ │ Container 3 │   │
│  │ ┌─────────┐ │ │ ┌─────────┐ │ │ ┌─────────┐ │   │
│  │ │ Go 1.21 │ │ │ │ Node 18 │ │ │ │ Python 3│ │   │
│  │ │ Your App│ │ │ │ Frontend│ │ │ │ ML Model│ │   │
│  │ └─────────┘ │ │ └─────────┘ │ │ └─────────┘ │   │
│  └─────────────┘ └─────────────┘ └─────────────┘   │
│  Each container is isolated and has everything it needs│
└─────────────────────────────────────────────────────┘
```

### Key Terms

| Term | What It Means | Real-World Analogy |
|------|---------------|-------------------|
| **Image** | A blueprint/template for creating containers | A recipe for a cake |
| **Container** | A running instance of an image | The actual cake made from the recipe |
| **Dockerfile** | Instructions to build an image | The recipe written down |
| **Docker Hub** | Online storage for images | A cookbook library |
| **Volume** | Persistent storage for containers | External hard drive |
| **Network** | Communication between containers | Phone lines between houses |

---

## Why Use Docker?

### Problem 1: "It Works on My Machine"

```
Developer's Machine:        Production Server:
├── Go 1.21                 ├── Go 1.19 (different!)
├── PostgreSQL 16           ├── PostgreSQL 14 (different!)
├── Linux Ubuntu            ├── CentOS (different!)
└── Specific libraries      └── Missing libraries

Result: App crashes in production! 😱
```

### Solution with Docker:

```
Same Docker Container runs EVERYWHERE:
├── Go 1.21 ✓
├── PostgreSQL 16 ✓
├── All dependencies ✓
└── Same configuration ✓

Works on: Mac, Windows, Linux, Cloud, Anywhere! 🎉
```

### Benefits Summary

| Benefit | Description |
|---------|-------------|
| **Consistency** | Same environment everywhere |
| **Isolation** | Apps don't interfere with each other |
| **Portability** | Run anywhere Docker is installed |
| **Scalability** | Easily run multiple copies |
| **Version Control** | Track changes to environment |
| **Fast Setup** | New team member? `docker-compose up`! |

---

## Installing Docker

### For macOS

1. **Download Docker Desktop**
   - Go to: https://www.docker.com/products/docker-desktop
   - Click "Download for Mac"
   - Choose Intel or Apple Silicon based on your Mac

2. **Install**
   - Open the downloaded `.dmg` file
   - Drag Docker to Applications folder
   - Open Docker from Applications
   - Wait for Docker to start (whale icon in menu bar)

3. **Verify Installation**
   ```bash
   docker --version
   # Output: Docker version 24.x.x, build xxxxx

   docker-compose --version
   # Output: Docker Compose version v2.x.x
   ```

### For Windows

1. **Requirements**
   - Windows 10/11 64-bit
   - Enable WSL 2 (Windows Subsystem for Linux)

2. **Download Docker Desktop**
   - Go to: https://www.docker.com/products/docker-desktop
   - Click "Download for Windows"

3. **Install**
   - Run the installer
   - Enable WSL 2 when prompted
   - Restart computer if required

4. **Verify**
   ```powershell
   docker --version
   docker-compose --version
   ```

### For Linux (Ubuntu/Debian)

```bash
# Update package index
sudo apt-get update

# Install prerequisites
sudo apt-get install ca-certificates curl gnupg

# Add Docker's official GPG key
sudo install -m 0755 -d /etc/apt/keyrings
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg
sudo chmod a+r /etc/apt/keyrings/docker.gpg

# Add repository
echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu \
  $(. /etc/os-release && echo "$VERSION_CODENAME") stable" | \
  sudo tee /etc/apt/sources.list.d/docker.list > /dev/null

# Install Docker
sudo apt-get update
sudo apt-get install docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin

# Add your user to docker group (so you don't need sudo)
sudo usermod -aG docker $USER

# Log out and back in, then verify
docker --version
```

---

## Docker Concepts

### Images vs Containers

```
┌─────────────────────────────────────────────────────────────┐
│                         IMAGE                                │
│  (Like a Class in programming or a Blueprint)               │
│                                                              │
│  ┌─────────────────────────────────────────────────────┐    │
│  │  - Operating System (Alpine Linux)                   │    │
│  │  - Go Runtime                                        │    │
│  │  - Your Application Code                             │    │
│  │  - Dependencies                                      │    │
│  │  - Configuration                                     │    │
│  └─────────────────────────────────────────────────────┘    │
│                                                              │
│  You can create MULTIPLE containers from ONE image          │
└─────────────────────────────────────────────────────────────┘
           │               │               │
           ▼               ▼               ▼
    ┌───────────┐   ┌───────────┐   ┌───────────┐
    │ Container │   │ Container │   │ Container │
    │     1     │   │     2     │   │     3     │
    │ (Running) │   │ (Running) │   │ (Stopped) │
    └───────────┘   └───────────┘   └───────────┘
```

### Layers Concept

Docker images are built in **layers**. Each instruction in Dockerfile creates a layer:

```
┌─────────────────────────────────────┐
│ Layer 5: Your App Binary            │ ← Smallest, changes often
├─────────────────────────────────────┤
│ Layer 4: Copy source code           │
├─────────────────────────────────────┤
│ Layer 3: Install dependencies       │
├─────────────────────────────────────┤
│ Layer 2: Install Go                 │
├─────────────────────────────────────┤
│ Layer 1: Alpine Linux OS            │ ← Largest, rarely changes
└─────────────────────────────────────┘

Benefits:
- Layers are cached
- If Layer 1-3 don't change, Docker reuses them
- Only rebuilds changed layers = FAST builds!
```

---

## Basic Docker Commands

### Most Common Commands

```bash
# ═══════════════════════════════════════════════════════════
# IMAGE COMMANDS
# ═══════════════════════════════════════════════════════════

# List all images on your machine
docker images

# Download an image from Docker Hub
docker pull nginx

# Build an image from Dockerfile
docker build -t my-app:v1 .

# Remove an image
docker rmi my-app:v1

# ═══════════════════════════════════════════════════════════
# CONTAINER COMMANDS
# ═══════════════════════════════════════════════════════════

# List running containers
docker ps

# List ALL containers (including stopped)
docker ps -a

# Run a container from an image
docker run nginx

# Run in background (detached mode)
docker run -d nginx

# Run with port mapping (host:container)
docker run -d -p 8080:80 nginx

# Run with a name
docker run -d --name my-nginx -p 8080:80 nginx

# Stop a container
docker stop my-nginx

# Start a stopped container
docker start my-nginx

# Remove a container
docker rm my-nginx

# Remove a running container (force)
docker rm -f my-nginx

# ═══════════════════════════════════════════════════════════
# DEBUGGING COMMANDS
# ═══════════════════════════════════════════════════════════

# View container logs
docker logs my-nginx

# Follow logs in real-time
docker logs -f my-nginx

# Execute command inside container
docker exec -it my-nginx bash

# View container details
docker inspect my-nginx

# View resource usage
docker stats

# ═══════════════════════════════════════════════════════════
# CLEANUP COMMANDS
# ═══════════════════════════════════════════════════════════

# Remove all stopped containers
docker container prune

# Remove all unused images
docker image prune

# Remove everything unused (careful!)
docker system prune -a
```

### Command Breakdown Example

```bash
docker run -d -p 8080:80 --name webserver -e MY_VAR=hello nginx:latest
```

| Part | Meaning |
|------|---------|
| `docker run` | Create and start a container |
| `-d` | Run in background (detached) |
| `-p 8080:80` | Map port 8080 on host to port 80 in container |
| `--name webserver` | Give container a name |
| `-e MY_VAR=hello` | Set environment variable |
| `nginx:latest` | Image name and tag |

---

## Understanding Dockerfile

Let's understand our project's Dockerfile line by line:

### Our Dockerfile Explained

```dockerfile
# ═══════════════════════════════════════════════════════════
# STAGE 1: BUILD STAGE
# ═══════════════════════════════════════════════════════════

# Start from Go image based on Alpine Linux
# "AS builder" names this stage so we can reference it later
FROM golang:1.21-alpine AS builder

# Why Alpine? It's tiny! ~5MB vs ~1GB for full OS
# golang:1.21        = ~1GB
# golang:1.21-alpine = ~300MB

# Install build tools needed for CGO (SQLite requires it)
RUN apk add --no-cache gcc musl-dev

# Set the working directory inside the container
# All following commands will run from /app
WORKDIR /app

# Copy go.mod and go.sum first (for caching)
# If these don't change, Docker uses cached dependencies
COPY go.mod go.sum ./

# Download all dependencies
# This layer is cached if go.mod/go.sum don't change
RUN go mod download

# Now copy the rest of the source code
COPY . .

# Build the Go application
# CGO_ENABLED=1: Enable CGO for SQLite
# GOOS=linux: Build for Linux
# -ldflags="-s -w": Strip debug info (smaller binary)
# -o /app/api: Output path
# ./cmd/api: Source path
RUN CGO_ENABLED=1 GOOS=linux go build -ldflags="-s -w" -o /app/api ./cmd/api

# At this point, we have a compiled binary at /app/api
# But the image is still large (~300MB) because it includes
# the entire Go toolchain. We don't need that to RUN the app!

# ═══════════════════════════════════════════════════════════
# STAGE 2: PRODUCTION STAGE
# ═══════════════════════════════════════════════════════════

# Start fresh from a minimal Alpine image
# This is called "multi-stage build"
FROM alpine:3.19

# Install only what we need to RUN (not build)
# ca-certificates: For HTTPS connections
# tzdata: For timezone support
RUN apk add --no-cache ca-certificates tzdata

# Create a non-root user for security
# Running as root is a security risk!
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# Set working directory
WORKDIR /app

# Copy ONLY the binary from the builder stage
# Not the source code, not Go, not dependencies
# Just the compiled binary!
COPY --from=builder /app/api .

# Copy config file template
COPY --from=builder /app/.env.example .env.example

# Change ownership to non-root user
RUN chown -R appuser:appgroup /app

# Switch to non-root user
USER appuser

# Tell Docker this container listens on port 8080
# This is documentation; doesn't actually publish the port
EXPOSE 8080

# Health check - Docker will check if app is healthy
# --interval=30s: Check every 30 seconds
# --timeout=3s: Fail if no response in 3 seconds
# --start-period=5s: Wait 5 seconds before first check
# --retries=3: Mark unhealthy after 3 failures
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/api/v1/health || exit 1

# Command to run when container starts
CMD ["./api"]
```

### Multi-Stage Build Benefit

```
Without Multi-Stage:
┌────────────────────────────┐
│ Final Image: ~500MB        │
│ ├── Go compiler            │
│ ├── All build tools        │
│ ├── Source code            │
│ └── Your binary (10MB)     │
└────────────────────────────┘

With Multi-Stage:
┌────────────────────────────┐
│ Final Image: ~30MB         │
│ ├── Alpine Linux           │
│ └── Your binary (10MB)     │
└────────────────────────────┘

Savings: 94% smaller image! 🎉
```

---

## Understanding docker-compose.yml

Docker Compose lets you define and run multi-container applications.

### Our docker-compose.yml Explained

```yaml
# Docker Compose file format version
version: '3.8'

# Define all services (containers) we need
services:

  # ═══════════════════════════════════════════════════════
  # SERVICE 1: Our Go API
  # ═══════════════════════════════════════════════════════
  api:
    # Build from Dockerfile in current directory
    build:
      context: .              # Where to look for Dockerfile
      dockerfile: Dockerfile  # Which Dockerfile to use

    # Port mapping: HOST:CONTAINER
    # Access http://localhost:8080 → container's port 8080
    ports:
      - "8080:8080"

    # Environment variables for the container
    environment:
      - APP_ENV=production
      - APP_PORT=8080
      - DB_DRIVER=postgres
      - DB_HOST=postgres      # Use service name as hostname!
      - DB_PORT=5432
      - DB_NAME=enterprise_db
      - DB_USER=postgres
      - DB_PASSWORD=postgres
      - DB_SSL_MODE=disable
      - JWT_SECRET=${JWT_SECRET:-your-super-secret-key-change-in-production-minimum-32-chars}

    # Don't start until postgres is healthy
    depends_on:
      postgres:
        condition: service_healthy

    # Restart policy
    restart: unless-stopped

    # Connect to this network
    networks:
      - app-network

  # ═══════════════════════════════════════════════════════
  # SERVICE 2: PostgreSQL Database
  # ═══════════════════════════════════════════════════════
  postgres:
    # Use official PostgreSQL image
    image: postgres:16-alpine

    # Database configuration
    environment:
      - POSTGRES_USER=postgres
      - POSTGRES_PASSWORD=postgres
      - POSTGRES_DB=enterprise_db

    # Persist data even if container is removed
    volumes:
      - postgres_data:/var/lib/postgresql/data

    # Expose port for external tools (optional)
    ports:
      - "5432:5432"

    # Health check for database
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres"]
      interval: 5s
      timeout: 5s
      retries: 5

    restart: unless-stopped
    networks:
      - app-network

  # ═══════════════════════════════════════════════════════
  # SERVICE 3: Redis Cache (Optional, for future use)
  # ═══════════════════════════════════════════════════════
  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data
    command: redis-server --appendonly yes
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 5s
      timeout: 5s
      retries: 5
    restart: unless-stopped
    networks:
      - app-network

# ═══════════════════════════════════════════════════════
# NETWORKS
# ═══════════════════════════════════════════════════════
networks:
  app-network:
    driver: bridge
    # All services on this network can communicate
    # using their service names as hostnames

# ═══════════════════════════════════════════════════════
# VOLUMES (Persistent Storage)
# ═══════════════════════════════════════════════════════
volumes:
  postgres_data:    # Database files survive container restart
  redis_data:       # Redis data survives container restart
```

### Network Visualization

```
┌─────────────────────────────────────────────────────────────┐
│                    app-network (bridge)                      │
│                                                              │
│  ┌─────────────┐    ┌─────────────┐    ┌─────────────┐     │
│  │     api     │    │  postgres   │    │    redis    │     │
│  │ :8080       │───▶│ :5432       │    │ :6379       │     │
│  │             │    │             │    │             │     │
│  └─────────────┘    └─────────────┘    └─────────────┘     │
│                                                              │
│  Containers can reach each other by service name:           │
│  - api can connect to postgres:5432                         │
│  - api can connect to redis:6379                            │
└─────────────────────────────────────────────────────────────┘
         │                   │                  │
         ▼                   ▼                  ▼
    localhost:8080     localhost:5432     localhost:6379
    (Your Browser)     (DB Tools)         (Redis Tools)
```

---

## Running This Project with Docker

### Method 1: Using Docker Compose (Recommended)

```bash
# Navigate to project directory
cd /Users/ahmadyar/Desktop/go-enterprise-api

# Start all services (build if needed)
docker-compose up -d

# View logs
docker-compose logs -f

# View logs for specific service
docker-compose logs -f api

# Stop all services
docker-compose down

# Stop and remove volumes (database data!)
docker-compose down -v
```

### Method 2: Manual Docker Commands

```bash
# 1. Create a network
docker network create go-api-network

# 2. Start PostgreSQL
docker run -d \
  --name postgres \
  --network go-api-network \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=enterprise_db \
  -p 5432:5432 \
  postgres:16-alpine

# 3. Wait for PostgreSQL to be ready
sleep 10

# 4. Build your application
docker build -t go-enterprise-api:latest .

# 5. Run your application
docker run -d \
  --name api \
  --network go-api-network \
  -p 8080:8080 \
  -e APP_ENV=production \
  -e DB_DRIVER=postgres \
  -e DB_HOST=postgres \
  -e DB_PORT=5432 \
  -e DB_NAME=enterprise_db \
  -e DB_USER=postgres \
  -e DB_PASSWORD=postgres \
  -e JWT_SECRET=your-super-secret-key-change-in-production-minimum-32-chars \
  go-enterprise-api:latest

# 6. Check if running
docker ps

# 7. Test the API
curl http://localhost:8080/api/v1/health
```

### Step-by-Step First Run

```bash
# 1. Navigate to project
cd /Users/ahmadyar/Desktop/go-enterprise-api

# 2. Build and start (first time will take a few minutes)
docker-compose up -d --build

# 3. Check status
docker-compose ps

# Expected output:
# NAME                  STATUS          PORTS
# go-enterprise-api-api-1       Up (healthy)   0.0.0.0:8080->8080/tcp
# go-enterprise-api-postgres-1  Up (healthy)   0.0.0.0:5432->5432/tcp
# go-enterprise-api-redis-1     Up (healthy)   0.0.0.0:6379->6379/tcp

# 4. Test the health endpoint
curl http://localhost:8080/api/v1/health

# Expected output:
# {"success":true,"data":{"status":"healthy"}}

# 5. View API logs
docker-compose logs -f api

# 6. When done, stop everything
docker-compose down
```

---

## Common Docker Operations

### Viewing Logs

```bash
# All services
docker-compose logs

# Follow logs (real-time)
docker-compose logs -f

# Specific service
docker-compose logs -f api

# Last 100 lines
docker-compose logs --tail=100 api

# With timestamps
docker-compose logs -t api
```

### Entering a Container

```bash
# Get shell access to running container
docker-compose exec api sh

# For containers with bash
docker-compose exec postgres bash

# Run a single command
docker-compose exec postgres psql -U postgres -d enterprise_db

# Example: Check database tables
docker-compose exec postgres psql -U postgres -d enterprise_db -c "\dt"
```

### Rebuilding After Code Changes

```bash
# Rebuild and restart just the API
docker-compose up -d --build api

# Rebuild everything
docker-compose up -d --build

# Force rebuild (no cache)
docker-compose build --no-cache
docker-compose up -d
```

### Managing Data

```bash
# View volumes
docker volume ls

# Inspect a volume
docker volume inspect go-enterprise-api_postgres_data

# Backup database
docker-compose exec postgres pg_dump -U postgres enterprise_db > backup.sql

# Restore database
cat backup.sql | docker-compose exec -T postgres psql -U postgres enterprise_db
```

---

## Troubleshooting

### Common Issues and Solutions

#### Issue 1: Port Already in Use

```bash
# Error: Bind for 0.0.0.0:8080 failed: port is already allocated

# Find what's using the port
lsof -i :8080

# Kill the process or change the port in docker-compose.yml
ports:
  - "8081:8080"  # Use 8081 instead
```

#### Issue 2: Container Keeps Restarting

```bash
# Check logs for errors
docker-compose logs api

# Check container status
docker-compose ps

# Common causes:
# - Database not ready (check depends_on)
# - Wrong environment variables
# - Application crash
```

#### Issue 3: Cannot Connect to Database

```bash
# Check if postgres is running
docker-compose ps postgres

# Check postgres logs
docker-compose logs postgres

# Test connection from api container
docker-compose exec api sh
# Inside container:
nc -zv postgres 5432
```

#### Issue 4: Changes Not Reflected

```bash
# Rebuild the image
docker-compose up -d --build

# If still not working, remove everything and start fresh
docker-compose down -v
docker-compose up -d --build
```

#### Issue 5: Out of Disk Space

```bash
# Check Docker disk usage
docker system df

# Clean up unused resources
docker system prune -a

# Remove unused volumes (careful - deletes data!)
docker volume prune
```

### Debug Checklist

```bash
# 1. Are all containers running?
docker-compose ps

# 2. Check logs for errors
docker-compose logs

# 3. Can containers reach each other?
docker-compose exec api ping postgres

# 4. Is the port accessible from host?
curl http://localhost:8080/api/v1/health

# 5. Check resource usage
docker stats
```

---

## Best Practices

### 1. Use .dockerignore

Create a `.dockerignore` file to exclude unnecessary files:

```
# .dockerignore
.git
.gitignore
*.md
.env
.env.local
tmp/
build/
*.test
coverage.*
.idea/
.vscode/
```

### 2. Never Store Secrets in Images

```yaml
# BAD - Secret in docker-compose.yml
environment:
  - JWT_SECRET=my-actual-secret-key

# GOOD - Use environment variable
environment:
  - JWT_SECRET=${JWT_SECRET}

# Then set it when running:
JWT_SECRET=my-secret docker-compose up -d
```

### 3. Use Specific Image Tags

```dockerfile
# BAD - Could change unexpectedly
FROM golang:latest

# GOOD - Predictable
FROM golang:1.21-alpine
```

### 4. Run as Non-Root User

```dockerfile
# Create and use non-root user
RUN adduser -D appuser
USER appuser
```

### 5. Use Health Checks

```dockerfile
HEALTHCHECK --interval=30s --timeout=3s \
    CMD wget --spider http://localhost:8080/health || exit 1
```

### 6. Minimize Image Size

```dockerfile
# Use multi-stage builds
# Use Alpine-based images
# Remove unnecessary files
# Combine RUN commands
RUN apk add --no-cache package1 package2 && \
    rm -rf /var/cache/apk/*
```

---

## Quick Reference Card

```bash
# ═══════════════════════════════════════════════════════════
# DOCKER COMPOSE COMMANDS
# ═══════════════════════════════════════════════════════════

docker-compose up -d          # Start all services
docker-compose down           # Stop all services
docker-compose ps             # List running services
docker-compose logs -f        # Follow all logs
docker-compose logs -f api    # Follow specific service logs
docker-compose exec api sh    # Shell into container
docker-compose build          # Rebuild images
docker-compose restart api    # Restart a service

# ═══════════════════════════════════════════════════════════
# DOCKER COMMANDS
# ═══════════════════════════════════════════════════════════

docker ps                     # List running containers
docker ps -a                  # List all containers
docker images                 # List images
docker logs <container>       # View logs
docker exec -it <container> sh # Shell into container
docker stop <container>       # Stop container
docker rm <container>         # Remove container
docker rmi <image>            # Remove image
docker system prune -a        # Clean up everything

# ═══════════════════════════════════════════════════════════
# THIS PROJECT
# ═══════════════════════════════════════════════════════════

# Start the project
cd /Users/ahmadyar/Desktop/go-enterprise-api
docker-compose up -d --build

# Test it works
curl http://localhost:8080/api/v1/health

# Stop the project
docker-compose down
```

---

## Next Steps

1. **Practice**: Run the project with Docker multiple times
2. **Experiment**: Modify docker-compose.yml and see what happens
3. **Learn More**:
   - Docker documentation: https://docs.docker.com/
   - Docker Compose docs: https://docs.docker.com/compose/
   - Play with Docker: https://labs.play-with-docker.com/

**Congratulations!** You now understand Docker well enough to deploy this Go application! 🎉
