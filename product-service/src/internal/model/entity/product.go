package entity

import "time"

type Product struct {
	ID          string    `db:"id" json:"id"`
	Name        string    `db:"name" json:"name"`
	Description string    `db:"description" json:"description"`
	Price       float64   `db:"price" json:"price"`
	SKU         string    `db:"sku" json:"sku"`
	IsActive    bool      `db:"is_active" json:"is_active"`
	Categories  []string  `db:"-" json:"categories,omitempty"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}
