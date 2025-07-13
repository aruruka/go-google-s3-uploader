#!/bin/bash

# 🚀 Go Google S3 Uploader - Local testing script

echo "🔧 Setting up Go Google S3 Uploader for local testing..."

# Check required environment variables
if [ -z "$GOOGLE_CLIENT_ID" ] || [ -z "$GOOGLE_CLIENT_SECRET" ]; then
    echo "⚠️  Missing OAuth credentials. Please set:"
    echo "   export GOOGLE_CLIENT_ID=your_client_id"
    echo "   export GOOGLE_CLIENT_SECRET=your_client_secret"
    echo ""
    echo "📁 Or copy credentials to OAuth-Credentials/ directory"
fi

# Set local development environment variables
export PORT=8080
export AWS_REGION=ap-northeast-1
export S3_BUCKET_NAME=raymond-go-s3-uploader-dev-2025
export ENVIRONMENT=development
export AUTH_SERVER_URL=http://localhost:8081
export APP_SERVER_URL=http://localhost:8080
export AWS_PROFILE=go-aws-sdk

echo "🌍 Environment configured:"
echo "   PORT: $PORT"
echo "   AWS_REGION: $AWS_REGION"
echo "   S3_BUCKET_NAME: $S3_BUCKET_NAME"
echo "   AWS_PROFILE: $AWS_PROFILE"
echo "   ENVIRONMENT: $ENVIRONMENT"

# Check Go version
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go 1.22+"
    exit 1
fi

GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
echo "✅ Go version: $GO_VERSION"

# Build application
echo "🔨 Building application..."
cd app-server
go mod tidy

if go build -o ../app-server-binary .; then
    echo "✅ Build successful!"
    cd ..
    
    echo "🚀 Starting application on :$PORT..."
    echo "🌐 Visit: http://localhost:$PORT"
    echo "🔒 OAuth: http://localhost:$PORT/auth/google"
    echo "❤️  Health: http://localhost:$PORT/health"
    echo ""
    echo "Press Ctrl+C to stop the server"
    
    ./app-server-binary
else
    echo "❌ Build failed!"
    exit 1
fi
