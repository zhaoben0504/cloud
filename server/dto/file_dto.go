package dto

import "mime/multipart"

type UploadFileDTO struct {
	Token    string                `json:"token" form:"token" validate:"required"`
	ID       string                `json:"id" form:"id" validate:"required"`
	File     *multipart.FileHeader `validate:"required"`
	FileName *string               `json:"file_name" form:"file_name" validate:"required"`
	Path     *string               `json:"path" form:"path" validate:"required"`
}

type DownloadFileDTO struct {
	Token string `json:"token" validate:"required"`
	ID    string `json:"string" validate:"required"`
}

type DeleteFileDTO struct {
	Token string `json:"token"  validate:"required"`
	ID    string `json:"string" validate:"required"`
}

type ListFileDTO struct {
	Token   string `json:"token" validate:"required"`
	Keyword string `json:"keyword" `
	Page    int    `json:"page" `
	Rows    int    `json:"rows"`
}
