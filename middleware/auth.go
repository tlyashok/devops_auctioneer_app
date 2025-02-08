package middleware

import (
	"context"
	"net/http"
	"strings"

	"auction-app/utils"
)

// Проверка авторизации
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Извлекаем заголовок авторизации
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Необходимо авторизоваться", http.StatusUnauthorized)
			return
		}

		// Парсим токен
		token := strings.TrimPrefix(authHeader, "Bearer ")
		userID, err := utils.ValidateToken(token) // Получаем только ID пользователя
		if err != nil {
			http.Error(w, "Неверный токен", http.StatusUnauthorized)
			return
		}

		// Сохраняем только userID в контексте запроса
		ctx := context.WithValue(r.Context(), "userID", userID)
		next.ServeHTTP(w, r.WithContext(ctx)) // Передаем управление следующему обработчику
	})
}
