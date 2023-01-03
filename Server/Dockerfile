# Stage 1 as builder
FROM golang:1.18-alpine as builder

# Set the Current Working Directory inside the container
WORKDIR /app/OES/server

# Copy Go Application Dependencies
COPY ./Server/go.mod .
COPY ./Server/go.sum .

# To keep only actual dependencies
RUN go mod tidy

# Download dependencies 
RUN go mod download

# Copy code files
COPY ./Server/ .

# Build the Go app
RUN go build -o main .

########### Stage 2 ###########

#Define the running image
FROM alpine:3.17

# For CDN to serve uploaded video files
RUN mkdir -p /app/OES/media/video

# Install git && ffmpeg
RUN apk update && \
    apk upgrade

# Set the Current Working Directory inside the container
WORKDIR /app/OES/server

COPY --from=builder /app/OES/server/config.json .

# Copy built binary application from 'builder' image
COPY --from=builder /app/OES/server/main .

# This container exposes port 8080 to the outside world
EXPOSE 9000

# Run the binary program produced by `go install`
CMD ["./main"]