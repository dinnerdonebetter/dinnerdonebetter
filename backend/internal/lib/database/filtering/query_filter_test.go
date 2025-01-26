package filtering

import (
	"context"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"testing"
	"time"

	"github.com/dinnerdonebetter/backend/internal/lib/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/lib/pointer"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQueryFilter_AttachToLogger(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		logger := logging.NewNoopLogger()

		qf := &QueryFilter{
			Query:           t.Name(),
			Page:            pointer.To(uint16(100)),
			Limit:           pointer.To(uint8(MaxQueryFilterLimit)),
			CreatedAfter:    pointer.To(time.Now().Truncate(time.Second)),
			CreatedBefore:   pointer.To(time.Now().Truncate(time.Second)),
			UpdatedAfter:    pointer.To(time.Now().Truncate(time.Second)),
			UpdatedBefore:   pointer.To(time.Now().Truncate(time.Second)),
			SortBy:          SortDescending,
			IncludeArchived: pointer.To(true),
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
			Query:           t.Name(),
			Page:            pointer.To(uint16(100)),
			Limit:           pointer.To(uint8(MaxQueryFilterLimit)),
			CreatedAfter:    pointer.To(tt),
			CreatedBefore:   pointer.To(tt),
			UpdatedAfter:    pointer.To(tt),
			UpdatedBefore:   pointer.To(tt),
			SortBy:          SortDescending,
			IncludeArchived: pointer.To(true),
		}

		exampleInput := url.Values{
			QueryKeySearch:          []string{t.Name()},
			QueryKeyPage:            []string{strconv.Itoa(int(*expected.Page))},
			QueryKeyLimit:           []string{strconv.Itoa(int(*expected.Limit))},
			QueryKeyCreatedBefore:   []string{expected.CreatedAfter.Format(time.RFC3339Nano)},
			QueryKeyCreatedAfter:    []string{expected.CreatedBefore.Format(time.RFC3339Nano)},
			QueryKeyUpdatedBefore:   []string{expected.UpdatedAfter.Format(time.RFC3339Nano)},
			QueryKeyUpdatedAfter:    []string{expected.UpdatedBefore.Format(time.RFC3339Nano)},
			QueryKeySortBy:          []string{*expected.SortBy},
			QueryKeyIncludeArchived: []string{strconv.FormatBool(true)},
		}

		actual.FromParams(exampleInput)

		assert.Equal(t, expected, actual)

		exampleInput[QueryKeySortBy] = []string{*SortAscending}

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
			Limit: pointer.To(uint8(10)),
			Page:  pointer.To(uint16(11)),
		}
		expected := uint16(100)
		actual := qf.QueryOffset()

		assert.Equal(t, expected, actual)
	})

	T.Run("with nil values", func(t *testing.T) {
		t.Parallel()

		qf := &QueryFilter{}
		expected := uint16(0)
		actual := qf.QueryOffset()

		assert.Equal(t, expected, actual)
	})

	T.Run("with max values", func(t *testing.T) {
		t.Parallel()

		qf := &QueryFilter{
			Page:            pointer.To(uint16(0)),
			Limit:           pointer.To(uint8(math.MaxUint8)),
			IncludeArchived: pointer.To(true),
		}
		expected := uint16(0)
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
			Query:           t.Name(),
			Page:            pointer.To(uint16(100)),
			Limit:           pointer.To(uint8(MaxQueryFilterLimit)),
			CreatedAfter:    pointer.To(tt),
			CreatedBefore:   pointer.To(tt),
			UpdatedAfter:    pointer.To(tt),
			UpdatedBefore:   pointer.To(tt),
			SortBy:          SortDescending,
			IncludeArchived: pointer.To(true),
		}

		expected := url.Values{
			QueryKeySearch:          []string{t.Name()},
			QueryKeyPage:            []string{strconv.Itoa(int(*qf.Page))},
			QueryKeyLimit:           []string{strconv.Itoa(int(*qf.Limit))},
			QueryKeyCreatedBefore:   []string{qf.CreatedAfter.Format(time.RFC3339Nano)},
			QueryKeyCreatedAfter:    []string{qf.CreatedBefore.Format(time.RFC3339Nano)},
			QueryKeyUpdatedBefore:   []string{qf.UpdatedAfter.Format(time.RFC3339Nano)},
			QueryKeyUpdatedAfter:    []string{qf.UpdatedBefore.Format(time.RFC3339Nano)},
			QueryKeyIncludeArchived: []string{strconv.FormatBool(*qf.IncludeArchived)},
			QueryKeySortBy:          []string{*qf.SortBy},
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
			Query:         t.Name(),
			Page:          pointer.To(uint16(100)),
			Limit:         pointer.To(uint8(MaxQueryFilterLimit)),
			CreatedAfter:  pointer.To(tt),
			CreatedBefore: pointer.To(tt),
			UpdatedAfter:  pointer.To(tt),
			UpdatedBefore: pointer.To(tt),
			SortBy:        SortDescending,
		}
		exampleInput := url.Values{
			QueryKeySearch:        []string{t.Name()},
			QueryKeyPage:          []string{strconv.Itoa(int(*expected.Page))},
			QueryKeyLimit:         []string{strconv.Itoa(int(*expected.Limit))},
			QueryKeyCreatedBefore: []string{expected.CreatedAfter.Format(time.RFC3339Nano)},
			QueryKeyCreatedAfter:  []string{expected.CreatedBefore.Format(time.RFC3339Nano)},
			QueryKeyUpdatedBefore: []string{expected.UpdatedAfter.Format(time.RFC3339Nano)},
			QueryKeyUpdatedAfter:  []string{expected.UpdatedBefore.Format(time.RFC3339Nano)},
			QueryKeySortBy:        []string{*expected.SortBy},
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://verygoodsoftwarenotvirus.ru", http.NoBody)
		assert.NoError(t, err)
		require.NotNil(t, req)

		req.URL.RawQuery = exampleInput.Encode()
		actual := ExtractQueryFilterFromRequest(req)
		assert.Equal(t, expected, actual)
	})

	T.Run("with missing values", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		expected := &QueryFilter{
			Page:   pointer.To(uint16(1)),
			Limit:  pointer.To(uint8(DefaultQueryFilterLimit)),
			SortBy: SortAscending,
		}
		exampleInput := url.Values{
			QueryKeyPage:  []string{"0"},
			QueryKeyLimit: []string{"0"},
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://verygoodsoftwarenotvirus.ru", http.NoBody)
		assert.NoError(t, err)
		require.NotNil(t, req)

		req.URL.RawQuery = exampleInput.Encode()
		actual := ExtractQueryFilterFromRequest(req)
		assert.Equal(t, expected, actual)
	})
}

func TestQueryFilter_ToPagination(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		qf := &QueryFilter{
			Page:  pointer.To(uint16(100)),
			Limit: pointer.To(uint8(MaxQueryFilterLimit)),
		}

		expected := Pagination{
			Page:  *qf.Page,
			Limit: *qf.Limit,
		}

		actual := qf.ToPagination()
		assert.Equal(t, expected, actual)
	})

	T.Run("with nil value", func(t *testing.T) {
		t.Parallel()

		qf := (*QueryFilter)(nil)

		actual := qf.ToPagination()
		assert.NotNil(t, actual)
	})
}
