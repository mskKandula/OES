# Stage 1 as builder
FROM golang:1.26-alpine AS builder

# Set the Current Working Directory inside the container
WORKDIR /build

# Copy Go Application Dependencies
COPY go.mod go.sum .

# Download dependencies 
RUN go mod download

# Copy code files
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o server .

########### Stage 2 ###########

#Define the running image
FROM alpine:3.23

# Create non-root user and group
RUN addgroup -S oesgroup && adduser -S oesuser -G oesgroup

# Set the Current Working Directory inside the container
WORKDIR /app

COPY --from=builder --chown=oesuser:oesgroup /build/config.json /build/server .

# Switch to non-root user
USER oesuser

# This container exposes port 9000 to the outside world
EXPOSE 9000

# Run the binary program produced by `go install`
CMD ["./server"]
