package integration

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/prixfixeco/backend/pkg/utils"
	"net/http"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v5"
	"github.com/pquerna/otp/totp"
	"github.com/stretchr/testify/require"

	"github.com/prixfixeco/backend/internal/observability/keys"
	"github.com/prixfixeco/backend/internal/observability/logging"
	logcfg "github.com/prixfixeco/backend/internal/observability/logging/config"
	"github.com/prixfixeco/backend/internal/observability/tracing"
	"github.com/prixfixeco/backend/pkg/apiclient"
	"github.com/prixfixeco/backend/pkg/types"
	"github.com/prixfixeco/backend/pkg/types/fakes"
	testutils "github.com/prixfixeco/backend/tests/utils"
)

func logJSON(t *testing.T, x any) {
	t.Helper()

	rawBytes, err := json.Marshal(x)
	require.NoError(t, err)

	t.Log(string(rawBytes))
}

func requireNotNilAndNoProblems(t *testing.T, i any, err error) {
	t.Helper()

	require.NoError(t, err)
	require.NotNil(t, i)
}

func createUserAndClientForTest(ctx context.Context, t *testing.T, input *types.UserRegistrationInput) (user *types.User, cookie *http.Cookie, cookieClient, pasetoClient *apiclient.Client) {
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

	cookieClient, err = initializeCookiePoweredClient(ctx, cookie)
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

func initializeCookiePoweredClient(ctx context.Context, cookie *http.Cookie) (*apiclient.Client, error) {
	if parsedURLToUse == nil {
		panic("url not set!")
	}

	logger, err := (&logcfg.Config{Provider: logcfg.ProviderZerolog}).ProvideLogger(ctx)
	if err != nil {
		return nil, err
	}

	c, err := apiclient.NewClient(parsedURLToUse,
		tracing.NewNoopTracerProvider(),
		apiclient.UsingLogger(logger),
		apiclient.UsingCookie(cookie),
	)
	if err != nil {
		return nil, err
	}

	if debug {
		if setOptionErr := c.SetOptions(apiclient.UsingDebug(true)); setOptionErr != nil {
			return nil, setOptionErr
		}
	}

	return c, nil
}

func initializePASETOPoweredClient(clientID string, secretKey []byte) (*apiclient.Client, error) {
	c, err := apiclient.NewClient(parsedURLToUse,
		tracing.NewNoopTracerProvider(),
		apiclient.UsingLogger(logging.NewNoopLogger()),
		apiclient.UsingPASETO(clientID, secretKey),
	)
	if err != nil {
		return nil, err
	}

	if debug {
		if setOptionErr := c.SetOptions(apiclient.UsingDebug(true)); setOptionErr != nil {
			return nil, setOptionErr
		}
	}

	return c, nil
}

func buildSimpleClient(t *testing.T) *apiclient.Client {
	t.Helper()

	c, err := apiclient.NewClient(parsedURLToUse, tracing.NewNoopTracerProvider())
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

func buildAdminCookieAndPASETOClients(ctx context.Context, t *testing.T) (cookieClient, pasetoClient *apiclient.Client) {
	t.Helper()

	ctx, span := tracing.StartSpan(ctx)
	defer span.End()

	u := utils.DetermineServiceURL()
	urlToUse = u.String()

	logger, err := (&logcfg.Config{Provider: logcfg.ProviderZerolog}).ProvideLogger(ctx)
	require.NoError(t, err)

	logger.WithValue(keys.URLKey, urlToUse).Info("checking server")
	utils.EnsureServerIsUp(ctx, urlToUse)

	adminCookie, err := testutils.GetLoginCookie(ctx, urlToUse, premadeAdminUser)
	require.NoError(t, err)

	cClient, err := initializeCookiePoweredClient(ctx, adminCookie)
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
