/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package exception_test

import (
	"testing"

	"github.com/thanhminhmr/go-exception"
)

func TestTemplate_Format(t *testing.T) {
	tests := []struct {
		name    string
		tmpl    exception.Template
		args    []any
		typ     string
		message string
		errStr  string
	}{
		{
			name:    "with format args",
			tmpl:    exception.Template("IOError: %s failed on %s"),
			args:    []any{"read", "file.txt"},
			typ:     "IOError",
			message: "read failed on file.txt",
			errStr:  "IOError: read failed on file.txt",
		},
		{
			name:    "no format args",
			tmpl:    exception.Template("StaticError: something happened"),
			typ:     "StaticError",
			message: "something happened",
			errStr:  "StaticError: something happened",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.tmpl.Format(tt.args...)
			requireExceptionFields(t, err, tt.typ, tt.message)
			if got := err.Error(); got != tt.errStr {
				t.Errorf("Error(): got %q, want %q", got, tt.errStr)
			}
		})
	}
}
