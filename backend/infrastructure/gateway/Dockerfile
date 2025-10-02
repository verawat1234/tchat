# Gateway Service Dockerfile
# Build context: backend/ (Railway root directory)
FROM golang:1.23-alpine

# Set working directory
WORKDIR /app

# Copy go.work and go.work.sum from backend root
COPY go.work* ./

# Copy shared module
COPY shared/ ./shared/

# Copy gateway service files
COPY infrastructure/gateway/ ./infrastructure/gateway/

# Set working directory to gateway
WORKDIR /app/infrastructure/gateway

# Build the application
RUN go mod tidy
RUN go build -o gateway-service .

# Expose port
EXPOSE 8080

# Start the service
CMD ["./gateway-service"]
