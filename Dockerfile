# Build stage
FROM golang:1.23-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /api ./cmd/api/main.go

# Run stage
FROM alpine:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /api /api

# Copy the .env file (if needed)
COPY .env .env

# Expose the application port
EXPOSE 8080

# Command to run the application
CMD ["/api"]