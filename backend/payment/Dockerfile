# Build for Payment Service - requires shared module
# Using Go 1.24 as required by go.mod
# Railway root directory should be: backend
FROM golang:1.24-alpine AS builder

# Set working directory
WORKDIR /app

# Copy shared module and payment service
COPY shared/ ./shared/
COPY payment/ ./payment/

# Build from payment directory
WORKDIR /app/payment
RUN go mod tidy
RUN go build -o payment-service ./main.go

# Runtime stage
FROM alpine:latest

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/payment/payment-service .

# Expose port
EXPOSE 8083

# Start the service
CMD ["./payment-service"]
