package processrepos

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/s3pweb/gitArchiveS3Report/config"
	"github.com/s3pweb/gitArchiveS3Report/utils/logger"
)

// uploadFile uploads a single file to the specified S3 bucket
func uploadFile(client *s3.Client, bucket, filePath, uploadKey string, logger *logger.Logger) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file %q: %v", filePath, err)
	}
	defer file.Close()

	// Get file size for logging
	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("failed to get file info: %v", err)
	}

	// Extract just the filename from the full path
	fileName := filepath.Base(filePath)

	// Construct the S3 key by joining the upload prefix with the filename
	s3Key := filepath.Join(uploadKey, fileName)
	// Replace Windows path separators with forward slashes for S3
	s3Key = filepath.ToSlash(s3Key)

	logger.Info("Starting upload of %s (%.2f MB) to s3://%s/%s",
		fileName,
		float64(fileInfo.Size())/1024/1024,
		bucket,
		s3Key)

	startTime := time.Now()

	_, err = client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(s3Key),
		Body:   file,
	})
	if err != nil {
		return fmt.Errorf("failed to upload file to s3://%s/%s: %v", bucket, s3Key, err)
	}

	duration := time.Since(startTime).Round(time.Millisecond)
	logger.Info("Successfully uploaded %s in %v", fileName, duration)

	return nil
}

// Upload uploads files to the specified S3 bucket
// If path is a directory, it uploads all files in the directory
// If path is a file, it uploads just that file
func Upload(path string) error {
	cfg := config.Get()

	logger, err := logger.NewLogger("S3Upload", cfg.Logger.Level)
	if err != nil {
		return fmt.Errorf("failed to initialize logger: %v", err)
	}

	if path == "" {
		return fmt.Errorf("path cannot be empty")
	}

	// Verify path exists
	fileInfo, err := os.Stat(path)
	if os.IsNotExist(err) {
		return fmt.Errorf("path %s does not exist", path)
	}
	if err != nil {
		return fmt.Errorf("error accessing path %s: %v", path, err)
	}

	// Determine if we're uploading a single file or a directory
	isSingleFile := !fileInfo.IsDir()

	logger.Info("Starting S3 upload process from: %s", path)
	startTime := time.Now()

	// Load AWS configuration
	awsCfg, err := awsConfig.LoadDefaultConfig(context.TODO(),
		awsConfig.WithRegion(cfg.AWS.Region),
		awsConfig.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				cfg.AWS.AccessKeyID,
				cfg.AWS.SecretAccessKey,
				"",
			),
		),
	)
	if err != nil {
		return fmt.Errorf("failed to load AWS config: %v", err)
	}

	// Create S3 client
	client := s3.NewFromConfig(awsCfg)

	// Verify bucket exists and is accessible
	_, err = client.HeadBucket(context.TODO(), &s3.HeadBucketInput{
		Bucket: aws.String(cfg.AWS.BucketName),
	})
	if err != nil {
		return fmt.Errorf("failed to access bucket %s: %v", cfg.AWS.BucketName, err)
	}

	logger.Info("Successfully connected to S3 bucket: %s", cfg.AWS.BucketName)

	// If it's a single file, upload it directly
	if isSingleFile {
		return uploadFile(client, cfg.AWS.BucketName, path, cfg.AWS.AWSUploadPath, logger)
	}

	// Otherwise, it's a directory - count files to upload
	var filesToUpload int
	err = filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			filesToUpload++
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to count files in directory: %v", err)
	}

	logger.Info("Found %d files to upload", filesToUpload)

	// Upload files
	filesUploaded := 0
	err = filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			if err := uploadFile(client, cfg.AWS.BucketName, filePath, cfg.AWS.AWSUploadPath, logger); err != nil {
				return err
			}
			filesUploaded++
			logger.Info("Progress: %d/%d files uploaded (%.1f%%)",
				filesUploaded,
				filesToUpload,
				float64(filesUploaded)/float64(filesToUpload)*100)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("upload process failed: %v", err)
	}

	duration := time.Since(startTime).Round(time.Second)
	logger.Info("Upload process completed successfully in %v", duration)
	logger.Info("Total files uploaded: %d", filesUploaded)

	return nil
}
