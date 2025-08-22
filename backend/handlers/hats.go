package handlers

import (
	"encoding/json"
	"net/http"
	"sync"
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

var pedidos []map[string]interface{}
var pedidosMutex sync.Mutex

// GetHats godoc
// @Summary Lista todos os chapéus
// @Description Retorna todos os chapéus disponíveis
// @Tags hats
// @Produce json
// @Success 200 {array} models.Hat
// @Router /hats [get]
func GetHats(w http.ResponseWriter, r *http.Request) {
	// Permitir requisições do frontend (CORS)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(hats)
}

// AddToCart godoc
// @Summary Adiciona um chapéu ao carrinho
// @Description Adiciona um chapéu ao carrinho do usuário
// @Tags cart
// @Accept json
// @Produce json
// @Param item body models.Hat true "Item do carrinho"
// @Success 200 {object} map[string]interface{}
// @Router /cart/add [post]
func AddToCart(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Adicionado ao carrinho"))
}

// UpdateCart godoc
// @Summary Atualiza o carrinho
// @Description Atualiza a quantidade de um item no carrinho
// @Tags cart
// @Accept json
// @Produce json
// @Param item body models.Hat true "Item do carrinho"
// @Success 200 {object} map[string]interface{}
// @Router /cart/update [put]
func UpdateCart(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Carrinho atualizado"))
}

// Checkout godoc
// @Summary Finaliza o pedido
// @Description Finaliza o pedido do usuário
// @Tags checkout
// @Accept json
// @Produce json
// @Param pedido body map[string]interface{} true "Dados do pedido"
// @Success 200 {object} map[string]interface{}
// @Router /checkout [post]
func Checkout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Compra finalizada"))
}

// RegistrarPedido godoc
// @Summary Registra um pedido
// @Description Registra os dados do pedido em memória
// @Tags pedido
// @Accept json
// @Produce json
// @Param pedido body map[string]interface{} true "Dados do pedido"
// @Success 200 {object} map[string]string
// @Router /pedido [post]
func RegistrarPedido(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method == "OPTIONS" {
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.WriteHeader(http.StatusOK)
		return
	}
	var pedido map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&pedido); err != nil {
		http.Error(w, "Pedido inválido", http.StatusBadRequest)
		return
	}
	pedidosMutex.Lock()
	pedidos = append(pedidos, pedido)
	pedidosMutex.Unlock()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
