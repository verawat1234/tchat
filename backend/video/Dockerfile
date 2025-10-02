# Build for Video Service - requires shared module
# Using Go 1.24 as required by go.mod
# Railway root directory should be: backend
FROM golang:1.24-alpine AS builder

# Set working directory
WORKDIR /app

# Copy shared module and video service
COPY shared/ ./shared/
COPY video/ ./video/

# Build from video directory
WORKDIR /app/video
RUN go mod tidy
RUN go build -o video-service ./main.go

# Runtime stage
FROM alpine:latest

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/video/video-service .

# Expose port
EXPOSE 8091

# Start the service
CMD ["./video-service"]
