// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package clix

import (
	"fmt"
	"io"
	"strings"

	flags "github.com/jessevdk/go-flags"
)

const optionHeader = "| Environment vars | Flags | Type | Description |\n| --- | --- | --- | --- |\n"

// GenerateMarkdown writes generated marakdown to the provided io.Writer.
func GenerateMarkdown(parser *flags.Parser, out io.Writer) {
	generateMarkdownRecursive(parser, out)
}

func generateMarkdownRecursive(parser *flags.Parser, out io.Writer, groups ...*flags.Group) {
	// TODO: commands?

	if groups == nil {
		groups = parser.Groups()
	}

	for _, group := range groups {
		if group.LongDescription != "" {
			fmt.Fprintf(out, "\n#### %s\n%s", group.LongDescription, optionHeader)
		} else if group.ShortDescription != "" {
			fmt.Fprintf(out, "\n#### %s\n%s", group.ShortDescription, optionHeader)
		}

		// print the options in this group first, then recursively continue into
		// each sub-group.
		options := group.Options()
		for _, option := range options {
			if option.Hidden {
				continue
			}

			environment := option.EnvKeyWithNamespace()
			if environment != "" {
				environment = "`" + environment + "`"
			} else {
				environment = "-"
			}

			description := option.Description

			if option.Required {
				description += " [**required**]"
			}

			if option.Default != nil {
				description += fmt.Sprintf(" [**default: %s**]", strings.Join(option.Default, ", "))
			}

			if option.Choices != nil {
				description += fmt.Sprintf(" [**choices: %s**]", strings.Join(option.Choices, ", "))
			}

			description = strings.ReplaceAll(description, "|", "\\|")

			_type := fmt.Sprintf("%T", option.Value())
			if strings.Contains(strings.ToLower(_type), "func") {
				_type = "-"
			}

			fmt.Fprintf(out, "| %s | `%s` | %s | %s |\n", environment, option.String(), _type, description)
		}

		groups := group.Groups()
		if len(groups) > 0 {
			generateMarkdownRecursive(parser, out, groups...)
		}
	}
}
