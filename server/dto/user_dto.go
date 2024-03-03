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

type EmailCodeDTO struct {
	Email string `json:"email"`
}

type UserInfoDTO struct {
	Token string `json:"token" validate:"required"`
	Id    string `json:"id"`
}
