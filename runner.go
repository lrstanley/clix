// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package clix

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"
)

type Runner func(ctx context.Context) error

func (r Runner) Invoke(ctx context.Context) func() error {
	fn := func() error {
		return r(ctx)
	}
	return fn
}

// Run invokes all runners concurrently, and listens for any termination signals
// (SIGINT, SIGTERM, SIGQUIT, etc).
//
// If any runners return an error, all runners will terminate (assuming they listen
// to the provided context), and the first known error will be returned.
func Run(runners ...Runner) error {
	return RunCtx(context.Background(), runners...)
}

// RunCtx is the same as Run, but with the provided context that can be used
// to externally cancel all runners.
func RunCtx(ctx context.Context, runners ...Runner) error {
	if len(runners) == 0 {
		panic("no runners provided")
	}

	var g *errgroup.Group
	g, ctx = errgroup.WithContext(ctx)

	g.Go(func() error {
		return signalListener(ctx)
	})

	for _, runner := range runners {
		g.Go(runner.Invoke(ctx))
	}

	return g.Wait()
}

func signalListener(ctx context.Context) error {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)

	select {
	case sig := <-quit:
		return fmt.Errorf("received signal: %v", sig)
	case <-ctx.Done():
		return nil
	}
}
