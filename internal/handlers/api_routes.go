package handlers

import (
	"net/http"

	"lab4/internal/middleware"
	"github.com/gorilla/mux"
)

// SetupAPIRoutes настраивает маршруты API
func SetupAPIRoutes(r *mux.Router) {
	// Создаем подроутер для API
	api := r.PathPrefix("/api").Subrouter()
	
	// === ПУБЛИЧНЫЕ УСЛУГИ (доступно всем) ===
	api.HandleFunc("/services", GetServicesHandler).Methods("GET")                    // 1. GET список услуг с фильтрацией
	api.HandleFunc("/services/{id}", GetServiceHandler).Methods("GET")               // 2. GET одна услуга
	
	// === АВТОРИЗАЦИЯ (публичные) ===
	api.HandleFunc("/auth/register", RegisterUserHandler).Methods("POST")             // Регистрация
	api.HandleFunc("/auth/login", LoginUserHandler).Methods("POST")                   // Вход
	api.HandleFunc("/auth/logout", LogoutUserHandler).Methods("POST")                 // Выход
	
	// === ВРЕМЕННО ПУБЛИЧНЫЕ (для демонстрации) ===
	api.HandleFunc("/orders/cart", GetCartIconHandler).Methods("GET")                 // Иконка корзины
	api.HandleFunc("/orders/services", AddServiceToOrderHandler).Methods("POST")      // Добавление услуги
	
	// === ЗАЩИЩЕННЫЕ МАРШРУТЫ (требуют авторизации) ===
	protected := api.PathPrefix("").Subrouter()
	protected.Use(middleware.AuthMiddleware)
	
	// Пользовательские данные
	protected.HandleFunc("/auth/me", GetUserHandler).Methods("GET")                   // Данные пользователя
	protected.HandleFunc("/auth/me", UpdateUserHandler).Methods("PUT")               // Обновление пользователя
	
	// Корзина и заявки пользователя
	protected.HandleFunc("/orders/cart", GetCartIconHandler).Methods("GET")          // Иконка корзины
	protected.HandleFunc("/orders", GetOrdersHandler).Methods("GET")                // Список заявок пользователя
	protected.HandleFunc("/orders/{id}", GetOrderHandler).Methods("GET")             // Одна заявка
	protected.HandleFunc("/orders/{id}", UpdateOrderHandler).Methods("PUT")         // Изменение заявки
	protected.HandleFunc("/orders/{id}/form", FormOrderHandler).Methods("PUT")       // Формирование заявки
	protected.HandleFunc("/orders/{id}", DeleteOrderHandler).Methods("DELETE")      // Удаление заявки
	
	// Добавление услуг в заявку
	protected.HandleFunc("/orders/services", AddServiceToOrderHandler).Methods("POST") // Добавление услуги
	protected.HandleFunc("/orders/{order_id}/services/{service_id}", DeleteOrderServiceHandler).Methods("DELETE") // Удаление услуги
	protected.HandleFunc("/orders/{order_id}/services/{service_id}", UpdateOrderServiceHandler).Methods("PUT") // Изменение связи
	
	// === АДМИНИСТРАТОРСКИЕ МАРШРУТЫ (только для модераторов) ===
	admin := api.PathPrefix("").Subrouter()
	admin.Use(middleware.AuthMiddleware)
	admin.Use(middleware.RequireModerator)
	
	// Управление услугами
	admin.HandleFunc("/services", CreateServiceHandler).Methods("POST")               // Создание услуги
	admin.HandleFunc("/services/{id}", UpdateServiceHandler).Methods("PUT")         // Изменение услуги
	admin.HandleFunc("/services/{id}", DeleteServiceHandler).Methods("DELETE")      // Удаление услуги
	admin.HandleFunc("/services/{id}/image", UploadServiceImageHandler).Methods("POST") // Загрузка изображения
	
	// Завершение заявок (только модераторы)
	admin.HandleFunc("/orders/{id}/complete", CompleteOrderHandler).Methods("PUT")    // Завершение заявки
	
	// Все заявки для модераторов
	admin.HandleFunc("/admin/orders", GetAllOrdersHandler).Methods("GET")             // Все заявки для модератора
	
	// CORS middleware
	api.Use(corsMiddleware)
}

// corsMiddleware добавляет CORS заголовки
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		
		next.ServeHTTP(w, r)
	})
}

