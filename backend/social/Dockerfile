# Build for Social Service - requires shared module
# Using Go 1.24 as required by go.mod
# Railway root directory should be: backend (not backend/social)
FROM golang:1.24-alpine

# Set working directory
WORKDIR /app

# Copy shared module and social service
COPY shared/ ./shared/
COPY social/ ./social/

# Build from social directory
WORKDIR /app/social
RUN go mod tidy
RUN go build -o social-service ./main.go

# Expose port
EXPOSE 8083

# Start the service
CMD ["./social-service"]
