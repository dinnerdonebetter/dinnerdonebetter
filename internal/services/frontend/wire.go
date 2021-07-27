package frontend

import (
	"gitlab.com/prixfixe/prixfixe/pkg/types"

	"github.com/google/wire"
)

var (
	// Providers is what we offer to dependency injection.
	Providers = wire.NewSet(
		ProvideService,
		ProvideAuthService,
		ProvideUsersService,
		ProvideValidIngredientsService,
		ProvideValidInstrumentsService,
		ProvideValidPreparationsService,
	)
)

// ProvideAuthService does what I hope one day wire figures out how to do.
func ProvideAuthService(x types.AuthService) AuthService {
	return x
}

// ProvideUsersService does what I hope one day wire figures out how to do.
func ProvideUsersService(x types.UserDataService) UsersService {
	return x
}

// ProvideValidIngredientsService does what I hope one day wire figures out how to do.
func ProvideValidIngredientsService(x types.ValidIngredientDataService) ValidIngredientsService {
	return x
}

// ProvideValidInstrumentsService does what I hope one day wire figures out how to do.
func ProvideValidInstrumentsService(x types.ValidInstrumentDataService) ValidInstrumentsService {
	return x
}

// ProvideValidPreparationsService does what I hope one day wire figures out how to do.
func ProvideValidPreparationsService(x types.ValidPreparationDataService) ValidPreparationsService {
	return x
}
