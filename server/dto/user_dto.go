package dto

type UserLoginDTO struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type UserRegisterDTO struct {
	Name     string `json:"name"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Code     string `json:"code"`
}
