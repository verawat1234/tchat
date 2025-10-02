# Build for Notification Service - requires shared module
# Using Go 1.24 as required by go.mod
# Railway root directory should be: backend
FROM golang:1.24-alpine AS builder

# Set working directory
WORKDIR /app

# Copy shared module and notification service
COPY shared/ ./shared/
COPY notification/ ./notification/

# Build from notification directory
WORKDIR /app/notification
RUN go mod tidy
RUN go build -o notification-service ./main.go

# Runtime stage
FROM alpine:latest

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/notification/notification-service .

# Expose port
EXPOSE 8087

# Start the service
CMD ["./notification-service"]
