package apiclient

import (
	"context"
	"net/http"
	"strconv"

	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/pkg/pointer"
	"github.com/dinnerdonebetter/backend/pkg/apiclient/generated"
	"github.com/dinnerdonebetter/backend/pkg/types"

	openapitypes "github.com/oapi-codegen/runtime/types"
)

// UserStatus fetches a user's status.
func (c *Client) UserStatus(ctx context.Context) (*types.UserStatusResponse, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	res, err := c.authedGeneratedClient.GetAuthStatus(ctx)
	if err != nil {
		return nil, observability.PrepareError(err, span, "user status")
	}
	defer c.closeResponseBody(ctx, res)

	var apiResponse *types.APIResponse[*types.UserStatusResponse]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return nil, observability.PrepareError(err, span, "retrieving user status")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	return apiResponse.Data, nil
}

// Login fetches a login cookie.
func (c *Client) Login(ctx context.Context, input *types.UserLoginInput) (*http.Cookie, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	// validating input here requires settings knowledge, so we regrettably don't bother

	body := generated.LoginJSONRequestBody{}
	c.copyType(&body, input)

	res, err := c.unauthedGeneratedClient.Login(ctx, body)
	if err != nil {
		return nil, observability.PrepareError(err, span, "executing login request")
	}
	defer c.closeResponseBody(ctx, res)

	var apiResponse *types.APIResponse[*types.UserStatusResponse]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return nil, observability.PrepareError(err, span, "parsing login JWT response")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	if cookies := res.Cookies(); len(cookies) > 0 {
		return cookies[0], nil
	}

	return nil, ErrNoCookiesReturned
}

// AdminLogin fetches a login cookie.
func (c *Client) AdminLogin(ctx context.Context, input *types.UserLoginInput) (*http.Cookie, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	// validating input here requires settings knowledge, so we regrettably don't bother

	body := generated.AdminLoginJSONRequestBody{}
	c.copyType(&body, input)

	res, err := c.unauthedGeneratedClient.AdminLogin(ctx, body)
	if err != nil {
		return nil, observability.PrepareError(err, span, "executing login request")
	}
	defer c.closeResponseBody(ctx, res)

	var apiResponse *types.APIResponse[*types.UserStatusResponse]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return nil, observability.PrepareError(err, span, "parsing login JWT response")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	if cookies := res.Cookies(); len(cookies) > 0 {
		return cookies[0], nil
	}

	return nil, ErrNoCookiesReturned
}

// LoginForJWT fetches a JWT for a user.
func (c *Client) LoginForJWT(ctx context.Context, input *types.UserLoginInput) (string, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return "", ErrNilInputProvided
	}

	// validating input here requires settings knowledge, so we regrettably don't bother

	body := generated.LoginForJWTJSONRequestBody{}
	c.copyType(&body, input)

	res, err := c.unauthedGeneratedClient.LoginForJWT(ctx, body)
	if err != nil {
		return "", observability.PrepareError(err, span, "executing login for jwt request")
	}
	defer c.closeResponseBody(ctx, res)

	var apiResponse *types.APIResponse[*types.JWTResponse]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return "", observability.PrepareError(err, span, "parsing login JWT response")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return "", err
	}

	return apiResponse.Data.Token, nil
}

// AdminLoginForJWT fetches a JWT for a user.
func (c *Client) AdminLoginForJWT(ctx context.Context, input *types.UserLoginInput) (string, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return "", ErrNilInputProvided
	}

	// validating input here requires settings knowledge, so we regrettably don't bother

	body := generated.AdminLoginForJWTJSONRequestBody{}
	c.copyType(&body, input)

	res, err := c.unauthedGeneratedClient.AdminLoginForJWT(ctx, body)
	if err != nil {
		return "", observability.PrepareError(err, span, "executing login for JWT request")
	}
	defer c.closeResponseBody(ctx, res)

	var apiResponse *types.APIResponse[*types.JWTResponse]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return "", observability.PrepareError(err, span, "parsing login JWT response")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return "", err
	}

	return apiResponse.Data.Token, nil
}

// Logout logs a user out.
func (c *Client) Logout(ctx context.Context) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	res, err := c.authedGeneratedClient.Logout(ctx)
	if err != nil {
		return observability.PrepareError(err, span, "executing logout request")
	}
	defer c.closeResponseBody(ctx, res)

	c.authedClient.Transport = newDefaultRoundTripper(c.authedClient.Timeout, c.impersonatedUserID, c.impersonatedHouseholdID)

	return nil
}

// ChangePassword changes a user's password.
func (c *Client) ChangePassword(ctx context.Context, cookie *http.Cookie, input *types.PasswordUpdateInput) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	if cookie == nil {
		return ErrCookieRequired
	}

	if input == nil {
		return ErrNilInputProvided
	}

	// validating here requires settings knowledge, so we do not do it.

	body := generated.UpdatePasswordJSONRequestBody{}
	c.copyType(&body, input)

	res, err := c.authedGeneratedClient.UpdatePassword(ctx, body)
	if err != nil {
		return observability.PrepareError(err, span, "changing password")
	}
	defer c.closeResponseBody(ctx, res)

	if res.StatusCode != http.StatusAccepted {
		return observability.PrepareError(errInvalidResponseCode, span, "invalid response code: %d", res.StatusCode)
	}

	return nil
}

// CycleTwoFactorSecret cycles a user's 2FA secret.
func (c *Client) CycleTwoFactorSecret(ctx context.Context, cookie *http.Cookie, input *types.TOTPSecretRefreshInput) (*types.TOTPSecretRefreshResponse, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	if cookie == nil {
		return nil, ErrCookieRequired
	}

	if input == nil {
		return nil, ErrNilInputProvided
	}

	if err := input.ValidateWithContext(ctx); err != nil {
		return nil, observability.PrepareError(err, span, "validating input")
	}

	body := generated.RefreshTOTPSecretJSONRequestBody{}
	c.copyType(&body, input)

	res, err := c.authedGeneratedClient.RefreshTOTPSecret(ctx, body)
	if err != nil {
		return nil, observability.PrepareError(err, span, "cycle two factor secret")
	}
	defer c.closeResponseBody(ctx, res)

	var apiResponse *types.APIResponse[*types.TOTPSecretRefreshResponse]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return nil, observability.PrepareError(err, span, "cycling two factor secret")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	return apiResponse.Data, nil
}

// VerifyTOTPSecret verifies a 2FA secret.
func (c *Client) VerifyTOTPSecret(ctx context.Context, userID, token string) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	if userID == "" {
		return ErrInvalidIDProvided
	}

	if _, err := strconv.ParseUint(token, 10, 64); token == "" || err != nil {
		return observability.PrepareError(err, span, "invalid token provided: %q", token)
	}

	body := generated.VerifyTOTPSecretJSONRequestBody{
		TotpToken: &token,
		UserID:    &userID,
	}

	res, err := c.authedGeneratedClient.VerifyTOTPSecret(ctx, body)
	if err != nil {
		return observability.PrepareError(err, span, "verifying two factor secret")
	}
	defer c.closeResponseBody(ctx, res)

	if res.StatusCode == http.StatusBadRequest {
		return ErrInvalidTOTPToken
	} else if res.StatusCode != http.StatusAccepted {
		return observability.PrepareError(errInvalidResponseCode, span, "erroneous response code when validating TOTP secret: %d", res.StatusCode)
	}

	return nil
}

// RequestPasswordResetToken requests a password reset token.
func (c *Client) RequestPasswordResetToken(ctx context.Context, emailAddress string) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	if emailAddress == "" {
		return ErrEmptyEmailAddressProvided
	}

	body := generated.RequestPasswordResetTokenJSONRequestBody{
		EmailAddress: pointer.To(openapitypes.Email(emailAddress)),
	}

	res, err := c.unauthedGeneratedClient.RequestPasswordResetToken(ctx, body)
	if err != nil {
		return observability.PrepareError(err, span, "requesting password reset token")
	}
	defer c.closeResponseBody(ctx, res)

	if res.StatusCode != http.StatusAccepted {
		return observability.PrepareError(errInvalidResponseCode, span, "erroneous response code when requesting password reset token: %d", res.StatusCode)
	}

	return nil
}

// RedeemPasswordResetToken redeems a password reset token.
func (c *Client) RedeemPasswordResetToken(ctx context.Context, input *types.PasswordResetTokenRedemptionRequestInput) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return ErrNilInputProvided
	}

	body := generated.RedeemPasswordResetTokenJSONRequestBody{}
	c.copyType(&body, input)

	res, err := c.unauthedGeneratedClient.RedeemPasswordResetToken(ctx, body)
	if err != nil {
		return observability.PrepareError(err, span, "requesting password reset token redemption")
	}
	defer c.closeResponseBody(ctx, res)

	if res.StatusCode != http.StatusAccepted {
		return observability.PrepareError(errInvalidResponseCode, span, "erroneous response code when redeeming password reset token: %d", res.StatusCode)
	}

	return nil
}
