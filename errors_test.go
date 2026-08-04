/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package exception_test

import (
	"errors"
	"testing"

	"github.com/thanhminhmr/go-exception"
)

type testError struct{ msg string }

func (e *testError) Error() string { return e.msg }

func TestErrorsIs_Cause_Matches(t *testing.T) {
	target := exception.String("TargetError: target")
	err := exception.String("TestError: msg").AddCause(target)
	if !errors.Is(err, target) {
		t.Error("expected errors.Is to find target in cause chain")
	}
}

func TestErrorsIs_Cause_DoesNotMatch(t *testing.T) {
	target := exception.String("TargetError: target")
	other := exception.String("OtherError: other")
	err := exception.String("TestError: msg").AddCause(other)
	if errors.Is(err, target) {
		t.Error("expected errors.Is to NOT find target")
	}
}

func TestErrorsIs_JoinMember_Matches(t *testing.T) {
	joined := exception.Join(one, two, three)
	for _, target := range []error{one, two, three} {
		if !errors.Is(joined, target) {
			t.Errorf("expected errors.Is to find %q in join", target.Error())
		}
	}
}

func TestErrorsIs_NestedJoin_Matches(t *testing.T) {
	inner := exception.Join(one, two)
	outer := exception.Join(inner, three)
	if !errors.Is(outer, one) {
		t.Error("expected errors.Is to find member of nested join")
	}
}

func TestErrorsIs_PromotedJoin_Matches(t *testing.T) {
	joined := exception.Join(one, two)
	promoted := joined.SetMessage("combined")
	if !errors.Is(promoted, one) {
		t.Error("expected errors.Is to find member of promoted join")
	}
}

func TestErrorsAs_LocatesTypedError(t *testing.T) {
	target := &testError{msg: "typed error"}
	err := exception.String("TestError: msg").AddCause(target)
	var found *testError
	if !errors.As(err, &found) {
		t.Fatal("expected errors.As to find *testError in cause chain")
	}
	//goland:noinspection GoDirectComparisonOfErrors
	if found != target {
		t.Errorf("found %v, want %v", found, target)
	}
}

func TestUnwrap_CauseOrder(t *testing.T) {
	err := exception.String("TestError: msg").AddCause(one, two, three)
	u, ok := err.(interface{ Unwrap() []error })
	if !ok {
		t.Fatal("expected Unwrap() []error method")
	}
	requireErrors(t, u.Unwrap(), one, two, three)
}
