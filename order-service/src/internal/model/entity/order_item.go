package entity

import "time"

type OrderItem struct {
	ID          string    `db:"id" json:"id"`
	OrderID     string    `db:"order_id" json:"order_id"`
	ProductID   string    `db:"product_id" json:"product_id"`
	ProductName string    `db:"product_name" json:"product_name"`
	Quantity    int       `db:"quantity" json:"quantity"`
	UnitPrice   float64   `db:"unit_price" json:"unit_price"`
	Subtotal    float64   `db:"subtotal" json:"subtotal"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
}
