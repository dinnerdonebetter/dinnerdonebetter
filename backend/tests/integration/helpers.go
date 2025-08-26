package integration

import (
	"context"
	"fmt"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/platform/random"
	"github.com/dinnerdonebetter/backend/pkg/client"
	"golang.org/x/oauth2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net/http"
	"testing"
)

const (
	grpcLocalServerAddress = ":8000"
)

func buildUnauthenticatedGRPCClientForTest(t *testing.T) (client.Client, error) {
	t.Helper()

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	return client.BuildClient(grpcLocalServerAddress, opts...)
}

func buildUnauthenticatedGRPCClient(address string) (client.Client, error) {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	return client.BuildClient(address, opts...)
}

func buildAuthedGRPCClient(ctx context.Context, scopes []string, httpServerAddress, grpcServerAddress, token string) client.Client {
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

	httpClient := tracing.BuildTracedHTTPClient()
	httpClient.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	res, err := httpClient.Do(req)
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

	c, err := client.BuildClient(grpcLocalServerAddress, opts...)
	if err != nil {
		panic(err)
	}

	return c
}

// Custom insecure OAuth2 credentials that skip security checks
type insecureOAuth struct {
	TokenSource oauth2.TokenSource
}

func (i *insecureOAuth) GetRequestMetadata(_ context.Context, _ ...string) (map[string]string, error) {
	token, err := i.TokenSource.Token()
	if err != nil {
		return nil, err
	}

	return map[string]string{"authorization": token.Type() + " " + token.AccessToken}, nil
}

func (i *insecureOAuth) RequireTransportSecurity() bool {
	return false // Explicitly allow insecure transport
}
