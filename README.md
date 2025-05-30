# Go Module Redirect Server

A simple HTTP server that provides Go module redirects for `go get`. Automatically creates the required meta tags for Go modules and redirects users to `pkg.go.dev`.

## Features

The server generates for each request to `/{package-name}`:

- `go-import` meta tag for `go get`
- `go-source` meta tag for source code navigation
- Automatic redirect to `pkg.go.dev`

## Installation

```bash
git clone <repository-url>
cd go-redirect-server
go build -o vanity
```

## Configuration

Set the following environment variables:

| Variable     | Description                    | Example               |
| ------------ | ------------------------------ | --------------------- |
| `BASE_URL`   | Your domain (without https://) | `example.com`         |
| `GITHUB_URL` | Your GitHub repository prefix  | `github.com/username` |

## Usage

### Development

```bash
# Set environment variables
export BASE_URL="example.com"
export GITHUB_URL="github.com/username"

# Start server
go run main.go
```

### Direct execution

```bash
BASE_URL="example.com" GITHUB_URL="github.com/username" go run main.go
```

### With compiled binary

```bash
# Build binary
go build -o vanity

# Run
BASE_URL="example.com" GITHUB_URL="github.com/username" ./vanity
```

## Example

With settings:

- `BASE_URL=mymodules.com`
- `GITHUB_URL=github.com/myuser`

Request to `http://localhost:1337/mypackage` generates:

```html
<meta
  name="go-import"
  content="mymodules.com/mypackage git https://github.com/myuser/mypackage"
/>
<meta
  name="go-source"
  content="mymodules.com/mypackage https://github.com/myuser/mypackage https://github.com/myuser/mypackage/tree/master{/dir} https://github.com/myuser/mypackage/tree/master{/dir}/{file}#L{line}"
/>
```

And redirects to: `https://pkg.go.dev/mymodules.com/mypackage`

## Go Module Integration

After setup, Go modules can be installed like this:

```bash
go get mymodules.com/mypackage
```

## Port

The server runs on port 1337 by default. For production, you should use a reverse proxy (nginx, Caddy).

## Deployment

### Docker (Linux)

```dockerfile
# syntax=docker/dockerfile:1

FROM golang:1.24.0 AS build

ENV DEBIAN_FRONTEND=noninteractive

# Set destination for COPY
WORKDIR /app

# Copy go mod files first for better caching
COPY go.mod go.sum* ./
RUN go mod download

# Install any required packages
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o vanity .

# Use a minimal base image for the final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from build stage
COPY --from=build /app/vanity .

EXPOSE 1337

# Run the application
CMD ["./vanity"]
```

#### Docker usage:

```bash
# Build image
docker build -t vanity-server .

# Run container
docker run -p 1337:1337 \
  -e BASE_URL="example.com" \
  -e GITHUB_URL="github.com/username" \
  vanity-server
```

### Docker Compose

```yaml
version: "3.8"
services:
  vanity:
    build: .
    ports:
      - "1337:1337"
    environment:
      - BASE_URL=example.com
      - GITHUB_URL=github.com/username
    restart: unless-stopped
```

### Systemd Service

```ini
[Unit]
Description=Go Module Redirect Server
After=network.target

[Service]
Type=simple
User=www-data
WorkingDirectory=/opt/vanity-server
ExecStart=/opt/vanity-server/vanity
Environment=BASE_URL=example.com
Environment=GITHUB_URL=github.com/username
Restart=always

[Install]
WantedBy=multi-user.target
```

## Development

### Requirements

- Go 1.21 or later
- Environment variables `BASE_URL` and `GITHUB_URL` must be set

### Testing

```bash
# Start the server
BASE_URL="localhost:1337" GITHUB_URL="github.com/test" go run main.go

# Test in another terminal
curl -v http://localhost:1337/testpackage
```

The response should include the correct meta tags and redirect to pkg.go.dev.
