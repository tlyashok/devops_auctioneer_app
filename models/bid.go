package models

import (
	"database/sql"
	"errors"
	"time"

	"auction-app/repository"
)

type Bid struct {
	ID        int
	AuctionID int
	UserID    int
	Amount    float64
	BidTime   time.Time
}

// Создание ставки
func CreateBid(auctionID, userID int, amount float64) (*Bid, error) {
	// Получаем текущую максимальную ставку
	var currentMaxBid float64
	var status int
	err := repository.DB.QueryRow(
		"SELECT max_bid, status FROM auctions WHERE id = $1",
		auctionID,
	).Scan(&currentMaxBid, &status)
	if err != nil {
		return nil, err
	}

	// Проверяем, что новая ставка выше текущей максимальной
	if amount <= currentMaxBid {
		return nil, errors.New("Cтавка должна быть выше текущей максимальной")
	}

	if status != StatusActive {
		return nil, errors.New("Данный лот уже продан")
	}

	// Вставляем новую ставку в таблицу bids
	var id int
	err = repository.DB.QueryRow(
		`INSERT INTO bids (auction_id, user_id, amount)
		VALUES ($1, $2, $3) RETURNING id`,
		auctionID, userID, amount,
	).Scan(&id)
	if err != nil {
		return nil, err
	}

	// Обновляем max_bid и winner_id в таблице auctions
	_, err = repository.DB.Exec(
		"UPDATE auctions SET max_bid = $1, winner_id = $2 WHERE id = $3",
		amount, userID, auctionID,
	)
	if err != nil {
		return nil, err
	}

	return GetBidByID(id)
}

func GetBidByID(id int) (*Bid, error) {
	b := &Bid{}
	err := repository.DB.QueryRow(
		`SELECT id, auction_id, user_id, amount, bid_time FROM bids WHERE id = $1`,
		id,
	).Scan(&b.ID, &b.AuctionID, &b.UserID, &b.Amount, &b.BidTime)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("Ставка не найдена")
		}
		return nil, err
	}
	return b, nil
}

func GetLastBid(auctionID int) (int, float64, error) {
	var userID int
	var amount float64
	err := repository.DB.QueryRow(
		`SELECT user_id, amount FROM bids WHERE auction_id = $1 ORDER BY amount DESC LIMIT 1`, // OFFSET 1 для пропуска максимальной ставки
		auctionID,
	).Scan(&userID, &amount)
	if err != nil {
		return 0, 0, err
	}
	return userID, amount, nil
}
