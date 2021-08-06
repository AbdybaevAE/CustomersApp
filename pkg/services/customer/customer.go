package customer

import (
	"context"
	"database/sql"

	"github.com/abdybaevae/customers-app/pkg/utils"

	"github.com/abdybaevae/customers-app/pkg/codes"
	"github.com/abdybaevae/customers-app/pkg/custval"
	"github.com/abdybaevae/customers-app/pkg/models"
	customerrepo "github.com/abdybaevae/customers-app/pkg/repos/customer"
	"github.com/abdybaevae/customers-app/pkg/services/customer/dto"

	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

const CustomersPerPage = 20

// Customer service interface, it can do below things. As data come to untrusted resources it will be better to validate
// data inside given service.
type CustomerService interface {
	// Create customer
	Create(ctx context.Context, customer *dto.CreateCustomerArguments) (err error)
	// Delete customer by id
	DeleteById(ctx context.Context, customerId int) (err error)
	// update customer(arguments includes version, which can handle properly overriding values)
	Update(ctx context.Context, args *dto.UpdateCustomerArguments) (err error)
	// query customers list(sorting by customer fields + search on firstName and lastName)
	QueryList(ctx context.Context, args *dto.ListCustomersArguments) (result *dto.ListCustomersResult, err error)
	// get detailed information by customer id(including version)
	GetById(ctx context.Context, customerId int) (customer *models.Customer, err error)
}

// Following documentation, it will be better to have single instance of validation that caches struct info
var validate = validator.New()

// Current implementation of customer service
type service struct {
	customerRepo customerrepo.CustomerRepo
	log          *logrus.Entry
}

// Main constructor for service, which applies customer repository as function arguments(di)
func New(customerRepo customerrepo.CustomerRepo, log *logrus.Entry) CustomerService {
	return &service{customerRepo: customerRepo,
		log: log,
	}
}
func (s *service) Create(ctx context.Context, customer *dto.CreateCustomerArguments) error {
	if err := validate.Struct(customer); err != nil {
		return codes.NewErr(codes.InvalidData, err.Error())
	}
	if !custval.IsValidBirthDate(customer.BirthDate) {
		return codes.NewErr(codes.InvalidData, codes.KnownMessageCustomerInvalidAge)
	}
	customerEntity := &models.Customer{
		FirstName: customer.FirstName,
		LastName:  customer.LastName,
		BirthDate: customer.BirthDate,
		Gender:    customer.Gender,
		Email:     customer.Email,
		Address:   customer.Address,
		Hash:      utils.GenCustomerHash(),
	}
	return s.customerRepo.Create(ctx, customerEntity)
}
func (s *service) DeleteById(ctx context.Context, customerId int) error {
	if err := s.customerRepo.DeleteById(ctx, customerId); err != nil {
		if err == sql.ErrNoRows {
			return codes.NewErr(codes.CustomerNotFound, codes.KnownCustomerNotFound)
		}
		return err
	}
	return nil
}
func (s *service) Update(ctx context.Context, customer *dto.UpdateCustomerArguments) error {
	if err := validate.Struct(customer); err != nil {
		return codes.NewErr(codes.InvalidData, err.Error())
	}
	if !custval.IsValidBirthDate(customer.BirthDate) {
		return codes.NewErr(codes.InvalidData, codes.KnownMessageCustomerInvalidAge)
	}
	customerEntity := &models.Customer{
		Id:        customer.Id,
		FirstName: customer.FirstName,
		LastName:  customer.LastName,
		Gender:    customer.Gender,
		Address:   customer.Address,
		Hash:      customer.Hash,
		BirthDate: customer.BirthDate,
	}
	if err := s.customerRepo.Update(ctx, customerEntity); err != nil {
		if err == codes.NoRowsModified {
			return codes.NewErr(codes.CustomerNotFound, codes.KnownMessageEditCustomerConflict)
		}
		return err
	}
	return nil
}
func (s *service) QueryList(ctx context.Context, args *dto.ListCustomersArguments) (*dto.ListCustomersResult, error) {
	if err := validate.Struct(args); err != nil {
		return nil, codes.NewErr(codes.InvalidData, err.Error())
	}
	var customers []models.Customer
	var err error
	if args.SearchValue == "" {
		customers, err = s.customerRepo.QueryList(ctx, args.Page*CustomersPerPage, args.OrderBy, args.OrderByValue, CustomersPerPage)
	} else {
		customers, err = s.customerRepo.SearchQueryList(ctx, args.Page*CustomersPerPage, args.OrderBy, args.OrderByValue, CustomersPerPage, args.SearchValue)
	}
	if err != nil {
		return nil, err
	}
	res := &dto.ListCustomersResult{
		Customers: []dto.ListCustomerResultItem{},
	}
	for _, v := range customers {
		res.Customers = append(res.Customers, dto.ListCustomerResultItem{
			Id:        v.Id,
			Email:     v.Email,
			FirstName: v.FirstName,
			LastName:  v.LastName,
			Gender:    v.Gender,
			BirthDate: utils.FormatBirthDate(v.BirthDate),
			Address:   v.Address,
		})
	}
	return res, nil
}
func (s *service) GetById(ctx context.Context, customerId int) (*models.Customer, error) {
	customer, err := s.customerRepo.GetById(ctx, customerId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, codes.NewErr(codes.CustomerNotFound, codes.KnownCustomerNotFound)
		} else {
			return nil, codes.NewErr(codes.ServerInternal, codes.KnownMessageSomethingWrongHappened)
		}
	}

	return customer, nil
}
