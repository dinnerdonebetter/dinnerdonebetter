package fakes

import (
	"time"

	fake "github.com/brianvoe/gofakeit/v5"

	"github.com/prixfixeco/api_server/pkg/types"
)

// BuildFleshedOutQueryFilter builds a fully fleshed out QueryFilter.
func BuildFleshedOutQueryFilter() *types.QueryFilter {
	return &types.QueryFilter{
		Page:          func(x uint64) *uint64 { return &x }(10),
		Limit:         func(x uint8) *uint8 { return &x }(20),
		CreatedAfter:  func(x time.Time) *time.Time { return &x }(fake.Date()),
		CreatedBefore: func(x time.Time) *time.Time { return &x }(fake.Date()),
		UpdatedAfter:  func(x time.Time) *time.Time { return &x }(fake.Date()),
		UpdatedBefore: func(x time.Time) *time.Time { return &x }(fake.Date()),
		SortBy:        types.SortAscending,
	}
}
