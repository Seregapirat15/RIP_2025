package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"lab2/internal/database"
	"lab2/internal/models"
)

// DeleteOrderServiceHandler - DELETE /api/orders/{order_id}/services/{service_id} - удаление услуги из заявки
func DeleteOrderServiceHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	// Извлечение ID из URL
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 6 {
		http.Error(w, "Неверный URL", http.StatusBadRequest)
		return
	}
	
	orderID, err := strconv.Atoi(pathParts[len(pathParts)-3])
	if err != nil {
		http.Error(w, "Неверный ID заявки", http.StatusBadRequest)
		return
	}
	
	serviceID, err := strconv.Atoi(pathParts[len(pathParts)-1])
	if err != nil {
		http.Error(w, "Неверный ID услуги", http.StatusBadRequest)
		return
	}
	
	// Проверяем права доступа
	order, err := database.GetOrderByID(orderID)
	if err != nil {
		http.Error(w, "Заявка не найдена", http.StatusNotFound)
		return
	}
	
	if order.CreatorID != FIXED_CREATOR_ID {
		http.Error(w, "Нет прав на изменение заявки", http.StatusForbidden)
		return
	}
	
	// Удаляем услугу из заявки
	err = database.DeleteOrderService(orderID, serviceID)
	if err != nil {
		log.Printf("Ошибка удаления услуги из заявки: %v", err)
		http.Error(w, "Ошибка удаления услуги из заявки", http.StatusInternalServerError)
		return
	}
	
	w.WriteHeader(http.StatusNoContent)
}

// UpdateOrderServiceHandler - PUT /api/orders/{order_id}/services/{service_id} - изменение связи м-м
func UpdateOrderServiceHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	// Извлечение ID из URL
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 6 {
		http.Error(w, "Неверный URL", http.StatusBadRequest)
		return
	}
	
	orderID, err := strconv.Atoi(pathParts[len(pathParts)-3])
	if err != nil {
		http.Error(w, "Неверный ID заявки", http.StatusBadRequest)
		return
	}
	
	serviceID, err := strconv.Atoi(pathParts[len(pathParts)-1])
	if err != nil {
		http.Error(w, "Неверный ID услуги", http.StatusBadRequest)
		return
	}
	
	var orderService models.OrderService
	if err := json.NewDecoder(r.Body).Decode(&orderService); err != nil {
		http.Error(w, "Неверный JSON", http.StatusBadRequest)
		return
	}
	
	// Проверяем права доступа
	order, err := database.GetOrderByID(orderID)
	if err != nil {
		http.Error(w, "Заявка не найдена", http.StatusNotFound)
		return
	}
	
	if order.CreatorID != FIXED_CREATOR_ID {
		http.Error(w, "Нет прав на изменение заявки", http.StatusForbidden)
		return
	}
	
	// Устанавливаем ID
	orderService.OrderID = orderID
	orderService.ServiceID = serviceID
	
	// Обновляем связь
	err = database.UpdateOrderService(orderService)
	if err != nil {
		log.Printf("Ошибка обновления связи заявка-услуга: %v", err)
		http.Error(w, "Ошибка обновления связи", http.StatusInternalServerError)
		return
	}
	
	json.NewEncoder(w).Encode(orderService)
}

