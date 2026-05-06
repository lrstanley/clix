// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

package clix

import (
	"os"
	"strings"
	"testing"

	"github.com/alecthomas/kong"
)

func TestWithAppInfoVersionInHelp(t *testing.T) {
	type Flags struct{}

	var buf strings.Builder
	oldArgs := os.Args
	t.Cleanup(func() { os.Args = oldArgs })
	os.Args = []string{"testapp", "--help"}

	New(
		WithKongOptions[Flags](
			kong.Writers(&buf, &buf),
			kong.Exit(func(int) {}),
		),
		WithAppInfo[Flags](AppInfo{
			Name:    "myapp",
			Version: "v9.9.9",
		}),
	)

	help := buf.String()
	if !strings.Contains(help, "v9.9.9") {
		t.Fatalf("expected help to include app version v9.9.9, got:\n%s", help)
	}
	if strings.Contains(help, "(devel)") {
		t.Fatalf("expected help not to show (devel) when AppInfo.Version is set, got:\n%s", help)
	}
}
