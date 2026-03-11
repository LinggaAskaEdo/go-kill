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
	UserEmail   string      `json:"user_email"`
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

type Notification struct {
	UserID    string                 `bson:"user_id"`
	Type      string                 `bson:"type"`
	Category  string                 `bson:"category"`
	Title     string                 `bson:"title"`
	Message   string                 `bson:"message"`
	Metadata  map[string]interface{} `bson:"metadata"`
	Status    string                 `bson:"status"`
	SentAt    time.Time              `bson:"sent_at,omitempty"`
	CreatedAt time.Time              `bson:"created_at"`
}

type NotificationPreferences struct {
	UserID       string `bson:"user_id"`
	EmailEnabled bool   `bson:"email_enabled"`
	SMSEnabled   bool   `bson:"sms_enabled"`
	PushEnabled  bool   `bson:"push_enabled"`
}
