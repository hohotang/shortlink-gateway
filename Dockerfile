# syntax=docker/dockerfile:1.4

# Base image with full compatibility
FROM golang:1.24.1

# Set working directory
WORKDIR /app

# Copy source code
COPY . .

# Download dependencies
RUN go mod download

# Build the application with verbose output
RUN go build -v -o gateway ./cmd/gateway && \
    chmod +x gateway && \
    ls -la

# Expose port
EXPOSE 8080

# Keep container running for debugging
CMD ["./gateway"]
