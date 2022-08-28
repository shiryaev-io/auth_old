package dto

// Сожержит access и refresh токены
type Tokens struct {
	Access  string `json:"asscess_token"`
	Refresh string `json:"refresh_token"`
}
