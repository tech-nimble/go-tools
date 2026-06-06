// SPDX-FileCopyrightText: 2025 Nimble Tech
// SPDX-License-Identifier: MIT

package errors

import (
	"context"
	"errors"
)

type errorContext interface {
	SetContext(ctx context.Context)
	GetContext() context.Context
}

func AddContext(ctx context.Context, err error) error {
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
