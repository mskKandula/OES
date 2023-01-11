FROM golang:1.16-alpine as builder1

# Install git
RUN apk update && \
    apk upgrade && \
    apk add git

# Set the Current Working Directory inside the container
WORKDIR /app/mqServer

# We want to populate the module cache based on the go.{mod,sum} files.
COPY ./MQServer/go.mod .

# To keep only actual dependencies
RUN go mod tidy

RUN go mod download

COPY ./MQServer/ .

# Build the Go app
RUN go build -o main .


########### Stage 2 ###########

#Define the running image
FROM alpine:3.17

# Install git && ffmpeg
RUN apk update && \
    apk upgrade && \
    apk add --no-cache ffmpeg

# Set the Current Working Directory inside the container
WORKDIR /app/mqServer

# Copy built binary application from 'builder' image
COPY --from=builder1 /app/mqServer/main .

# This container exposes port 5672 to the outside world
# EXPOSE 5672

# Run the binary program produced by `go install`
CMD ["./main"]