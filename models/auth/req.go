package auth

type LoginReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginFromGoogleAuthReq struct {
	AuthCode string `json:"auth_code" binding:"required"`
}
