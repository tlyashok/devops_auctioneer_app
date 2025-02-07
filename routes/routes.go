package routes

import (
	"auction-app/controllers"
	"auction-app/middleware"
	"log"
	"net/http"
	"time"
)

func logRequest(r *http.Request, statusCode int) {
	log.Printf("%s %s %s %d %s",
		time.Now().Format(time.RFC1123),
		r.Method,
		r.URL.Path,
		statusCode,
		r.UserAgent(),
	)
}

func InitRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/users/register", logRequestHandler(controllers.CreateUser)) // POST /users/reg
	mux.HandleFunc("/users/login", logRequestHandler(controllers.LoginUser))     // POST /users/login

	// Защищенные маршруты, которые требуют авторизации
	mux.Handle("/auctions", middleware.AuthMiddleware(http.HandlerFunc(auctionsHandler)))                         // POST /auctions, GET список аукционов
	mux.Handle("/auction/bid", middleware.AuthMiddleware(http.HandlerFunc(controllers.CreateBid)))                // POST /auction/bid?auction_id=1
	mux.Handle("/auctions/get", middleware.AuthMiddleware(http.HandlerFunc(controllers.GetAuctionByIDHandler)))   // GET /auctions/get?id=1
	mux.Handle("/auctions/end", middleware.AuthMiddleware(http.HandlerFunc(controllers.EndAuctionHandler)))       // POST /auctions/end?id=1
	mux.Handle("/auctions/delete", middleware.AuthMiddleware(http.HandlerFunc(controllers.DeleteAuctionHandler))) // DELETE /auctions/delete?id=1
	mux.Handle("/users/me", middleware.AuthMiddleware(http.HandlerFunc(controllers.GetUserInfo)))
}

func logRequestHandler(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rr := &responseRecorder{w, http.StatusOK}
		next.ServeHTTP(rr, r)
		logRequest(r, rr.statusCode)
	}
}

type responseRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (r *responseRecorder) WriteHeader(statusCode int) {
	r.statusCode = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

func auctionsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		controllers.CreateAuction(w, r)
	case http.MethodGet:
		controllers.GetAuctions(w, r)
	default:
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
	}
}
