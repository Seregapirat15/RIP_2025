package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"lab4/internal/auth"
	"lab4/internal/database"
	"lab4/internal/middleware"
	"lab4/internal/models"
	"lab4/internal/session"
)

// RegisterUserRequest представляет запрос на регистрацию
type RegisterUserRequest struct {
	Login     string `json:"login"`
	Password  string `json:"password"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Role      string `json:"role"`
}

// RegisterUserHandler - POST /api/users/register - регистрация пользователя
func RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	var req RegisterUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Неверный JSON", http.StatusBadRequest)
		return
	}
	
	// Валидация обязательных полей
	if req.Login == "" || req.Password == "" || req.Email == "" || req.FirstName == "" || req.LastName == "" {
		http.Error(w, "Заполните все обязательные поля", http.StatusBadRequest)
		return
	}
	
	// Хешируем пароль
	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		log.Printf("Ошибка хеширования пароля: %v", err)
		http.Error(w, "Ошибка обработки пароля", http.StatusInternalServerError)
		return
	}
	
	// Устанавливаем системные поля
	user := models.User{
		Login:        req.Login,
		Email:        req.Email,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		PasswordHash: hashedPassword,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		IsActive:     true,
		Role:         "user", // По умолчанию пользователь
	}
	
	// Если указана роль, используем её
	if req.Role != "" {
		user.Role = req.Role
	}
	
	userID, err := database.CreateUser(user)
	if err != nil {
		log.Printf("Ошибка создания пользователя: %v", err)
		http.Error(w, "Ошибка создания пользователя", http.StatusInternalServerError)
		return
	}
	
	user.ID = userID
	// Не возвращаем хеш пароля
	user.PasswordHash = ""
	
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

// GetUserHandler - GET /api/users/me - получение данных пользователя
func GetUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	// Получаем ID пользователя из контекста (установлен middleware)
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "Пользователь не авторизован", http.StatusUnauthorized)
		return
	}
	
	user, err := database.GetUserByID(userID)
	if err != nil {
		log.Printf("Ошибка получения пользователя: %v", err)
		http.Error(w, "Пользователь не найден", http.StatusNotFound)
		return
	}
	
	// Не возвращаем хеш пароля
	user.PasswordHash = ""
	
	json.NewEncoder(w).Encode(user)
}

// UpdateUserHandler - PUT /api/users/me - обновление пользователя
func UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	// Получаем ID пользователя из контекста
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "Пользователь не авторизован", http.StatusUnauthorized)
		return
	}
	
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Неверный JSON", http.StatusBadRequest)
		return
	}
	
	// Получаем текущего пользователя для сохранения системных полей
	currentUser, err := database.GetUserByID(userID)
	if err != nil {
		http.Error(w, "Пользователь не найден", http.StatusNotFound)
		return
	}
	
	// Обновляем только разрешенные поля
	user.ID = userID
	user.CreatedAt = currentUser.CreatedAt // Сохраняем оригинальную дату
	user.UpdatedAt = time.Now()
	user.Role = currentUser.Role // Сохраняем оригинальную роль
	user.PasswordHash = currentUser.PasswordHash // Сохраняем хеш пароля
	
	err = database.UpdateUser(user)
	if err != nil {
		log.Printf("Ошибка обновления пользователя: %v", err)
		http.Error(w, "Ошибка обновления пользователя", http.StatusInternalServerError)
		return
	}
	
	// Не возвращаем хеш пароля
	user.PasswordHash = ""
	
	json.NewEncoder(w).Encode(user)
}

// LoginUserHandler - POST /api/auth/login - аутентификация
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
	
	// Получаем пользователя по логину
	user, err := database.GetUserByLogin(credentials.Login)
	if err != nil {
		http.Error(w, "Неверные учетные данные", http.StatusUnauthorized)
		return
	}
	
	// Проверяем пароль
	if !auth.CheckPasswordHash(credentials.Password, user.PasswordHash) {
		http.Error(w, "Неверные учетные данные", http.StatusUnauthorized)
		return
	}
	
	// Проверяем, что пользователь активен
	if !user.IsActive {
		http.Error(w, "Аккаунт заблокирован", http.StatusUnauthorized)
		return
	}
	
	// Генерируем JWT токен
	token, err := auth.GenerateToken(user.ID, user.Role, user.Login)
	if err != nil {
		log.Printf("Ошибка генерации токена: %v", err)
		http.Error(w, "Ошибка авторизации", http.StatusInternalServerError)
		return
	}
	
	// Создаем сессию в Redis
	sessionID := uuid.New().String()
	sessionData := session.SessionData{
		UserID: user.ID,
		Role:   user.Role,
		Login:  user.Login,
	}
	
	// Сохраняем сессию на 24 часа
	err = session.CreateSession(sessionID, sessionData, 24*time.Hour)
	if err != nil {
		log.Printf("Ошибка создания сессии в Redis: %v", err)
		// Не прерываем авторизацию, если Redis недоступен
	} else {
		// Устанавливаем куки с session_id
		cookie := &http.Cookie{
			Name:     "session_id",
			Value:    sessionID,
			Path:     "/",
			MaxAge:   86400, // 24 часа в секундах
			HttpOnly: true,  // Защита от XSS
			SameSite: http.SameSiteLaxMode,
		}
		http.SetCookie(w, cookie)
		log.Printf("Создана сессия в Redis: session_id=%s, user_id=%d", sessionID, user.ID)
	}
	
	// Не возвращаем хеш пароля
	user.PasswordHash = ""
	
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user":    user,
		"token":   token,
		"message": "Успешная авторизация",
	})
}

// LogoutUserHandler - POST /api/auth/logout - деавторизация
func LogoutUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	// Получаем session_id из куки
	cookie, err := r.Cookie("session_id")
	if err == nil && cookie != nil {
		// Удаляем сессию из Redis
		err = session.DeleteSession(cookie.Value)
		if err != nil {
			log.Printf("Ошибка удаления сессии из Redis: %v", err)
		} else {
			log.Printf("Удалена сессия из Redis: session_id=%s", cookie.Value)
		}
		
		// Удаляем куки
		cookie.MaxAge = -1
		cookie.Path = "/"
		http.SetCookie(w, cookie)
	}
	
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Успешный выход из системы",
	})
}
