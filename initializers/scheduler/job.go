// SPDX-FileCopyrightText: 2025 Nimble Tech
// SPDX-License-Identifier: MIT

package scheduler

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog/log"
)

const (
	// DefaultJobTimeout caps a job that does not set a limit of its own.
	DefaultJobTimeout = 5 * time.Minute

	// DefaultTransientFailures is how many failures in a row a job survives quietly. A short
	// outage of the database or the network is cured by the next tick, so an alert is only
	// worth raising once the work has been stuck that many runs.
	DefaultTransientFailures = 3
)

// Job describes a scheduled task: when it runs, how long it may take and what it does.
type Job struct {
	Name              string
	Spec              string
	Timeout           time.Duration
	Quiet             bool
	TransientFailures int
	Run               func(ctx context.Context) error
}

type failureTracker struct {
	limit int

	mu    sync.Mutex
	inRow int
}

// Register adds jobs to the scheduler. The context belongs to the application, not to a single
// run: canceling it stops the jobs that are in flight.
func Register(ctx context.Context, cronScheduler *cron.Cron, jobs ...Job) {
	for _, job := range jobs {
		register(ctx, cronScheduler, job)
	}
}

func newFailureTracker(limit int) *failureTracker {
	if limit <= 0 {
		limit = DefaultTransientFailures
	}

	return &failureTracker{limit: limit}
}

func (t *failureTracker) failed(err error) (int, bool) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.inRow++

	transient := errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled)

	return t.inRow, !transient || t.inRow >= t.limit
}

func (t *failureTracker) succeeded() {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.inRow = 0
}

func register(ctx context.Context, cronScheduler *cron.Cron, job Job) {
	timeout := job.Timeout
	if timeout <= 0 {
		timeout = DefaultJobTimeout
	}

	tracker := newFailureTracker(job.TransientFailures)

	if _, err := cronScheduler.AddFunc(job.Spec, func() {
		jobCtx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		if !job.Quiet {
			log.Info().Str("op", job.Name).Str("phase", "start").Msg("cron: job started")
		}

		if err := job.Run(jobCtx); err != nil {
			inRow, alarming := tracker.failed(err)

			event := log.Warn()
			if alarming {
				event = log.Error()
			}

			event.Err(err).
				Str("op", job.Name).
				Str("phase", "end").
				Str("result", "failed").
				Int("failuresInRow", inRow).
				Msg("cron: job failed")

			return
		}

		tracker.succeeded()

		if !job.Quiet {
			log.Info().Str("op", job.Name).Str("phase", "end").Str("result", "success").Msg("cron: job finished")
		}
	}); err != nil {
		log.Error().Err(err).
			Str("op", job.Name).
			Str("spec", job.Spec).
			Msg("cron: failed to register job")

		return
	}

	log.Info().Str("op", job.Name).Str("spec", job.Spec).Msg("cron: job registered")
}
