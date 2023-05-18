package fakes

import (
	"github.com/dinnerdonebetter/backend/internal/pkg/pointers"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

// BuildFleshedOutQueryFilter builds a fully fleshed out QueryFilter.
func BuildFleshedOutQueryFilter() *types.QueryFilter {
	return &types.QueryFilter{
		Page:          pointers.Pointer(uint16(10)),
		Limit:         pointers.Pointer(uint8(20)),
		CreatedAfter:  pointers.Pointer(BuildFakeTime()),
		CreatedBefore: pointers.Pointer(BuildFakeTime()),
		UpdatedAfter:  pointers.Pointer(BuildFakeTime()),
		UpdatedBefore: pointers.Pointer(BuildFakeTime()),
		SortBy:        types.SortAscending,
	}
}
