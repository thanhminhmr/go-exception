/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package exception

import (
	"fmt"
	"maps"
	"slices"
)

// type check
var _ Exception = fullException{}

type fullException struct {
	Type       string
	Message    string
	Cause      []error
	Suppressed []error
	Recovered  any
	StackTrace []StackFrame
	Extras     map[string]any
}

func (e fullException) Error() string {
	switch {
	case e.Type == "":
		return e.Message
	case e.Message == "":
		return e.Type
	default:
		return e.Type + ": " + e.Message
	}
}

func (e fullException) GetType() string {
	return e.Type
}

func (e fullException) GetMessage() string {
	return e.Message
}

func (e fullException) SetMessage(message string, parameters ...any) Exception {
	if message == "" || len(parameters) == 0 {
		e.Message = message
	} else {
		e.Message = fmt.Sprintf(message, parameters...)
	}
	return e
}

func (e fullException) GetCause() []error {
	return e.Cause
}

func (e fullException) AddCause(errors ...error) Exception {
	concat(&e.Cause, errors...)
	return e
}

func (e fullException) GetSuppressed() []error {
	return e.Suppressed
}

func (e fullException) AddSuppressed(errors ...error) Exception {
	concat(&e.Suppressed, errors...)
	return e
}

func (e fullException) GetRecovered() any {
	return e.Recovered
}

func (e fullException) SetRecovered(recovered any) Exception {
	e.Recovered = recovered
	return e
}

func (e fullException) GetStackTrace() StackFrames {
	return e.StackTrace
}

func (e fullException) FillStackTrace(skip int) Exception {
	e.StackTrace = StackTrace(skip + 1)
	return e
}

func (e fullException) GetExtras() map[string]any {
	return e.Extras
}

func (e fullException) SetExtras(extras map[string]any) Exception {
	e.Extras = extras
	return e
}

func (e fullException) GetExtra(key string) (any, bool) {
	extra, exists := e.Extras[key]
	return extra, exists
}

func (e fullException) SetExtra(key string, value any) Exception {
	if value != nil {
		e.Extras[key] = value
	} else {
		delete(e.Extras, key)
	}
	return e
}

func (e fullException) Clone() Exception {
	return fullException{
		Type:       e.Type,
		Message:    e.Message,
		Cause:      slices.Clone(e.Cause),
		Suppressed: slices.Clone(e.Suppressed),
		Recovered:  e.Recovered,
		StackTrace: slices.Clone(e.StackTrace),
		Extras:     maps.Clone(e.Extras),
	}
}

func (e fullException) __() {}

func (e fullException) Unwrap() []error {
	return e.Cause
}
