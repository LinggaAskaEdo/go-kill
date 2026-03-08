package entity

import "time"

type Payment struct {
	ID            string    `db:"id" json:"id"`
	OrderID       string    `db:"order_id" json:"order_id"`
	PaymentMethod string    `db:"payment_method" json:"payment_method"`
	Amount        float64   `db:"amount" json:"amount"`
	Status        string    `db:"status" json:"status"`
	TransactionID *string   `db:"transaction_id" json:"transaction_id,omitempty"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time `db:"updated_at" json:"updated_at"`
}
