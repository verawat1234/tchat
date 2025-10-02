# Simple build for Gateway Service (matching auth-final and social-fixed pattern)
FROM golang:1.23-alpine

# Set working directory
WORKDIR /app

# Copy all files from gateway directory (Railway builds from backend/infrastructure/gateway)
COPY . .

# Build the application
RUN go mod tidy
RUN go build -o gateway-service ./main.go

# Expose port
EXPOSE 8080

# Start the service
CMD ["./gateway-service"]
