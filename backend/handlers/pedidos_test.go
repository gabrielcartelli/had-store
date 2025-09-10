package handlers

import (
	"bytes"
	"encoding/json"
	"hat-store-training/backend/models"
	"net/http"
	"net/http/httptest"
	"testing"
)

func pedidoValido() models.Pedido {
	return models.Pedido{
		Nome:      "Cliente Teste",
		CPF:       "123.456.789-00",
		Email:     "cliente@teste.com",
		Telefone:  "(11) 91234-5678",
		CEP:       "12345-678",
		UF:        "SP",
		Pagamento: "pix",
		Itens:     []models.Hat{{ID: 1, Nome: "Chapéu Panamá", Price: 100.00, Quantidade: 1}},
		Cupom:     "",
	}
}

func TestCriarPedido_Sucesso(t *testing.T) {
	pedidosFeitos = make(map[string]bool)
	hatOffUsadoPorCPF = make(map[string]bool)
	pedido := pedidoValido()
	jsonBody, _ := json.Marshal(pedido)
	req := httptest.NewRequest(http.MethodPost, "/pedido", bytes.NewBuffer(jsonBody))
	w := httptest.NewRecorder()
	CriarPedido(w, req)
	resp := w.Result()
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("esperado status 201, recebeu %d", resp.StatusCode)
	}
}

func TestCriarPedido_CupomHATOFFPrimeiraVez(t *testing.T) {
	pedidosFeitos = make(map[string]bool)
	hatOffUsadoPorCPF = make(map[string]bool)
	pedido := pedidoValido()
	pedido.Cupom = "HATOFF"
	jsonBody, _ := json.Marshal(pedido)
	req := httptest.NewRequest(http.MethodPost, "/pedido", bytes.NewBuffer(jsonBody))
	w := httptest.NewRecorder()
	CriarPedido(w, req)
	resp := w.Result()
	var respPedido models.Pedido
	json.NewDecoder(resp.Body).Decode(&respPedido)
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("esperado status 201, recebeu %d", resp.StatusCode)
	}
	if respPedido.Total != 80.0 {
		t.Errorf("esperado desconto de 20%%, recebeu %.2f", respPedido.Total)
	}
}

func TestCriarPedido_CupomHATOFFSegundaVez(t *testing.T) {
	pedidosFeitos = make(map[string]bool)
	hatOffUsadoPorCPF = make(map[string]bool)
	pedido := pedidoValido()
	pedido.Cupom = "HATOFF"
	jsonBody, _ := json.Marshal(pedido)
	req := httptest.NewRequest(http.MethodPost, "/pedido", bytes.NewBuffer(jsonBody))
	w := httptest.NewRecorder()
	CriarPedido(w, req)
	// Segunda tentativa com mesmo CPF e cupom
	jsonBody2, _ := json.Marshal(pedido)
	req2 := httptest.NewRequest(http.MethodPost, "/pedido", bytes.NewBuffer(jsonBody2))
	w2 := httptest.NewRecorder()
	CriarPedido(w2, req2)
	resp2 := w2.Result()
	if resp2.StatusCode != http.StatusForbidden {
		t.Errorf("esperado status 403, recebeu %d", resp2.StatusCode)
	}
}

func TestCriarPedido_NomeInvalido(t *testing.T) {
	pedido := pedidoValido()
	pedido.Nome = "A"
	jsonBody, _ := json.Marshal(pedido)
	req := httptest.NewRequest(http.MethodPost, "/pedido", bytes.NewBuffer(jsonBody))
	w := httptest.NewRecorder()
	CriarPedido(w, req)
	resp := w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("esperado status 400, recebeu %d", resp.StatusCode)
	}
}

func TestCriarPedido_CPFInvalido(t *testing.T) {
	pedido := pedidoValido()
	pedido.CPF = "12345678900"
	jsonBody, _ := json.Marshal(pedido)
	req := httptest.NewRequest(http.MethodPost, "/pedido", bytes.NewBuffer(jsonBody))
	w := httptest.NewRecorder()
	CriarPedido(w, req)
	resp := w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("esperado status 400, recebeu %d", resp.StatusCode)
	}
}

func TestCriarPedido_EmailInvalido(t *testing.T) {
	pedido := pedidoValido()
	pedido.Email = "emailinvalido"
	jsonBody, _ := json.Marshal(pedido)
	req := httptest.NewRequest(http.MethodPost, "/pedido", bytes.NewBuffer(jsonBody))
	w := httptest.NewRecorder()
	CriarPedido(w, req)
	resp := w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("esperado status 400, recebeu %d", resp.StatusCode)
	}
}

func TestCriarPedido_TelefoneInvalido(t *testing.T) {
	pedido := pedidoValido()
	pedido.Telefone = "11912345678"
	jsonBody, _ := json.Marshal(pedido)
	req := httptest.NewRequest(http.MethodPost, "/pedido", bytes.NewBuffer(jsonBody))
	w := httptest.NewRecorder()
	CriarPedido(w, req)
	resp := w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("esperado status 400, recebeu %d", resp.StatusCode)
	}
}

func TestCriarPedido_CEPInvalido(t *testing.T) {
	pedido := pedidoValido()
	pedido.CEP = "12345678"
	jsonBody, _ := json.Marshal(pedido)
	req := httptest.NewRequest(http.MethodPost, "/pedido", bytes.NewBuffer(jsonBody))
	w := httptest.NewRecorder()
	CriarPedido(w, req)
	resp := w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("esperado status 400, recebeu %d", resp.StatusCode)
	}
}

func TestCriarPedido_UFInvalido(t *testing.T) {
	pedido := pedidoValido()
	pedido.UF = "SaoPaulo"
	jsonBody, _ := json.Marshal(pedido)
	req := httptest.NewRequest(http.MethodPost, "/pedido", bytes.NewBuffer(jsonBody))
	w := httptest.NewRecorder()
	CriarPedido(w, req)
	resp := w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("esperado status 400, recebeu %d", resp.StatusCode)
	}
}

func TestCriarPedido_PagamentoInvalido(t *testing.T) {
	pedido := pedidoValido()
	pedido.Pagamento = "cartao"
	jsonBody, _ := json.Marshal(pedido)
	req := httptest.NewRequest(http.MethodPost, "/pedido", bytes.NewBuffer(jsonBody))
	w := httptest.NewRecorder()
	CriarPedido(w, req)
	resp := w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("esperado status 400, recebeu %d", resp.StatusCode)
	}
}
