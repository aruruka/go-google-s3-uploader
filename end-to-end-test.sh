#!/bin/bash

# End-to-end testing script
# 🎯 Go Google S3 Uploader - Complete workflow testing

echo "🎯 End-to-end testing begins"
echo "============================"

# 设置颜色
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Set environment variables
# ⚠️ OAuth credentials need to be loaded from environment variables or .env file  
if [[ -z "$GOOGLE_CLIENT_ID" ]]; then
    echo "⚠️  GOOGLE_CLIENT_ID not set"
    echo "   Please set environment variables or copy configuration from auth-server/.env.example"
    echo "   export GOOGLE_CLIENT_ID=\"your_client_id\""
fi
if [[ -z "$GOOGLE_CLIENT_SECRET" ]]; then
    echo "⚠️  GOOGLE_CLIENT_SECRET not set"
    echo "   Please set environment variables or copy configuration from auth-server/.env.example"
    echo "   export GOOGLE_CLIENT_SECRET=\"your_client_secret\""
fi

export AUTH_SERVER_URL="http://localhost:8081"
export APP_SERVER_URL="http://localhost:8080"
export AWS_PROFILE="go-aws-sdk"
export AWS_REGION="ap-northeast-1"
export S3_BUCKET_NAME="raymond-go-s3-uploader-dev-2025"

echo -e "${BLUE}📋 Test environment configuration:${NC}"
echo "  Auth Server: $AUTH_SERVER_URL"
echo "  App Server: $APP_SERVER_URL"
echo "  S3 Bucket: $S3_BUCKET_NAME"
echo "  AWS Profile: $AWS_PROFILE"
echo ""

# Clean up old processes
echo -e "${YELLOW}🧹 Cleaning up old processes...${NC}"
pkill -f "auth-server" 2>/dev/null || true
pkill -f "app-server" 2>/dev/null || true
sleep 2

echo -e "${BLUE}🔐 启动 Auth Server (8081)...${NC}"
cd auth-server
go run main.go &
AUTH_PID=$!
cd ..

echo -e "${BLUE}📱 启动 App Server (8080)...${NC}"
cd app-server  
go run main.go &
APP_PID=$!
cd ..

# Wait for servers to start
echo -e "${YELLOW}⏳ Waiting for servers to start...${NC}"
sleep 5

# Check server status
echo -e "${BLUE}🔍 Checking server status:${NC}"
curl -s http://localhost:8081/health > /dev/null && echo "✅ Auth Server (8081) - Running normally" || echo "❌ Auth Server (8081) - Failed to start"
curl -s http://localhost:8080/health > /dev/null && echo "✅ App Server (8080) - Running normally" || echo "❌ App Server (8080) - Failed to start"

echo ""
echo -e "${GREEN}🎉 Servers started! Please follow these steps for manual testing:${NC}"
echo ""
echo -e "${BLUE}📋 Testing steps:${NC}"
echo "1. Open browser and visit: http://localhost:8080"
echo "2. Click 'Upload File' or visit: http://localhost:8080/upload"
echo "3. System will redirect to Google OAuth login page"
echo "4. Log in with your Google account"
echo "5. Authorize the application to access your basic information"
echo "6. After successful login, you'll be redirected back to the upload page"
echo "7. Select an image file to upload"
echo "8. After successful upload, the S3 URL of the file will be displayed"
echo ""
echo -e "${YELLOW}⚠️  Important notes:${NC}"
echo "• Only supports .jpg, .jpeg, .png, .gif image formats"
echo "• File size limit: 10MB"
echo "• Ensure your Google account has access permissions"
echo ""
echo -e "${BLUE}🌐 Quick access links:${NC}"
echo "• App homepage: http://localhost:8080"
echo "• Direct upload: http://localhost:8080/upload"
echo "• Auth login: http://localhost:8081/login"
echo ""
echo -e "${GREEN}After testing is complete, press Ctrl+C to stop all servers${NC}"

# Keep script running and wait for user interruption
trap 'echo -e "\n🛑 Stopping servers..."; kill $AUTH_PID $APP_PID 2>/dev/null; exit 0' INT

# Continuously monitor server status
while true; do
    sleep 30
    if ! ps -p $AUTH_PID > /dev/null 2>&1; then
        echo -e "${YELLOW}⚠️  Auth Server has stopped${NC}"
    fi
    if ! ps -p $APP_PID > /dev/null 2>&1; then
        echo -e "${YELLOW}⚠️  App Server has stopped${NC}"
    fi
done
