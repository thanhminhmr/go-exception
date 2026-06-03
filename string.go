/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package exception

import (
	"fmt"
	"strings"
)

// separator between type and message
const separator = ": "

// type check
var _ Exception = String("")

// String is a string-based [Exception]. It behaves like a simple error
// containing only a type and a message, with no causes, suppressed errors,
// recovered value, or stack trace.
//
// Type and message are separated by a separator sequence: a colon followed by a
// space (": "). If the separator is missing, the [Exception] is considered to
// have an empty message.
//
// [String] is often used as a starting point for building a full exception with
// additional context. When causes, suppressed errors, or stack traces are added,
// a new [Exception] will be created that keeps the type and includes the added
// details:
//
//	err := exception.String("IOError: read failed").FillStackTrace(0)
//
// [String] can also be used as a constant error value, for example:
//
//	const ErrRead = exception.String("IOError: read failed")
type String string

// Error returns a string representation of this exception in the form of "Type:
// Message"
func (e String) Error() string {
	// removing the separator when either the type or the message is missing
	switch strings.Index(string(e), separator) {
	case 0:
		return string(e)[len(separator):]
	case len(e) - len(separator):
		return string(e[:len(e)-len(separator)])
	default:
		return string(e)
	}
}

// GetType returns the type of this exception.
func (e String) GetType() string {
	if i := strings.Index(string(e), separator); i >= 0 {
		return string(e)[:i]
	}
	return string(e)
}

// GetMessage returns the message of this exception.
func (e String) GetMessage() string {
	if i := strings.Index(string(e), separator); i >= 0 {
		return string(e)[i+len(separator):]
	}
	return ""
}

// SetMessage stores a message inside this exception.
//
// Note: This method may modify the current exception or return a new one. Always
// use the returned [Exception].
func (e String) SetMessage(message string, parameters ...any) Exception {
	eType := e.GetType()
	if message == "" {
		return String(eType)
	}
	if len(parameters) > 0 {
		message = fmt.Sprintf(message, parameters...)
	}
	if message == "" {
		return String(eType)
	}
	return String(eType + separator + message)
}

// GetCause returns the list of underlying causes associated with this exception.
// The slice may be empty if no causes have been specified.
func (e String) GetCause() []error {
	return nil
}

// AddCause attaches one or more underlying causes to this exception. Causes are
// typically used to represent the root errors that led to this exception being
// raised.
//
// Note: This method may modify the current exception or return a new one. Always
// use the returned [Exception].
func (e String) AddCause(errors ...error) Exception {
	var cause multipleErrors
	cause.append(errors...)
	if len(cause) > 0 {
		f := toFullException(e)
		f.Cause = cause
		return f
	}
	return e
}

// GetSuppressed returns the list of suppressed errors that were intentionally
// ignored or deferred while handling this exception. This can be useful when
// multiple errors occur, but only one is chosen as the primary failure.
func (e String) GetSuppressed() []error {
	return nil
}

// AddSuppressed attaches one or more suppressed errors to this exception.
//
// Note: This method may modify the current exception or return a new one. Always
// use the returned [Exception].
func (e String) AddSuppressed(errors ...error) Exception {
	var suppressed multipleErrors
	suppressed.append(errors...)
	if len(suppressed) > 0 {
		f := toFullException(e)
		f.Suppressed = suppressed
		return f
	}
	return e
}

// GetRecovered returns the value captured from a panic recovery, if any. It
// returns nil if no value was recovered.
func (e String) GetRecovered() any {
	return nil
}

// SetRecovered stores a recovered panic value inside this exception.
//
// Note: This method may modify the current exception or return a new one. Always
// use the returned [Exception].
func (e String) SetRecovered(recovered any) Exception {
	if recovered == nil {
		return e
	}
	f := toFullException(e)
	f.Recovered = recovered
	return f
}

// GetStackTrace returns the stack trace captured for this exception, represented
// as [StackFrames]. The result may be nil if no stack trace was filled.
func (e String) GetStackTrace() StackFrames {
	return nil
}

// FillStackTrace captures the current call stack starting from the caller of
// [FillStackTrace] itself and attaches it to this exception.
//
// The skip parameter controls how many additional stack frames are omitted. A
// value of 0 includes the caller of [FillStackTrace], a value of 1 skips that
// frame, and higher values skip more.
//
// Note: This method may modify the current exception or return a new one. Always
// use the returned [Exception].
func (e String) FillStackTrace(skip int) Exception {
	f := toFullException(e)
	f.StackTrace = StackTrace(skip + 1)
	return f
}

// GetExtras returns the additional metadata associated with this exception as a
// key-value map.
//
// Extras can be used to attach arbitrary contextual information such as request
// identifiers, user information, diagnostic values, or application-specific
// state.
//
// The returned map may be nil if no extras have been set.
func (e String) GetExtras() map[string]any {
	return nil
}

// SetExtras stores the additional metadata associated with this exception as a
// key-value map.
//
// Extras can be used to attach arbitrary contextual information such as request
// identifiers, user information, diagnostic values, or application-specific
// state.
//
// Note: This method may modify the current exception or return a new one. Always
// use the returned [Exception].
func (e String) SetExtras(extras map[string]any) Exception {
	if extras == nil {
		return e
	}
	f := toFullException(e)
	f.Extras = extras
	return f
}

// GetExtra retrieves a single extra value associated with the given key.
//
// Extras can be used to attach arbitrary contextual information such as request
// identifiers, user information, diagnostic values, or application-specific
// state.
//
// The returned boolean reports whether the key exists.
func (e String) GetExtra(string) (any, bool) {
	return nil, false
}

// SetExtra stores an additional metadata value inside this exception under the
// specified key.
//
// Extras can be used to attach arbitrary contextual information such as request
// identifiers, user information, diagnostic values, or application-specific
// state.
//
// If a value already exists for the key, it is replaced.
//
// Note: This method may modify the current exception or return a new one. Always
// use the returned [Exception].
func (e String) SetExtra(key string, value any) Exception {
	if value == nil {
		return e
	}
	f := toFullException(e)
	f.Extras = map[string]any{key: value}
	return f
}

// Clone returns an independent copy of this exception.
//
// After cloning, changes made through the methods of one [Exception] do not
// affect the other.
//
// Referenced values are not deeply cloned and may still be shared.
func (e String) Clone() Exception {
	return e
}

func (e String) __() {}

func toFullException(s String) fullException {
	t, m, _ := strings.Cut(string(s), separator)
	return fullException{
		Type:    t,
		Message: m,
	}
}
