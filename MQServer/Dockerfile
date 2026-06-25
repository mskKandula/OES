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

# Install ffmpeg
RUN apk add --no-cache ffmpeg

# Create non-root user and group
RUN addgroup -S oesgroup && adduser -S oesuser -G oesgroup

# Set the Current Working Directory inside the container
WORKDIR /app/OES/mqServer

# Copy built binary application from 'builder' image
COPY --from=builder1 /app/OES/mqServer/main .

# Own the workdir by the non-root user
RUN chown -R oesuser:oesgroup /app/OES/mqServer

# Switch to non-root user
USER oesuser

# Run the binary program produced by `go install`
CMD ["./main"]
