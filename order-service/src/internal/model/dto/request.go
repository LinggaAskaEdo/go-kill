package dto

import "time"

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

type ProductDetails struct {
	ID    string
	Name  string
	Price float64
}

type OrderEvent struct {
	EventID   string    `json:"event_id"`
	EventType string    `json:"event_type"`
	Version   string    `json:"version"`
	Timestamp time.Time `json:"timestamp"`
	Source    string    `json:"source"`
	Data      OrderData `json:"data"`
}

type OrderData struct {
	OrderID     string      `json:"order_id"`
	OrderNumber string      `json:"order_number"`
	UserID      string      `json:"user_id"`
	UserEmail   string      `json:"user_email"`
	TotalAmount float64     `json:"total_amount"`
	Status      string      `json:"status"`
	Items       []OrderItem `json:"items"`
}

type OrderItem struct {
	ProductID   string  `json:"product_id"`
	ProductName string  `json:"product_name"`
	Quantity    int32   `json:"quantity"`
	UnitPrice   float64 `json:"unit_price"`
}
