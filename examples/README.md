# Goober Usage Examples

This directory contains example configurations for using Goober with different types of Go projects.

## Web Server Example

```bash
# Basic web server
goober -build "go build -o server ./cmd/server" -run "./server"

# With environment variables
goober -build "go build -o server ./cmd/server" -run "PORT=8080 ./server"

# With build tags
goober -build "go build -tags dev -o server ./cmd/server" -run "./server --env dev"
```

## Microservice Example

```bash
# API Gateway
goober -dir ./gateway -build "go build -o gateway ./cmd/gateway" -run "./gateway --port 8080"

# Auth Service
goober -dir ./auth -build "go build -o auth ./cmd/auth" -run "./auth --port 8081"

# User Service
goober -dir ./users -build "go build -o users ./cmd/users" -run "./users --port 8082"
```

## CLI Tool Example

```bash
# CLI application
goober -build "go build -o mytool ./cmd/mytool" -run "./mytool --help"

# With test data
goober -build "go build -o mytool ./cmd/mytool" -run "./mytool process --input testdata/sample.json"
```

## gRPC Service Example

```bash
# gRPC server
goober -build "go build -o grpc-server ./cmd/grpc" -run "./grpc-server --port 9000"

# With TLS
goober -build "go build -o grpc-server ./cmd/grpc" -run "./grpc-server --port 9000 --tls-cert cert.pem --tls-key key.pem"
```

## Multi-module Project Example

```bash
# Root module watcher
goober -dir ./services/api -build "cd services/api && go build -o ../../bin/api" -run "./bin/api"

# Specific service
goober -dir ./internal/userservice -build "go build -o userservice ./cmd/userservice" -run "./userservice"
```

## Docker Development Example

```bash
# Build and run in Docker
goober -build "docker build -t myapp ." -run "docker run --rm -p 8080:8080 myapp"

# Docker Compose
goober -build "docker-compose build app" -run "docker-compose up app"
```

## With Makefile Example

```bash
# Use Makefile commands
goober -build "make build" -run "make run"

# Development target
goober -build "make build-dev" -run "make run-dev"
```

## Testing Example

```bash
# Run tests on change
goober -build "go test ./..." -run "echo 'Tests completed'"

# With coverage
goober -build "go test -cover ./..." -run "go tool cover -html=coverage.out"
```

## Configuration Files

You can also create project-specific scripts:

### `goober.sh` (Project script)

```bash
#!/bin/bash
goober -dir ./src \
       -build "go build -ldflags '-X main.version=$(git describe --tags)' -o bin/myapp ./cmd/myapp" \
       -run "./bin/myapp --config config/dev.yaml"
```

### `package.json` equivalent

```json
{
  "scripts": {
    "dev": "goober -build 'go build -o app' -run './app'",
    "dev:api": "goober -dir ./api -build 'go build -o api ./cmd/api' -run './api --port 8080'",
    "dev:worker": "goober -dir ./worker -build 'go build -o worker ./cmd/worker' -run './worker'"
  }
}
```

## Environment-specific Examples

### Development

```bash
goober -build "go build -tags dev -o app" -run "APP_ENV=dev ./app"
```

### Staging

```bash
goober -build "go build -tags staging -o app" -run "APP_ENV=staging ./app"
```

### With Hot Reload for Templates

```bash
# Watch templates too (if using template files)
goober -build "go build -o webapp ./cmd/webapp" -run "TEMPLATE_DEV=true ./webapp"
```
