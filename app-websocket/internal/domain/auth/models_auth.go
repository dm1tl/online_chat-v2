package auth

//Models for register user
type CreateUserReq struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateUserResp struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
}

//Models for login user
type LoginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResp struct {
	Token string `json:"token"`
}

//Models for token validation
type ValidateTokenReq struct {
	Token string `json:"token"`
}

type ValidateTokenResp struct {
	ID int64 `json:"id"`
}
