/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package exception

import (
	"reflect"
	"runtime"
)

// StackFrame represents a single frame in a stack trace. It contains the
// function name, source file, and line number for a point in the call stack.
type StackFrame struct {
	Function string `json:"function,omitempty"`
	File     string `json:"file,omitempty"`
	Line     int    `json:"line,omitempty"`
}

// StackFrames is a slice of [StackFrame] values. It represents a complete stack
// trace.
type StackFrames []StackFrame

// StackTrace captures the current call stack as [StackFrames], starting from the
// caller of [StackTrace] itself.
//
// The skip parameter controls how many additional stack frames are omitted. A
// value of 0 includes the caller of [StackTrace], a value of 1 skips that frame,
// and higher values skip more.
func StackTrace(skip int) StackFrames {
	// get stack trace
	const depth = 64
	var programCounters [depth]uintptr
	programCountersLength := runtime.Callers(2+skip, programCounters[:])
	frames := runtime.CallersFrames(programCounters[:programCountersLength])
	// create stack frames
	stack := make([]StackFrame, 0, programCountersLength)
	for {
		frame, more := frames.Next()
		stack = append(stack, StackFrame{
			Function: frame.Function,
			File:     frame.File,
			Line:     frame.Line,
		})
		if !more {
			break
		}
	}
	return stack
}

func Function(fn any) (frame StackFrame, ok bool) {
	if fn == nil {
		return
	}
	value := reflect.ValueOf(fn)
	for value.Kind() == reflect.Pointer || value.Kind() == reflect.Interface {
		value = value.Elem()
	}
	if value.Kind() != reflect.Func {
		return
	}
	info := runtime.FuncForPC(value.Pointer())
	if info == nil {
		return
	}
	file, line := info.FileLine(info.Entry())
	return StackFrame{Function: info.Name(), File: file, Line: line}, true
}
