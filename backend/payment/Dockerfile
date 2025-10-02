# Simple build for Payment Service (matching social-fixed pattern)
FROM golang:1.23-alpine

# Set working directory
WORKDIR /app

# Copy all files from payment directory (Railway builds from backend/payment)
COPY . .

# Build the application
RUN go mod tidy
RUN go build -o payment-service ./main.go

# Expose port
EXPOSE 8083

# Start the service
CMD ["./payment-service"]
