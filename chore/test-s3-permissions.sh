#!/bin/bash

echo "🧪 S3 Integration Test - Production Environment Permission Verification"
echo "====================================================================="

# Set environment variables
export S3_BUCKET_NAME="raymond-go-s3-uploader-dev-2025"
export AWS_REGION="ap-northeast-1"

echo "📋 Test configuration:"
echo "   S3 Bucket: $S3_BUCKET_NAME"
echo "   Region: $AWS_REGION"
echo ""

echo "🔐 Test 1: Administrator Permissions (local-dev-admin)"
echo "Purpose: Development environment, full permissions"
export AWS_PROFILE="local-dev-admin"
echo "   Using Profile: $AWS_PROFILE"
(cd app-server && go run -tags=integration test-s3-integration.go)
echo ""

echo "🎯 Test 2: Production Equivalent Permissions (go-aws-sdk)"
echo "Purpose: Simulate App Runner application permissions, principle of least privilege"
export AWS_PROFILE="go-aws-sdk"
echo "   Using Profile: $AWS_PROFILE" 
(cd app-server && go run -tags=integration test-s3-integration.go)
echo ""

echo "✅ Permission verification complete!"
echo "If both tests pass, it means:"
echo "  1. Development environment is configured correctly"
echo "  2. Production environment permissions follow principle of least privilege"
echo "  3. App Runner deployment permissions have been verified"
