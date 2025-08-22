package main

import (
	"hat-store-training/backend/routes"
	"log"
	"net/http"
)

func main() {
	r := routes.InitializeRoutes()

	// Serve arquivos estáticos do frontend
	fs := http.FileServer(http.Dir("../frontend"))
	http.Handle("/", fs)
	http.Handle("/api/", http.StripPrefix("/api", r)) // API em /api

	log.Println("Server is running on port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
