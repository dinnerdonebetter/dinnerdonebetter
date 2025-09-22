package identity

import (
	"github.com/dinnerdonebetter/backend/internal/domain/identity"

	"github.com/google/wire"
)

var (
	Providers = wire.NewSet(
		ProvideIdentityRepository,
		ProvideUserDataManager,
	)
)

func ProvideUserDataManager(r identity.Repository) identity.UserDataManager {
	return r
}
