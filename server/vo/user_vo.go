package vo

type UserLoginVo struct {
	Token string `json:"token"`
	Code  int    `json:"code"`
}
