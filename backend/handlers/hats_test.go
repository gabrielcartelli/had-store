package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetHats(t *testing.T) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/hats", nil)
	GetHats(w, req)
	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("esperado status 200, recebeu %d", resp.StatusCode)
	}
	type HatComEstoque struct {
		ID         int     `json:"id"`
		Nome       string  `json:"nome"`
		Price      float64 `json:"price"`
		Quantidade int     `json:"quantidade"`
		TemEstoque bool    `json:"temEstoque"`
		Categoria  string  `json:"categoria"`
	}
	var result []HatComEstoque
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Errorf("erro ao decodificar resposta: %v", err)
	}
	if len(result) == 0 {
		t.Errorf("esperado lista de chapéus, recebeu vazia")
	}
	for _, hat := range result {
		if hat.Quantidade > 0 && !hat.TemEstoque {
			t.Errorf("chapéu %s deveria ter estoque, mas temEstoque=false", hat.Nome)
		}
		if hat.Quantidade == 0 && hat.TemEstoque {
			t.Errorf("chapéu %s não deveria ter estoque, mas temEstoque=true", hat.Nome)
		}
		if hat.Categoria == "" {
			t.Errorf("chapéu %s deveria ter categoria definida", hat.Nome)
		}
	}
	// Testa filtro por categoria nacional
	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/hats?categoria=nacional", nil)
	GetHats(w, req)
	resp = w.Result()
	var resultNacional []HatComEstoque
	if err := json.NewDecoder(resp.Body).Decode(&resultNacional); err != nil {
		t.Errorf("erro ao decodificar resposta: %v", err)
	}
	for _, hat := range resultNacional {
		if hat.Categoria != "nacional" {
			t.Errorf("esperado apenas chapéus nacionais, recebeu %s", hat.Categoria)
		}
	}
	// Testa filtro múltiplo
	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/hats?categoria=nacional,importado", nil)
	GetHats(w, req)
	resp = w.Result()
	var resultMulti []HatComEstoque
	if err := json.NewDecoder(resp.Body).Decode(&resultMulti); err != nil {
		t.Errorf("erro ao decodificar resposta: %v", err)
	}
	for _, hat := range resultMulti {
		if hat.Categoria != "nacional" && hat.Categoria != "importado" {
			t.Errorf("esperado apenas chapéus nacionais ou importados, recebeu %s", hat.Categoria)
		}
	}
	// Testa filtro inexistente
	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/hats?categoria=inexistente", nil)
	GetHats(w, req)
	resp = w.Result()
	var resultInexistente []HatComEstoque
	if err := json.NewDecoder(resp.Body).Decode(&resultInexistente); err != nil {
		t.Errorf("erro ao decodificar resposta: %v", err)
	}
	if len(resultInexistente) != 0 {
		t.Errorf("esperado lista vazia para categoria inexistente, recebeu %d chapéus", len(resultInexistente))
	}
	// Testa filtro por valor mínimo
	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/hats?min=100", nil)
	GetHats(w, req)
	resp = w.Result()
	var resultMin []HatComEstoque
	if err := json.NewDecoder(resp.Body).Decode(&resultMin); err != nil {
		t.Errorf("erro ao decodificar resposta: %v", err)
	}
	for _, hat := range resultMin {
		if hat.Price < 100 {
			t.Errorf("esperado apenas chapéus com preço >= 100, recebeu %.2f", hat.Price)
		}
	}
	// Testa filtro por valor máximo
	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/hats?max=50", nil)
	GetHats(w, req)
	resp = w.Result()
	var resultMax []HatComEstoque
	if err := json.NewDecoder(resp.Body).Decode(&resultMax); err != nil {
		t.Errorf("erro ao decodificar resposta: %v", err)
	}
	for _, hat := range resultMax {
		if hat.Price > 50 {
			t.Errorf("esperado apenas chapéus com preço <= 50, recebeu %.2f", hat.Price)
		}
	}
	// Testa filtro por valor mínimo e máximo
	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/hats?min=50&max=100", nil)
	GetHats(w, req)
	resp = w.Result()
	var resultFaixa []HatComEstoque
	if err := json.NewDecoder(resp.Body).Decode(&resultFaixa); err != nil {
		t.Errorf("erro ao decodificar resposta: %v", err)
	}
	for _, hat := range resultFaixa {
		if hat.Price < 50 || hat.Price > 100 {
			t.Errorf("esperado apenas chapéus com preço entre 50 e 100, recebeu %.2f", hat.Price)
		}
	}
}
func TestRegistrarPedido_SemEstoque(t *testing.T) {
	// Zera o estoque do chapéu 1
	for i := range hats {
		if hats[i].ID == 1 {
			hats[i].Quantidade = 0
		}
	}
	pedido := Pedido{
		Nome:      "Cliente Teste",
		CPF:       "12345678900",
		Email:     "cliente@teste.com",
		Pagamento: "pix",
		Itens:     []HatPedido{{ID: 1, Nome: "Chapéu Panamá", Price: 120.00, Quantidade: 1}},
		Cupom:     "",
	}
	jsonBody, _ := json.Marshal(pedido)
	req := httptest.NewRequest(http.MethodPost, "/pedido", bytes.NewBuffer(jsonBody))
	w := httptest.NewRecorder()
	RegistrarPedido(w, req)
	resp := w.Result()
	if resp.StatusCode != http.StatusConflict {
		t.Errorf("esperado status 409 (sem estoque), recebeu %d", resp.StatusCode)
	}
	// Restaura estoque para outros testes
	for i := range hats {
		if hats[i].ID == 1 {
			hats[i].Quantidade = 10
		}
	}
}

func TestRegistrarPedido_Sucesso(t *testing.T) {
	pedidos = nil
	pedido := Pedido{
		Nome:      "Cliente Teste",
		CPF:       "12345678900",
		Email:     "cliente@teste.com",
		Pagamento: "pix",
		Itens:     []HatPedido{{ID: 1, Nome: "Chapéu Panamá", Price: 120.00, Quantidade: 1}},
		Cupom:     "",
	}
	jsonBody, _ := json.Marshal(pedido)
	req := httptest.NewRequest(http.MethodPost, "/pedido", bytes.NewBuffer(jsonBody))
	w := httptest.NewRecorder()
	RegistrarPedido(w, req)
	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("esperado status 200, recebeu %d", resp.StatusCode)
	}
}

func TestRegistrarPedido_BodyInvalido(t *testing.T) {
	pedidos = nil
	req := httptest.NewRequest(http.MethodPost, "/pedido", bytes.NewBuffer([]byte("{invalido}")))
	w := httptest.NewRecorder()
	RegistrarPedido(w, req)
	resp := w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("esperado status 400, recebeu %d", resp.StatusCode)
	}
}

func TestRegistrarPedido_DescontoPrimeiroPedido(t *testing.T) {
	pedidos = nil
	pedido := Pedido{
		Nome:      "Cliente Teste",
		CPF:       "12345678900",
		Email:     "cliente@teste.com",
		Pagamento: "pix",
		Itens:     []HatPedido{{ID: 1, Nome: "Chapéu Panamá", Price: 100.00, Quantidade: 1}},
		Cupom:     "HATOFF",
	}
	jsonBody, _ := json.Marshal(pedido)
	req := httptest.NewRequest(http.MethodPost, "/pedido", bytes.NewBuffer(jsonBody))
	w := httptest.NewRecorder()
	RegistrarPedido(w, req)
	resp := w.Result()
	var respBody map[string]string
	json.NewDecoder(resp.Body).Decode(&respBody)
	if resp.StatusCode != http.StatusOK {
		t.Errorf("esperado status 200, recebeu %d", resp.StatusCode)
	}
	if len(pedidos) == 0 || pedidos[0].Total != 80.0 {
		t.Errorf("esperado desconto de 20%%, recebeu %.2f", pedidos[0].Total)
	}
}

func TestRegistrarPedido_SegundoPedidoSemDesconto(t *testing.T) {
	pedidos = nil
	pedido1 := Pedido{
		Nome:      "Cliente Teste",
		CPF:       "12345678900",
		Email:     "cliente@teste.com",
		Pagamento: "pix",
		Itens:     []HatPedido{{ID: 1, Nome: "Chapéu Panamá", Price: 100.00, Quantidade: 1}},
		Cupom:     "HATOFF",
	}
	jsonBody1, _ := json.Marshal(pedido1)
	req1 := httptest.NewRequest(http.MethodPost, "/pedido", bytes.NewBuffer(jsonBody1))
	w1 := httptest.NewRecorder()
	RegistrarPedido(w1, req1)

	pedido2 := Pedido{
		Nome:      "Cliente Teste",
		CPF:       "12345678900",
		Email:     "cliente@teste.com",
		Pagamento: "pix",
		Itens:     []HatPedido{{ID: 1, Nome: "Chapéu Panamá", Price: 100.00, Quantidade: 1}},
		Cupom:     "HATOFF",
	}
	jsonBody2, _ := json.Marshal(pedido2)
	req2 := httptest.NewRequest(http.MethodPost, "/pedido", bytes.NewBuffer(jsonBody2))
	w2 := httptest.NewRecorder()
	RegistrarPedido(w2, req2)

	if len(pedidos) < 2 {
		t.Fatalf("esperado 2 pedidos, recebeu %d", len(pedidos))
	}
	if pedidos[1].Total != 100.0 {
		t.Errorf("esperado total sem desconto no segundo pedido, recebeu %.2f", pedidos[1].Total)
	}
}

func TestConsultarPedidos_Sucesso(t *testing.T) {
	pedidos = nil
	pedido := Pedido{
		Nome:      "Cliente Teste",
		CPF:       "12345678900",
		Email:     "cliente@teste.com",
		Pagamento: "pix",
		Itens:     []HatPedido{{ID: 1, Nome: "Chapéu Panamá", Price: 120.00, Quantidade: 1}},
		Cupom:     "",
		Total:     120.00,
	}
	pedidos = append(pedidos, pedido)
	req := httptest.NewRequest(http.MethodGet, "/pedidos?cpf=12345678900", nil)
	w := httptest.NewRecorder()
	ConsultarPedidos(w, req)
	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("esperado status 200, recebeu %d", resp.StatusCode)
	}
	var result []Pedido
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Errorf("erro ao decodificar resposta: %v", err)
	}
	if len(result) != 1 {
		t.Errorf("esperado 1 pedido, recebeu %d", len(result))
	}
}

func TestConsultarPedidos_CPFNaoInformado(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/pedidos", nil)
	w := httptest.NewRecorder()
	ConsultarPedidos(w, req)
	resp := w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("esperado status 400, recebeu %d", resp.StatusCode)
	}
}
