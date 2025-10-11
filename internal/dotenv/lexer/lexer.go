// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

package lexer

import (
	"errors"
	"fmt"
	"iter"
	"strings"
	"unicode/utf8"
)

// snapshot represents a snapshot of the lexer's state that can be restored later.
type snapshot struct {
	start        int
	pos          int
	width        int
	line         int
	prevLineCols int
	col          int
}

// Token represents a token ("identifier") in the input.
type Token string

type QuoteType int

// Reference represents a reference to a token in the input, and the associated
// context.
type Reference struct {
	Token     Token     // Token that represents the literal.
	Position  int       // Character position of the token in the input.
	Line      int       // Line number of the token in the input.
	Column    int       // Column number of the token in the input.
	Value     string    // Literal value of the token.
	QuoteType QuoteType // Type of quote used, if any.
}

const (
	EOF        Token = "EOF"
	Whitespace Token = "WHITESPACE"
	Newline    Token = "NEWLINE"
	Comment    Token = "COMMENT"
	Export     Token = "EXPORT"
	Key        Token = "KEY"
	Equals     Token = "EQUALS"
	Value      Token = "VALUE"

	eof = -1

	rawCommentStartHash = '#'
	rawNewline          = '\n'
	rawSeparator        = '='
	rawExport           = "export "
)

const (
	QuoteTypeNone = iota
	QuoteTypeSingle
	QuoteTypeDouble
)

func (t Token) String() string {
	return string(t)
}

// StatefulLexer is a lexer that maintains state between calls to [Iter], and only
// applies specific idents when applicable. We keep this state/context validation
// mostly because dot env files are pretty loose, and can be written in a variety
// of ways with very little strictness.
type StatefulLexer struct {
	input string // Input string to lex.

	// Current state.
	line         int // Current line number.
	prevLineCols int // Previous line column count.
	col          int // Current column number.
	start        int // Start position of the current token.
	pos          int // Current position in the input.
	width        int // Width of the current token.

	// Historical state.
	lastToken Token // Last token read, ignoring whitespace.
}

// New creates a new StatefulLexer. Note that all windows-style carriage returns
// are converted to newlines.
func New(input string) *StatefulLexer {
	return &StatefulLexer{
		input:     strings.ReplaceAll(input, "\r\n", "\n"),
		line:      1,
		col:       0,
		pos:       0,
		start:     0,
		width:     0,
		lastToken: EOF,
	}
}

// snapshot creates a snapshot of the lexer's state that can be restored later.
func (l *StatefulLexer) snapshot() snapshot {
	return snapshot{
		start:        l.start,
		pos:          l.pos,
		width:        l.width,
		line:         l.line,
		prevLineCols: l.prevLineCols,
		col:          l.col,
	}
}

// restore restores the lexer's state to the snapshot.
func (l *StatefulLexer) restore(s snapshot) {
	l.start = s.start
	l.pos = s.pos
	l.width = s.width
	l.line = s.line
	l.prevLineCols = s.prevLineCols
	l.col = s.col
}

// read reads the next rune from the buffered reader. Returns eof if an error
// occurs (or io.EOF is returned).
func (l *StatefulLexer) next() rune {
	if l.pos >= len(l.input) {
		l.width = 0
		return eof
	}

	r, w := utf8.DecodeRuneInString(l.input[l.pos:])
	l.width = w
	l.pos += l.width

	if r == rawNewline {
		l.line++
		l.prevLineCols = l.col
		l.col = 0
	} else {
		l.col++
	}

	return r
}

// backup steps back one rune. Should only be called once per call of next.
func (l *StatefulLexer) backup() {
	l.pos -= l.width

	if l.col == 0 {
		l.col = l.prevLineCols
		l.line--
	} else {
		l.col--
	}
}

// peek steps forward one rune, reads, and backs up again.
func (l *StatefulLexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

func (l *StatefulLexer) mark() {
	l.start = l.pos
}

// scan continues to read forward as long as all of the provided functions continue
// to return true. It returns when one of the functions returns false (and doesn't
// read the rune that caused the function to return false).
func (l *StatefulLexer) scan(fns ...func(r rune) bool) {
	if len(fns) == 0 {
		return
	}

	for {
		found := false
		for _, fn := range fns {
			if fn(l.peek()) {
				l.next()
				found = true
				break
			}
		}
		if !found {
			return
		}
	}
}

// scanUntil scans until one of the functions returns true. It returns true if one
// of the functions returns true (not reading that rune), false otherwise.
func (l *StatefulLexer) scanUntil(fns ...func(r rune) bool) {
	if len(fns) == 0 {
		return
	}

	for {
		r := l.peek()
		if r == eof {
			return
		}
		for _, fn := range fns {
			if fn(r) {
				return
			}
		}
		l.next()
	}
}

// scanExact scans for an exact string. It returns true if the string is found (and
// keeps its position), false otherwise (and restores the snapshot to the original
// position).
func (l *StatefulLexer) scanExact(v string) bool {
	s := l.snapshot()
	for _, char := range v {
		if r := l.peek(); r != char {
			// Back to original state.
			l.restore(s)
			return false
		}
		l.next()
	}
	return true
}

// value returns the current value of the lexer, between the last marked position,
// and the current position.
func (l *StatefulLexer) value() string {
	return l.input[l.start:l.pos]
}

// Iter iterates over the lexer's input and yields references to the tokens found.
func (l *StatefulLexer) Iter() iter.Seq2[*Reference, error] { //nolint:gocognit
	return func(yield func(*Reference, error) bool) {
		var ref *Reference
		var snap snapshot
		var r rune
		var err error
		for {
			err = nil
			l.mark()
			r = l.peek()
			snap = l.snapshot()

			switch {
			case r == eof:
				return
			case isWhitespace(r):
				l.scan(isWhitespace)
				ref = &Reference{
					Token: Whitespace,
					Value: l.value(),
				}
			case r == rawNewline:
				l.next()
				if l.lastToken == Equals {
					ref = &Reference{
						Token: Value,
						Value: "",
					}
				} else {
					ref = &Reference{
						Token: Newline,
						Value: l.value(),
					}
				}
			case r == rawCommentStartHash:
				l.scanUntil(func(r rune) bool { return r == rawNewline })
				ref = &Reference{
					Token: Comment,
					Value: strings.TrimRightFunc(l.value(), isWhitespace),
				}
			case (l.lastToken == EOF || l.lastToken == Value) && l.scanExact(rawExport):
				ref = &Reference{
					Token: Export,
					Value: strings.TrimSpace(l.value()),
				}
			case (l.lastToken == EOF || l.lastToken == Export || l.lastToken == Value) && isValidKeyStart(r):
				l.next()
				l.scan(isValidKeyChar)
				ref = &Reference{
					Token: Key,
					Value: l.value(),
				}
			case l.lastToken == Equals && l.scanExact("\"\"\""):
				var value string
				value, err = scanTripleQuote(l, QuoteTypeDouble)
				if err != nil {
					break
				}
				ref = &Reference{
					Token:     Value,
					Value:     value,
					QuoteType: QuoteTypeDouble,
				}
			case l.lastToken == Equals && l.scanExact("'''"):
				var value string
				value, err = scanTripleQuote(l, QuoteTypeSingle)
				if err != nil {
					break
				}
				ref = &Reference{
					Token:     Value,
					Value:     value,
					QuoteType: QuoteTypeSingle,
				}
			case l.lastToken == Equals && r == '"':
				l.next()
				var value string
				value, err = scanQuotes(l, QuoteTypeDouble)
				if err != nil {
					break
				}
				ref = &Reference{
					Token:     Value,
					Value:     value,
					QuoteType: QuoteTypeDouble,
				}
			case l.lastToken == Equals && r == '\'':
				l.next()
				var value string
				value, err = scanQuotes(l, QuoteTypeSingle)
				if err != nil {
					break
				}
				ref = &Reference{
					Token:     Value,
					Value:     value,
					QuoteType: QuoteTypeSingle,
				}
			case l.lastToken == Key && r == rawSeparator:
				l.next()
				ref = &Reference{
					Token: Equals,
					Value: l.value(),
				}
			case l.lastToken == Equals:
				l.next()
				l.scanUntil(func(r rune) bool { return r == rawNewline || isWhitespace(r) })
				ref = &Reference{
					Token: Value,
					Value: strings.TrimSpace(l.value()),
				}
			default:
				err = &GenericError{
					Err:    fmt.Errorf("invalid start of token: %q", r),
					Line:   l.line,
					Column: l.col,
				}
			}

			if err != nil {
				if !yield(nil, err) {
					return
				}
				continue
			}

			if ref.Token != Whitespace && ref.Token != Newline && ref.Token != Comment {
				l.lastToken = ref.Token
			}
			if ref != nil {
				ref.Line = snap.line
				ref.Column = snap.col
				ref.Position = l.start
			}
			if !yield(ref, nil) {
				return
			}
		}
	}
}

// isWhitespace checks if the rune is a whitespace character.
func isWhitespace(r rune) bool {
	return r == '\t' || r == '\v' || r == '\f' || r == ' '
}

// isValidKeyStart checks if the rune is a valid start of an environment variable key.
func isValidKeyStart(r rune) bool {
	// a-z, A-Z, _
	return r == '_' || (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')
}

// isValidKeyChar checks if the rune is a valid character in an environment variable key.
func isValidKeyChar(r rune) bool {
	// a-z, A-Z, 0-9, _
	return r == '_' || (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9')
}

// scanQuotes scans a single/double quote string. It assumes the first quote has already
// been read.
func scanQuotes(l *StatefulLexer, _type QuoteType) (value string, err error) {
	for {
		r := l.next()
		switch {
		case r == eof:
			return "", &GenericError{
				Err:    errors.New("unterminated quotes"),
				Line:   l.line,
				Column: l.col,
			}
		case _type == QuoteTypeDouble && r == '"':
			// Don't include the quotes, and unescape any escaped characters.
			return strings.ReplaceAll(l.input[l.start+1:l.pos-1], `\"`, `"`), nil
		case _type == QuoteTypeSingle && r == '\'':
			// Don't include the quotes, and unescape any escaped characters.
			return strings.ReplaceAll(l.input[l.start+1:l.pos-1], `\'`, "'"), nil
		case r == '\\':
			l.next()
		}
	}
}

// scanTripleQuote scans a triple quote string. It assumes the first '"""' has already
// been read.
func scanTripleQuote(l *StatefulLexer, _type QuoteType) (value string, err error) {
	count := 0
	l.scanUntil(func(r rune) bool {
		if (_type == QuoteTypeDouble && r == '"') || (_type == QuoteTypeSingle && r == '\'') {
			count++
		} else {
			count = 0
		}
		return count == 3
	})

	if count != 3 {
		return "", &GenericError{
			Err:    errors.New("unterminated quote block"),
			Line:   l.line,
			Column: l.col,
		}
	}

	// when scanUntils func returns true, it only peeked at the next rune, so
	// move it forward one more.
	l.next()

	q := l.input[l.start+3 : l.pos-3]
	q = strings.TrimPrefix(q, "\n") // Trim only 1, if it exists.
	q = strings.TrimSuffix(q, "\n") // Trim only 1, if it exists.

	switch _type {
	case QuoteTypeDouble:
		return strings.ReplaceAll(q, `\"\"\"`, `"""`), nil
	case QuoteTypeSingle:
		return strings.ReplaceAll(q, `\'\'\'`, `'''`), nil
	default:
		return "", &GenericError{
			Err:    errors.New("invalid quote type"),
			Line:   l.line,
			Column: l.col,
		}
	}
}
