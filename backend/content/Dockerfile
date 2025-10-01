# Multi-stage build for Content Service
FROM golang:1.23-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy module manifests for dependency caching
COPY shared/go.mod shared/go.sum ./shared/
COPY content/go.mod content/go.sum ./content/

# Resolve module dependencies
RUN cd shared && go mod download
RUN cd content && go mod download

# Copy source code
COPY shared/ ./shared/
COPY content/ ./content/

# Build the binary
WORKDIR /app/content
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o content-service ./main.go

# Production stage
FROM alpine:3.18

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates tzdata

# Create non-root user
RUN addgroup -g 1001 -S tchat && \
    adduser -S tchat -u 1001

# Set working directory
WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /app/content/content-service .

# Change ownership to non-root user
RUN chown tchat:tchat /app/content-service

# Switch to non-root user
USER tchat

# Expose port
EXPOSE 8086

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8086/health || exit 1

# Run the service
CMD ["./content-service"]
