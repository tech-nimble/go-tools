package errors

type errorCode interface {
	SetCode(code int)
	GetCode() int
}

func AddErrorCode(err error, code int) error {
	if err == nil {
		return nil
	}

	if extendedErr, ok := err.(errorCode); ok {
		extendedErr.SetCode(code)

		return err
	}

	return &extendedError{id: generateID(), errorType: Domain, err: err, code: code}
}

func GetErrorCode(err error) int {
	if extendedErr, ok := err.(errorCode); ok {
		return extendedErr.GetCode()
	}

	return 0
}
