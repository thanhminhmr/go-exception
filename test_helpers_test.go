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

func requireExceptionFields(t *testing.T, err exception.Exception, typ, message string) {
	t.Helper()
	if err.GetType() != typ {
		t.Errorf("GetType(): got %q, want %q", err.GetType(), typ)
	}
	if err.GetMessage() != message {
		t.Errorf("GetMessage(): got %q, want %q", err.GetMessage(), message)
	}
}

func requireErrors(t *testing.T, actual []error, expected ...error) {
	t.Helper()
	if len(actual) != len(expected) {
		t.Fatalf("error count: got %d, want %d (got=%v, want=%v)", len(actual), len(expected), actual, expected)
	}
	for i, want := range expected {
		if actual[i].Error() != want.Error() {
			t.Errorf("error[%d]: got %v, want %v", i, actual[i], want)
		}
	}
}

func requireValidStackTrace(t *testing.T, trace exception.StackFrames) {
	t.Helper()
	if len(trace) == 0 {
		t.Fatal("expected non-empty stack trace")
	}
	for i, frame := range trace {
		if frame.Function == "" || frame.File == "" || frame.Line == 0 {
			t.Fatalf("frame %d incomplete: %#v", i, frame)
		}
	}
}

func requireFirstFrameSuffix(t *testing.T, trace exception.StackFrames, suffix string) {
	t.Helper()
	requireValidStackTrace(t, trace)
	if !strings.HasSuffix(trace[0].Function, suffix) {
		t.Fatalf("first frame function: got %q, want suffix %q", trace[0].Function, suffix)
	}
}

func requireNilGetter(t *testing.T, err exception.Exception) {
	t.Helper()
	if suppressed := err.GetSuppressed(); suppressed != nil {
		t.Errorf("GetSuppressed(): got %v, want nil", suppressed)
	}
	if recovered := err.GetRecovered(); recovered != nil {
		t.Errorf("GetRecovered(): got %v, want nil", recovered)
	}
	if trace := err.GetStackTrace(); trace != nil {
		t.Errorf("GetStackTrace(): got %v, want nil", trace)
	}
	if extras := err.GetExtras(); extras != nil {
		t.Errorf("GetExtras(): got %v, want nil", extras)
	}
	if v, ok := err.GetExtra("anything"); ok || v != nil {
		t.Errorf("GetExtra(): got (%v, %v), want (nil, false)", v, ok)
	}
}
