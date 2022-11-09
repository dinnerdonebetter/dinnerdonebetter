package fakes

import (
	fake "github.com/brianvoe/gofakeit/v5"

	"github.com/prixfixeco/backend/internal/pointers"
	"github.com/prixfixeco/backend/pkg/types"
)

// BuildFleshedOutQueryFilter builds a fully fleshed out QueryFilter.
func BuildFleshedOutQueryFilter() *types.QueryFilter {
	return &types.QueryFilter{
		Page:          pointers.Uint16(10),
		Limit:         pointers.Uint8(20),
		CreatedAfter:  pointers.Time(fake.Date()),
		CreatedBefore: pointers.Time(fake.Date()),
		UpdatedAfter:  pointers.Time(fake.Date()),
		UpdatedBefore: pointers.Time(fake.Date()),
		SortBy:        types.SortAscending,
	}
}
