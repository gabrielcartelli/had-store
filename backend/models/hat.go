package models

type Hat struct {
	ID         int     `json:"id"`
	Nome       string  `json:"nome"`
	Price      float64 `json:"price"`
	Quantidade int     `json:"quantidade"`
	Categoria  string  `json:"categoria"` // nacional, importado, crescer
}
