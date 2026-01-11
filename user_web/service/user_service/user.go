package user_service

type UserLoginRequest struct {
	Mobile    string `json:"mobile" binding:"required,mobile" ` //以使用 binding:"mobile" 这样的标签  自动调用验证函数
	Password  string `json:"password" binding:"required"`
	CaptchaId string `json:"captcha_id" binding:"required"`
	Answer    string `json:"answer" binding:"required"`
}

type UserRegisterRequest struct {
	Mobile   string `json:"mobile" binding:"required,mobile" `
	Password string `json:"password" binding:"required"`
	Code     string `json:"code" binding:"required"`
	Role     int32  `json:"role" binding:"required"`
}

type UserUpdateRequest struct {
	Id       int32  `json:"id,omitempty"`
	Password string `json:"password,omitempty"`
	NickName string `json:"nick_name,omitempty"`
	BirthDay uint64 `json:"birth_day,omitempty"`
	Gender   string `json:"gender,omitempty"`
	Role     int32  `json:"role,omitempty"`
}

type UserListResponse struct {
	Id       int32  `json:"id"`
	Password string `json:"password"`
	Mobile   string `json:"mobile"`
	NickName string `json:"nick_name"`
	BirthDay string `json:"birth_day"`
	Gender   string `json:"gender"`
	Role     int    `json:"role"`
}
