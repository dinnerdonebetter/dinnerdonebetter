package integration

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"hash/fnv"
	"log"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/dinnerdonebetter/backend/internal/grpc/messages"
	"github.com/dinnerdonebetter/backend/internal/grpc/service"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/keys"
	loggingcfg "github.com/dinnerdonebetter/backend/internal/lib/observability/logging/config"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/lib/random"
	"github.com/dinnerdonebetter/backend/pkg/types"

	"github.com/pquerna/otp/totp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func requireNotNilAndNoProblems(t *testing.T, i any, err error) {
	t.Helper()

	require.NoError(t, err)
	require.NotNil(t, i)
}

func hashStringToNumber(s string) uint64 {
	h := fnv.New64a()

	if _, err := h.Write([]byte(s)); err != nil {
		panic(err)
	}

	return h.Sum64()
}

func assertJSONEquality[T1, T2 any](t *testing.T, expected T1, actual T2) {
	t.Helper()

	json1, err := json.Marshal(expected)
	require.NoError(t, err)

	json2, err := json.Marshal(actual)
	require.NoError(t, err)

	assert.JSONEq(t, string(json1), string(json2))
}

func createUserAndClientForTest(ctx context.Context, t *testing.T, httpServerAddress, grpcServerAddress string, input *messages.UserRegistrationInput) (user *types.User, oauthedClient service.EatingServiceClient) {
	t.Helper()

	if input == nil {
		input = &messages.UserRegistrationInput{
			Birthday:              timestamppb.New(time.Now().UTC()),
			EmailAddress:          fmt.Sprintf("test+%d@whatever.com", hashStringToNumber(t.Name())),
			FirstName:             fmt.Sprintf("test_%d", hashStringToNumber(t.Name())),
			HouseholdName:         fmt.Sprintf("test_%d", hashStringToNumber(t.Name())),
			LastName:              fmt.Sprintf("test_%d", hashStringToNumber(t.Name())),
			Password:              fmt.Sprintf("test_%d", hashStringToNumber(t.Name())),
			Username:              fmt.Sprintf("test_%d", hashStringToNumber(t.Name())),
			AcceptedPrivacyPolicy: true,
			AcceptedTOS:           true,
		}
	}

	user, err := createServiceUser(ctx, grpcServerAddress, input)
	require.NoError(t, err)

	t.Logf("created user %s", user.ID)

	code, err := totp.GenerateCode(strings.ToUpper(user.TwoFactorSecret), time.Now().UTC())
	require.NoError(t, err)

	loginInput := &messages.UserLoginInput{
		Username:  user.Username,
		Password:  user.HashedPassword,
		TOTPToken: code,
	}

	oauthedClient, err = initializeOAuth2PoweredClient(ctx, httpServerAddress, grpcServerAddress, loginInput)
	require.NoError(t, err)

	return user, oauthedClient
}

func initializeOAuth2PoweredClient(ctx context.Context, httpServerAddress, grpcServerAddress string, input *messages.UserLoginInput) (service.EatingServiceClient, error) {
	c := buildUnauthenticatedGRPCClient(grpcServerAddress)

	tokenResponse, err := c.LoginForToken(ctx, &messages.LoginForTokenRequest{Input: input})
	if err != nil {
		return nil, err
	}

	return buildAuthedGRPCClient(ctx, []string{"household_member"}, httpServerAddress, grpcServerAddress, tokenResponse.Result.AccessToken), nil
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

func buildAdminCookieAndOAuthedClients(ctx context.Context, httpServerAddress, grpcServerAddress string, t *testing.T) (oauthedClient service.EatingServiceClient) {
	t.Helper()

	logger, err := (&loggingcfg.Config{Provider: loggingcfg.ProviderSlog}).ProvideLogger(ctx)
	require.NoError(t, err)
	logger.WithValue(keys.URLKey, grpcServerAddress).Info("checking server")

	adminCode, err := totp.GenerateCode(strings.ToUpper(premadeAdminUser.TwoFactorSecret), time.Now().UTC())
	require.NoError(t, err)

	loginInput := &messages.UserLoginInput{
		Username:  premadeAdminUser.Username,
		Password:  premadeAdminUser.HashedPassword,
		TOTPToken: adminCode,
	}

	oauthedClient, err = initializeOAuth2PoweredClient(ctx, httpServerAddress, grpcServerAddress, loginInput)
	require.NoError(t, err)

	return oauthedClient
}

func buildUnauthenticatedGRPCClient(address string) service.EatingServiceClient {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	conn, err := grpc.NewClient(address, opts...)
	if err != nil {
		panic(err)
	}

	return service.NewEatingServiceClient(conn)
}

func buildAuthedGRPCClient(ctx context.Context, scopes []string, httpServerAddress, grpcServerAddress, token string) service.EatingServiceClient {
	state, err := random.GenerateBase64EncodedString(ctx, 32)
	if err != nil {
		panic(err)
	}

	oauth2Config := oauth2.Config{
		ClientID:     createdClientID,
		ClientSecret: createdClientSecret,
		Scopes:       scopes,
		RedirectURL:  httpServerAddress,
		Endpoint: oauth2.Endpoint{
			AuthStyle: oauth2.AuthStyleInParams,
			AuthURL:   httpServerAddress + "/oauth2/authorize",
			TokenURL:  httpServerAddress + "/oauth2/token",
		},
	}

	authCodeURL := oauth2Config.AuthCodeURL(
		state,
		oauth2.SetAuthURLParam("code_challenge_method", "plain"),
	)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		authCodeURL,
		http.NoBody,
	)
	if err != nil {
		panic(fmt.Errorf("failed to get oauth2 code: %w", err))
	}

	req.Header.Set("Authorization", "Bearer "+token)

	client := tracing.BuildTracedHTTPClient()
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	res, err := client.Do(req)
	if err != nil {
		panic(fmt.Errorf("failed to get oauth2 code: %w", err))
	}
	defer func() {
		if err = res.Body.Close(); err != nil {
			log.Println("failed to close oauth2 response body", err)
		}
	}()

	const (
		codeKey = "code"
	)

	rl, err := res.Location()
	if err != nil {
		panic(err)
	}

	code := rl.Query().Get(codeKey)
	if code == "" {
		panic("code not returned from oauth2 redirect")
	}

	oauth2Token, err := oauth2Config.Exchange(ctx, code)
	if err != nil {
		panic(err)
	}

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithPerRPCCredentials(&insecureOAuth{
			TokenSource: oauth2Config.TokenSource(ctx, oauth2Token),
		}),
	}

	conn, err := grpc.NewClient(grpcServerAddress, opts...)
	if err != nil {
		panic(err)
	}

	return service.NewEatingServiceClient(conn)
}

// createServiceUser creates a user.
func createServiceUser(ctx context.Context, address string, in *messages.UserRegistrationInput) (*types.User, error) {
	c := buildUnauthenticatedGRPCClient(address)

	ucr, err := c.CreateUser(ctx, &messages.CreateUserRequest{Input: in})
	if err != nil {
		return nil, fmt.Errorf("creating user: %w", err)
	}

	token, err := totp.GenerateCode(ucr.TwoFactorSecret, time.Now().UTC())
	if err != nil {
		return nil, fmt.Errorf("generating totp code: %w", err)
	}

	if _, err = c.VerifyTOTPSecret(ctx, &messages.VerifyTOTPSecretRequest{
		TOTPToken: token,
		UserID:    ucr.CreatedUserID,
	}); err != nil {
		return nil, fmt.Errorf("verifying totp code: %w", err)
	}

	u := &types.User{
		ID:              ucr.CreatedUserID,
		Username:        ucr.Username,
		EmailAddress:    ucr.EmailAddress,
		TwoFactorSecret: ucr.TwoFactorSecret,
		CreatedAt:       ucr.CreatedAt.AsTime(),
		// this is a dirty trick to reuse most of this type
		HashedPassword: in.Password,
	}

	return u, nil
}
