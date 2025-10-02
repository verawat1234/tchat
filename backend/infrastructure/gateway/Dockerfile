# Simple build for Gateway Service (matching auth-final and social-fixed pattern)
# Using Go 1.24 as required by go.mod
FROM golang:1.24-alpine

# Set working directory
WORKDIR /app

# Copy all files from gateway directory (Railway builds from backend/infrastructure/gateway)
COPY . .

# Build the application
RUN go mod tidy
RUN go build -o gateway-service .

# Expose port
EXPOSE 8080

# Start the service with absolute path
CMD ["/app/gateway-service"]
