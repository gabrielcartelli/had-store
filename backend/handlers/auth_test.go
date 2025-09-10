package handlers

import (
	"bytes"
	"encoding/json"
	"hat-store-training/backend/models"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRegisterHandler_Success(t *testing.T) {
	// Limpa o mapa de usuários antes do teste
	users = make(map[string]models.User)
	userIDCounter = 1

	body := models.RegisterRequest{
		Email:    "teste@exemplo.com",
		Password: "senha123",
	}
	jsonBody, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	RegisterHandler(w, req)
	resp := w.Result()

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("esperado status 201, recebeu %d", resp.StatusCode)
	}
}

func TestRegisterHandler_UsuarioExistente(t *testing.T) {
	users = make(map[string]models.User)
	userIDCounter = 1
	users["teste@exemplo.com"] = models.User{ID: 1, Email: "teste@exemplo.com", Password: "hash"}

	body := models.RegisterRequest{
		Email:    "teste@exemplo.com",
		Password: "senha123",
	}
	jsonBody, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	RegisterHandler(w, req)
	resp := w.Result()

	if resp.StatusCode != http.StatusConflict {
		t.Errorf("esperado status 409, recebeu %d", resp.StatusCode)
	}
}

func TestRegisterHandler_BodyInvalido(t *testing.T) {
	users = make(map[string]models.User)
	userIDCounter = 1

	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer([]byte("{invalido}")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	RegisterHandler(w, req)
	resp := w.Result()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("esperado status 400, recebeu %d", resp.StatusCode)
	}
}

func TestLoginHandler_Sucesso(t *testing.T) {
	users = make(map[string]models.User)
	userIDCounter = 1
	// Cria usuário válido
	RegisterHandler(httptest.NewRecorder(), httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBufferString(`{"email":"login@exemplo.com","password":"senha123"}`)))

	body := models.LoginRequest{
		Email:    "login@exemplo.com",
		Password: "senha123",
	}
	jsonBody, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	LoginHandler(w, req)
	resp := w.Result()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("esperado status 200, recebeu %d", resp.StatusCode)
	}
}

func TestLoginHandler_EmailOuSenhaInvalido(t *testing.T) {
	users = make(map[string]models.User)
	userIDCounter = 1
	RegisterHandler(httptest.NewRecorder(), httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBufferString(`{"email":"login@exemplo.com","password":"senha123"}`)))

	body := models.LoginRequest{
		Email:    "login@exemplo.com",
		Password: "errada",
	}
	jsonBody, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	LoginHandler(w, req)
	resp := w.Result()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("esperado status 401, recebeu %d", resp.StatusCode)
	}
}

func TestLoginHandler_BodyInvalido(t *testing.T) {
	users = make(map[string]models.User)
	userIDCounter = 1

	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer([]byte("{invalido}")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	LoginHandler(w, req)
	resp := w.Result()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("esperado status 400, recebeu %d", resp.StatusCode)
	}
}
