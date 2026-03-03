package dto

import (
	"database/sql"

	x "github.com/linggaaskaedo/go-kill/common/pkg/errors"
	"github.com/linggaaskaedo/go-kill/user-service/src/internal/model/entity"
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

type UserResp struct {
	ID        string `json:"id"  extensions:"x-order=0"`
	Email     string `json:"email"  extensions:"x-order=1"`
	FirstName string `json:"first_name"  extensions:"x-order=2"`
	LastName  string `json:"last_name"  extensions:"x-order=3"`
}

type UserRegResp struct {
	ID        string `json:"id"`
	AutdID    string `json:"auth_id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type UserActivity struct {
	Success    bool                   `json:"success"`
	Data       []*entity.UserActivity `json:"data"`
	Pagination Pagination             `json:"pagination"`
}

type Address struct {
	ID            string         `json:"id"`
	AddressType   string         `json:"address_type"`
	StreetAddress string         `json:"street_address"`
	City          string         `json:"city"`
	State         sql.NullString `json:"state,omitempty"`
	PostalCode    string         `json:"postal_code"`
	Country       string         `json:"country"`
	IsDefault     bool           `json:"is_default"`
}
