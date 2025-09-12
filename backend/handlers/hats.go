package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/gorilla/mux"
)

// Verifica UUID de desenvolvimento
func checkDevAuth(r *http.Request) bool {
	uuid := r.Header.Get("X-Dev-UUID")
	expected := os.Getenv("DEV_UUID")
	return uuid != "" && uuid == expected
}

// parseFloat faz o parse seguro de string para float64

// ListarEstoque godoc
// @Summary Lista o estoque de chapéus
// @Description Retorna o estoque atual de cada chapéu
// @Tags estoque
// @Produce json
// @Param X-Dev-UUID header string false "UUID de desenvolvimento"
// @Success 200 {array} Hat
// @Router /estoque [get]
func ListarEstoque(w http.ResponseWriter, r *http.Request) {
	if os.Getenv("ENVIRONMENT") == "development" && !checkDevAuth(r) {
		http.Error(w, "Acesso não autorizado: forneça o UUID", http.StatusUnauthorized)
		return
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	type EstoqueHat struct {
		ID         int    `json:"id"`
		Nome       string `json:"nome"`
		Quantidade int    `json:"quantidade"`
	}
	estoque := make([]EstoqueHat, len(hats))
	for i, h := range hats {
		estoque[i] = EstoqueHat{
			ID:         h.ID,
			Nome:       h.Nome,
			Quantidade: h.Quantidade,
		}
	}
	json.NewEncoder(w).Encode(estoque)
}

// Estrutura do chapéu (mantida)
type Hat struct {
	ID         int     `json:"id"`
	Nome       string  `json:"nome"`
	Price      float64 `json:"price"`
	Quantidade int     `json:"quantidade"`
	Categoria  string  `json:"categoria"` // nacional, importado, crescer
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
	{ID: 1, Nome: "Chapéu Panamá", Price: 120.00, Quantidade: 100, Categoria: "importado"},
	{ID: 2, Nome: "Chapéu Fedora", Price: 150.00, Quantidade: 80, Categoria: "importado"},
	{ID: 3, Nome: "Chapéu Bucket", Price: 49.90, Quantidade: 150, Categoria: "importado"},
	{ID: 4, Nome: "Chapéu Cowboy", Price: 109.90, Quantidade: 5, Categoria: "importado"},
	{ID: 5, Nome: "Chapéu Floppy", Price: 79.90, Quantidade: 120, Categoria: "importado"},
	{ID: 6, Nome: "Chapéu Bowler", Price: 69.90, Quantidade: 7, Categoria: "importado"},
	{ID: 7, Nome: "Chapéu Beanie", Price: 39.90, Quantidade: 20, Categoria: "importado"},
	{ID: 8, Nome: "Chapéu Pork Pie", Price: 59.90, Quantidade: 0, Categoria: "importado"},
	{ID: 9, Nome: "Chapéu Trilby", Price: 84.90, Quantidade: 9, Categoria: "importado"},
	{ID: 10, Nome: "Chapéu Snapback", Price: 44.90, Quantidade: 0, Categoria: "nacional"},
	{ID: 11, Nome: "Chapéu Sertanejo", Price: 99.90, Quantidade: 110, Categoria: "nacional"},
	{ID: 12, Nome: "Chapéu Gaúcho", Price: 129.90, Quantidade: 130, Categoria: "nacional"},
	{ID: 13, Nome: "Chapéu Cangaceiro", Price: 139.90, Quantidade: 40, Categoria: "nacional"},
	{ID: 14, Nome: "Chapéu de Pescador", Price: 29.90, Quantidade: 16, Categoria: "nacional"},
	{ID: 15, Nome: "Chapéu Gustavo Carvalho", Price: 60.00, Quantidade: 1000, Categoria: "crescer"},
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
// @Param X-Dev-UUID header string false "UUID de desenvolvimento"
// @Success 200 {array} models.Hat
// @Router /hats [get]
func GetHats(w http.ResponseWriter, r *http.Request) {
	// Proteção por UUID no ambiente de desenvolvimento
	if os.Getenv("ENVIRONMENT") == "development" && !checkDevAuth(r) {
		http.Error(w, "Acesso não autorizado: forneça o UUID", http.StatusUnauthorized)
		return
	}
	// Permitir requisições do frontend (CORS)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	// Adiciona flag de estoque antes de retornar
	type HatComEstoque struct {
		ID         int     `json:"id"`
		Nome       string  `json:"nome"`
		Price      float64 `json:"price"`
		Quantidade int     `json:"quantidade"`
		TemEstoque bool    `json:"temEstoque"`
		Categoria  string  `json:"categoria"`
	}
	// Filtro por categoria via query param: ?categoria=nacional,importado,crescer
	categorias := r.URL.Query().Get("categoria")
	var filtro []string
	if categorias != "" {
		filtro = strings.Split(categorias, ",")
	}
	// Filtro por faixa de valor
	minStr := r.URL.Query().Get("min")
	maxStr := r.URL.Query().Get("max")
	var min, max float64
	var errMin, errMax error
	if minStr != "" {
		min, errMin = parseFloat(minStr)
	}
	if maxStr != "" {
		max, errMax = parseFloat(maxStr)
	}
	hatsComEstoque := make([]HatComEstoque, 0)
	for _, h := range hats {
		// Filtro de categoria
		if len(filtro) > 0 && !containsCategoria(filtro, h.Categoria) {
			continue
		}
		// Filtro de valor mínimo
		if minStr != "" && errMin == nil && h.Price < min {
			continue
		}
		// Filtro de valor máximo
		if maxStr != "" && errMax == nil && h.Price > max {
			continue
		}
		hatsComEstoque = append(hatsComEstoque, HatComEstoque{
			ID:         h.ID,
			Nome:       h.Nome,
			Price:      h.Price,
			Quantidade: h.Quantidade,
			TemEstoque: h.Quantidade > 0,
			Categoria:  h.Categoria,
		})
	}
	json.NewEncoder(w).Encode(hatsComEstoque)
}

// parseFloat faz o parse seguro de string para float64
func parseFloat(s string) (float64, error) {
	return stringsToFloat(s)
}

// stringsToFloat converte string para float64 usando padrão brasileiro e internacional
func stringsToFloat(s string) (float64, error) {
	s = strings.ReplaceAll(s, ",", ".")
	return strconv.ParseFloat(s, 64)
}

func containsCategoria(filtros []string, categoria string) bool {
	for _, f := range filtros {
		if strings.ToLower(strings.TrimSpace(f)) == categoria {
			return true
		}
	}
	return false
}

// AddToCart godoc
// @Summary Adiciona um chapéu ao carrinho
// @Description Adiciona um chapéu ao carrinho do usuário
// @Tags cart
// @Accept json
// @Produce json
// @Param Authorization header string false "JWT token"
// @Param item body models.Hat true "Item do carrinho"
// @Success 200 {object} map[string]interface{}
// @Router /cart/add [post]
func AddToCart(w http.ResponseWriter, r *http.Request) {
	if os.Getenv("ENVIRONMENT") == "development" && !checkDevAuth(r) {
		http.Error(w, "Acesso não autorizado: forneça o UUID", http.StatusUnauthorized)
		return
	}
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
// @Param Authorization header string false "JWT token"
// @Param item body models.Hat true "Item do carrinho"
// @Success 200 {object} map[string]interface{}
// @Router /cart/update [put]
func UpdateCart(w http.ResponseWriter, r *http.Request) {
	if os.Getenv("ENVIRONMENT") == "development" && !checkDevAuth(r) {
		http.Error(w, "Acesso não autorizado: forneça o UUID", http.StatusUnauthorized)
		return
	}
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
// @Param Authorization header string false "JWT token"
// @Param pedido body map[string]interface{} true "Dados do pedido"
// @Success 200 {object} map[string]interface{}
// @Router /checkout [post]
func Checkout(w http.ResponseWriter, r *http.Request) {
	if os.Getenv("ENVIRONMENT") == "development" && !checkDevAuth(r) {
		http.Error(w, "Acesso não autorizado: forneça o UUID", http.StatusUnauthorized)
		return
	}
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
// @Param Authorization header string false "JWT token"
// @Param pedido body Pedido true "Dados do pedido"
// @Success 200 {object} map[string]string
// @Router /pedido [post]
func RegistrarPedido(w http.ResponseWriter, r *http.Request) {
	if os.Getenv("ENVIRONMENT") == "development" && !checkDevAuth(r) {
		http.Error(w, "Acesso não autorizado: forneça o UUID", http.StatusUnauthorized)
		return
	}
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

	// Verifica estoque e desconta ao registrar o pedido
	for _, item := range pedido.Itens {
		encontrado := false
		for i := range hats {
			if hats[i].ID == item.ID {
				encontrado = true
				if hats[i].Quantidade < item.Quantidade {
					http.Error(w, "Sem estoque suficiente para o chapéu: "+hats[i].Nome, http.StatusConflict)
					return
				}
				// Desconta o estoque
				hats[i].Quantidade -= item.Quantidade
			}
		}
		if !encontrado {
			http.Error(w, "Chapéu não encontrado: "+item.Nome, http.StatusBadRequest)
			return
		}
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
	if os.Getenv("ENVIRONMENT") == "development" && !checkDevAuth(r) {
		http.Error(w, "Acesso não autorizado: forneça o UUID", http.StatusUnauthorized)
		return
	}
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

// PATCH /api/hats/{id}/estoque
// EditarEstoqueHat godoc
// @Summary Edita o estoque de um chapéu
// @Description Atualiza a quantidade de estoque de um chapéu pelo ID
// @Tags hats
// @Accept json
// @Produce json
// @Param id path int true "ID do chapéu"
// @Param body body handlers.EstoquePayload true "Nova quantidade de estoque (0 a 200)"
// @Success 200 {object} Hat
// @Failure 400 {string} string "Quantidade inválida ou erro de requisição"
// @Failure 404 {string} string "Chapéu não encontrado"
// @Router /hats/{id}/estoque [patch]
func EditarEstoqueHat(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}
	var p EstoquePayload
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}
	if p.Quantidade < 0 || p.Quantidade > 200 {
		http.Error(w, "Quantidade deve ser entre 0 e 200", http.StatusBadRequest)
		return
	}
	for i := range hats {
		if hats[i].ID == id {
			hats[i].Quantidade = p.Quantidade
			json.NewEncoder(w).Encode(hats[i])
			return
		}
	}
	http.Error(w, "Chapéu não encontrado", http.StatusNotFound)
}

type EstoquePayload struct {
	Quantidade int `json:"quantidade"`
}
