package types

import (
	"database/sql"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/prixfixeco/api_server/internal/observability/keys"
	"github.com/prixfixeco/api_server/internal/observability/logging"
	"github.com/prixfixeco/api_server/internal/pointers"
)

const (
	// MaxLimit is the maximum value for list queries.
	MaxLimit = 250
	// DefaultLimit represents how many results we return in a response by default.
	DefaultLimit = 20

	// SearchQueryKey is the query param key to find search queries in requests.
	SearchQueryKey = "q"
	// LimitQueryKey is the query param key to specify a limit in a query.
	LimitQueryKey = "limit"

	pageQueryKey            = "page"
	createdBeforeQueryKey   = "createdBefore"
	createdAfterQueryKey    = "createdAfter"
	updatedBeforeQueryKey   = "updatedBefore"
	updatedAfterQueryKey    = "updatedAfter"
	includeArchivedQueryKey = "includeArchived"
	sortByQueryKey          = "sortBy"
)

// QueryFilter represents all the filters a User could apply to a list query.
type QueryFilter struct {
	_ struct{}

	SortBy          *string    `json:"sortBy"`
	Page            *uint64    `json:"page"`
	CreatedAfter    *time.Time `json:"createdBefore,omitempty"`
	CreatedBefore   *time.Time `json:"createdAfter,omitempty"`
	UpdatedAfter    *time.Time `json:"updatedBefore,omitempty"`
	UpdatedBefore   *time.Time `json:"updatedAfter,omitempty"`
	Limit           *uint8     `json:"limit"`
	IncludeArchived *bool      `json:"includeArchived,omitempty"`
}

// DefaultQueryFilter builds the default query filter.
func DefaultQueryFilter() *QueryFilter {
	return &QueryFilter{
		Page:   pointers.Uint64(1),
		Limit:  pointers.Uint8(DefaultLimit),
		SortBy: SortAscending,
	}
}

type QueryFilterDatabaseArgs struct {
	CreatedAfter  sql.NullTime
	CreatedBefore sql.NullTime
	UpdatedAfter  sql.NullTime
	UpdatedBefore sql.NullTime
}

func (qf *QueryFilter) ToDatabaseArgs() QueryFilterDatabaseArgs {
	args := QueryFilterDatabaseArgs{
		CreatedAfter:  sql.NullTime{},
		CreatedBefore: sql.NullTime{},
		UpdatedAfter:  sql.NullTime{},
		UpdatedBefore: sql.NullTime{},
	}

	if qf.CreatedAfter != nil {
		args.CreatedAfter = sql.NullTime{Time: *qf.CreatedAfter}
	}

	if qf.CreatedBefore != nil {
		args.CreatedBefore = sql.NullTime{Time: *qf.CreatedBefore}
	}

	if qf.UpdatedAfter != nil {
		args.UpdatedAfter = sql.NullTime{Time: *qf.UpdatedAfter}
	}

	if qf.UpdatedAfter != nil {
		args.UpdatedAfter = sql.NullTime{Time: *qf.UpdatedAfter}
	}

	return args
}

// AttachToLogger attaches a QueryFilter's values to a logging.Logger.
func (qf *QueryFilter) AttachToLogger(logger logging.Logger) logging.Logger {
	l := logging.EnsureLogger(logger).Clone()

	if qf == nil {
		return l.WithValue(keys.FilterIsNilKey, true)
	}

	if qf.Page != nil {
		l = l.WithValue(pageQueryKey, qf.Page)
	}

	if qf.Limit != nil {
		l = l.WithValue(LimitQueryKey, qf.Limit)
	}

	if qf.SortBy != nil {
		l = l.WithValue(sortByQueryKey, qf.SortBy)
	}

	if qf.CreatedBefore != nil {
		l = l.WithValue(createdBeforeQueryKey, qf.CreatedBefore)
	}

	if qf.CreatedAfter != nil {
		l = l.WithValue(createdAfterQueryKey, qf.CreatedAfter)
	}

	if qf.UpdatedBefore != nil {
		l = l.WithValue(updatedBeforeQueryKey, qf.UpdatedBefore)
	}

	if qf.UpdatedAfter != nil {
		l = l.WithValue(updatedAfterQueryKey, qf.UpdatedAfter)
	}

	return l
}

// FromParams overrides the core QueryFilter values with values retrieved from url.Params.
func (qf *QueryFilter) FromParams(params url.Values) {
	if i, err := strconv.ParseUint(params.Get(pageQueryKey), 10, 64); err == nil {
		qf.Page = uint64Pointer(uint64(math.Max(float64(i), 1)))
	}

	if i, err := strconv.ParseUint(params.Get(LimitQueryKey), 10, 64); err == nil {
		qf.Limit = uint8Pointer(uint8(math.Min(math.Max(float64(i), 0), MaxLimit)))
	}

	if t, err := time.Parse(time.RFC3339, params.Get(createdBeforeQueryKey)); err == nil {
		qf.CreatedBefore = &t
	}

	if t, err := time.Parse(time.RFC3339, params.Get(createdAfterQueryKey)); err == nil {
		qf.CreatedAfter = &t
	}

	if t, err := time.Parse(time.RFC3339, params.Get(updatedBeforeQueryKey)); err == nil {
		qf.UpdatedBefore = &t
	}

	if t, err := time.Parse(time.RFC3339, params.Get(updatedAfterQueryKey)); err == nil {
		qf.UpdatedAfter = &t
	}

	if i, err := strconv.ParseBool(params.Get(includeArchivedQueryKey)); err == nil {
		qf.IncludeArchived = &i
	}

	switch strings.ToLower(params.Get(sortByQueryKey)) {
	case "asc":
		qf.SortBy = SortAscending
	case "desc":
		qf.SortBy = SortDescending
	}
}

// SetPage sets the current page with certain constraints.
func (qf *QueryFilter) SetPage(page *uint64) {
	if page != nil {
		qf.Page = uint64Pointer(uint64(math.Max(1, float64(*page))))
	}
}

// QueryPage calculates a query page from the current filter values.
func (qf *QueryFilter) QueryPage() uint64 {
	if qf.Limit != nil && qf.Page != nil {
		return uint64(*qf.Limit) * (*qf.Page - 1)
	}
	return 0
}

// ToValues returns a url.Values from a QueryFilter.
func (qf *QueryFilter) ToValues() url.Values {
	if qf == nil {
		return DefaultQueryFilter().ToValues()
	}

	v := url.Values{}

	if qf.Page != nil {
		v.Set(pageQueryKey, strconv.FormatUint(*qf.Page, 10))
	}

	if qf.Limit != nil {
		v.Set(LimitQueryKey, strconv.FormatUint(uint64(*qf.Limit), 10))
	}

	if qf.SortBy != nil {
		v.Set(sortByQueryKey, *qf.SortBy)
	}

	if qf.CreatedBefore != nil {
		v.Set(createdBeforeQueryKey, qf.CreatedBefore.Format(time.RFC3339))
	}

	if qf.CreatedAfter != nil {
		v.Set(createdAfterQueryKey, qf.CreatedAfter.Format(time.RFC3339))
	}

	if qf.UpdatedBefore != nil {
		v.Set(updatedBeforeQueryKey, qf.UpdatedBefore.Format(time.RFC3339))
	}

	if qf.UpdatedAfter != nil {
		v.Set(updatedAfterQueryKey, qf.UpdatedAfter.Format(time.RFC3339))
	}

	if qf.IncludeArchived != nil {
		v.Set(includeArchivedQueryKey, strconv.FormatBool(*qf.IncludeArchived))
	}

	return v
}

// ExtractQueryFilterFromRequest can extract a QueryFilter from a request.
func ExtractQueryFilterFromRequest(req *http.Request) *QueryFilter {
	qf := &QueryFilter{}
	qf.FromParams(req.URL.Query())

	if qf.Page != nil {
		if *qf.Page == 0 {
			qf.Page = uint64Pointer(1)
		}
	}

	if qf.Limit != nil {
		if *qf.Limit == 0 {
			qf.Limit = uint8Pointer(DefaultLimit)
		}
	}

	return qf
}
