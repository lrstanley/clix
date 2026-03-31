// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

package clix

import (
	"strings"
	"testing"
)

func TestGetVersionInfo(t *testing.T) {
	version := GetVersionInfo(&AppInfo{
		Name:    "example-name",
		Version: "v1.2.3",
		Commit:  "abcd1234",
		Links: []Link{
			{Name: "example-link", URL: "https://example.com"},
		},
	})

	data := version.String()

	expected := []string{
		"example-name",
		"v1.2.3",
		"abcd1234",
		"https://example.com",
	}

	for _, e := range expected {
		if !strings.Contains(data, e) {
			t.Fatalf("expected %q to be in version output", e)
		}
	}
}
