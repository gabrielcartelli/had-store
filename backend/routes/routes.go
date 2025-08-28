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
var jwtKey = []byte("e3e6c6c2-9b7d-4c5e-8c1a-2f7b8f8e2a1d")

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
			log.Printf("[ERRO] Autenticação falhou: Cabeçalho de autorização ausente em %s %s", r.Method, r.URL.Path)
			http.Error(w, "Cabeçalho de autorização ausente", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims := &jwt.RegisteredClaims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			log.Printf("[ERRO] Token inválido: %v | Rota: %s %s", err, r.Method, r.URL.Path)
			http.Error(w, "Token inválido", http.StatusUnauthorized)
			return
		}

		log.Printf("[INFO] Autenticação bem-sucedida para rota %s %s", r.Method, r.URL.Path)
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

	// 3. Sub-roteador para rotas da API públicas
	apiPublic := router.PathPrefix("/api").Subrouter()
	apiPublic.HandleFunc("/hats", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[INFO] Requisição recebida em /api/hats: %s %s", r.Method, r.URL.Path)
		handlers.GetHats(w, r)
	}).Methods("GET")

	// 4. Sub-roteador para rotas da API protegidas
	apiProtected := router.PathPrefix("/api").Subrouter()
	apiProtected.Use(authMiddleware)
	apiProtected.HandleFunc("/pedido", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[INFO] Requisição recebida em /api/pedido: %s %s", r.Method, r.URL.Path)
		handlers.RegistrarPedido(w, r)
	}).Methods("POST")
	apiProtected.HandleFunc("/pedidos", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[INFO] Requisição recebida em /api/pedidos: %s %s", r.Method, r.URL.Path)
		handlers.ConsultarPedidos(w, r)
	}).Methods("GET")

	return router
}
