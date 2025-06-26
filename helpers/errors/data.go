package errors

import (
	"errors"
)

type errorDataReaderWriter interface {
	SetData(key string, data any)
	GetData(key string) (any, bool)
}

func AddErrorData(err error, field string, data any) error {
	if err == nil {
		return nil
	}

	if extErr, ok := err.(errorDataReaderWriter); ok {
		extErr.SetData(field, data)
	}

	return err
}

func GetErrorData(err error, key string) any {
	if extendedErr, ok := err.(errorDataReaderWriter); ok {
		if v, ok := extendedErr.GetData(key); ok {
			return v
		}
	}

	if err := errors.Unwrap(err); err != nil {
		return GetErrorData(err, key)
	}

	return nil
}
