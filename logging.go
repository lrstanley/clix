// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package clix

import (
	"os"

	"github.com/apex/log"
	"github.com/apex/log/handlers/discard"
	"github.com/apex/log/handlers/json"
	"github.com/apex/log/handlers/logfmt"
	"github.com/apex/log/handlers/text"
)

// LoggerConfig are the flags that define how log entries are processed/returned.
// If using github.com/jessevdk/go-flags, you can defined a LoggerConfig in
// your struct, and then call LoggerConfig.New(isDebug), directly, so you
// don't have to define additional flags. See the struct tags for what it will
// default to in terms of environment variables.
//
// Example (where you can set LOG_LEVEL as an environment variable, for example):
//
//   type Flags struct {
//   	Debug    bool               `long:"debug" env:"DEBUG" description:"enable debugging"`
//   	Log      *chix.LoggerConfig `group:"Logging Options" namespace:"log" env-namespace:"LOG"`
//   }
//   [...]
//   cli.Log.New(cli.Debug)
type LoggerConfig struct {
	// Quiet disables all logging.
	Quiet bool `env:"QUIET" long:"quiet" description:"disable logging to stdout (also: see levels)"`

	// Level is the minimum level of log messages to output, must be one of info|warn|error|debug|fatal.
	Level string `env:"LEVEL" long:"level" default:"info" choice:"debug" choice:"info" choice:"warn" choice:"error" choice:"fatal" description:"logging level"`

	// JSON enables JSON logging.
	JSON bool `env:"JSON" long:"json" description:"output logs in JSON format"`

	// Pretty enables cli-friendly logging.
	Pretty bool `env:"PRETTY" long:"pretty" description:"output logs in a pretty colored format (cannot be easily parsed)"`
}

// new parses LoggerConfig and creates a new structured logger with the
// provided configuration.
func (cli *CLI[T]) newLogger() {
	cli.Logger = &log.Logger{}

	if cli.logConfig.Level == "" {
		cli.Logger.Level = log.InfoLevel
	}

	if cli.Debug {
		cli.Logger.Level = log.DebugLevel
	} else if cli.logConfig.Level == "" {
		cli.Logger.Level = log.InfoLevel
	} else {
		cli.Logger.Level = log.MustParseLevel(cli.logConfig.Level)
	}

	switch {
	case cli.logConfig.Quiet:
		cli.Logger.Handler = discard.New()
	case cli.logConfig.JSON:
		cli.Logger.Handler = json.New(os.Stdout)
	case cli.logConfig.Pretty:
		cli.Logger.Handler = text.New(os.Stdout)
	default:
		cli.Logger.Handler = logfmt.New(os.Stdout)
	}

	if cli.options&OptDisableGlobalLogger == 0 {
		log.SetLevel(cli.Logger.Level)
		log.SetHandler(cli.Logger.Handler)
	}
}
