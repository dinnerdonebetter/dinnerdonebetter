package types

import (
	"fmt"
)

const (
	// sortAscendingString is the pre-determined Ascending sortType for external use.
	sortAscendingString = "asc"
	// sortDescendingString is the pre-determined Descending sortType for external use.
	sortDescendingString = "desc"
	// QueryFilterRequestContextKey is the non-string type we use for referencing QueryFilter objects.
	QueryFilterRequestContextKey ContextKey = "query_filter"
)

var (
	// SortAscending is the pre-determined Ascending string for external use.
	SortAscending = func(x string) *string { return &x }(sortAscendingString)
	// SortDescending is the pre-determined Descending string for external use.
	SortDescending = func(x string) *string { return &x }(sortDescendingString)
)

type (
	// ContextKey represents strings to be used in Context objects. From the docs:
	// 	"The provided key must be comparable and should not be of type string or
	// 	 any other built-in type to avoid collisions between packages using context."
	ContextKey string

	// Pagination represents a pagination request.
	Pagination struct {
		_ struct{}

		Page          uint16 `json:"page"`
		Limit         uint8  `json:"limit"`
		FilteredCount uint64 `json:"filteredCount"`
		TotalCount    uint64 `json:"totalCount"`
	}

	QueryFilteredResult[T any] struct {
		Data []*T `json:"data"`
		Pagination
	}

	// ErrorResponse represents a response we might send to the User in the event of an error.
	ErrorResponse struct {
		_ struct{}

		Message string `json:"message"`
		Code    int    `json:"code"`
	}
)

var _ error = (*ErrorResponse)(nil)

func (er *ErrorResponse) Error() string {
	return fmt.Sprintf("%d: %s", er.Code, er.Message)
}
