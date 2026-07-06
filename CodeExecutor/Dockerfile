# =============================================================================
# code-executor — Multi-stage Dockerfile
# Stage 1: compile the Go binary
# Stage 2: minimal alpine runtime image
# =============================================================================

# ── Stage 1: Build ────────────────────────────────────────────────────────────
FROM golang:1.26-alpine AS builder

WORKDIR /build

# Copy dependency manifests first for layer caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source and build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-s -w" \
    -o code-executor \
    ./main.go

# ── Stage 2: Runtime ─────────────────────────────────────────────────────────
FROM alpine:3.23

WORKDIR /app
COPY --from=builder /build/code-executor .

# Use the existing nobody user (UID 65534) — already present in alpine
USER nobody

# code-executor has no listening port — it only makes outbound connections
# (RabbitMQ consumer, Redis publisher, Kubernetes API client)
ENTRYPOINT ["./code-executor"]
