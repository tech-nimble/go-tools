// SPDX-FileCopyrightText: 2025 Nimble Tech
// SPDX-License-Identifier: MIT

package scheduler

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"
)

func TestRegister_RunsWithOwnTimeout(t *testing.T) {
	t.Parallel()

	const timeout = time.Minute

	var (
		jobCtxErr   error
		jobDeadline time.Time
		hasDeadline bool
	)

	cronScheduler := cron.New()
	Register(context.Background(), cronScheduler, Job{
		Name:    "test.job",
		Spec:    "0 3 * * *",
		Timeout: timeout,
		Run: func(ctx context.Context) error {
			jobCtxErr = ctx.Err()
			jobDeadline, hasDeadline = ctx.Deadline()

			return nil
		},
	})

	entries := cronScheduler.Entries()
	require.Len(t, entries, 1)
	entries[0].Job.Run()

	require.NoError(t, jobCtxErr)
	require.True(t, hasDeadline)
	require.WithinDuration(t, time.Now().Add(timeout), jobDeadline, time.Minute)
}

func TestRegister_ZeroTimeoutFallsBackToDefault(t *testing.T) {
	t.Parallel()

	var (
		jobDeadline time.Time
		hasDeadline bool
	)

	cronScheduler := cron.New()
	Register(context.Background(), cronScheduler, Job{
		Name: "test.job",
		Spec: "0 3 * * *",
		Run: func(ctx context.Context) error {
			jobDeadline, hasDeadline = ctx.Deadline()

			return nil
		},
	})

	entries := cronScheduler.Entries()
	require.Len(t, entries, 1)
	entries[0].Job.Run()

	require.True(t, hasDeadline)
	require.WithinDuration(t, time.Now().Add(DefaultJobTimeout), jobDeadline, time.Minute)
}

func TestRegister_StopsWithApplicationContext(t *testing.T) {
	t.Parallel()

	appCtx, stopJobs := context.WithCancel(context.Background())

	var jobCtxErr error

	cronScheduler := cron.New()
	Register(appCtx, cronScheduler, Job{
		Name:    "test.job",
		Spec:    "0 3 * * *",
		Timeout: time.Minute,
		Run: func(ctx context.Context) error {
			jobCtxErr = ctx.Err()

			return nil
		},
	})
	stopJobs()

	entries := cronScheduler.Entries()
	require.Len(t, entries, 1)
	entries[0].Job.Run()

	require.ErrorIs(t, jobCtxErr, context.Canceled)
}

func TestRegister_InvalidSpecDoesNotRegister(t *testing.T) {
	t.Parallel()

	cronScheduler := cron.New()
	Register(context.Background(), cronScheduler, Job{
		Name:    "test.job",
		Spec:    "not a spec",
		Timeout: time.Minute,
		Run: func(context.Context) error {
			return errors.New("must not run")
		},
	})

	require.Empty(t, cronScheduler.Entries())
}

func TestFailureTracker_Failed(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		limit         int
		errs          []error
		wantAlarmedAt []bool
	}{
		{
			name:          "timeouts stay quiet until limit",
			limit:         3,
			errs:          []error{context.DeadlineExceeded, context.DeadlineExceeded, context.DeadlineExceeded},
			wantAlarmedAt: []bool{false, false, true},
		},
		{
			name:          "canceled context is transient too",
			limit:         2,
			errs:          []error{context.Canceled, context.Canceled},
			wantAlarmedAt: []bool{false, true},
		},
		{
			name:          "wrapped timeout is still transient",
			limit:         2,
			errs:          []error{fmt.Errorf("%w, while expiring offers", context.DeadlineExceeded)},
			wantAlarmedAt: []bool{false},
		},
		{
			name:          "other errors alarm at once",
			limit:         3,
			errs:          []error{errors.New("broken query")},
			wantAlarmedAt: []bool{true},
		},
		{
			name:          "zero limit alarms on the first failure",
			limit:         0,
			errs:          []error{context.DeadlineExceeded},
			wantAlarmedAt: []bool{true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tracker := newFailureTracker(tt.limit)

			for i, err := range tt.errs {
				inRow, alarming := tracker.failed(err)

				require.Equal(t, i+1, inRow)
				require.Equal(t, tt.wantAlarmedAt[i], alarming)
			}
		})
	}
}

func TestFailureTracker_SuccessResetsRow(t *testing.T) {
	t.Parallel()

	tracker := newFailureTracker(2)

	inRow, alarming := tracker.failed(context.DeadlineExceeded)
	require.Equal(t, 1, inRow)
	require.False(t, alarming)

	tracker.succeeded()

	inRow, alarming = tracker.failed(context.DeadlineExceeded)
	require.Equal(t, 1, inRow)
	require.False(t, alarming)
}

func TestRegister_SingleTimeoutIsNotReportedAsError(t *testing.T) {
	var buf bytes.Buffer

	originalLogger := log.Logger
	log.Logger = zerolog.New(&buf)

	t.Cleanup(func() { log.Logger = originalLogger })

	cronScheduler := cron.New()
	Register(context.Background(), cronScheduler, Job{
		Name:                "test.job",
		Spec:                "0 3 * * *",
		Timeout:             time.Minute,
		Quiet:               true,
		FailuresBeforeAlarm: 2,
		Run: func(context.Context) error {
			return fmt.Errorf("%w, while expiring offers", context.DeadlineExceeded)
		},
	})

	entries := cronScheduler.Entries()
	require.Len(t, entries, 1)

	buf.Reset()
	entries[0].Job.Run()
	require.Contains(t, buf.String(), `"level":"warn"`)
	require.Contains(t, buf.String(), `"failuresInRow":1`)

	buf.Reset()
	entries[0].Job.Run()
	require.Contains(t, buf.String(), `"level":"error"`)
	require.Contains(t, buf.String(), `"failuresInRow":2`)
}

func TestRunOnce(t *testing.T) {
	t.Parallel()

	jobErr := errors.New("broken query")

	tests := []struct {
		name    string
		run     string
		jobRun  func(context.Context) error
		wantErr error
		wantRan bool
	}{
		{
			name:    "runs the named job",
			run:     "test.job",
			jobRun:  func(context.Context) error { return nil },
			wantRan: true,
		},
		{
			name:    "reports what the job returned",
			run:     "test.job",
			jobRun:  func(context.Context) error { return jobErr },
			wantErr: jobErr,
			wantRan: true,
		},
		{
			name:    "unknown name is an error",
			run:     "test.missing",
			jobRun:  func(context.Context) error { return nil },
			wantErr: ErrJobNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var ran bool

			err := RunOnce(context.Background(), tt.run, Job{
				Name:    "test.job",
				Spec:    "0 3 * * *",
				Timeout: time.Minute,
				Run: func(ctx context.Context) error {
					ran = true

					return tt.jobRun(ctx)
				},
			})

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tt.wantRan, ran)
		})
	}
}

func TestRunOnce_UsesJobTimeout(t *testing.T) {
	t.Parallel()

	var (
		deadline    time.Time
		hasDeadline bool
	)

	err := RunOnce(context.Background(), "test.job", Job{
		Name:    "test.job",
		Spec:    "0 3 * * *",
		Timeout: time.Minute,
		Run: func(ctx context.Context) error {
			deadline, hasDeadline = ctx.Deadline()

			return nil
		},
	})

	require.NoError(t, err)
	require.True(t, hasDeadline)
	require.WithinDuration(t, time.Now().Add(time.Minute), deadline, time.Minute)
}

func TestNames(t *testing.T) {
	t.Parallel()

	names := Names(
		Job{Name: "first.job"},
		Job{Name: "second.job"},
	)

	require.Equal(t, []string{"first.job", "second.job"}, names)
}
