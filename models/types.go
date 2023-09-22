package models

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Loginresponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}
