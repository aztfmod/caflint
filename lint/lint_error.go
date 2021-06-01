package lint

import (
	"fmt"
)

type LintError struct {
	StatusCode int
	Err        error
}

func NewLintError(statusCode int, format string, a ...interface{}) *LintError {
	l := new(LintError)
	l.StatusCode = statusCode
	l.Err = fmt.Errorf(format, a...)
	return l
}
func (l *LintError) Error() string {
	return fmt.Sprintf("%v", l.Err)
}
