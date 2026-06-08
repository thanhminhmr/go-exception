/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package exception

import "strings"

// PanicError is the default type for exceptions created by [Panic] and
// [Recover].
const PanicError = String("panicked")

// Panic behaves like the built-in panic, but always panics with an [Exception]
// whose type is [PanicError].
//
// If the value already implements [Exception] and its type is [PanicError], it
// is re-panicked directly. This allows [Panic] to be used in a chain of recover
// handlers without changing the original panic state.
//
// Otherwise, [Panic] creates a new [Exception] that uses [PanicError] as its
// type, keeps the recovered value, and records the stack trace starting from the
// caller.
//
// Typical usage together with [Recover]:
//
//	defer exception.Recover(func(recovered exception.Exception) {
//	    // handle recovered exception
//	})
//
//	if somethingWrong {
//	    exception.Panic("bad state")
//	}
func Panic(recovered any) {
	//goland:noinspection GoTypeAssertionOnErrors
	if err, ok := recovered.(Exception); !ok || err.GetType() != string(PanicError) {
		recovered = fullException{
			Type:       string(PanicError),
			Recovered:  recovered,
			StackTrace: StackTrace(1),
		}
	}
	panic(recovered)
}

// Recover recovers from a panic and passes the recovered value to callback as an
// [Exception].
//
// If no panic occurred, [Recover] does nothing.
//
// If the recovered value already implements [Exception] and its type is
// [PanicError], it is passed to callback directly. This allows multiple recover
// handlers to work together: the first one captures the panic state, and later
// ones can observe or rethrow the same [Exception] without modification.
//
// Otherwise, [Recover] creates a new [Exception] that uses [PanicError] as its
// type, keeps the recovered value, and records the stack trace starting from
// the location where the panic occurred.
//
// [Recover] is intended to be used with defer:
//
//	defer exception.Recover(func(ex exception.Exception) {
//	    // handle recovered exception
//	})
//
//	if somethingWrong {
//	    exception.Panic("bad state")
//	}
func Recover(callback func(recovered Exception)) {
	if callback == nil {
		panic("BUG: callback is nil")
	}
	// try to recover
	recovered := recover()
	if recovered == nil {
		return
	}
	// check if chained panic
	//goland:noinspection GoTypeAssertionOnErrors
	if ex, ok := recovered.(Exception); ok && ex.GetType() == string(PanicError) {
		callback(ex)
		return
	}
	// skip to panic frame if exists
	trace := StackTrace(1)
	for len(trace) > 0 {
		frame := &trace[0]
		if function, found := strings.CutPrefix(frame.Function, "runtime."); found {
			// this is a hack at best, but it works on Linux and Windows at least
			if strings.Contains(function, "panic") {
				trace = trace[1:]
				continue
			}
		}
		break
	}
	callback(fullException{
		Type:       string(PanicError),
		Recovered:  recovered,
		StackTrace: trace,
	})
}
