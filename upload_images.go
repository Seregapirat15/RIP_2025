package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func main() {
	// Конфигурация MinIO
	endpoint := "localhost:9000"
	accessKeyID := "minioadmin123"
	secretAccessKey := "minioadmin123456"
	useSSL := false
	bucketName := "telescope-images"

	// Создание клиента MinIO
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatal("Ошибка создания клиента MinIO:", err)
	}

	// Создание bucket если не существует
	ctx := context.Background()
	exists, err := minioClient.BucketExists(ctx, bucketName)
	if err != nil {
		log.Fatal("Ошибка проверки bucket:", err)
	}

	if !exists {
		err = minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			log.Fatal("Ошибка создания bucket:", err)
		}
		fmt.Printf("Bucket '%s' создан\n", bucketName)
	}

	// Список изображений для загрузки
	images := []string{
		"harps.jpg",
		"james_webb.jpg", 
		"espresso.jpg",
		"spirou.jpg",
	}

	// Загрузка изображений
	for _, imageName := range images {
		// Проверяем, существует ли файл
		if _, err := os.Stat(imageName); os.IsNotExist(err) {
			fmt.Printf("Файл %s не найден, пропускаем\n", imageName)
			continue
		}

		// Загружаем файл
		_, err = minioClient.FPutObject(ctx, bucketName, imageName, imageName, minio.PutObjectOptions{})
		if err != nil {
			log.Printf("Ошибка загрузки %s: %v", imageName, err)
			continue
		}

		fmt.Printf("Изображение %s успешно загружено\n", imageName)
	}

	fmt.Println("Загрузка завершена!")
}
