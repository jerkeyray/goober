# 🚀 Hosting and Using Goober Guide

This guide shows you how to host Goober and use it across different projects.

## 📦 Step 1: Prepare for Hosting

### 1.1 Update Module Name

Replace `yourusername` in the following files with your actual GitHub username:

- `go.mod` - Update module name to `github.com/yourusername/goober`
- `cmd/goober/main.go` - Update import paths
- `tui/tui.go` - Update import paths
- `README.md` - Update installation URLs
- `.goreleaser.yml` - Update homebrew tap

### 1.2 Initialize Git Repository

```bash
cd /path/to/goober
git init
git add .
git commit -m "Initial commit: Goober TUI file watcher"
```

## 🌐 Step 2: Host on GitHub

### 2.1 Create GitHub Repository

1. Go to [github.com/new](https://github.com/new)
2. Repository name: `goober`
3. Make it public
4. Don't initialize with README (we already have one)

### 2.2 Push to GitHub

```bash
git remote add origin https://github.com/yourusername/goober.git
git branch -M main
git push -u origin main
```

### 2.3 Create First Release

```bash
git tag v1.0.0
git push origin v1.0.0
```

This will trigger the GitHub Actions workflow to build cross-platform binaries.

## 📱 Step 3: Using Goober in Other Projects

### Method 1: Global Installation (Recommended)

```bash
# Install globally
go install github.com/yourusername/goober/cmd/goober@latest

# Use anywhere
cd ~/my-go-project
goober
```

### Method 2: Project-specific Installation

Create `tools.go` in your project:

```go
//go:build tools

package tools

import _ "github.com/yourusername/goober/cmd/goober"
```

Add to your project:

```bash
go mod tidy
go install github.com/yourusername/goober/cmd/goober
```

### Method 3: Direct Binary Download

Users can download from GitHub releases:

```bash
# Linux/macOS
curl -L https://github.com/yourusername/goober/releases/latest/download/goober_Linux_x86_64.tar.gz | tar xz
sudo mv goober /usr/local/bin/

# Or use the install script (create this)
curl -sSL https://raw.githubusercontent.com/yourusername/goober/main/install.sh | bash
```

## 🛠️ Step 4: Project Usage Examples

### Basic Web API Project

```bash
cd my-api-project
goober -build "go build -o api ./cmd/api" -run "./api --port 8080"
```

### Microservices Project

```bash
# Terminal 1 - User Service
cd user-service
goober -build "go build -o user-service" -run "./user-service --port 8001"

# Terminal 2 - Order Service
cd order-service
goober -build "go build -o order-service" -run "./order-service --port 8002"
```

### CLI Tool Development

```bash
cd my-cli-tool
goober -build "go build -o mytool ./cmd/mytool" -run "./mytool --help"
```

## 📋 Step 5: Create Project Templates

### 5.1 Create `.goober` config file support (optional enhancement)

```json
{
  "build": "go build -o app ./cmd/app",
  "run": "./app --dev",
  "dir": "./src",
  "debounce": "1s"
}
```

### 5.2 Create project scripts

**`scripts/dev.sh`**:

```bash
#!/bin/bash
goober -build "go build -ldflags '-X main.version=dev' -o bin/app ./cmd/app" \
       -run "./bin/app --config config/dev.yaml"
```

## 🎯 Step 6: Distribution Strategies

### Homebrew (macOS/Linux)

The `.goreleaser.yml` will create a Homebrew formula:

```bash
# After setting up homebrew-tap repository
brew tap yourusername/tap
brew install goober
```

### Docker Distribution

Create `Dockerfile`:

```dockerfile
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o goober ./cmd/goober

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/goober .
CMD ["./goober"]
```

### NPM-like Installation Script

Create `install.sh`:

```bash
#!/bin/bash
set -e

REPO="yourusername/goober"
LATEST_RELEASE=$(curl -s "https://api.github.com/repos/$REPO/releases/latest" | grep -o '"tag_name": "[^"]*"' | cut -d'"' -f4)

OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case $ARCH in
  x86_64) ARCH="x86_64" ;;
  arm64) ARCH="arm64" ;;
  *) echo "Unsupported architecture: $ARCH"; exit 1 ;;
esac

case $OS in
  linux) OS="Linux" ;;
  darwin) OS="Darwin" ;;
  *) echo "Unsupported OS: $OS"; exit 1 ;;
esac

DOWNLOAD_URL="https://github.com/$REPO/releases/download/$LATEST_RELEASE/goober_${OS}_${ARCH}.tar.gz"

echo "Downloading Goober $LATEST_RELEASE for $OS/$ARCH..."
curl -L "$DOWNLOAD_URL" | tar xz
sudo mv goober /usr/local/bin/
echo "Goober installed successfully!"
```

## 📚 Step 7: Documentation and Examples

### Create usage examples for common frameworks:

**Gin Web Framework**:

```bash
goober -build "go build -o server ./cmd/server" -run "./server"
```

**Fiber Framework**:

```bash
goober -build "go build -tags fiber -o app" -run "./app --port 3000"
```

**gRPC Service**:

```bash
goober -build "go build -o grpc-server ./cmd/grpc" -run "./grpc-server --grpc-port 9000"
```

## 🔄 Step 8: Maintenance and Updates

### Version Management

```bash
# Create new release
git tag v1.1.0
git push origin v1.1.0
```

### User Update Instructions

```bash
# Update global installation
go install github.com/yourusername/goober/cmd/goober@latest

# Or download latest binary
curl -L https://github.com/yourusername/goober/releases/latest/download/goober_Linux_x86_64.tar.gz | tar xz
```

## 🎉 Step 9: Promotion and Community

1. **Submit to awesome-go**: Add to [awesome-go](https://github.com/avelino/awesome-go)
2. **Create demo video**: Show the TUI in action
3. **Write blog post**: Compare with Air, Nodemon, etc.
4. **Social media**: Share on Twitter, Reddit r/golang
5. **Documentation**: Create GitHub wiki with advanced usage

## 🔧 Troubleshooting

### Common Issues:

1. **Import path errors**: Ensure all import paths use the correct GitHub URL
2. **Build failures**: Check Go version compatibility in go.mod
3. **Permission issues**: Users may need `sudo` for global installation

### Support Resources:

- GitHub Issues for bug reports
- GitHub Discussions for questions
- Wiki for advanced configuration

---

**Your Goober project is now ready for the world! 🎯🚀**
