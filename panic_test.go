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

func TestRecover_NoPanic(t *testing.T) {
	called := false
	exception.Recover(func(recovered exception.Exception) {
		called = true
	})
	if called {
		t.Error("callback should not be called when no panic occurred")
	}
}

func TestRecover_RawPanic_CapturesPanicSite(t *testing.T) {
	defer exception.Recover(func(recovered exception.Exception) {
		requireFirstFrameSuffix(t, recovered.GetStackTrace(),
			"/go-exception_test.TestRecover_RawPanic_CapturesPanicSite")
	})
	panic("Test")
}

func TestRecover_RawPanic_StoresRecoveredValue(t *testing.T) {
	defer exception.Recover(func(recovered exception.Exception) {
		if recovered.GetType() != string(exception.PanicError) {
			t.Errorf("GetType(): got %q, want %q",
				recovered.GetType(), exception.PanicError)
		}
		if recovered.GetRecovered() != "raw panic" {
			t.Errorf("GetRecovered(): got %v, want 'raw panic'",
				recovered.GetRecovered())
		}
	})
	panic("raw panic")
}

func TestPanic_WrapsNormalValue(t *testing.T) {
	defer exception.Recover(func(recovered exception.Exception) {
		if recovered.GetType() != string(exception.PanicError) {
			t.Errorf("GetType(): got %q, want %q",
				recovered.GetType(), exception.PanicError)
		}
		if recovered.GetRecovered() != "test value" {
			t.Errorf("GetRecovered(): got %v, want 'test value'",
				recovered.GetRecovered())
		}
		requireFirstFrameSuffix(t, recovered.GetStackTrace(),
			"/go-exception_test.TestPanic_WrapsNormalValue")
	})
	exception.Panic("test value")
}

func TestPanic_RethrowsExistingPanicError(t *testing.T) {
	defer exception.Recover(func(recovered exception.Exception) {
		if recovered.GetRecovered() != "original" {
			t.Errorf("GetRecovered(): got %v, want 'original'", recovered.GetRecovered())
		}
		requireFirstFrameSuffix(t, recovered.GetStackTrace(),
			"/go-exception_test.TestPanic_RethrowsExistingPanicError")
	})

	defer exception.Recover(func(recovered exception.Exception) {
		exception.Panic(recovered)
	})
	exception.Panic("original")
}

func TestRecover_NonPanicErrorException_StoresAsRecovered(t *testing.T) {
	defer exception.Recover(func(recovered exception.Exception) {
		if recovered.GetType() != string(exception.PanicError) {
			t.Errorf("GetType(): got %q, want %q",
				recovered.GetType(), exception.PanicError)
		}
		ex, ok := recovered.GetRecovered().(exception.Exception)
		if !ok {
			t.Fatalf("expected recovered value to be an Exception, got %T",
				recovered.GetRecovered())
		}
		if ex.GetType() != "CustomError" {
			t.Errorf("recovered exception type: got %q, want %q",
				ex.GetType(), "CustomError")
		}
	})
	panic(exception.String("CustomError: something"))
}

func TestRecover_NilCallback_Panics(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("expected panic for nil callback")
		}
		msg, ok := r.(string)
		if !ok || msg != "BUG: callback is nil" {
			t.Errorf("expected 'BUG: callback is nil', got %v", r)
		}
	}()
	exception.Recover(nil)
}
