package types

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/prixfixeco/api_server/internal/observability/logging"
	"github.com/prixfixeco/api_server/internal/pointers"
)

func TestQueryFilter_AttachToLogger(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		logger := logging.NewNoopLogger()

		qf := &QueryFilter{
			Page:            pointers.Uint64(100),
			Limit:           pointers.Uint8(MaxLimit),
			CreatedAfter:    pointers.Time(time.Now().Truncate(time.Second)),
			CreatedBefore:   pointers.Time(time.Now().Truncate(time.Second)),
			UpdatedAfter:    pointers.Time(time.Now().Truncate(time.Second)),
			UpdatedBefore:   pointers.Time(time.Now().Truncate(time.Second)),
			SortBy:          SortDescending,
			IncludeArchived: boolPointer(true),
		}

		assert.NotNil(t, qf.AttachToLogger(logger))
	})

	T.Run("with nil", func(t *testing.T) {
		t.Parallel()

		logger := logging.NewNoopLogger()

		assert.NotNil(t, (*QueryFilter)(nil).AttachToLogger(logger))
	})
}

func TestQueryFilter_FromParams(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		tt, err := time.Parse(time.RFC3339Nano, time.Now().UTC().Truncate(time.Second).Format(time.RFC3339Nano))
		require.NoError(t, err)

		actual := &QueryFilter{}
		expected := &QueryFilter{
			Page:            func(x uint64) *uint64 { return &x }(100),
			Limit:           func(x uint8) *uint8 { return &x }(MaxLimit),
			CreatedAfter:    pointers.Time(tt),
			CreatedBefore:   pointers.Time(tt),
			UpdatedAfter:    pointers.Time(tt),
			UpdatedBefore:   pointers.Time(tt),
			SortBy:          SortDescending,
			IncludeArchived: boolPointer(true),
		}

		exampleInput := url.Values{
			pageQueryKey:            []string{strconv.Itoa(int(*expected.Page))},
			LimitQueryKey:           []string{strconv.Itoa(int(*expected.Limit))},
			createdBeforeQueryKey:   []string{expected.CreatedAfter.Format(time.RFC3339Nano)},
			createdAfterQueryKey:    []string{expected.CreatedBefore.Format(time.RFC3339Nano)},
			updatedBeforeQueryKey:   []string{expected.UpdatedAfter.Format(time.RFC3339Nano)},
			updatedAfterQueryKey:    []string{expected.UpdatedBefore.Format(time.RFC3339Nano)},
			sortByQueryKey:          []string{*expected.SortBy},
			includeArchivedQueryKey: []string{strconv.FormatBool(true)},
		}

		actual.FromParams(exampleInput)
		actual.CreatedAfter.Location()

		assert.Equal(t, expected, actual)

		exampleInput[sortByQueryKey] = []string{*SortAscending}

		actual.FromParams(exampleInput)
		assert.Equal(t, SortAscending, actual.SortBy)
	})
}

func TestQueryFilter_SetPage(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		expected := uint64(123)
		qf := &QueryFilter{}
		qf.SetPage(&expected)

		assert.Equal(t, expected, *qf.Page)
	})
}

func TestQueryFilter_QueryPage(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()
		qf := &QueryFilter{Limit: func(x uint8) *uint8 { return &x }(10), Page: func(x uint64) *uint64 { return &x }(11)}
		expected := uint64(100)
		actual := qf.QueryPage()

		assert.Equal(t, expected, actual)
	})
}

func TestQueryFilter_ToValues(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		tt, err := time.Parse(time.RFC3339Nano, time.Now().UTC().Truncate(time.Second).Format(time.RFC3339Nano))
		require.NoError(t, err)

		qf := &QueryFilter{
			Page:            func(x uint64) *uint64 { return &x }(100),
			Limit:           func(x uint8) *uint8 { return &x }(MaxLimit),
			CreatedAfter:    pointers.Time(tt),
			CreatedBefore:   pointers.Time(tt),
			UpdatedAfter:    pointers.Time(tt),
			UpdatedBefore:   pointers.Time(tt),
			SortBy:          SortDescending,
			IncludeArchived: boolPointer(true),
		}

		expected := url.Values{
			pageQueryKey:            []string{strconv.Itoa(int(*qf.Page))},
			LimitQueryKey:           []string{strconv.Itoa(int(*qf.Limit))},
			createdBeforeQueryKey:   []string{qf.CreatedAfter.Format(time.RFC3339Nano)},
			createdAfterQueryKey:    []string{qf.CreatedBefore.Format(time.RFC3339Nano)},
			updatedBeforeQueryKey:   []string{qf.UpdatedAfter.Format(time.RFC3339Nano)},
			updatedAfterQueryKey:    []string{qf.UpdatedBefore.Format(time.RFC3339Nano)},
			includeArchivedQueryKey: []string{strconv.FormatBool(*qf.IncludeArchived)},
			sortByQueryKey:          []string{*qf.SortBy},
		}

		actual := qf.ToValues()
		assert.Equal(t, expected, actual)
	})

	T.Run("with nil", func(t *testing.T) {
		t.Parallel()
		qf := (*QueryFilter)(nil)
		expected := DefaultQueryFilter().ToValues()
		actual := qf.ToValues()
		assert.Equal(t, expected, actual)
	})
}

func TestExtractQueryFilter(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		tt, err := time.Parse(time.RFC3339Nano, time.Now().UTC().Truncate(time.Second).Format(time.RFC3339Nano))
		require.NoError(t, err)

		expected := &QueryFilter{
			Page:          func(x uint64) *uint64 { return &x }(100),
			Limit:         func(x uint8) *uint8 { return &x }(MaxLimit),
			CreatedAfter:  pointers.Time(tt),
			CreatedBefore: pointers.Time(tt),
			UpdatedAfter:  pointers.Time(tt),
			UpdatedBefore: pointers.Time(tt),
			SortBy:        SortDescending,
		}
		exampleInput := url.Values{
			pageQueryKey:          []string{strconv.Itoa(int(*expected.Page))},
			LimitQueryKey:         []string{strconv.Itoa(int(*expected.Limit))},
			createdBeforeQueryKey: []string{expected.CreatedAfter.Format(time.RFC3339Nano)},
			createdAfterQueryKey:  []string{expected.CreatedBefore.Format(time.RFC3339Nano)},
			updatedBeforeQueryKey: []string{expected.UpdatedAfter.Format(time.RFC3339Nano)},
			updatedAfterQueryKey:  []string{expected.UpdatedBefore.Format(time.RFC3339Nano)},
			sortByQueryKey:        []string{*expected.SortBy},
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://verygoodsoftwarenotvirus.ru", http.NoBody)
		assert.NoError(t, err)
		require.NotNil(t, req)

		req.URL.RawQuery = exampleInput.Encode()
		actual := ExtractQueryFilterFromRequest(req)
		assert.Equal(t, expected, actual)
	})
}
