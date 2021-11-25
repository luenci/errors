package errors

import (
	"fmt"
	"io"
)

type withCode struct {
	cause error
	code  int
	*stack
}

func WithCode(code int, err error) error {
	if err == nil {
		return nil
	}

	if newErr, ok := err.(*withCode); ok {
		newErr.code = code
		return newErr
	}

	if newErr, ok := err.(*withStack); ok {
		return &withCode{
			code:  code,
			cause: newErr.error,
			stack: newErr.stack,
		}
	}

	return &withCode{
		code:  code,
		cause: err,
		stack: callers(),
	}
}

func (w *withCode) Error() string { return w.cause.Error() }
func (w *withCode) Cause() error  { return w.cause }
func (w *withCode) Unwrap() error { return w.cause }
func (w *withCode) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			_, _ = fmt.Fprintf(s, "%+v", w.Cause())
			w.stack.Format(s, verb)
			return
		}
		fallthrough
	case 's':
		_, _ = io.WriteString(s, w.Error())
	case 'q':
		_, _ = fmt.Fprintf(s, "%q", w.Error())
	}
}
