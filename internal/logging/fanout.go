// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

package logging

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
)

var _ slog.Handler = (*FanoutHandler)(nil) // Ensure we implement the [log/slog.Handler] interface.

// FanoutHandler distributes log records to multiple [log/slog.Handler] instances.
type FanoutHandler struct {
	handlers []slog.Handler
}

// Fanout creates a new fanout handler that distributes records to multiple
// [log/slog.Handler] instances.
func Fanout(handlers ...slog.Handler) slog.Handler {
	return &FanoutHandler{
		handlers: handlers,
	}
}

// Enabled checks if any of the underlying handlers are enabled for the given
// level. The handler is considered enabled if at least one of its child handlers
// is enabled for the specified level. If at least one handler can process the log,
// the fanout handler will attempt to distribute it.
func (h *FanoutHandler) Enabled(ctx context.Context, l slog.Level) bool {
	for i := range h.handlers {
		if h.handlers[i].Enabled(ctx, l) {
			return true
		}
	}
	return false
}

// Handle distributes a log record to all enabled handlers.
func (h *FanoutHandler) Handle(ctx context.Context, r slog.Record) error {
	var errs []error
	for i := range h.handlers {
		if h.handlers[i].Enabled(ctx, r.Level) {
			err := try(func() error {
				return h.handlers[i].Handle(ctx, r.Clone())
			})
			if err != nil {
				errs = append(errs, err)
			}
		}
	}
	return errors.Join(errs...)
}

// WithAttrs creates a new handler with additional attributes added to all child
// handlers.
func (h *FanoutHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	handlers := make([]slog.Handler, len(h.handlers))
	for i := range h.handlers {
		handlers[i] = h.handlers[i].WithAttrs(attrs)
	}
	return Fanout(handlers...)
}

// WithGroup creates a new handler with a group name applied to all child
// handlers.
func (h *FanoutHandler) WithGroup(name string) slog.Handler {
	// https://cs.opensource.google/go/go/+/master:src/log/slog/handler.go;drc=3fd729b2a14a7efcf08465cbea60a74da5457f06;l=90
	if name == "" {
		return h
	}
	handlers := make([]slog.Handler, len(h.handlers))
	for i := range h.handlers {
		handlers[i] = h.handlers[i].WithGroup(name)
	}
	return Fanout(handlers...)
}

func try(callback func() error) (err error) {
	defer func() {
		if r := recover(); r != nil {
			if perr, ok := r.(error); ok {
				err = perr
			} else {
				err = fmt.Errorf("panic: %+v", r)
			}
		}
	}()
	err = callback()
	return err
}
