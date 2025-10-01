// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

package clix

import (
	"bytes"
	"fmt"
	"os"
	"slices"
	"text/template"

	"github.com/alecthomas/kong"
)

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
	return func(cli *CLI[T]) {
		if cli.checkAlreadyInit("markdown") {
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
			// Do this to bypass all required flags that might be set by the end
			// application.
			kong.PostBuild(func(ctx *kong.Kong) error {
				if !slices.Contains(os.Args, "generate-markdown") {
					return nil
				}

				err := cmd.BeforeApply(ctx, cli.app, cli.version)
				if err != nil {
					return err
				}
				os.Exit(0)
				return nil
			}),
		)
	}
}

type MarkdownCommand struct{}

func (m *MarkdownCommand) BeforeApply(
	ctx *kong.Kong,
	appInfo *AppInfo,
	version *Version,
) error {
	var output string
	var err error

	if v := os.Getenv("CLIX_TEMPLATE_PATH"); v == "" {
		output, err = m.GenerateMarkdown(ctx.Model, templates, appInfo, version)
	} else {
		var tmpl *template.Template
		tmpl, err = template.New("").
			Funcs(tmplFuncMap).
			ParseFS(os.DirFS(v), templatePaths...)
		if err != nil {
			return err
		}

		output, err = m.GenerateMarkdown(ctx.Model, tmpl, appInfo, version)
	}

	if err != nil {
		return err
	}

	if v := os.Getenv("CLIX_OUTPUT_PATH"); v == "-" || v == "" {
		fmt.Fprint(os.Stdout, output) //nolint:forbidigo
	} else {
		err = os.WriteFile(v, []byte(output), 0o600)
		if err != nil {
			return err
		}
	}

	os.Exit(0)
	return nil
}

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
