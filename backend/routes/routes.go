package routes

import (
	"hat-store-training/backend/handlers"
	"log"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
)

// Chave secreta para os tokens. Precisa ser a mesma usada no handlers/auth.go
// O ideal seria compartilhar isso de um pacote de configuração, mas por simplicidade vamos redeclará-la.
var jwtKey = []byte("minha_chave_super_secreta")

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

// authMiddleware verifica o token JWT
func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Cabeçalho de autorização ausente", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims := &jwt.RegisteredClaims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Token inválido", http.StatusUnauthorized)
			return
		}
		
		next.ServeHTTP(w, r)
	})
}

func InitializeRoutes() *mux.Router {
	// 1. Roteador Principal
	router := mux.NewRouter()
	// Aplica o logging em TODAS as requisições
	router.Use(loggingMiddleware)

	// 2. Sub-roteador para rotas de AUTENTICAÇÃO (PÚBLICAS)
	// Caminho: /auth/...
	// Middlewares: Nenhum (além do logging)
	authRouter := router.PathPrefix("/auth").Subrouter()
	authRouter.HandleFunc("/register", handlers.RegisterHandler).Methods("POST")
	authRouter.HandleFunc("/login", handlers.LoginHandler).Methods("POST")

	// 3. Sub-roteador para rotas da API (PROTEGIDAS)
	// Caminho: /api/...
	// Middlewares: uuidMiddleware e authMiddleware (aplicados em sequência)
	apiRouter := router.PathPrefix("/api").Subrouter()
	apiRouter.Use(uuidMiddleware)
	apiRouter.Use(authMiddleware) 
	
	apiRouter.HandleFunc("/hats", handlers.GetHats).Methods("GET")
	apiRouter.HandleFunc("/pedido", handlers.CriarPedido).Methods("POST") // Rota correta para criar pedido
	apiRouter.HandleFunc("/pedidos", handlers.ConsultarPedidos).Methods("GET")

	return router
}