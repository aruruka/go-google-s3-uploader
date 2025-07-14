#!/bin/bash

# 🔧 Docker Test Environment Variable Configuration Script
# Run this script before running test-docker-local.sh to set necessary environment variables

echo "🔧 Setting up Docker test environment variables..."
echo "================================================="

# AWS basic configuration
export AWS_REGION="ap-northeast-1"
export S3_BUCKET_NAME="raymond-go-s3-uploader-dev-2025"

# Check AWS credentials configuration
echo "📋 Checking AWS credentials configuration..."

# Option 1: Use AWS Profile (recommended for local development)
if [[ -f "$HOME/.aws/credentials" ]] && grep -q "\[go-aws-sdk\]" "$HOME/.aws/credentials"; then
    export AWS_PROFILE="go-aws-sdk"
    echo "✅ Using AWS Profile: go-aws-sdk"
    echo "   Configuration file: $HOME/.aws/credentials"
elif [[ -f "$HOME/.aws/credentials" ]] && grep -q "\[local-dev-admin\]" "$HOME/.aws/credentials"; then
    export AWS_PROFILE="local-dev-admin"
    echo "✅ Using AWS Profile: local-dev-admin"
    echo "   Configuration file: $HOME/.aws/credentials"
else
    echo "⚠️  AWS Profile configuration not found"
    echo "   Please set environment variables or configure AWS Profile:"
    echo "   aws configure --profile go-aws-sdk"
    echo ""
    echo "   Or set environment variables:"
    echo "   export AWS_ACCESS_KEY_ID=your_access_key"
    echo "   export AWS_SECRET_ACCESS_KEY=your_secret_key"
fi

# OAuth configuration (optional, for full functionality testing)
echo ""
echo "🔑 OAuth configuration (optional)..."

if [[ -f "OAuth-Credentials/client_secret.json" ]]; then
    echo "✅ OAuth credentials file found"
    # Here we could extract credentials from file, but for security reasons, let user set manually
    echo "   Please manually set OAuth environment variables:"
    echo "   export GOOGLE_CLIENT_ID=\"your_client_id\""
    echo "   export GOOGLE_CLIENT_SECRET=\"your_client_secret\""
    echo "   export REDIRECT_URL=\"http://localhost:8080/auth/callback\""
elif [[ -n "$GOOGLE_CLIENT_ID" ]]; then
    echo "✅ OAuth environment variables already set"
    export REDIRECT_URL="http://localhost:8080/auth/callback"
    export APP_SERVER_URL="http://localhost:8080"
else
    echo "⚠️  OAuth credentials not set"
    echo "   Application will not provide full authentication functionality"
    echo "   But basic container functionality testing can still proceed"
fi

echo ""
echo "📊 Current environment variable configuration:"
echo "   AWS_REGION: $AWS_REGION"
echo "   S3_BUCKET_NAME: $S3_BUCKET_NAME"
echo "   AWS_PROFILE: ${AWS_PROFILE:-Not set}"
echo "   GOOGLE_CLIENT_ID: ${GOOGLE_CLIENT_ID:+Set}"
echo "   GOOGLE_CLIENT_SECRET: ${GOOGLE_CLIENT_SECRET:+Set}"

echo ""
echo "🚀 Ready! Now you can run Docker tests:"
echo "   ./test-docker-local.sh"
echo ""
echo "💡 Tips:"
echo "   - Container startup testing only: No OAuth configuration needed"
echo "   - File upload testing: OAuth configuration required"
echo "   - S3 integration testing: AWS credentials required"
