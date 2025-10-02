# Simple build for Content Service (matching social-fixed pattern)
FROM golang:1.24-alpine

# Set working directory
WORKDIR /app

# Copy all files from content directory (Railway builds from backend/content)
COPY . .

# Build the application
RUN go mod tidy
RUN go build -o content-service ./main.go

# Expose port
EXPOSE 8086

# Start the service
CMD ["./content-service"]
