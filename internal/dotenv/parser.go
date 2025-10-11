// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

package dotenv

import (
	"fmt"
	"maps"
	"os"
	"slices"
	"strings"

	"github.com/lrstanley/clix/v2/internal/dotenv/lexer"
)

const defaultResolveMaxDepth = 20

// Parser is a parser for dotenv files.
type Parser struct {
	refs       []*lexer.Reference
	pos        int
	vars       map[string]string
	quoteTypes map[string]lexer.QuoteType
}

// New creates a new Parser.
func New() *Parser {
	return &Parser{
		refs:       make([]*lexer.Reference, 0),
		vars:       make(map[string]string),
		quoteTypes: make(map[string]lexer.QuoteType),
	}
}

// Parse parses the provided value and stores the results in the Parser.
func (p *Parser) Parse(value string) error {
	p.refs = nil
	p.pos = 0

	lex := lexer.New(value)

	for ref, err := range lex.Iter() {
		if err != nil {
			return err
		}
		p.refs = append(p.refs, ref)
	}

	if len(p.refs) == 0 {
		return nil
	}

	for {
		p.skip(lexer.Whitespace, lexer.Newline, lexer.Comment, lexer.Export)
		r := p.next()

		if r.Token == lexer.EOF {
			break
		}

		if r.Token != lexer.Key {
			return &ParseError{
				Content: value,
				Line:    r.Line,
				Column:  r.Column,
				Err:     fmt.Errorf("expected KEY, got %s (%q)", r.Token, r.Value),
			}
		}

		key := r.Value

		p.skip(lexer.Whitespace)

		r = p.next()
		if r.Token != lexer.Equals {
			return &ParseError{
				Content: value,
				Line:    r.Line,
				Column:  r.Column,
				Err:     fmt.Errorf("expected '=', got %s (%q)", r.Token, r.Value),
			}
		}

		p.skip(lexer.Whitespace)

		r = p.next()

		switch r.Token { //nolint:exhaustive
		case lexer.Newline, lexer.EOF:
			p.vars[key] = ""
			p.quoteTypes[key] = lexer.QuoteTypeNone
		case lexer.Value:
			p.vars[key] = r.Value
			p.quoteTypes[key] = r.QuoteType
		default:
			return &ParseError{
				Content: value,
				Line:    r.Line,
				Column:  r.Column,
				Err:     fmt.Errorf("expected VALUE, got %s (%q)", r.Token, r.Value),
			}
		}
	}

	return nil
}

// Values returns a copy of the variables parsed by the Parser.
func (p *Parser) Values() map[string]string {
	return maps.Clone(p.vars)
}

// next returns the next reference from the Parser.
func (p *Parser) next() *lexer.Reference {
	if p.pos >= len(p.refs) {
		return &lexer.Reference{
			Token:    lexer.EOF,
			Line:     p.refs[p.pos-1].Line,
			Column:   p.refs[p.pos-1].Column,
			Position: p.refs[p.pos-1].Position,
		}
	}

	ref := p.refs[p.pos]
	p.pos++
	return ref
}

// skip skips the next reference from the Parser if it matches any of the provided
// tokens.
func (p *Parser) skip(ts ...lexer.Token) {
	if len(ts) == 0 {
		return
	}

	for p.pos < len(p.refs) {
		if !slices.Contains(ts, p.refs[p.pos].Token) {
			break
		}
		p.pos++
	}
}

// ParseFiles parses the provided files and returns the variables parsed.
func ParseFiles(paths ...string) (map[string]string, error) {
	parser := New()

	for _, path := range paths {
		content, err := os.ReadFile(path)
		if err != nil {
			return nil, &FileAccessError{
				Path: path,
				Err:  err,
			}
		}

		err = parser.Parse(string(content))
		if err != nil {
			return nil, err
		}
	}

	envVars := make(map[string]string)
	for _, envVar := range os.Environ() {
		parts := strings.SplitN(envVar, "=", 2)
		envVars[parts[0]] = parts[1]
	}
	parser.ExpandVariables(defaultResolveMaxDepth, envVars)

	return parser.Values(), nil
}

// ParseStrings parses the provided strings and returns the variables parsed.
func ParseStrings(values ...string) (map[string]string, error) {
	parser := New()

	for _, value := range values {
		err := parser.Parse(value)
		if err != nil {
			return nil, err
		}
	}

	envVars := make(map[string]string)
	for _, envVar := range os.Environ() {
		parts := strings.SplitN(envVar, "=", 2)
		envVars[parts[0]] = parts[1]
	}
	parser.ExpandVariables(defaultResolveMaxDepth, envVars)

	return parser.Values(), nil
}
