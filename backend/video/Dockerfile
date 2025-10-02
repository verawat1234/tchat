# Simple build for Video Service (matching social-fixed pattern)
FROM golang:1.24-alpine

# Set working directory
WORKDIR /app

# Copy all files from video directory (Railway builds from backend/video)
COPY . .

# Build the application
RUN go mod tidy
RUN go build -o video-service ./main.go

# Expose port
EXPOSE 8091

# Start the service
CMD ["./video-service"]
