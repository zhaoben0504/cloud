package server

import (
	"encoding/json"
	"fmt"
)

const (
	ServiceID           = "CLOUD_ID"
	OkCode              = 0
	InternalErrCode     = 100001
	ParamErrCode        = 100002
	TokenInvalidErrCode = 100006
	PermissionErrCode   = 100007
	SQLErrCode          = 200001
)

type Empty struct {
}

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"msg"`
}

func (e Error) Error() string {
	return fmt.Sprintf("Error{code: %d, message: %s}", e.Code, e.Message)
}

// Response response
type Response struct {
	Error
	Data interface{} `json:"data"` // 数据
}

// String string
func (r *Response) String() string {
	b, err := json.Marshal(r)
	if nil != err {
		return string(b)
	}
	return fmt.Sprintf("%+v", *r)
}

func NewOK(lang string, data interface{}) Response {
	return Response{
		Error: Error{
			Code:    OkCode,
			Message: server.bundle.GetMsgByCode(lang, OkCode),
		},
		Data: data,
	}
}

func NewError(lang string, code int) Response {
	return Response{
		Error: Error{
			Code:    code,
			Message: server.bundle.GetMsgByCode(lang, code),
		},
		Data: Empty{},
	}
}

type AddBaseVO struct {
	ID int64 `json:"id,string"`
}

type RedisUserInfo struct {
	Platform    *int     `json:"platform"`
	AppVersion  string   `json:"app_version"`
	SysVersion  string   `json:"sys_version"`
	DeviceID    string   `json:"device_id"`
	Language    string   `json:"language"`
	UID         string   `json:"uid"`
	Name        string   `json:"name"`
	Account     string   `json:"account"`
	Company     string   `json:"company"`
	Phone       string   `json:"phone"`
	Permissions []string `json:"permissions"`
	SysModel    string   `json:"sys_model"`
	IP          string   `json:"ip"`
	Type        int      `json:"type"`
	Index       int64    `json:"index"`
}
