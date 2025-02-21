package types

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// DataChangeMessage represents an event that asks a worker to write data to the datastore.
	DataChangeMessage struct {
		_ struct{} `json:"-"`

		EventType   string         `json:"eventType"`
		Context     map[string]any `json:"context,omitempty"`
		UserID      string         `json:"userID"`
		HouseholdID string         `json:"householdID,omitempty"`
	}
)

func (d *DataChangeMessage) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(ctx, &d,
		validation.Field(&d.UserID, validation.When(d.HouseholdID == "", validation.Required)),
		validation.Field(&d.HouseholdID, validation.When(d.UserID == "", validation.Required)),
		validation.Field(&d.Context, validation.Required),
	)
}
