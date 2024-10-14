package testutils

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/apiclient/generated/v2"
	"github.com/dinnerdonebetter/backend/pkg/types"

	"github.com/pquerna/otp/totp"
)

var (
	errEmptyAddressUnallowed = errors.New("empty address not allowed")
)

// CreateServiceUser creates a user.
func CreateServiceUser(ctx context.Context, address string, in *types.UserRegistrationInput) (*types.User, error) {
	if address == "" {
		return nil, errEmptyAddressUnallowed
	}

	parsedAddress, err := url.Parse(address)
	if err != nil {
		return nil, err
	}

	c, err := apiclient.NewClient(
		parsedAddress,
		tracing.NewNoopTracerProvider(),
		apiclient.UsingTracingProvider(tracing.NewNoopTracerProvider()),
		apiclient.UsingURL(address),
	)
	if err != nil {
		return nil, fmt.Errorf("initializing client: %w", err)
	}

	ucr, err := c.CreateUser(ctx, in)
	if err != nil {
		return nil, err
	}

	token, tokenErr := totp.GenerateCode(ucr.TwoFactorSecret, time.Now().UTC())
	if tokenErr != nil {
		return nil, fmt.Errorf("generating totp code: %w", tokenErr)
	}

	if _, validationErr := c.VerifyTOTPSecret(ctx, &types.TOTPSecretVerificationInput{
		TOTPToken: token,
		UserID:    ucr.CreatedUserID,
	}); validationErr != nil {
		return nil, fmt.Errorf("verifying totp code: %w", validationErr)
	}

	u := &types.User{
		ID:              ucr.CreatedUserID,
		Username:        ucr.Username,
		EmailAddress:    ucr.EmailAddress,
		TwoFactorSecret: ucr.TwoFactorSecret,
		CreatedAt:       ucr.CreatedAt,
		// this is a dirty trick to reuse most of this model,
		HashedPassword: in.Password,
	}

	return u, nil
}
