/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package exception

import (
	"bytes"
	"encoding/json"
	"reflect"
	"unsafe"
)

// MarshalJSON marshals this [Exception] as a JSON object, so that any logger
// that support JSON object dump could work seamlessly.
func (e String) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.Error())
}

var fullExceptionStruct = reflect.StructOf(reflect.VisibleFields(reflect.TypeFor[fullException]()))

func (e fullException) MarshalJSON() ([]byte, error) {
	data := reflect.NewAt(fullExceptionStruct, unsafe.Pointer(&e)).Interface()
	buffer := bytes.Buffer{}
	encoder := json.NewEncoder(&buffer)
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(data); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func (e multipleErrors) MarshalJSON() ([]byte, error) {
	switch len(e) {
	case 0:
		return nil, nil
	case 1:
		return json.Marshal(e[0])
	default:
		return json.Marshal([]error(e))
	}
}

// MarshalJSON marshals this [StackFrames] as a JSON object.
func (s StackFrames) MarshalJSON() ([]byte, error) {
	switch len(s) {
	case 0:
		return nil, nil
	case 1:
		return json.Marshal(s[0])
	default:
		return json.Marshal([]StackFrame(s))
	}
}
