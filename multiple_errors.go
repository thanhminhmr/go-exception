package exception

import "fmt"

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
	var multiple []error
	if !combine(&multiple, errors...) {
		return nil
	}
	return multipleErrors(multiple)
}

// type check
var _ Exception = multipleErrors{}

type multipleErrors []error

func (e multipleErrors) Error() string {
	return ""
}

func (e multipleErrors) GetType() string {
	return ""
}

func (e multipleErrors) GetMessage() string {
	return ""
}

func (e multipleErrors) SetMessage(message string, parameters ...any) Exception {
	switch {
	case message == "":
		return e
	case len(parameters) == 0:
		return fullException{
			Message: message,
			Cause:   e,
		}
	default:
		return fullException{
			Message: fmt.Sprintf(message, parameters...),
			Cause:   e,
		}
	}
}

func (e multipleErrors) GetCause() []error {
	return e
}

func (e multipleErrors) AddCause(errors ...error) Exception {
	concat((*[]error)(&e), errors...)
	return e
}

func (e multipleErrors) GetSuppressed() []error {
	return nil
}

func (e multipleErrors) AddSuppressed(errors ...error) Exception {
	var suppressed []error
	if combine(&suppressed, errors...) {
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

func (e multipleErrors) __() {}

func (e multipleErrors) Unwrap() []error {
	return e
}
