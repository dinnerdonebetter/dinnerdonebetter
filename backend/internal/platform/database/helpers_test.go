package database

import (
	"testing"

	"github.com/dinnerdonebetter/backend/internal/platform/database/filtering"
	"github.com/dinnerdonebetter/backend/internal/platform/pointer"

	"github.com/stretchr/testify/assert"
)

func TestFetchAllRows(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleData := []string{"one", "two", "three"}

		queryFilter := &filtering.QueryFilter{
			Page:  pointer.To(uint16(1)),
			Limit: pointer.To(uint8(1)),
		}

		actual, err := FetchAllRows(func(filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[string], error) {
			page := pointer.Dereference(filter.Page)

			return &filtering.QueryFilteredResult[string]{
				Data: []*string{&exampleData[page-1]},
				Pagination: filtering.Pagination{
					Page:          page + 1,
					Limit:         1,
					FilteredCount: 1,
					TotalCount:    uint64(len(exampleData)),
				},
			}, nil
		}, queryFilter)

		assert.Equal(t, exampleData, actual)
		assert.NoError(t, err)
	})
}
