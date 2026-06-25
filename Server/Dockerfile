# Stage 1 as builder
FROM golang:1.26-alpine AS builder

# Set the Current Working Directory inside the container
WORKDIR /app/OES/server

# Copy Go Application Dependencies
COPY ./Server/go.mod ./Server/go.sum ./

# Download dependencies 
RUN go mod download

# Copy code files
COPY ./Server/ .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o main .

########### Stage 2 ###########

#Define the running image
FROM alpine:3.23

# Set the Current Working Directory inside the container
WORKDIR /app/OES/server

COPY --from=builder /app/OES/server/config.json .

# Copy built binary application from 'builder' image
COPY --from=builder /app/OES/server/main .

# This container exposes port 8080 to the outside world
EXPOSE 9000

# Run the binary program produced by `go install`
CMD ["./main"]
