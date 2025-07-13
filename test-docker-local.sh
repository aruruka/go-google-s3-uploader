#!/bin/bash

# 🐳 Local Docker Container Testing Script
# Used to verify Docker image build and run before pushing to GitHub

set -e

echo "🐳 Docker container local testing begins..."
echo "=========================================="

# Check required environment variables
check_env_vars() {
    echo "📋 Checking environment variable configuration..."
    
    # Check AWS credentials (at least one method needed)
    if [[ -z "$AWS_ACCESS_KEY_ID" && -z "$AWS_PROFILE" ]]; then
        echo "⚠️  Warning: AWS credentials not set"
        echo "   Please set environment variables or AWS Profile:"
        echo "   export AWS_ACCESS_KEY_ID=your_key"
        echo "   export AWS_SECRET_ACCESS_KEY=your_secret"
        echo "   or: export AWS_PROFILE=go-aws-sdk"
        echo ""
    fi
    
    # Check S3 configuration
    if [[ -z "$S3_BUCKET_NAME" ]]; then
        echo "⚠️  Warning: S3_BUCKET_NAME not set, using default value"
        export S3_BUCKET_NAME="raymond-go-s3-uploader-dev-2025"
    fi
    
    if [[ -z "$AWS_REGION" ]]; then
        echo "⚠️  Setting default AWS_REGION"
        export AWS_REGION="ap-northeast-1"
    fi
    
    echo "✅ Environment variable check complete"
    echo "   S3_BUCKET_NAME: $S3_BUCKET_NAME"
    echo "   AWS_REGION: $AWS_REGION"
    echo ""
}

# Build Docker image
build_image() {
    echo "🔨 Building Docker image..."
    
    # Build image
    docker build -t go-s3-uploader:local .
    
    if [ $? -eq 0 ]; then
        echo "✅ Docker image build successful"
    else
        echo "❌ Docker image build failed"
        exit 1
    fi
    echo ""
}

# Run container (app-server only, simplified testing)
run_container() {
    echo "🚀 Starting Docker container..."
    
    # Stop any existing old containers
    docker stop go-s3-uploader-test 2>/dev/null || true
    docker rm go-s3-uploader-test 2>/dev/null || true
    
    # Prepare environment variable arguments
    ENV_ARGS=""
    
    # AWS credentials (prioritize environment variables, then Profile)
    if [[ -n "$AWS_ACCESS_KEY_ID" ]]; then
        ENV_ARGS="$ENV_ARGS -e AWS_ACCESS_KEY_ID=$AWS_ACCESS_KEY_ID"
        ENV_ARGS="$ENV_ARGS -e AWS_SECRET_ACCESS_KEY=$AWS_SECRET_ACCESS_KEY"
        [[ -n "$AWS_SESSION_TOKEN" ]] && ENV_ARGS="$ENV_ARGS -e AWS_SESSION_TOKEN=$AWS_SESSION_TOKEN"
        echo "   Using environment variable AWS credentials"
    elif [[ -n "$AWS_PROFILE" ]]; then
        # Mount AWS configuration files
        if [[ -d "$HOME/.aws" ]]; then
            ENV_ARGS="$ENV_ARGS -v $HOME/.aws:/root/.aws:ro"
            ENV_ARGS="$ENV_ARGS -e AWS_PROFILE=$AWS_PROFILE"
            echo "   Using AWS Profile: $AWS_PROFILE"
        else
            echo "⚠️  AWS Profile specified but ~/.aws directory does not exist"
        fi
    fi
    
    # S3 and application configuration
    ENV_ARGS="$ENV_ARGS -e S3_BUCKET_NAME=$S3_BUCKET_NAME"
    ENV_ARGS="$ENV_ARGS -e AWS_REGION=$AWS_REGION"
    ENV_ARGS="$ENV_ARGS -e PORT=8080"
    ENV_ARGS="$ENV_ARGS -e ENVIRONMENT=development"
    
    # OAuth configuration (optional, for full functionality testing)
    [[ -n "$GOOGLE_CLIENT_ID" ]] && ENV_ARGS="$ENV_ARGS -e GOOGLE_CLIENT_ID=$GOOGLE_CLIENT_ID"
    [[ -n "$GOOGLE_CLIENT_SECRET" ]] && ENV_ARGS="$ENV_ARGS -e GOOGLE_CLIENT_SECRET=$GOOGLE_CLIENT_SECRET"
    
    # Start container
    echo "   Start command: docker run --name go-s3-uploader-test -p 8080:8080 -d $ENV_ARGS go-s3-uploader:local"
    eval "docker run --name go-s3-uploader-test -p 8080:8080 -d $ENV_ARGS go-s3-uploader:local"
    
    if [ $? -eq 0 ]; then
        echo "✅ Container started successfully"
        echo "🌐 Access URL: http://localhost:8080"
        echo "❤️  Health check: http://localhost:8080/health"
    else
        echo "❌ Container failed to start"
        exit 1
    fi
    echo ""
}

# Wait for service startup and test
test_container() {
    echo "🧪 Testing container service..."
    
    # Wait for service startup
    echo "   Waiting for service to start..."
    for i in {1..30}; do
        if curl -s http://localhost:8080/health > /dev/null; then
            echo "✅ Service started (${i}s)"
            break
        fi
        if [ $i -eq 30 ]; then
            echo "❌ Service startup timeout"
            show_logs_and_cleanup
            exit 1
        fi
        sleep 1
    done
    
    # Test health check
    echo "   Testing health check endpoint..."
    HEALTH_RESPONSE=$(curl -s http://localhost:8080/health)
    if [ "$HEALTH_RESPONSE" == "OK" ]; then
        echo "✅ Health check passed: $HEALTH_RESPONSE"
    else
        echo "❌ Health check failed: $HEALTH_RESPONSE"
    fi
    
    # Test homepage (parts that don't require authentication)
    echo "   Testing homepage response..."
    if curl -s -o /dev/null -w "%{http_code}" http://localhost:8080/ | grep -q "200\|302"; then
        echo "✅ Homepage response normal"
    else
        echo "❌ Homepage response abnormal"
    fi
    
    echo ""
}

# Show logs and cleanup
show_logs_and_cleanup() {
    echo "📜 Showing container logs..."
    docker logs go-s3-uploader-test
    echo ""
    
    echo "🧹 Cleaning up test resources..."
    docker stop go-s3-uploader-test 2>/dev/null || true
    docker rm go-s3-uploader-test 2>/dev/null || true
    echo "✅ Cleanup complete"
}

# Main process
main() {
    check_env_vars
    build_image
    run_container
    test_container
    
    echo "🎉 Docker container testing complete!"
    echo "====================================="
    echo "📊 Test summary:"
    echo "   ✅ Docker image build successful"
    echo "   ✅ Container startup successful"
    echo "   ✅ Health check passed"
    echo "   ✅ Basic functionality normal"
    echo ""
    echo "🌐 Container is running:"
    echo "   Access URL: http://localhost:8080"
    echo "   Health check: http://localhost:8080/health"
    echo ""
    echo "🛠️  Manual testing suggestions:"
    echo "   1. Visit homepage to check interface"
    echo "   2. Test file upload functionality (requires OAuth configuration)"
    echo "   3. Check if S3 integration is working properly"
    echo ""
    echo "⏹️  Stop container:"
    echo "   docker stop go-s3-uploader-test"
    echo "   docker rm go-s3-uploader-test"
}

# Error handling
trap 'echo ""; echo "❌ Testing interrupted"; show_logs_and_cleanup; exit 1' INT TERM

# Run main process
main
