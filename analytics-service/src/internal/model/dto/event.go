package dto

import "time"

type OrderEvent struct {
	EventID   string    `json:"event_id"`
	EventType string    `json:"event_type"`
	Timestamp time.Time `json:"timestamp"`
	Data      OrderData `json:"data"`
}

type OrderData struct {
	OrderID     string      `json:"order_id"`
	OrderNumber string      `json:"order_number"`
	UserID      string      `json:"user_id"`
	TotalAmount float64     `json:"total_amount"`
	Status      string      `json:"status"`
	Items       []OrderItem `json:"items"`
}

type OrderItem struct {
	ProductID   string  `json:"product_id"`
	ProductName string  `json:"product_name"`
	Quantity    int     `json:"quantity"`
	UnitPrice   float64 `json:"unit_price"`
}

type OrderAnalytics struct {
	Date            time.Time       `bson:"date"`
	Metrics         Metrics         `bson:"metrics"`
	HourlyBreakdown []HourlyMetrics `bson:"hourly_breakdown"`
	UpdatedAt       time.Time       `bson:"updated_at"`
}

type Metrics struct {
	TotalOrders     int     `bson:"total_orders"`
	TotalRevenue    float64 `bson:"total_revenue"`
	AverageOrderVal float64 `bson:"average_order_value"`
	CancelledOrders int     `bson:"cancelled_orders"`
	CompletedOrders int     `bson:"completed_orders"`
}

type HourlyMetrics struct {
	Hour    int     `bson:"hour"`
	Orders  int     `bson:"orders"`
	Revenue float64 `bson:"revenue"`
}
