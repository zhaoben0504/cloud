package vo

type UploadFileVO struct {
	FileAddr string `json:"file_addr"`
}

// FileListVO File list
type FileListVO struct {
	Total int64      `json:"total"`
	List  []FileList `json:"list,omitempty"`
}

type FileList struct {
	ID        *int64  `json:"id" xorm:"pk"`
	FileName  *int64  `json:"file_name"`
	Path      *int    `json:"path"`
	Md5       *int    `json:"MD5"`
	Uid       *string `json:"uid"`
	UserName  *string `json:"user_name"`
	CreatedAt *int64  `json:"created_at"`
	UpdatedAt *int64  `json:"updated_at"`
	Deleted   *int    `json:"-"`
}
