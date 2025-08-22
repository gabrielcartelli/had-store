package main

import (
	"hat-store-training/backend/routes"
	"log"
	"net/http"
)

func main() {
	r := routes.InitializeRoutes()

	log.Println("Server is running on port 8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
