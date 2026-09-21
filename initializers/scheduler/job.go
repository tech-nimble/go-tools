// SPDX-FileCopyrightText: 2025 Nimble Tech
// SPDX-License-Identifier: MIT

package scheduler

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog/log"
)

// DefaultJobTimeout caps a job that does not set a limit of its own.
const DefaultJobTimeout = 5 * time.Minute

// ErrJobNotFound is returned when no job carries the requested name.
var ErrJobNotFound = errors.New("job not found")

// Job describes a scheduled task: when it runs, how long it may take and what it does.
//
// FailuresBeforeAlarm is how many failures in a row a frequent job survives at warning level: a
// job ticking every few seconds recovers from a short database or network outage by itself, and
// one timeout is not worth an alert. It stays 1 unless the job says otherwise, so a rare job
// reports the very first failure — three silent nights of a nightly job is not a trade worth making.
type Job struct {
	Name                string
	Spec                string
	Timeout             time.Duration
	Quiet               bool
	FailuresBeforeAlarm int
	Run                 func(ctx context.Context) error
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
		limit = 1
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

// Names lists the job names, in the order the service declared them.
func Names(jobs ...Job) []string {
	names := make([]string, 0, len(jobs))
	for _, job := range jobs {
		names = append(names, job.Name)
	}

	return names
}

// RunOnce runs a single job by name, outside of any schedule, and reports what it returned. It
// exists so a service can try a job from the command line instead of waiting for its next tick.
func RunOnce(ctx context.Context, name string, jobs ...Job) error {
	for _, job := range jobs {
		if job.Name != name {
			continue
		}

		jobCtx, cancel := context.WithTimeout(ctx, jobTimeout(job))
		defer cancel()

		log.Info().Str("op", job.Name).Str("phase", "start").Msg("cron: job started once")

		if err := job.Run(jobCtx); err != nil {
			log.Error().Err(err).
				Str("op", job.Name).
				Str("phase", "end").
				Str("result", "failed").
				Msg("cron: job failed")

			return fmt.Errorf("%w, while running job %s", err, job.Name)
		}

		log.Info().Str("op", job.Name).Str("phase", "end").Str("result", "success").Msg("cron: job finished")

		return nil
	}

	return fmt.Errorf("%w: %s", ErrJobNotFound, name)
}

func jobTimeout(job Job) time.Duration {
	if job.Timeout <= 0 {
		return DefaultJobTimeout
	}

	return job.Timeout
}

func register(ctx context.Context, cronScheduler *cron.Cron, job Job) {
	tracker := newFailureTracker(job.FailuresBeforeAlarm)

	if _, err := cronScheduler.AddFunc(job.Spec, func() {
		jobCtx, cancel := context.WithTimeout(ctx, jobTimeout(job))
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
