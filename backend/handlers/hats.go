package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"sync"
)

// Estrutura do chapéu (mantida)
type Hat struct {
	ID    int     `json:"id"`
	Nome  string  `json:"nome"`
	Price float64 `json:"price"`
}

// Estrutura do pedido
type Pedido struct {
	Nome      string      `json:"nome"`
	CPF       string      `json:"cpf"`
	Email     string      `json:"email"`
	Pagamento string      `json:"pagamento"`
	Itens     []HatPedido `json:"itens"`
	Total     float64     `json:"total"`
	Cupom     string      `json:"cupom"`
}

// Estrutura dos itens do pedido
type HatPedido struct {
	ID         int     `json:"id"`
	Nome       string  `json:"nome"`
	Price      float64 `json:"price"`
	Quantidade int     `json:"quantidade"`
}

// Variável global para armazenar chapéus (mantida)
var hats = []Hat{
	{ID: 1, Nome: "Chapéu Panamá", Price: 120.00},
	{ID: 2, Nome: "Chapéu Fedora", Price: 150.00},
	{ID: 3, Nome: "Chapéu Bucket", Price: 49.90},
	{ID: 4, Nome: "Chapéu Cowboy", Price: 109.90},
	{ID: 5, Nome: "Chapéu Floppy", Price: 79.90},
	{ID: 6, Nome: "Chapéu Bowler", Price: 69.90},
	{ID: 7, Nome: "Chapéu Beanie", Price: 39.90},
	{ID: 8, Nome: "Chapéu Pork Pie", Price: 59.90},
	{ID: 9, Nome: "Chapéu Trilby", Price: 84.90},
	{ID: 10, Nome: "Chapéu Snapback", Price: 44.90},
	{ID: 11, Nome: "Chapéu Sertanejo", Price: 99.90},
	{ID: 12, Nome: "Chapéu Gaúcho", Price: 129.90},
	{ID: 13, Nome: "Chapéu Cangaceiro", Price: 139.90},
	{ID: 14, Nome: "Chapéu Nordestino", Price: 34.90},
	{ID: 15, Nome: "Chapéu de Pescador", Price: 29.90},
	{ID: 16, Nome: "Chapéu de Festa Junina", Price: 19.90},
	{ID: 17, Nome: "Chapéu de Lavrador", Price: 27.90},
	{ID: 18, Nome: "Chapéu de Vaqueiro", Price: 89.90},
	{ID: 19, Nome: "Chapéu de Tropeiro", Price: 69.90},
	{ID: 20, Nome: "Chapéu Gustavo Carvalho", Price: 60.00},
}

// Variável global para armazenar pedidos
var (
	pedidos     []Pedido
	pedidosLock sync.Mutex
)

// Handler para listar chapéus (mantido)
// GetHats godoc
// @Summary Lista todos os chapéus
// @Description Retorna a lista de chapéus disponíveis
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
// @Param pedido body Pedido true "Dados do pedido"
// @Success 200 {object} map[string]string
// @Router /pedido [post]
func RegistrarPedido(w http.ResponseWriter, r *http.Request) {
	var pedido Pedido
	err := json.NewDecoder(r.Body).Decode(&pedido)
	ip := strings.Split(r.RemoteAddr, ":")[0]
	user := r.Header.Get("X-User-Email")
	if err != nil {
		if user != "" {
			log.Printf("[ERROR][%s][%s] Pedido inválido: %v", ip, user, err)
		} else {
			log.Printf("[ERROR][%s] Pedido inválido: %v", ip, err)
		}
		http.Error(w, "Dados inválidos", http.StatusBadRequest)
		return
	}

	// Calcular o total do pedido
	total := 0.0
	for _, item := range pedido.Itens {
		total += item.Price * float64(item.Quantidade)
	}

	// Verifica se é o primeiro pedido do CPF
	primeiroPedido := true
	pedidosLock.Lock()
	for _, p := range pedidos {
		if p.CPF == pedido.CPF {
			primeiroPedido = false
			break
		}
	}
	pedidosLock.Unlock()

	// Aplica desconto HATOFF só se for o primeiro pedido desse CPF
	if pedido.Cupom == "HATOFF" && primeiroPedido {
		total = total * 0.80
	}
	pedido.Total = total

	pedidosLock.Lock()
	pedidos = append(pedidos, pedido)
	pedidosLock.Unlock()

	if user != "" {
		log.Printf("[INFO][%s][%s] Pedido registrado | Nome: %s | CPF: %s | Total: %.2f | Pagamento: %s | Itens: %d", ip, user, pedido.Nome, pedido.CPF, pedido.Total, pedido.Pagamento, len(pedido.Itens))
	} else {
		log.Printf("[INFO][%s] Pedido registrado | Nome: %s | CPF: %s | Total: %.2f | Pagamento: %s | Itens: %d", ip, pedido.Nome, pedido.CPF, pedido.Total, pedido.Pagamento, len(pedido.Itens))
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"mensagem": "Pedido registrado com sucesso",
	})
}

// Handler para consultar pedidos por CPF
func ConsultarPedidos(w http.ResponseWriter, r *http.Request) {
	ip := strings.Split(r.RemoteAddr, ":")[0]
	user := r.Header.Get("X-User-Email")
	cpf := r.URL.Query().Get("cpf")
	if cpf == "" {
		if user != "" {
			log.Printf("[WARN][%s][%s] Consulta de pedidos sem CPF informado", ip, user)
		} else {
			log.Printf("[WARN][%s] Consulta de pedidos sem CPF informado", ip)
		}
		http.Error(w, "CPF não informado", http.StatusBadRequest)
		return
	}

	pedidosLock.Lock()
	defer pedidosLock.Unlock()

	var pedidosFiltrados []Pedido
	for _, pedido := range pedidos {
		if pedido.CPF == cpf {
			pedidosFiltrados = append(pedidosFiltrados, pedido)
		}
	}

	if user != "" {
		log.Printf("[INFO][%s][%s] Consulta de pedidos | CPF: %s | Quantidade: %d", ip, user, cpf, len(pedidosFiltrados))
	} else {
		log.Printf("[INFO][%s] Consulta de pedidos | CPF: %s | Quantidade: %d", ip, cpf, len(pedidosFiltrados))
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pedidosFiltrados)
}
