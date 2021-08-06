package server

import (
	"net/http"

	customerservice "github.com/abdybaevae/customers-app/pkg/services/customer"
	"github.com/abdybaevae/customers-app/pkg/utils"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	_ "github.com/urfave/negroni"
)

func NewHandler(customerService customerservice.CustomerService, log *logrus.Entry) http.Handler {
	router := mux.NewRouter()
	templates := utils.LoadTemplates()
	h := &handler{
		customerService: customerService,
		templates:       templates,
		log:             log,
	}
	fs := http.FileServer(http.Dir("./ui/static"))
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))
	router.HandleFunc("/", h.showListPage).Methods(http.MethodGet)
	router.HandleFunc("/", h.queryList).Methods(http.MethodPost)
	router.HandleFunc("/customers/add", h.addCustomerPage).Methods(http.MethodGet)
	router.HandleFunc("/customers/add", h.handleAddCustomer).Methods(http.MethodPost)
	router.HandleFunc("/customers/{customerId}/edit", h.editCustomerPage).Methods(http.MethodGet)
	router.HandleFunc("/customers/{customerId}/edit", h.handleUpdateCustomer).Methods(http.MethodPost)
	router.HandleFunc("/customers/{customerId}/delete", h.handleDeleteCustomer).Methods(http.MethodPost)
	return router
}
