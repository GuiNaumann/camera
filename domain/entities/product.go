package entities

import "time"

// Product - Entity for Product
type Product struct {

	// Id - Identificador do Produto
	Id int64 `json:"id"`

	// UserID - Identificador do Usuário relacionado ao Produto
	UserID int64 `json:"userID"`

	// Name - Nome do Produto
	Name string `json:"name"`

	// Description - Descrição do Produto
	Description string `json:"description"`

	// Quantidade - Quantidade de Produtos
	Quantidade int `json:"quantity"`

	// Preco - Preço do Produto
	Preco float64 `json:"price"`

	// Tamanho - Tamanho do Produto (P, M, G, etc.)
	Tamanho string `json:"size"`

	// ImageBase64 - Imagem em formato base64
	ImageBase64 string `json:"imageBase64,omitempty"`

	// ImageURL - URL da Imagem do Produto
	ImageURL string `json:"imageURL,omitempty"`

	// IsActive - Indica se o Produto está ativo
	IsActive bool `json:"isActive"`

	// StatusCode - Status do Produto
	StatusCode int `json:"statusCode"`

	// CreatedAt - Data de Criação do Produto
	CreatedAt time.Time `json:"createdAt"`

	// ModifiedAt - Data da Última Modificação do Produto
	ModifiedAt time.Time `json:"modifiedAt"`
}
