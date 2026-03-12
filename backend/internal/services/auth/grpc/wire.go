package grpc

import (
	"github.com/dinnerdonebetter/backend/internal/authentication/webauthn"
	webauthncfg "github.com/dinnerdonebetter/backend/internal/authentication/webauthn/config"
	"github.com/dinnerdonebetter/backend/internal/branding"
	"github.com/dinnerdonebetter/backend/internal/config"

	"github.com/google/wire"
)

// ProvideSessionStoreConfig extracts the session store config from the API service config.
func ProvideSessionStoreConfig(cfg *config.APIServiceConfig) *webauthncfg.Config {
	return &cfg.Auth.SessionStore
}

// ProvidePasskeyConfig extracts passkey config from the API service config.
// When RPID is empty (e.g. local dev), returns localhost defaults.
func ProvidePasskeyConfig(cfg *config.APIServiceConfig) webauthn.Config {
	p := cfg.Auth.Passkey
	if p.RPID == "" {
		return webauthn.Config{
			RPID:          "localhost",
			RPDisplayName: branding.CompanyName,
			RPOrigins:     []string{"https://localhost:8080", "http://localhost:8080"},
		}
	}
	origins := p.RPOrigins
	if len(origins) == 0 {
		origins = []string{}
	}
	return webauthn.Config{
		RPID:          p.RPID,
		RPDisplayName: orDefault(p.RPDisplayName, branding.CompanyName),
		RPOrigins:     origins,
	}
}

func orDefault(s, def string) string {
	if s != "" {
		return s
	}
	return def
}

var (
	AuthSvcProviders = wire.NewSet(
		NewAuthService,
		ProvideMethodPermissions,
		ProvidePasskeyService,
		ProvidePasskeyConfig,
		ProvideSessionStoreConfig,
		webauthncfg.ProvideSessionStore,
	)
)
