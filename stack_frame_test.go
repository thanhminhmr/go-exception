/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package exception_test

import (
	"io"
	"runtime"
	"strings"
	"testing"

	"github.com/thanhminhmr/go-exception"
)

func referenceStackTrace(skip int) exception.StackFrames {
	const depth = 64
	var pcs [depth]uintptr
	n := runtime.Callers(2+skip, pcs[:])
	frames := runtime.CallersFrames(pcs[:n])
	stack := make([]exception.StackFrame, 0, n)
	for {
		frame, more := frames.Next()
		stack = append(stack, exception.StackFrame{
			Function: frame.Function,
			File:     frame.File,
			Line:     frame.Line,
		})
		if !more {
			break
		}
	}
	return stack
}

func TestStackTrace_CapturesCallerFrame(t *testing.T) {
	trace := exception.StackTrace(0)
	requireFirstFrameSuffix(t, trace, "/go-exception_test.TestStackTrace_CapturesCallerFrame")
}

func TestStackTrace_FromPanic(t *testing.T) {
	defer exception.Recover(func(recovered exception.Exception) {
		requireFirstFrameSuffix(t, recovered.GetStackTrace(),
			"/go-exception_test.TestStackTrace_FromPanic")
	})
	panic("DIE")
}

var globalInterface io.Closer

func TestStackTrace_FromNilPointerDereference(t *testing.T) {
	defer exception.Recover(func(recovered exception.Exception) {
		requireFirstFrameSuffix(t, recovered.GetStackTrace(),
			"/go-exception_test.TestStackTrace_FromNilPointerDereference")
	})
	_ = globalInterface.Close()
}

func TestStackTrace_PreservesAllFrames(t *testing.T) {
	actual := exception.StackTrace(0)
	expected := referenceStackTrace(0)

	if len(actual) != len(expected) {
		t.Errorf("StackTrace returned %d frames, expected %d", len(actual), len(expected))
		return
	}
	if len(expected) > 0 {
		lastExpected := expected[len(expected)-1]
		lastActual := actual[len(actual)-1]
		if lastActual.Function != lastExpected.Function {
			t.Errorf("last frame function: got %q, want %q",
				lastActual.Function, lastExpected.Function)
		}
	}
}

func TestStackTrace_SkipPreservesAllFrames(t *testing.T) {
	actual := exception.StackTrace(1)
	expected := referenceStackTrace(1)

	if len(actual) != len(expected) {
		t.Errorf("StackTrace(1) returned %d frames, expected %d", len(actual), len(expected))
	}
}

func sampleFunc() {}

func TestFunction(t *testing.T) {
	t.Run("direct function", func(t *testing.T) {
		frame, ok := exception.Function(sampleFunc)
		if !ok {
			t.Fatal("expected ok=true for direct function")
		}
		if !strings.HasSuffix(frame.Function, ".sampleFunc") {
			t.Errorf("function: got %q, want suffix .sampleFunc", frame.Function)
		}
		if frame.File == "" || frame.Line == 0 {
			t.Errorf("expected file and line, got %#v", frame)
		}
	})

	t.Run("interface-wrapped function", func(t *testing.T) {
		var fn any = sampleFunc
		frame, ok := exception.Function(fn)
		if !ok {
			t.Fatal("expected ok=true for interface-wrapped function")
		}
		if !strings.HasSuffix(frame.Function, ".sampleFunc") {
			t.Errorf("function: got %q, want suffix .sampleFunc", frame.Function)
		}
	})

	t.Run("pointer to function", func(t *testing.T) {
		f := sampleFunc
		frame, ok := exception.Function(&f)
		if !ok {
			t.Fatal("expected ok=true for pointer to function")
		}
		if !strings.HasSuffix(frame.Function, ".sampleFunc") {
			t.Errorf("function: got %q, want suffix .sampleFunc", frame.Function)
		}
	})

	t.Run("nil", func(t *testing.T) {
		if _, ok := exception.Function(nil); ok {
			t.Error("expected ok=false for nil")
		}
	})

	t.Run("non-function", func(t *testing.T) {
		if _, ok := exception.Function("hello world"); ok {
			t.Error("expected ok=false for string")
		}
	})

	t.Run("typed nil pointer", func(t *testing.T) {
		var f *func()
		if _, ok := exception.Function(f); ok {
			t.Error("expected ok=false for typed nil pointer")
		}
	})
}
