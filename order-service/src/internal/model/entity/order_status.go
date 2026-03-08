package entity

import "time"

type OrderStatusHistory struct {
	ID        string    `db:"id" json:"id"`
	OrderID   string    `db:"order_id" json:"order_id"`
	Status    string    `db:"status" json:"status"`
	Note      *string   `db:"note" json:"note,omitempty"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}
