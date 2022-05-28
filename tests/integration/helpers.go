package integration

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v5"
	"github.com/pquerna/otp/totp"
	"github.com/stretchr/testify/require"

	"github.com/prixfixeco/api_server/internal/observability/keys"
	"github.com/prixfixeco/api_server/internal/observability/logging"
	logcfg "github.com/prixfixeco/api_server/internal/observability/logging/config"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/pkg/client/httpclient"
	"github.com/prixfixeco/api_server/pkg/types"
	"github.com/prixfixeco/api_server/pkg/types/fakes"
	testutils "github.com/prixfixeco/api_server/tests/utils"
)

const (
	creationTimeout = 10 * time.Second
	waitPeriod      = 1000 * time.Millisecond
)

func stringPointer(s string) *string {
	return &s
}

func requireNotNilAndNoProblems(t *testing.T, i interface{}, err error) {
	t.Helper()

	require.NoError(t, err)
	require.NotNil(t, i)
}

func createUserAndClientForTest(ctx context.Context, t *testing.T, input *types.UserRegistrationInput) (user *types.User, cookie *http.Cookie, cookieClient, pasetoClient *httpclient.Client) {
	t.Helper()

	if input == nil {
		input = &types.UserRegistrationInput{
			EmailAddress: gofakeit.Email(),
			Username:     fakes.BuildFakeUser().Username,
			Password:     gofakeit.Password(true, true, true, true, false, 64),
		}
	}

	user, err := testutils.CreateServiceUser(ctx, urlToUse, input)
	require.NoError(t, err)

	t.Logf("created user %q with email address %s: %q", user.ID, user.EmailAddress, user.Username)

	cookie, err = testutils.GetLoginCookie(ctx, urlToUse, user)
	require.NoError(t, err)

	cookieClient, err = initializeCookiePoweredClient(cookie)
	require.NoError(t, err)

	apiClient, err := cookieClient.CreateAPIClient(ctx, cookie, &types.APIClientCreationRequestInput{
		Name: t.Name(),
		UserLoginInput: types.UserLoginInput{
			Username:  user.Username,
			Password:  user.HashedPassword,
			TOTPToken: generateTOTPTokenForUser(t, user),
		},
	})
	require.NoError(t, err)

	secretKey, err := base64.RawURLEncoding.DecodeString(apiClient.ClientSecret)
	require.NoError(t, err)

	pasetoClient, err = initializePASETOPoweredClient(apiClient.ClientID, secretKey)
	require.NoError(t, err)

	return user, cookie, cookieClient, pasetoClient
}

func initializeCookiePoweredClient(cookie *http.Cookie) (*httpclient.Client, error) {
	if parsedURLToUse == nil {
		panic("url not set!")
	}

	logger := (&logcfg.Config{Provider: logcfg.ProviderZerolog}).ProvideLogger()

	c, err := httpclient.NewClient(parsedURLToUse,
		tracing.NewNoopTracerProvider(),
		httpclient.UsingLogger(logger),
		httpclient.UsingCookie(cookie),
	)
	if err != nil {
		return nil, err
	}

	if debug {
		if setOptionErr := c.SetOptions(httpclient.UsingDebug(true)); setOptionErr != nil {
			return nil, setOptionErr
		}
	}

	return c, nil
}

func initializePASETOPoweredClient(clientID string, secretKey []byte) (*httpclient.Client, error) {
	c, err := httpclient.NewClient(parsedURLToUse,
		tracing.NewNoopTracerProvider(),
		httpclient.UsingLogger(logging.NewNoopLogger()),
		httpclient.UsingPASETO(clientID, secretKey),
	)
	if err != nil {
		return nil, err
	}

	if debug {
		if setOptionErr := c.SetOptions(httpclient.UsingDebug(true)); setOptionErr != nil {
			return nil, setOptionErr
		}
	}

	return c, nil
}

func buildSimpleClient(t *testing.T) *httpclient.Client {
	t.Helper()

	c, err := httpclient.NewClient(parsedURLToUse, tracing.NewNoopTracerProvider())
	require.NoError(t, err)

	return c
}

func generateTOTPTokenForUser(t *testing.T, u *types.User) string {
	t.Helper()

	code, err := totp.GenerateCode(u.TwoFactorSecret, time.Now().UTC())
	require.NotEmpty(t, code)
	require.NoError(t, err)

	return code
}

func buildAdminCookieAndPASETOClients(ctx context.Context, t *testing.T) (cookieClient, pasetoClient *httpclient.Client) {
	t.Helper()

	ctx, span := tracing.StartSpan(ctx)
	defer span.End()

	u := testutils.DetermineServiceURL()
	urlToUse = u.String()
	logger := (&logcfg.Config{Provider: logcfg.ProviderZerolog}).ProvideLogger()

	logger.WithValue(keys.URLKey, urlToUse).Info("checking server")
	testutils.EnsureServerIsUp(ctx, urlToUse)

	adminCookie, err := testutils.GetLoginCookie(ctx, urlToUse, premadeAdminUser)
	require.NoError(t, err)

	cClient, err := initializeCookiePoweredClient(adminCookie)
	require.NoError(t, err)

	code, err := totp.GenerateCode(premadeAdminUser.TwoFactorSecret, time.Now().UTC())
	require.NoError(t, err)

	apiClient, err := cClient.CreateAPIClient(ctx, adminCookie, &types.APIClientCreationRequestInput{
		Name: fmt.Sprintf("admin_paseto_client_%d", time.Now().UnixNano()),
		UserLoginInput: types.UserLoginInput{
			Username:  premadeAdminUser.Username,
			Password:  premadeAdminUser.HashedPassword,
			TOTPToken: code,
		},
	})
	require.NoError(t, err)

	secretKey, err := base64.RawURLEncoding.DecodeString(apiClient.ClientSecret)
	require.NoError(t, err)

	PASETOClient, err := initializePASETOPoweredClient(apiClient.ClientID, secretKey)
	require.NoError(t, err)

	return cClient, PASETOClient
}
