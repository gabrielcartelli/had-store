package routes

import (
	"hat-store-training/backend/handlers"
	"log"
	"net/http"
    "strings" 
    "github.com/golang-jwt/jwt/v5"
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

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1. Pega o cabeçalho de autorização 
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Cabeçalho de autorização ausente", http.StatusUnauthorized)
			return
		}

		// 2. O carimbo vem no formato "Bearer <token>", então pegamos só o token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// 3. Verifica se é válido e não é falso
		claims := &jwt.RegisteredClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Token inválido", http.StatusUnauthorized)
			return
		}

		// Se o for válido, efetua o login!
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
	router.HandleFunc("/api/pedidos", handlers.CriarPedido).Methods("POST")

	publicRouter := router.PathPrefix("/auth").Subrouter()
	publicRouter.HandleFunc("/register", handlers.RegisterHandler).Methods("POST")
	publicRouter.HandleFunc("/login", handlers.LoginHandler).Methods("POST")

	apiRouter := router.PathPrefix("/api").Subrouter()
	apiRouter.Use(uuidMiddleware)
	//apiRouter.Use(authMiddleware) // <<-- DESCOMENTE ESTA LINHA QUANDO O FRONTEND ESTIVER PRONTO!
	apiRouter.HandleFunc("/hats", handlers.GetHats).Methods("GET")
	apiRouter.HandleFunc("/pedido", handlers.CriarPedido).Methods("POST") // Usando a função correta
	apiRouter.HandleFunc("/pedidos", handlers.ConsultarPedidos).Methods("GET")

	

	return router
}
