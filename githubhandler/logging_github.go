// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

package githubhandler

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"slices"
	"sync"

	"github.com/apex/log"
	"github.com/sethvargo/go-githubactions"
)

var (
	// Default handler outputting to stderr.
	Default = New(os.Stderr)

	githubCommandFields = []string{
		"title",
		"file",
		"col",
		"endColumn",
		"line",
		"endLine",
	}

	// Strings mapping.
	Strings = [...]string{
		log.DebugLevel: "DEBUG",
		log.InfoLevel:  "INFO",
		log.WarnLevel:  "WARN",
		log.ErrorLevel: "ERROR",
		log.FatalLevel: "FATAL",
	}
)

// Handler implementation.
type Handler struct {
	pool   sync.Pool
	mu     sync.Mutex
	Writer io.Writer
}

// New handler.
func New(w io.Writer) *Handler {
	return &Handler{
		pool: sync.Pool{
			New: func() any {
				return new(bytes.Buffer)
			},
		},
		Writer: w,
	}
}

// HandleLog implements log.Handler.
func (h *Handler) HandleLog(e *log.Entry) error {
	level := Strings[e.Level]
	names := e.Fields.Names()

	buf, _ := h.pool.Get().(*bytes.Buffer)
	defer h.pool.Put(buf)
	buf.Reset()

	githubSpecific := map[string]string{}

	if e.Level != log.InfoLevel {
		for i, name := range names {
			if !slices.Contains(githubCommandFields, name) {
				continue
			}

			v := fmt.Sprintf("%v", e.Fields.Get(name))
			if v == "" {
				continue
			}

			githubSpecific[name] = v
			names = slices.Delete(names, i, i+1) // Delete from slice.
		}
	}

	if e.Level != log.InfoLevel {
		fmt.Fprintf(buf, e.Message)
	} else {
		fmt.Fprintf(buf, "[%s] %s", level, e.Message)
	}

	for _, name := range names {
		fmt.Fprintf(buf, " %s=%v", name, e.Fields.Get(name))
	}

	if e.Level != log.InfoLevel {
		action := githubactions.WithFieldsMap(githubSpecific)

		switch e.Level {
		case log.DebugLevel:
			action.Debugf(buf.String())
		case log.InfoLevel:
			action.Infof(buf.String())
		case log.WarnLevel:
			action.Warningf(buf.String())
		case log.ErrorLevel:
			action.Errorf(buf.String())
		case log.FatalLevel:
			action.Fatalf(buf.String())
		case log.InvalidLevel:
			// Invalid level, do nothing.
		}

		return nil
	}

	fmt.Fprintf(buf, "\n")

	h.mu.Lock()
	_, err := h.Writer.Write(buf.Bytes())
	h.mu.Unlock()

	return err
}
