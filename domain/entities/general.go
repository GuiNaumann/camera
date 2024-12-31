package entities

type StatusCode int64

const (
	StatusExist StatusCode = 0

	StatusIncomplete StatusCode = 3

	StatusDeleted StatusCode = 2

	StatusNeedsRelation = 4
)

type Error struct {
	Code    string `json:"code"`
	Message string `json:"Message"`
}

type Cookie struct {
	Name       string
	Value      string
	RawExpires string

	MaxAge   int
	Secure   bool
	HttpOnly bool
	Raw      string
	Unparsed []string
}

type Type struct {
	Code  int    `json:"id"`
	Label string `json:"label"`
}

type PasswordEntropy struct {
	Entropy float64 `json:"entropy"`

	MinEntropy int `json:"minEntropy"`
}
