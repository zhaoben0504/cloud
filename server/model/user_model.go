package model

type UserBasic struct {
	Id        *int64  `json:"id"`
	Identity  *string `json:"identity"`
	Name      *string `json:"name"`
	Password  *string `json:"password"`
	Email     *string `json:"email"`
	CreatedAt *int64  `json:"created"`
	UpdatedAt *int64  `json:"updated_at"`
	DeletedAt *int64  `json:"deleted_at"`
}
