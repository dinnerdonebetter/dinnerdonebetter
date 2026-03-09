package auth

import (
	"context"
	"fmt"

	"github.com/dinnerdonebetter/backend/internal/config"
	"github.com/dinnerdonebetter/backend/internal/localdev"
	"github.com/dinnerdonebetter/backend/pkg/client"
)

// BuildAuthedClient constructs an authenticated gRPC client from an access token.
// When developingLocally is true, uses insecure (plaintext) gRPC; otherwise uses TLS.
func BuildAuthedClient(
	ctx context.Context,
	conn config.APIServiceOAuth2ConnectionConfig,
	accessToken string,
	developingLocally bool,
) (client.Client, error) {
	if developingLocally {
		return localdev.BuildInsecureOAuthedGRPCClient(
			ctx,
			conn.OAuth2APIClientID,
			conn.OAuth2APIClientSecret,
			conn.HTTPAPIServerURL,
			conn.GRPCAPIServerURL,
			accessToken,
		)
	}

	oauthOpt, err := client.WithOAuth2Credentials(
		ctx,
		conn.HTTPAPIServerURL,
		conn.OAuth2APIClientID,
		conn.OAuth2APIClientSecret,
		accessToken,
	)
	if err != nil {
		return nil, fmt.Errorf("building OAuth2 credentials: %w", err)
	}

	return client.BuildTLSGRPCClient(conn.GRPCAPIServerURL, oauthOpt)
}
