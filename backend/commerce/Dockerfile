# Build stage
FROM golang:1.24-alpine AS builder

# Set working directory
WORKDIR /app

# Copy all files (including shared module copied into commerce directory)
COPY . .

# Build the application
RUN GOWORK=off go mod download
RUN GOWORK=off go build -o commerce-service ./main.go

# Runtime stage
FROM alpine:latest

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/commerce-service .

# Expose port
EXPOSE 8084

# Start the service
CMD ["./commerce-service"]
