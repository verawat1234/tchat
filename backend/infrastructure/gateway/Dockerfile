# Build stage
FROM golang:1.24-bookworm AS builder

# Set working directory
WORKDIR /app

# Copy all files (including shared module copied into gateway directory)
COPY . .

# Build the application with all .go files
RUN GOWORK=off go mod download
RUN GOWORK=off CGO_ENABLED=1 GOOS=linux go build -o gateway-service .

# Runtime stage - using Debian slim for glibc compatibility
FROM debian:bookworm-slim

# Install CA certificates and other runtime dependencies
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/gateway-service .

# Expose port
EXPOSE 8080

# Start the service
CMD ["./gateway-service"]
