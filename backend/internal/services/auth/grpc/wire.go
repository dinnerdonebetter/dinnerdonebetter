package grpc

import (
	webauthncfg "github.com/dinnerdonebetter/backend/internal/authentication/webauthn/config"
	"github.com/dinnerdonebetter/backend/internal/config"

	"github.com/google/wire"
)

// ProvideSessionStoreConfig extracts the session store config from the API service config.
func ProvideSessionStoreConfig(cfg *config.APIServiceConfig) *webauthncfg.Config {
	return &cfg.Auth.SessionStore
}

var (
	AuthSvcProviders = wire.NewSet(
		NewAuthService,
		ProvideMethodPermissions,
		ProvidePasskeyService,
		ProvideSessionStoreConfig,
		webauthncfg.ProvideSessionStore,
	)
)
