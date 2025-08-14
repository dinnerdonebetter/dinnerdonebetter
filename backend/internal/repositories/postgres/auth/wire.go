package auth

import (
	"github.com/dinnerdonebetter/backend/internal/domain/auth"
	"github.com/google/wire"
)

var (
	Providers = wire.NewSet(
		ProvideAuthRepository,
		ProvidePasswordResetTokenDataManager,
	)
)

func ProvidePasswordResetTokenDataManager(r auth.Repository) auth.PasswordResetTokenDataManager {
	return r
}
