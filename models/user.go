package models

type UserParams struct {
	Email           string `json:"email"`
	Password        string `json:"password"`
	ExpireInSeconds *int   `json:"expires_in_seconds"`
}
type User struct {
	ResponseUser
	Password []byte `json:"password"`
}
type ResponseUser struct {
	Email string `json:"email"`
	ID    int    `json:"id"`
}
