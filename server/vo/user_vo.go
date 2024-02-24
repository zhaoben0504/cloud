package vo

type EmailCodeVo struct {
	Token string `json:"token"`
	Code  int    `json:"code"`
}
