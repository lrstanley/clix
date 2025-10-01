// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

package clix

import (
	"fmt"
	"strings"
	"sync/atomic"

	"github.com/alecthomas/kong"
)

// WithAppInfo adds the provided application information to the CLI.
func WithAppInfo[T any](app AppInfo) Option[T] {
	var initialized atomic.Bool
	return func(cli *CLI[T]) {
		if initialized.Swap(true) {
			return
		}

		if app.Name != "" {
			cli.kongOptions = append(cli.kongOptions, kong.Name(app.Name))
		}

		if app.Description != "" {
			cli.kongOptions = append(cli.kongOptions, kong.Description(app.Description))
		}

		cli.app.Inherit(app)
	}
}

type AppInfo struct {
	Name        string `json:"name"`          // Application name. Defaults to the main module path.
	Description string `json:"description"`   // Application description.
	Version     string `json:"build_version"` // Build version. Uses VCS info if available.
	Commit      string `json:"build_commit"`  // VCS commit SHA. Uses VCS info if available.
	Date        string `json:"build_date"`    // VCS commit date. Uses VCS info if available.

	Links []Link `json:"links,omitempty"` // Links to the project's website, support, issues, security, etc.
}

// Inherit inherits the values from the provided app.
func (a *AppInfo) Inherit(app AppInfo) {
	if app.Name != "" {
		a.Name = app.Name
	}
	if app.Description != "" {
		a.Description = app.Description
	}
	if app.Version != "" {
		a.Version = app.Version
	}
	if app.Commit != "" {
		a.Commit = app.Commit
	}
	if app.Date != "" {
		a.Date = app.Date
	}
	if len(app.Links) > 0 {
		a.Links = app.Links
	}
}

// IsModuleName returns true if the name looks like a module path-style format,
// otherwise false.
func (a *AppInfo) IsModuleName() bool {
	return strings.Contains(a.Name, "/")
}

// Link allows you to define a link to be included in the version and usage
// output.
type Link struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// GithubLinks return an opinonated set of links for the project, using
// common Github layout conventions.
func GithubLinks(repo, branch, homepage string) []Link {
	repo = strings.TrimPrefix(repo, "https://")
	repo = strings.TrimSuffix(repo, "/")

	links := []Link{}

	if branch == "" {
		branch = "master"
	}

	if homepage != "" {
		links = append(links, Link{
			Name: "homepage",
			URL:  homepage,
		})
	}

	links = append(links, []Link{
		{Name: "github", URL: "https://" + repo},
		{Name: "issues", URL: fmt.Sprintf("https://%s/issues/new/choose", repo)},
		{Name: "support", URL: fmt.Sprintf("https://%s/blob/%s/.github/SUPPORT.md", repo, branch)},
		{Name: "contributing", URL: fmt.Sprintf("https://%s/blob/%s/.github/CONTRIBUTING.md", repo, branch)},
		{Name: "security", URL: fmt.Sprintf("https://%s/security/policy", repo)},
	}...)

	return links
}
