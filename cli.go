// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

package clix

import (
	"log/slog"
	"os"
	"strconv"
	"sync/atomic"

	"github.com/alecthomas/kong"
	"github.com/joho/godotenv"
)

// Option is a function that can be used to configure the CLI.
type Option[T any] func(cli *CLI[T])

// WithKongOptions adds the provided kong options to the CLI.
func WithKongOptions[T any](opts ...kong.Option) Option[T] {
	return func(cli *CLI[T]) {
		cli.kongOptions = append(cli.kongOptions, opts...)
	}
}

// WithGithubDebug enables debug mode if the RUNNER_DEBUG environment variable is
// set, which is set when running a Github Action in debug mode.
func WithGithubDebug[T any]() Option[T] {
	var initialized atomic.Bool
	return func(cli *CLI[T]) {
		if initialized.Load() {
			return
		}
		cli.kongOptions = append(cli.kongOptions, kong.WithBeforeApply(func(cli *CLI[T]) error {
			if initialized.Swap(true) {
				return nil
			}
			if runnerDebug, _ := strconv.ParseBool(os.Getenv("RUNNER_DEBUG")); runnerDebug {
				cli.Debug = true
			}
			return nil
		}))
	}
}

// WithEnvFiles loads environment variables from ".env" style files from the
// provided paths. If no paths are provided, it will load from the current
// working directory.
func WithEnvFiles[T any](paths ...string) Option[T] {
	var initialized atomic.Bool
	return func(cli *CLI[T]) {
		if initialized.Load() {
			return
		}
		cli.kongOptions = append(cli.kongOptions, kong.WithBeforeReset(func() error {
			if initialized.Swap(true) {
				return nil
			}
			err := godotenv.Load(paths...)
			if err != nil && len(paths) > 0 {
				// Only throw an error if they explicitly provided paths.
				return err
			}
			return nil
		}))
	}
}

// CLI is the main construct for clix, obtained via [New] or [NewWithDefaults].
type CLI[T any] struct {
	kong.Plugins               // Kong-specific plugins.
	kongOptions  []kong.Option `kong:"-"`
	version      *Version      `kong:"-"`
	app          *AppInfo      `kong:"-"`
	logHandler   slog.Handler  `kong:"-"`
	logger       *slog.Logger  `kong:"-"`

	// Context is the context returned by kong after initial parsing.
	Context *kong.Context `kong:"-"`

	// Flags are the user-provided flags.
	Flags *T `embed:""`

	// Debug can be used to enable/disable debugging as a global flag. Also
	// sets the log level to debug.
	Debug bool `short:"D" name:"debug" help:"enables debug mode"`
}

// New executes the cli parser, with the provided options (no defaults are
// configured). Order of options is important. See also [NewWithDefaults].
//
// Supported struct tags: https://github.com/alecthomas/kong#supported-tags
func New[T any](options ...Option[T]) *CLI[T] {
	cli := &CLI[T]{
		Flags: new(T),
	}

	cli.app = &AppInfo{}
	cli.version = GetVersionInfo(cli.app)
	cli.kongOptions = []kong.Option{
		kong.ConfigureHelp(kong.HelpOptions{
			Tree:      true,
			FlagsLast: true,
		}),
		kong.UsageOnError(),
		kong.Description(cli.app.Description),
		kong.Bind(cli.version),
		kong.Bind(cli.app),
		kong.Bind(cli),
	}

	for _, opt := range options {
		if opt == nil {
			continue
		}
		opt(cli)
	}

	cli.Context = kong.Parse(cli, cli.kongOptions...)
	return cli
}

// Defaults returns the default options for clix.
func Defaults[T any]() []Option[T] {
	return []Option[T]{
		WithEnvFiles[T](),
		WithGithubDebug[T](),
		WithLoggingPlugin[T](true),
		WithVersionPlugin[T](),
		WithMarkdownPlugin[T](),
	}
}

// NewWithDefaults executes the cli parser, with the provided options and associated
// recommended defaults. Order of options is important. See also [New] and [Defaults].
//
// Supported struct tags: https://github.com/alecthomas/kong#supported-tags
//
// Initialize a new CLI like so:
//
//	type Flags struct {
//		// Normal kong flags.
//		SomeFlag string `env:"SOME_FLAG" help:"some flag"`
//	}
//
//	var cli = clix.NewWithDefaults[Flags]() // Where Flags is a user-provided type (struct).
//
//	func main() {
//		// [...]
//	}
func NewWithDefaults[T any](options ...Option[T]) *CLI[T] {
	return New(
		append(
			Defaults[T](),
			options...,
		)...,
	)
}
