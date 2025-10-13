package service

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"os"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var Client *minio.Client

func MinioInit() error {
	endpoint := os.Getenv("MINIO_ENDPOINT")
	accessKey := os.Getenv("MINIO_ACCESS_KEY")
	secretKey := os.Getenv("MINIO_SECRET_KEY")
	bucket := os.Getenv("MINIO_BUCKET")

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: false,
	})
	if err != nil {
		return err
	}

	Client = client

	ctx := context.Background()
	exists, err := Client.BucketExists(ctx, bucket)
	if err != nil {
		return err
	}

	if !exists {
		if err := Client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{}); err != nil {
			return err
		}
		fmt.Println("Bucket created:", bucket)
	}

	return nil
}
func MinioUpload(filePath string, file io.Reader, size int64, contentType string) (string, error) {
	bucket := os.Getenv("MINIO_BUCKET")
	ctx := context.Background()
	opts := minio.PutObjectOptions{
		ContentType: contentType,
	}
	if os.Getenv("MINIO_PUBLIC") == "true" {
		if opts.UserMetadata == nil {
			opts.UserMetadata = make(map[string]string)
		}
		opts.UserMetadata["x-amz-acl"] = "public-read"
	}
	_, err := Client.PutObject(ctx, bucket, filePath, file, size, opts)
	if err != nil {
		return "", err
	}
	return filePath, nil
}
func MinioRemove(filePath string) error {
	if Client == nil {
		return fmt.Errorf("minio client not initialized")
	}
	bucket := os.Getenv("MINIO_BUCKET")
	ctx := context.Background()
	err := Client.RemoveObject(ctx, bucket, filePath, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to remove file %s: %w", filePath, err)
	}
	fmt.Println("File removed from MinIO:", filePath)
	return nil
}
func MinioGetURL(filePath string, duration time.Duration) (string, error) {
	if os.Getenv("MINIO_PUBLIC") == "true" {
		return fmt.Sprintf("http://%s/%s/%s", os.Getenv("MINIO_ENDPOINT"), os.Getenv("MINIO_BUCKET"), filePath), nil
	}
	ctx := context.Background()
	bucket := os.Getenv("MINIO_BUCKET")
	reqParams := make(url.Values)
	presignedURL, err := Client.PresignedGetObject(ctx, bucket, filePath, duration, reqParams)
	if err != nil {
		return "", err
	}
	return presignedURL.String(), nil
}
