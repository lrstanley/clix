// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

package clix

import (
	"log/slog"
	"os"
	"strconv"

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
	return WithKongOptions[T](
		kong.WithBeforeApply(func(cli *CLI[T]) error {
			if runnerDebug, _ := strconv.ParseBool(os.Getenv("RUNNER_DEBUG")); runnerDebug {
				cli.Debug = true
			}
			return nil
		}),
	)
}

// WithEnvFiles loads environment variables from ".env" style files from the
// provided paths. If no paths are provided, it will load from the current
// working directory.
func WithEnvFiles[T any](paths ...string) Option[T] {
	return WithKongOptions[T](
		kong.WithBeforeResolve(func() error {
			err := godotenv.Load(paths...)
			if err != nil && len(paths) > 0 {
				// Only throw an error if they explicitly provided paths.
				return err
			}
			return nil
		}),
	)
}

type Application struct {
	Name        string `json:"name"`          // Application name. Defaults to the main module path.
	Description string `json:"description"`   // Application description.
	Version     string `json:"build_version"` // Build version. Uses VCS info if available.
	Commit      string `json:"build_commit"`  // VCS commit SHA. Uses VCS info if available.
	Date        string `json:"build_date"`    // VCS commit date. Uses VCS info if available.

	Links []Link `json:"links,omitempty"` // Links to the project's website, support, issues, security, etc.
}

// CLI is the main construct for clix. Do not manually set any fields until
// you've called [Parse] or [ParseWithDefaults].
//
// Supported struct tags: https://github.com/alecthomas/kong#supported-tags
//
// Initialize a new CLI like so:
//
//	var (
//		cli    = &clix.CLI[Flags]{} // Where Flags is a user-provided type (struct).
//	)
//
//	type Flags struct {
//		// Normal kong flags.
//		SomeFlag string `env:"SOME_FLAG" help:"some flag"`
//	}
//
//	func main() {
//		cli.ParseWithDefaults()
//		// [...]
//	}
type CLI[T any] struct {
	kong.Plugins               // Kong-specific plugins.
	kongOptions  []kong.Option `kong:"-"`
	version      *Version      `kong:"-"`

	// Application should contain basic information about your application.
	Application Application `kong:"-"`

	// Context is the context returned by kong after initial parsing.
	Context *kong.Context `kong:"-"`

	// Flags are the user-provided flags.
	Flags *T `embed:""`

	// Debug can be used to enable/disable debugging as a global flag. Also
	// sets the log level to debug.
	Debug bool `short:"D" name:"debug" help:"enables debug mode"`

	// LogHandler is the generated logger handler. Usually don't need to use this
	// directly. Not used if logging configuration is disabled.
	LogHandler slog.Handler `kong:"-"`

	// Logger is the generated logger. Not used if logging configuration is disabled.
	Logger *slog.Logger `kong:"-"`
}

// Parse executes the cli parser, with the provided options (no defaults are
// configured). Order of options is important. See also [CLI.ParseWithDefaults].
func (cli *CLI[T]) Parse(options ...Option[T]) {
	if cli.Flags == nil {
		cli.Flags = new(T)
	}

	cli.version = GetVersionInfo(cli.Application)
	cli.Application = cli.version.Application // The version info can also help fill out the app info.
	cli.kongOptions = []kong.Option{
		kong.ConfigureHelp(kong.HelpOptions{
			Tree:      true,
			FlagsLast: true,
		}),
		kong.UsageOnError(),
		kong.Name(cli.Application.Name),
		kong.Description(cli.Application.Description),
		kong.Bind(cli.version),
		kong.Bind(cli.Application),
		kong.Bind(cli),
	}

	for _, opt := range options {
		if opt == nil {
			continue
		}
		opt(cli)
	}

	cli.Context = kong.Parse(cli, cli.kongOptions...)
}

// ParseWithDefaults executes the cli parser, with the provided options and associated
// recommended defaults. Order of options is important. See also [CLI.Parse].
func (cli *CLI[T]) ParseWithDefaults(options ...Option[T]) {
	cli.Parse(
		append(
			[]Option[T]{
				WithEnvFiles[T](),
				WithGithubDebug[T](),
				WithLoggingPlugin[T](true),
				WithVersionPlugin[T](),
				WithMarkdownPlugin[T](),
			},
			options...,
		)...,
	)
}
