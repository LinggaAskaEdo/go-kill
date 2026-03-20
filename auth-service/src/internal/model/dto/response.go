package dto

import (
	x "github.com/linggaaskaedo/go-kill/common/pkg/errors"
)

type Meta struct {
	Path       string      `json:"path" extensions:"x-order=0"`
	StatusCode int         `json:"status_code" extensions:"x-order=1"`
	Status     string      `json:"status" extensions:"x-order=2"`
	Message    string      `json:"message" extensions:"x-order=3"`
	Error      *x.AppError `json:"error,omitempty" swaggertype:"primitive,object" extensions:"x-order=4"`
	Timestamp  string      `json:"timestamp" extensions:"x-order=5"`
}

type HttpSuccessResp struct {
	Meta       Meta        `json:"metadata" extensions:"x-order=0"`
	Data       any         `json:"data,omitempty" extensions:"x-order=1"`
	Pagination *Pagination `json:"pagination,omitempty" extensions:"x-order=2"`
}

type HTTPErrorResp struct {
	Meta Meta `json:"metadata"`
}

type CreateAuthUserResponse struct {
	Success bool   `json:"success"`
	AuthId  string `json:"auth_id"`
}

type LoginResponse struct {
	Success      bool   `json:"success"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

type ValidateTokenResponse struct {
	Valid  bool   `json:"valid"`
	UserId string `json:"user_id"`
	Email  string `json:"email"`
}

type RefreshTokenResponse struct {
	Success      bool   `json:"success"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

type LogoutResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
