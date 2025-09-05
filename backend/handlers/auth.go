package handlers

import (
	"encoding/json"
	"hat-store-training/backend/models"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// Nossa "lista de membros". A chave do mapa é o email do usuário.
var users = make(map[string]models.User)
var usersMutex sync.Mutex
var userIDCounter = 1

// A chave secreta para criar e verificar nossos "carimbos" (tokens).
// Em um projeto real, isso viria de uma variável de ambiente!
var jwtKey = []byte("e3e6c6c2-9b7d-4c5e-8c1a-2f7b8f8e2a1d")

// Estrutura para rastrear tentativas de login falhas
type loginAttempt struct {
	Count       int
	LastAttempt time.Time
}

// Nosso "caderninho" para anotar quem errou o aperto de mão
var failedAttempts = make(map[string]loginAttempt)
var attemptsMutex sync.Mutex

// RegisterHandler godoc
// @Summary Registra um novo usuário
// @Description Cria uma nova conta de usuário com email e senha
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body models.RegisterRequest true "Credenciais de Registro"
// @Success 201 {object} map[string]string "Usuário criado com sucesso!"
// @Failure 400 {string} string "Requisição inválida"
// @Failure 409 {string} string "Usuário já existe"
// @Router /auth/register [post]
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var creds models.RegisterRequest

	// 1. Pega o email e a senha que o usuário enviou
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		log.Printf("[ERROR] Registro inválido: %v", err)
		http.Error(w, "Requisição inválida", http.StatusBadRequest)
		return
	}

	usersMutex.Lock()
	defer usersMutex.Unlock()

	// 2. Verifica se o membro já existe na nossa lista
	if _, exists := users[creds.Email]; exists {
		log.Printf("[WARN] Tentativa de registro para email já existente: %s", creds.Email)
		http.Error(w, "Usuário já existe", http.StatusConflict)
		return
	}

	// 3. O PASSO MAIS IMPORTANTE: Criando o "aperto de mão secreto codificado" (Hash)
	// Nós NUNCA guardamos a senha real. Guardamos uma versão embaralhada dela.
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("[ERROR] Erro ao processar senha para %s: %v", creds.Email, err)
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

	ip := strings.Split(r.RemoteAddr, ":")[0]
	log.Printf("[INFO][%s][%s] Novo usuário registrado", ip, creds.Email)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Usuário criado com sucesso!"})
}

// LoginHandler godoc
// @Summary Autentica um usuário
// @Description Loga um usuário com email e senha e retorna um token JWT
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body models.LoginRequest true "Credenciais de Login"
// @Success 200 {object} map[string]string "Token JWT"
// @Failure 400 {string} string "Requisição inválida"
// @Failure 401 {string} string "Email ou senha inválidos"
// @Failure 429 {string} string "Muitas tentativas de login"
// @Router /auth/login [post]
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var creds models.LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		log.Printf("[ERROR] Login inválido: %v", err)
		http.Error(w, "Requisição inválida", http.StatusBadRequest)
		return
	}

	// --- CRITÉRIO 4: Bloqueio por 5 tentativas ---
	attemptsMutex.Lock()
	attempt, ok := failedAttempts[creds.Email]
	// Se já errou 5 vezes e a última tentativa foi nos últimos 15 minutos, bloqueia.
	if ok && attempt.Count >= 5 && time.Since(attempt.LastAttempt) < 15*time.Minute {
		attemptsMutex.Unlock()
		log.Printf("[WARN] Login bloqueado por excesso de tentativas: %s", creds.Email)
		http.Error(w, "Muitas tentativas de login. Tente novamente mais tarde.", http.StatusTooManyRequests)
		return
	}
	attemptsMutex.Unlock()

	// --- CRITÉRIO 1: Validar credenciais ---
	usersMutex.Lock()
	user, exists := users[creds.Email]
	usersMutex.Unlock()

	// Se o usuário não existe OU o "aperto de mão" está errado...
	if !exists || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password)) != nil {
		log.Printf("[WARN] Login falhou para usuário: %s", creds.Email)

		attemptsMutex.Lock()
		if time.Since(attempt.LastAttempt) >= 15*time.Minute {
			attempt.Count = 0
		}
		attempt.Count++
		attempt.LastAttempt = time.Now()
		failedAttempts[creds.Email] = attempt
		attemptsMutex.Unlock()

		http.Error(w, "Email ou senha inválidos", http.StatusUnauthorized)
		return
	}

	// Se chegou aqui, o aperto de mão está CORRETO!

	// Limpa o caderninho de erros para este usuário
	attemptsMutex.Lock()
	delete(failedAttempts, creds.Email)
	attemptsMutex.Unlock()

	// --- CRITÉRIO 5: "Lembrar-me" ---
	// Define a validade do carimbo. Se "Lembrar-me" estiver marcado, dura 30 dias. Senão, 8 horas.
	expirationTime := time.Now().Add(8 * time.Hour)
	if creds.RememberMe {
		expirationTime = time.Now().Add(30 * 24 * time.Hour) // 30 dias
	}

	// Cria o "carimbo" (Token JWT)
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
	ip := strings.Split(r.RemoteAddr, ":")[0]
	log.Printf("[INFO][%s][%s] Login bem-sucedido", ip, creds.Email)

	// Envia o carimbo de volta para o frontend
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"token": tokenString,
	})
}
