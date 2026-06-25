FROM golang:1.26-alpine AS builder1

# Set the Current Working Directory inside the container
WORKDIR /app/OES/fileserver

# We want to populate the module cache based on the go.{mod} files.
COPY ./FileServer/go.mod .

RUN go mod download

COPY ./FileServer/ .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o files .


########### Stage 2 ###########

#Define the running image
FROM alpine:3.23

# Create non-root user and group
RUN addgroup -S oesgroup && adduser -S oesuser -G oesgroup

# Set the Current Working Directory inside the container
WORKDIR /app/OES/fileserver

# Copy built binary application from 'builder' image
COPY --from=builder1 /app/OES/fileserver/files .

# Own the workdir by the non-root user
RUN chown -R oesuser:oesgroup /app/OES/fileserver

# Switch to non-root user
USER oesuser

# This container exposes port 8887 to the outside world
EXPOSE 8887

# Run the binary program produced by `go install`
CMD ["./files"]
