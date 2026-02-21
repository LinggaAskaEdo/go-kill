package entity

type User struct {
	ID        string `db:"id" json:"id"`
	AutdID    string `db:"autd_id" json:"autd_id"`
	Email     string `db:"email" json:"email"`
	FirstName string `db:"first_name" json:"first_name"`
	LastName  string `db:"last_name" json:"last_name"`
	CreatedAt string `db:"created_at" json:"created_at"`
	UpdatedAt string `db:"updated_at" json:"updated_at"`
}
