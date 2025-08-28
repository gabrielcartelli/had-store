// backend/models/api_models.go
package models

// RegisterRequest define a estrutura do corpo da requisição de registro.
type RegisterRequest struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}

// LoginRequest define a estrutura do corpo da requisição de login.
type LoginRequest struct {
    Email      string `json:"email"`
    Password   string `json:"password"`
    RememberMe bool   `json:"rememberMe"`
}