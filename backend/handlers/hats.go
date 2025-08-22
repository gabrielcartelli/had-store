package handlers

import (
	"encoding/json"
	"net/http"
)

type Hat struct {
	ID    int     `json:"id"`
	Nome  string  `json:"nome"`
	Preco float64 `json:"preco"`
}

var hats = []Hat{
	{ID: 1, Nome: "Chapéu Panamá", Preco: 99.90},
	{ID: 2, Nome: "Chapéu Fedora", Preco: 89.90},
	{ID: 3, Nome: "Chapéu Bucket", Preco: 49.90},
	{ID: 4, Nome: "Chapéu Cowboy", Preco: 109.90},
	{ID: 5, Nome: "Chapéu Floppy", Preco: 79.90},
	{ID: 6, Nome: "Chapéu Bowler", Preco: 69.90},
	{ID: 7, Nome: "Chapéu Beanie", Preco: 39.90},
	{ID: 8, Nome: "Chapéu Pork Pie", Preco: 59.90},
	{ID: 9, Nome: "Chapéu Trilby", Preco: 84.90},
	{ID: 10, Nome: "Chapéu Snapback", Preco: 44.90},
	{ID: 11, Nome: "Chapéu Beret", Preco: 54.90},
	{ID: 12, Nome: "Chapéu Cloche", Preco: 64.90},
	{ID: 13, Nome: "Chapéu Top Hat", Preco: 119.90},
	{ID: 14, Nome: "Chapéu Sun Hat", Preco: 49.90},
	{ID: 15, Nome: "Chapéu Newsboy", Preco: 59.90},
	{ID: 16, Nome: "Chapéu Visor", Preco: 29.90},
	{ID: 17, Nome: "Chapéu Boater", Preco: 74.90},
	{ID: 18, Nome: "Chapéu Bucket Estampado", Preco: 54.90},
	{ID: 19, Nome: "Chapéu Aviador", Preco: 89.90},
	{ID: 20, Nome: "Chapéu Militar", Preco: 69.90},
	{ID: 21, Nome: "Chapéu Safari", Preco: 79.90},
	{ID: 22, Nome: "Chapéu Pescador", Preco: 39.90},
	{ID: 23, Nome: "Chapéu Trapper", Preco: 99.90},
	{ID: 24, Nome: "Chapéu Sombrero", Preco: 129.90},
	{ID: 25, Nome: "Chapéu Turbante", Preco: 59.90},
	{ID: 26, Nome: "Chapéu Balaclava", Preco: 49.90},
	{ID: 27, Nome: "Chapéu de Palha", Preco: 34.90},
	{ID: 28, Nome: "Chapéu de Feltro", Preco: 64.90},
	{ID: 29, Nome: "Chapéu de Lã", Preco: 54.90},
	{ID: 30, Nome: "Chapéu de Couro", Preco: 109.90},
}

func GetHats(w http.ResponseWriter, r *http.Request) {
	// Permitir requisições do frontend (CORS)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(hats)
}

// Os handlers abaixo são apenas exemplos para evitar erro de rota.
// Implemente a lógica real conforme sua necessidade.
func AddToCart(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Adicionado ao carrinho"))
}

func UpdateCart(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Carrinho atualizado"))
}

func Checkout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Compra finalizada"))
}
