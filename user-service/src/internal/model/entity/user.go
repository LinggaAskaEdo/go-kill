package entity

import (
	"database/sql"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type User struct {
	ID        string `db:"id" json:"id"`
	AuthID    string `db:"auth_id" json:"auth_id"`
	Email     string `db:"email" json:"email"`
	FirstName string `db:"first_name" json:"first_name"`
	LastName  string `db:"last_name" json:"last_name"`
	CreatedAt string `db:"created_at" json:"created_at"`
	UpdatedAt string `db:"updated_at" json:"updated_at"`
}

type UserActivity struct {
	ID           bson.ObjectID          `bson:"_id,omitempty"`
	UserID       string                 `bson:"user_id"`
	ActivityType string                 `bson:"activity_type"`
	Metadata     map[string]interface{} `bson:"metadata"`
	Timestamp    time.Time              `bson:"timestamp"`
}

type UserAddress struct {
	ID            string         `db:"id" json:"id"`
	UserID        string         `db:"user_id" json:"user_id"`
	AddressType   string         `db:"address_type" json:"address_type"`
	StreetAddress string         `db:"street_address" json:"street_address"`
	City          string         `db:"city" json:"city"`
	State         sql.NullString `db:"state" json:"state"`
	PostalCode    string         `db:"postal_code" json:"postal_code"`
	Country       string         `db:"country" json:"country"`
	IsDefault     bool           `db:"is_default" json:"is_default"`
}
