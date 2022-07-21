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
		LastUpdatedOn *uint64 `json:"lastUpdatedOn"`
		ArchivedOn    *uint64 `json:"archivedOn"`
		ID            string  `json:"id"`
		Token         string  `json:"token"`
		BelongsToUser string  `json:"belongsToUser"`
		ExpiresAt     uint64  `json:"expiresAt"`
		CreatedOn     uint64  `json:"createdOn"`
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

	// PasswordResetTokenDataManager describes a structure capable of storing password reset tokens permanently.
	PasswordResetTokenDataManager interface {
		GetPasswordResetToken(ctx context.Context, passwordResetTokenID string) (*PasswordResetToken, error)
		GetTotalPasswordResetTokenCount(ctx context.Context) (uint64, error)
		CreatePasswordResetToken(ctx context.Context, input *PasswordResetTokenDatabaseCreationInput) (*PasswordResetToken, error)
		ArchivePasswordResetToken(ctx context.Context, passwordResetTokenID string) error
	}

	// PasswordResetTokenDataService describes a structure capable of serving traffic related to password reset tokens.
	PasswordResetTokenDataService interface {
		CreateHandler(res http.ResponseWriter, req *http.Request)
		ReadHandler(res http.ResponseWriter, req *http.Request)
		ArchiveHandler(res http.ResponseWriter, req *http.Request)
	}
)

var _ validation.ValidatableWithContext = (*PasswordResetTokenCreationRequestInput)(nil)

// ValidateWithContext validates a PasswordResetTokenCreationRequestInput.
func (x *PasswordResetTokenCreationRequestInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.EmailAddress, validation.Required),
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
