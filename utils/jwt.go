package utils

import (
	"auction-app/models"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"time"
)

var secretKey = []byte("секретный_ключ")

type Claims struct {
	UserID int `json:"user_id"`
	jwt.StandardClaims
}

// Генерация токена
func GenerateToken(userID int) (string, error) {
	claims := Claims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

// Проверка токена
func ValidateToken(tokenStr string) (int, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil || !token.Valid {
		return 0, errors.New("Неверный токен")
	}

	// Декодируем данные пользователя из токена
	claims, ok := token.Claims.(*Claims)
	if !ok {
		return 0, errors.New("Ошибка декодирования токена")
	}

	// Получаем пользователя из базы данных
	user, err := models.GetUserByID(claims.UserID)
	if err != nil {
		return 0, err
	}

	return user.ID, nil
}
