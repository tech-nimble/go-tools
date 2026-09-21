// SPDX-FileCopyrightText: 2025 Nimble Tech
// SPDX-License-Identifier: MIT

package scheduler

import (
	"time"

	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog/log"
)

// DefaultTickInterval is used when a configured interval is missing or malformed.
const DefaultTickInterval = 5 * time.Second

type cronLogger struct{}

// Initialize builds a scheduler that never starts a second copy of a job that is still running:
// a task ticking every few seconds would otherwise overlap itself and process the same rows twice.
//
// The seconds field is optional, so both "15 3 * * *" and "0 30 3 * * *" are valid specs and a
// service is free to schedule a job below the minute without rewriting the rest of its specs.
func Initialize() *cron.Cron {
	parser := cron.NewParser(
		cron.SecondOptional | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor,
	)

	return cron.New(
		cron.WithParser(parser),
		cron.WithChain(cron.SkipIfStillRunning(cronLogger{})),
	)
}

// EverySpec turns an interval into a cron spec. A non-positive interval means a broken
// configuration, and it must not turn into a job running every second, so the caller says what
// to fall back to; a non-positive fallback lands on DefaultTickInterval.
func EverySpec(interval, fallback time.Duration) string {
	if interval <= 0 {
		interval = fallback
	}

	if interval <= 0 {
		interval = DefaultTickInterval
	}

	return "@every " + interval.String()
}

func (cronLogger) Info(msg string, keysAndValues ...any) {
	log.Debug().Fields(keysAndValues).Msg("cron: " + msg)
}

func (cronLogger) Error(err error, msg string, keysAndValues ...any) {
	log.Error().Err(err).Fields(keysAndValues).Msg("cron: " + msg)
}
