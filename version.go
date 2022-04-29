// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package clix

import (
	"fmt"
	"runtime"
	"runtime/debug"
)

func (cli *CLI[T]) PrintVersion(version, commit, date string) {
	if version != "" {
		if commit == "" {
			commit = "unknown"
		}
		if date == "" {
			date = "unknown"
		}

		fmt.Printf("%s/%s %s (built with: %s)\n", version, commit, date, runtime.Version())
		return
	}

	build, ok := debug.ReadBuildInfo()
	if !ok {
		fmt.Printf("build information not available\n")
		return
	}

	fmt.Printf("%s/%s (%s)\n", build.Main.Path, build.Main.Version, build.GoVersion)

	if !cli.IsSet(OptDisableBuildSettings) {
		fmt.Printf("\nbuild settings:\n")
		for _, option := range build.Settings {
			fmt.Printf("  %s: %s\n", option.Key, option.Value)
		}
	}

	if !cli.IsSet(OptDisableDeps) {
		fmt.Printf("\ndependencies:\n")
		for _, dep := range build.Deps {
			if dep.Replace != nil {
				dep = dep.Replace
			}
			fmt.Printf("  %s :: %s :: %s\n", dep.Sum, dep.Path, dep.Version)
		}
	}
}
