package errors

import (
	"github.com/tech-nimble/go-i18n"
)

type translatedError struct {
	*extendedError
	localizer   *i18n.Localize
	msgID       string
	namedArgs   map[string]any
	pluralCount any
}

func (err translatedError) Error() string {
	msg := err.msg

	if (msg == "" && err.msgID == "") || err.localizer == nil {
		return err.extendedError.Error()
	}

	if err.msgID != "" {
		msg = err.msgID
	}

	if len(err.namedArgs) > 0 {
		return err.localizer.Tf(msg, err.namedArgs, err.pluralCount)
	}

	return err.localizer.T(msg)
}

func (err translatedError) GetType() ErrType {
	return err.errorType
}

func (t ErrType) NewT(msg, msgID string) error {
	return &translatedError{
		extendedError: &extendedError{
			id:        generateID(),
			msg:       msg,
			errorType: t,
		},
		msgID: msgID,
	}
}

func (t ErrType) NewTf(msg, msgID string, namedArgs map[string]any, pluralCount any) error {
	return &translatedError{
		extendedError: &extendedError{
			id:        generateID(),
			msg:       msg,
			errorType: t,
		},
		msgID:       msgID,
		namedArgs:   namedArgs,
		pluralCount: pluralCount,
	}
}

func (t ErrType) WrapT(err error, msg, msgID string) error {
	return &translatedError{
		extendedError: &extendedError{
			id:        generateID(),
			msg:       msg,
			err:       err,
			errorType: t,
		},
		msgID: msgID,
	}
}

func (t ErrType) WrapTf(err error, msg, msgID string, namedArgs map[string]any, pluralCount any) error {
	return &translatedError{
		extendedError: &extendedError{
			id:        generateID(),
			msg:       msg,
			err:       err,
			errorType: t,
		},
		msgID:       msgID,
		namedArgs:   namedArgs,
		pluralCount: pluralCount,
	}
}

func AddErrLocalize(err error, localize *i18n.Localize) error {
	if err == nil {
		return nil
	}

	if trErr, ok := err.(*translatedError); ok {
		translatedErr := *trErr
		translatedErr.localizer = localize

		return &translatedErr
	}

	return err
}

func AddErrLocalizeData(err error, data map[string]any) error {
	if err == nil {
		return nil
	}

	if trErr, ok := err.(*translatedError); ok {
		translatedErr := *trErr
		translatedErr.namedArgs = data

		return &translatedErr
	}

	return err
}

func AddErrLocalizePluralCount(err error, pluralCount any) error {
	if err == nil {
		return nil
	}

	if trErr, ok := err.(*translatedError); ok {
		translatedErr := *trErr
		translatedErr.pluralCount = pluralCount

		return &translatedErr
	}

	return err
}
