package main

import (
	"auction-app/tasks"
	"log"
	"net/http"
	"os"

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

	// Создаём новый роутер
	mux := http.NewServeMux()

	// Обработчик статики с 404-логированием
	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("/", staticOr404(fs))

	// Инициализация маршрутов
	routes.InitRoutes(mux)

	// Фоновый обработчик статусов лотов
	go tasks.MonitorAuctions()

	// Запускаем сервер
	log.Println("Сервер запущен на порту 8000...")
	if err := http.ListenAndServe(":8000", mux); err != nil {
		log.Fatalf("Ошибка сервера: %v", err)
	}
}
