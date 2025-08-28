package main

import (
	_ "hat-store-training/backend/docs" // Importa os docs gerados pelo swag
	"hat-store-training/backend/routes"
	"log"
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Had Store API
// @version 1.0
// @description API da loja de chapéus Had Store.
// @host localhost:8080
// @BasePath /

// AVISO: Mudei o @BasePath para / porque o gorilla/mux agora controla a raiz.

func main() {
	// 1. Inicializa o nosso roteador gorilla/mux que já conhece TODAS as rotas (/api, /auth)
	r := routes.InitializeRoutes()

	// 2. Registra a rota da documentação Swagger DIRETAMENTE no roteador gorilla/mux
	// O PathPrefix garante que ele capture todos os subcaminhos como /swagger/index.html
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	// 3. Define a rota "catch-all" para servir os arquivos do frontend.
	// O gorilla/mux é inteligente: ele primeiro tenta casar as rotas mais específicas (/api, /auth, /swagger)
	// Se nenhuma casar, ele usa esta regra genérica.
	fs := http.FileServer(http.Dir("./frontend"))
	r.PathPrefix("/").Handler(fs)

	// 4. Inicia o servidor e diz a ele para usar o roteador 'r' para TODAS as requisições.
	// A MUDANÇA MAIS IMPORTANTE: Passamos 'r' em vez de 'nil' para ListenAndServe.
	log.Println("Servidor iniciado em http://localhost:8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}