// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

package clix

import (
	"log/slog"
	"os"
	"slices"
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
	return func(cli *CLI[T]) {
		if cli.checkAlreadyInit("github-debug") {
			return
		}
		cli.kongOptions = append(cli.kongOptions, kong.WithBeforeApply(func(cli *CLI[T]) error {
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
	return func(cli *CLI[T]) {
		if cli.checkAlreadyInit("env-files") {
			return
		}
		cli.kongOptions = append(cli.kongOptions, kong.WithBeforeReset(func() error {
			err := godotenv.Load(paths...)
			if err != nil && len(paths) > 0 {
				// Only throw an error if they explicitly provided paths.
				return err
			}
			return nil
		}))
	}
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
	kong.Plugins                     // Kong-specific plugins.
	kongOptions        []kong.Option `kong:"-"`
	pluginsInitialized []string      `kong:"-"`
	version            *Version      `kong:"-"`
	app                *AppInfo      `kong:"-"`
	logHandler         slog.Handler  `kong:"-"`
	logger             *slog.Logger  `kong:"-"`

	// Context is the context returned by kong after initial parsing.
	Context *kong.Context `kong:"-"`

	// Flags are the user-provided flags.
	Flags *T `embed:""`

	// Debug can be used to enable/disable debugging as a global flag. Also
	// sets the log level to debug.
	Debug bool `short:"D" name:"debug" help:"enables debug mode"`
}

// Parse executes the cli parser, with the provided options (no defaults are
// configured). Order of options is important. See also [CLI.ParseWithDefaults].
func (cli *CLI[T]) Parse(options ...Option[T]) {
	if cli.Flags == nil {
		cli.Flags = new(T)
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

func (cli *CLI[T]) checkAlreadyInit(plugin string) bool {
	if slices.Contains(cli.pluginsInitialized, plugin) {
		return true
	}
	cli.pluginsInitialized = append(cli.pluginsInitialized, plugin)
	return false
}
