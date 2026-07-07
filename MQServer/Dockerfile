FROM golang:1.26-alpine AS builder

# Set the Current Working Directory inside the container
WORKDIR /build
# We want to populate the module cache based on the go.{mod} files.
COPY go.mod go.sum .
RUN go mod download

COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o mq-server .

########### Stage 2 ###########

#Define the running image
FROM alpine:3.23

# Install ffmpeg
RUN apk add --no-cache ffmpeg

# Create non-root user and group
RUN addgroup -S oesgroup && adduser -S oesuser -G oesgroup

# Set the Current Working Directory inside the container
WORKDIR /app

COPY --from=builder --chown=oesuser:oesgroup /build/mq-server .

# Switch to non-root user
USER oesuser

# Run the binary program produced by `go install`
CMD ["./mq-server"]
