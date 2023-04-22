package types

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
	"testing"
	"time"

	"github.com/prixfixeco/backend/internal/observability/logging"
	"github.com/prixfixeco/backend/internal/pkg/pointers"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQueryFilter_AttachToLogger(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		logger := logging.NewNoopLogger()

		qf := &QueryFilter{
			Page:            pointers.Pointer(uint16(100)),
			Limit:           pointers.Pointer(uint8(MaxLimit)),
			CreatedAfter:    pointers.Pointer(time.Now().Truncate(time.Second)),
			CreatedBefore:   pointers.Pointer(time.Now().Truncate(time.Second)),
			UpdatedAfter:    pointers.Pointer(time.Now().Truncate(time.Second)),
			UpdatedBefore:   pointers.Pointer(time.Now().Truncate(time.Second)),
			SortBy:          SortDescending,
			IncludeArchived: pointers.Pointer(true),
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
			Page:            pointers.Pointer(uint16(100)),
			Limit:           pointers.Pointer(uint8(MaxLimit)),
			CreatedAfter:    pointers.Pointer(tt),
			CreatedBefore:   pointers.Pointer(tt),
			UpdatedAfter:    pointers.Pointer(tt),
			UpdatedBefore:   pointers.Pointer(tt),
			SortBy:          SortDescending,
			IncludeArchived: pointers.Pointer(true),
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

		expected := uint16(123)
		qf := &QueryFilter{}
		qf.SetPage(&expected)

		assert.Equal(t, expected, *qf.Page)
	})
}

func TestQueryFilter_QueryPage(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()
		qf := &QueryFilter{
			Limit: pointers.Pointer(uint8(10)),
			Page:  pointers.Pointer(uint16(11)),
		}
		expected := uint16(100)
		actual := qf.QueryOffset()

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
			Page:            pointers.Pointer(uint16(100)),
			Limit:           pointers.Pointer(uint8(MaxLimit)),
			CreatedAfter:    pointers.Pointer(tt),
			CreatedBefore:   pointers.Pointer(tt),
			UpdatedAfter:    pointers.Pointer(tt),
			UpdatedBefore:   pointers.Pointer(tt),
			SortBy:          SortDescending,
			IncludeArchived: pointers.Pointer(true),
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
			Page:          pointers.Pointer(uint16(100)),
			Limit:         pointers.Pointer(uint8(MaxLimit)),
			CreatedAfter:  pointers.Pointer(tt),
			CreatedBefore: pointers.Pointer(tt),
			UpdatedAfter:  pointers.Pointer(tt),
			UpdatedBefore: pointers.Pointer(tt),
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
