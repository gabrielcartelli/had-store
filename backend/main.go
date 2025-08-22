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
// @BasePath /api

func main() {
	r := routes.InitializeRoutes()

	// Documentação Swagger acessível em /swagger/
	http.Handle("/swagger/", httpSwagger.WrapHandler)

	// Serve arquivos estáticos do frontend
	fs := http.FileServer(http.Dir("./frontend"))
	http.Handle("/", fs)
	http.Handle("/api/", http.StripPrefix("/api", r))

	log.Println("Server is running on port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
