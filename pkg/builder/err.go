package builder

import (
	"fmt"

	"github.com/cockroachdb/errors"
)

type BuildkitdErr struct {
	err error
}

func (e *BuildkitdErr) Error() string {
	return e.err.Error()
}
func (e *BuildkitdErr) Format(s fmt.State, verb rune) { errors.FormatError(e, s, verb) }
