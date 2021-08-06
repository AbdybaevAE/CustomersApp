package customer

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/abdybaevae/customers-app/pkg/codes"
	"github.com/abdybaevae/customers-app/pkg/models"
	"github.com/abdybaevae/customers-app/pkg/utils"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

// Simple customer repository that works with customer entity.
type CustomerRepo interface {
	Create(ctx context.Context, data *models.Customer) (err error)
	Update(ctx context.Context, data *models.Customer) (err error)
	DeleteById(ctx context.Context, customerId int) (err error)
	GetById(ctx context.Context, customerId int) (customer *models.Customer, err error)
	// query customers without search pattern
	QueryList(ctx context.Context, offset int, orderBy string, orderByValue string, limit int) ([]models.Customer, error)
	// query customers with search pattern
	SearchQueryList(ctx context.Context, offset int, orderBy string, orderByValue string, limit int, pattern string) ([]models.Customer, error)
}
type repo struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) CustomerRepo {
	return &repo{
		db,
	}
}

// predefined string queries to db.
const createCustomerQuery = `
insert into customers
	(
		customer_first_name,
		customer_last_name,
		customer_birth_date,
		customer_gender,
		customer_email,
		customer_address,
		customer_hash
	) values 
	(
		:customer_first_name,
		:customer_last_name,
		:customer_birth_date,
		:customer_gender,
		:customer_email,
		:customer_address,
		:customer_hash
	)	
`

const getByIdQuery = `
select * from customers where customer_id = $1
`

func (r *repo) GetById(ctx context.Context, customerId int) (*models.Customer, error) {
	customer := &models.Customer{}
	err := r.db.GetContext(ctx, customer, getByIdQuery, customerId)
	return customer, err
}
func (r *repo) Create(ctx context.Context, customer *models.Customer) error {
	_, err := r.db.NamedExecContext(ctx, createCustomerQuery, customer)
	// this is email duplication error code, reject customer creation
	if err, ok := err.(*pq.Error); ok && err.Code == "23505" {
		return codes.UniqueConstraintViolation
	}
	return err
}

// prevented overwrite update query
const updateCustomerQuery = `
update customers
set 
	customer_first_name = $1,
	customer_last_name = $2,
	customer_gender = $3,
	customer_address = $4,
	customer_birth_date = $5,
	customer_hash = $6
where 
	customer_id = $7 
	and 
	customer_hash = $8
`

func (r *repo) Update(ctx context.Context, customer *models.Customer) error {
	res, err := r.db.ExecContext(ctx, updateCustomerQuery, customer.FirstName, customer.LastName, customer.Gender,
		customer.Address, customer.BirthDate, utils.GenCustomerHash(), customer.Id, customer.Hash)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	// tricky part, it means either customer do not exist with given id or hash already changed(but customer exists), it's possible to create solution to differentiate
	// this situations, but let's think this isn't our case and just return user already chaged error.
	if count == 0 {
		return codes.NoRowsModified
	}
	return nil

}

const deleteCustomerQuery = `
delete 
from customers
where customer_id = $1
`

func (r *repo) DeleteById(ctx context.Context, customerId int) error {
	res, err := r.db.ExecContext(ctx, deleteCustomerQuery, customerId)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return codes.NoRowsModified
	}
	return nil
}

const queryListQuery = `
	SELECT * FROM customers 
	ORDER BY %s %s
	OFFSET $1
	LIMIT $2
`

func (r *repo) QueryList(ctx context.Context, offset int, orderBy string, orderByValue string, limit int) ([]models.Customer, error) {
	customers := []models.Customer{}
	actQuery := fmt.Sprintf(queryListQuery, orderBy, orderByValue)
	err := r.db.SelectContext(ctx, &customers, actQuery, offset, limit)
	return customers, err
}

// It makes sense that if you wanna find customer by firstName or lastName you exactly start searching by typing starts of their names
// It's means there are less chances when someone wants to search customers by providing substring that is not prefix of firstname or lastname
// that is why I decided to choose simple postgresql search by pattern function
// Generally we can use full text search by postgres(tsvector?) or elasticsearch to make efficient search mechanism or create custom solution...
// but for current implementation I decided to use simple search approach.
// Suppose user want to search customers by providing string with spaces. It seems reasonably that user want to search customers(or list of customers in this criteria)
// So, let's split this string and search for all occurences.
func (r *repo) SearchQueryList(ctx context.Context, offset int, orderBy string, orderByValue string, limit int, pattern string) ([]models.Customer, error) {
	if orderByValue == "" {
		return nil, codes.BadSearchCriteria
	}
	tokens := strings.Split(pattern, " ")
	var sb strings.Builder
	sb.WriteString("select * from customers where ")
	tokenLen := len(tokens)
	args := []interface{}{}
	seenTokens := map[string]bool{}
	for i, token := range tokens {
		token = strings.ToLower(strings.TrimSpace(token))
		if token == "" || seenTokens[token] {
			continue
		}
		if i != 0 {
			sb.WriteString(" or ")
		}
		sb.WriteString("customer_first_name ilike '%' || $" + strconv.Itoa(2*i+1) + " || '%' or ")
		sb.WriteString("customer_last_name ilike '%' || $" + strconv.Itoa(2*i+2) + " || '%'")
		args = append(args, token, token)
	}
	sb.WriteString(" order by " + orderBy + " " + orderByValue + " offset $" + strconv.Itoa(2*tokenLen+1) + " limit $" + strconv.Itoa(2*tokenLen+2))
	args = append(args, offset, limit)
	ret := []models.Customer{}
	err := r.db.SelectContext(ctx, &ret, sb.String(), args...)
	return ret, err
}
