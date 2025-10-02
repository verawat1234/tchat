# Build for Auth Service - requires shared module
# Using Go 1.24 as required by go.mod
# Railway root directory should be: backend (not backend/auth)
FROM golang:1.24-alpine

# Set working directory
WORKDIR /app

# Copy shared module and auth service
COPY shared/ ./shared/
COPY auth/ ./auth/

# Build from auth directory
WORKDIR /app/auth
RUN go mod tidy
RUN go build -o auth-service ./main.go

# Expose port
EXPOSE 8081

# Start the service
CMD ["./auth-service"]
