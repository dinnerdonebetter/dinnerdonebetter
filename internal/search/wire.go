package search

import (
	"github.com/dinnerdonebetter/backend/pkg/types"

	"github.com/google/wire"
)

var (
	Providers = wire.NewSet(
		ValidIngredientSearchSubsetIndexSearcherFromIndex,
	)
)

func ValidIngredientSearchSubsetIndexSearcherFromIndex(index Index[types.ValidIngredientSearchSubset]) IndexSearcher[types.ValidIngredientSearchSubset] {
	return index
}
