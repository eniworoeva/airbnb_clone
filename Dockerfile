# Build stage
FROM golang:1.24.1-alpine AS builder

# Set working directory
WORKDIR /app

# Install git and ca-certificates (needed for git operations and HTTPS)
RUN apk add --no-cache git ca-certificates

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main cmd/server/main.go

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

# Create non-root user
RUN addgroup -g 1001 -S appuser && \
    adduser -u 1001 -S appuser -G appuser

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/main .

# Copy any config files if needed
COPY --from=builder /app/.env.example .env.example

# Change ownership to non-root user
RUN chown -R appuser:appuser /root/
USER appuser

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Run the application
CMD ["./main"]