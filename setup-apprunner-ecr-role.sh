#!/bin/bash

# AWS App Runner ECR Access Role Setup Script
# This script creates the necessary IAM role for App Runner to access private ECR

set -e

# Configuration
ACCOUNT_ID="925867211284"
REGION="ap-northeast-1"
ECR_REPOSITORY="go-s3-uploader"
ACCESS_ROLE_NAME="AppRunnerECRAccessRole"

echo "🚀 Setting up App Runner ECR Access Role..."

# Step 1: Create the access role for ECR
echo "📝 Creating App Runner ECR access role..."
aws iam create-role \
    --role-name $ACCESS_ROLE_NAME \
    --assume-role-policy-document file://apprunner-ecr-access-role.json

# Step 2: Attach the AWSAppRunnerServicePolicyForECRAccess managed policy
echo "🔗 Attaching ECR access policy..."
aws iam attach-role-policy \
    --role-name $ACCESS_ROLE_NAME \
    --policy-arn arn:aws:iam::aws:policy/service-role/AWSAppRunnerServicePolicyForECRAccess

# Step 3: Get and display the role ARN
ROLE_ARN=$(aws iam get-role --role-name $ACCESS_ROLE_NAME --query 'Role.Arn' --output text)

echo "✅ App Runner ECR access role created successfully!"
echo "📋 Role ARN: $ROLE_ARN"
echo ""
echo "🔧 Next steps:"
echo "1. Add this role ARN to your GitHub Actions secrets as 'APP_RUNNER_ACCESS_ROLE_ARN'"
echo "2. Update your workflow to use this role in the App Runner service configuration"
echo ""
echo "GitHub Secret command:"
echo "gh secret set APP_RUNNER_ACCESS_ROLE_ARN --body \"$ROLE_ARN\""
