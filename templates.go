// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

package clix

import (
	"embed"
	"fmt"
	"reflect"
	"regexp"
	"slices"
	"strings"
	"text/template"

	"github.com/alecthomas/kong"
)

var (
	templatePaths = []string{
		"templates/*.gotmpl",
		"templates/helpers/*.gotmpl",
	}

	//go:embed templates
	templateDir embed.FS
	templates   = template.Must(
		template.New("").
			Funcs(tmplFuncMap).
			ParseFS(
				templateDir,
				templatePaths...,
			),
	)
	tmplFuncMap = template.FuncMap{
		"children_by_type": func(node *kong.Node, typ string) []*kong.Node {
			var results []*kong.Node
			for _, child := range node.Children {
				switch typ {
				case "command":
					if child.Type == kong.CommandNode {
						results = append(results, child)
					}
				case "argument":
					if child.Type == kong.ArgumentNode {
						results = append(results, child)
					}
				case "application":
					if child.Type == kong.ApplicationNode {
						results = append(results, child)
					}
				default:
					panic(fmt.Sprintf("unknown node type: %s", typ))
				}
			}
			return results
		},
		"node_is_type": func(node *kong.Node, typ string) bool {
			switch typ {
			case "command":
				return node.Type == kong.CommandNode
			case "argument":
				return node.Type == kong.ArgumentNode
			case "application":
				return node.Type == kong.ApplicationNode
			default:
				panic(fmt.Sprintf("unknown node type: %s", typ))
			}
		},
		"node_has_unhidden_children": func(node *kong.Node) bool {
			for _, child := range node.Children {
				if !child.Hidden {
					return true
				}
			}
			return false
		},
		"flag_groups": func(flags []*kong.Flag) []*kong.Group {
			var results []*kong.Group
			for _, flag := range flags {
				if flag.Hidden || flag.Group == nil {
					continue
				}
				if !slices.ContainsFunc(results, func(g *kong.Group) bool { return g.Key == flag.Group.Key }) {
					results = append(results, flag.Group)
				}
			}
			return results
		},
		"flags_by_group": func(flags []*kong.Flag, groupKey string) []*kong.Flag {
			var results []*kong.Flag
			for _, flag := range flags {
				if flag.Hidden {
					continue
				}
				if groupKey == "" && flag.Group == nil {
					results = append(results, flag)
				}
				if flag.Group != nil && (flag.Group.Key == groupKey || flag.Group.Title == groupKey) {
					results = append(results, flag)
				}
			}
			return results
		},
		"slug": func(input any) string {
			re := regexp.MustCompile(`[^a-zA-Z0-9_-]+`)
			slug := fmt.Sprintf("%v", input)
			slug = strings.ToLower(slug)
			slug = strings.TrimSpace(slug)
			slug = re.ReplaceAllString(slug, "-")
			slug = strings.ReplaceAll(slug, "--", "-")
			slug = strings.Trim(slug, "-")
			return slug
		},
		"sanitize_md": sanizeMarkdown,
		"quote_code": func(input any) any {
			return applyStringFn(input, func(s string) string {
				if s == "" {
					return "-"
				}
				return "`" + strings.ReplaceAll(s, "`", "\\`") + "`"
			})
		},
		"bold": func(input any) any {
			return applyStringFn(input, func(s string) string {
				if s == "" {
					return "-"
				}
				return "**" + strings.ReplaceAll(s, "**", "\\*") + "**"
			})
		},
		"italic": func(input any) any {
			return applyStringFn(input, func(s string) string {
				if s == "" {
					return "-"
				}
				return "*" + strings.ReplaceAll(s, "*", "\\*") + "*"
			})
		},
		"add_int": func(a, b int) int {
			return a + b
		},
		"mult_int": func(a, b int) int {
			return a * b
		},
		"join":  strings.Join,
		"trim":  strings.TrimSpace,
		"upper": strings.ToUpper,
		"lower": strings.ToLower,
		"title": strings.Title,
		// Switch order so that "foo" | repeat 5
		"repeat": func(count int, str string) string { return strings.Repeat(str, count) },
		"dict": func(args ...any) map[string]any {
			dict := map[string]any{}
			lenv := len(args)
			if lenv%2 != 0 {
				panic("dict must have an even number of arguments")
			}
			for i := 0; i < lenv; i += 2 {
				key := fmt.Sprintf("%v", args[i])
				if i+1 >= lenv {
					dict[key] = ""
					continue
				}
				dict[key] = args[i+1]
			}
			return dict
		},
		"slice": func(args ...any) []any {
			return args
		},
		"bool_or": func(input, v1, v2 any) any {
			if reflectBool(input) {
				return v1
			}
			return v2
		},
		"slice_append": func(slice []any, v any) []any {
			tp := reflect.TypeOf(slice).Kind()
			switch tp { //nolint:exhaustive
			case reflect.Slice, reflect.Array:
				l2 := reflect.ValueOf(slice)
				l := l2.Len()
				nl := make([]any, l)
				for i := range l {
					nl[i] = l2.Index(i).Interface()
				}
				return append(nl, v) //nolint:makezero
			default:
				panic(fmt.Sprintf("cannot push on type %s", tp))
			}
		},
		"table": func(rowsRaw []any, sanitize bool) string {
			rows := make([][]string, len(rowsRaw))
			for i, row := range rowsRaw {
				v := row.([]any) //nolint:errcheck
				rows[i] = make([]string, len(v))
				for j, vv := range v {
					if sanitize {
						rows[i][j] = sanizeMarkdown(fmt.Sprintf("%v", vv))
					} else {
						rows[i][j] = fmt.Sprintf("%v", vv)
					}
				}
			}

			if len(rows) < 2 {
				panic("table must have at least 2 rows")
			}

			// Make sure all rows have the same number of columns.
			for _, row := range rows {
				if len(row) != len(rows[0]) {
					panic("all rows must have the same number of columns")
				}
			}

			colSizes := make([]int, len(rows[0]))
			for _, row := range rows {
				for i, cell := range row {
					colSizes[i] = max(colSizes[i], len(cell))
				}
			}

			headers := rows[0]
			rows = rows[1:]

			var buf strings.Builder
			for i, header := range headers {
				buf.WriteString("| ")
				buf.WriteString(header)
				buf.WriteString(strings.Repeat(" ", max(0, colSizes[i]-len(header)+1)))
			}
			buf.WriteString("|")
			buf.WriteString("\n")

			for i := range headers {
				buf.WriteString("|")
				buf.WriteString(strings.Repeat("-", max(0, colSizes[i]+2)))
			}
			buf.WriteString("|")
			buf.WriteString("\n")

			for _, row := range rows {
				for i, cell := range row {
					buf.WriteString("| ")
					buf.WriteString(cell)
					buf.WriteString(strings.Repeat(" ", max(0, colSizes[i]-len(cell)+1)))
				}
				buf.WriteString("|")
				buf.WriteString("\n")
			}
			return buf.String()
		},
	}
)

func reflectBool(input any) bool {
	rv := reflect.ValueOf(input)

	if !rv.IsValid() || rv.IsZero() {
		return false
	}

	switch rv.Kind() { //nolint:exhaustive
	case reflect.Slice, reflect.Array:
		if rv.Cap() == 0 || rv.Len() == 0 {
			return false
		}

		// kong can sometimes set the first value in an array to an empty string.
		if rv.Len() == 1 && rv.Index(0).String() == "" {
			return false
		}

		return true
	case reflect.Map:
		return rv.Len() > 0
	case reflect.Bool:
		return rv.Bool()
	case reflect.String:
		return rv.String() != ""
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return rv.Int() != 0
	case reflect.Float32, reflect.Float64:
		return rv.Float() != 0
	case reflect.Pointer:
		return reflectBool(rv.Elem().Interface())
	default:
		panic(fmt.Sprintf("cannot reflect bool on type %s", rv.Kind()))
	}
}

func applyStringFn(input any, fn func(string) string) any {
	// Use reflection to see if it's a slice. If it is, do it for each item,
	// otherwise, do it for the single item.
	rv := reflect.ValueOf(input)
	if rv.Kind() == reflect.Slice || rv.Kind() == reflect.Array {
		// Handle slice/array - quote each item.
		var results []string
		for i := range rv.Len() {
			results = append(results, fn(fmt.Sprintf("%v", rv.Index(i).Interface())))
		}
		return results
	}
	// Handle single item.
	return fn(fmt.Sprintf("%v", input))
}

func sanizeMarkdown(s any) string {
	replacer := strings.NewReplacer(
		`\`, `\\`,
		`|`, `\|`,
		`*`, `\*`,
		`_`, `\_`,
		`<`, `\<`,
		`>`, `\>`,
		`#`, `\#`,
		`[`, `\[`,
		`]`, `\]`,
		`(`, `\(`,
		`)`, `\)`,
		`+`, `\+`,
		`-`, `\-`,
		`*`, `\*`,
		`#`, `\#`,
		`{`, `\{`,
		`}`, `\}`,
		`!`, `\!`,
		`?`, `\?`,
		"`", `\`+"`",
		"\n", `<br>`,
	)
	return replacer.Replace(fmt.Sprintf("%v", s))
}
