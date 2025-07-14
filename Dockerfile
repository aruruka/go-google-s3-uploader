# Multi-stage build - minimize final image size

# Builder for app-server
FROM golang:1.24-alpine AS app-server-builder

WORKDIR /app-server

# Copy dependency files for app-server
COPY app-server/go.mod app-server/go.sum ./
RUN go mod download

# Copy app-server source code
COPY app-server/ ./
COPY shared/ ../shared/

# Build app-server binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app-server-main .

# Builder for auth-server
FROM golang:1.24-alpine AS auth-server-builder

WORKDIR /auth-server

# Copy dependency files for auth-server
COPY auth-server/go.mod auth-server/go.sum ./
RUN go mod download

# Copy auth-server source code
COPY auth-server/ ./
COPY shared/ ../shared/

# Build auth-server binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o auth-server-main .

# Final stage - minimal image
FROM alpine:latest

# Install CA certificates (required for HTTPS)
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy binaries from builder stages
COPY --from=app-server-builder /app-server/app-server-main ./app-server/app-server-main
COPY --from=auth-server-builder /auth-server/auth-server-main ./auth-server/auth-server-main

# Copy shared static files
COPY --from=app-server-builder /shared ./shared

# Copy the start script
COPY start.sh .
RUN chmod +x start.sh

# Expose ports (assuming both services run on different ports)
EXPOSE 8080
EXPOSE 8081

# Set environment variables (these will be overridden by App Runner)
ENV PORT_APP_SERVER=8080
ENV PORT_AUTH_SERVER=8081
ENV AWS_REGION=ap-northeast-1
ENV S3_BUCKET_NAME=raymond-go-s3-uploader-dev-2025

# Start both applications
CMD ["/app/start.sh"]
