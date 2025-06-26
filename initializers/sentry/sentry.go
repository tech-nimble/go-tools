package sentry

import (
	"strconv"

	zlogsentry "github.com/archdx/zerolog-sentry"
	"github.com/getsentry/sentry-go"
	"github.com/gobuffalo/envy"
	"github.com/rs/zerolog/log"
)

const productionEnv = "production"

type AppVersion string

func InitializeSentry(v AppVersion) {
	debug, err := strconv.ParseBool(envy.Get("APP_DEBUG", "false"))
	if err != nil {
		debug = false
	}

	if err := sentry.Init(sentry.ClientOptions{
		Dsn:         envy.Get("SENTRY_DSN", ""),
		Environment: envy.Get("APP_ENV", productionEnv),
		Release:     string(v),
		Debug:       debug,
	}); err != nil {
		log.Info().Err(err).Msg("sentry initialization failed")
	}
}

func InitializeZerologSentry(v AppVersion) (*zlogsentry.Writer, error) {
	debug, err := strconv.ParseBool(envy.Get("APP_DEBUG", "false"))
	if err != nil {
		debug = false
	}

	options := []zlogsentry.WriterOption{
		zlogsentry.WithEnvironment(envy.Get("APP_ENV", productionEnv)),
		zlogsentry.WithRelease(string(v)),
	}

	if debug {
		options = append(options, zlogsentry.WithDebug())
	}

	return zlogsentry.New(envy.Get("SENTRY_DSN", ""), options...)
}
