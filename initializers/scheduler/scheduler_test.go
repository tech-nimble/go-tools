// SPDX-FileCopyrightText: 2025 Nimble Tech
// SPDX-License-Identifier: MIT

package scheduler

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestInitialize_SupportsSecondsAndSkipsOverlap(t *testing.T) {
	t.Parallel()

	var (
		running atomic.Bool
		overlap atomic.Bool
		runs    atomic.Int32
	)

	cronScheduler := Initialize()
	Register(context.Background(), cronScheduler, Job{
		Name:    "test.job",
		Spec:    "@every 1s",
		Timeout: time.Minute,
		Quiet:   true,
		Run: func(context.Context) error {
			if !running.CompareAndSwap(false, true) {
				overlap.Store(true)
			}

			runs.Add(1)
			time.Sleep(2500 * time.Millisecond)
			running.Store(false)

			return nil
		},
	})

	require.Len(t, cronScheduler.Entries(), 1)

	cronScheduler.Start()
	time.Sleep(4 * time.Second)
	<-cronScheduler.Stop().Done()

	require.False(t, overlap.Load())
	require.Less(t, runs.Load(), int32(3))
}

func TestEverySpec(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		interval time.Duration
		want     string
	}{
		{name: "positive interval", interval: 30 * time.Second, want: "@every 30s"},
		{name: "zero falls back", interval: 0, want: "@every 5s"},
		{name: "negative falls back", interval: -time.Second, want: "@every 5s"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, tt.want, EverySpec(tt.interval))
		})
	}
}
