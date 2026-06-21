package project

import "fmt"

// RegistrationError marks mutate failures that should reject registration (CLI exit 1).
type RegistrationError struct {
	Message string
}

func (e *RegistrationError) Error() string {
	return e.Message
}

func registrationErrorf(format string, args ...any) error {
	return &RegistrationError{Message: fmt.Sprintf(format, args...)}
}
