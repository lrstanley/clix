// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

package clix

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/alecthomas/kong"
	"github.com/lmittmann/tint"
)

// WithLoggingHandlerOptions allows customizing how slog handlers are created. This
// only applies when using the [WithLoggingPlugin] plugin. By default, the handler
// options are set to:
//
//   - Level: values passed from --log.level flag.
//   - AddSource: true
//   - ReplaceAttr: nil
func WithLoggingHandlerOptions[T any](opts *slog.HandlerOptions) Option[T] {
	return func(cli *CLI[T]) {
		cli.logHandlerOptions = opts
	}
}

// WithLoggingPlugin adds the logging plugin to the CLI. This includes flags
// for controlling [log/slog] logging levels, logging to files, JSON output, and
// supports setting the global slog logger. You can access the resulting
// [log/slog.Handler] via [CLI.GetLogHandler] and the [log/slog.Logger] via
// [CLI.GetLogger]. opts is optional, and can also be set using [WithLoggingHandlerOptions].
func WithLoggingPlugin[T any](global bool, opts *slog.HandlerOptions) Option[T] {
	var initialized atomic.Bool

	return func(cli *CLI[T]) {
		if initialized.Load() {
			return
		}

		if opts != nil {
			cli.logHandlerOptions = opts
		}

		var flags struct {
			Logging *LoggingPlugin `embed:"" group:"Logging flags"`
		}

		flags.Logging = &LoggingPlugin{}
		cli.Plugins = append(cli.Plugins, &flags)
		cli.kongOptions = append(cli.kongOptions, kong.WithAfterApply(func(kctx *kong.Context) error {
			if initialized.Swap(true) || cli.logHandler != nil {
				return nil
			}

			logger, err := flags.Logging.CreateHandler(cli.Debug, global, cli.logHandlerOptions)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error creating logger: %v\n", err)
				os.Exit(1)
			}

			cli.logHandler = logger
			cli.logger = slog.New(logger)

			kctx.Bind(cli.logHandler)
			kctx.Bind(cli.logger)

			cli.logger.Debug(
				"logger initialized",
				"name", cli.app.Name,
				"version", cli.app.Version,
				"commit", cli.app.Commit,
				"go_version", cli.version.GoVersion,
				"os", cli.version.OS,
				"arch", cli.version.Arch,
			)

			return nil
		}))
	}
}

// GetLogger returns the generated logger if enabled, nil otherwise. If global
// logging is enabled, you can also use [log/slog.Default].
func (c *CLI[T]) GetLogger() *slog.Logger {
	return c.logger
}

// GetLogHandler returns the generated logger handler if enabled, nil otherwise.
// Usually don't need to use this directly. Not used if logging configuration is
// disabled.
func (c *CLI[T]) GetLogHandler() slog.Handler {
	return c.logHandler
}

// LoggingPlugin are the flags that define how log entries are processed/returned.
type LoggingPlugin struct {
	// Level is the minimum level of log messages to output, must be one of none|debug|info|warn|error.
	Level string `name:"log.level" env:"LOG_LEVEL" default:"info" enum:"none,debug,info,warn,error" help:"logging level (none: disables logging)"`

	// JSON enables JSON logging.
	JSON bool `name:"log.json" env:"LOG_JSON" help:"output logs in JSON format"`

	// Path is the path to the log file.
	Path string `name:"log.path" env:"LOG_PATH" type:"path" help:"path to log file (disables stderr logging)"`
}

func (l *LoggingPlugin) GetLevel() slog.Level {
	switch l.Level {
	case "none":
		return -1
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

// CreateHandler creates a new [log/slog.Handler] with the provided configuration.
func (l *LoggingPlugin) CreateHandler(isDebug, setGlobal bool, opts *slog.HandlerOptions) (handler slog.Handler, err error) {
	level := l.GetLevel()

	if isDebug {
		level = slog.LevelDebug
	}

	if opts == nil {
		opts = &slog.HandlerOptions{
			Level:     level,
			AddSource: true,
		}
	}

	if opts.Level == nil {
		opts.Level = level
	}

	noColor, _ := strconv.ParseBool(os.Getenv("NO_COLOR"))

	switch {
	case l.Path != "":
		var f *os.File
		f, err = os.OpenFile(l.Path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0o600)
		if err != nil {
			return nil, err
		}

		// We can't really close the file here.
		handler = slog.NewJSONHandler(f, opts)
	case level == -1:
		handler = slog.DiscardHandler
	case l.JSON:
		handler = slog.NewJSONHandler(os.Stderr, opts)
	case noColor:
		handler = slog.NewTextHandler(os.Stderr, opts)
	default:
		handler = tint.NewHandler(
			os.Stderr,
			&tint.Options{
				Level:      opts.Level,
				AddSource:  opts.AddSource,
				TimeFormat: time.TimeOnly,
				NoColor:    noColor,
				ReplaceAttr: func(groups []string, attr slog.Attr) slog.Attr {
					if attr.Value.Kind() == slog.KindAny {
						if _, ok := attr.Value.Any().(error); ok {
							attr = tint.Attr(9, attr)
						}
					}
					if opts.ReplaceAttr != nil {
						return opts.ReplaceAttr(groups, attr)
					}
					return attr
				},
			},
		)
	}

	if setGlobal {
		_ = slog.SetLogLoggerLevel(level)
		slog.SetDefault(slog.New(handler))
	}

	return handler, nil
}
