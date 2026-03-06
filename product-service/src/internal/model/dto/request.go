package dto

type CreateProductRequest struct {
	Name        string   `json:"name" binding:"required"`
	Description string   `json:"description" binding:"required"`
	Price       float64  `json:"price" binding:"required"`
	SKU         string   `json:"sku" binding:"required"`
	IsActive    bool     `json:"is_active"`
	Categories  []string `json:"categories,omitempty"`
}
