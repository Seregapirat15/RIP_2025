package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"lab2/pkg/config"
)

// MinIOTelescopeImagesClient представляет клиент MinIO для хранения изображений телескопов
var MinIOTelescopeImagesClient *minio.Client

// InitMinIOTelescopeImagesConnection инициализирует подключение к MinIO
func InitMinIOTelescopeImagesConnection(cfg *config.ExoplanetCalculationConfig) error {
	var err error
	MinIOTelescopeImagesClient, err = minio.New(cfg.MinIOEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.MinIOAccessKeyID, cfg.MinIOSecretAccessKey, ""),
		Secure: cfg.MinIOUseSSL,
	})
	if err != nil {
		return fmt.Errorf("ошибка инициализации MinIO клиента: %w", err)
	}

	// Создание bucket если не существует
	ctx := context.Background()
	exists, err := MinIOTelescopeImagesClient.BucketExists(ctx, cfg.MinIOBucketName)
	if err != nil {
		return fmt.Errorf("ошибка проверки существования bucket: %w", err)
	}

	if !exists {
		err = MinIOTelescopeImagesClient.MakeBucket(ctx, cfg.MinIOBucketName, minio.MakeBucketOptions{})
		if err != nil {
			return fmt.Errorf("ошибка создания bucket: %w", err)
		}
		log.Printf("Bucket '%s' создан успешно", cfg.MinIOBucketName)
	}

	log.Println("Подключение к MinIO установлено")
	return nil
}

// GetTelescopeImagePresignedURL генерирует presigned URL для изображения телескопа
func GetTelescopeImagePresignedURL(objectName string) (string, error) {
	if objectName == "" {
		return "", nil
	}

	ctx := context.Background()
	
	// Создаем новый клиент с внешним адресом для генерации URL
	externalClient, err := minio.New("localhost:9000", &minio.Options{
		Creds:  credentials.NewStaticV4("minioadmin123", "minioadmin123456", ""),
		Secure: false,
		Region: "us-east-1",
	})
	if err != nil {
		return "", err
	}
	
	presignedURL, err := externalClient.PresignedGetObject(ctx, "telescope-images", objectName, time.Hour*24, nil)
	if err != nil {
		return "", err
	}

	return presignedURL.String(), nil
}
