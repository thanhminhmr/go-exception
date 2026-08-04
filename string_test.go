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

func TestString_Parsing(t *testing.T) {
	tests := []struct {
		name    string
		input   exception.String
		typ     string
		message string
		errStr  string
	}{
		{"type only", "Test", "Test", "", "Test"},
		{"type and message", "Test: Message", "Test", "Message", "Test: Message"},
		{"message only", ": Message", "", "Message", "Message"},
		{"trailing separator", "Type: ", "Type", "", "Type"},
		{"empty", "", "", "", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requireExceptionFields(t, tt.input, tt.typ, tt.message)
			if got := tt.input.Error(); got != tt.errStr {
				t.Errorf("Error(): got %q, want %q", got, tt.errStr)
			}
		})
	}
}

func TestString_SetMessage_FormatsParameters(t *testing.T) {
	err := exception.String("TestError").SetMessage("failed at %s:%d", "file.go", 42)
	requireExceptionFields(t, err, "TestError", "failed at file.go:42")
}

func TestString_SetMessage_EmptyClearsMessage(t *testing.T) {
	tests := []struct {
		name   string
		input  exception.String
		typ    string
		errStr string
	}{
		{"empty message on typed", "TestError: msg", "TestError", "TestError"},
		{"empty format result on untyped", "", "", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.input.SetMessage("%s", "")
			requireExceptionFields(t, err, tt.typ, "")
			if got := err.Error(); got != tt.errStr {
				t.Errorf("Error(): got %q, want %q", got, tt.errStr)
			}
		})
	}
}

func TestString_NoOp_PreservesOriginal(t *testing.T) {
	original := exception.String("TestError: msg")
	noOps := []struct {
		name string
		fn   func(exception.String) exception.Exception
	}{
		{"SetRecovered(nil)", func(e exception.String) exception.Exception { return e.SetRecovered(nil) }},
		{"SetExtras(nil)", func(e exception.String) exception.Exception { return e.SetExtras(nil) }},
		{"SetExtra(key, nil)", func(e exception.String) exception.Exception { return e.SetExtra("key", nil) }},
	}
	for _, tt := range noOps {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.fn(original)
			if result.Error() != original.Error() {
				t.Errorf("%s changed the error: got %q, want %q", tt.name, result.Error(), original.Error())
			}
		})
	}
}

func TestString_EmptyGetters_ReturnZeroValues(t *testing.T) {
	err := exception.String("TestError: msg")
	if causes := err.GetCause(); causes != nil {
		t.Errorf("GetCause(): got %v, want nil", causes)
	}
	requireNilGetter(t, err)
}

func TestString_AddCause_PromotesAndStores(t *testing.T) {
	cause := exception.String("CauseError: root")
	err := exception.String("TestError: msg").AddCause(cause)
	requireExceptionFields(t, err, "TestError", "msg")
	requireErrors(t, err.GetCause(), cause)
}

func TestString_AddSuppressed_PromotesAndStores(t *testing.T) {
	suppressed := exception.String("SuppressedError: hidden")
	err := exception.String("TestError: msg").AddSuppressed(suppressed)
	requireExceptionFields(t, err, "TestError", "msg")
	requireErrors(t, err.GetSuppressed(), suppressed)
}

func TestString_SetRecovered_PromotesAndStores(t *testing.T) {
	err := exception.String("TestError: msg").SetRecovered(42)
	requireExceptionFields(t, err, "TestError", "msg")
	if v := err.GetRecovered(); v != 42 {
		t.Errorf("GetRecovered(): got %v, want 42", v)
	}
}

func TestString_FillStackTrace_PromotesAndCaptures(t *testing.T) {
	err := exception.String("TestError: msg").FillStackTrace(0)
	requireExceptionFields(t, err, "TestError", "msg")
	requireFirstFrameSuffix(t, err.GetStackTrace(),
		"/go-exception_test.TestString_FillStackTrace_PromotesAndCaptures")
}

func TestString_SetExtras_PromotesAndStores(t *testing.T) {
	err := exception.String("TestError: msg").SetExtras(map[string]any{"key1": "value1", "key2": 42})
	requireExceptionFields(t, err, "TestError", "msg")
	if v, ok := err.GetExtra("key1"); !ok || v != "value1" {
		t.Errorf("GetExtra(key1): got (%v, %v), want (value1, true)", v, ok)
	}
	if v, ok := err.GetExtra("key2"); !ok || v != 42 {
		t.Errorf("GetExtra(key2): got (%v, %v), want (42, true)", v, ok)
	}
}

func TestString_SetExtra_PromotesAndStores(t *testing.T) {
	err := exception.String("TestError: msg").SetExtra("key", "value")
	requireExceptionFields(t, err, "TestError", "msg")
	if v, ok := err.GetExtra("key"); !ok || v != "value" {
		t.Errorf("GetExtra(key): got (%v, %v), want (value, true)", v, ok)
	}
}

func TestString_Clone_ReturnsEqual(t *testing.T) {
	original := exception.String("TestError: msg")
	clone := original.Clone()
	if clone.Error() != original.Error() {
		t.Errorf("Clone().Error(): got %q, want %q", clone.Error(), original.Error())
	}
}
