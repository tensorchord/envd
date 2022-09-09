package errorutil

import (
	"fmt"
	"strings"

	"github.com/cockroachdb/errors"
)

type TraceError struct {
	sourceErr     error
	numTracePrint int
	errors        []error
}

func extractPrefix(err, cause error) string {
	causeSuffix := cause.Error()
	errMsg := err.Error()

	if strings.HasSuffix(errMsg, causeSuffix) {
		prefix := errMsg[:len(errMsg)-len(causeSuffix)]
		if strings.HasSuffix(prefix, ": ") {
			return prefix[:len(prefix)-2]
		}
	}
	return ""
}

func NewTraceError(err error, numTracePrint int) *TraceError {
	e := &TraceError{
		sourceErr:     err,
		numTracePrint: numTracePrint,
	}
	for i := 0; i < numTracePrint; i++ {
		e.errors = append(e.errors, err)
		err = errors.Unwrap(err)
		if err == nil {
			e.numTracePrint = i + 1
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
		// Verbose mode. Make fmt ask the cause
		// to print itself verbosely.
		// for i := 0; i < len(e.errors); i++ {
		// 	fmt.Fprintf(s, "%s\n", e.errors[i])
		// }

		for i := len(e.errors) - 1; i >= 1; i-- {
			pref := extractPrefix(e.errors[i], e.errors[i-1])
			fmt.Fprintf(s, "%s\n", pref)
		}
		fmt.Fprintf(s, "%s\n", e.errors[0])
	} else {
		// Simple mode. Make fmt ask the cause
		// to print itself simply.
		fmt.Fprint(s, e.sourceErr)
	}
}
