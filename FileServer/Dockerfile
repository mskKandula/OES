FROM golang:1.16-alpine as builder1

# Install git
RUN apk update && \
    apk upgrade && \
    apk add git

# Set the Current Working Directory inside the container
WORKDIR /app/OES/fileserver

# We want to populate the module cache based on the go.{mod,sum} files.
COPY ./FileServer/go.mod .

# To keep only actual dependencies
RUN go mod tidy

RUN go mod download

COPY ./FileServer/ .

# Build the Go app
RUN go build -o files .


########### Stage 2 ###########

#Define the running image
FROM alpine:3.17

# Install
RUN apk update && \
    apk upgrade

# Set the Current Working Directory inside the container
WORKDIR /app/OES/fileserver

# Copy built binary application from 'builder' image
COPY --from=builder1 /app/OES/fileserver/files .

# This container exposes port 8080 to the outside world
EXPOSE 8887

# Run the binary program produced by `go install`
CMD ["./files"]