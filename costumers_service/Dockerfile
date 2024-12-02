# Example Dockerfile for a Go microservice
# Base image
FROM golang:1.21 as builder

# Set working directory
WORKDIR /app

# Copy Go modules and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy application source code
COPY . .

# Build the application
RUN go build -o main .

# Start a new stage for a clean image
FROM debian:bullseye-slim

# Set working directory
WORKDIR /app

# Copy the compiled binary from the builder stage
COPY --from=builder /app/main ./main

# Expose the application port (to match docker-compose service ports)
EXPOSE 8080

# Command to run the application
CMD ["./main"]
