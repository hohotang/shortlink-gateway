# Build stage
FROM golang:1.24.1-alpine AS builder

# Set working directory
WORKDIR /build

# Copy go.mod and go.sum first to leverage Docker cache
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build with static linking
RUN CGO_ENABLED=0 GOOS=linux go build  -o gateway ./cmd/gateway

# Final lightweight stage
FROM alpine:latest

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /build/gateway .
COPY config.yaml .

# Expose port
EXPOSE 8080

# Run the application
CMD ["./gateway"]
