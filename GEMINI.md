# GEMINI.md

## Project Overview

This project is an authentication service written in Go. It uses gRPC for communication and is designed to be run as a microservice.

The main technologies used are:
- **Go**: The programming language used for the application.
- **gRPC**: For the API layer. The protobuf definitions are managed in an external repository (`github.com/Weit145/proto-repo`).
- **cleanenv**: For configuration management, loading settings from a YAML file.
- **slog**: Go's standard library for structured logging.

The project follows a standard Go project layout, with the main application entrypoint in the `cmd` directory and the core logic in the `internal` directory. The configuration is clearly separated, and there are utility packages for logging and standardizing responses.

## Building and Running

### Prerequisites
- Go 1.25.5 or later installed.
- Access to the private protobuf repository `github.com/Weit145/proto-repo`.

### Running the application

To run the application, execute the following command from the root of the project:

```bash
go run cmd/main.go
```

The server will start and listen for gRPC connections on port `:50051` as configured in `config/local.yaml`.

### Configuration

The application is configured via a YAML file. By default, it uses `config/local.yaml`. You can specify a different configuration file by setting the `CONFIG_PATH` environment variable:

```bash
CONFIG_PATH=path/to/your/config.yaml go run cmd/main.go
```

## Development Conventions

### Logging

The project uses Go's structured logging library, `slog`. 
- In the `local` environment, logs are output as text to `stdout`.
- In the `prod` environment, logs are output as JSON to `stdout`.

When logging errors, use the `logger.Err()` helper function to create a structured error field:

```go
import "github.com/Weit145/Auth_golang/internal/lib/logger"

// ...

log.Error("an error occurred", logger.Err(err))
```

### Responses

The `internal/lib/response` package provides a standardized structure for API responses, which suggests that a REST/HTTP gateway might be a planned feature.

### gRPC Implementation

The gRPC services are defined in the `internal/grpc` directory. The actual service logic should be implemented in the `server` struct in `internal/grpc/gateway/server.go`, which currently uses `pb.UnimplementedAuthServer`.
