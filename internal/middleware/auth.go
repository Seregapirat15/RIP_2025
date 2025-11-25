package middleware

import (
	"context"
	"log"
	"net/http"

	"lab4/internal/auth"
	"lab4/internal/session"
)

// AuthMiddleware проверяет JWT токен в заголовке Authorization или сессию в куки
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var userID int
		var role string
		var login string
		var authenticated bool

		// Приоритет 1: Проверяем JWT токен в заголовке Authorization
		authHeader := r.Header.Get("Authorization")
		if authHeader != "" {
			tokenString, err := auth.ExtractTokenFromHeader(authHeader)
			if err == nil {
				claims, err := auth.ValidateToken(tokenString)
				if err == nil {
					userID = claims.UserID
					role = claims.Role
					login = claims.Login
					authenticated = true
					log.Printf("AuthMiddleware: Авторизация через JWT, user_id=%d", userID)
				}
			}
		}

		// Приоритет 2: Если JWT не найден, проверяем сессию в куки
		if !authenticated {
			cookie, err := r.Cookie("session_id")
			if err == nil && cookie != nil {
				sessionData, err := session.GetSession(cookie.Value)
				if err == nil {
					userID = sessionData.UserID
					role = sessionData.Role
					login = sessionData.Login
					authenticated = true
					log.Printf("AuthMiddleware: Авторизация через куки, session_id=%s, user_id=%d", cookie.Value, userID)
				}
			}
		}

		// Если не авторизован ни одним способом
		if !authenticated {
			http.Error(w, "Требуется авторизация", http.StatusUnauthorized)
			return
		}

		// Добавляем информацию о пользователе в контекст
		ctx := context.WithValue(r.Context(), "user_id", userID)
		ctx = context.WithValue(ctx, "user_role", role)
		ctx = context.WithValue(ctx, "user_login", login)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequireRole проверяет, что пользователь имеет определенную роль
func RequireRole(requiredRole string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userRole := r.Context().Value("user_role")
			if userRole == nil {
				http.Error(w, "Требуется авторизация", http.StatusUnauthorized)
				return
			}

			role := userRole.(string)
			if role != requiredRole && role != "admin" {
				http.Error(w, "Недостаточно прав", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// RequireModerator проверяет права модератора
func RequireModerator(next http.Handler) http.Handler {
	return RequireRole("moderator")(next)
}

// RequireUser проверяет, что пользователь авторизован (любая роль)
func RequireUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value("user_id")
		if userID == nil {
			http.Error(w, "Требуется авторизация", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// GetUserID извлекает ID пользователя из контекста
func GetUserID(ctx context.Context) (int, bool) {
	userID, ok := ctx.Value("user_id").(int)
	return userID, ok
}

// GetUserRole извлекает роль пользователя из контекста
func GetUserRole(ctx context.Context) (string, bool) {
	role, ok := ctx.Value("user_role").(string)
	return role, ok
}
