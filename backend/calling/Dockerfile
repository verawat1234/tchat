# Build for Calling Service - requires shared module
# Using Go 1.24 as required by go.mod
# Railway root directory should be: backend
FROM golang:1.24-alpine AS builder

# Set working directory
WORKDIR /app

# Copy shared module and calling service
COPY shared/ ./shared/
COPY calling/ ./calling/

# Build from calling directory
WORKDIR /app/calling
RUN go mod tidy
RUN go build -o calling-service ./main.go

# Runtime stage
FROM alpine:latest

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/calling/calling-service .

# Expose port
EXPOSE 8093

# Start the service
CMD ["./calling-service"]
