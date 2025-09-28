# Fresh Messaging Service Build
FROM golang:1.23-alpine

# Set working directory
WORKDIR /app

# Copy all files
COPY . .

# Build the application
RUN go mod tidy
RUN go build -o messaging-service ./main.go

# Expose the port
EXPOSE 8082

# Start the service
CMD ["./messaging-service"]