package model

// UploadFile db fle object
type UploadFile struct {
	ID        *int64  `json:"id,string" xorm:"pk"`
	FileName  *string `json:"file_name"`
	Path      *string `json:"path"`
	Uid       *int64  `json:"uid"`
	CreatedAt *int64  `json:"created_at"`
	UpdatedAt *int64  `json:"updated_at"`
	DeletedAt *int64  `json:"deleted_at"`
}
