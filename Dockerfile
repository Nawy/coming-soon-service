# Build stage
FROM golang:1.25.1-alpine AS builder

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY main.go ./

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o coming-soon-service .

# Runtime stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# Set working directory
WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/coming-soon-service .

# Create a non-root user to run the application
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser && \
    chown -R appuser:appuser /app

# Switch to non-root user
USER appuser

# Expose port 8080
EXPOSE 8080

# Set environment variable placeholder (must be overridden at runtime)
ENV SECRET_TOKEN=""

# Volume for emails.txt file (mounted from host)
VOLUME ["/app/data"]

# Set the emails file path to use the mounted volume
ENV EMAIL_FILE_PATH=/app/data/emails.txt

# Run the application
CMD ["./coming-soon-service"]
