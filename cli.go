// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package clix

import (
	"fmt"
	"os"

	"github.com/apex/log"
	flags "github.com/jessevdk/go-flags"
)

// Options allows overriding default logic.
type Options int

const (
	OptDisableLogging       Options = 1 << iota // Disable logging initialization.
	OptDisableVersion                           // Disable version printing (must handle manually).
	OptDisableDeps                              // Disable dependency printing in version output.
	OptDisableBuildSettings                     // Disable build info printing in version output.
	OptDisableGlobalLogger                      // Disable setting the global logger for apex/log.
)

// New returns a new CLI instance, initializes and executes the go-flags
// parser. Make sure to call Parse() after invoking New().
//
// * Use cli.Logger as a apex/log log.Interface.
// * Use cli.Parser to change/add configuration to the go-flags parser.
// * Use cli.Args to get the remaining arguments provided to the program.
//
// Example usage:
//  var logger log.Interface
//  type Flags struct {
//  	SomeFlag string `long:"some-flag" description:"some flag"`
//  }
//  cli := clix.New[Flags](clix.OptDisableGlobalLogger|clix.OptDisableBuildSettings).Parse()
//  logger = cli.Logger
//
// Use it like so:
//  cli.Flags.SomeFlag  // some value
//  cli.Debug           // true|false
//  logger = cli.Logger // initialize your logger using the one from CLI.
func New[T any](options Options) (cli *CLI[T]) {
	cli = &CLI[T]{
		options: options,
		Flags:   new(T),
	}

	cli.Parser = flags.NewParser(cli, flags.HelpFlag|flags.PassDoubleDash)

	cli.Parser.NamespaceDelimiter = "."
	cli.Parser.EnvNamespaceDelimiter = "_"

	return cli
}

// CLI wraps user-provided Flags, with a logger, and Version/Debug helpers.
// Args contains the arguments passed to the binary, after flags have been
// parsed.
type CLI[T any] struct {
	// Flags are the user-provided flags.
	Flags *T

	// Args are the remaining arguments after parsing.
	Args []string

	// Version can be used to print the version information to console.
	Version bool `short:"v" long:"version" description:"prints version information and exits"`

	// Debug can be used to enable/disable debugging as a global flag. Also
	// sets the log level to debug.
	Debug bool `short:"D" long:"debug" env:"DEBUG" description:"enables debug mode"`

	// Logger is the generated logger.
	Logger *log.Logger `json:"-"`

	// Parser is the go-flags parser, which can be used to add/change parser
	// configurations.
	Parser *flags.Parser `json:"-"`

	logConfig LoggerConfig `group:"Logging Options" namespace:"log" env-namespace:"LOG"`
	options   Options      `json:"-"`
}

// Parse executes the go-flags parser, returns the remaining arguments, as
// well as initializes a new logger. If cli.Version is set, it will print
// the version information (unless disabled).
func (cli *CLI[T]) Parse() *CLI[T] {
	args, err := cli.Parser.Parse()
	if err != nil {
		if FlagErr, ok := err.(*flags.Error); ok && FlagErr.Type == flags.ErrHelp {
			os.Exit(0)
		}

		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	cli.Args = args

	// Initialize the logger.
	if !cli.IsSet(OptDisableLogging) {
		cli.newLogger()
	}

	if cli.Version && !cli.IsSet(OptDisableVersion) {
		cli.PrintVersion("", "", "")
		os.Exit(1)
	}

	return cli
}

// IsSet returns true if the given option is set.
func (cli *CLI[T]) IsSet(options Options) bool {
	return cli.options&options != 0
}

// Set sets the given option.
func (cli *CLI[T]) Set(option Options) {
	cli.options |= option
}
