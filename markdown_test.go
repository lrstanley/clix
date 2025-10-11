// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

package clix

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestWithMarkdownPlugin(t *testing.T) {
	type Flags struct {
		Foo string `name:"foo" short:"f" env:"FOO_ENV_VAR" help:"foo"`
	}

	dir := t.TempDir()
	fn := filepath.Join(dir, "generated.md")
	t.Setenv("CLIX_OUTPUT_PATH", fn)
	os.Args = []string{"clix", "generate-markdown"}

	markdownShouldExit = false
	cli := NewWithDefaults[Flags]()
	if cli.Context.Error != nil {
		t.Fatal(cli.Context.Error)
	}

	b, err := os.ReadFile(fn)
	if err != nil {
		t.Fatal(err)
	}

	data := string(b)

	fmt.Printf("markdown:\n%s\n", data) //nolint:forbidigo

	expected := []string{
		"CLI Usage Documentation: clix",
		"--help",
		"-v, --version",
		"--version-json",
		"-D, --debug",
		"-f, --foo=",
		"FOO_ENV_VAR",
		"--log.level=",
		"LOG_LEVEL",
		"--log.json",
		"LOG_JSON",
		"--log.path=",
		"LOG_PATH",
	}

	for _, e := range expected {
		if !strings.Contains(data, e) {
			t.Fatalf("expected %q to be in generated markdown", e)
		}
	}
}
