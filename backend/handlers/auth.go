// backend/handlers/auth.go
package handlers

import (
	"hat-store-training/backend/models"
	"sync"
	"time"
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"github.com/golang-jwt/jwt/v5"
)


var users = make(map[string]models.User)
var usersMutex sync.Mutex
var userIDCounter = 1


var jwtKey = []byte("minha_chave_super_secreta")

// Estrutura para rastrear tentativas de login falhas
type loginAttempt struct {
	Count      int
	LastAttempt time.Time
}

var failedAttempts = make(map[string]loginAttempt)
var attemptsMutex sync.Mutex

// RegisterHandler cuida do registro de novos usuários
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var creds struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// 1. Pega o email e a senha que o usuário enviou
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Requisição inválida", http.StatusBadRequest)
		return
	}

	usersMutex.Lock()
	defer usersMutex.Unlock()

	// 2. Verifica se o membro já existe na nossa lista
	if _, exists := users[creds.Email]; exists {
		http.Error(w, "Usuário já existe", http.StatusConflict)
		return
	}

	// 3. O PASSO MAIS IMPORTANTE: Criando o  (Hash)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Erro ao processar a senha", http.StatusInternalServerError)
		return
	}

	// 4. Adiciona o novo membro à nossa lista
	newUser := models.User{
		ID:       userIDCounter,
		Email:    creds.Email,
		Password: string(hashedPassword),
	}
	users[creds.Email] = newUser
	userIDCounter++

	log.Printf("AUDITORIA: Novo usuário registrado com sucesso: %s", creds.Email)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Usuário criado com sucesso!"})
}

// LoginHandler cuida da autenticação dos usuários
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var creds struct {
		Email      string `json:"email"`
		Password   string `json:"password"`
		RememberMe bool   `json:"rememberMe"` // Para a opção "Lembrar-me"
	}

	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Requisição inválida", http.StatusBadRequest)
		return
	}

	// --- CRITÉRIO 4: Bloqueio por 5 tentativas ---
	attemptsMutex.Lock()
	attempt, ok := failedAttempts[creds.Email]
	// Se já errou 5 vezes e a última tentativa foi nos últimos 15 minutos, bloqueia.
	if ok && attempt.Count >= 5 && time.Since(attempt.LastAttempt) < 15*time.Minute {
		attemptsMutex.Unlock()
		log.Printf("AUDITORIA: Tentativa de login bloqueada para o usuário: %s", creds.Email)
		http.Error(w, "Muitas tentativas de login. Tente novamente mais tarde.", http.StatusTooManyRequests)
		return
	}
	attemptsMutex.Unlock()


	// --- CRITÉRIO 1: Validar credenciais ---
	usersMutex.Lock()
	user, exists := users[creds.Email]
	usersMutex.Unlock()

	// Se o usuário não existe OU a senha está errada...
	if !exists || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password)) != nil {
		// --- CRITÉRIO 6: Log de auditoria ---
		log.Printf("AUDITORIA: Tentativa de login falhou para o usuário: %s", creds.Email)
		
		// Anota o erro 
		attemptsMutex.Lock()
		attempt.Count++
		attempt.LastAttempt = time.Now()
		failedAttempts[creds.Email] = attempt
		attemptsMutex.Unlock()

		http.Error(w, "Email ou senha inválidos", http.StatusUnauthorized)
		return
	}
	
	// Limpa as tentativas para este usuário
	attemptsMutex.Lock()
	delete(failedAttempts, creds.Email)
	attemptsMutex.Unlock()

	// --- CRITÉRIO 5: "Lembrar-me" ---
	// Define a validade do token. Se "Lembrar-me" estiver marcado, dura 30 dias. Senão, 8 horas.
	expirationTime := time.Now().Add(8 * time.Hour)
	if creds.RememberMe {
		expirationTime = time.Now().Add(30 * 24 * time.Hour) // 30 dias
	}

	// Cria o  (Token JWT)
	claims := &jwt.RegisteredClaims{
		Subject:   user.Email,
		ExpiresAt: jwt.NewNumericDate(expirationTime),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		http.Error(w, "Erro ao gerar token", http.StatusInternalServerError)
		return
	}

	// --- CRITÉRIOS 2 e 6: Acessa o site e loga sucesso ---
	log.Printf("AUDITORIA: Login bem-sucedido para o usuário: %s", creds.Email)

	// Envia o token de volta para o frontend
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"token": tokenString,
	})
}