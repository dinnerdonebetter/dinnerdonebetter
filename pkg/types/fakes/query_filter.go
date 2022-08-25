package fakes

import (
	fake "github.com/brianvoe/gofakeit/v5"

	"github.com/prixfixeco/api_server/pkg/types"
)

// BuildFleshedOutQueryFilter builds a fully fleshed out QueryFilter.
func BuildFleshedOutQueryFilter() *types.QueryFilter {
	return &types.QueryFilter{
		Page:          func(x uint64) *uint64 { return &x }(10),
		Limit:         func(x uint8) *uint8 { return &x }(20),
		CreatedAfter:  func(x uint64) *uint64 { return &x }(uint64(uint32(fake.Date().Unix()))),
		CreatedBefore: func(x uint64) *uint64 { return &x }(uint64(uint32(fake.Date().Unix()))),
		UpdatedAfter:  func(x uint64) *uint64 { return &x }(uint64(uint32(fake.Date().Unix()))),
		UpdatedBefore: func(x uint64) *uint64 { return &x }(uint64(uint32(fake.Date().Unix()))),
		SortBy:        types.SortAscending,
	}
}
