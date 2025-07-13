# 🛡️ AWS IAM and S3 Permission Configuration Guide

## 📋 Overview

This project requires configuring AWS IAM user and S3 bucket permissions. For security reasons, real configuration files are not included in the repository. Please create your own configuration according to the following guide.

## 🔧 Required Configuration Files

### 1. S3 Bucket Policy (`s3-bucket-policy.json`)

```bash
# Copy template file
cp s3-bucket-policy.json.example s3-bucket-policy.json

# Edit file, replace the following placeholders:
# - YOUR_AWS_ACCOUNT_ID: Your AWS Account ID
# - YOUR_DEV_USER: Development environment IAM username  
# - YOUR_APP_USER: Application IAM username
# - YOUR_S3_BUCKET_NAME: Your S3 bucket name
```

### 2. S3 Write Policy (`s3-write-policy.json`)

```bash
# Copy template file
cp s3-write-policy.json.example s3-write-policy.json

# Edit file, replace:
# - YOUR_S3_BUCKET_NAME: Your S3 bucket name
```

## 🚀 Deployment Steps

### Step 1: Create S3 Bucket

```bash
# Use AWS CLI to create bucket
aws s3 mb s3://your-bucket-name --region ap-northeast-1

# Apply bucket policy
aws s3api put-bucket-policy --bucket your-bucket-name --policy file://s3-bucket-policy.json
```

### Step 2: Create IAM User and Policy

```bash
# Create IAM policy
aws iam create-policy --policy-name S3UploaderWritePolicy --policy-document file://s3-write-policy.json

# Create IAM user (for application)
aws iam create-user --user-name your-app-user

# Attach policy to user
aws iam attach-user-policy --user-name your-app-user --policy-arn arn:aws:iam::YOUR_ACCOUNT_ID:policy/S3UploaderWritePolicy

# Create access key
aws iam create-access-key --user-name your-app-user
```

### Step 3: Configure App Runner Service Role

```bash
# Create service role (for App Runner)
aws iam create-role --role-name AppRunnerInstanceRole --assume-role-policy-document '{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Service": "tasks.apprunner.amazonaws.com"
      },
      "Action": "sts:AssumeRole"
    }
  ]
}'

# Attach S3 policy to service role
aws iam attach-role-policy --role-name AppRunnerInstanceRole --policy-arn arn:aws:iam::YOUR_ACCOUNT_ID:policy/S3UploaderWritePolicy
```

## 📚 Permission Explanation

### S3 Bucket Policy Permissions
- `s3:PutObject` - Upload files
- `s3:PutObjectAcl` - Set file permissions
- `s3:GetObject` - Read files
- `s3:DeleteObject` - Delete files
- `s3:ListBucket` - List bucket contents
- `s3:GetObjectVersion` - Get file versions

### Principle of Least Privilege
Configuration follows the principle of least privilege, granting only the minimum permissions required for application operation.

## ⚠️ Security Considerations

1. **Don't commit real configurations**: Real AWS configuration files are excluded by `.gitignore`
2. **Rotate keys regularly**: Regularly update IAM user access keys
3. **Monitor usage**: Use AWS CloudTrail to monitor S3 access
4. **Environment isolation**: Use different buckets and policies for different environments (dev/test/prod)

## 🔧 Environment Variable Configuration

After configuration is complete, set the following environment variables:

```bash
export AWS_REGION="ap-northeast-1"
export S3_BUCKET_NAME="your-bucket-name"
export AWS_ACCESS_KEY_ID="your-access-key"
export AWS_SECRET_ACCESS_KEY="your-secret-key"
```

Or use AWS Profile:

```bash
export AWS_PROFILE="your-profile-name"
```
