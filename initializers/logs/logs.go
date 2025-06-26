package logs

import (
	"io"
	"os"
	"strconv"

	zlogsentry "github.com/archdx/zerolog-sentry"
	"github.com/gobuffalo/envy"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	LogLevelEnv       = "LOG_LEVEL"
	EnableJSONLogsEnv = "ENABLE_JSON_LOGS"
	DefaultLogLevel   = 0
)

func InitializeLogs(sentry *zlogsentry.Writer) zerolog.Logger {
	logLevel, err := strconv.Atoi(envy.Get(LogLevelEnv, "0"))
	if err != nil {
		logLevel = DefaultLogLevel
	}

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.Level(logLevel))

	if sentry != nil {
		log.Logger = log.Output(io.MultiWriter(sentry, os.Stderr))
	}

	// json or human readable output
	if envy.Get(EnableJSONLogsEnv, "true") == "false" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	// add filepath+row num to log(app/cmd/serve.go:37)
	log.Logger = log.With().Caller().Logger()

	return log.Logger
}
