# Simple build for Commerce Service (matching social-fixed pattern)
FROM golang:1.24-alpine

# Set working directory
WORKDIR /app

# Copy all files from commerce directory (Railway builds from backend/commerce)
COPY . .

# Build the application
RUN go mod tidy
RUN go build -o commerce-service ./main.go

# Expose port
EXPOSE 8084

# Start the service
CMD ["./commerce-service"]
