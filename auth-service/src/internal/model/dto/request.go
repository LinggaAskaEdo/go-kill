package dto

type UserLoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type CreateAuthUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type LoginRequest struct {
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=8"`
	IpAddress string `json:"ip_address"`
	UserAgent string `json:"user_agent"`
}

type ValidateTokenRequest struct {
	Token string `json:"token" binding:"required"`
}

type LogoutRequest struct {
	Token  string `json:"token" binding:"required"`
	UserId string `json:"user_id" binding:"required"`
}
