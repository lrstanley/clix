// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

package clix

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"time"

	"github.com/alecthomas/kong"
	"github.com/lmittmann/tint"
	"github.com/lrstanley/x/logging/handlers"
)

// WithLoggingPlugin adds the logging plugin to the CLI. This includes flags
// for controlling log/slog logging levels, logging to files, JSON output, and
// supports setting the global slog logger. You can access the resulting
// [log/slog.Handler] via [CLI.GetLogHandler] and the [log/slog.Logger] via
// [CLI.GetLogger].
func WithLoggingPlugin[T any](global bool) Option[T] {
	return func(cli *CLI[T]) {
		if cli.checkAlreadyInit("logging") {
			return
		}

		var flags struct {
			Logging *LoggingPlugin `embed:"" group:"Logging flags"`
		}

		flags.Logging = &LoggingPlugin{}
		cli.Plugins = append(cli.Plugins, &flags)
		cli.kongOptions = append(cli.kongOptions, kong.WithAfterApply(func() error {
			logger, err := flags.Logging.CreateHandler(cli.Debug, global)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error creating logger: %v\n", err)
				os.Exit(1)
			}

			cli.logHandler = logger
			cli.logger = slog.New(logger)

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
func (l *LoggingPlugin) CreateHandler(isDebug, setGlobal bool) (handler slog.Handler, err error) {
	level := l.GetLevel()

	if isDebug {
		level = slog.LevelDebug
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
		handler = slog.NewJSONHandler(
			f,
			&slog.HandlerOptions{
				Level:     level,
				AddSource: true,
			},
		)
	case level == -1:
		handler = handlers.NewDiscard()
	case l.JSON:
		handler = slog.NewJSONHandler(
			os.Stderr,
			&slog.HandlerOptions{
				Level:     level,
				AddSource: true,
			},
		)
	case noColor:
		handler = slog.NewTextHandler(
			os.Stderr,
			&slog.HandlerOptions{
				Level:     level,
				AddSource: true,
			},
		)
	default:
		handler = tint.NewHandler(
			os.Stderr,
			&tint.Options{
				Level:      level,
				AddSource:  true,
				TimeFormat: time.RFC3339,
				NoColor:    noColor,
			},
		)
	}

	if setGlobal {
		_ = slog.SetLogLoggerLevel(level)
		slog.SetDefault(slog.New(handler))
	}

	return handler, nil
}
