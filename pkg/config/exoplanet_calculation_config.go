package config

import "os"

// ExoplanetCalculationConfig содержит конфигурацию системы расчета экзопланет
type ExoplanetCalculationConfig struct {
	// PostgreSQL
	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string
	
	// MinIO для хранения изображений телескопов
	MinIOEndpoint        string
	MinIOAccessKeyID     string
	MinIOSecretAccessKey string
	MinIOUseSSL          bool
	MinIOBucketName      string
	
	// Server
	ServerPort string
}

// GetExoplanetCalculationConfig возвращает конфигурацию по умолчанию
func GetExoplanetCalculationConfig() *ExoplanetCalculationConfig {
	// Проверяем, запущено ли приложение в Docker
	isDocker := os.Getenv("DOCKER_ENV") == "true"
	
	config := &ExoplanetCalculationConfig{
		// PostgreSQL
		DBUser:     "postgres",
		DBPassword: "postgres123",
		DBName:     "exoplanet_calculations",
		DBPort:     5432,
		
		// MinIO для изображений телескопов
		MinIOAccessKeyID:     "minioadmin123",
		MinIOSecretAccessKey: "minioadmin123456",
		MinIOUseSSL:          false,
		MinIOBucketName:      "telescope-images",
		
		// Server
		ServerPort: "8081",
	}
	
	if isDocker {
		// Конфигурация для Docker
		config.DBHost = "postgres"
		config.MinIOEndpoint = "minio:9000"
	} else {
		// Конфигурация для локального запуска
		config.DBHost = "localhost"
		config.MinIOEndpoint = "localhost:9000"
	}
	
	return config
}
