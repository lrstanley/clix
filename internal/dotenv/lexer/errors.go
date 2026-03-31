// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

package lexer

import (
	"errors"
	"fmt"
)

// GenericError is a generic error that can be returned by the lexer.
type GenericError struct {
	Err    error `json:"error"`
	Line   int   `json:"line"`
	Column int   `json:"column"`
}

func (e *GenericError) Unwrap() error {
	return e.Err
}

func (e *GenericError) Error() string {
	return fmt.Sprintf("line %d, column %d: %v", e.Line, e.Column, e.Err)
}

// IsGenericError checks if the error is a [GenericError].
func IsGenericError(err error) (error, bool) { //nolint:revive
	if err == nil {
		return nil, false
	}
	e := &GenericError{}
	ok := errors.As(err, &e)
	return e, ok
}
