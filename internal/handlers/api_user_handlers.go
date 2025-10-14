package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"lab2/internal/database"
	"lab2/internal/models"
)

// RegisterUserHandler - POST /api/users/register - регистрация пользователя
func RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Неверный JSON", http.StatusBadRequest)
		return
	}
	
	// Валидация обязательных полей
	if user.Login == "" || user.Email == "" || user.FirstName == "" || user.LastName == "" {
		http.Error(w, "Заполните все обязательные поля", http.StatusBadRequest)
		return
	}
	
	// Устанавливаем системные поля
	user.CreatedAt = time.Now()
	user.IsActive = true
	user.Role = "user" // По умолчанию пользователь
	
	userID, err := database.CreateUser(user)
	if err != nil {
		log.Printf("Ошибка создания пользователя: %v", err)
		http.Error(w, "Ошибка создания пользователя", http.StatusInternalServerError)
		return
	}
	
	user.ID = userID
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

// GetUserHandler - GET /api/users/me - получение данных пользователя
func GetUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	// В реальном приложении ID пользователя получается из токена
	userID := FIXED_CREATOR_ID
	
	user, err := database.GetUserByID(userID)
	if err != nil {
		log.Printf("Ошибка получения пользователя: %v", err)
		http.Error(w, "Пользователь не найден", http.StatusNotFound)
		return
	}
	
	json.NewEncoder(w).Encode(user)
}

// UpdateUserHandler - PUT /api/users/me - обновление пользователя
func UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	// В реальном приложении ID пользователя получается из токена
	userID := FIXED_CREATOR_ID
	
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Неверный JSON", http.StatusBadRequest)
		return
	}
	
	// Запрещаем изменение системных полей
	user.ID = userID
	user.CreatedAt = time.Now() // Сохраняем оригинальную дату
	user.Role = "creator"       // Сохраняем оригинальную роль
	
	err := database.UpdateUser(user)
	if err != nil {
		log.Printf("Ошибка обновления пользователя: %v", err)
		http.Error(w, "Ошибка обновления пользователя", http.StatusInternalServerError)
		return
	}
	
	json.NewEncoder(w).Encode(user)
}

// LoginUserHandler - POST /api/users/login - аутентификация
func LoginUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	var credentials struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		http.Error(w, "Неверный JSON", http.StatusBadRequest)
		return
	}
	
	// В реальном приложении здесь проверка пароля
	user, err := database.GetUserByLogin(credentials.Login)
	if err != nil {
		http.Error(w, "Неверные учетные данные", http.StatusUnauthorized)
		return
	}
	
	// В реальном приложении здесь генерация JWT токена
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user": user,
		"token": "fake-jwt-token-" + strconv.Itoa(user.ID),
	})
}

// LogoutUserHandler - POST /api/users/logout - деавторизация
func LogoutUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	// В реальном приложении здесь инвалидация токена
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Успешный выход из системы",
	})
}
