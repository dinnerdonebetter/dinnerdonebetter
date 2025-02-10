package integration

import (
	"context"
	"errors"
	"fmt"
	"hash/fnv"
	"strings"
	"testing"
	"time"

	"github.com/dinnerdonebetter/backend/internal/grpc/messages"
	"github.com/dinnerdonebetter/backend/internal/grpc/service"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/keys"
	loggingcfg "github.com/dinnerdonebetter/backend/internal/lib/observability/logging/config"
	"github.com/dinnerdonebetter/backend/pkg/types"

	"github.com/pquerna/otp/totp"
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	grpcoauth "google.golang.org/grpc/credentials/oauth"
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
	c := buildUnauthedGRPCClient(grpcServerAddress)

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

func buildUnauthedGRPCClient(address string) service.EatingServiceClient {
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

	opts := []grpc.DialOption{
		grpc.WithPerRPCCredentials(grpcoauth.TokenSource{
			TokenSource: oauth2Config.TokenSource(ctx, &oauth2.Token{AccessToken: token}),
		}),
		grpc.WithTransportCredentials(NewFakeTransportCredentials()),
	}

	conn, err := grpc.NewClient(grpcServerAddress, opts...)
	if err != nil {
		panic(err)
	}

	return service.NewEatingServiceClient(conn)
}

// createServiceUser creates a user.
func createServiceUser(ctx context.Context, address string, in *messages.UserRegistrationInput) (*types.User, error) {
	c := buildUnauthedGRPCClient(address)

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
