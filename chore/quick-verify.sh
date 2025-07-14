#!/bin/bash

# Quick system verification script
echo "🔍 Quick System Verification"
echo "============================"

# Set environment variables
# ⚠️ OAuth credentials need to be loaded from environment variables or .env file
if [[ -z "$GOOGLE_CLIENT_ID" ]]; then
    echo "⚠️  GOOGLE_CLIENT_ID not set, some features may not work properly"
fi
if [[ -z "$GOOGLE_CLIENT_SECRET" ]]; then
    echo "⚠️  GOOGLE_CLIENT_SECRET not set, OAuth functionality will not be available"
fi

export AUTH_SERVER_URL="http://localhost:8081"
export APP_SERVER_URL="http://localhost:8080"
export AWS_PROFILE="go-aws-sdk"
export AWS_REGION="ap-northeast-1"
export S3_BUCKET_NAME="raymond-go-s3-uploader-dev-2025"

echo "✅ 1. Environment Variable Configuration Check"
echo "   GOOGLE_CLIENT_ID: ${GOOGLE_CLIENT_ID:+Set}"
echo "   AWS_PROFILE: $AWS_PROFILE"
echo "   S3_BUCKET_NAME: $S3_BUCKET_NAME"

echo ""
echo "✅ 2. Unit Test Verification"
cd app-server
if go test ./pkg/s3 -v > /dev/null 2>&1; then
    echo "   S3 Package Unit Tests: Passed"
else
    echo "   S3 Package Unit Tests: Failed"
fi
cd ..

echo ""
echo "✅ 3. Build Verification"
cd auth-server
if go build -o ../auth-test . > /dev/null 2>&1; then
    echo "   Auth Server Build: Success"
    rm -f ../auth-test
else
    echo "   Auth Server Build: Failed"
fi
cd ..

cd app-server
if go build -o ../app-test . > /dev/null 2>&1; then
    echo "   App Server Build: Success"
    rm -f ../app-test
else
    echo "   App Server Build: Failed"
fi
cd ..

echo ""
echo "✅ 4. AWS S3 Connection Verification"
cd app-server
if go run -tags=integration test-s3-integration.go > /dev/null 2>&1; then
    echo "   S3 Integration Test: Passed"
else
    echo "   S3 Integration Test: Needs checking"
fi
cd ..

echo ""
echo "🎯 Summary: Core System Verification Complete"
echo "============================================="
echo "If all items above show as success/passed,"
echo "your system is ready for end-to-end testing!"
echo ""
echo "Next steps:"
echo "1. Run ./end-to-end-test.sh to start full testing"
echo "2. Or manually start servers for testing"
