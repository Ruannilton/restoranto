# Base image
FROM golang:1.23 as builder

# Set working directory
WORKDIR /app

# Copy Go modules and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy application source code
COPY . .

# Build the application
RUN go build -o main .

# Use a compatible base image with an updated glibc version
FROM debian:bookworm-slim

# Set working directory
WORKDIR /app

# Install necessary runtime dependencies
RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates && \
    rm -rf /var/lib/apt/lists/*

# Copy the compiled binary from the builder stage
COPY --from=builder /app/main ./main

# Expose the application port (to match docker-compose service ports)
EXPOSE 8080

# Command to run the application
CMD ["./main"]
