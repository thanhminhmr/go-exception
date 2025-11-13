/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package exception

import "fmt"

// Template represents a reusable message pattern for creating [String]
// exceptions. It behaves like a format string that can be expanded with
// parameters to produce consistent exception messages.
//
// For example:
//
//	const FileIOError = exception.Template("IOError: %s failed")
type Template string

// Format applies the given parameters to this template using [fmt.Sprintf] and
// returns a new [String] containing the formatted message.
func (t Template) Format(parameters ...any) String {
	return String(fmt.Sprintf(string(t), parameters...))
}
