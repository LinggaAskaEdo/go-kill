package entity

import "time"

type OrderStatus string

const (
	StatusPending    OrderStatus = "pending"
	StatusConfirmed  OrderStatus = "confirmed"
	StatusProcessing OrderStatus = "processing"
	StatusShipped    OrderStatus = "shipped"
	StatusDelivered  OrderStatus = "delivered"
	StatusCancelled  OrderStatus = "cancelled"
)

type Order struct {
	ID                string       `db:"id" json:"id"`
	UserID            string       `db:"user_id" json:"user_id"`
	OrderNumber       string       `db:"order_number" json:"order_number"`
	Status            OrderStatus  `db:"status" json:"status"`
	TotalAmount       float64      `db:"total_amount" json:"total_amount"`
	ShippingAddressID *string      `db:"shipping_address_id" json:"shipping_address_id,omitempty"`
	BillingAddressID  *string      `db:"billing_address_id" json:"billing_address_id,omitempty"`
	CreatedAt         time.Time    `db:"created_at" json:"created_at"`
	UpdatedAt         time.Time    `db:"updated_at" json:"updated_at"`
	Items             []*OrderItem `db:"-" json:"items,omitempty"`
}
