# Stage 1: Build the Go binary
FROM golang:1.21-alpine AS builder

# Install git and necessary tools
RUN apk add --no-cache git

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go.mod and go.sum and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the Go app
RUN go build -o main .

# Stage 2: Create a minimal final image
FROM alpine:latest

# Set working directory
WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/main .

# Expose port (optional, e.g., 8080)
EXPOSE 8080

# Command to run the executable
CMD ["./main"]
