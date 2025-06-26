package errors

import (
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
	sentryHelper "github.com/tech-nimble/go-tools/helpers/sentry"
)

const ServerErrTitle = "Ошибка сервера, попробуйте позже"

var handler ErrHandler

type ErrHandler struct {
	isDebug bool
	logger  zerolog.Logger
	sentry  *sentry.Client
}

func InitHandler(debug bool, logger zerolog.Logger, client *sentry.Client) ErrHandler {
	handler = ErrHandler{
		isDebug: debug,
		logger:  logger,
		sentry:  client,
	}

	return handler
}

func Handle(err error) error {
	return handler.Handle(err)
}

func Log(err error) {
	handler.Log(err)
}

func (h ErrHandler) Handle(err error) error {
	if err == nil {
		return nil
	}

	errType := GetType(err)

	switch errType {
	case NoType, Runtime:
		h.Log(err)

		err = Wrap(err, ServerErrTitle)
	case Domain:
		h.Log(err)
	case Debug:
		h.Log(err)

		err = nil
	}

	return err
}

func (h ErrHandler) Log(err error) {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	errType := GetType(err)

	switch errType {
	case NoType, Runtime:
		h.logger.Error().Stack().Err(err).Send()
		h.SentryLog(err, sentry.LevelError)
	case Domain:
		h.logger.Info().Err(err).Send()
	case Debug:
		h.logger.Debug().Err(err).Send()
	}
}

func (h ErrHandler) SentryLog(err error, lvl sentry.Level) {
	if h.sentry == nil {
		return
	}

	var hub *sentry.Hub
	event := sentry.Event{
		Timestamp: time.Now().UTC(),
		Level:     lvl,
	}

	ctx := GetContext(err)
	if ctx != nil {
		hub = sentryHelper.GetHubFromContext(ctx)
	}

	event.Message = err.Error()
	event.Exception = append(event.Exception, sentry.Exception{
		Value:      err.Error(),
		Stacktrace: sentry.ExtractStacktrace(err),
	})

	if hub != nil {
		hub.CaptureEvent(&event)
	} else {
		h.sentry.CaptureEvent(&event, nil, nil)
	}
}
