package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"lab2/pkg/config"
)

// PostgreSQLConnection представляет подключение к PostgreSQL
var PostgreSQLConnection *sql.DB

// InitPostgreSQLConnection инициализирует подключение к PostgreSQL
func InitPostgreSQLConnection(cfg *config.ExoplanetCalculationConfig) error {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable connect_timeout=10",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)
	
	var err error
	PostgreSQLConnection, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		return fmt.Errorf("ошибка подключения к PostgreSQL: %w", err)
	}
	
	// Проверка подключения
	err = PostgreSQLConnection.Ping()
	if err != nil {
		return fmt.Errorf("ошибка ping PostgreSQL: %w", err)
	}
	
	log.Println("Подключение к PostgreSQL установлено")
	return nil
}

// ClosePostgreSQLConnection закрывает подключение к PostgreSQL
func ClosePostgreSQLConnection() error {
	if PostgreSQLConnection != nil {
		return PostgreSQLConnection.Close()
	}
	return nil
}

