package errorutil

import (
	"fmt"

	"github.com/cockroachdb/errors"
)

type TraceError struct {
	sourceErr     error
	numTracePrint int
	errors        []error
}

func NewTraceError(err error, numTracePrint int) *TraceError {
	e := &TraceError{
		sourceErr:     err,
		numTracePrint: numTracePrint,
	}

	for err != nil {
		e.errors = append(e.errors, err)
		err = errors.Unwrap(err)
		if err == nil {
			break
		}
	}
	return e
}

// Error implements the error interface.

func (e *TraceError) Error() string {
	return e.sourceErr.Error()
}

// Format implements the fmt.Formatter interface.
func (e *TraceError) Format(s fmt.State, verb rune) {
	if verb == 'v' && s.Flag('+') {
		fmt.Fprintf(s, "%+v\n", e.errors[len(e.errors)-e.numTracePrint])
	} else {
		// Simple mode. Make fmt ask the cause
		// to print itself simply.
		fmt.Fprint(s, e.sourceErr)
	}
}
