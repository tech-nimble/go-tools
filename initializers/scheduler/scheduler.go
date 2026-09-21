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

// Initialize builds a scheduler that understands seconds in a cron spec and never starts a second
// copy of a job that is still running: a task ticking every few seconds would otherwise overlap
// itself and process the same rows twice.
func Initialize() *cron.Cron {
	return cron.New(
		cron.WithSeconds(),
		cron.WithChain(cron.SkipIfStillRunning(cronLogger{})),
	)
}

// EverySpec turns an interval into a cron spec. A non-positive interval means a broken
// configuration, and it must not turn into a job running every second.
func EverySpec(interval time.Duration) string {
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
