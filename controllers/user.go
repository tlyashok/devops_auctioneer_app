package controllers

import (
	"auction-app/models"
	"auction-app/utils"
	"encoding/json"
	"net/http"
	"strings"
)

// Регистрация пользователя
func CreateUser(w http.ResponseWriter, r *http.Request) {
	var userRequest struct {
		Username string  `json:"username"`
		Password string  `json:"password"`
		Balance  float64 `json:"balance"`
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&userRequest); err != nil {
		utils.SendError(w, "Невалидные данные", http.StatusBadRequest)
		return
	}

	// Проверка, что пароль не пустой и имеет минимальную длину
	if strings.TrimSpace(userRequest.Password) == "" || len(userRequest.Password) < 6 {
		utils.SendError(w, "Пароль должен быть не менее 6 символов", http.StatusBadRequest)
		return
	}

	// Проверка уникальности имени пользователя
	existingUser, err := models.GetUserByUsername(userRequest.Username)
	if err != nil {
		utils.SendError(w, "Ошибка при проверке имени пользователя", http.StatusInternalServerError)
		return
	}
	if existingUser != nil {
		utils.SendError(w, "Пользователь с таким именем уже существует", http.StatusBadRequest)
		return
	}

	// Логика создания пользователя
	user, err := models.CreateUser(userRequest.Username, userRequest.Password, userRequest.Balance)
	if err != nil {
		utils.SendError(w, "Ошибка создания пользователя", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// Проверка существования пользователя
func userExists(username string) bool {
	// Логика для проверки существования пользователя в базе данных через модель
	user, err := models.GetUserByUsername(username)
	return user != nil && err == nil
}

// Авторизация пользователя
func LoginUser(w http.ResponseWriter, r *http.Request) {
	var loginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&loginRequest); err != nil {
		utils.SendError(w, "Невалидные данные", http.StatusBadRequest)
		return
	}

	user, err := models.AuthenticateUser(loginRequest.Username, loginRequest.Password)
	if err != nil {
		utils.SendError(w, "Неверное имя пользователя или пароль", http.StatusUnauthorized)
		return
	}

	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		utils.SendError(w, "Ошибка при генерации токена", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"user":  user,
		"token": token,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Получение пользователя по ID
func GetUserInfo(w http.ResponseWriter, r *http.Request) {
	// Извлекаем userID из контекста
	userID, ok := r.Context().Value("userID").(int)
	if !ok || userID == 0 {
		utils.SendError(w, "Пользователь не найден", http.StatusBadRequest)
		return
	}

	// Получаем данные пользователя из базы данных по userID
	user, err := models.GetUserByID(userID)
	if err != nil {
		utils.SendError(w, "Пользователь не найден: "+err.Error(), http.StatusNotFound)
		return
	}

	// Создание объекта с информацией для ответа
	userInfo := &models.UserInfo{
		ID:       user.ID,
		Username: user.Username,
		Balance:  user.Balance,
	}

	// Ответ в формате JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userInfo)
}
