package handlers

import (
	"encoding/json"
	"hat-store-training/backend/models"
	"log"
	"net/http"
	"os"
	"regexp"
	"sync"
)

// Simulação de armazenamento em memória dos CPFs que já fizeram pedido
var pedidosFeitos = make(map[string]bool)
var pedidosMutex sync.Mutex

// Armazena CPFs que já usaram o cupom HATOFF
var hatOffUsadoPorCPF = make(map[string]bool)
var hatOffMutex sync.Mutex

func PedidoJaExiste(cpf string) bool {
	pedidosMutex.Lock()
	defer pedidosMutex.Unlock()
	return pedidosFeitos[cpf]
}

func CPFUsouHatOff(cpf string) bool {
	hatOffMutex.Lock()
	defer hatOffMutex.Unlock()
	return hatOffUsadoPorCPF[cpf]
}

func CriarPedido(w http.ResponseWriter, r *http.Request) {
	if os.Getenv("ENVIRONMENT") == "development" {
		uuid := r.Header.Get("X-Dev-UUID")
		expected := os.Getenv("DEV_UUID")
		if uuid == "" || uuid != expected {
			http.Error(w, "Acesso não autorizado: forneça o UUID", http.StatusUnauthorized)
			return
		}
	}
	var pedido models.Pedido
	var err error
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&pedido)
	if err != nil {
		log.Printf("[ERROR] CriarPedido: Dados inválidos | Erro: %v", err)
		http.Error(w, "Dados inválidos", http.StatusBadRequest)
		return
	}

	// Controle de estoque
	for _, item := range pedido.Itens {
		encontrado := false
		for i := range hats {
			if hats[i].ID == item.ID {
				encontrado = true
				if hats[i].Quantidade < item.Quantidade {
					http.Error(w, "Produto sem estoque suficiente: "+hats[i].Nome, http.StatusConflict)
					return
				}
			}
		}
		if !encontrado {
			http.Error(w, "Produto não encontrado: "+item.Nome, http.StatusNotFound)
			return
		}
	}

	// Só retorna erro se o cupom for HATOFF e o CPF já usou HATOFF antes
	if pedido.Cupom == "HATOFF" && CPFUsouHatOff(pedido.CPF) {
		log.Printf("[WARN] CriarPedido: Cupom HATOFF já utilizado | CPF: %s", pedido.CPF)
		http.Error(w, "Cupom HATOFF já utilizado por este CPF.", http.StatusForbidden)
		return
	}

	// Calcula total e desconta estoque
	total := 0.0
	for _, item := range pedido.Itens {
		total += item.Price * float64(item.Quantidade)
		for i := range hats {
			if hats[i].ID == item.ID {
				hats[i].Quantidade -= item.Quantidade
			}
		}
	}

	// Aplica desconto HATOFF só se for a primeira vez desse CPF
	if pedido.Cupom == "HATOFF" {
		total = total * 0.80
	}

	pedido.Total = total

	// Validações (mantidas do seu código)
	if len(pedido.Nome) < 3 {
		log.Printf("[WARN] CriarPedido: Nome inválido | Nome: %s", pedido.Nome)
		http.Error(w, "Nome inválido", http.StatusBadRequest)
		return
	}
	cpfRegex := regexp.MustCompile(`^\d{3}\.\d{3}\.\d{3}-\d{2}$`)
	if !cpfRegex.MatchString(pedido.CPF) {
		log.Printf("[WARN] CriarPedido: CPF inválido | CPF: %s", pedido.CPF)
		http.Error(w, "CPF inválido", http.StatusBadRequest)
		return
	}
	emailRegex := regexp.MustCompile(`^[^@\s]+@[^@\s]+\.[^@\s]+$`)
	if !emailRegex.MatchString(pedido.Email) {
		log.Printf("[WARN] CriarPedido: Email inválido | Email: %s", pedido.Email)
		http.Error(w, "Email inválido", http.StatusBadRequest)
		return
	}
	telRegex := regexp.MustCompile(`^\(\d{2}\) \d{5}-\d{4}$`)
	if !telRegex.MatchString(pedido.Telefone) {
		log.Printf("[WARN] CriarPedido: Telefone inválido | Telefone: %s", pedido.Telefone)
		http.Error(w, "Telefone inválido", http.StatusBadRequest)
		return
	}
	cepRegex := regexp.MustCompile(`^\d{5}-\d{3}$`)
	if !cepRegex.MatchString(pedido.CEP) {
		log.Printf("[WARN] CriarPedido: CEP inválido | CEP: %s", pedido.CEP)
		http.Error(w, "CEP inválido", http.StatusBadRequest)
		return
	}
	ufRegex := regexp.MustCompile(`^[A-Za-z]{2}$`)
	if !ufRegex.MatchString(pedido.UF) {
		log.Printf("[WARN] CriarPedido: UF inválido | UF: %s", pedido.UF)
		http.Error(w, "UF inválido", http.StatusBadRequest)
		return
	}
	if pedido.Pagamento != "pix" && pedido.Pagamento != "boleto" {
		log.Printf("[WARN] CriarPedido: Pagamento inválido | Pagamento: %s", pedido.Pagamento)
		http.Error(w, "Forma de pagamento inválida", http.StatusBadRequest)
		return
	}

	// Marca o CPF como tendo usado HATOFF, apenas se o cupom foi usado
	if pedido.Cupom == "HATOFF" {
		hatOffMutex.Lock()
		hatOffUsadoPorCPF[pedido.CPF] = true
		hatOffMutex.Unlock()
	}

	log.Printf("[INFO] Pedido criado | Nome: %s | CPF: %s | Total: %.2f | Pagamento: %s | Itens: %d", pedido.Nome, pedido.CPF, pedido.Total, pedido.Pagamento, len(pedido.Itens))
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(pedido)
}
