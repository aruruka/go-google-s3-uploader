//go:build integration
// +build integration

// Integration test for S3 functionality
// This file tests real AWS S3 operations with actual credentials
// Use: go run -tags=integration test-s3-integration.go
//
// Purpose:
// 1. Validate S3 integration works with real AWS credentials
// 2. Test permissions of IAM roles/users
// 3. Verify bucket access and file operations
// 4. Integration testing before deployment

package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/aruruka/go-google-s3-uploader/app-server/pkg/s3"
)

func main() {
	currentProfile := os.Getenv("AWS_PROFILE")
	if currentProfile == "" {
		fmt.Println("❌ AWS_PROFILE environment variable not set")
		os.Exit(1)
	}

	fmt.Printf("🧪 S3 Integration Test - Profile: %s\n", currentProfile)
	fmt.Println("=" + strings.Repeat("=", 50))

	testS3WithCurrentProfile()
}

func testS3WithCurrentProfile() {
	currentProfile := os.Getenv("AWS_PROFILE")

	// 验证必需的环境变量
	bucketName := os.Getenv("S3_BUCKET_NAME")
	region := os.Getenv("AWS_REGION")

	if bucketName == "" || region == "" {
		fmt.Println("❌ Required environment variables not set:")
		fmt.Println("   S3_BUCKET_NAME:", bucketName)
		fmt.Println("   AWS_REGION:", region)
		os.Exit(1)
	}

	fmt.Printf("📋 配置信息:\n")
	fmt.Printf("   AWS Profile: %s\n", currentProfile)
	fmt.Printf("   S3 Bucket: %s\n", bucketName)
	fmt.Printf("   Region: %s\n", region)
	fmt.Println()

	// 创建S3客户端
	s3Client, err := s3.NewS3Client()
	if err != nil {
		fmt.Printf("❌ Failed to create S3 client: %v\n", err)
		os.Exit(1)
	}

	// 测试文件内容
	testContent := fmt.Sprintf("Test file from profile: %s\nUploaded at: %s",
		currentProfile, time.Now().Format("2006-01-02 15:04:05"))
	testKey := fmt.Sprintf("integration-test/%s/test-file.txt", currentProfile)

	fmt.Printf("📤 Uploading test file: %s\n", testKey)

	// 上传测试文件
	ctx := context.Background()
	err = s3Client.UploadFile(ctx, testKey, strings.NewReader(testContent), "text/plain")
	if err != nil {
		fmt.Printf("❌ Upload failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✅ Upload successful with profile %s\n", currentProfile)

	// 获取文件URL
	fileURL := s3Client.GetFileURL(testKey)
	fmt.Printf("🌐 File URL: %s\n", fileURL)

	// 测试列出文件（权限验证）
	files, err := s3Client.ListFiles(ctx, "integration-test/")
	if err != nil {
		fmt.Printf("❌ List files failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("📂 Found %d files in integration-test/\n", len(files))
	for _, file := range files {
		if strings.Contains(file, currentProfile) {
			fmt.Printf("   ✅ %s\n", file)
		}
	}

	fmt.Printf("\n🎉 Profile %s - All tests passed!\n", currentProfile)
}
