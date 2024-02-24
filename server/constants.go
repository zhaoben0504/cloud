package server

import (
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt"
)

// jwt加密，如果是自定义结构体需要实现jwt.StandardClaims结构体，里面可以定义过期时间颁发者等等
type UserClaim struct {
	jwt.StandardClaims
	Id       int64
	Identity string
	Name     string
}

var JwtKey = "my_cloud_disk"

// 验证码长度
var EmailCodeLen = 6

// 验证码过期时间
var CodeExprie = 300

// 腾讯云
var CloudKey = "TECENTCLOUDSECRETKEY"
var CloudId = "TECENTCLOUDSECRETID"
var COSADDR = "https://my-storage-1306331535.cos.ap-nanjing.myqcloud.com"

// 默认的分页
var Pagesize int = 20

var DateTime = "2006-01-02 15:04:5"

var TokenExpire int64 = 3600 * 12
var RefreshTokenExpire int64 = 3600 * 24

const (
	ServiceID       = "CLOUD_ID"
	OkCode          = 0
	InternalErrCode = 100001
	ParamErrCode    = 100002
	SQLErrCode      = 100003

	TokenInvalidErrCode     = 200001
	PermissionErrCode       = 200002
	VerificationCodeErrCode = 200003
	UserAlreadyExistErrCode = 200004
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
