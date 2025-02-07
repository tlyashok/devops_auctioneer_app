package models

import (
	"time"

	"auction-app/repository"
)

const (
	StatusActive = iota
	StatusCompleted
)

func StatusToString(status int) string {
	switch status {
	case StatusActive:
		return "Active"
	case StatusCompleted:
		return "Completed"
	}
	return "Undefined"
}

type Auction struct {
	ID            int
	Title         string
	Description   string
	StartingPrice float64
	Status        int
	StartTime     time.Time
	EndTime       time.Time
	CreatorID     int
	WinnerID      int
	MaxBid        float64
}

type CreateAuctionRequest struct {
	Title         string  `json:"title"`
	Description   string  `json:"description"`
	Duration      int     `json:"duration"`       // продолжительность аукциона в минутах
	CreatorID     int     `json:"creator_id"`     // ID создателя аукциона
	StartingPrice float64 `json:"starting_price"` // Добавлено поле для максимальной ставки
}

type AuctionResponse struct {
	ID              int       `json:"id"`
	Title           string    `json:"title"`
	Description     string    `json:"description"`
	StartTime       time.Time `json:"start_time"`
	EndTime         time.Time `json:"end_time,omitempty"`
	CreatorUsername string    `json:"creator_username"`
	Status          string    `json:"status"`
	WinnerUsername  string    `json:"winner_username"`
	MaxBid          float64   `json:"max_bid"`
	StartingPrice   float64   `json:"starting_price"`
}

func CreateAuction(auc CreateAuctionRequest) (*Auction, error) {
	startTime := time.Now()
	endTime := startTime.Add(time.Duration(auc.Duration) * time.Minute)
	var id int
	err := repository.DB.QueryRow(
		`INSERT INTO auctions 
    	(title, description, starting_price, status, 
    	 start_time, end_time, creator_id, winner_id, max_bid)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NULL, $8) RETURNING id`,
		auc.Title, auc.Description, auc.StartingPrice, StatusActive,
		startTime, endTime, auc.CreatorID, auc.StartingPrice,
	).Scan(&id)
	if err != nil {
		return nil, err
	}
	return GetAuctionByID(id)
}

func GetAuctionByID(id int) (*Auction, error) {
	var auction Auction
	query := `SELECT id, title, description, start_time, end_time, creator_id, 
              COALESCE(winner_id, -1) AS winner_id, starting_price 
              FROM auctions WHERE id = $1`
	err := repository.DB.QueryRow(query, id).Scan(
		&auction.ID, &auction.Title, &auction.Description, &auction.StartTime,
		&auction.EndTime, &auction.CreatorID, &auction.WinnerID, &auction.StartingPrice,
	)
	if err != nil {
		return nil, err
	}
	return &auction, nil
}

func GetAllAuctions() ([]Auction, error) {
	rows, err := repository.DB.Query(
		`SELECT 
    	 id, title, description, start_time, 
    	 end_time, creator_id, status, COALESCE(winner_id, -1), max_bid, starting_price  
		 FROM auctions`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var auctions []Auction
	for rows.Next() {
		var a Auction
		err := rows.Scan(
			&a.ID,
			&a.Title,
			&a.Description,
			&a.StartTime,
			&a.EndTime,
			&a.CreatorID,
			&a.Status,
			&a.WinnerID,
			&a.MaxBid,
			&a.StartingPrice,
		)
		if err != nil {
			return nil, err
		}
		auctions = append(auctions, a)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return auctions, nil
}

func GetMaxBid(auctionID int) (float64, error) {
	var maxBid float64
	err := repository.DB.QueryRow("SELECT MAX(amount) FROM bids WHERE auction_id = $1", auctionID).
		Scan(&maxBid)
	if err != nil {
		return 0, err
	}
	return maxBid, nil
}

func EndAuction(auctionID int) error {
	// Получаем максимальную ставку
	maxBid, err := GetMaxBid(auctionID)
	if err != nil && maxBid != 0 {
		return err
	}
	// Обновляем аукцион, устанавливая победителя, максимальную ставку и статус
	_, err = repository.DB.Exec(
		`UPDATE auctions SET status = $1 WHERE id = $2`,
		StatusCompleted, auctionID,
	)
	return err
}

func DeleteAuction(auctionID int) error {
	_, err := repository.DB.Exec(
		"DELETE FROM auctions WHERE id = $1",
		auctionID,
	)
	return err
}
