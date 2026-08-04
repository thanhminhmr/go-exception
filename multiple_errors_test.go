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

const (
	one   = exception.String("One: Error")
	two   = exception.String("Two: Error")
	three = exception.String("Three: Error")
	four  = exception.String("Four: Error")
)

var oneToFour = []error{one, two, three, four}

var (
	oneSqr   = exception.Join(nil, one, nil, one, nil)
	twoSqr   = exception.Join(nil, two, nil, two, nil)
	threeSqr = exception.Join(nil, three, nil, three, nil)
	fourSqr  = exception.Join(nil, four, nil, four, nil)
)

var oneToFourSqr = []error{one, one, two, two, three, three, four, four}

func TestJoin_NilHandling(t *testing.T) {
	tests := []struct {
		name       string
		args       []error
		wantNil    bool
		wantCauseN int
	}{
		{"empty", nil, true, 0},
		{"all nil", []error{nil, nil, nil}, true, 0},
		{"single", []error{one}, false, 1},
		{"multiple", []error{one, two, three}, false, 3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := exception.Join(tt.args...)
			if tt.wantNil {
				if result != nil {
					t.Errorf("Join(): got %v, want nil", result)
				}
				return
			}
			if result == nil {
				t.Fatal("Join(): got nil, want non-nil")
			}
			if n := len(result.GetCause()); n != tt.wantCauseN {
				t.Errorf("cause count: got %d, want %d", n, tt.wantCauseN)
			}
		})
	}
}

func TestJoin_FiltersNilAndFlattens(t *testing.T) {
	t.Run("filters nil", func(t *testing.T) {
		result := exception.Join(nil, one, nil, two, three, nil, nil, four, nil, nil)
		requireErrors(t, result.GetCause(), oneToFour...)
	})
	t.Run("flattens unmodified joins", func(t *testing.T) {
		result := exception.Join(nil, oneSqr, nil, twoSqr, threeSqr, nil, nil, fourSqr, nil, nil)
		requireErrors(t, result.GetCause(), oneToFourSqr...)
	})
}

func TestJoin_PreservesModifiedNestedJoin(t *testing.T) {
	joined := exception.Join(one, two)
	modified := joined.SetMessage("combined")
	result := exception.Join(modified, three)
	if n := len(result.GetCause()); n != 2 {
		t.Fatalf("cause count: got %d, want 2 (modified join preserved as one cause)", n)
	}
	requireErrors(t, result.GetCause(), modified, three)
}

func TestMultipleErrors_EmptyGetters(t *testing.T) {
	err := exception.Join(one, two)
	if err.GetType() != "" {
		t.Errorf("GetType(): got %q, want %q", err.GetType(), "")
	}
	if err.GetMessage() != "" {
		t.Errorf("GetMessage(): got %q, want %q", err.GetMessage(), "")
	}
	requireNilGetter(t, err)
}

func TestMultipleErrors_NoOp_PreservesSelf(t *testing.T) {
	joined := exception.Join(one, two)
	originalErr := joined.Error()
	noOps := []struct {
		name string
		fn   func(exception.Exception) exception.Exception
	}{
		{"SetMessage(empty)", func(e exception.Exception) exception.Exception { return e.SetMessage("") }},
		{"SetMessage(empty format)", func(e exception.Exception) exception.Exception {
			return e.SetMessage("%s", "")
		}},
		{"AddSuppressed(all nil)", func(e exception.Exception) exception.Exception {
			return e.AddSuppressed(nil, nil)
		}},
		{"SetRecovered(nil)", func(e exception.Exception) exception.Exception { return e.SetRecovered(nil) }},
		{"SetExtras(nil)", func(e exception.Exception) exception.Exception { return e.SetExtras(nil) }},
		{"SetExtra(key, nil)", func(e exception.Exception) exception.Exception { return e.SetExtra("key", nil) }},
	}
	for _, tt := range noOps {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.fn(joined)
			if result.Error() != originalErr {
				t.Errorf("%s: Error() changed: got %q, want %q", tt.name, result.Error(), originalErr)
			}
			if len(result.GetCause()) != 2 {
				t.Errorf("%s: cause count: got %d, want 2", tt.name, len(result.GetCause()))
			}
		})
	}
}

func TestMultipleErrors_SetMessage_Promotes(t *testing.T) {
	joined := exception.Join(one, two)
	result := joined.SetMessage("combined error")
	requireExceptionFields(t, result, "", "combined error")
	requireErrors(t, result.GetCause(), one, two)
}

func TestMultipleErrors_AddCause_AppendsInOrder(t *testing.T) {
	joined := exception.Join(one)
	result := joined.AddCause(two, nil, three)
	requireErrors(t, result.GetCause(), one, two, three)
}

func TestMultipleErrors_AddSuppressed_Promotes(t *testing.T) {
	joined := exception.Join(one)
	result := joined.AddSuppressed(two)
	requireErrors(t, result.GetCause(), one)
	requireErrors(t, result.GetSuppressed(), two)
}

func TestMultipleErrors_SetRecovered_Promotes(t *testing.T) {
	joined := exception.Join(one)
	result := joined.SetRecovered(42)
	requireErrors(t, result.GetCause(), one)
	if v := result.GetRecovered(); v != 42 {
		t.Errorf("GetRecovered(): got %v, want 42", v)
	}
}

func TestMultipleErrors_FillStackTrace_Promotes(t *testing.T) {
	joined := exception.Join(one)
	result := joined.FillStackTrace(0)
	requireErrors(t, result.GetCause(), one)
	requireValidStackTrace(t, result.GetStackTrace())
}

func TestMultipleErrors_SetExtras_Promotes(t *testing.T) {
	joined := exception.Join(one)
	result := joined.SetExtras(map[string]any{"key": "value"})
	requireErrors(t, result.GetCause(), one)
	if v, ok := result.GetExtra("key"); !ok || v != "value" {
		t.Errorf("GetExtra(key): got (%v, %v), want (value, true)", v, ok)
	}
}

func TestMultipleErrors_SetExtra_Promotes(t *testing.T) {
	joined := exception.Join(one)
	result := joined.SetExtra("key", "value")
	requireErrors(t, result.GetCause(), one)
	if v, ok := result.GetExtra("key"); !ok || v != "value" {
		t.Errorf("GetExtra(key): got (%v, %v), want (value, true)", v, ok)
	}
}

func TestMultipleErrors_Clone_IndependentSlice(t *testing.T) {
	joined := exception.Join(one, two)
	clone := joined.Clone()
	requireErrors(t, clone.GetCause(), one, two)

	clone.GetCause()[0] = three
	if original := joined.GetCause()[0]; original.Error() != one.Error() {
		t.Errorf("modifying clone affected original: got %v, want %v", original, one)
	}
}

func TestMultipleErrors_Error_ExactOutput(t *testing.T) {
	joined := exception.Join(one, two)
	want := "[One: Error Two: Error]"
	if got := joined.Error(); got != want {
		t.Errorf("Error(): got %q, want %q", got, want)
	}
}
