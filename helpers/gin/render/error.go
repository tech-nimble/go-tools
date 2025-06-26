package render

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gobuffalo/envy"
	"github.com/google/jsonapi"
	"github.com/google/uuid"
	"github.com/tech-nimble/go-tools/helpers/errors"
)

const (
	UnknownErrCode = 400
)

type JSONAPIErrors struct {
	Errors     []*jsonapi.ErrorObject
	httpStatus int
	isDebug    bool
}

func (e *JSONAPIErrors) Render(w http.ResponseWriter) error {
	e.WriteContentType(w)

	return json.NewEncoder(w).Encode(jsonapi.ErrorsPayload{Errors: e.Errors})
}

func (e *JSONAPIErrors) WriteContentType(w http.ResponseWriter) {
	header := w.Header()
	if val := header["Content-Type"]; len(val) == 0 {
		header["Content-Type"] = []string{jsonapi.MediaType}
	}
}

func NewJSONAPIErrors() *JSONAPIErrors {
	isDebug := envy.Get("APP_DEBUG", "false") == "true"

	return &JSONAPIErrors{
		isDebug: isDebug,
	}
}

func (e *JSONAPIErrors) AddErrors(errs []error) *JSONAPIErrors {
	for _, err := range errs {
		e.AddError(err)
	}

	return e
}

func (e *JSONAPIErrors) AddError(err error) *JSONAPIErrors {
	httpStatus := e.getStatus(err)

	e.Errors = append(e.Errors, &jsonapi.ErrorObject{
		ID:     e.getId(),
		Title:  e.getTitle(err),
		Detail: e.getDetail(err),
		Status: strconv.Itoa(httpStatus),
		Code:   strconv.Itoa(e.getCode(err)),
	})

	if e.httpStatus < 500 {
		e.httpStatus = httpStatus
	}

	return e
}

func (e *JSONAPIErrors) getTitle(err error) string {
	errType := errors.GetType(err)
	if errType == errors.Domain || e.isDebug {
		return err.Error()
	}

	return errors.ServerErrTitle
}

func (e *JSONAPIErrors) getDetail(err error) string {
	errType := errors.GetType(err)
	if errType == errors.Domain || e.isDebug {
		return err.Error()
	}

	return errors.ServerErrTitle
}

func (e *JSONAPIErrors) getId() string {
	return uuid.New().String()
}

func (e *JSONAPIErrors) getStatus(err error) int {
	httpStatus := http.StatusInternalServerError

	errType := errors.GetType(err)
	if errType == errors.Domain {
		httpStatus = http.StatusBadRequest
	}

	return httpStatus
}

func (e *JSONAPIErrors) getCode(err error) int {
	code := http.StatusInternalServerError

	errType := errors.GetType(err)
	if errType == errors.Domain {
		if code = errors.GetErrorCode(err); code == 0 {
			code = UnknownErrCode
		}
	}

	return code
}

func (e *JSONAPIErrors) GetHttpStatus() int {
	if e.httpStatus == 0 {
		return http.StatusInternalServerError
	}

	return e.httpStatus
}
