# 🚀 AWS App Runner Automated Deployment Guide

## Project Overview
This is a Go language file upload application that integrates Google OAuth and AWS S3 storage, using GitHub Actions for automated deployment to AWS App Runner.

## 🏗️ Architecture Design

### Deployment Architecture
```
GitHub Repository → GitHub Actions → AWS App Runner → S3 Storage
     ↓                    ↓               ↓            ↓
Source Management     CI/CD Process   App Hosting   File Storage
```

### Cost Optimization
- **App Runner**: ~$2-5/month (0.25 vCPU, 0.5GB RAM)
- **S3 Storage**: ~$0.023/GB/month
- **Total Cost**: < $10/month (suitable for demos and small projects)

## 🔧 Completed Infrastructure

### AWS Resources
1. **S3 Bucket**: `raymond-go-s3-uploader-dev-2025`
   - Region: ap-northeast-1
   - Encryption: AES256
   - Access Control: Private

2. **IAM Roles**:
   - `GitHubActions-AppRunner-Role`: GitHub Actions OIDC authentication
   - `AppRunnerInstanceRole`: Application runtime S3 access permissions

3. **IAM Policies**:
   - `AppRunnerS3WritePolicy`: S3 read/write permissions
   - OIDC Web Identity trust relationship

### Application Configuration
1. **Docker Multi-stage Build**: Optimize image size
2. **Health Check Endpoint**: `/health` (App Runner requirement)
3. **Environment Variables**: Support dynamic port configuration
4. **Static File Service**: Integrated frontend resources

## 📋 Deployment Steps

### 1. GitHub Repository Setup
```bash
# Initialize Git repository
git init
git add .
git commit -m "Initial commit: OAuth + S3 integration"

# Add remote repository
git remote add origin https://github.com/your-username/go-google-s3-uploader.git
git push -u origin main
```

### 2. GitHub Secrets Configuration
Add the following Secrets in GitHub repository settings:

```
AWS_ROLE_ARN=arn:aws:iam::925867211284:role/GitHubActions-AppRunner-Role
```

### 3. Automated Deployment Trigger
- **Push to main branch**: Automatically triggers deployment
- **Pull Request**: Only runs tests
- **Manual Trigger**: Through GitHub Actions interface

## 🔄 CI/CD Process

### GitHub Actions Workflow
1. **Testing Phase**:
   - Go code compilation check
   - Unit test execution
   - Code quality verification

2. **Deployment Phase**:
   - OIDC identity authentication
   - Docker image build
   - App Runner service creation/update
   - Health check verification

### Deployment Strategy
- **Zero Downtime**: App Runner rolling updates
- **Automatic Rollback**: Rollback on health check failure
- **Monitoring Integration**: CloudWatch logs and metrics

## 📂 Project Structure

```
go-google-s3-uploader/
├── .github/workflows/
│   └── deploy.yml              # GitHub Actions workflow
├── app-server/                 # Go application server
│   ├── main.go                # Main program (supports dynamic port)
│   └── pkg/                   # Business logic packages
├── auth-server/               # OAuth authentication server
├── shared/                    # Shared resources
│   └── static/               # Frontend static files
├── Dockerfile                 # Multi-stage build configuration
├── apprunner.yaml            # App Runner configuration
├── s3-bucket-policy.json     # S3 access policy
└── s3-write-policy.json      # IAM policy document
```

## 🔒 Security Configuration

### IAM Permission Minimization
- **GitHub Actions**: Can only create/update App Runner services
- **App Runner**: Can only access specified S3 bucket
- **No Hard-coded Credentials**: Use OIDC and IAM roles

### HTTPS Enforcement
- App Runner automatically provides SSL/TLS
- All communications encrypted in transit
- OAuth callback URLs use HTTPS

## 🎯 Future Development Plans

### Phase 3: S3 Integration (Current)
- [x] AWS infrastructure setup
- [x] CI/CD pipeline configuration
- [ ] S3 file upload functionality implementation
- [ ] File type validation and security checks
- [ ] Upload progress display

### Phase 4: Domain Configuration
- [ ] CloudFlare DNS configuration
- [ ] `app.shouneng.website` domain mapping
- [ ] SSL certificate management

### Phase 5: Monitoring Operations
- [ ] CloudWatch log aggregation
- [ ] Performance metrics monitoring
- [ ] Error alerting configuration

## 🚨 Troubleshooting

### Common Issues
1. **Deployment Failure**: Check GitHub Secrets configuration
2. **Permission Error**: Verify IAM roles and policies
3. **Health Check Failure**: Confirm `/health` endpoint response
4. **Static File 404**: Check file paths and Docker COPY instructions

### Debug Commands
```bash
# Local testing
docker build -t go-s3-uploader .
docker run -p 8080:8080 go-s3-uploader

# AWS service status check
aws apprunner describe-service --service-arn <SERVICE_ARN>
aws s3 ls s3://raymond-go-s3-uploader-dev-2025/
```

## 📊 Performance Metrics

### Expected Metrics
- **Cold Start Time**: < 10 seconds
- **Response Time**: < 500ms
- **Availability**: > 99.5%
- **Concurrent Processing**: 10+ users

### Monitoring Metrics
- HTTP request count and latency
- S3 operation success rate
- Error rate and exception logs
- Resource utilization (CPU/Memory)

---
## 📞 Support Contact

**Technology Stack**: Go + AWS App Runner + GitHub Actions + OAuth + S3
**Cost Budget**: < $10/month
**Deployment Mode**: Fully automated Zero-Touch deployment
✅ GitHub repository configured with AWS deployment secrets
