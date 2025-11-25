package session

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	RedisClient *redis.Client
	ctx         = context.Background()
)

// SessionData представляет данные сессии
type SessionData struct {
	UserID int    `json:"user_id"`
	Role   string `json:"role"`
	Login  string `json:"login"`
}

// InitRedis инициализирует подключение к Redis
func InitRedis() error {
	redisHost := os.Getenv("REDIS_HOST")
	if redisHost == "" {
		redisHost = "localhost"
	}

	redisPort := os.Getenv("REDIS_PORT")
	if redisPort == "" {
		redisPort = "6379"
	}

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", redisHost, redisPort),
		Password: "", // нет пароля
		DB:       0,  // используем базу данных по умолчанию
	})

	// Проверяем подключение
	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		return fmt.Errorf("ошибка подключения к Redis: %v", err)
	}

	return nil
}

// CreateSession создает новую сессию в Redis
func CreateSession(sessionID string, data SessionData, expiration time.Duration) error {
	key := fmt.Sprintf("session:%s", sessionID)
	
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("ошибка сериализации данных сессии: %v", err)
	}

	err = RedisClient.Set(ctx, key, jsonData, expiration).Err()
	if err != nil {
		return fmt.Errorf("ошибка сохранения сессии в Redis: %v", err)
	}

	// Также сохраняем обратную связь: user_id -> session_id
	userKey := fmt.Sprintf("user:%d:session", data.UserID)
	err = RedisClient.Set(ctx, userKey, sessionID, expiration).Err()
	if err != nil {
		return fmt.Errorf("ошибка сохранения связи пользователя с сессией: %v", err)
	}

	return nil
}

// GetSession получает данные сессии из Redis
func GetSession(sessionID string) (*SessionData, error) {
	key := fmt.Sprintf("session:%s", sessionID)
	
	val, err := RedisClient.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, fmt.Errorf("сессия не найдена")
	} else if err != nil {
		return nil, fmt.Errorf("ошибка получения сессии из Redis: %v", err)
	}

	var data SessionData
	err = json.Unmarshal([]byte(val), &data)
	if err != nil {
		return nil, fmt.Errorf("ошибка десериализации данных сессии: %v", err)
	}

	return &data, nil
}

// DeleteSession удаляет сессию из Redis
func DeleteSession(sessionID string) error {
	key := fmt.Sprintf("session:%s", sessionID)
	
	// Получаем данные сессии для удаления связи пользователя
	data, err := GetSession(sessionID)
	if err == nil && data != nil {
		userKey := fmt.Sprintf("user:%d:session", data.UserID)
		RedisClient.Del(ctx, userKey)
	}

	err = RedisClient.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("ошибка удаления сессии из Redis: %v", err)
	}

	return nil
}

// GetAllSessions возвращает все сессии (для отладки)
func GetAllSessions() (map[string]*SessionData, error) {
	keys, err := RedisClient.Keys(ctx, "session:*").Result()
	if err != nil {
		return nil, fmt.Errorf("ошибка получения ключей сессий: %v", err)
	}

	sessions := make(map[string]*SessionData)
	for _, key := range keys {
		val, err := RedisClient.Get(ctx, key).Result()
		if err != nil {
			continue
		}

		var data SessionData
		if err := json.Unmarshal([]byte(val), &data); err != nil {
			continue
		}

		// Извлекаем session_id из ключа (session:xxx -> xxx)
		sessionID := key[8:] // убираем "session:"
		sessions[sessionID] = &data
	}

	return sessions, nil
}

// GetUserSession получает сессию пользователя по user_id
func GetUserSession(userID int) (string, error) {
	userKey := fmt.Sprintf("user:%d:session", userID)
	sessionID, err := RedisClient.Get(ctx, userKey).Result()
	if err == redis.Nil {
		return "", fmt.Errorf("сессия пользователя не найдена")
	} else if err != nil {
		return "", fmt.Errorf("ошибка получения сессии пользователя: %v", err)
	}

	return sessionID, nil
}



