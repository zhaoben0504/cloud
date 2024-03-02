package dto

import "mime/multipart"

type UploadFileDTO struct {
	Token    string                `json:"token"`
	ID       int64                 `json:"id,string" form:"id" validate:"required"`
	File     *multipart.FileHeader `validate:"required"`
	FileName *string               `json:"file_name" form:"file_name" validate:"required"`
	Path     *string               `json:"path" form:"path" validate:"required"`
}

type DownloadFileDTO struct {
	Token string `json:"token"`
	ID    string `json:"string" form:"id" validate:"required"`
}
