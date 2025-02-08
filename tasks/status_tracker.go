package tasks

import (
	"auction-app/models"
	"auction-app/repository"
	"time"
)

func MonitorAuctions() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		rows, err := repository.DB.Query(
			"SELECT id FROM auctions WHERE status = $1 AND end_time <= $2",
			models.StatusActive, time.Now(),
		)
		if err != nil {
			continue
		}
		defer rows.Close()

		for rows.Next() {
			var auctionID int
			if err := rows.Scan(&auctionID); err == nil {
				_ = models.EndAuction(auctionID) // Завершаем аукцион
			}
		}
	}
}
