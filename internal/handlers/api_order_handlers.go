package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"lab2/internal/database"
	"lab2/internal/models"
)

// GetCartIconHandler - GET /api/orders/cart - иконка корзины
func GetCartIconHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	// Получаем заявку-черновик пользователя
	order, err := database.GetDraftOrder(FIXED_CREATOR_ID)
	if err != nil {
		// Если заявки нет, возвращаем пустую корзину
		json.NewEncoder(w).Encode(models.CartIcon{
			OrderID:       0,
			CalculationID: 0,
			ServicesCount: 0,
		})
		return
	}
	
	// Получаем количество услуг в заявке
	serviceCount, err := database.GetOrderServiceCount(order.ID)
	if err != nil {
		log.Printf("Ошибка получения количества услуг: %v", err)
		serviceCount = 0
	}
	
	json.NewEncoder(w).Encode(models.CartIcon{
		OrderID:       order.ID,
		CalculationID: order.ID, // В данном контексте calculation_id = order_id
		ServicesCount: serviceCount,
	})
}

// GetOrdersHandler - GET /api/orders - список заявок с фильтрацией
func GetOrdersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	// Парсинг параметров фильтрации
	filter := models.OrderFilter{
		Status: r.URL.Query().Get("status"),
	}
	
	// Парсинг дат
	if formationFrom := r.URL.Query().Get("formation_from"); formationFrom != "" {
		if t, err := time.Parse("2006-01-02", formationFrom); err == nil {
			filter.FormationFrom = &t
		}
	}
	
	if formationTo := r.URL.Query().Get("formation_to"); formationTo != "" {
		if t, err := time.Parse("2006-01-02", formationTo); err == nil {
			filter.FormationTo = &t
		}
	}
	
	orders, err := database.GetOrdersWithFilter(filter)
	if err != nil {
		log.Printf("Ошибка получения заявок: %v", err)
		http.Error(w, "Ошибка получения заявок", http.StatusInternalServerError)
		return
	}
	
	json.NewEncoder(w).Encode(orders)
}

// GetOrderHandler - GET /api/orders/{id} - одна заявка с услугами
func GetOrderHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	// Извлечение ID из URL
	pathParts := strings.Split(r.URL.Path, "/")
	orderID, err := strconv.Atoi(pathParts[len(pathParts)-1])
	if err != nil {
		http.Error(w, "Неверный ID заявки", http.StatusBadRequest)
		return
	}
	
	order, err := database.GetOrderWithServices(orderID)
	if err != nil {
		log.Printf("Ошибка получения заявки: %v", err)
		http.Error(w, "Заявка не найдена", http.StatusNotFound)
		return
	}
	
	json.NewEncoder(w).Encode(order)
}

// UpdateOrderHandler - PUT /api/orders/{id} - изменение заявки
func UpdateOrderHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	// Извлечение ID из URL
	pathParts := strings.Split(r.URL.Path, "/")
	orderID, err := strconv.Atoi(pathParts[len(pathParts)-1])
	if err != nil {
		http.Error(w, "Неверный ID заявки", http.StatusBadRequest)
		return
	}
	
	var order models.Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		http.Error(w, "Неверный JSON", http.StatusBadRequest)
		return
	}
	
	// Проверяем права доступа
	if order.CreatorID != FIXED_CREATOR_ID {
		http.Error(w, "Нет прав на изменение заявки", http.StatusForbidden)
		return
	}
	
	order.ID = orderID
	err = database.UpdateOrder(order)
	if err != nil {
		log.Printf("Ошибка обновления заявки: %v", err)
		http.Error(w, "Ошибка обновления заявки", http.StatusInternalServerError)
		return
	}
	
	json.NewEncoder(w).Encode(order)
}

// FormOrderHandler - PUT /api/orders/{id}/form - сформировать заявку
func FormOrderHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	// Извлечение ID из URL
	pathParts := strings.Split(r.URL.Path, "/")
	orderID, err := strconv.Atoi(pathParts[len(pathParts)-2])
	if err != nil {
		http.Error(w, "Неверный ID заявки", http.StatusBadRequest)
		return
	}
	
	// Проверяем, что заявка принадлежит пользователю
	order, err := database.GetOrderByID(orderID)
	if err != nil {
		http.Error(w, "Заявка не найдена", http.StatusNotFound)
		return
	}
	
	if order.CreatorID != FIXED_CREATOR_ID {
		http.Error(w, "Нет прав на формирование заявки", http.StatusForbidden)
		return
	}
	
	// Проверяем статус
	if order.Status != "черновик" {
		http.Error(w, "Можно формировать только черновики", http.StatusBadRequest)
		return
	}
	
	// Проверяем обязательные поля
	if err := database.ValidateOrderForForming(orderID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	// Формируем заявку
	err = database.FormOrder(orderID)
	if err != nil {
		log.Printf("Ошибка формирования заявки: %v", err)
		http.Error(w, "Ошибка формирования заявки", http.StatusInternalServerError)
		return
	}
	
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Заявка сформирована",
		"status":  "сформирован",
	})
}

// CompleteOrderHandler - PUT /api/orders/{id}/complete - завершить/отклонить заявку
func CompleteOrderHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	// Извлечение ID из URL
	pathParts := strings.Split(r.URL.Path, "/")
	orderID, err := strconv.Atoi(pathParts[len(pathParts)-2])
	if err != nil {
		http.Error(w, "Неверный ID заявки", http.StatusBadRequest)
		return
	}
	
	var request struct {
		Action string `json:"action"` // complete, reject
		Result string `json:"result,omitempty"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Неверный JSON", http.StatusBadRequest)
		return
	}
	
	// Проверяем статус заявки
	order, err := database.GetOrderByID(orderID)
	if err != nil {
		http.Error(w, "Заявка не найдена", http.StatusNotFound)
		return
	}
	
	if order.Status != "сформирован" {
		http.Error(w, "Можно завершать только сформированные заявки", http.StatusBadRequest)
		return
	}
	
	// Завершаем заявку
	err = database.CompleteOrder(orderID, request.Action, request.Result)
	if err != nil {
		log.Printf("Ошибка завершения заявки: %v", err)
		http.Error(w, "Ошибка завершения заявки", http.StatusInternalServerError)
		return
	}
	
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Заявка " + request.Action + "d",
		"status":  request.Action + "d",
	})
}

// DeleteOrderHandler - DELETE /api/orders/{id} - удаление заявки
func DeleteOrderHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	// Извлечение ID из URL
	pathParts := strings.Split(r.URL.Path, "/")
	orderID, err := strconv.Atoi(pathParts[len(pathParts)-1])
	if err != nil {
		http.Error(w, "Неверный ID заявки", http.StatusBadRequest)
		return
	}
	
	// Проверяем права доступа
	order, err := database.GetOrderByID(orderID)
	if err != nil {
		http.Error(w, "Заявка не найдена", http.StatusNotFound)
		return
	}
	
	if order.CreatorID != FIXED_CREATOR_ID {
		http.Error(w, "Нет прав на удаление заявки", http.StatusForbidden)
		return
	}
	
	err = database.DeleteOrder(orderID)
	if err != nil {
		log.Printf("Ошибка удаления заявки: %v", err)
		http.Error(w, "Ошибка удаления заявки", http.StatusInternalServerError)
		return
	}
	
	w.WriteHeader(http.StatusNoContent)
}

