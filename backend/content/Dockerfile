# Build for Content Service - requires shared module
# Using Go 1.24 as required by go.mod
# Railway root directory should be: backend
FROM golang:1.24-alpine AS builder

# Set working directory
WORKDIR /app

# Copy shared module and content service
COPY shared/ ./shared/
COPY content/ ./content/

# Build from content directory
WORKDIR /app/content
RUN go mod tidy
RUN go build -o content-service ./main.go

# Runtime stage
FROM alpine:latest

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/content/content-service .

# Expose port
EXPOSE 8086

# Start the service
CMD ["./content-service"]
