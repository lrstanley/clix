// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

// Example app that uses all of the clix defaults, and includes support for multiple
// commands.
package main

import (
	"os"
	"time"

	"github.com/lrstanley/clix/v2"
)

type Flags struct {
	RM struct {
		User      string        `help:"Run as user." short:"u" default:"default"`
		Force     bool          `help:"Force removal." short:"f"`
		Recursive bool          `help:"Recursively remove files." short:"r"`
		Delay     time.Duration `help:"testing time.Duration." default:"1s"`
		Paths     []string      `arg:"" help:"Paths to remove." type:"path" name:"path"`
	} `cmd:"" help:"Remove files."`

	LS struct {
		Paths []string `arg:"" optional:"" help:"Paths to list." placeholder:"<paths>" type:"path"`
	} `cmd:"" help:"List paths."`

	Status struct{} `cmd:"" default:"" help:"Get status information."`
}

var cli = clix.NewWithDefaults(
	clix.WithAppInfo[Flags](clix.AppInfo{
		Name:        "simple-app",
		Description: "a simple app that supports multiple commands",
		Links:       clix.GithubLinks("github.com/lrstanley/clix", "master", "https://liam.sh/"),
	}),
	// You can provide any normal kong options here like so:
	// clix.WithKongOptions[Flags](
	// 	kong.Name("simple-app"),
	// ),
)

func main() {
	logger := cli.GetLogger()

	switch cli.Context.Command() {
	case "rm <path>":
		logger.Info("removing path(s)", "paths", cli.Flags.RM.Paths)
		// [...]
	case "ls":
		wd, err := os.Getwd()
		if err != nil {
			logger.Error("failed to get working directory", "error", err)
			return
		}
		logger.Info("listing paths", "paths", wd)
		// [...]
	case "ls <paths>":
		logger.Info("listing paths", "paths", cli.Flags.LS.Paths)
		// [...]
	case "status":
		logger.Info("getting status information...")
		// [...]
	default:
		panic("unknown command")
	}
}
