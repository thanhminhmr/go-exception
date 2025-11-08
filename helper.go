/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package exception

func is(source Exception, target error) bool {
	if targetException, ok := target.(Exception); ok {
		return source.GetType() == targetException.GetType()
	}
	return false
}

func as(source Exception, target any) bool {
	if targetException, ok := target.(*Exception); ok {
		if source.GetType() == (*targetException).GetType() {
			*targetException = source
			return true
		}
	}
	return false
}

// ========================================

func combine(result *[]error, errors ...error) (changed bool) {
	// assert result != nil
	for _, err := range errors {
		combineAdd(result, &changed, err)
	}
	return
}

func combineAdd(result *[]error, changed *bool, err error) {
	if err == nil {
		return
	}
	if multiple, ok := err.(multipleErrors); ok {
		for _, inner := range multiple {
			combineAdd(result, changed, inner)
		}
	} else {
		*result = append(*result, err)
		*changed = true
	}
}

func concat(result *[]error, errors ...error) {
	// assert result != nil
	for _, err := range errors {
		concatAdd(result, err)
	}
	return
}

func concatAdd(result *[]error, err error) {
	if err == nil {
		return
	}
	if multiple, ok := err.(multipleErrors); ok {
		for _, inner := range multiple {
			concatAdd(result, inner)
		}
	} else {
		*result = append(*result, err)
	}
}
