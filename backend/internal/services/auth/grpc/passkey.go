package grpc

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/authentication/webauthn"
	"github.com/dinnerdonebetter/backend/internal/domain/identity"
	identitymanager "github.com/dinnerdonebetter/backend/internal/domain/identity/manager"
)

// passkeyUserStore adapts identityDataManager to webauthn.UserStore.
type passkeyUserStore struct {
	identityDataManager identitymanager.IdentityDataManager
}

func (s *passkeyUserStore) GetUserByID(ctx context.Context, userID string) (*identity.User, error) {
	return s.identityDataManager.GetUser(ctx, userID)
}

func (s *passkeyUserStore) GetUserByUsername(ctx context.Context, username string) (*identity.User, error) {
	return s.identityDataManager.GetUserByUsername(ctx, username)
}

// ProvidePasskeyService creates a WebAuthn passkey service.
func ProvidePasskeyService(
	identityDataManager identitymanager.IdentityDataManager,
	identityRepo identity.Repository,
	sessionStore webauthn.SessionStore,
) (*webauthn.Service, error) {
	// TODO: read RPID and RPOrigins from config
	cfg := webauthn.Config{
		RPID:          "localhost",
		RPDisplayName: "Dinner Done Better",
		RPOrigins:     []string{"https://localhost:8080", "http://localhost:8080"},
	}
	userStore := &passkeyUserStore{identityDataManager: identityDataManager}
	return webauthn.NewService(cfg, identityRepo, userStore, sessionStore)
}
