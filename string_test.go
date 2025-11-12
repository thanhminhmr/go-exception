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

func TestStringNoMessage(t *testing.T) {
	const StringError = exception.String("Test")
	if StringError.GetType() != "Test" {
		t.Errorf("Expected to have type \"Test\" but got \"%s\"", StringError.GetType())
	}
	if StringError.GetMessage() != "" {
		t.Errorf("Expected to have empty message but got \"%s\"", StringError.GetMessage())
	}
	if StringError.Error() != "Test" {
		t.Errorf("Expected to have error string \"Test\" but got \"%s\"", StringError.Error())
	}
}

func TestStringWithMessage(t *testing.T) {
	const StringError = exception.String("Test: Message")
	if StringError.GetType() != "Test" {
		t.Errorf("Expected to have type \"Test\" but got \"%s\"", StringError.GetType())
	}
	if StringError.GetMessage() != "Message" {
		t.Errorf("Expected to have message \"Message\" but got \"%s\"", StringError.GetMessage())
	}
	if StringError.Error() != "Test: Message" {
		t.Errorf("Expected to have error string \"Test: Message\" but got \"%s\"", StringError.Error())
	}
}

func TestStringNoType(t *testing.T) {
	const StringError = exception.String(": Message")
	if StringError.GetType() != "" {
		t.Errorf("Expected to have empty type but got \"%s\"", StringError.GetType())
	}
	if StringError.GetMessage() != "Message" {
		t.Errorf("Expected to have message \"Message\" but got \"%s\"", StringError.GetMessage())
	}
	if StringError.Error() != "Message" {
		t.Errorf("Expected to have error string \"Message\" but got \"%s\"", StringError.Error())
	}
}

func TestStringEmpty(t *testing.T) {
	const StringError = exception.String("")
	if StringError.GetType() != "" {
		t.Errorf("Expected to have empty type but got \"%s\"", StringError.GetType())
	}
	if StringError.GetMessage() != "" {
		t.Errorf("Expected to have empty message but got \"%s\"", StringError.GetMessage())
	}
	if StringError.Error() != "" {
		t.Errorf("Expected to have empty error string but got \"%s\"", StringError.Error())
	}
}
