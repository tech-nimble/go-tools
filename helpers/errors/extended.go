package errors

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
)

type extendedError struct {
	id        string
	msg       string
	args      []any
	err       error
	errorType ErrType
	context   context.Context
	data      map[string]any
	code      int
}

func (err extendedError) Error() string {
	msg := err.msg

	if msg == "" && err.err != nil {
		return err.err.Error()
	}

	if len(err.args) > 0 {
		msg = fmt.Sprintf(err.msg, err.args...)
	}

	switch err.errorType {
	case NoType, Debug, Runtime:
		if err.err == nil {
			return msg
		}

		return fmt.Sprintf("%s: %s", msg, err.err.Error())
	default:
		return msg
	}
}

func (err extendedError) Is(target error) bool {
	extTarget, ok := target.(*extendedError)
	if !ok {
		return false
	}

	return err.id == extTarget.id
}

func (err extendedError) Unwrap() error {
	return err.err
}

func (err *extendedError) SetCode(code int) {
	err.code = code
}

func (err *extendedError) GetCode() int {
	return err.code
}

func (err *extendedError) SetData(key string, data any) {
	if err.data == nil {
		err.data = map[string]any{}
	}

	err.data[key] = data
}

func (err *extendedError) GetData(key string) (any, bool) {
	v, ok := err.data[key]

	if ok {
		return v, true
	}

	return nil, false
}

func (err *extendedError) SetContext(ctx context.Context) {
	err.context = ctx
}

func (err *extendedError) GetContext() context.Context {
	return err.context
}

func (err *extendedError) GetType() ErrType {
	return err.errorType
}

func (err extendedError) Cause() error {
	return err.err
}

func (err extendedError) StackTrace() errors.StackTrace {
	type stackTracer interface {
		StackTrace() errors.StackTrace
	}

	sterr, ok := err.err.(stackTracer)
	if !ok {
		return nil
	}

	return sterr.StackTrace()
}
