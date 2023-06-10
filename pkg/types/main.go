package types

import (
	"fmt"
	"time"
)

const (
	// sortAscendingString is the pre-determined Ascending sortType for external use.
	sortAscendingString = "asc"
	// sortDescendingString is the pre-determined Descending sortType for external use.
	sortDescendingString = "desc"
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

	Identifiable interface {
		Identify() string
	}

	DatabaseRecord struct {
		_ struct{}

		CreatedAt     time.Time  `json:"createdAt"`
		ArchivedAt    *time.Time `json:"archivedAt"`
		LastUpdatedAt *time.Time `json:"lastUpdatedAt"`
		ID            string     `json:"id"`
	}

	DatabaseRecords []Identifiable

	APIMeta struct {
		_ struct{}

		UserID      string `json:"userID"`
		HouseholdID string `json:"householdID"`
	}

	APIResponse[T any] struct {
		_    struct{}
		Data *T      `json:"data"`
		Meta APIMeta `json:"meta"`
	}

	// APIError represents a response we might send to the User in the event of an error.
	APIError struct {
		_ struct{}

		Message string `json:"message"`
		Code    int    `json:"code"`
	}

	// NumberRange represents a range of numbers.
	NumberRange struct {
		_ struct{}

		Max *float32 `json:"max"`
		Min float32  `json:"min"`
	}

	NamedID struct {
		ID   string `json:"id,omitempty"`
		Name string `json:"name,omitempty"`
	}
)

func (e *APIError) Error() string {
	return fmt.Sprintf("%d: %s", e.Code, e.Message)
}

// Identify returns the ID of the DatabaseRecord.
func (r DatabaseRecord) Identify() string {
	return r.ID
}

// Len implements the sort.Sort interface.
func (m DatabaseRecords) Len() int {
	return len(m)
}

// Less implements the sort.Sort interface.
func (m DatabaseRecords) Less(i, j int) bool {
	return m[i].Identify() < m[j].Identify()
}

// Swap implements the sort.Sort interface.
func (m DatabaseRecords) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}
