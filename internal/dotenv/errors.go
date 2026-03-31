// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

package dotenv

import (
	"errors"
	"fmt"
)

type FileAccessError struct {
	Path string `json:"path"`
	Err  error  `json:"error"`
}

func (e *FileAccessError) Unwrap() error {
	return e.Err
}

func (e *FileAccessError) Error() string {
	return fmt.Sprintf("error accessing file %q: %v", e.Path, e.Err)
}

func IsFileAccessError(err error) (error, bool) { //nolint:revive
	if err == nil {
		return nil, false
	}
	e := &FileAccessError{}
	ok := errors.As(err, &e)
	return e, ok
}

type ParseError struct {
	Content string `json:"content"`
	Line    int    `json:"line"`
	Column  int    `json:"column"`
	Err     error  `json:"error"`
}

func (e *ParseError) Unwrap() error {
	return e.Err
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("failed to parse: %v", e.Err)
}

func IsParseError(err error) (error, bool) { //nolint:revive
	if err == nil {
		return nil, false
	}
	e := &ParseError{}
	ok := errors.As(err, &e)
	return e, ok
}
