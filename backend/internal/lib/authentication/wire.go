package authentication

import (
	"github.com/dinnerdonebetter/backend/internal/database"

	"github.com/google/wire"
)

var (
	// AuthProviders are what we offer to dependency injection.
	AuthProviders = wire.NewSet(
		NewManager,
		ProvideUserAuthDataManager,
		ProvideArgon2Authenticator,
	)
)

func ProvideUserAuthDataManager(db database.DataManager) UserAuthDataManager {
	return db
}
