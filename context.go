// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

package clix

import (
	"context"
	"net/http"
)

type contextKey string

const (
	contextKeyCLI contextKey = "clix"
)

// NewContext returns a new context with the CLI[T] injected. It can be accessed
// from the context with [FromContext].
func (cli *CLI[T]) NewContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, contextKeyCLI, cli)
}

// NewHTTPContext is an http middleware that injects the CLI[T] into the context.
func (cli *CLI[T]) NewHTTPContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r.WithContext(cli.NewContext(r.Context())))
	})
}

// FromContext returns the CLI[T] from the context, or nil if it is not present.
func FromContext[T any](ctx context.Context) *CLI[T] {
	v := ctx.Value(contextKeyCLI)
	if v == nil {
		return nil
	}
	return v.(*CLI[T]) //nolint:errcheck
}
