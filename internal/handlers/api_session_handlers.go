package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"lab4/internal/session"
)

// GetSessionsHandler - GET /api/admin/sessions - получение всех сессий (для отладки)
func GetSessionsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	sessions, err := session.GetAllSessions()
	if err != nil {
		log.Printf("Ошибка получения сессий: %v", err)
		http.Error(w, "Ошибка получения сессий", http.StatusInternalServerError)
		return
	}

	// Форматируем ответ для удобного просмотра
	result := make(map[string]interface{})
	result["total"] = len(sessions)
	result["sessions"] = sessions

	json.NewEncoder(w).Encode(result)
}



