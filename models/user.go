package models

import (
	"auction-app/repository"
	"database/sql"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type User struct {
	ID        int
	Username  string
	Balance   float64
	Password  string
	CreatedAt time.Time
}

type UserInfo struct {
	ID       int     `json:"id"`
	Username string  `json:"username"`
	Balance  float64 `json:"balance"`
}

// Хеширование пароля
func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// Проверка пароля
func checkPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

// Создание пользователя
func CreateUser(username, password string, balance float64) (*User, error) {
	hashedPassword, err := hashPassword(password)
	if err != nil {
		return nil, err
	}

	var id int
	err = repository.DB.QueryRow(
		"INSERT INTO users (username, password, balance) VALUES ($1, $2, $3) RETURNING id",
		username, hashedPassword, balance,
	).Scan(&id)
	if err != nil {
		return nil, err
	}
	return GetUserByID(id)
}

// Получение пользователя по ID
func GetUserByID(id int) (*User, error) {
	u := &User{}
	err := repository.DB.QueryRow(
		"SELECT id, username, password, balance, created_at FROM users WHERE id = $1",
		id,
	).Scan(&u.ID, &u.Username, &u.Password, &u.Balance, &u.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("пользователь не найден")
		}
		return nil, err
	}
	return u, nil
}

// Получение пользователя по имени пользователя
func GetUserByUsername(username string) (*User, error) {
	u := &User{}
	err := repository.DB.QueryRow(
		"SELECT id, username, password, balance, created_at FROM users WHERE username = $1",
		username,
	).Scan(&u.ID, &u.Username, &u.Password, &u.Balance, &u.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Пользователь не найден
		}
		return nil, err
	}
	return u, nil
}

// Обновление баланса
func UpdateUserBalance(userID int, newBalance float64) error {
	_, err := repository.DB.Exec("UPDATE users SET balance = $1 WHERE id = $2", newBalance, userID)
	return err
}

// Авторизация пользователя
func AuthenticateUser(username, password string) (*User, error) {
	u := &User{}
	err := repository.DB.QueryRow(
		"SELECT id, username, password, balance, created_at FROM users WHERE username = $1",
		username,
	).Scan(&u.ID, &u.Username, &u.Password, &u.Balance, &u.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("Пользователь не найден")
		}
		return nil, err
	}

	if !checkPassword(u.Password, password) {
		return nil, errors.New("Неверный пароль")
	}

	return u, nil
}
