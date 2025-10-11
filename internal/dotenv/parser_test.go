// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

package dotenv

import (
	"reflect"
	"testing"
)

var testCases = []struct {
	name     string
	input    string
	err      bool
	expected map[string]string
}{
	{
		name:  "simple",
		input: `FOO=bar`,
		expected: map[string]string{
			"FOO": "bar",
		},
	},
	{
		name: "simple-export",
		input: `
		export FOO=bar
		`,
		expected: map[string]string{
			"FOO": "bar",
		},
	},
	{
		// joho/godotenv test case has OPTION_H as "1 2" without quotes, but this isn't valid for most
		// dotenv parsers, and wouldn't work with bash if you were to source the file, so we won't
		// support it.
		name: "simple-misc-spacing",
		input: `
		OPTION_A=1
		OPTION_B=2
		OPTION_C= 3
		OPTION_D =4
		OPTION_E = 5

		OPTION_F =
		OPTION_G=
		OPTION_H=   6
		`,
		expected: map[string]string{
			"OPTION_A": "1",
			"OPTION_B": "2",
			"OPTION_C": "3",
			"OPTION_D": "4",
			"OPTION_E": "5",
			"OPTION_F": "",
			"OPTION_G": "",
			"OPTION_H": "6",
		},
	},
	{
		name: "various-comments",
		input: `
		# Full line comment
		qux=thud # fred # other
		thud=fred#qux # other
		fred=qux#baz # other # more
		foo=bar # baz
		bar=foo#baz
		baz="foo"#bar`,
		expected: map[string]string{
			"qux":  "thud",
			"thud": "fred#qux",
			"fred": "qux#baz",
			"foo":  "bar",
			"bar":  "foo#baz",
			"baz":  "foo",
		},
	},
	{
		name: "exported-newline",
		input: `
		export OPTION_A=2
		export OPTION_B='\n'
		`,
		expected: map[string]string{
			"OPTION_A": "2",
			"OPTION_B": "\\n",
		},
	},
	{
		name: "invalid-line",
		input: `
		INVALID LINE
		foo=bar
		`,
		err: true,
	},
	{
		name: "misc-quoting",
		input: `
OPTION_A='1'
OPTION_B='2'
OPTION_C=''
OPTION_D='\n'
OPTION_E="1"
OPTION_F="2"
OPTION_G=""
OPTION_H="\n"
OPTION_I = "echo 'asd'"
OPTION_J='line 1
line 2'
OPTION_K='line one
this is \'quoted\'
one more line'
OPTION_L="line 1
line 2"
OPTION_M="line one
this is \"quoted\"
one more line"
		`,
		expected: map[string]string{
			"OPTION_A": "1",
			"OPTION_B": "2",
			"OPTION_C": "",
			"OPTION_D": "\\n",
			"OPTION_E": "1",
			"OPTION_F": "2",
			"OPTION_G": "",
			"OPTION_H": "\\n",
			"OPTION_I": "echo 'asd'",
			"OPTION_J": "line 1\nline 2",
			"OPTION_K": "line one\nthis is 'quoted'\none more line",
			"OPTION_L": "line 1\nline 2",
			"OPTION_M": "line one\nthis is \"quoted\"\none more line",
		},
	},
	{
		name: "env-expansion",
		input: `
		OPTION_A=1
		OPTION_B=${OPTION_A}
		OPTION_C=$OPTION_B
		OPTION_D=${OPTION_A}${OPTION_B}
		OPTION_E=${OPTION_NOT_DEFINED}
		OPTION_F="$FOO"
		OPTION_G="${OPTION_NOT_DEFINED}"
		OPTION_H='${OPTION_A}'
		OPTION_I="\${OPTION_A}"
		`,
		expected: map[string]string{
			"OPTION_A": "1",
			"OPTION_B": "1",
			"OPTION_C": "1",
			"OPTION_D": "11",
			"OPTION_E": "",
			"OPTION_F": "",
			"OPTION_G": "",
			"OPTION_H": "${OPTION_A}",
			"OPTION_I": "${OPTION_A}",
		},
	},
}

func TestParseStrings(t *testing.T) {
	t.Parallel()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			vars, err := ParseStrings(tc.input)
			if tc.err && err == nil {
				t.Fatal("expected error, got nil")
			}
			if !tc.err && err != nil {
				t.Fatalf("expected no error, got error: %v", err)
			}
			if tc.err {
				return
			}

			// Compare the expected and actual variables.
			if !reflect.DeepEqual(vars, tc.expected) {
				t.Fatalf("expected %#v, got %#v", tc.expected, vars)
			}
		})
	}
}
