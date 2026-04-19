package dto

type UserRegistration struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserCredentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
