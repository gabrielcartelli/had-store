package routes

import (
	"hat-store-training/backend/handlers"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

const apiUUID = "e3e6c6c2-9b7d-4c5e-8c1a-2f7b8f8e2a1d"

// Middleware para validar o UUID no header
func uuidMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		uuid := r.Header.Get("X-API-UUID")
		if uuid != apiUUID {
			http.Error(w, "Acesso não autorizado", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

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
	router.Use(uuidMiddleware)

	router.HandleFunc("/hats", handlers.GetHats).Methods("GET")
	router.HandleFunc("/pedido", handlers.RegistrarPedido).Methods("POST", "OPTIONS")
	router.HandleFunc("/pedidos", handlers.ConsultarPedidos).Methods("GET")

	return router
}
