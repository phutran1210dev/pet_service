package storage

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"pet-service/config"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioClient struct {
	client *minio.Client
	bucket string
}

var minioInstance *MinioClient

func InitMinio() error {
	cfg := config.AppConfig

	client, err := minio.New(cfg.MinioEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.MinioAccessKey, cfg.MinioSecretKey, ""),
		Secure: cfg.MinioUseSSL,
	})
	if err != nil {
		return err
	}

	// Check if bucket exists, create if not
	ctx := context.Background()
	exists, err := client.BucketExists(ctx, cfg.MinioBucket)
	if err != nil {
		return err
	}

	if !exists {
		err = client.MakeBucket(ctx, cfg.MinioBucket, minio.MakeBucketOptions{})
		if err != nil {
			return err
		}
		log.Printf("Bucket %s created successfully", cfg.MinioBucket)
	}

	minioInstance = &MinioClient{
		client: client,
		bucket: cfg.MinioBucket,
	}

	log.Println("MinIO connection established")
	return nil
}

func GetMinioClient() *MinioClient {
	return minioInstance
}

func (m *MinioClient) UploadFile(objectName string, data []byte, contentType string) (string, error) {
	ctx := context.Background()

	reader := bytes.NewReader(data)
	size := int64(len(data))

	_, err := m.client.PutObject(ctx, m.bucket, objectName, reader, size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", err
	}

	// Return the URL
	url := fmt.Sprintf("%s/%s/%s", config.AppConfig.MinioEndpoint, m.bucket, objectName)
	if config.AppConfig.MinioUseSSL {
		url = fmt.Sprintf("https://%s/%s/%s", config.AppConfig.MinioEndpoint, m.bucket, objectName)
	} else {
		url = fmt.Sprintf("http://%s/%s/%s", config.AppConfig.MinioEndpoint, m.bucket, objectName)
	}

	return url, nil
}

func (m *MinioClient) DownloadFile(objectName string) ([]byte, error) {
	ctx := context.Background()

	object, err := m.client.GetObject(ctx, m.bucket, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	defer object.Close()

	data, err := io.ReadAll(object)
	if err != nil {
		return nil, err
	}

	return data, nil
}
