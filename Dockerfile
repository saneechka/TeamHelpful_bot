# Multi-stage build for minimal image size
FROM golang:1.21-alpine AS builder

# Install necessary dependencies
RUN apk add --no-cache git gcc musl-dev

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o helpbot .

# Final image
FROM alpine:3.16

# Install runtime dependencies
RUN apk add --no-cache ca-certificates

# Set working directory
WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /app/helpbot .

# Copy any necessary files
COPY --from=builder /app/users.db ./users.db

# Run the application
CMD ["./helpbot"] 