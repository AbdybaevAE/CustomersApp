package main

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/abdybaevae/customers-app/conf"
	"github.com/abdybaevae/customers-app/internal/server"

	"github.com/abdybaevae/customers-app/internal/db"
	customerrepo "github.com/abdybaevae/customers-app/pkg/repos/customer"
	customerservice "github.com/abdybaevae/customers-app/pkg/services/customer"

	_ "github.com/brianvoe/gofakeit/v6"
	"github.com/sirupsen/logrus"
)

func main() {
	// define root context that can be cancelled with releasing all resources
	ctx, cancel := context.WithCancel(context.Background())
	log := logrus.WithContext(ctx)
	defer cancel()
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)

	cfg := conf.Load()
	dbConn := db.Connect(cfg)

	customerRepo := customerrepo.New(dbConn)
	customerService := customerservice.New(customerRepo, log)
	handler := server.NewHandler(customerService, log)

	// Run migrations
	if err := db.HandleMigrations(cfg, customerService, dbConn); err != nil {
		log.Fatal(err)
	}

	srv := &http.Server{
		Addr:    cfg.ServerAddress,
		Handler: handler,
		BaseContext: func(l net.Listener) context.Context {
			return ctx
		},
	}

	go func() {
		sig := <-ch
		log.Infof("handle signal %v, exiting", sig)
		srv.Shutdown(ctx)
		cancel()
	}()
	log.Printf("Srn service started on port %v", cfg.ServerAddress)

	log.Fatal(ctx, srv.ListenAndServe())

}
