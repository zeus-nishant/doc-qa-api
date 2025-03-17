# Build stage
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /api ./cmd/api/main.go

# Run stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /api /api
COPY .env .env
EXPOSE 8080
CMD ["/api"]