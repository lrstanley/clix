// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

package clix

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"regexp"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/alecthomas/kong"
)

type testExitCodeError struct{ code int }

func (e testExitCodeError) Error() string { return fmt.Sprintf("exit code %d", e.code) }
func (e testExitCodeError) ExitCode() int { return e.code }

// newExecuteTestCLI returns a CLI with a captured exit code and log output.
func newExecuteTestCLI(t *testing.T) (cli *CLI[struct{}], code *int, logs *strings.Builder) {
	t.Helper()

	oldArgs := os.Args
	t.Cleanup(func() { os.Args = oldArgs })
	os.Args = []string{"testapp"}

	code = new(int)
	*code = -1
	cli = New(WithKongOptions[struct{}](kong.Exit(func(c int) { *code = c })))

	logs = &strings.Builder{}
	cli.logger = slog.New(slog.NewTextHandler(logs, &slog.HandlerOptions{AddSource: true}))
	return cli, code, logs
}

func TestExecute(t *testing.T) {
	tests := []struct {
		name     string
		fn       func(ctx context.Context, logger *slog.Logger) error
		wantCode int
		wantLogs []string
	}{
		{
			name:     "success",
			fn:       func(context.Context, *slog.Logger) error { return nil },
			wantCode: 0,
		},
		{
			name:     "error",
			fn:       func(context.Context, *slog.Logger) error { return errors.New("boom") },
			wantCode: exitCodeError,
			wantLogs: []string{"fatal error", "boom"},
		},
		{
			name: "exit-coder",
			fn: func(context.Context, *slog.Logger) error {
				return fmt.Errorf("wrapped: %w", testExitCodeError{code: 42})
			},
			wantCode: 42,
			wantLogs: []string{"fatal error", "exit code 42"},
		},
		{
			name:     "panic",
			fn:       func(context.Context, *slog.Logger) error { panic("kaboom") },
			wantCode: exitCodePanic,
			wantLogs: []string{"panic", "kaboom", "execute_test.go"},
		},
		{
			name: "canceled-without-signal",
			fn: func(context.Context, *slog.Logger) error {
				return fmt.Errorf("wrapped: %w", context.Canceled)
			},
			wantCode: exitCodeError,
			wantLogs: []string{"fatal error", "context canceled"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cli, code, logs := newExecuteTestCLI(t)
			cli.Execute(tt.fn)

			if *code != tt.wantCode {
				t.Fatalf("expected exit code %d, got %d (logs: %s)", tt.wantCode, *code, logs)
			}
			for _, want := range tt.wantLogs {
				if !strings.Contains(logs.String(), want) {
					t.Fatalf("expected logs to contain %q, got:\n%s", want, logs)
				}
			}
			if len(tt.wantLogs) == 0 && logs.Len() > 0 {
				t.Fatalf("expected no logs, got:\n%s", logs)
			}
		})
	}
}

func TestExecuteLogSource(t *testing.T) {
	var nilMap map[string]int

	tests := []struct {
		name string
		run  func(cli *CLI[struct{}]) (line int)
	}{
		{
			name: "error-uses-execute-caller",
			run: func(cli *CLI[struct{}]) int {
				_, _, line, _ := runtime.Caller(0)
				cli.Execute(func(context.Context, *slog.Logger) error { return errors.New("boom") })
				return line + 1
			},
		},
		{
			name: "panic-uses-panic-location",
			run: func(cli *CLI[struct{}]) (line int) {
				cli.Execute(func(context.Context, *slog.Logger) error {
					_, _, line, _ = runtime.Caller(0)
					panic("kaboom")
				})
				return line + 1
			},
		},
		{
			name: "runtime-panic-uses-panic-location",
			run: func(cli *CLI[struct{}]) (line int) {
				cli.Execute(func(context.Context, *slog.Logger) error {
					_, _, line, _ = runtime.Caller(0)
					nilMap["foo"] = 1
					return nil
				})
				return line + 1
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cli, _, logs := newExecuteTestCLI(t)
			line := tt.run(cli)

			want := regexp.MustCompile(fmt.Sprintf(`(?m)^.* source=\S*/execute_test\.go:%d `, line))
			if !want.MatchString(logs.String()) {
				t.Fatalf("expected log source to be execute_test.go:%d, got:\n%s", line, logs)
			}
		})
	}
}

func TestExecuteDeferredCleanup(t *testing.T) {
	cli, code, _ := newExecuteTestCLI(t)

	var cleaned bool
	cli.Execute(func(context.Context, *slog.Logger) error {
		defer func() { cleaned = true }()
		panic("kaboom")
	})

	if !cleaned {
		t.Fatal("expected deferred cleanup in fn to run")
	}
	if *code != exitCodePanic {
		t.Fatalf("expected exit code %d, got %d", exitCodePanic, *code)
	}
}

func TestExecuteContext(t *testing.T) {
	cli, _, _ := newExecuteTestCLI(t)

	cli.Execute(func(ctx context.Context, logger *slog.Logger) error {
		if FromContext[struct{}](ctx) != cli {
			t.Error("expected context to contain the CLI")
		}
		if logger != cli.logger {
			t.Error("expected logger to be the CLI logger")
		}
		return nil
	})
}

func TestExecuteSignal(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("sending interrupt signals is not supported on windows")
	}

	cli, code, logs := newExecuteTestCLI(t)

	cli.Execute(func(ctx context.Context, _ *slog.Logger) error {
		proc, err := os.FindProcess(os.Getpid())
		if err != nil {
			return err
		}
		if sigErr := proc.Signal(os.Interrupt); sigErr != nil {
			return sigErr
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(5 * time.Second):
			return errors.New("timed out waiting for context cancellation")
		}
	})

	if *code != 0 {
		t.Fatalf("expected exit code 0, got %d (logs: %s)", *code, logs)
	}
	for _, want := range []string{"shutting down", "interrupt signal received"} {
		if !strings.Contains(logs.String(), want) {
			t.Fatalf("expected logs to contain %q, got:\n%s", want, logs)
		}
	}
}
