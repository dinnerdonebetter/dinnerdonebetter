package types

import (
	"context"
	"fmt"

	"github.com/dinnerdonebetter/backend/internal/lib/database/filtering"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// ResponseDetails represents details about the response.
	ResponseDetails struct {
		_ struct{} `json:"-"`

		CurrentHouseholdID string `json:"currentHouseholdID"`
		TraceID            string `json:"traceID"`
	}

	// APIResponse represents a response we might send to the user.
	APIResponse[T any] struct {
		_ struct{} `json:"-"`

		Data       T                     `json:"data,omitempty"`
		Pagination *filtering.Pagination `json:"pagination,omitempty"`
		Error      *APIError             `json:"error,omitempty"`
		Details    ResponseDetails       `json:"details"`
	}

	// APIError represents a response we might send to the User in the event of an error.
	APIError struct {
		_ struct{} `json:"-"`

		Message string    `json:"message"`
		Code    ErrorCode `json:"code"`
	}

	NamedID struct {
		_ struct{} `json:"-"`

		ID   string `json:"id,omitempty"`
		Name string `json:"name,omitempty"`
	}

	Range[T comparable] struct {
		Min T `json:"min"`
		Max T `json:"max"`
	}

	RangeWithOptionalUpperBound[T comparable] struct {
		Min T  `json:"min"`
		Max *T `json:"max"`
	}

	RangeUpdateRequestInput[T comparable] struct {
		Min *T `json:"min"`
		Max *T `json:"max"`
	}

	OptionalRange[T comparable] struct {
		Min *T `json:"min"`
		Max *T `json:"max"`
	}

	// OptionalFloat32Range should be replaced with a generic Range type.
	OptionalFloat32Range struct {
		Max *float32 `json:"max,omitempty"`
		Min *float32 `json:"min,omitempty"`
	}

	// Float32RangeWithOptionalMax should be replaced with a generic Range type.
	Float32RangeWithOptionalMax struct {
		Max *float32 `json:"max,omitempty"`
		Min float32  `json:"min"`
	}

	// Float32RangeWithOptionalMaxUpdateRequestInput should be replaced with a generic Range type.
	Float32RangeWithOptionalMaxUpdateRequestInput struct {
		Min *float32 `json:"min,omitempty"`
		Max *float32 `json:"max,omitempty"`
	}

	// Uint16RangeWithOptionalMax should be replaced with a generic Range type.
	Uint16RangeWithOptionalMax struct {
		Max *uint16 `json:"max,omitempty"`
		Min uint16  `json:"min"`
	}

	// Uint16RangeWithOptionalMaxUpdateRequestInput should be replaced with a generic Range type.
	Uint16RangeWithOptionalMaxUpdateRequestInput struct {
		Min *uint16 `json:"min,omitempty"`
		Max *uint16 `json:"max,omitempty"`
	}

	// OptionalUint32Range should be replaced with a generic Range type.
	OptionalUint32Range struct {
		Max *uint32 `json:"max,omitempty"`
		Min *uint32 `json:"min,omitempty"`
	}

	// Uint32RangeWithOptionalMax should be replaced with a generic Range type.
	Uint32RangeWithOptionalMax struct {
		Max *uint32 `json:"max,omitempty"`
		Min uint32  `json:"min"`
	}

	// Uint32RangeWithOptionalMaxUpdateRequestInput should be replaced with a generic Range type.
	Uint32RangeWithOptionalMaxUpdateRequestInput struct {
		Min *uint32 `json:"min,omitempty"`
		Max *uint32 `json:"max,omitempty"`
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

// NewAPIErrorResponse returns a new APIResponse with an error field.
func NewAPIErrorResponse(issue string, code ErrorCode, details ResponseDetails) *APIResponse[any] {
	return &APIResponse[any]{
		Error: &APIError{
			Message: issue,
			Code:    code,
		},
		Details: details,
	}
}

var _ validation.ValidatableWithContext = (*Float32RangeWithOptionalMax)(nil)

func (x *Float32RangeWithOptionalMax) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.Min, validation.Required),
	)
}

var _ validation.ValidatableWithContext = (*Uint16RangeWithOptionalMax)(nil)

func (x *Uint16RangeWithOptionalMax) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.Min, validation.Required),
	)
}

var _ validation.ValidatableWithContext = (*Uint32RangeWithOptionalMax)(nil)

func (x *Uint32RangeWithOptionalMax) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.Min, validation.Required),
	)
}

var _ validation.ValidatableWithContext = (*Range[string])(nil)

func (x *Range[T]) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.Min, validation.Required),
		validation.Field(&x.Max, validation.Required),
	)
}

var _ validation.ValidatableWithContext = (*RangeWithOptionalUpperBound[string])(nil)

func (x *RangeWithOptionalUpperBound[T]) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.Min, validation.Required),
	)
}
