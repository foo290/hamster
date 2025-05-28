# Hamster

A Go project with a clean architecture structure.

## Project Structure

```
hamster/
├── cmd/
│   └── server/
│       └── main.go       # Application entry point
├── internal/
│   └── deploy/
│       └── handler.go    # Internal deployment handlers
├── .env                  # Environment variables (optional)
├── go.mod               # Go module definition
└── README.md            # This file
```

## Getting Started

1. Ensure you have Go installed (1.16+)
2. Clone this repository
3. Run `go mod tidy` to install dependencies
4. Run `go run cmd/server/main.go` to start the server

## Development

This project follows standard Go project layout conventions. 