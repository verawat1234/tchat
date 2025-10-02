# Build for Commerce Service - requires shared module
# Using Go 1.24 as required by go.mod
# Railway root directory should be: backend
FROM golang:1.24-alpine AS builder

# Set working directory
WORKDIR /app

# Copy shared module and commerce service
COPY shared/ ./shared/
COPY commerce/ ./commerce/

# Build from commerce directory
WORKDIR /app/commerce
RUN go mod tidy
RUN go build -o commerce-service ./main.go

# Runtime stage
FROM alpine:latest

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/commerce/commerce-service .

# Expose port
EXPOSE 8084

# Start the service
CMD ["./commerce-service"]
