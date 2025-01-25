package integration

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/dinnerdonebetter/backend/internal/lib/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/logging/config"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/lib/server/http/utils"
	"github.com/dinnerdonebetter/backend/internal/testutils"
	"github.com/dinnerdonebetter/backend/pkg/apiclient"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"

	"github.com/pquerna/otp/totp"
	"github.com/stretchr/testify/require"
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

func createUserAndClientForTest(ctx context.Context, t *testing.T, input *types.UserRegistrationInput) (user *types.User, oauthedClient *apiclient.Client) {
	t.Helper()

	if input == nil {
		input = fakes.BuildFakeUserRegistrationInput()
	}

	user, err := testutils.CreateServiceUser(ctx, urlToUse, input)
	require.NoError(t, err)

	t.Logf("created user %s", user.ID)

	code, err := totp.GenerateCode(strings.ToUpper(user.TwoFactorSecret), time.Now().UTC())
	require.NoError(t, err)

	loginInput := &types.UserLoginInput{
		Username:  user.Username,
		Password:  user.HashedPassword,
		TOTPToken: code,
	}

	oauthedClient, err = initializeOAuth2PoweredClient(ctx, loginInput)
	require.NoError(t, err)

	return user, oauthedClient
}

func initializeOAuth2PoweredClient(ctx context.Context, input *types.UserLoginInput) (*apiclient.Client, error) {
	if parsedURLToUse == nil {
		panic("url not set!")
	}

	logger, err := (&loggingcfg.Config{Provider: loggingcfg.ProviderSlog}).ProvideLogger(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not create logger: %v", err)
	}

	c, err := apiclient.NewClient(
		parsedURLToUse,
		tracing.NewNoopTracerProvider(),
		apiclient.UsingLogger(logger),
		apiclient.UsingTracerProvider(tracing.NewNoopTracerProvider()),
		apiclient.UsingURL(urlToUse),
	)
	if err != nil {
		return nil, err
	}

	tokenResponse, err := c.LoginForToken(ctx, input)
	if err != nil {
		return nil, err
	}

	if err = c.SetOptions(apiclient.UsingOAuth2(ctx, createdClientID, createdClientSecret, []string{"household_member"}, tokenResponse.Token)); err != nil {
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

	c, err := apiclient.NewClient(
		parsedURLToUse,
		tracing.NewNoopTracerProvider(),
		apiclient.UsingTracerProvider(tracing.NewNoopTracerProvider()),
		apiclient.UsingURL(urlToUse),
	)
	require.NoError(t, err)

	return c
}

func generateTOTPTokenForUserWithoutTest(u *types.User) (string, error) {
	if u.TwoFactorSecret == "" {
		return "", errors.New("empty two factor secret")
	}

	return totp.GenerateCode(u.TwoFactorSecret, time.Now().UTC())
}

func generateTOTPTokenForUser(t *testing.T, u *types.User) string {
	t.Helper()

	code, err := generateTOTPTokenForUserWithoutTest(u)
	require.NotEmpty(t, code)
	require.NoError(t, err)

	return code
}

func buildAdminCookieAndOAuthedClients(ctx context.Context, t *testing.T) (oauthedClient *apiclient.Client) {
	t.Helper()

	u := serverutils.DetermineServiceURL()
	urlToUse = u.String()

	logger, err := (&loggingcfg.Config{Provider: loggingcfg.ProviderSlog}).ProvideLogger(ctx)
	require.NoError(t, err)
	logger.WithValue(keys.URLKey, urlToUse).Info("checking server")

	serverutils.EnsureServerIsUp(ctx, urlToUse)

	adminCode, err := totp.GenerateCode(strings.ToUpper(premadeAdminUser.TwoFactorSecret), time.Now().UTC())
	require.NoError(t, err)

	loginInput := &types.UserLoginInput{
		Username:  premadeAdminUser.Username,
		Password:  premadeAdminUser.HashedPassword,
		TOTPToken: adminCode,
	}

	oauthedClient, err = initializeOAuth2PoweredClient(ctx, loginInput)
	require.NoError(t, err)

	return oauthedClient
}
