// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

package clix

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"os/signal"
	"runtime"
	"runtime/debug"
	"strings"
	"syscall"
	"time"

	"github.com/alecthomas/kong"
)

const (
	exitCodeError = 1
	exitCodePanic = 2 // Matches the exit code used by the Go runtime for panics.
)

// Execute is an optional helper for wrapping the body of an application's main()
// function, so it can exit cleanly. It is only useful as the last call in main(),
// as it never returns (unless the exit function was overridden via [kong.Exit]).
// Usage:
//
//	func main() {
//		cli.Execute(run)
//	}
//
//	func run(ctx context.Context, logger *slog.Logger) error {
//		db, err := openDB(ctx)
//		if err != nil {
//			return err
//		}
//		defer db.Close() // Runs even if run() returns an error or panics.
//		// [...]
//	}
//
// Behavior:
//
//   - The provided context is canceled when SIGINT or SIGTERM is received (see
//     [signal.NotifyContext]). After the first signal, default signal handling is
//     restored, so a second signal forcefully terminates the application. The
//     context also contains the CLI (see [FromContext]).
//   - The provided logger is the logger from [CLI.GetLogger], or [slog.Default]
//     if the logging plugin is disabled.
//   - If fn returns nil, or returns [context.Canceled] after a signal was
//     received, the application exits with code 0.
//   - If fn returns an error, it is logged and the application exits with code 1,
//     or the code from the error if it implements [kong.ExitCoder].
//   - If fn panics, the panic and its stack trace are logged and the application
//     exits with code 2. Note that only panics from the goroutine running fn can
//     be recovered.
//
// Log entries from Execute use the location of the panic as their source, or
// otherwise the location where Execute was called, rather than Execute itself.
//
// Because fn returns rather than calling [os.Exit] directly, deferred calls in fn
// always run before the application exits, unlike with log.Fatal-style helpers.
//
// Alternatively, kong has built-in ways of achieving something similar:
//
//   - Command structs can implement a Run(...) error method, which is invoked
//     through [kong.Context.Run] (using cli.Context.Run(...)) for the selected
//     command. Its arguments are injected from bound values (see [kong.Bind],
//     [kong.BindTo], and [kong.Context.Bind]), and the logger is already bound
//     by the logging plugin. Hooks like AfterApply can also be used to run logic
//     before any Run method.
//   - [kong.Kong.FatalIfErrorf] (using cli.Context.FatalIfErrorf(err)) prints
//     the error to stderr (not through the logger), and exits with code 1, or
//     the code from [kong.ExitCoder]. Deferred calls are skipped, and panics are
//     not handled.
//
// These can be combined with Execute, by returning the result of cli.Context.Run
// from fn.
func (c *CLI[T]) Execute(fn func(ctx context.Context, logger *slog.Logger) error) {
	logger := c.logger
	if logger == nil {
		logger = slog.Default()
	}

	var pcs [1]uintptr
	runtime.Callers(2, pcs[:]) // Skip [runtime.Callers] and Execute.
	callerPC := pcs[0]

	ctx, stop := signal.NotifyContext(c.NewContext(context.Background()), os.Interrupt, syscall.SIGTERM)

	signalDone := make(chan struct{})
	context.AfterFunc(ctx, func() {
		defer close(signalDone)
		// Signal causes also match context.Canceled via errors.Is, so compare
		// directly to only skip the cancellation from calling stop().
		if cause := context.Cause(ctx); cause != context.Canceled { //nolint:errorlint
			logAt(ctx, logger, callerPC, slog.LevelInfo, "shutting down", "reason", cause.Error())
		}
		stop() // Restore default signal handling, so another signal forces an exit.
	})

	code := execute(ctx, logger, callerPC, fn)

	// Cancel the context (if not already) and wait for the above to finish, so
	// any shutdown log entry is written before exiting.
	stop()
	<-signalDone

	if c.Context != nil && c.Context.Exit != nil {
		c.Context.Exit(code)
		return
	}
	os.Exit(code)
}

func execute(
	ctx context.Context,
	logger *slog.Logger,
	callerPC uintptr,
	fn func(ctx context.Context, logger *slog.Logger) error,
) (code int) {
	defer func() {
		if r := recover(); r != nil {
			pc := panicPC()
			if pc == 0 {
				pc = callerPC
			}
			logAt(ctx, logger, pc, slog.LevelError, "panic", "error", r, "stack", string(debug.Stack()))
			code = exitCodePanic
		}
	}()

	err := fn(ctx, logger)
	if err == nil || (errors.Is(err, context.Canceled) && ctx.Err() != nil) {
		return 0
	}

	logAt(ctx, logger, callerPC, slog.LevelError, "fatal error", "error", err)

	var exitCoder kong.ExitCoder
	if errors.As(err, &exitCoder) {
		return exitCoder.ExitCode()
	}
	return exitCodeError
}

// logAt logs a message using pc as the source location, rather than the caller.
func logAt(ctx context.Context, logger *slog.Logger, pc uintptr, level slog.Level, msg string, args ...any) {
	if !logger.Enabled(ctx, level) {
		return
	}
	r := slog.NewRecord(time.Now(), level, msg, pc)
	r.Add(args...)
	_ = logger.Handler().Handle(ctx, r)
}

// panicPC returns the program counter of the frame that panicked, or 0 if it
// can't be found. It must be called from a deferred function while panicking.
func panicPC() uintptr {
	pcs := make([]uintptr, 64)
	frames := runtime.CallersFrames(pcs[:runtime.Callers(1, pcs)])

	var panicking bool
	for {
		frame, more := frames.Next()
		switch {
		case frame.Function == "runtime.gopanic":
			panicking = true
		case panicking && !strings.HasPrefix(frame.Function, "runtime."):
			return frame.PC
		}
		if !more {
			return 0
		}
	}
}
