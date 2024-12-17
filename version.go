// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

package clix

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"strings"

	"github.com/gookit/color"
)

// Module represents a module.
type Module struct {
	Path    string  `json:"path,omitempty"`     // module path
	Version string  `json:"version,omitempty"`  // module version
	Sum     string  `json:"sum,omitempty"`      // checksum
	Replace *Module `json:"replaces,omitempty"` // replaced by this module
}

func (m Module) String() string {
	if m.Replace != nil {
		m = *m.Replace
	}

	return fmt.Sprintf("%s :: %s :: %s", m.Sum, m.Path, m.Version)
}

// BuildSetting describes a setting that may be used to understand how the
// binary was built. For example, VCS commit and dirty status is stored here.
type BuildSetting struct {
	// Key and Value describe the build setting.
	// Key must not contain an equals sign, space, tab, or newline.
	// Value must not contain newlines ('\n').
	Key   string `json:"key"`
	Value string `json:"value"`
}

func (s BuildSetting) String() string {
	return fmt.Sprintf("%s: %s", s.Key, s.Value)
}

// VersionInfo represents the version information for the CLI.
type VersionInfo[T any] struct {
	Name         string         `json:"name"`                     // Name of cli tool.
	Version      string         `json:"build_version"`            // Build version.
	Commit       string         `json:"build_commit"`             // VCS commit SHA.
	Date         string         `json:"build_date"`               // VCS commit date.
	Settings     []BuildSetting `json:"build_settings,omitempty"` // Other information about the build.
	Dependencies []Module       `json:"dependencies,omitempty"`   // Module dependencies.

	Command   string `json:"command"`    // Executable name where the command was called from.
	GoVersion string `json:"go_version"` // Version of Go that produced this binary.
	OS        string `json:"os"`         // Operating system for this build.
	Arch      string `json:"arch"`       // CPU Architecture for this build.

	// Items hoisted from the parent CLI. Do not change this.
	Links []Link `json:"links,omitempty"`

	cli *CLI[T] `json:"-"`
}

// NonSensitiveVersion represents the version information for the CLI.
type NonSensitiveVersion struct {
	Name    string `json:"name"`          // Name of cli tool.
	Version string `json:"build_version"` // Build version.
	Commit  string `json:"build_commit"`  // VCS commit SHA.
	Date    string `json:"build_date"`    // VCS commit date.

	Command   string `json:"command"`    // Executable name where the command was called from.
	GoVersion string `json:"go_version"` // Version of Go that produced this binary.
	OS        string `json:"os"`         // Operating system for this build.
	Arch      string `json:"arch"`       // CPU Architecture for this build.

	// Items hoisted from the parent CLI. Do not change this.
	Links []Link `json:"links,omitempty"`
}

// NonSensitive returns a copy of VersionInfo with sensitive information removed.
func (v *VersionInfo[T]) NonSensitive() *NonSensitiveVersion {
	return &NonSensitiveVersion{
		Name:    v.Name,
		Version: v.Version,
		Commit:  v.Commit,
		Date:    v.Date,

		Command:   v.Command,
		GoVersion: v.GoVersion,
		OS:        v.OS,
		Arch:      v.Arch,

		Links: v.Links,
	}
}

// GetSetting returns the value of the setting with the given key, otherwise
// defaults to defaultValue.
func (v *VersionInfo[T]) GetSetting(key, defaultValue string) string {
	if v.Settings == nil {
		return defaultValue
	}

	for _, s := range v.Settings {
		if s.Key == key {
			return s.Value
		}
	}

	return defaultValue
}

func (v *VersionInfo[T]) stringBase() string {
	w := &bytes.Buffer{}

	fmt.Fprintf(w, "<cyan>%s</> :: <yellow>%s</>\n", v.Name, v.Version)
	fmt.Fprintf(w, "|  build commit :: <green>%s</>\n", v.Commit)
	fmt.Fprintf(w, "|    build date :: <green>%s</>\n", v.Date)
	fmt.Fprintf(w, "|    go version :: <green>%s %s/%s</>\n", v.GoVersion, v.OS, v.Arch)

	if len(v.Links) > 0 {
		var longest int
		for _, l := range v.Links {
			if len(l.Name) > longest {
				longest = len(l.Name)
			}
		}

		fmt.Fprintf(w, "\n<cyan>helpful links:</>\n")
		for _, l := range v.Links {
			fmt.Fprintf(
				w, "|  %s%s :: <magenta>%s</>\n",
				strings.Repeat(" ", longest-len(l.Name)),
				l.Name, l.URL,
			)
		}
	}

	return w.String()
}

func (v *VersionInfo[T]) String() string {
	w := &bytes.Buffer{}

	w.WriteString(v.stringBase())

	if !v.cli.IsSet(OptDisableBuildSettings) {
		var longest int
		for _, s := range v.Settings {
			if len(s.Key) > longest {
				longest = len(s.Key)
			}
		}

		fmt.Fprintf(w, "\n<cyan>build options:</>\n")
		for _, s := range v.Settings {
			fmt.Fprintf(
				w, "|  %s%s :: <magenta>%s</>\n",
				strings.Repeat(" ", longest-len(s.Key)),
				s.Key, s.Value,
			)
		}
	}

	if !v.cli.IsSet(OptDisableDeps) {
		fmt.Fprintf(w, "\n<cyan>dependencies:</>\n")
		for _, m := range v.Dependencies {
			if m.Replace != nil {
				m = *m.Replace
			}

			if m.Sum == "" {
				m.Sum = "unknown"
			}

			fmt.Fprintf(w, "  %47s :: <cyan>%s</> :: <yellow>%s</>\n", m.Sum, m.Path, m.Version)
		}
	}

	return color.Sprint(w.String())
}

// GetVersionInfo returns the version information for the CLI.
func (cli *CLI[T]) GetVersionInfo() *VersionInfo[T] {
	v := VersionInfo[T]{}

	if cli.VersionInfo != nil {
		v.Name = cli.VersionInfo.Name
		v.Version = cli.VersionInfo.Version
		v.Commit = cli.VersionInfo.Commit
		v.Date = cli.VersionInfo.Date
	}

	v.cli = cli
	v.GoVersion = runtime.Version()
	v.Command = filepath.Base(os.Args[0])
	v.OS = runtime.GOOS
	v.Arch = runtime.GOARCH
	v.Links = cli.Links

	build, ok := debug.ReadBuildInfo()
	if ok {
		if v.Settings == nil {
			v.Settings = make([]BuildSetting, 0, len(build.Settings))
			for _, setting := range build.Settings {
				v.Settings = append(v.Settings, BuildSetting{
					Key:   setting.Key,
					Value: setting.Value,
				})
			}
		}

		if v.Dependencies == nil {
			v.Dependencies = make([]Module, 0, len(build.Deps))
			for _, dep := range build.Deps {
				v.Dependencies = append(v.Dependencies, Module{
					Path:    dep.Path,
					Version: dep.Version,
					Sum:     dep.Sum,
				})
			}
		}

		if v.Name == "" {
			v.Name = build.Main.Path
		}

		if v.Version == "" {
			v.Version = build.Main.Version
		}

		if v.Commit == "" {
			v.Commit = v.GetSetting("vcs.revision", build.Main.Sum)
		}

		if v.Date == "" {
			v.Date = v.GetSetting("vcs.time", "unknown")
		}
	}

	if v.Name == "" {
		v.Name = v.Command
	}

	if v.Version == "" {
		v.Version = "unknown"
	}

	if v.Commit == "" {
		v.Commit = "unknown"
	}

	if v.Date == "" {
		v.Date = "unknown"
	}

	return &v
}
