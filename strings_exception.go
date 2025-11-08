/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package exception

import "fmt"

// type check
var _ Exception = stringsException{}

// WithMessage returns an [Exception] that has only type and message, with no
// causes, suppressed errors, recovered value, or stack trace.
func WithMessage(eType, message string, parameters ...any) Exception {
	if message == "" || len(parameters) == 0 {
		return stringsException{eType, message}
	}
	return stringsException{eType, fmt.Sprintf(message, parameters...)}
}

type stringsException [2]string

func (e stringsException) Error() string {
	switch {
	case e[0] == "":
		return e[1]
	case e[1] == "":
		return e[0]
	default:
		return e[0] + ": " + e[1]
	}
}

func (e stringsException) GetType() string {
	return e[0]
}

func (e stringsException) GetMessage() string {
	return e[1]
}

func (e stringsException) SetMessage(message string, parameters ...any) Exception {
	if message == "" {
		return e
	}
	if len(parameters) == 0 {
		e[1] = message
	} else {
		e[1] = fmt.Sprintf(message, parameters...)
	}
	return e
}

func (e stringsException) GetCause() []error {
	return nil
}

func (e stringsException) AddCause(errors ...error) Exception {
	var cause []error
	if combine(&cause, errors...) {
		return fullException{
			Type:    e[0],
			Message: e[1],
			Cause:   cause,
		}
	}
	return e
}

func (e stringsException) GetSuppressed() []error {
	return nil
}

func (e stringsException) AddSuppressed(errors ...error) Exception {
	var suppressed []error
	if combine(&suppressed, errors...) {
		return fullException{
			Type:       e[0],
			Message:    e[1],
			Suppressed: suppressed,
		}
	}
	return e
}

func (e stringsException) GetRecovered() any {
	return nil
}

func (e stringsException) SetRecovered(recovered any) Exception {
	if recovered == nil {
		return e
	}
	return fullException{
		Type:      e[0],
		Message:   e[1],
		Recovered: recovered,
	}
}

func (e stringsException) GetStackTrace() StackFrames {
	return nil
}

func (e stringsException) FillStackTrace(skip int) Exception {
	return fullException{
		Type:       e[0],
		Message:    e[1],
		StackTrace: StackTrace(skip + 1),
	}
}

func (e stringsException) __() {}

func (e stringsException) Is(target error) bool {
	return is(e, target)
}

func (e stringsException) As(target any) bool {
	return as(e, target)
}
