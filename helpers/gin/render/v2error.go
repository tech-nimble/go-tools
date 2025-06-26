package render

import "github.com/tech-nimble/go-tools/helpers/errors"

const (
	ValidationErrCode     = 101
	InternalServerErrCode = 500
)

func NewResponseWithError(err error) Response {
	msg := errors.ServerErrTitle
	code := InternalServerErrCode
	errType := errors.GetType(err)

	if errType == errors.Domain {
		msg = err.Error()

		if code = errors.GetErrorCode(err); code == 0 {
			code = UnknownErrCode
		}
	}

	return NewResponse(nil, code, "error", msg)
}
