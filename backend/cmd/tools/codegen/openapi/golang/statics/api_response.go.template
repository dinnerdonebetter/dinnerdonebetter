package apiclient

import (
	"fmt"
)

type (
	ErrorCode string

	// APIResponse represents a response we might send to the user.
	APIResponse[T any] struct {
		_ struct{} `json:"-"`

		Data       T               `json:"data"`
		Pagination *Pagination     `json:"pagination"`
		Error      *APIError       `json:"error"`
		Details    ResponseDetails `json:"details"`
	}
)

// Error returns the error message.
func (e *APIError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// AsError returns the error message.
func (e *APIError) AsError() error {
	if e == nil {
		return nil
	}
	return e
}
