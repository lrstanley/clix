// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

// A simple example using the out-of-the-box functionality with clix, which includes
// logging, version flags, loading of .env files, etc. Take a look at the help output
// with:
//
//	$ go run . --help
//
// Or use the markdown generation logic using:
//
//	$ go run . generate-markdown > USAGE.md
package main

import (
	"github.com/lrstanley/clix/v2"
)

type Flags struct {
	Name string `name:"name" default:"world" help:"name to print"`
}

var cli = clix.NewWithDefaults[Flags]()

func main() {
	logger := cli.GetLogger()

	if cli.Debug {
		logger.Debug("thinking really hard...")
	}

	logger.Info("hello", "name", cli.Flags.Name)
}
