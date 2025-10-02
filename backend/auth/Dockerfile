# Simple build for Auth Service (matching social-fixed pattern)
# Using Go 1.24 as required by go.mod
FROM golang:1.24-alpine

# Set working directory
WORKDIR /app

# Copy all files from auth directory (Railway builds from backend/auth)
COPY . .

# Build the application
RUN go mod tidy
RUN go build -o auth-service ./main.go

# Expose port
EXPOSE 8081

# Start the service
CMD ["./auth-service"]
