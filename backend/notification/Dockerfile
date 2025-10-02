# Simple build for Notification Service (matching social-fixed pattern)
FROM golang:1.23-alpine

# Set working directory
WORKDIR /app

# Copy all files from notification directory (Railway builds from backend/notification)
COPY . .

# Build the application
RUN go mod tidy
RUN go build -o notification-service ./main.go

# Expose port
EXPOSE 8085

# Start the service
CMD ["./notification-service"]
