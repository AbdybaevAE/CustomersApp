package codes

import (
	"errors"
	"net/http"
)

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

// and reverse mapping to http status int
var codeToStatus = map[Code]int{
	InvalidData:      http.StatusBadRequest,
	OverwriteData:    http.StatusConflict,
	ServerInternal:   http.StatusInternalServerError,
	BadRequest:       http.StatusBadRequest,
	EmailTaken:       http.StatusBadRequest,
	Ok:               http.StatusOK,
	NotFound:         http.StatusNotFound,
	Created:          http.StatusCreated,
	CustomerNotFound: http.StatusNotFound,
}

func StatusCode(code Code) int {
	return codeToStatus[code]
}

// This is general error interface which includes typed message Code that understandable by service consumers(obvioulsy programs, machines)
// and message string(human readable and understandable error message)
// Its' expected that ErrorCode are used as error values which can handle both customers-app server and his consumers(frontend app, mobile app, another third party libraries)
type ErrorCode interface {
	error
	Code() Code
	Message() string
}

// Simple implementation of ErrorCode
type errorCodeImpl struct {
	CodeValue    Code
	MessageValue string
}

func (e *errorCodeImpl) Error() string {
	return string(e.CodeValue) + " " + e.MessageValue
}
func NewErr(code Code, message string) ErrorCode {
	return &errorCodeImpl{
		CodeValue:    code,
		MessageValue: message,
	}
}
func (e *errorCodeImpl) Code() Code {
	return e.CodeValue
}
func (e *errorCodeImpl) Message() string {
	return e.MessageValue
}

// Custom sql repository errors

var NoRowsModified = errors.New("No rows modified")
var UniqueConstraintViolation = errors.New("Unique constraint vialation")
var BadSearchCriteria = errors.New("Bad search criteria")
