package database

import (
	"fmt"
	"math"

	"github.com/dinnerdonebetter/backend/internal/platform/database/filtering"
	"github.com/dinnerdonebetter/backend/internal/platform/pointer"
)

func FetchAllRows[T any](fetchFunc func(*filtering.QueryFilter) (*filtering.QueryFilteredResult[T], error), initialFilter *filtering.QueryFilter) ([]T, error) {
	var done bool
	allData := []T{}

	var filter = initialFilter
	if filter == nil {
		filter = &filtering.QueryFilter{
			Page:            pointer.To(uint16(1)),
			Limit:           pointer.To(uint8(math.MaxUint8)),
			IncludeArchived: pointer.To(true),
		}
	}

	for !done {
		data, err := fetchFunc(filter)
		if err != nil {
			return nil, fmt.Errorf("getting data: %w", err)
		}

		for _, x := range data.Data {
			if x != nil {
				allData = append(allData, *x)
			}
		}

		if data.TotalCount <= uint64(len(allData)) {
			done = true
		}
		filter.Page = pointer.To(*filter.Page + 1)
	}

	return allData, nil
}
