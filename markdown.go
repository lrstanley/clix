// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

package clix

import "github.com/alecthomas/kong"

// WithMarkdownPlugin adds a hidden "generate-markdown" command that allows
// generating markdown documentation for the CLI.
func WithMarkdownPlugin[T any]() Option[T] {
	return func(cli *CLI[T]) {
		cli.kongOptions = append(cli.kongOptions, kong.DynamicCommand(
			"generate-markdown",
			"generate markdown documentation and write to stdout",
			"",
			&MarkdownCommand{},
			"hidden",
		))
	}
}

type MarkdownCommand struct {
	OutputPath    string `short:"o" name:"output-path" type:"path" default:"-" help:"path to write the markdown to, or '-' to write to stdout"`
	TemplatePath  string `short:"t" name:"template-path" type:"path" help:"optional path to a template file to use for the markdown"`
	IncludeHidden bool   `name:"include-hidden" help:"include hidden flags and commands"`
}

func (m *MarkdownCommand) BeforeApply(ctx *kong.Context) error {
	panic("not implemented")
}

// const optionHeader = "| Environment vars | Flags | Type | Description |\n| --- | --- | --- | --- |\n"

// // Markdown writes generated marakdown to the provided io.Writer.
// func (cli *CLI[T]) Markdown(out io.Writer) {
// 	cli.generateRecursive(out)
// }

// func (cli *CLI[T]) generateRecursive(out io.Writer, groups ...*flags.Group) {
// 	// TODO: commands?

// 	parser := cli.newParser()

// 	if groups == nil {
// 		groups = parser.Groups()
// 	}

// 	for _, group := range groups {
// 		if group.LongDescription != "" {
// 			fmt.Fprintf(out, "\n#### %s\n%s", group.LongDescription, optionHeader)
// 		} else if group.ShortDescription != "" {
// 			fmt.Fprintf(out, "\n#### %s\n%s", group.ShortDescription, optionHeader)
// 		}

// 		// print the options in this group first, then recursively continue into
// 		// each sub-group.
// 		options := group.Options()
// 		for _, option := range options {
// 			if option.Hidden {
// 				continue
// 			}

// 			environment := option.EnvKeyWithNamespace()
// 			if environment != "" {
// 				environment = "`" + environment + "`"
// 			} else {
// 				environment = "-"
// 			}

// 			description := option.Description

// 			if option.Required {
// 				description += " [**required**]"
// 			}

// 			if option.Default != nil {
// 				description += fmt.Sprintf(" [**default: %s**]", strings.Join(option.Default, ", "))
// 			}

// 			if option.Choices != nil {
// 				description += fmt.Sprintf(" [**choices: %s**]", strings.Join(option.Choices, ", "))
// 			}

// 			description = strings.ReplaceAll(description, "|", "\\|")

// 			_type := fmt.Sprintf("%T", option.Value())
// 			if strings.Contains(strings.ToLower(_type), "func") {
// 				_type = "-"
// 			}

// 			fmt.Fprintf(out, "| %s | `%s` | %s | %s |\n", environment, option.String(), _type, description)
// 		}

// 		children := group.Groups()
// 		if len(children) > 0 {
// 			cli.generateRecursive(out, children...)
// 		}
// 	}
// }
