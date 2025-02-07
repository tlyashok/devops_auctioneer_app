package repository

import (
	"database/sql"
	"log"

	"auction-app/config"
	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB(cfg *config.Config) {
	connStr := cfg.GetDBConnectionString()
	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Ошибка подключения к базе: %v", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatalf("Невозможно достучаться до базы: %v", err)
	}
}
