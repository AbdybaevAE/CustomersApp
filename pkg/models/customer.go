package models

import (
	"time"
)

// Customer entity
// customer has hash string which generates randomly every time user was edited.
// With this simple approach it's possbile
// prevent overwrite errors. During customer editing frontend part receives current customer hash, make changes
// and send's back to customer service. Customer service updates customer by id and received customer hash. If hash was changed
// (someone already edited user) update query don't match customer and update will be cancelled.
type Customer struct {
	Id        int       `db:"customer_id"`
	FirstName string    `db:"customer_first_name"`
	LastName  string    `db:"customer_last_name"`
	BirthDate time.Time `db:"customer_birth_date"`
	Gender    string    `db:"customer_gender"`
	Email     string    `db:"customer_email"`
	Address   string    `db:"customer_address"`
	CreatedAt time.Time `db:"customer_created_at"`
	UpdatedAt time.Time `db:"customer_updated_at"`
	Hash      string    `db:"customer_hash"`
}
