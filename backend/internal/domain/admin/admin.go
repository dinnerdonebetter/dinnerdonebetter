package admin

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// AdminUserDataManager contains administrative User functions that we don't necessarily want to expose broadly
	AdminUserDataManager interface {
		UpdateUserAccountStatus(ctx context.Context, userID string, input *UserAccountStatusUpdateInput) error
	}

	// UserAccountStatusUpdateInput represents what an admin User could provide as input for changing statuses.
	UserAccountStatusUpdateInput struct {
		_ struct{} `json:"-"`

		NewStatus    string `json:"newStatus"`
		Reason       string `json:"reason"`
		TargetUserID string `json:"targetUserID"`
	}
)

var _ validation.ValidatableWithContext = (*UserAccountStatusUpdateInput)(nil)

// ValidateWithContext ensures our struct is validatable.
func (i *UserAccountStatusUpdateInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(ctx, i,
		validation.Field(&i.NewStatus, validation.Required),
		validation.Field(&i.Reason, validation.Required),
		validation.Field(&i.TargetUserID, validation.Required),
	)
}
