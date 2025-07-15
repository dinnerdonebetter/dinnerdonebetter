package grocerylistpreparation

import (
	"github.com/google/wire"
)

var (
	ProvidersGroceryListPreparation = wire.NewSet(
		NewGroceryListCreator,
	)
)
