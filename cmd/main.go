package main

import (
	"auction-app/tasks"
	"auction-app/utils"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"auction-app/config"
	"auction-app/repository"
	"auction-app/routes"
)

// Обработчик для статики + логирование 404
func staticOr404(fs http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Проверяем, существует ли запрашиваемый файл
		filePath := "./static" + r.URL.Path
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			log.Printf("404: %s %s %s", r.Method, r.URL.Path, r.UserAgent())
			http.NotFound(w, r)
			return
		}
		fs.ServeHTTP(w, r) // Если файл есть, отдаём его
	}
}

func main() {
	// Загружаем конфигурацию
	cfg := config.LoadConfig()

	// Инициализируем БД
	repository.InitDB(cfg)
	defer repository.DB.Close()

	// Загружаем секретный ключ для jwt токенов
	utils.InitSecretKey(cfg)

	// Создаём новый роутер
	mux := http.NewServeMux()

	// Обработчик статики с 404-логированием
	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("/", staticOr404(fs))

	// Инициализация маршрутов
	routes.InitRoutes(mux)

	// Фоновый обработчик статусов лотов
	go tasks.MonitorAuctions()

	log.Println("Сервер запущен на порту", 8000)
	port, _ := strconv.Atoi(cfg.AppPort)
	if port == 0 {
		port = 8080 // Устанавливаем порт по умолчанию, если ошибка
	}

	// Запуск сервера
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), mux); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %v\n", err)
	}
}
