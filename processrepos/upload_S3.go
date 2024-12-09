package processrepos

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/s3pweb/gitArchiveS3Report/utils"
)

// uploadFile uploads a single file to the specified S3 bucket
func uploadFile(client *s3.Client, bucket, filePath string, uploadkey string) error {
	// Open the local file to upload
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file %q, %v", filePath, err)
	}
	defer file.Close()

	// Upload the file to S3
	_, err = client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(uploadkey + filePath), // 'key' represents the path in S3 where the file will be stored
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
	// Load environment variables from the .secrets file
	secrets, err := utils.ReadConfigFile(".secrets")
	if err != nil {
		return fmt.Errorf("unable to load .secrets file, %v", err)
	}

	// Retrieve values from the .secrets file
	region := secrets["AWS_REGION"]
	bucket := secrets["AWS_BUCKET_NAME"]
	awsKey := secrets["AWS_ACCESS_KEY_ID"]
	awsSecret := secrets["AWS_SECRET_ACCESS_KEY"]
	uploadkey := secrets["UPLOAD_KEY"]

	// Load AWS configuration using the credentials directly
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(awsKey, awsSecret, ""),
		),
	)
	if err != nil {
		return fmt.Errorf("unable to load SDK config, %v", err)
	}

	// Create an S3 client
	client := s3.NewFromConfig(cfg)

	// Recursively walk through the directory and upload each file
	err = filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			err = uploadFile(client, bucket, path, uploadkey)
			if err != nil {
				return err
			}
			return nil
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to upload files from directory %q, %v", dirPath, err)
	}

	return nil
}
