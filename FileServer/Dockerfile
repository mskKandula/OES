FROM golang:1.26-alpine AS builder

# Set the Current Working Directory inside the container
WORKDIR /build

# We want to populate the module cache based on the go.{mod} files.
COPY go.mod go.sum .

RUN go mod download

COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o file-server .


########### Stage 2 ###########

#Define the running image
FROM alpine:3.23

# Create non-root user and group
RUN addgroup -S oesgroup && adduser -S oesuser -G oesgroup

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy built binary application from 'builder' image
COPY --from=builder --chown=oesuser:oesgroup /build/file-server .

# Switch to non-root user
USER oesuser

# This container exposes port 8887 to the outside world
EXPOSE 8887

# Run the binary program produced by `go install`
CMD ["./file-server"]
