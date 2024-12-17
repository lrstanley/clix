// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

package clix

import (
	"os"

	"github.com/apex/log"
	logcli "github.com/apex/log/handlers/cli"
	"github.com/apex/log/handlers/discard"
	"github.com/apex/log/handlers/json"
	"github.com/apex/log/handlers/logfmt"
	"github.com/apex/log/handlers/text"
	"github.com/lrstanley/clix/githubhandler"
)

// LoggerConfig are the flags that define how log entries are processed/returned.
// If using github.com/jessevdk/go-flags, you can defined a LoggerConfig in
// your struct, and then call LoggerConfig.New(isDebug), directly, so you
// don't have to define additional flags. See the struct tags for what it will
// default to in terms of environment variables.
//
// Example (where you can set LOG_LEVEL as an environment variable, for example):
//
//	type Flags struct {
//		Debug    bool               `long:"debug" env:"DEBUG" description:"enable debugging"`
//		Log      *chix.LoggerConfig `group:"Logging Options" namespace:"log" env-namespace:"LOG"`
//	}
//	[...]
//	cli.Log.New(cli.Debug)
type LoggerConfig struct {
	// Quiet disables all logging.
	Quiet bool `env:"QUIET" long:"quiet" description:"disable logging to stdout (also: see levels)"`

	// Level is the minimum level of log messages to output, must be one of info|warn|error|debug|fatal.
	Level string `env:"LEVEL" long:"level" default:"info" choice:"debug" choice:"info" choice:"warn" choice:"error" choice:"fatal" description:"logging level"`

	// JSON enables JSON logging.
	JSON bool `env:"JSON" long:"json" description:"output logs in JSON format"`

	// Github enables GitHub Actions logging.
	Github bool `env:"GITHUB" long:"github" description:"output logs in GitHub Actions format"`

	// Pretty enables cli-friendly logging.
	Pretty bool `env:"PRETTY" long:"pretty" description:"output logs in a pretty colored format (cannot be easily parsed)"`

	// Path is the path to the log file.
	Path string `env:"PATH" long:"path" description:"path to log file (disables stdout logging)"`
}

// new parses LoggerConfig and creates a new structured logger with the
// provided configuration.
func (cli *CLI[T]) newLogger() error {
	cli.Logger = &log.Logger{}

	if cli.LoggerConfig.Level == "" {
		cli.Logger.Level = log.InfoLevel
	}

	if cli.Debug {
		cli.Logger.Level = log.DebugLevel
	} else if cli.LoggerConfig.Level == "" {
		cli.Logger.Level = log.InfoLevel
	} else {
		cli.Logger.Level = log.MustParseLevel(cli.LoggerConfig.Level)
	}

	switch {
	case cli.LoggerConfig.Path != "":
		f, err := os.OpenFile(cli.LoggerConfig.Path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0o644)
		if err != nil {
			return err
		}

		// We can't really close the file here.

		cli.Logger.Handler = logcli.New(f)
	case cli.LoggerConfig.Github:
		// Since debug is by default masked unless debugging is enabled in Actions.
		cli.Logger.Level = log.DebugLevel
		cli.Logger.Handler = githubhandler.New(os.Stdout)
	case cli.LoggerConfig.Quiet:
		cli.Logger.Handler = discard.New()
	case cli.LoggerConfig.JSON:
		cli.Logger.Handler = json.New(os.Stdout)
	case cli.LoggerConfig.Pretty:
		cli.Logger.Handler = text.New(os.Stdout)
	default:
		cli.Logger.Handler = logfmt.New(os.Stdout)
	}

	if cli.options&OptDisableGlobalLogger == 0 {
		log.SetLevel(cli.Logger.Level)
		log.SetHandler(cli.Logger.Handler)
	}

	return nil
}
