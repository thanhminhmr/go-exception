/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package exception_test

import (
	"io"
	"strings"
	"testing"

	"github.com/thanhminhmr/go-exception"
)

func TestStackTrace(t *testing.T) {
	trace := exception.StackTrace(0)
	if len(trace) == 0 {
		t.Fatalf("expected non-empty stack trace")
	}
	for _, frame := range trace {
		if frame.Function == "" || frame.File == "" || frame.Line == 0 {
			t.Fatalf("expected function, file, and line populated, got %#v", frame)
		}
	}
	if !strings.HasSuffix(trace[0].Function, "/go-exception_test.TestStackTrace") {
		t.Fatalf("expected first function is this function, got %#v", trace[0])
	}
}

func TestStackTracePanic(t *testing.T) {
	defer exception.Recover(func(recovered exception.Exception) {
		trace := recovered.GetStackTrace()
		if len(trace) == 0 {
			t.Fatalf("expected non-empty stack trace")
		}
		for _, frame := range trace {
			if frame.Function == "" || frame.File == "" || frame.Line == 0 {
				t.Fatalf("expected function, file, and line populated, got %#v", frame)
			}
		}
		if !strings.HasSuffix(trace[0].Function, "/go-exception_test.TestStackTracePanic") {
			t.Fatalf("expected first function is this function, got %#v", trace[0])
		}
	})
	panic("DIE")
}

var globalInterface io.Closer

func TestStackTraceDereference(t *testing.T) {
	defer exception.Recover(func(recovered exception.Exception) {
		trace := recovered.GetStackTrace()
		if len(trace) == 0 {
			t.Fatalf("expected non-empty stack trace")
		}
		for _, frame := range trace {
			if frame.Function == "" || frame.File == "" || frame.Line == 0 {
				t.Fatalf("expected function, file, and line populated, got %#v", frame)
			}
		}
		if !strings.HasSuffix(trace[0].Function, "/go-exception_test.TestStackTraceDereference") {
			t.Fatalf("expected first function is this function, got %#v", trace[0])
		}
	})
	_ = globalInterface.Close()
}

func TestFunction(t *testing.T) {
	if frame, ok := exception.Function(nil); ok {
		t.Errorf("Expected to be not ok, got a frame %#v", frame)
	}
	if frame, ok := exception.Function("hello world"); ok {
		t.Errorf("Expected to be not ok, got a frame %#v", frame)
	}
	if frame, ok := exception.Function(TestFunction); !ok ||
		!strings.HasSuffix(frame.Function, "/go-exception_test.TestFunction") {
		t.Errorf("Expected to be this function, got %#v", frame)
	}
}
