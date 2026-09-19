// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

package clix

import (
	"os"
	"path/filepath"
	"testing"
)

func TestWithEnvFilesNoDotEnv(t *testing.T) {
	const (
		envKey      = "CLIX_TEST_ENV_FILE_VALUE"
		envValue    = "from-env-file"
		loadedKey   = "CLIX_TEST_ENV_FILE_ONLY"
		loadedValue = "loaded-from-env-file"
	)

	envFile := filepath.Join(t.TempDir(), ".env")
	if err := os.WriteFile(
		envFile,
		[]byte(envKey+"="+envValue+"\n"+loadedKey+"="+loadedValue+"\n"),
		0o600,
	); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name       string
		args       []string
		noDotEnv   string
		wantValue  string
		wantLoaded bool
	}{
		{
			name:       "loads env file",
			noDotEnv:   "false",
			wantValue:  "before-env-file",
			wantLoaded: true,
		},
		{
			name:      "flag disables env file",
			args:      []string{"--no-dot-env"},
			noDotEnv:  "false",
			wantValue: "before-env-file",
		},
		{
			name:      "env disables env file",
			noDotEnv:  "true",
			wantValue: "before-env-file",
		},
		{
			name:       "flag overrides env",
			args:       []string{"--no-dot-env=false"},
			noDotEnv:   "true",
			wantValue:  "before-env-file",
			wantLoaded: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldArgs := os.Args
			t.Cleanup(func() { os.Args = oldArgs })
			os.Args = append([]string{"testapp"}, tt.args...)

			t.Setenv("NO_DOT_ENV", tt.noDotEnv)
			t.Setenv(envKey, "before-env-file")
			if err := os.Unsetenv(loadedKey); err != nil {
				t.Fatal(err)
			}
			t.Cleanup(func() {
				if err := os.Unsetenv(loadedKey); err != nil {
					t.Errorf("restore %s: %v", loadedKey, err)
				}
			})
			t.Chdir(filepath.Dir(envFile))

			New(WithEnvFiles[struct{}]())

			value, ok := os.LookupEnv(envKey)
			if !ok {
				t.Fatal("expected test environment variable to exist")
			}
			if value != tt.wantValue {
				t.Fatalf("expected env file value %q, got %q", tt.wantValue, value)
			}

			loaded, ok := os.LookupEnv(loadedKey)
			if tt.wantLoaded && (!ok || loaded != loadedValue) {
				t.Fatalf("expected %s=%q, got %q", loadedKey, loadedValue, loaded)
			}
			if !tt.wantLoaded && ok {
				t.Fatalf("expected %s to remain unset, got %q", loadedKey, loaded)
			}
		})
	}
}
