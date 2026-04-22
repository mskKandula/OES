FROM golang:1.26-alpine AS builder1

# Set the Current Working Directory inside the container
WORKDIR /app/OES/mqServer

# We want to populate the module cache based on the go.{mod} files.
COPY ./MQServer/go.mod ./MQServer/go.sum ./
RUN go mod download

COPY ./MQServer/ .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o main .

########### Stage 2 ###########

#Define the running image
FROM alpine:3.21


# Update CA certificates first to fix SSL in corporate/proxy environments
RUN apk add --no-cache ca-certificates && update-ca-certificates

# Install ffmpeg
RUN apk add --no-cache ffmpeg

# Set the Current Working Directory inside the container
WORKDIR /app/OES/mqServer

# Copy built binary application from 'builder' image
COPY --from=builder1 /app/OES/mqServer/main .

# Run the binary program produced by `go install`
CMD ["./main"]
