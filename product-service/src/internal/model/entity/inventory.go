package entity

type Inventory struct {
	ID               string `json:"id"`
	ProductID        string `json:"product_id"`
	Quantity         int    `json:"quantity"`
	ReservedQuantity int    `json:"reserved_quantity"`
	UpdatedAt        string `db:"updated_at" json:"updated_at"`
}
