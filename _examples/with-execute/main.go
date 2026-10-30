// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

// An example that uses [clix.CLI.Execute] to wrap the application's main()
// function. Execute provides a context that is canceled on SIGINT/SIGTERM, logs
// returned errors and panics (with stack traces), and exits with an appropriate
// exit code, all while still running any deferred calls. Try it with:
//
//	$ go run .               # runs until Ctrl+C, then exits with code 0.
//	$ go run . --fail        # returns an error, exits with code 1.
//	$ go run . --exit-code 3 # returns an error with a custom exit code.
//	$ go run . --panic       # panics, exits with code 2.
package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/lrstanley/clix/v2"
)

type Flags struct {
	Fail     bool `name:"fail" help:"return an error after starting"`
	ExitCode int  `name:"exit-code" help:"return an error with a custom exit code after starting"`
	Panic    bool `name:"panic" help:"panic after starting"`
}

var cli = clix.NewWithDefaults(
	clix.WithAppInfo[Flags](clix.AppInfo{
		Name:        "with-execute",
		Description: "an example app that uses clix.Execute to exit cleanly",
		Links:       clix.GithubLinks("github.com/lrstanley/clix", "master", "https://liam.sh/"),
	}),
)

func main() {
	cli.Execute(run)
}

func run(ctx context.Context, cli *clix.CLI[Flags], logger *slog.Logger) error {
	conn, err := openConnection(ctx, logger)
	if err != nil {
		return fmt.Errorf("failed to open connection: %w", err)
	}
	// Runs before the application exits, even if run() returns an error or panics.
	defer conn.Close()

	switch {
	case cli.Flags.Fail:
		return errors.New("something went wrong")
	case cli.Flags.ExitCode != 0:
		return &exitCodeError{code: cli.Flags.ExitCode}
	case cli.Flags.Panic:
		panic("something went very wrong")
	}

	logger.InfoContext(ctx, "running, press Ctrl+C to stop")
	for {
		select {
		case <-ctx.Done():
			// Returning the context error after a signal results in a clean exit.
			return ctx.Err()
		case <-time.After(5 * time.Second):
			logger.InfoContext(ctx, "still running")
		}
	}
}

// exitCodeError implements [kong.ExitCoder], which Execute uses to determine
// the exit code.
type exitCodeError struct {
	code int
}

func (e *exitCodeError) Error() string {
	return fmt.Sprintf("failed with exit code %d", e.code)
}

func (e *exitCodeError) ExitCode() int {
	return e.code
}

type connection struct {
	logger *slog.Logger
}

func openConnection(ctx context.Context, logger *slog.Logger) (*connection, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	logger.InfoContext(ctx, "opening connection")
	return &connection{logger: logger}, nil
}

func (c *connection) Close() {
	c.logger.Info("closing connection")
}
