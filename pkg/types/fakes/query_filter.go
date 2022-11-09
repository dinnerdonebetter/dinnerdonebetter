package fakes

import (
	fake "github.com/brianvoe/gofakeit/v5"

	"github.com/prixfixeco/api_server/internal/pointers"
	"github.com/prixfixeco/api_server/pkg/types"
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
