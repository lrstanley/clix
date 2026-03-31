// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

package main

import (
	"time"

	"github.com/alecthomas/kong"
	"github.com/lrstanley/clix/v2"
)

type Flags struct {
	Name   string `name:"name" default:"world" help:"name to print: ${foo}"`
	Hidden string `name:"hidden" hidden:"" help:"hidden flag"`

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

	Foo struct {
		Bar struct {
			Config map[string]string `arg:"" help:"config"`
		} `cmd:"" help:"sub-command bar"`
	} `cmd:"" help:"sub-command test"`

	Test    string `name:"test" required:"" enum:"foo,bar" help:"test flag"`
	TestFoo string `name:"test-foo" default:"foo" env:"FOO,BAR" enum:"foo,bar" help:"test flag"`
}

var cli = clix.NewWithDefaults(
	clix.WithAppInfo[Flags](clix.AppInfo{
		Name:        "simple-app",
		Description: "a simple app that prints hello world",
		Links:       clix.GithubLinks("github.com/lrstanley/clix", "master", "https://liam.sh/"),
	}),
	clix.WithKongOptions[Flags](
		kong.Vars{"foo": "bar"},
	),
)

func main() {
	cli.GetLogger().Info("hello", "name", cli.Flags.Name)
}
