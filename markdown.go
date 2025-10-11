// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

package clix

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sync/atomic"
	"text/template"

	"github.com/alecthomas/kong"
)

var markdownShouldExit = true

// WithMarkdownPlugin adds a hidden "generate-markdown" command that allows
// generating markdown documentation for the CLI. To make it so this command
// ignores any other required flags, it's invoked before kong applies additional
// restrictions, which means it does not support special flags. To adjust the
// behavior, you can use environment variables:
//
//   - CLIX_TEMPLATE_PATH: optional path to a directory containing template files to
//     use for the markdown.
//   - CLIX_OUTPUT_PATH: path to write the markdown to, or '-' to write to stdout
//     (defaults to stdout).
func WithMarkdownPlugin[T any]() Option[T] {
	var initialized atomic.Bool

	return func(cli *CLI[T]) {
		if initialized.Load() {
			return
		}

		cmd := &MarkdownCommand{}

		cli.kongOptions = append(
			cli.kongOptions, kong.DynamicCommand(
				"generate-markdown",
				"generate markdown documentation and write to stdout",
				"",
				cmd,
				"hidden",
			),
		)
	}
}

type MarkdownCommand struct{}

func (m *MarkdownCommand) BeforeReset(
	ctx *kong.Kong,
	appInfo *AppInfo,
	version *Version,
) error {
	var output string
	var err error

	if v := os.Getenv("CLIX_TEMPLATE_PATH"); v == "" {
		output, err = m.GenerateMarkdown(ctx.Model, templates, appInfo, version)
	} else {
		var files []string
		err = filepath.Walk(v, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}
			files = append(files, path)
			return nil
		})
		if err != nil {
			return fmt.Errorf("failed to walk template directory: %w", err)
		}

		var tmpl *template.Template
		tmpl, err = templates.ParseFiles(files...)
		if err != nil {
			return fmt.Errorf("failed to parse templates: %w", err)
		}

		output, err = m.GenerateMarkdown(ctx.Model, tmpl, appInfo, version)
	}

	if err != nil {
		return fmt.Errorf("failed to generate markdown: %w", err)
	}

	if v := os.Getenv("CLIX_OUTPUT_PATH"); v == "-" || v == "" {
		fmt.Fprint(os.Stdout, output) //nolint:forbidigo
	} else {
		err = os.WriteFile(v, []byte(output), 0o600)
		if err != nil {
			return err
		}
	}

	if markdownShouldExit {
		os.Exit(0)
	}
	return nil
}

// GenerateMarkdown generates the markdown documentation for the CLI, returning the
// markdown as a string.
func (m *MarkdownCommand) GenerateMarkdown(
	model *kong.Application,
	tmpl *template.Template,
	appInfo *AppInfo,
	version *Version,
) (string, error) {
	buf := bytes.NewBuffer(nil)

	err := tmpl.ExecuteTemplate(buf, "main.gotmpl", map[string]any{
		"Model":   model,
		"AppInfo": appInfo,
		"Config":  m,
		"Version": version,
	})
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
