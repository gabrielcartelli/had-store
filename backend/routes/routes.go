package routes

import (
	"hat-store-training/backend/handlers"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Middleware para logar cada request
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Recebida requisição: %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func InitializeRoutes() *mux.Router {
	router := mux.NewRouter()

	router.Use(loggingMiddleware)

	router.HandleFunc("/hats", handlers.GetHats).Methods("GET")
	router.HandleFunc("/cart/add", handlers.AddToCart).Methods("POST")
	router.HandleFunc("/cart/update", handlers.UpdateCart).Methods("PUT")
	router.HandleFunc("/checkout", handlers.Checkout).Methods("POST")
	router.HandleFunc("/pedido", handlers.RegistrarPedido).Methods("POST", "OPTIONS")

	return router
}
