# Google S3 Uploader Project

## Project Overview

This is an image upload application built with Go, demonstrating typical architecture patterns of modern cloud-native applications:

- **Microservice Architecture**: Separating authentication service and application service
- **OAuth 2.0 Authentication**: Using Google as identity provider
- **Cloud Storage Integration**: Uploading images to AWS S3
- **Standard Library First**: Using Go official and community standard libraries

## Architecture Design

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   User Browser  │    │   App Server    │    │   Auth Server   │
│     :8080       │    │     :8080       │    │     :8081       │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         │                       │                       │
         │ 1. Upload Image       │                       │
         ├──────────────────────>│                       │
         │                       │ 2. Validate Token    │
         │                       ├──────────────────────>│
         │                       │ 3. User Info          │
         │                       │<──────────────────────┤
         │                       │                       │
         │                       │ 4. Upload to S3       │
         │                       ├──────────────────────>│
         │ 5. Success Response   │                    ┌───┴───┐
         │<──────────────────────┤                    │AWS S3 │
         │                       │                    └───────┘
```

## Service Description

### Auth Server (:8081)
- Responsible for Google OAuth 2.0 authentication flow
- Uses `golang.org/x/oauth2` and `coreos/go-oidc`
- Provides token validation service to App Server

### App Server (:8080)
- Responsible for image upload functionality
- Communicates with Auth Server to validate user identity
- Uses AWS SDK v2 to upload files to S3

## Development Phases

### Phase 1: Skeleton Setup ✅
- [x] Create project directory structure
- [x] Create basic web server
- [x] Define routes and placeholder handlers

### Phase 2: Authentication Service 🚧
- [x] Implement Google OAuth 2.0 flow
- [x] Integrate `golang.org/x/oauth2`
- [x] Integrate `coreos/go-oidc`
- [ ] Token validation API --> skipped, I am using Cookie-based authentication instead, this project is mainly for DevOps.

### Phase 3: File Upload ✅
- [x] Inter-service authentication validation (cookie-based)
- [x] AWS S3 integration
- [x] File upload handling

### Phase 4: Optimization and Deployment ✅
- [ ] Error handling --> skipped, this is a toy project focused on DevOps demonstration
- [ ] Logging --> skipped, basic logging already sufficient for demo
- [ ] Configuration management --> skipped, current env-based config sufficient
- [x] Dockerization

## Quick Start

1. Start Auth Server:
```bash
cd auth-server
go run main.go
```

2. Start App Server:
```bash
cd app-server  
go run main.go
```

3. Access the application:
- App Server: http://localhost:8080
- Auth Server: http://localhost:8081

## Technology Stack

- **Language**: Go 1.21+
- **Authentication**: OAuth 2.0 / OIDC
- **Cloud Service**: AWS S3
- **Architecture**: Microservices
- **Standard Libraries**:
  - `golang.org/x/oauth2`
  - `github.com/coreos/go-oidc/v3`
  - `github.com/aws/aws-sdk-go-v2`

## End-to-End DevOps Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│  Local Dev      │    │  GitHub Actions │    │  AWS App Runner │
│  Environment    │    │  CI/CD Pipeline │    │  Production     │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         │ git push              │                       │
         ├──────────────────────>│                       │
         │                       │ 1. Build Docker       │
         │                       │ 2. Push to ECR        │
         │                       │ 3. Deploy to App Runner│
         │                       ├──────────────────────>│
         │                       │                       │
         │                       │                       │
    ┌────┴────┐              ┌───┴───┐              ┌────┴────┐
    │Auth:8081│              │ Docker│              │Combined │
    │App :8080│              │ Image │              │Service  │
    └─────────┘              └───────┘              └─────────┘
```

### 🌟 **Code Flexibility Highlight**

**Same Codebase, Multiple Deployment Models:**

- **Local Development**: Run as separate microservices (auth-server:8081, app-server:8080)
- **App Runner Production**: Run as combined service (single container, all routes)
- **Smart Environment Detection**: Code automatically adapts based on environment variables
- **Zero Code Changes**: Deploy to production without modifying application logic

## Deployment Challenges & Solutions

### 1. OAuth Redirect URL Issue ❌➡️✅
**Problem**: Initial deployment redirected to `http://localhost:8081/auth/callback`
- **Root Cause**: Environment not set to production
- **Solution**: Added `ENV=production` in GitHub Actions workflow
- **Fix**: Environment-aware URL configuration

### 2. "Too Many Redirects" Loop ❌➡️✅
**Problem**: Redirect loop between `/` and `/login` on App Runner
- **Root Cause**: App Runner single-container limitation vs microservice architecture
- **Solution**: Smart redirect detection for same-domain scenarios
- **Implementation**: Combined service with internal routing logic

### 3. S3 Permission Denied ❌➡️✅
**Problem**: Application couldn't upload files to S3
- **Root Cause**: App Runner instance role not configured for S3 access
- **Solution**: Created separate IAM roles (ECR access + S3 instance role)
- **AWS Config**: Proper trust policies for `tasks.apprunner.amazonaws.com`

### 4. Missing Environment Variables ❌➡️✅
**Problem**: Authentication flow broken due to missing `AUTH_SERVER_URL`
- **Root Cause**: App Runner environment variables incomplete
- **Solution**: Added all required environment variables via GitHub Actions
- **Variables**: `AUTH_SERVER_URL`, `APP_SERVER_URL`, `REDIRECT_URL`

## Architecture Evolution

### Local Development (Microservices)
```
┌─────────────────┐    ┌─────────────────┐
│   Auth Server   │    │   App Server    │
│   Port 8081     │    │   Port 8080     │
│   ┌─────────┐   │    │   ┌─────────┐   │
│   │OAuth    │   │    │   │S3 Upload│   │
│   │Cookie   │   │    │   │Validate │   │
│   └─────────┘   │    │   └─────────┘   │
└─────────────────┘    └─────────────────┘
```

### App Runner Deployment (Combined)
```
┌─────────────────────────────────────────────┐
│         Combined Service (Port 8080)        │
│   ┌─────────┐           ┌─────────┐        │
│   │OAuth    │           │S3 Upload│        │
│   │Cookie   │  +        │Validate │        │
│   │Routes   │           │Routes   │        │
│   └─────────┘           └─────────┘        │
└─────────────────────────────────────────────┘
```

# Deployment Test
# Test deployment with updated IAM permissions

Permission Update Details:
- Added `iam:CreateServiceLinkedRole` permission
- Allows App Runner to automatically create the necessary service-linked role when the first service is created in the account
- Fixed previous access denied error

# ECR Integration Test
# Testing Docker build and push to private ECR registry

ECR Permission Update:
- Added ECR push permissions: `ecr:PutImage`, `ecr:InitiateLayerUpload`, `ecr:UploadLayerPart`, `ecr:CompleteLayerUpload`
- Enables GitHub Actions to build and push custom Go application Docker images
- Replaces hello-app-runner with actual application deployment

Service Recreation:
- Deleted existing ECR_PUBLIC service to avoid repository type conflict
- Will create new service with private ECR image from scratch

You should now be able to successfully create App Runner services.
