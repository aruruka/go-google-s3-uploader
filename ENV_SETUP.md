# 🔧 Environment Variable Setup Guide

## Required Environment Variables

### 1. AWS Configuration
```bash
export AWS_REGION="ap-northeast-1"
export S3_BUCKET_NAME="raymond-go-s3-uploader-dev-2025"
export AWS_PROFILE="go-aws-sdk"  # or use AWS_ACCESS_KEY_ID/AWS_SECRET_ACCESS_KEY
```

### 2. Google OAuth Configuration
```bash
export GOOGLE_CLIENT_ID="your_google_client_id"
export GOOGLE_CLIENT_SECRET="your_google_client_secret"
export REDIRECT_URL="http://localhost:8080/auth/callback"
export APP_SERVER_URL="http://localhost:8080"
```

## Quick Setup Methods

### Method 1: Use .env File (Recommended)
```bash
# Copy example file
cp auth-server/.env.example auth-server/.env

# Edit .env file, fill in real credentials
nano auth-server/.env

# Load environment variables
source auth-server/.env
```

### Method 2: Set Environment Variables Directly
```bash
# Set all variables at once
export GOOGLE_CLIENT_ID="your_client_id"
export GOOGLE_CLIENT_SECRET="your_client_secret"
export AWS_REGION="ap-northeast-1"
export S3_BUCKET_NAME="raymond-go-s3-uploader-dev-2025"
export AWS_PROFILE="go-aws-sdk"
```

### Method 3: Use Script Setup
```bash
# Run environment setup script
source ./setup-docker-env.sh
```

## Security Considerations

⚠️ **Never**:
- Hard-code credentials in scripts
- Commit .env files containing credentials to Git
- Share real API keys in public places

✅ **Recommended practices**:
- Use environment variables
- Use AWS Profile
- Use .env.example as template
- Use IAM roles in production environment

## Verify Setup

Run the following command to check if environment variables are set correctly:

```bash
./quick-verify.sh
```

If all checks pass, you can run the complete test:

```bash
./end-to-end-test.sh
```
