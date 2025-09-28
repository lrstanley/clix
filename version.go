// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

package clix

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"strings"
)

// WithVersionPlugin adds the version plugin to the CLI. This includes flags
// for --version, --version-json, etc.
func WithVersionPlugin[T any]() Option[T] {
	return func(cli *CLI[T]) {
		if cli.checkAlreadyInit("version") {
			return
		}
		cli.Plugins = append(cli.Plugins, &VersionPlugin{})
	}
}

// GetVersion returns the version information for the CLI, which will be populated
// after parsing.
func (cli *CLI[T]) GetVersion() *Version {
	return cli.version
}

type VersionPlugin struct {
	Version     VersionFlag     `short:"v" name:"version" help:"prints version information and exits"`
	VersionJSON VersionJSONFlag `name:"version-json" help:"prints version information in JSON format and exits"`
}

type VersionFlag bool

func (v VersionFlag) AfterApply(ver *Version) error {
	if v {
		fmt.Println(ver.String()) //nolint:forbidigo
		os.Exit(0)
	}
	return nil
}

type VersionJSONFlag bool

func (v VersionJSONFlag) BeforeApply(ver *Version) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "    ")
	if err := enc.Encode(ver); err != nil {
		panic(err)
	}
	os.Exit(0)
	return nil
}

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

// Version represents the version information for the CLI.
type Version struct {
	AppInfo      *AppInfo       `json:"app_info,omitempty"`       // Application information.
	Settings     []BuildSetting `json:"build_settings,omitempty"` // Other information about the build.
	Dependencies []Module       `json:"dependencies,omitempty"`   // Module dependencies.

	Command   string `json:"command"`    // Executable name where the command was called from.
	GoVersion string `json:"go_version"` // Version of Go that produced this binary.
	OS        string `json:"os"`         // Operating system for this build.
	Arch      string `json:"arch"`       // CPU Architecture for build build.
}

// NonSensitiveVersion represents the version information for the CLI.
type NonSensitiveVersion struct {
	AppInfo   *AppInfo `json:"app_info,omitempty"` // Application information.
	Command   string   `json:"command"`            // Executable name where the command was called from.
	GoVersion string   `json:"go_version"`         // Version of Go that produced this binary.
	OS        string   `json:"os"`                 // Operating system for this build.
	Arch      string   `json:"arch"`               // CPU Architecture for build build.
}

// NonSensitive returns a copy of Version with sensitive information removed.
func (v *Version) NonSensitive() *NonSensitiveVersion {
	return &NonSensitiveVersion{
		AppInfo:   v.AppInfo,
		Command:   v.Command,
		GoVersion: v.GoVersion,
		OS:        v.OS,
		Arch:      v.Arch,
	}
}

// GetSetting returns the value of the setting with the given key, otherwise
// defaults to defaultValue.
func (v *Version) GetSetting(key, defaultValue string) string {
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

func (v *Version) stringBase() string {
	w := &bytes.Buffer{}

	fmt.Fprintf(w, "%s :: %s\n", v.AppInfo.Name, v.AppInfo.Version)
	fmt.Fprintf(w, "|  build commit :: %s\n", v.AppInfo.Commit)
	fmt.Fprintf(w, "|    build date :: %s\n", v.AppInfo.Date)
	fmt.Fprintf(w, "|    go version :: %s %s/%s\n", v.GoVersion, v.OS, v.Arch)

	if len(v.AppInfo.Links) > 0 {
		var longest int
		for _, l := range v.AppInfo.Links {
			if len(l.Name) > longest {
				longest = len(l.Name)
			}
		}

		fmt.Fprintf(w, "\nhelpful links:\n")
		for _, l := range v.AppInfo.Links {
			fmt.Fprintf(
				w, "|  %s%s :: %s\n",
				strings.Repeat(" ", longest-len(l.Name)),
				l.Name, l.URL,
			)
		}
	}

	return w.String()
}

func (v *Version) String() string {
	w := &bytes.Buffer{}

	w.WriteString(v.stringBase())

	var longest int
	for _, s := range v.Settings {
		if len(s.Key) > longest {
			longest = len(s.Key)
		}
	}

	fmt.Fprintf(w, "\nbuild options:\n")
	for _, s := range v.Settings {
		fmt.Fprintf(
			w, "|  %s%s :: %s\n",
			strings.Repeat(" ", longest-len(s.Key)),
			s.Key, s.Value,
		)
	}

	fmt.Fprintf(w, "\ndependencies:\n")
	for _, m := range v.Dependencies {
		if m.Replace != nil {
			m = *m.Replace
		}

		if m.Sum == "" {
			m.Sum = "unknown" //nolint:goconst
		}

		fmt.Fprintf(w, "  %47s :: %s :: %s\n", m.Sum, m.Path, m.Version)
	}

	return w.String()
}

// GetVersionInfo returns the version information for the CLI.
func GetVersionInfo(app *AppInfo) *Version {
	v := &Version{
		AppInfo:   app,
		GoVersion: runtime.Version(),
		Command:   filepath.Base(os.Args[0]),
		OS:        runtime.GOOS,
		Arch:      runtime.GOARCH,
	}

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

		if v.AppInfo.Name == "" {
			v.AppInfo.Name = build.Main.Path
		}

		if v.AppInfo.Version == "" {
			v.AppInfo.Version = build.Main.Version
		}

		if v.AppInfo.Commit == "" {
			v.AppInfo.Commit = v.GetSetting("vcs.revision", build.Main.Sum)
		}

		if v.AppInfo.Date == "" {
			v.AppInfo.Date = v.GetSetting("vcs.time", "unknown")
		}
	}

	if v.AppInfo.Name == "" {
		v.AppInfo.Name = v.Command
	}

	if v.AppInfo.Version == "" {
		v.AppInfo.Version = "unknown"
	}

	if v.AppInfo.Commit == "" {
		v.AppInfo.Commit = "unknown"
	}

	if v.AppInfo.Date == "" {
		v.AppInfo.Date = "unknown"
	}

	if v.AppInfo.Description == "" {
		// TODO: https://github.com/alecthomas/kong/issues/376
		v.AppInfo.Description = strings.ReplaceAll(v.stringBase(), "|", "")
	}

	return v
}
