# Build stage
FROM golang:1.24-alpine AS builder

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/server ./cmd/server

# Final stage - use Alpine for smaller, secure image
FROM alpine:3.19

# Define non-root user IDs for easier maintenance
ARG NONROOT_UID=65532
ARG NONROOT_GID=65532

# Create non-root user for security
RUN addgroup -g ${NONROOT_GID} -S nonroot && \
    adduser -u ${NONROOT_UID} -S nonroot -G nonroot

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/server .

# Change ownership to non-root user
RUN chown -R nonroot:nonroot /app

# Switch to non-root user
USER nonroot

# Expose port
EXPOSE 8080

# Run the application
ENTRYPOINT ["./server"]
