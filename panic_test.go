/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package exception_test

import (
	"strings"
	"testing"

	"github.com/thanhminhmr/go-exception"
)

func checkStackTrace(t *testing.T, trace exception.StackFrames, suffix string) {
	if len(trace) == 0 {
		t.Fatalf("expected non-empty stack trace")
	}
	for _, frame := range trace {
		if frame.Function == "" || frame.File == "" || frame.Line == 0 {
			t.Fatalf("expected function, file, and line populated, got %#v", frame)
		}
	}
	if !strings.HasSuffix(trace[0].Function, suffix) {
		t.Fatalf("expected first function is this function, got %#v", trace[0])
	}
}

func TestPanicRecoverPair(t *testing.T) {
	defer func() {
		if recovered := exception.Recover(recover()); recovered != nil {
			checkStackTrace(t, recovered.GetStackTrace(), "/go-exception_test.TestPanicRecoverPair")
		}
	}()
	exception.Panic("Test")
}

func TestRecoverRawPanic(t *testing.T) {
	defer func() {
		if recovered := exception.Recover(recover()); recovered != nil {
			checkStackTrace(t, recovered.GetStackTrace(), "/go-exception_test.TestRecoverRawPanic")
		}
	}()
	panic("Test")
}
