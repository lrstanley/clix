// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

// A simple example using the out-of-the-box functionality with clix, which includes
// logging, version flags, loading of .env files, etc, but also shows using the
// github.com/lrstanley/x/scheduler package to run a set of jobs along side your
// application.
//
// This example also shows how to break up your flags into different groups,
// which could be pulled from other sub-packages.
package main

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/lrstanley/clix/v2"
	"github.com/lrstanley/x/scheduler"
)

type Flags struct {
	Service *ServiceConfig `prefix:"service" envprefix:"SERVICE_" group:"Service flags"`
}

var cli = clix.NewWithDefaults[Flags]()

func main() {
	logger := cli.GetLogger()

	ctx := context.TODO()

	// This is an example of using the github.com/lrstanley/x/scheduler package
	// to run a set of jobs. It's not technically tied to clix, but it pairs well
	// with it. You can use the scheduler for cron/interval based jobs, or
	// background jobs that you want to always run (and if they error, exit
	// the application). These jobs also listen to signals like SIGINT, SIGTERM,
	// etc. and will stop all jobs and exit the application when they are received.
	err := scheduler.Run(
		ctx,
		// Interval based cron job using the [scheduler.Job] interface.
		scheduler.NewCron("fetch-svc", &fetchService{logger: logger}).
			WithInterval(30*time.Second).
			WithImmediate(true).
			WithExitOnError(true).
			WithLogger(logger),

		// Crontab-style cron job using the [scheduler.JobFunc] wrapper.
		scheduler.NewCron("test", scheduler.JobFunc(hourlyCron)).
			WithSchedule("0 * * * *"),

		// Long-running job (that doesn't get invoked at an interval, just runs
		// in the background).
		scheduler.JobFunc(longRunningJob),
	)
	if err != nil {
		logger.Error("error running services", "error", err)
		os.Exit(1)
	}
}

type ServiceConfig struct {
	Interval time.Duration `default:"30s" help:"interval to run the service"`
}

type fetchService struct {
	logger *slog.Logger
}

func (s *fetchService) Invoke(ctx context.Context) error {
	s.logger.InfoContext(ctx, "fetching data")
	time.Sleep(5 * time.Second)
	s.logger.InfoContext(ctx, "data fetched")
	return nil
}

func hourlyCron(ctx context.Context) error {
	slog.InfoContext(ctx, "hello world") //nolint:sloglint
	return nil
}

func longRunningJob(ctx context.Context) error {
	slog.InfoContext(ctx, "long running job started") //nolint:sloglint
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(10 * time.Second):
			slog.InfoContext(ctx, "long running job... still running") //nolint:sloglint
		}
	}
}
