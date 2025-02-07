package controllers

import (
	"auction-app/models"
	"auction-app/utils"
	"encoding/json"
	"net/http"
	"strconv"
)

// Создание аукциона: POST /auctions
func CreateAuction(w http.ResponseWriter, r *http.Request) {
	var createAuction models.CreateAuctionRequest
	if err := json.NewDecoder(r.Body).Decode(&createAuction); err != nil {
		utils.SendError(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	createAuction.CreatorID = r.Context().Value("userID").(int)

	auction, err := models.CreateAuction(createAuction)
	if err != nil {
		utils.SendError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(auction)
}

func CreateBid(w http.ResponseWriter, r *http.Request) {
	// Декодирование входных данных
	var input struct {
		AuctionID int     `json:"auction_id"` // Правильное имя для аукциона
		Amount    float64 `json:"amount"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.SendError(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		utils.SendError(w, "Не удалось определить пользователя по ID", http.StatusUnauthorized)
		return
	}
	// Проверка наличия пользователя
	user, err := models.GetUserByID(userID)
	if err != nil {
		utils.SendError(w, "Пользователь не найден", http.StatusNotFound)
		return
	}

	// Проверка баланса
	if user.Balance < input.Amount {
		utils.SendError(w, "Недостаточно средств", http.StatusBadRequest)
		return
	}

	// Получаем последнюю ставку, если она есть
	previousUserID, previousAmount, err := models.GetLastBid(input.AuctionID)
	if err != nil && err.Error() != "sql: no rows in result set" { // Если ставка была, то ошибку можно игнорировать
		utils.SendError(w, "Ошибка при получении последней ставки", http.StatusInternalServerError)
		return
	}

	// Если предыдущий пользователь существует и ставка перебита, возвращаем ему деньги
	if previousUserID != 0 && previousAmount > 0 && previousAmount != input.Amount {
		// Обновление баланса предыдущего пользователя
		previousUser, err := models.GetUserByID(previousUserID)
		if err != nil {
			utils.SendError(w, "Ошибка при получении информации о предыдущем пользователе", http.StatusInternalServerError)
			return
		}

		// Обновляем баланс предыдущего пользователя
		previousNewBalance := previousUser.Balance + previousAmount
		if err := models.UpdateUserBalance(previousUserID, previousNewBalance); err != nil {
			utils.SendError(w, "Ошибка при обновлении баланса предыдущего пользователя", http.StatusInternalServerError)
			return
		}
	}

	// Проверка наличия пользователя
	user, err = models.GetUserByID(userID)

	// Обновление баланса текущего пользователя
	newBalance := user.Balance - input.Amount
	if err := models.UpdateUserBalance(user.ID, newBalance); err != nil {
		utils.SendError(w, "Ошибка при обновлении баланса", http.StatusInternalServerError)
		return
	}

	// Создание ставки
	bid, err := models.CreateBid(input.AuctionID, userID, input.Amount)
	if err != nil {
		utils.SendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Ответ с информацией о ставке
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bid)
}

// Получение списка аукционов: GET /auctions
func GetAuctions(w http.ResponseWriter, r *http.Request) {
	auctions, err := models.GetAllAuctions()
	if err != nil {
		utils.SendError(w, "Ошибка при получении списка аукционов", http.StatusInternalServerError)
		return
	}

	auctionsResponse := make([]models.AuctionResponse, len(auctions))

	for i, auction := range auctions {
		creator, err := models.GetUserByID(auction.CreatorID)
		if err != nil {
			creator = &models.User{
				Username: "Unknown",
			}
		}
		winner, err := models.GetUserByID(auction.WinnerID)
		if err != nil {
			winner = &models.User{
				Username: "Unknown",
			}
		}

		auctionsResponse[i] = models.AuctionResponse{
			ID:              auction.ID,
			Title:           auction.Title,
			Description:     auction.Description,
			StartTime:       auction.StartTime,
			EndTime:         auction.EndTime,
			CreatorUsername: creator.Username,
			Status:          models.StatusToString(auction.Status),
			WinnerUsername:  winner.Username,
			MaxBid:          auction.MaxBid,
			StartingPrice:   auction.StartingPrice,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(auctionsResponse)
}

// Получение аукциона по ID: GET /auctions/{id}
func GetAuctionByIDHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.SendError(w, "Неверный ID аукциона", http.StatusBadRequest)
		return
	}

	auction, err := models.GetAuctionByID(id)
	if err != nil {
		utils.SendError(w, err.Error(), http.StatusNotFound)
		return
	}

	// Получаем текущую максимальную ставку для аукциона
	maxBid, err := models.GetMaxBid(id)
	if err != nil {
		utils.SendError(w, "Ошибка при получении максимальной ставки", http.StatusInternalServerError)
		return
	}

	// Добавляем максимальную ставку в ответ
	auction.MaxBid = maxBid

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(auction)
}

// Завершение аукциона: POST /auctions/{id}/end
func EndAuctionHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.SendError(w, "Неверный ID аукциона", http.StatusBadRequest)
		return
	}

	err = models.EndAuction(id)
	if err != nil {
		utils.SendError(w, "Ошибка при завершении аукциона", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Аукцион завершен"))
}

// Удаление аукциона: DELETE /auctions/{id}
func DeleteAuctionHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.SendError(w, "Неверный ID аукциона", http.StatusBadRequest)
		return
	}

	err = models.DeleteAuction(id)
	if err != nil {
		utils.SendError(w, "Ошибка при удалении аукциона", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Аукцион удален"))
}
