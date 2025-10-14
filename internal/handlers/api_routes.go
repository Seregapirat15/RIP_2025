package handlers

import (
	"net/http"

	"github.com/gorilla/mux"
)

// SetupAPIRoutes настраивает маршруты API
func SetupAPIRoutes(r *mux.Router) {
	// Создаем подроутер для API
	api := r.PathPrefix("/api").Subrouter()
	
	// === УСЛУГИ ===
	api.HandleFunc("/services", GetServicesHandler).Methods("GET")                    // 1. GET список услуг с фильтрацией
	api.HandleFunc("/services/{id}", GetServiceHandler).Methods("GET")               // 2. GET одна услуга
	api.HandleFunc("/services", CreateServiceHandler).Methods("POST")               // 3. POST добавление услуги
	api.HandleFunc("/services/{id}", UpdateServiceHandler).Methods("PUT")          // 4. PUT изменение услуги
	api.HandleFunc("/services/{id}", DeleteServiceHandler).Methods("DELETE")        // 5. DELETE удаление услуги
	api.HandleFunc("/services/{id}/image", UploadServiceImageHandler).Methods("POST") // 6. POST добавление изображения
	
	// === ЗАЯВКИ ===
	api.HandleFunc("/orders/cart", GetCartIconHandler).Methods("GET")              // 7. GET иконка корзины
	api.HandleFunc("/orders", GetOrdersHandler).Methods("GET")                      // 8. GET список заявок с фильтрацией
	api.HandleFunc("/orders/{id}", GetOrderHandler).Methods("GET")                  // 9. GET одна заявка с услугами
	api.HandleFunc("/orders/{id}", UpdateOrderHandler).Methods("PUT")                // 10. PUT изменение заявки
	api.HandleFunc("/orders/{id}/form", FormOrderHandler).Methods("PUT")            // 11. PUT сформировать заявку
	api.HandleFunc("/orders/{id}/complete", CompleteOrderHandler).Methods("PUT")    // 12. PUT завершить/отклонить заявку
	api.HandleFunc("/orders/{id}", DeleteOrderHandler).Methods("DELETE")            // 13. DELETE удаление заявки
	
	// === СВЯЗИ ЗАЯВКА-УСЛУГА ===
	api.HandleFunc("/orders/services", AddServiceToOrderHandler).Methods("POST")    // 14. POST добавление услуги в заявку
	api.HandleFunc("/orders/{order_id}/services/{service_id}", DeleteOrderServiceHandler).Methods("DELETE") // 15. DELETE удаление услуги из заявки
	api.HandleFunc("/orders/{order_id}/services/{service_id}", UpdateOrderServiceHandler).Methods("PUT") // 16. PUT изменение связи м-м
	
	// === ПОЛЬЗОВАТЕЛИ ===
	api.HandleFunc("/users/register", RegisterUserHandler).Methods("POST")          // 17. POST регистрация пользователя
	api.HandleFunc("/users/me", GetUserHandler).Methods("GET")                      // 18. GET данные пользователя
	api.HandleFunc("/users/me", UpdateUserHandler).Methods("PUT")                  // 19. PUT обновление пользователя
	api.HandleFunc("/users/login", LoginUserHandler).Methods("POST")               // 20. POST аутентификация
	api.HandleFunc("/users/logout", LogoutUserHandler).Methods("POST")             // 21. POST деавторизация
	
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

