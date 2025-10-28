package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"lab4/internal/database"
	"lab4/internal/middleware"
	"lab4/internal/models"
)

// Константа для зафиксированного пользователя-создателя
const FIXED_CREATOR_ID = 1

// GetServicesHandler - GET /api/services - список услуг с фильтрацией
func GetServicesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	// Получение параметров фильтрации
	filter := models.ServiceFilter{
		Type:   r.URL.Query().Get("type"),
		Status: r.URL.Query().Get("status"),
		Search: r.URL.Query().Get("search"),
	}
	
	services, err := database.GetServicesWithFilter(filter)
	if err != nil {
		log.Printf("Ошибка получения услуг: %v", err)
		http.Error(w, "Ошибка получения услуг", http.StatusInternalServerError)
		return
	}
	
	json.NewEncoder(w).Encode(services)
}

// GetServiceHandler - GET /api/services/{id} - одна услуга
func GetServiceHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	// Извлечение ID из URL
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		http.Error(w, "Неверный URL", http.StatusBadRequest)
		return
	}
	
	serviceID, err := strconv.Atoi(pathParts[len(pathParts)-1])
	if err != nil {
		http.Error(w, "Неверный ID услуги", http.StatusBadRequest)
		return
	}
	
	service, err := database.GetServiceByID(serviceID)
	if err != nil {
		log.Printf("Ошибка получения услуги: %v", err)
		http.Error(w, "Услуга не найдена", http.StatusNotFound)
		return
	}
	
	json.NewEncoder(w).Encode(service)
}

// CreateServiceHandler - POST /api/services - добавление услуги
func CreateServiceHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	var service models.Service
	if err := json.NewDecoder(r.Body).Decode(&service); err != nil {
		http.Error(w, "Неверный JSON", http.StatusBadRequest)
		return
	}
	
	// Устанавливаем системные поля
	service.CreatedAt = time.Now()
	service.IsDeleted = false
	service.Status = "Активен"
	
	serviceID, err := database.CreateService(service)
	if err != nil {
		log.Printf("Ошибка создания услуги: %v", err)
		http.Error(w, "Ошибка создания услуги", http.StatusInternalServerError)
		return
	}
	
	service.ID = serviceID
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(service)
}

// UpdateServiceHandler - PUT /api/services/{id} - изменение услуги
func UpdateServiceHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	// Извлечение ID из URL
	pathParts := strings.Split(r.URL.Path, "/")
	serviceID, err := strconv.Atoi(pathParts[len(pathParts)-1])
	if err != nil {
		http.Error(w, "Неверный ID услуги", http.StatusBadRequest)
		return
	}
	
	var service models.Service
	if err := json.NewDecoder(r.Body).Decode(&service); err != nil {
		http.Error(w, "Неверный JSON", http.StatusBadRequest)
		return
	}
	
	service.ID = serviceID
	err = database.UpdateService(service)
	if err != nil {
		log.Printf("Ошибка обновления услуги: %v", err)
		http.Error(w, "Ошибка обновления услуги", http.StatusInternalServerError)
		return
	}
	
	json.NewEncoder(w).Encode(service)
}

// DeleteServiceHandler - DELETE /api/services/{id} - удаление услуги
func DeleteServiceHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	// Извлечение ID из URL
	pathParts := strings.Split(r.URL.Path, "/")
	serviceID, err := strconv.Atoi(pathParts[len(pathParts)-1])
	if err != nil {
		http.Error(w, "Неверный ID услуги", http.StatusBadRequest)
		return
	}
	
	err = database.DeleteService(serviceID)
	if err != nil {
		log.Printf("Ошибка удаления услуги: %v", err)
		http.Error(w, "Ошибка удаления услуги", http.StatusInternalServerError)
		return
	}
	
	w.WriteHeader(http.StatusNoContent)
}

// AddServiceToOrderHandler - POST /api/orders/services - добавление услуги в заявку
func AddServiceToOrderHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	// Получаем ID пользователя из контекста
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "Пользователь не авторизован", http.StatusUnauthorized)
		return
	}
	
	var orderService models.OrderService
	if err := json.NewDecoder(r.Body).Decode(&orderService); err != nil {
		http.Error(w, "Неверный JSON", http.StatusBadRequest)
		return
	}
	
	// Получаем или создаем заявку-черновик
	order, err := database.GetOrCreateDraftOrder(userID)
	if err != nil {
		log.Printf("Ошибка получения/создания заявки: %v", err)
		http.Error(w, "Ошибка работы с заявкой", http.StatusInternalServerError)
		return
	}
	
	orderService.OrderID = order.ID
	err = database.AddServiceToOrder(orderService)
	if err != nil {
		log.Printf("Ошибка добавления услуги в заявку: %v", err)
		
		// Проверяем, является ли ошибка дубликатом
		if err.Error() == "услуга уже добавлена в заявку" {
			http.Error(w, "Эта услуга уже добавлена в заявку", http.StatusConflict)
			return
		}
		
		http.Error(w, "Ошибка добавления услуги в заявку", http.StatusInternalServerError)
		return
	}
	
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"order_id": order.ID,
		"service_id": orderService.ServiceID,
		"message": "Услуга успешно добавлена в заявку",
	})
}

// UploadServiceImageHandler - POST /api/services/{id}/image - добавление изображения
func UploadServiceImageHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	// Извлечение ID из URL
	pathParts := strings.Split(r.URL.Path, "/")
	serviceID, err := strconv.Atoi(pathParts[len(pathParts)-2])
	if err != nil {
		http.Error(w, "Неверный ID услуги", http.StatusBadRequest)
		return
	}
	
	// Парсинг multipart/form-data
	err = r.ParseMultipartForm(10 << 20) // 10MB
	if err != nil {
		http.Error(w, "Ошибка парсинга формы", http.StatusBadRequest)
		return
	}
	
	file, _, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Файл не найден", http.StatusBadRequest)
		return
	}
	defer file.Close()
	
	// Генерация имени файла на латинице
	fileName := fmt.Sprintf("service_%d_%d.jpg", serviceID, time.Now().Unix())
	
	// Загрузка в MinIO (упрощенная версия)
	imageURL := fmt.Sprintf("http://localhost:9000/service-images/%s", fileName)
	
	// Обновляем URL в базе
	err = database.UpdateServiceImageURL(serviceID, imageURL)
	if err != nil {
		log.Printf("Ошибка обновления URL изображения: %v", err)
		http.Error(w, "Ошибка обновления URL изображения", http.StatusInternalServerError)
		return
	}
	
	json.NewEncoder(w).Encode(map[string]string{
		"image_url": imageURL,
	})
}

