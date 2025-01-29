package processrepos

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	config "github.com/s3pweb/gitArchiveS3Report/config"
)

// uploadFile uploads a single file to the specified S3 bucket
func uploadFile(client *s3.Client, bucket, filePath string, uploadkey string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file %q, %v", filePath, err)
	}
	defer file.Close()

	_, err = client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(uploadkey + filePath),
		Body:   file,
	})
	if err != nil {
		return fmt.Errorf("unable to upload file to %q, %v", bucket, err)
	}

	fmt.Printf("Successfully uploaded %q to %q\n", filePath, bucket)
	return nil
}

// Upload uploads all files in the specified directory to the specified S3 bucket
func Upload(dirPath string) error {
	cfg := config.Get()

	// Charger la configuration AWS
	awsConfig, err := awsConfig.LoadDefaultConfig(context.TODO(),
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
		return fmt.Errorf("unable to load SDK config, %v", err)
	}

	// create an S3 client
	client := s3.NewFromConfig(awsConfig)

	// Walk through the directory and upload all files
	err = filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			err = uploadFile(client, cfg.AWS.BucketName, path, cfg.AWS.UploadKey)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to upload files from directory %q, %v", dirPath, err)
	}

	return nil
}
