# Build stage
FROM golang:1.23-alpine AS builder

# Install build dependencies
RUN apk add --no-cache gcc musl-dev postgresql-dev

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -a -ldflags="-s -w" -o validra-engine ./src

# Final stage
FROM alpine:3.19

# Install PostgreSQL client and CA certificates
RUN apk add --no-cache ca-certificates postgresql-client

# Create a non-root user and group
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/validra-engine .

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 8080

# Set health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 CMD wget --no-verbose --spider http://localhost:8080/health || exit 1

# Run the application
CMD ["./validra-engine"]