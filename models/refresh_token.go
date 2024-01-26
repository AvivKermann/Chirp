package models

type RefreshToken struct {
	Token    string `json:"token"`
	IsActive bool   `json:"is_active"`
}
