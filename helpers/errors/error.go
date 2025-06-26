package errors

import (
	"github.com/gofrs/uuid"
)

const (
	NoType = ErrType(iota)
	Runtime
	Domain
	Debug
)

type ErrType uint

func (t ErrType) New(msg string) error {
	return &extendedError{
		id:        generateID(),
		msg:       msg,
		errorType: t,
	}
}

func (t ErrType) Newf(msg string, args ...any) error {
	return &extendedError{
		id:        generateID(),
		msg:       msg,
		args:      args,
		errorType: t,
	}
}

func (t ErrType) Wrap(err error, msg string) error {
	return t.Wrapf(err, msg)
}

func (t ErrType) Wrapf(err error, msg string, args ...any) error {
	return &extendedError{
		id:        generateID(),
		msg:       msg,
		args:      args,
		err:       err,
		errorType: t,
	}
}

func New(msg string) error {
	return NoType.New(msg)
}

func Newf(msg string, args ...any) error {
	return NoType.Newf(msg, args...)
}

func Wrap(err error, msg string) error {
	return NoType.Wrap(err, msg)
}

func Wrapf(err error, msg string, args ...any) error {
	return NoType.Wrapf(err, msg, args...)
}

func GetType(err error) ErrType {
	if e, ok := err.(interface{ GetType() ErrType }); ok {
		return e.GetType()
	}

	if extendedErr, ok := err.(*extendedError); ok {
		return extendedErr.errorType
	}

	return NoType
}

func IsDomainError(err error) bool {
	return GetType(err) == Domain
}

func IsRuntimeError(err error) bool {
	return GetType(err) == Runtime
}

func IsDebugError(err error) bool {
	return GetType(err) == Debug
}

func generateID() string {
	id, _ := uuid.NewGen().NewV4()

	return id.String()
}
