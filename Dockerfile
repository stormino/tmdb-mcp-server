# Build stage
FROM golang:1.25.2-alpine3.22 AS builder
WORKDIR /build
RUN apk add --no-cache ca-certificates git

# Copy go files
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY *.go ./

# Build the application with security flags
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -trimpath \
    -ldflags="-w -s -extldflags '-static'" \
    -a \
    -o tmdb-mcp-server .

# Run stages: using distroless for minimal attack surface

# HTTP mode stage
FROM gcr.io/distroless/static-debian12:nonroot AS http
WORKDIR /app

# Copy CA certificates and binary from builder
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /build/tmdb-mcp-server /app/tmdb-mcp-server

# distroless images already run as non-root user (uid/gid: 65532)
USER nonroot:nonroot

# Default HTTP port (can be overridden with -e HTTP_PORT=xxxx)
ENV HTTP_PORT=8080

EXPOSE 8080

ENTRYPOINT ["/app/tmdb-mcp-server"]
CMD ["--mode", "http"]

# STDIO mode stage (default)
FROM gcr.io/distroless/static-debian12:nonroot AS stdio
WORKDIR /app

# Copy CA certificates and binary from builder
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /build/tmdb-mcp-server /app/tmdb-mcp-server

# distroless images already run as non-root user (uid/gid: 65532)
USER nonroot:nonroot

ENV MCP_MODE=stdio

ENTRYPOINT ["/app/tmdb-mcp-server"]
CMD ["--mode", "stdio"]
