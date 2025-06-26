package errors

const ValidationErrCode = 101

func DomainErrWithDefaultCode(msg string) error {
	return DomainErrWithCode(msg, ValidationErrCode)
}
func DomainErrWithDefaultCodeT(msg, msgID string) error {
	return DomainErrWithCodeT(msg, msgID, ValidationErrCode)
}

func DomainErrWithCode(msg string, code int) error {
	err := Domain.New(msg)

	return AddErrorCode(err, code)
}

func DomainErrWithCodeT(msg, msgID string, code int) error {
	err := Domain.NewT(msg, msgID)

	return AddErrorCode(err, code)
}

func NewfDomainErrWidthCode(msg string, code int, args ...any) error {
	err := Domain.Newf(msg, args...)

	return AddErrorCode(err, code)
}

func NewTfDomainErrWidthCode(msg, msgID string, code int, namedArgs map[string]any, pluralCount any) error {
	err := Domain.NewTf(msg, msgID, namedArgs, pluralCount)

	return AddErrorCode(err, code)
}
