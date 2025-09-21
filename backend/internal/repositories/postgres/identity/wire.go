package identity

import (
	"github.com/dinnerdonebetter/backend/internal/domain/identity"

	"github.com/google/wire"
)

var (
	IDRepoProviders = wire.NewSet(
		ProvideIdentityRepository,
		ProvideUserDataManager,
	)
)

func ProvideUserDataManager(r identity.Repository) identity.UserDataManager {
	return r
}
