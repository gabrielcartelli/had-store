package handlers

import (
	"encoding/json"
	"hat-store-training/backend/models"
	"net/http"
	"regexp"
	"sync"
)

// Simulação de armazenamento em memória dos CPFs que já fizeram pedido
var pedidosFeitos = make(map[string]bool)
var pedidosMutex sync.Mutex

// Armazena CPFs que já usaram o cupom HAD10
var had10UsadoPorCPF = make(map[string]bool)
var had10Mutex sync.Mutex

func PedidoJaExiste(cpf string) bool {
	pedidosMutex.Lock()
	defer pedidosMutex.Unlock()
	return pedidosFeitos[cpf]
}

func CPFUsouHad10(cpf string) bool {
	had10Mutex.Lock()
	defer had10Mutex.Unlock()
	return had10UsadoPorCPF[cpf]
}

func CriarPedido(w http.ResponseWriter, r *http.Request) {
	var pedido models.Pedido

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&pedido)
	if err != nil {
		http.Error(w, "Dados inválidos", http.StatusBadRequest)
		return
	}

	// Só retorna erro se o cupom for HAD10 e o CPF já usou HAD10 antes
	if pedido.Cupom == "HAD10" && CPFUsouHad10(pedido.CPF) {
		http.Error(w, "Cupom HAD10 já utilizado por este CPF.", http.StatusForbidden)
		return
	}

	// Calcula total
	total := 0.0
	for _, item := range pedido.Itens {
		total += item.Price * float64(item.Quantidade)
	}

	// Aplica desconto HAD10 só se for a primeira vez desse CPF
	if pedido.Cupom == "HAD10" {
		total = total * 0.9
	}

	pedido.Total = total

	// Validações (mantidas do seu código)
	if len(pedido.Nome) < 3 {
		http.Error(w, "Nome inválido", http.StatusBadRequest)
		return
	}
	cpfRegex := regexp.MustCompile(`^\d{3}\.\d{3}\.\d{3}-\d{2}$`)
	if !cpfRegex.MatchString(pedido.CPF) {
		http.Error(w, "CPF inválido", http.StatusBadRequest)
		return
	}
	emailRegex := regexp.MustCompile(`^[^@\s]+@[^@\s]+\.[^@\s]+$`)
	if !emailRegex.MatchString(pedido.Email) {
		http.Error(w, "Email inválido", http.StatusBadRequest)
		return
	}
	telRegex := regexp.MustCompile(`^\(\d{2}\) \d{5}-\d{4}$`)
	if !telRegex.MatchString(pedido.Telefone) {
		http.Error(w, "Telefone inválido", http.StatusBadRequest)
		return
	}
	cepRegex := regexp.MustCompile(`^\d{5}-\d{3}$`)
	if !cepRegex.MatchString(pedido.CEP) {
		http.Error(w, "CEP inválido", http.StatusBadRequest)
		return
	}
	ufRegex := regexp.MustCompile(`^[A-Za-z]{2}$`)
	if !ufRegex.MatchString(pedido.UF) {
		http.Error(w, "UF inválido", http.StatusBadRequest)
		return
	}
	if pedido.Pagamento != "pix" && pedido.Pagamento != "boleto" {
		http.Error(w, "Forma de pagamento inválida", http.StatusBadRequest)
		return
	}

	// Marca o CPF como tendo usado HAD10, apenas se o cupom foi usado
	if pedido.Cupom == "HAD10" {
		had10Mutex.Lock()
		had10UsadoPorCPF[pedido.CPF] = true
		had10Mutex.Unlock()
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(pedido)
}
