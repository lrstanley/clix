// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

package dotenv

import (
	"regexp"

	"github.com/lrstanley/clix/v2/internal/dotenv/lexer"
)

var expandVarRegex = regexp.MustCompile(`(\\)?(\$)(\()?\{?([a-zA-Z_][a-zA-Z0-9_]+)?\}?`)

// ExpandVariables resolves ${VAR} references in the variables parsed by the Parser.
// ${VAR} format can be replaced in combination with any other mixed text, however,
// $VAR will only be replace if that's the entire value, e.g. "FOO=$VAR". Additionally,
// $VAR/${VAR} references in single quotes will not be replaced, which includes
// triple-single quote blocks.
func (p *Parser) ExpandVariables(maxDepth int, includeVars map[string]string) { //nolint:gocognit
	maxDepth = max(1, maxDepth)
	if includeVars == nil {
		includeVars = make(map[string]string)
	}

	// Loop over the variables, replacing ${VAR} references with the actual value of the
	// variable in question, until we have no more replacements.
	attempts := 0
	for {
		attempts++

		if attempts > maxDepth {
			return
		}

		hadChanges := false

		for key := range p.vars {
			if p.quoteTypes[key] == lexer.QuoteTypeSingle {
				continue
			}
			p.vars[key] = expandVarRegex.ReplaceAllStringFunc(p.vars[key], func(s string) string {
				submatch := expandVarRegex.FindStringSubmatch(s)
				if submatch == nil {
					return s
				}
				// Don't do anything with escaped references.
				if submatch[1] == "\\" || submatch[2] == "(" {
					return s
				}
				if submatch[4] != "" {
					// If $VAR is used, vs ${VAR}, only replace if that's the entire
					// value.
					if s[1] != '{' && s[len(s)-1] != '}' && s != p.vars[key] {
						return s
					}

					hadChanges = true
					if val, ok := p.vars[submatch[4]]; ok {
						return val
					}
					if val, ok := includeVars[submatch[4]]; ok {
						return val
					}
					return ""
				}
				return s
			})
		}

		if !hadChanges {
			break
		}
	}

	// Unescape any escaped ${VAR} references.
	for key := range p.vars {
		p.vars[key] = expandVarRegex.ReplaceAllStringFunc(p.vars[key], func(s string) string {
			submatch := expandVarRegex.FindStringSubmatch(s)
			if submatch == nil {
				return s
			}
			if submatch[1] == "\\" || submatch[2] == "(" {
				return submatch[0][1:]
			}
			return s
		})
	}
}
