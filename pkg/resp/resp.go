package resp

import (
	"encoding/json"
	"net/http"
	"sync"
	"text/template"

	"github.com/abdybaevae/customers-app/pkg/codes"
	"github.com/abdybaevae/customers-app/pkg/utils"
	"github.com/sirupsen/logrus"
)

// Some known error
// With this approach we can handle multi language error handling(depending or user lang headers)
// but let's focus on en:)
const (
	KnownMessageSomethingWrongHappened      = "Something went wrong, please try later."
	KnownMessageInvalidPageProvided         = "Invalid page provided, page must be positive integer."
	KnownMessageGivenEmailBusyUseAnotherOne = "Provided email address already in use, please provide another one."
	KnownMessageNotFoundPage                = "Given page doesn't exist."
	KnownMessageBadRequest                  = "Wrong request."
	KnownMessageCustomerCreated             = "Customer was successfully created."
	KnownMessageCustomerEdited              = "Customer was successfully edited."
	KnownMessageCustomerDeleted             = "Customer was successfully deleted."
	KnownMessageCustomerInvalidBirthDate    = "Customer birthdate must be of format yyyy-MM-dd"
)

// It's simple interface for handling errors and success answers(when resource are creating or update and so on)

type ResponseFactory interface {
	// respond to client with appropriate message(success or not depending on code mapping status)
	CodeMessage(rw http.ResponseWriter, code codes.Code, message string)
	// try reflect error to known error code otherwise response with server internal error.
	Error(rw http.ResponseWriter, err error)
}
type jsonResponseFactoryImpl struct {
	log *logrus.Logger
}
type htmlResponseFactoryImpl struct {
	log       *logrus.Logger
	templates *template.Template
}

var jsonRsOnce sync.Once
var jsonRsFactoryInstance ResponseFactory

// Singleton factory to respond with json format
func GetJsonResponseFactory() ResponseFactory {
	jsonRsOnce.Do(func() {
		jsonRsFactoryInstance = &jsonResponseFactoryImpl{
			log: logrus.New(),
		}
	})
	return jsonRsFactoryInstance
}

var htmlRsOnce sync.Once
var htmlRsFactoryInstance ResponseFactory

// singleton factory for html/template response
func GetHtmlResponseFactory() ResponseFactory {
	htmlRsOnce.Do(func() {
		htmlRsFactoryInstance = &htmlResponseFactoryImpl{
			log:       logrus.New(),
			templates: utils.LoadTemplates(),
		}
	})
	return htmlRsFactoryInstance
}

// Default Message template data
type rsMessage struct {
	Message   string `json:"message"`
	Code      string `json:"code"`
	IsSuccess bool   `json:"isSuccess"`
}

func (o *jsonResponseFactoryImpl) CodeMessage(rw http.ResponseWriter, code codes.Code, message string) {
	rw.Header().Set("Content-Type", "application/json")
	res := &rsMessage{
		Code:    string(code),
		Message: message,
	}
	rw.WriteHeader(codes.StatusCode(code))
	json.NewEncoder(rw).Encode(res)
}
func (o *jsonResponseFactoryImpl) Error(rw http.ResponseWriter, err error) {
	if errCode, ok := err.(codes.ErrorCode); ok {
		o.CodeMessage(rw, errCode.Code(), errCode.Message())
	} else {
		// print warning message for unhandled message
		o.log.Warnf("unhandled exception %v", err)
		o.CodeMessage(rw, codes.ServerInternal, KnownMessageSomethingWrongHappened)
	}
}
func (h *htmlResponseFactoryImpl) CodeMessage(rw http.ResponseWriter, code codes.Code, message string) {
	rw.Header().Set("Content-Type", "text/html")
	status := codes.StatusCode(code)
	rw.WriteHeader(status)
	h.templates.ExecuteTemplate(rw, "message", &rsMessage{
		Code:      string(code),
		Message:   message,
		IsSuccess: status >= 200 && status <= 299,
	})

}
func (h *htmlResponseFactoryImpl) Error(rw http.ResponseWriter, err error) {
	if errCode, ok := err.(codes.ErrorCode); ok {
		h.CodeMessage(rw, errCode.Code(), errCode.Message())
	} else {
		h.log.Warnf("unhandled exception %v", err)
		h.CodeMessage(rw, codes.ServerInternal, KnownMessageSomethingWrongHappened)
	}
}
