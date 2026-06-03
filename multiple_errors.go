/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package exception

import (
	"fmt"
	"slices"
)

// Join combines multiple errors into a single [Exception] with an empty type.
//
// Nil values are ignored. If no errors remain, [Join] returns nil.
//
// If any of the provided errors are Exceptions produced by [Join] and have not
// been modified further (other than adding more causes), their causes are
// automatically unwrapped and merged into the new [Exception].
//
// The resulting [Exception] exposes all non-nil errors, including those from
// unboxed joins, as its causes. Other details such as the message, suppressed
// errors, recovered value, and stack trace are left empty.
func Join(errors ...error) Exception {
	var multiple multipleErrors
	multiple.append(errors...)
	if len(multiple) > 0 {
		return multiple
	}
	return nil
}

// type check
var _ Exception = multipleErrors{}

type multipleErrors []error

func (e multipleErrors) Error() string {
	return fmt.Sprintf("%v", []error(e))
}

func (e multipleErrors) GetType() string {
	return ""
}

func (e multipleErrors) GetMessage() string {
	return ""
}

func (e multipleErrors) SetMessage(message string, parameters ...any) Exception {
	if message == "" {
		return e
	}
	if len(parameters) > 0 {
		message = fmt.Sprintf(message, parameters...)
	}
	if message == "" {
		return e
	}
	return fullException{
		Message: message,
		Cause:   e,
	}
}

func (e multipleErrors) GetCause() []error {
	return e
}

func (e multipleErrors) AddCause(errors ...error) Exception {
	e.append(errors...)
	return e
}

func (e multipleErrors) GetSuppressed() []error {
	return nil
}

func (e multipleErrors) AddSuppressed(errors ...error) Exception {
	var suppressed multipleErrors
	suppressed.append(errors...)
	if len(suppressed) > 0 {
		return fullException{
			Cause:      e,
			Suppressed: suppressed,
		}
	}
	return e
}

func (e multipleErrors) GetRecovered() any {
	return nil
}

func (e multipleErrors) SetRecovered(recovered any) Exception {
	if recovered == nil {
		return e
	}
	return fullException{
		Cause:     e,
		Recovered: recovered,
	}
}

func (e multipleErrors) GetStackTrace() StackFrames {
	return nil
}

func (e multipleErrors) FillStackTrace(skip int) Exception {
	return fullException{
		Cause:      e,
		StackTrace: StackTrace(skip + 1),
	}
}

func (e multipleErrors) GetExtras() map[string]any {
	return nil
}

func (e multipleErrors) SetExtras(extras map[string]any) Exception {
	if extras == nil {
		return e
	}
	return fullException{
		Cause:  e,
		Extras: extras,
	}
}

func (e multipleErrors) GetExtra(string) (any, bool) {
	return nil, false
}

func (e multipleErrors) SetExtra(key string, value any) Exception {
	if value == nil {
		return e
	}
	return fullException{
		Cause:  e,
		Extras: map[string]any{key: value},
	}
}

func (e multipleErrors) Clone() Exception {
	return slices.Clone(e)
}

func (e multipleErrors) __() {}

func (e multipleErrors) Unwrap() []error {
	return e
}

//goland:noinspection GoMixedReceiverTypes
func (e *multipleErrors) append(errors ...error) {
	for _, err := range errors {
		if err == nil {
			continue
		}
		//goland:noinspection GoTypeAssertionOnErrors
		if multiple, ok := err.(multipleErrors); ok {
			*e = append(*e, multiple...)
		} else {
			*e = append(*e, err)
		}
	}
}
