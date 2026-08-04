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

func TestFullException_SetExtra_Semantics(t *testing.T) {
	tests := []struct {
		name   string
		setup  func() exception.Exception
		key    string
		value  any
		wantV  any
		wantOk bool
	}{
		{
			name:  "allocate on nil map",
			setup: func() exception.Exception { return exception.String("E: m").FillStackTrace(0) },
			key:   "key", value: "value",
			wantV: "value", wantOk: true,
		},
		{
			name:  "replace existing value",
			setup: func() exception.Exception { return exception.String("E: m").SetExtra("key", "old") },
			key:   "key", value: "new",
			wantV: "new", wantOk: true,
		},
		{
			name:  "delete on nil value",
			setup: func() exception.Exception { return exception.String("E: m").SetExtra("key", "value") },
			key:   "key", value: nil,
			wantV: nil, wantOk: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.setup().SetExtra(tt.key, tt.value)
			if v, ok := err.GetExtra(tt.key); ok != tt.wantOk || v != tt.wantV {
				t.Errorf("GetExtra(%q): got (%v, %v), want (%v, %v)", tt.key, v, ok, tt.wantV, tt.wantOk)
			}
		})
	}
}

func TestFullException_SetExtra_WorksAfterEachPromotion(t *testing.T) {
	base := exception.String("TestError: msg")
	promotions := []struct {
		name string
		fn   func(exception.String) exception.Exception
	}{
		{"AddCause", func(s exception.String) exception.Exception {
			return s.AddCause(exception.String("Cause: root"))
		}},
		{"AddSuppressed", func(s exception.String) exception.Exception {
			return s.AddSuppressed(exception.String("Suppressed: hidden"))
		}},
		{"SetRecovered", func(s exception.String) exception.Exception { return s.SetRecovered(42) }},
		{"FillStackTrace", func(s exception.String) exception.Exception { return s.FillStackTrace(0) }},
		{"SetExtras", func(s exception.String) exception.Exception {
			return s.SetExtras(map[string]any{"old": "val"})
		}},
		{"SetExtra", func(s exception.String) exception.Exception { return s.SetExtra("old", "val") }},
	}
	for _, p := range promotions {
		t.Run(p.name, func(t *testing.T) {
			err := p.fn(base).SetExtra("key", "value")
			if v, ok := err.GetExtra("key"); !ok || v != "value" {
				t.Errorf("GetExtra(key): got (%v, %v), want (value, true)", v, ok)
			}
		})
	}
}

func TestFullException_GetExtra_OnNilMap(t *testing.T) {
	err := exception.String("TestError: msg").FillStackTrace(0)
	if v, ok := err.GetExtra("anything"); ok || v != nil {
		t.Errorf("GetExtra on nil map: got (%v, %v), want (nil, false)", v, ok)
	}
}

func TestFullException_SetExtras_ReplacesMap(t *testing.T) {
	err := exception.String("TestError: msg").
		SetExtras(map[string]any{"old": "value"}).
		SetExtras(map[string]any{"new": "value"})
	if v, ok := err.GetExtra("old"); ok {
		t.Errorf("expected old key gone, got %v", v)
	}
	if v, ok := err.GetExtra("new"); !ok || v != "value" {
		t.Errorf("GetExtra(new): got (%v, %v), want (value, true)", v, ok)
	}
}

func TestFullException_SetMessage_FormatsParameters(t *testing.T) {
	err := exception.String("TestError: msg").
		AddCause(exception.String("Cause: root")).
		SetMessage("failed at %s", "step1")
	requireExceptionFields(t, err, "TestError", "failed at step1")
	if len(err.GetCause()) != 1 {
		t.Errorf("cause count: got %d, want 1", len(err.GetCause()))
	}
}

func TestFullException_AddCause_AppendsInOrder(t *testing.T) {
	err := exception.String("TestError: msg").
		AddCause(one).
		AddCause(two, three).
		AddCause(four)
	requireErrors(t, err.GetCause(), one, two, three, four)
}

func TestFullException_AddCause_FiltersNilAndFlattensJoins(t *testing.T) {
	joined := exception.Join(two, three)
	err := exception.String("TestError: msg").AddCause(nil, one, nil, joined, nil, four)
	requireErrors(t, err.GetCause(), one, two, three, four)
}

func TestFullException_AddSuppressed_AppendsRepeatedly(t *testing.T) {
	err := exception.String("TestError: msg").
		AddSuppressed(one).
		AddSuppressed(two, three).
		AddSuppressed(four)
	requireErrors(t, err.GetSuppressed(), one, two, three, four)
}

func TestFullException_SetRecovered_ReplacesAndClears(t *testing.T) {
	err := exception.String("TestError: msg").
		SetRecovered(42).
		SetRecovered("replaced")
	if v := err.GetRecovered(); v != "replaced" {
		t.Errorf("GetRecovered(): got %v, want 'replaced'", v)
	}
	err = err.SetRecovered(nil)
	if v := err.GetRecovered(); v != nil {
		t.Errorf("GetRecovered(): got %v, want nil", v)
	}
}

func TestFullException_FillStackTrace_Replaces(t *testing.T) {
	err := exception.String("TestError: msg").FillStackTrace(0)
	firstTrace := err.GetStackTrace()
	err = err.FillStackTrace(0)
	secondTrace := err.GetStackTrace()
	if len(firstTrace) == 0 || len(secondTrace) == 0 {
		t.Fatal("expected non-empty stack traces")
	}
	if len(firstTrace) != len(secondTrace) {
		t.Errorf("trace length changed: first=%d, second=%d", len(firstTrace), len(secondTrace))
	}
}

func TestFullException_ChainedMutations_PreserveAllFields(t *testing.T) {
	err := exception.String("TestError: msg").
		AddCause(exception.String("Cause1: a"), exception.String("Cause2: b")).
		AddSuppressed(exception.String("Suppressed1: x")).
		SetRecovered(42).
		SetExtras(map[string]any{"key": "value"}).
		FillStackTrace(0)

	requireExceptionFields(t, err, "TestError", "msg")
	requireErrors(t, err.GetCause(),
		exception.String("Cause1: a"), exception.String("Cause2: b"))
	requireErrors(t, err.GetSuppressed(),
		exception.String("Suppressed1: x"))
	if v := err.GetRecovered(); v != 42 {
		t.Errorf("GetRecovered(): got %v, want 42", v)
	}
	if v, ok := err.GetExtra("key"); !ok || v != "value" {
		t.Errorf("GetExtra(key): got (%v, %v), want (value, true)", v, ok)
	}
	requireValidStackTrace(t, err.GetStackTrace())
}

func TestFullException_Clone_IndependentSlicesAndMaps(t *testing.T) {
	original := exception.String("TestError: msg").
		AddCause(exception.String("Cause1: a")).
		AddSuppressed(exception.String("Suppressed1: x")).
		SetExtras(map[string]any{"key": "value"}).
		FillStackTrace(0)

	clone := original.Clone()

	requireExceptionFields(t, clone, "TestError", "msg")
	requireErrors(t, clone.GetCause(), exception.String("Cause1: a"))
	requireErrors(t, clone.GetSuppressed(), exception.String("Suppressed1: x"))
	if v, ok := clone.GetExtra("key"); !ok || v != "value" {
		t.Errorf("GetExtra(key): got (%v, %v), want (value, true)", v, ok)
	}
	requireValidStackTrace(t, clone.GetStackTrace())

	clone.GetCause()[0] = exception.String("ModifiedCause")
	if original.GetCause()[0].Error() == "ModifiedCause" {
		t.Error("modifying clone's cause slice affected original")
	}

	clone.GetSuppressed()[0] = exception.String("ModifiedSuppressed")
	if original.GetSuppressed()[0].Error() == "ModifiedSuppressed" {
		t.Error("modifying clone's suppressed slice affected original")
	}

	clone.GetStackTrace()[0] = exception.StackFrame{Function: "ModifiedFunc"}
	if original.GetStackTrace()[0].Function == "ModifiedFunc" {
		t.Error("modifying clone's stack trace affected original")
	}

	clone.GetExtras()["newKey"] = "newValue"
	if _, ok := original.GetExtra("newKey"); ok {
		t.Error("modifying clone's extras map affected original")
	}
}

func TestFullException_Error_Formatting(t *testing.T) {
	tests := []struct {
		name     string
		err      exception.Exception
		expected string
	}{
		{"type and message", exception.String("TestError: message"), "TestError: message"},
		{"type only", exception.String("TestError"), "TestError"},
		{"message only via promotion",
			exception.String(": message").AddCause(exception.String("cause")), "message"},
		{"empty via promotion",
			exception.String("").AddCause(exception.String("cause")), ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.expected {
				t.Errorf("Error(): got %q, want %q", got, tt.expected)
			}
		})
	}
}
