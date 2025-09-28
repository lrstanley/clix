// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

package clix

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"
)

// Runner is a generic runnable entity. See also [RunnerFunc].
type Runner interface {
	Invoke(ctx context.Context) error
}

var _ Runner = (*RunnerFunc)(nil)

// RunnerFunc is a function that can be used to create a [Runner].
type RunnerFunc func(ctx context.Context) error

func (r RunnerFunc) Invoke(ctx context.Context) error {
	return r(ctx)
}

// Run invokes all runners concurrently, and listens for any termination signals
// (SIGINT, SIGTERM, SIGQUIT, etc).
//
// If any runners return an error, all runners will terminate (assuming they listen
// to the provided context), and the first known error will be returned.
func Run(ctx context.Context, runners ...Runner) error {
	if len(runners) == 0 {
		panic("no runners provided")
	}

	ctx, cancel := signal.NotifyContext(
		ctx,
		os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	defer cancel()

	var g *errgroup.Group
	g, ctx = errgroup.WithContext(ctx)

	for _, runner := range runners {
		g.Go(func() error {
			return runner.Invoke(ctx)
		})
	}

	return g.Wait()
}

var _ Runner = (*Cron)(nil)

type Cron struct {
	name        string
	interval    time.Duration
	immediate   bool
	exitOnError bool
	runner      Runner
	logger      *slog.Logger
}

// NewCron creates a new cron runner with the provided name and underlying runner.
// The cron runner will run the runner at the provided interval, and will exit
// on error if the exitOnError flag is set. The default interval is 5 minutes,
// which can be changed with [Cron.WithInterval].
func NewCron(name string, runner Runner) *Cron {
	return &Cron{
		name:     name,
		runner:   runner,
		interval: 5 * time.Minute,
		logger:   slog.Default(),
	}
}

// WithInterval sets the interval at which the cron runner will run the underlying
// runner. Defaults to 5 minutes.
func (c *Cron) WithInterval(interval time.Duration) *Cron {
	if interval > 5*time.Millisecond {
		c.interval = interval
	}
	return c
}

// WithImmediate sets whether the cron runner should run the underlying runner
// immediately upon creation. This defaults to false. If true, the runner will
// also exit on error if the initial immediate run fails.
func (c *Cron) WithImmediate(enabled bool) *Cron {
	c.immediate = enabled
	return c
}

// WithExitOnError sets whether the cron runner should exit on error. This defaults
// to false. If true, the runner will exit if the underlying runner returns an error.
func (c *Cron) WithExitOnError(enabled bool) *Cron {
	c.exitOnError = enabled
	return c
}

// WithLogger sets the logger for the cron runner. This defaults to the default
// logger.
func (c *Cron) WithLogger(logger *slog.Logger) *Cron {
	if logger != nil {
		c.logger = logger
	}
	return c
}

// Invoke runs the cron runner. This is typically not called directly, but rather
// via [Run] or [RunContext].
func (c *Cron) Invoke(ctx context.Context) error {
	l := c.logger.With(
		"cron", c.name,
		"interval", c.interval,
		"exit_on_error", c.exitOnError,
	)

	var lastRun time.Time

	if c.immediate {
		lastRun = time.Now()
		l.InfoContext(ctx, "invoking cron")
		if err := c.runner.Invoke(ctx); err != nil {
			l.ErrorContext(
				ctx,
				"cron failed",
				"error", err,
				"duration", time.Since(lastRun),
			)
			return err
		}
		l.InfoContext(
			ctx,
			"cron complete",
			"duration", time.Since(lastRun),
		)
	}

	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			lastRun = time.Now()
			l.InfoContext(ctx, "invoking cron")
			if err := c.runner.Invoke(ctx); err != nil {
				l.ErrorContext(
					ctx,
					"cron failed",
					"error", err,
					"duration", time.Since(lastRun),
				)
				if c.exitOnError {
					return err
				}
			}
			l.InfoContext(
				ctx,
				"cron complete",
				"duration", time.Since(lastRun),
			)
		}
	}
}
