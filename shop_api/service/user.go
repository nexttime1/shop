package service

type UserLoginRequest struct {
	Mobile   string `json:"mobile" binding:"required,mobile" ` //以使用 binding:"mobile" 这样的标签  自动调用验证函数
	Password string `json:"password" binding:"required"`
}
