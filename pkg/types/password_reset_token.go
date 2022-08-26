package types

import (
	"context"
	"encoding/gob"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const (
	// PasswordResetTokenDataType indicates an event is related to a password reset token.
	PasswordResetTokenDataType dataType = "valid_ingredient"
)

func init() {
	gob.Register(new(PasswordResetToken))
	gob.Register(new(PasswordResetTokenCreationRequestInput))
}

type (
	// PasswordResetToken represents a password reset token.
	PasswordResetToken struct {
		_             struct{}
		LastUpdatedAt *uint64 `json:"lastUpdatedAt"`
		RedeemedAt    *uint64 `json:"archivedAt"`
		ID            string  `json:"id"`
		Token         string  `json:"token"`
		BelongsToUser string  `json:"belongsToUser"`
		ExpiresAt     uint64  `json:"expiresAt"`
		CreatedAt     uint64  `json:"createdAt"`
	}

	// UsernameReminderRequestInput represents what a user could set as input for creating password reset tokens.
	UsernameReminderRequestInput struct {
		_            struct{}
		EmailAddress string `json:"emailAddress"`
	}

	// PasswordResetTokenCreationRequestInput represents what a user could set as input for creating password reset tokens.
	PasswordResetTokenCreationRequestInput struct {
		_            struct{}
		EmailAddress string `json:"emailAddress"`
	}

	// PasswordResetTokenDatabaseCreationInput represents what a user could set as input for creating password reset tokens.
	PasswordResetTokenDatabaseCreationInput struct {
		_             struct{}
		ID            string `json:"id"`
		Token         string `json:"token"`
		BelongsToUser string `json:"belongsToUser"`
		ExpiresAt     uint64 `json:"expiresAt"`
	}

	// PasswordResetTokenRedemptionRequestInput represents what a user could set as input for creating password reset tokens.
	PasswordResetTokenRedemptionRequestInput struct {
		_           struct{}
		Token       string `json:"token"`
		NewPassword string `json:"newPassword"`
	}

	// PasswordResetTokenDataManager describes a structure capable of storing password reset tokens permanently.
	PasswordResetTokenDataManager interface {
		GetPasswordResetTokenByToken(ctx context.Context, passwordResetTokenID string) (*PasswordResetToken, error)
		CreatePasswordResetToken(ctx context.Context, input *PasswordResetTokenDatabaseCreationInput) (*PasswordResetToken, error)
		RedeemPasswordResetToken(ctx context.Context, passwordResetTokenID string) error
	}

	// PasswordResetTokenDataService describes a structure capable of serving traffic related to password reset tokens.
	PasswordResetTokenDataService interface {
		CreateHandler(res http.ResponseWriter, req *http.Request)
		ReadHandler(res http.ResponseWriter, req *http.Request)
		ArchiveHandler(res http.ResponseWriter, req *http.Request)
	}
)

var _ validation.ValidatableWithContext = (*UsernameReminderRequestInput)(nil)

// ValidateWithContext validates a UsernameReminderRequestInput.
func (x *UsernameReminderRequestInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.EmailAddress, validation.Required),
	)
}

var _ validation.ValidatableWithContext = (*PasswordResetTokenCreationRequestInput)(nil)

// ValidateWithContext validates a PasswordResetTokenCreationRequestInput.
func (x *PasswordResetTokenCreationRequestInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.EmailAddress, validation.Required),
	)
}

var _ validation.ValidatableWithContext = (*PasswordResetTokenRedemptionRequestInput)(nil)

// ValidateWithContext validates a PasswordResetTokenRedemptionRequestInput.
func (x *PasswordResetTokenRedemptionRequestInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.Token, validation.Required),
		validation.Field(&x.NewPassword, validation.Required),
	)
}

var _ validation.ValidatableWithContext = (*PasswordResetTokenDatabaseCreationInput)(nil)

// ValidateWithContext validates a PasswordResetTokenDatabaseCreationInput.
func (x *PasswordResetTokenDatabaseCreationInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.ID, validation.Required),
		validation.Field(&x.Token, validation.Required),
		validation.Field(&x.ExpiresAt, validation.Required),
		validation.Field(&x.BelongsToUser, validation.Required),
	)
}
