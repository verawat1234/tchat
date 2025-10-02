# Simple build for Calling Service (matching social-fixed pattern)
FROM golang:1.23-alpine

# Set working directory
WORKDIR /app

# Copy all files from calling directory (Railway builds from backend/calling)
COPY . .

# Build the application
RUN go mod tidy
RUN go build -o calling-service ./main.go

# Expose port
EXPOSE 8093

# Start the service
CMD ["./calling-service"]
