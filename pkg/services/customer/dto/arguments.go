package dto

import (
	"time"
)

// It's always better to have data transfer object per each layer. services has their own dto, repos deals with entities
// more flexibility when you gonna change some endpoints / entities.
type CustomerItem struct {
	FirstName string    `validate:"required,max=100"`
	LastName  string    `validate:"required,max=100"`
	BirthDate time.Time `validate:"required"`
	Gender    string    `validate:"required,oneof=female male"`
	Email     string    `validate:"required,email"`
	Address   string
	Hash      string
}
type CreateCustomerArguments struct {
	CustomerItem
}
type UpdateCustomerArguments struct {
	Id        int       `validate:"required"`
	FirstName string    `validate:"required,max=100"`
	LastName  string    `validate:"required,max=100"`
	BirthDate time.Time `validate:"required"`
	Gender    string    `validate:"required,oneof=female male"`
	// random hash string length must be syncronized here too
	Hash    string `validate:"required,len=20"`
	Address string
}
type ListCustomersArguments struct {
	Page         int    `validate:"min=0"`
	SearchValue  string `validate:"max=100"`
	OrderBy      string `validate:"required,oneof=customer_first_name customer_last_name customer_birth_date customer_address customer_email"`
	OrderByValue string `validate:"required,oneof=asc desc"`
}
type ListCustomerResultItem struct {
	Id        int
	Email     string
	FirstName string
	LastName  string
	Address   string
	BirthDate string
	Gender    string
}
type ListCustomersResult struct {
	Customers []ListCustomerResultItem
}
type GetByIdResult struct {
	CustomerItem
}
