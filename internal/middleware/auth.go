package middleware

import (
	"context"
	"net/http"

	"lab4/internal/auth"
)

// AuthMiddleware проверяет JWT токен в заголовке Authorization
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		
		// Логируем для отладки
		// log.Printf("AuthMiddleware: Path=%s, Method=%s, AuthHeader=%s", r.URL.Path, r.Method, authHeader)
		
		tokenString, err := auth.ExtractTokenFromHeader(authHeader)
		if err != nil {
			// log.Printf("AuthMiddleware: Ошибка извлечения токена: %v", err)
			http.Error(w, "Требуется авторизация", http.StatusUnauthorized)
			return
		}

		claims, err := auth.ValidateToken(tokenString)
		if err != nil {
			// log.Printf("AuthMiddleware: Ошибка валидации токена: %v", err)
			http.Error(w, "Неверный токен", http.StatusUnauthorized)
			return
		}

		// Добавляем информацию о пользователе в контекст
		ctx := context.WithValue(r.Context(), "user_id", claims.UserID)
		ctx = context.WithValue(ctx, "user_role", claims.Role)
		ctx = context.WithValue(ctx, "user_login", claims.Login)

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
