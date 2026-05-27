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

func TestMultipleErrorsAppend(t *testing.T) {
	multiple := exception.Join(nil, one, nil, two, three, nil, nil, four, nil, nil)
	if len(multiple.GetCause()) != len(oneToFour) {
		t.Fatal(multiple)
	}
	for i, v := range multiple.GetCause() {
		//goland:noinspection GoDirectComparisonOfErrors
		if oneToFour[i] != v {
			t.Fatal(multiple)
		}
	}
}

var (
	oneSqr   = exception.Join(nil, one, nil, one, nil)
	twoSqr   = exception.Join(nil, two, nil, two, nil)
	threeSqr = exception.Join(nil, three, nil, three, nil)
	fourSqr  = exception.Join(nil, four, nil, four, nil)
)

var oneToFourSqr = []error{one, one, two, two, three, three, four, four}

func TestMultipleErrorsAppendRecursive(t *testing.T) {
	multiple := exception.Join(nil, oneSqr, nil, twoSqr, threeSqr, nil, nil, fourSqr, nil, nil)
	if len(multiple.GetCause()) != len(oneToFourSqr) {
		t.Fatal(multiple)
	}
	for i, v := range multiple.GetCause() {
		//goland:noinspection GoDirectComparisonOfErrors
		if oneToFourSqr[i] != v {
			t.Fatal(multiple)
		}
	}
}
