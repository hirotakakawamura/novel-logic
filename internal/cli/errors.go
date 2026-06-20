package cli

import "fmt"

// ExitError carries a process exit code for cobra RunE handlers.
type ExitError struct {
	Code int
	Err  error
}

func (e *ExitError) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return "exit"
}

func exitErr(code int, err error) error {
	return &ExitError{Code: code, Err: err}
}

func exitErrf(code int, format string, args ...any) error {
	return exitErr(code, fmt.Errorf(format, args...))
}