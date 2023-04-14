package fakes

import (
	"github.com/prixfixeco/backend/internal/pkg/pointers"
	"github.com/prixfixeco/backend/pkg/types"
)

// BuildFleshedOutQueryFilter builds a fully fleshed out QueryFilter.
func BuildFleshedOutQueryFilter() *types.QueryFilter {
	return &types.QueryFilter{
		Page:          pointers.Uint16(10),
		Limit:         pointers.Uint8(20),
		CreatedAfter:  pointers.Time(BuildFakeTime()),
		CreatedBefore: pointers.Time(BuildFakeTime()),
		UpdatedAfter:  pointers.Time(BuildFakeTime()),
		UpdatedBefore: pointers.Time(BuildFakeTime()),
		SortBy:        types.SortAscending,
	}
}
