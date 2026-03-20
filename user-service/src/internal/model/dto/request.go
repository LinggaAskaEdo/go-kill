package dto

type RegisterUserRequest struct {
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=8"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
}

type CreateUserAddress struct {
	AddressType   string `json:"address_type" binding:"required,oneof=shipping billing both"`
	StreetAddress string `json:"street_address" binding:"required"`
	City          string `json:"city" binding:"required"`
	State         string `json:"state"`
	PostalCode    string `json:"postal_code" binding:"required"`
	Country       string `json:"country" binding:"required"`
	IsDefault     bool   `json:"is_default"`
}
