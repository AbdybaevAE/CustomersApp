package server

import (
	"net/http"
	"strconv"
	"text/template"
	"time"

	"github.com/abdybaevae/customers-app/conf"
	"github.com/abdybaevae/customers-app/pkg/codes"
	"github.com/abdybaevae/customers-app/pkg/custval"
	"github.com/abdybaevae/customers-app/pkg/resp"
	customerservice "github.com/abdybaevae/customers-app/pkg/services/customer"
	"github.com/abdybaevae/customers-app/pkg/services/customer/dto"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

const birthDateLayout = "2006-01-02"

type handler struct {
	customerService customerservice.CustomerService
	templates       *template.Template
	log             *logrus.Entry
	Cfg             *conf.Config
}

type pages int

const (
	CustomerList pages = iota
	AddCustomer
	DeleteCustomer
	EditCustomer
	MessagePage
)

var respFactory = resp.GetHtmlResponseFactory()

type queryListData struct {
	Customers   []dto.ListCustomerResultItem
	Next        bool
	Prev        bool
	NextValue   int
	PrevValue   int
	SearchValue string
}

func (h *handler) showListPage(rw http.ResponseWriter, r *http.Request) {
	args := &dto.ListCustomersArguments{
		OrderBy:      "customer_first_name",
		Page:         0,
		OrderByValue: "asc",
	}
	data, err := h.customerService.QueryList(r.Context(), args)
	if err != nil {
		respFactory.Error(rw, err)
		return
	}
	tempData := &queryListData{
		Customers: data.Customers,
		Next:      len(data.Customers) == customerservice.CustomersPerPage,
		Prev:      false,
		NextValue: args.Page + 1,
	}
	h.templates.ExecuteTemplate(rw, "customers_list", tempData)
}
func (h *handler) queryList(rw http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		respFactory.CodeMessage(rw, codes.BadRequest, codes.KnownMessageBadRequest)
		return
	}
	queryArgs := &dto.ListCustomersArguments{
		OrderBy:      r.FormValue("orderBy"),
		SearchValue:  r.FormValue("searchValue"),
		OrderByValue: r.FormValue("orderByValue"),
	}
	if queryArgs.OrderByValue == "" {
		queryArgs.OrderByValue = "asc"
	}
	page, err := strconv.Atoi(r.FormValue("page"))
	if err != nil {
		page = 0
	}
	queryArgs.Page = page
	data, err := h.customerService.QueryList(r.Context(), queryArgs)
	if err != nil {
		respFactory.Error(rw, err)
		return
	}
	tempData := &queryListData{
		Customers:   data.Customers,
		Next:        len(data.Customers) == customerservice.CustomersPerPage,
		Prev:        queryArgs.Page != 0,
		NextValue:   queryArgs.Page + 1,
		PrevValue:   queryArgs.Page - 1,
		SearchValue: queryArgs.SearchValue,
	}
	h.templates.ExecuteTemplate(rw, "customers_list", tempData)
}

type AddCustomerPageData struct {
	MinDate string
	MaxDate string
}

func (h *handler) addCustomerPage(rw http.ResponseWriter, r *http.Request) {
	min, max := custval.ComputeBirthDateRange()
	data := &AddCustomerPageData{
		MinDate: min.Format(birthDateLayout),
		MaxDate: max.Format(birthDateLayout),
	}
	h.templates.ExecuteTemplate(rw, "create_customer", data)
}
func (h *handler) handleAddCustomer(rw http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		respFactory.CodeMessage(rw, codes.BadRequest, codes.KnownMessageBadRequest)
		return
	}
	birthDate, err := time.Parse(birthDateLayout, r.PostForm.Get("birthDate"))
	if err != nil {
		h.log.Error(err)
		respFactory.CodeMessage(rw, codes.InvalidData, codes.KnownMessageCustomerInvalidBirthDate)
		return
	}
	addArgs := &dto.CreateCustomerArguments{
		CustomerItem: dto.CustomerItem{
			FirstName: r.PostForm.Get("firstName"),
			LastName:  r.PostForm.Get("lastName"),
			BirthDate: birthDate,
			Gender:    r.PostForm.Get("gender"),
			Address:   r.PostForm.Get("address"),
			Email:     r.PostForm.Get("email"),
		},
	}
	err = h.customerService.Create(r.Context(), addArgs)
	if err == codes.UniqueConstraintViolation {
		respFactory.CodeMessage(rw, codes.BadRequest, codes.KnownMessageGivenEmailBusyUseAnotherOne)
		return
	}
	if err != nil {
		respFactory.Error(rw, err)
		return
	}
	respFactory.CodeMessage(rw, codes.Created, codes.KnownMessageCustomerCreated)
}

type EditCustomerPageData struct {
	Id        int
	FirstName string
	LastName  string
	BirthDate string
	Gender    string
	Address   string
	Hash      string
	MaxDate   string
	MinDate   string
}

func (h *handler) editCustomerPage(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	customerId, err := strconv.Atoi(vars["customerId"])
	if err != nil {
		respFactory.CodeMessage(rw, codes.NotFound, codes.KnownMessageNotFoundPage)
		return
	}
	customer, err := h.customerService.GetById(r.Context(), customerId)
	if err != nil {
		respFactory.Error(rw, err)
		return
	}

	min, max := custval.ComputeBirthDateRange()
	data := EditCustomerPageData{
		Id:        customer.Id,
		FirstName: customer.FirstName,
		LastName:  customer.LastName,
		BirthDate: customer.BirthDate.Format(birthDateLayout),
		Gender:    customer.Gender,
		Address:   customer.Address,
		Hash:      customer.Hash,
		MinDate:   min.Format(birthDateLayout),
		MaxDate:   max.Format(birthDateLayout),
	}

	h.templates.ExecuteTemplate(rw, "edit_customer", data)
}
func (h *handler) handleUpdateCustomer(rw http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		respFactory.CodeMessage(rw, codes.BadRequest, codes.KnownMessageBadRequest)
		return
	}
	vars := mux.Vars(r)
	customerId, err := strconv.Atoi(vars["customerId"])
	if err != nil {
		respFactory.CodeMessage(rw, codes.NotFound, codes.KnownMessageNotFoundPage)
		return
	}
	birthDate, err := time.Parse(birthDateLayout, r.PostForm.Get("birthDate"))
	if err != nil {
		respFactory.CodeMessage(rw, codes.InvalidData, codes.KnownMessageCustomerInvalidBirthDate)
		return
	}
	editArgs := &dto.UpdateCustomerArguments{
		Id:        customerId,
		FirstName: r.PostForm.Get("firstName"),
		LastName:  r.PostForm.Get("lastName"),
		Gender:    r.PostForm.Get("gender"),
		BirthDate: birthDate,
		Address:   r.PostForm.Get("address"),
		Hash:      r.PostForm.Get("hash"),
	}
	if err := h.customerService.Update(r.Context(), editArgs); err != nil {
		respFactory.Error(rw, err)
		return
	}
	respFactory.CodeMessage(rw, codes.Ok, codes.KnownMessageCustomerEdited)
}
func (h *handler) handleDeleteCustomer(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	customerId, err := strconv.Atoi(vars["customerId"])
	if err != nil {
		resp.GetHtmlResponseFactory().CodeMessage(rw, codes.NotFound, codes.KnownCustomerNotFound)
		return
	}
	if err := h.customerService.DeleteById(r.Context(), customerId); err != nil {
		respFactory.Error(rw, err)
		return
	}
	respFactory.CodeMessage(rw, codes.Ok, codes.KnownMessageCustomerDeleted)
}
