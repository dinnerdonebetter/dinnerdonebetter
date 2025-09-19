package client

import (
	"context"
	"fmt"
	"log"
	"net/http"

	auditgrpc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/audit"
	authgrpc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/auth"
	dataprivacygrpc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/dataprivacy"
	identitygrpc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/identity"
	internalopsgrpc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/internalops"
	mealplanninggrpc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/mealplanning"
	notificationsgrpc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/notifications"
	oauthgrpc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/oauth"
	settingsgrpc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/settings"
	webhooksgrpc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/webhooks"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/platform/random"

	"golang.org/x/oauth2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/oauth"
	"google.golang.org/grpc/metadata"
)

const (
	zuckModeUserHeader    = "X-Zuck-Mode-User"
	zuckModeAccountHeader = "X-Zuck-Mode-Account"
)

type Client interface {
	authgrpc.AuthServiceClient
	identitygrpc.IdentityServiceClient
	auditgrpc.AuditServiceClient
	dataprivacygrpc.DataPrivacyServiceClient
	internalopsgrpc.InternalOperationsClient
	mealplanninggrpc.MealPlanningServiceClient
	notificationsgrpc.UserNotificationsServiceClient
	oauthgrpc.OAuthServiceClient
	settingsgrpc.SettingsServiceClient
	webhooksgrpc.WebhooksServiceClient
}

type client struct {
	authgrpc.AuthServiceClient
	identitygrpc.IdentityServiceClient
	auditgrpc.AuditServiceClient
	dataprivacygrpc.DataPrivacyServiceClient
	internalopsgrpc.InternalOperationsClient
	mealplanninggrpc.MealPlanningServiceClient
	notificationsgrpc.UserNotificationsServiceClient
	oauthgrpc.OAuthServiceClient
	settingsgrpc.SettingsServiceClient
	webhooksgrpc.WebhooksServiceClient
}

// BuildClient builds a new Client.
func BuildClient(grpcServerAddress string, opts ...grpc.DialOption) (Client, error) {
	conn, err := grpc.NewClient(grpcServerAddress, opts...)
	if err != nil {
		return nil, fmt.Errorf("building grpc client: %w", err)
	}

	c := &client{
		AuthServiceClient:              authgrpc.NewAuthServiceClient(conn),
		IdentityServiceClient:          identitygrpc.NewIdentityServiceClient(conn),
		AuditServiceClient:             auditgrpc.NewAuditServiceClient(conn),
		DataPrivacyServiceClient:       dataprivacygrpc.NewDataPrivacyServiceClient(conn),
		InternalOperationsClient:       internalopsgrpc.NewInternalOperationsClient(conn),
		MealPlanningServiceClient:      mealplanninggrpc.NewMealPlanningServiceClient(conn),
		UserNotificationsServiceClient: notificationsgrpc.NewUserNotificationsServiceClient(conn),
		OAuthServiceClient:             oauthgrpc.NewOAuthServiceClient(conn),
		SettingsServiceClient:          settingsgrpc.NewSettingsServiceClient(conn),
		WebhooksServiceClient:          webhooksgrpc.NewWebhooksServiceClient(conn),
	}

	return c, nil
}

func WithOAuth2Credentials(
	ctx context.Context,
	authServerAddress,
	clientID,
	clientSecret,
	authToken string,
	scopes []string,
) []grpc.DialOption {
	state, err := random.GenerateBase64EncodedString(ctx, 32)
	if err != nil {
		panic(err)
	}

	oauth2Config := oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes:       scopes,
		RedirectURL:  authServerAddress,
		Endpoint: oauth2.Endpoint{
			AuthStyle: oauth2.AuthStyleInParams,
			AuthURL:   authServerAddress + "/oauth2/authorize",
			TokenURL:  authServerAddress + "/oauth2/token",
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
		panic(fmt.Errorf("failed to build oauth2 code retrieval request: %w", err))
	}

	req.Header.Set("Authorization", "Bearer "+authToken)

	c := tracing.BuildTracedHTTPClient()
	c.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	res, err := c.Do(req)
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

	ts := oauth2.ReuseTokenSource(oauth2Token, oauth2Config.TokenSource(ctx, oauth2Token))

	return []grpc.DialOption{
		grpc.WithPerRPCCredentials(oauth.TokenSource{
			TokenSource: ts,
		}),
	}
}

func ImpersonateUseAndAccountContext(ctx context.Context, userID, accountID string) context.Context {
	return metadata.NewOutgoingContext(ctx, metadata.New(map[string]string{
		zuckModeUserHeader:    userID,
		zuckModeAccountHeader: accountID,
	}))
}
