package models

type Pedido struct {
	ID        int     `json:"id"`
	Nome      string  `json:"nome"`
	CPF       string  `json:"cpf"`
	Email     string  `json:"email"`
	Telefone  string  `json:"telefone"`
	Endereco  string  `json:"endereco"`
	Numero    string  `json:"numero"`
	Bairro    string  `json:"bairro"`
	CEP       string  `json:"cep"`
	Cidade    string  `json:"cidade"`
	UF        string  `json:"uf"`
	Pagamento string  `json:"pagamento"` // "pix" ou "boleto"
	Itens     []Hat   `json:"itens"`
	Total     float64 `json:"total"`
	Cupom     string  `json:"cupom"`
}
