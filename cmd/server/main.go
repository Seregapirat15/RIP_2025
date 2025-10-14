package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"lab2/internal/database"
	"lab2/internal/handlers"
	"lab2/pkg/config"
)

func main() {
	// Получение конфигурации
	cfg := config.GetExoplanetCalculationConfig()
	
	// Отладочная информация
	fmt.Printf("Конфигурация подключения к PostgreSQL: %s:%d\n", cfg.DBHost, cfg.DBPort)
	fmt.Printf("Конфигурация подключения к MinIO: %s\n", cfg.MinIOEndpoint)

	// Инициализация подключений
	err := database.InitPostgreSQLConnection(cfg)
	if err != nil {
		log.Fatal("Ошибка подключения к PostgreSQL:", err)
	}
	defer database.ClosePostgreSQLConnection()

	err = database.InitMinIOTelescopeImagesConnection(cfg)
	if err != nil {
		log.Fatal("Ошибка подключения к MinIO:", err)
	}

	// Настройка роутера
	r := mux.NewRouter()

	// Статические файлы
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Маршруты
	r.HandleFunc("/", handlers.AstronomicalTelescopeInstrumentsHandler).Methods("GET")
	r.HandleFunc("/tool/{id}", handlers.AstronomicalTelescopeInstrumentDetailHandler).Methods("GET")
	r.HandleFunc("/calculation/{id}", handlers.ExoplanetMassCalculationHandler).Methods("GET")
	r.HandleFunc("/add-to-calculation", handlers.AddAstronomicalTelescopeInstrumentToExoplanetMassCalculationHandler).Methods("POST")
	r.HandleFunc("/delete-calculation/{id}", handlers.DeleteExoplanetMassCalculationHandler).Methods("POST")

	fmt.Println("Сервер системы расчета массы экзопланет запущен на http://localhost:8081")
	fmt.Println("PostgreSQL подключен к:", cfg.DBHost+":"+fmt.Sprintf("%d", cfg.DBPort))
	fmt.Println("MinIO подключен к:", cfg.MinIOEndpoint)
	fmt.Println("Adminer доступен на: http://localhost:8080")
	fmt.Println("Доступные страницы:")
	fmt.Println("  GET / - Список астрономических инструментов")
	fmt.Println("  GET /tool/{id} - Детали инструмента")
	fmt.Println("  GET /calculation/{id} - Просмотр заявки")
	fmt.Println("  POST /add-to-calculation - Добавить инструмент в заявку")
	fmt.Println("  POST /delete-calculation/{id} - Удалить заявку")

	err = http.ListenAndServe(":"+cfg.ServerPort, r)
	if err != nil {
		log.Fatal("Ошибка запуска сервера:", err)
	}
}
