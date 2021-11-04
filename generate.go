// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.
package goflagsmarkdown

import (
	"fmt"
	"io"
	"strings"

	flags "github.com/jessevdk/go-flags"
)

const optionHeader = "| Environment vars | Flags | Description |\n| --- | --- | --- |\n"

func Generate(parser *flags.Parser, out io.Writer) {
	generateRecursive(parser, out)

}

func generateRecursive(parser *flags.Parser, out io.Writer, groups ...*flags.Group) {
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
				environment = "N/A"
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

			description = strings.Replace(description, "|", "\\|", -1)

			fmt.Fprintf(out, "| %s | `%s` | %s |\n", environment, option.String(), description)
		}

		groups := group.Groups()
		if len(groups) > 0 {
			generateRecursive(parser, out, groups...)
		}
	}
}
