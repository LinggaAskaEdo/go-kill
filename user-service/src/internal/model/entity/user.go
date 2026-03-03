package entity

import (
	"database/sql"
	"time"
)

type User struct {
	ID        string `db:"id" json:"id"`
	AutdID    string `db:"auth_id" json:"auth_id"`
	Email     string `db:"email" json:"email"`
	FirstName string `db:"first_name" json:"first_name"`
	LastName  string `db:"last_name" json:"last_name"`
	CreatedAt string `db:"created_at" json:"created_at"`
	UpdatedAt string `db:"updated_at" json:"updated_at"`
}

type UserActivity struct {
	UserID       string                 `bson:"user_id"`
	ActivityType string                 `bson:"activity_type"`
	Metadata     map[string]interface{} `bson:"metadata"`
	Timestamp    time.Time              `bson:"timestamp"`
	CreatedAt    time.Time              `bson:"created_at"`
}

type UserAddress struct {
	ID            string         `db:"id" json:"id"`
	AddressType   string         `db:"address_type" json:"address_type"`
	StreetAddress string         `db:"street_address" json:"street_address"`
	City          string         `db:"city" json:"city"`
	State         sql.NullString `db:"state" json:"state"`
	PostalCode    string         `db:"postal_code" json:"postal_code"`
	Country       string         `db:"country" json:"country"`
	IsDefault     bool           `db:"is_default" json:"is_default"`
}
