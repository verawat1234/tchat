# Fresh Social Service Build
FROM golang:1.23-alpine

# Set working directory
WORKDIR /app

# Copy all files
COPY . .

# Build the application
RUN go mod tidy
RUN go build -o social-service ./main.go

# Expose the port
EXPOSE 8083

# Start the service
CMD ["./social-service"]
