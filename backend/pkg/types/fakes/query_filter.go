package fakes

import (
	"github.com/dinnerdonebetter/backend/internal/lib/database/filtering"
	"github.com/dinnerdonebetter/backend/internal/lib/pointer"
)

// BuildFleshedOutQueryFilter builds a fully fleshed out QueryFilter.
func BuildFleshedOutQueryFilter() *filtering.QueryFilter {
	return &filtering.QueryFilter{
		Page:          pointer.To(uint16(10)),
		Limit:         pointer.To(uint8(20)),
		CreatedAfter:  pointer.To(BuildFakeTime()),
		CreatedBefore: pointer.To(BuildFakeTime()),
		UpdatedAfter:  pointer.To(BuildFakeTime()),
		UpdatedBefore: pointer.To(BuildFakeTime()),
		SortBy:        filtering.SortAscending,
	}
}
