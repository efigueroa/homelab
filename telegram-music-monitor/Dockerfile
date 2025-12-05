# Multi-stage Dockerfile for Telegram Music Monitor
# This uses two stages: build and runtime
# The build stage compiles the Go application
# The runtime stage creates a minimal image with just the binary

# Stage 1: Build stage
# Use official Go image as the build environment
FROM golang:1.21-alpine AS builder

# Install build dependencies
# git is needed for go mod download
RUN apk add --no-cache git

# Set working directory inside the container
WORKDIR /build

# Copy go.mod and go.sum files first
# This allows Docker to cache the dependency layer
COPY go.mod go.sum ./

# Download dependencies
# This layer is cached unless go.mod or go.sum changes
RUN go mod download

# Copy the entire source code
COPY . .

# Build the application
# CGO_ENABLED=0 creates a statically linked binary (no C dependencies)
# -ldflags="-w -s" strips debug information to reduce binary size
# -o /app/music-monitor specifies the output binary path
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/music-monitor .

# Stage 2: Runtime stage
# Use a minimal base image for the final container
FROM alpine:latest

# Install CA certificates for HTTPS requests
# tzdata for timezone support
RUN apk --no-cache add ca-certificates tzdata

# Create a non-root user to run the application
# This is a security best practice
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser

# Create data directory for database and session files
# This will be mounted as a volume
RUN mkdir -p /data && chown appuser:appuser /data

# Set working directory
WORKDIR /app

# Copy the compiled binary from the builder stage
COPY --from=builder /app/music-monitor .

# Change ownership to the non-root user
RUN chown appuser:appuser /app/music-monitor

# Switch to non-root user
USER appuser

# Expose the API port
# This doesn't actually publish the port, it's documentation
EXPOSE 8080

# Set the entrypoint to run the application
ENTRYPOINT ["/app/music-monitor"]
