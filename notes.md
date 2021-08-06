Unfortunately there are only 3 unit tests per project, I was spended a lot of time on frontend part and html/template part.
As we use templates we need to add csrf.
Something like this:
```go
http.ListenAndServe(":8000",
        csrf.Protect([]byte("32-byte-long-auth-key"))(r))
    

// and this

t.ExecuteTemplate(w, "signup_form.tmpl", map[string]interface{}{
        csrf.TemplateTag: csrf.TemplateField(r),
    })
```
And then send token in hidden field(and validaate on server). In order to protect from attacks. But I didn't implemented it.
I wasn't worked with html/template.
You can start server with docker compose(in project root dir)
```sh
docker-compose up
```

Generally Service has following abstractions:
- Customer service 
- Customer repository
- Customer entity

I'm using gorilla multiplexer to handle routes. Every "controller"(except some that don't need data at all) reads data from form(we can change data source to json/xml and so on...)
and then pass that to service function call, where data is validating first and then do what it need's to do.
Including that I decided to have simple message template to show succes message, error message. 
Services can perform actions and return values and error. Typed error(like birthdate wrong value or email exist error) are returned from services as <code>ErrorCode</code> value which is <code>error</code> interface.
```go
type ErrorCode interface {
	error
	Code() Code
	Message() string
}
```
where Code - error string short code and message - human understandable error
Data validation was performed on both frontend and backend(as default behaviour)

Besides that I have following error codes and known messages:
```go
const (
	KnownMessageSomethingWrongHappened      = "Something went wrong, please try later."
	KnownMessageInvalidPageProvided         = "Invalid page provided, page must be positive integer."
	KnownMessageGivenEmailBusyUseAnotherOne = "Provided email address already in use, please provide another one."
	KnownMessageNotFoundPage                = "Given page doesn't exist."
	KnownMessageBadRequest                  = "Wrong request."
	KnownMessageCustomerCreated             = "Customer was successfully created."
	KnownMessageCustomerEdited              = "Customer was successfully edited."
	KnownMessageCustomerDeleted             = "Customer was successfully deleted."
	KnownMessageCustomerInvalidBirthDate    = "Customer birthdate must be of format yyyy-MM-dd."
	KnownMessageCustomerInvalidAge          = "Customer age must be between 18 and 60 inclusively."
	KnownCustomerNotFound                   = "Give customer do not exist."
	KnownMessageEditCustomerConflict        = "Given customer already edited, please load last data."
)

// This is custom error code
type Code string

// Some known codes
const (
	ServerInternal   Code = "ServerInternal"
	InvalidData      Code = "InvalidData"
	OverwriteData    Code = "OverwriteData"
	BadRequest       Code = "BadRequest"
	EmailTaken       Code = "EmailTaken"
	Ok               Code = "Ok"
	Created          Code = "Created"
	NotFound         Code = "NotFound"
	CustomerNotFound Code = "CustomerNotFound"
	ResourceNotFound Code = "ResourceNotFound"
)
```

For configuration I'm using <code>viper</code> that reads data from <code>.env</code>(not only) and from environment variables

```sh
package conf

import "github.com/spf13/viper"

type Config struct {
	ServerAddress string `mapstructure:"SERVER_ADDRESS"`
	DbUser        string `mapstructure:"POSTGRES_USER"`
	DbPassword    string `mapstructure:"POSTGRES_PASSWORD"`
	DbName        string `mapstructure:"POSTGRES_DB"`
	DbHost        string `mapstructure:"POSTGRES_HOST"`
}

// function panics if config cannot not be initialized
func Load() *Config {
	conf := &Config{}
	viper.AddConfigPath("./resources/")
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	// overwrite config values from envrironment variables
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		panic("cannot read config from app.env file " + err.Error())
	}
	if err := viper.Unmarshal(conf); err != nil {
		panic("error unmarsha config")
	}
	return conf
}
```

For db migration I'm using this library:
```sh
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
  
  
  	m, err := migrate.New(
		"file://resources/db/migrations",
		fmt.Sprintf("postgresql://%s:%s@%s/%s?sslmode=disable", cfg.DbUser, cfg.DbPassword, cfg.DbHost, cfg.DbName))
	if err != nil {
		return err
	}
	if err := m.Steps(1); err != nil && err != migrate.ErrNoChange {
		// just print migration error
		log.Errorf("migration error %v", err)
	}
```

Data validation was implemented with <code>go playgrond validate library and meta tags</code>:
```go
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
```
Some notes could be found in src, like this:
```go
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

```
and this:
```go

// Customer service interface, it can do below things. As data come to untrusted resources it will be better to validate
// data inside given service.
type CustomerService interface {
	// Create customer
	Create(ctx context.Context, customer *dto.CreateCustomerArguments) (err error)
	// Delete customer by id
	DeleteById(ctx context.Context, customerId int) (err error)
	// update customer(arguments includes hash, which can handle properly overriding values)
	Update(ctx context.Context, args *dto.UpdateCustomerArguments) (err error)
	// query customers list(sorting by customer fields + search on firstName and lastName)
	QueryList(ctx context.Context, args *dto.ListCustomersArguments) (result *dto.ListCustomersResult, err error)
	// get detailed information by customer id(including hash)
	GetById(ctx context.Context, customerId int) (customer *models.Customer, err error)
}
```


