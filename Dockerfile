# Multi-stage build for combined service

# Builder stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy all module files first for better layer caching
COPY go.mod go.sum ./
COPY auth-server/go.mod auth-server/go.sum ./auth-server/
COPY app-server/go.mod app-server/go.sum ./app-server/
COPY shared/go.mod ./shared/

# Download dependencies
RUN go mod download
RUN cd auth-server && go mod download
RUN cd app-server && go mod download
RUN cd shared && go mod download

# Copy all source code
COPY . .

# Build the combined service
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o combined-service ./main.go

# Final stage - minimal image
FROM alpine:latest

# Install CA certificates (required for HTTPS)
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy the combined service binary
COPY --from=builder /app/combined-service ./combined-service

# Copy template files
COPY --from=builder /app/auth-server/templates ./auth-server/templates
COPY --from=builder /app/app-server/templates ./app-server/templates
COPY --from=builder /app/shared ./shared

# Expose port 8080 (App Runner will set PORT environment variable)
EXPOSE 8080

# Set default environment variables (these will be overridden by App Runner)
ENV PORT=8080
ENV AWS_REGION=ap-northeast-1

# Start the combined service
CMD ["./combined-service"]
