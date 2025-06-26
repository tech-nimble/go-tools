package errors

import (
	"context"
	"errors"
)

type errorContext interface {
	SetContext(ctx context.Context)
	GetContext() context.Context
}

func AddContext(err error, ctx context.Context) error {
	if err == nil {
		return nil
	}

	if extendedErr, ok := err.(errorContext); ok {
		extendedErr.SetContext(ctx)

		return err
	}

	return &extendedError{id: generateID(), errorType: NoType, err: err, context: ctx}
}

func GetContext(err error) context.Context {
	if extendedErr, ok := err.(errorContext); ok && extendedErr.GetContext() != nil {
		return extendedErr.GetContext()
	}

	if err := errors.Unwrap(err); err != nil {
		return GetContext(err)
	}

	return nil
}
