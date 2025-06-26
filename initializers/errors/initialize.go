package errors

import (
	"strconv"

	"github.com/getsentry/sentry-go"
	"github.com/gobuffalo/envy"
	"github.com/rs/zerolog"
	"github.com/tech-nimble/go-tools/helpers/errors"
)

func Initialize(logger zerolog.Logger, client *sentry.Client) errors.ErrHandler {
	debug, err := strconv.ParseBool(envy.Get("APP_DEBUG", "false"))
	if err != nil {
		debug = false
	}

	return errors.InitHandler(debug, logger, client)
}
