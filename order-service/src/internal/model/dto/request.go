package dto

type CreateOrderRequest struct {
	UserID            string       `json:"user_id"`
	ShippingAddressID string       `json:"shipping_address_id"`
	BillingAddressID  string       `json:"billing_address_id"`
	PaymentMethod     string       `json:"patment_method"`
	Items             []*OrderItem `json:"items"`
}

type GetOrderRequest struct {
	OrderID string `json:"order_id"`
	UserID  string `json:"user_id"`
}

type ListOrderRequest struct {
	UserID string `json:"user_id"`
	Limit  int32  `json:"limit"`
	Offset int32  `json:"offset"`
}

type CancelOrderRequest struct {
	OrderID string `json:"order_id"`
	UserID  string `json:"user_id"`
	Reason  string `json:"reason"`
}
