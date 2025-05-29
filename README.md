# Hamster

A Go service that provides a secure endpoint to trigger Docker Compose deployments. Perfect for simple deployment automation of Docker Compose projects.

## Project Structure

```
hamster/
├── cmd/
│   └── server/
│       └── main.go       # Application entry point
├── internal/
│   └── deploy/
│       └── handler.go    # Deployment handler
├── .env                  # Environment variables
├── go.mod               # Go module definition
└── README.md            # This file
```

## Getting Started

1. Ensure you have Go installed (1.16+)
2. Clone this repository
3. Run `go mod tidy` to install dependencies
4. Set up your `.env` file (see Environment Variables section)
5. Run `go run cmd/server/main.go` to start the server

## Environment Variables

Create a `.env` file in the root directory with these variables:

```
PORT=8080                        # Server port (optional, defaults to 8080)
DEPLOY_TOKEN=your-token          # Authentication token for deployment
COMPOSE_PROJECT_DIR=/path/to/dir # Directory containing your docker-compose.yml
```

## API Endpoints

### POST /deploy
Triggers a Docker Compose deployment by running:
1. `docker compose pull` to get latest images
2. `docker compose up -d --build` to update services

**Headers Required:**
- `Authorization: Bearer <DEPLOY_TOKEN>`

**Response:**
- 200: Successful deployment
- 403: Invalid authentication token
- 500: Deployment error or missing configuration

## Requirements

- Go 1.16+
- Docker and Docker Compose installed on the host machine
- Access to Docker socket
