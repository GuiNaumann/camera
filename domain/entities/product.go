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

	// IsActive - Indica se o Produto está ativo
	IsActive bool `json:"isActive"`

	// StatusCode - Status do Produto
	StatusCode int `json:"statusCode"`

	// CreatedAt - Data de Criação do Produto
	CreatedAt time.Time `json:"createdAt"`

	// ModifiedAt - Data da Última Modificação do Produto
	ModifiedAt time.Time `json:"modifiedAt"`

	//todo: add ---------------------------------------------

	IPAddress string `json:"ip_address"`

	Port int `json:"port"`

	Username string `json:"username"`

	Password string `json:"password"`

	StreamPath string `json:"stream_path"`

	CameraType string `json:"camera_type"` // Intelbras ou Yoosee

	LocalID int64 `json:"localID"`
}

type Local struct {

	// Id - Identificador do Local
	Id int64 `json:"id"`

	// UserID - Identificador do Usuário relacionado ao Local
	UserID int64 `json:"userID"`

	// Name - Nome do Local
	Name string `json:"name"`

	// Description - Descrição do Local
	Description string `json:"description,omitempty"`

	// State - Estado onde o Local está
	State string `json:"state"`

	// City - Cidade onde o Local está
	City string `json:"city"`

	// Street - Rua onde o Local está
	Street string `json:"street"`

	// Number - Número do Local
	Number string `json:"number"`

	IsActive bool `json:"isActive"`

	StatusCode int `json:"statusCode"`

	// ZipCode - CEP do Local
	ZipCode string `json:"zipCode,omitempty"`

	// Complement - Complemento do Endereço
	Complement string `json:"complement,omitempty"`

	// CreatedAt - Data de Criação do Local
	CreatedAt time.Time `json:"createdAt"`

	// ModifiedAt - Data da Última Modificação do Local
	ModifiedAt time.Time `json:"modifiedAt"`
}

type CameraStream struct {
	CameraID    int64  `json:"camera_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	StreamURL   string `json:"stream_url"`
}
