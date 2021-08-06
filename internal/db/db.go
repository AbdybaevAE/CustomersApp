package db

import (
	"context"
	"fmt"
	"time"

	"github.com/abdybaevae/customers-app/conf"
	customerservice "github.com/abdybaevae/customers-app/pkg/services/customer"
	"github.com/abdybaevae/customers-app/pkg/services/customer/dto"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

func Connect(cfg *conf.Config) *sqlx.DB {
	var connStr = fmt.Sprintf("postgresql://%s:%s@%s/%s?sslmode=disable", cfg.DbUser, cfg.DbPassword, cfg.DbHost, cfg.DbName)
	return sqlx.MustConnect("postgres", connStr)
}
func HandleMigrations(cfg *conf.Config, customerService customerservice.CustomerService, db *sqlx.DB) error {
	log := logrus.New()
	log.Infof("start migrations")
	m, err := migrate.New(
		"file://resources/db/migrations",
		fmt.Sprintf("postgresql://%s:%s@%s/%s?sslmode=disable", cfg.DbUser, cfg.DbPassword, cfg.DbHost, cfg.DbName))
	if err != nil {
		return err
	}
	if err := m.Steps(1); err != nil && err != migrate.ErrNoChange {
		log.Errorf("migration error %v", err)
	}
	var count int

	if err := db.QueryRow("select count(*) from customers").Scan(&count); err != nil {
		return err
	}
	log.Infof("count of customer is %v", count)
	if count == 0 {
		log.Info("Create fake customers...")
		for i := 0; i < 1000; i++ {
			from, to := time.Now().AddDate(-59, 0, 0), time.Now().AddDate(-20, 0, 0)
			if err := customerService.Create(context.Background(), &dto.CreateCustomerArguments{
				CustomerItem: dto.CustomerItem{
					FirstName: gofakeit.Person().FirstName,
					LastName:  gofakeit.Person().LastName,
					Gender:    gofakeit.Person().Gender,
					BirthDate: gofakeit.DateRange(from, to),
					Email:     gofakeit.Email(),
					Address:   gofakeit.Person().Address.Address,
				},
			}); err != nil {
				log.Errorf("error or creating customer %v", err)
				return err
			}
		}
	}
	return nil
}
