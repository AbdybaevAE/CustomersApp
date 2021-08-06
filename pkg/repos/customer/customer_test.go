package customer

import (
	"context"
	"log"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/abdybaevae/customers-app/pkg/models"
	"github.com/jmoiron/sqlx"

	"testing"
)

var newCustomer = &models.Customer{
	FirstName: "FirstName",
	LastName:  "LastName",
	Address:   "Address",
	BirthDate: time.Now(),
	Gender:    "male",
	Hash:      "hash",
	Email:     "email@gmail.com",
}

func TestCreateCustomer(t *testing.T) {
	db, mock := conn()
	repo := New(db)
	defer db.Close()
	mock.ExpectExec("insert into customers").WithArgs(newCustomer.FirstName,
		newCustomer.LastName,
		newCustomer.BirthDate,
		newCustomer.Gender,
		newCustomer.Email,
		newCustomer.Address,
		newCustomer.Hash,
	).WillReturnResult(sqlmock.NewResult(1, 1))
	if err := repo.Create(context.Background(), newCustomer); err != nil {
		t.Error("error while inserting", err)
	}
}

// func TestQueryList(t *testing.T) {
// 	db, mock := conn()
// 	defer db.Close()
// 	repo := New(db)
// 	offset, orderBy, orderByValue, limit := 10, "customer_first_name", "asc", 30
// 	mock.ExpectQuery("SELECT * FROM customers ORDER BY customer_first_name asc OFFSET \\? LIMIT \\?").WithArgs(offset, limit).WillReturnRows(sqlmock.NewRows([]string{"customer_id", "customer_first_name"}))
// 	err, res := repo.QueryList(context.Background(), offset, orderBy, orderByValue, limit)
// 	if err != nil {
// 		t.Error("error while query list", err)
// 	}
// 	t.Log(res)
// }
func conn() (*sqlx.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("en error %s was not expted ", err)
	}
	return sqlx.NewDb(db, "sqlmock"), mock
}
