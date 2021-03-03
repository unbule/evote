package model

type Register struct {
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
	Email    string `json:"email"  form:"email"`
	Status   string `json:"status"  form:"status"`
}
