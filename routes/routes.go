package routes

import (
	"auction-app/controllers"
	"auction-app/middleware"
	"net/http"
)

// ChainMiddleware - объединяет несколько мидлварей.
func ChainMiddleware(h http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h)
	}
	return h
}

func InitRoutes(mux *http.ServeMux) {
	// Открытые маршруты
	mux.Handle("/users/register", ChainMiddleware(
		http.HandlerFunc(controllers.CreateUser),
		middleware.LoggingMiddleware,
	))
	mux.Handle("/users/login", ChainMiddleware(
		http.HandlerFunc(controllers.LoginUser),
		middleware.LoggingMiddleware,
	))

	// Защищенные маршруты
	authenticated := []struct {
		pattern string
		handler http.HandlerFunc
	}{
		{"/auctions", controllers.AuctionsHandler},
		{"/auction/bid", controllers.CreateBid},
		{"/auctions/get", controllers.GetAuctionByIDHandler},
		{"/auctions/end", controllers.EndAuctionHandler},
		{"/auctions/delete", controllers.DeleteAuctionHandler},
		{"/users/me", controllers.GetUserInfo},
	}

	for _, route := range authenticated {
		mux.Handle(route.pattern, ChainMiddleware(
			route.handler,
			middleware.AuthMiddleware,
			middleware.LoggingMiddleware,
		))
	}
}
