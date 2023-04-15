package errors

import "fmt"

// ChuxModelsError is a custom error type
// that wraps an error and adds a message
// to the error.
// This is the error that is returned by
// all functions in chux-models that return
// an error.
type ChuxModelsError struct {
	// Message is the message that is
	// given by chux-models when an error
	// occurs.
	// This message is used to provide
	// more context to the error.
	// The Err field contains the actual
	// error that occurred.
	Message  string
	InnerErr error
}

// NewChuxParserError returns a new ChuxModelsError
func NewChuxModelsError(message string, err error) *ChuxModelsError {
	return &ChuxModelsError{
		Message:  message,
		InnerErr: err,
	}
}

func (e *ChuxModelsError) Error() string {
	return e.Message
}

// Unwrap returns the underlying error without
// the message added by chux-parser.
func (e *ChuxModelsError) Unwrap() error {
	return e.InnerErr
}

// handleError is a helper function that handles
// errors occuring in chux-models. This means
// that it prints the error message and the
// underlying error. It will also log the error
func handleError(err error) {
	fmt.Printf("Error: %v\n", err)
}
