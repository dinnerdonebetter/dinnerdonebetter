package apiclient

import (
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/dinnerdonebetter/backend/internal/lib/pointer"
	textsearch "github.com/dinnerdonebetter/backend/internal/lib/search/text"
)

const (
	// sortAscendingString is the pre-determined Ascending sortType for external use.
	sortAscendingString = "asc"
	// sortDescendingString is the pre-determined Descending sortType for external use.
	sortDescendingString = "desc"
)

var (
	// SortAscending is the pre-determined Ascending string for external use.
	SortAscending = pointer.To(sortAscendingString)
	// SortDescending is the pre-determined Descending string for external use.
	SortDescending = pointer.To(sortDescendingString)
)

const (
	// MaxQueryFilterLimit is the maximum value for list queries.
	MaxQueryFilterLimit = 250
	// DefaultQueryFilterLimit represents how many results we return in a response by default.
	DefaultQueryFilterLimit = 50

	// QueryKeySearchWithDatabase is the query param key to find search queries in requests.
	QueryKeySearchWithDatabase = "useDB"

	// QueryKeyLimit is the query param key to specify a limit in a query.
	QueryKeyLimit = "limit"
	// QueryKeyPage is the query param key for specifying which page the user would in a list query.
	QueryKeyPage = "page"
	// QueryKeyCreatedBefore is the query param key for a creation time limit in a list query.
	QueryKeyCreatedBefore = "createdBefore"
	// QueryKeyCreatedAfter is the query param key for a creation time limit in a list query.
	QueryKeyCreatedAfter = "createdAfter"
	// QueryKeyUpdatedBefore is the query param key for a creation time limit in a list query.
	QueryKeyUpdatedBefore = "updatedBefore"
	// QueryKeyUpdatedAfter is the query param key for a creation time limit in a list query.
	QueryKeyUpdatedAfter = "updatedAfter"
	// QueryKeyIncludeArchived is the query param key for including archived results in a query.
	QueryKeyIncludeArchived = "includeArchived"
	// QueryKeySortBy is the query param key for sort order in a query.
	QueryKeySortBy = "sortBy"
)

// Pagination represents a pagination request.
type Pagination struct {
	_ struct{} `json:"-"`

	Page          uint16 `json:"page"`
	Limit         uint8  `json:"limit"`
	FilteredCount uint64 `json:"filteredCount"`
	TotalCount    uint64 `json:"totalCount"`
}

type QueryFilteredResult[T any] struct {
	_ struct{} `json:"-"`

	Data []*T `json:"data"`
	Pagination
}

// QueryFilter represents all the filters a User could apply to a list query.
type QueryFilter struct {
	_ struct{} `json:"-"`

	SortBy          *string    `json:"sortBy,omitempty"`
	Page            *uint16    `json:"page,omitempty"`
	CreatedAfter    *time.Time `json:"createdBefore,omitempty"`
	CreatedBefore   *time.Time `json:"createdAfter,omitempty"`
	UpdatedAfter    *time.Time `json:"updatedBefore,omitempty"`
	UpdatedBefore   *time.Time `json:"updatedAfter,omitempty"`
	Limit           *uint8     `json:"limit,omitempty"`
	IncludeArchived *bool      `json:"includeArchived,omitempty"`
	Query           string     `json:"q,omitempty"`
}

// DefaultQueryFilter builds the default query filter.
func DefaultQueryFilter() *QueryFilter {
	return &QueryFilter{
		Page:   pointer.To(uint16(1)),
		Limit:  pointer.To(uint8(DefaultQueryFilterLimit)),
		SortBy: SortAscending,
	}
}

// QueryOffset calculates a query page from the current filter values.
func (qf *QueryFilter) QueryOffset() uint16 {
	if qf != nil && qf.Limit != nil && qf.Page != nil {
		page := *qf.Page
		if page == 0 {
			page = 1
		}

		return uint16(*qf.Limit) * (page - 1)
	}
	return 0
}

// ToValues returns a url.Values from a QueryFilter.
func (qf *QueryFilter) ToValues() url.Values {
	if qf == nil {
		return DefaultQueryFilter().ToValues()
	}

	v := url.Values{}

	if qf.Query != "" {
		v.Set(textsearch.QueryKeySearch, qf.Query)
	}

	if qf.Page != nil {
		v.Set(QueryKeyPage, strconv.FormatUint(uint64(*qf.Page), 10))
	}

	if qf.Limit != nil {
		v.Set(QueryKeyLimit, strconv.FormatUint(uint64(*qf.Limit), 10))
	}

	if qf.SortBy != nil {
		v.Set(QueryKeySortBy, *qf.SortBy)
	}

	if qf.CreatedBefore != nil {
		v.Set(QueryKeyCreatedBefore, qf.CreatedBefore.Format(time.RFC3339Nano))
	}

	if qf.CreatedAfter != nil {
		v.Set(QueryKeyCreatedAfter, qf.CreatedAfter.Format(time.RFC3339Nano))
	}

	if qf.UpdatedBefore != nil {
		v.Set(QueryKeyUpdatedBefore, qf.UpdatedBefore.Format(time.RFC3339Nano))
	}

	if qf.UpdatedAfter != nil {
		v.Set(QueryKeyUpdatedAfter, qf.UpdatedAfter.Format(time.RFC3339Nano))
	}

	if qf.IncludeArchived != nil {
		v.Set(QueryKeyIncludeArchived, strconv.FormatBool(*qf.IncludeArchived))
	}

	return v
}

// ToPagination returns a Pagination from a QueryFilter.
func (qf *QueryFilter) ToPagination() Pagination {
	if qf == nil {
		return DefaultQueryFilter().ToPagination()
	}

	x := Pagination{}

	if qf.Page != nil {
		x.Page = *qf.Page
	}

	if qf.Limit != nil {
		x.Limit = *qf.Limit
	}

	return x
}

// FromParams overrides the core QueryFilter values with values retrieved from url.Params.
func (qf *QueryFilter) FromParams(params url.Values) {
	if i := params.Get(textsearch.QueryKeySearch); i != "" {
		qf.Query = i
	}

	if i, err := strconv.ParseUint(params.Get(QueryKeyPage), 10, 64); err == nil {
		qf.Page = pointer.To(uint16(math.Max(float64(i), 1)))
	}

	if i, err := strconv.ParseUint(params.Get(QueryKeyLimit), 10, 64); err == nil {
		qf.Limit = pointer.To(uint8(math.Min(math.Max(float64(i), 0), MaxQueryFilterLimit)))
	}

	if t, err := time.Parse(time.RFC3339Nano, params.Get(QueryKeyCreatedBefore)); err == nil {
		qf.CreatedBefore = &t
	}

	if t, err := time.Parse(time.RFC3339Nano, params.Get(QueryKeyCreatedAfter)); err == nil {
		qf.CreatedAfter = &t
	}

	if t, err := time.Parse(time.RFC3339Nano, params.Get(QueryKeyUpdatedBefore)); err == nil {
		qf.UpdatedBefore = &t
	}

	if t, err := time.Parse(time.RFC3339Nano, params.Get(QueryKeyUpdatedAfter)); err == nil {
		qf.UpdatedAfter = &t
	}

	if i, err := strconv.ParseBool(params.Get(QueryKeyIncludeArchived)); err == nil {
		qf.IncludeArchived = &i
	}

	switch strings.ToLower(params.Get(QueryKeySortBy)) {
	case sortAscendingString:
		qf.SortBy = SortAscending
	case sortDescendingString:
		qf.SortBy = SortDescending
	}
}

// ExtractQueryFilterFromRequest can extract a QueryFilter from a request.
func ExtractQueryFilterFromRequest(req *http.Request) *QueryFilter {
	qf := DefaultQueryFilter()
	qf.FromParams(req.URL.Query())

	if qf.Limit != nil {
		if *qf.Limit == 0 {
			qf.Limit = pointer.To(uint8(DefaultQueryFilterLimit))
		}
	}

	return qf
}
