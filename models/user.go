package models

type UserParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type User struct {
	ResponseUser
	Password []byte `json:"password"`
}

type ResponseUser struct {
	Email string `json:"email"`
	ID    int    `json:"id"`
}
