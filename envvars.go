// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

package clix

import (
	"os"
	"sync/atomic"

	"github.com/alecthomas/kong"
	"github.com/lrstanley/clix/v2/internal/dotenv"
)

// WithEnvFiles loads environment variables from ".env" style files from the
// provided paths. If no paths are provided, it will load from the current
// working directory as ".env", but will not return an error if the file has
// access issues/doesn't exist.
func WithEnvFiles[T any](paths ...string) Option[T] {
	var initialized atomic.Bool
	return func(cli *CLI[T]) {
		if initialized.Load() {
			return
		}
		cli.kongOptions = append(cli.kongOptions, kong.WithBeforeReset(func() error {
			if initialized.Swap(true) {
				return nil
			}
			var vars map[string]string
			var err error

			if len(paths) > 0 {
				vars, err = dotenv.ParseFiles(paths...)
				if err != nil {
					return err
				}
			} else {
				vars, err = dotenv.ParseFiles(".env")
				if err != nil {
					if _, ok := dotenv.IsFileAccessError(err); ok {
						return nil
					}
					return err
				}
			}
			for k, v := range vars {
				err = os.Setenv(k, v)
				if err != nil {
					return err
				}
			}
			return nil
		}))
	}
}
