package chi

import (
	"context"
	"strconv"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	testutils "github.com/dinnerdonebetter/backend/tests/utils"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func TestNewRouteParamManager(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		assert.NotNil(t, NewRouteParamManager())
	})
}

func Test_BuildRouteParamIDFetcher(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		r := &chiRouteParamManager{}

		ctx := context.Background()
		exampleKey := "blah"
		fn := r.BuildRouteParamIDFetcher(logging.NewNoopLogger(), exampleKey, "thing")
		expected := uint64(123)
		req := testutils.BuildTestRequest(t).WithContext(
			context.WithValue(
				ctx,
				chi.RouteCtxKey,
				&chi.Context{
					URLParams: chi.RouteParams{
						Keys:   []string{exampleKey},
						Values: []string{strconv.FormatUint(expected, 10)},
					},
				},
			),
		)

		actual := fn(req)
		assert.Equal(t, expected, actual)
	})

	T.Run("with invalid value somehow", func(t *testing.T) {
		// NOTE: This will probably never happen in dev or production
		t.Parallel()

		r := &chiRouteParamManager{}

		ctx := context.Background()
		exampleKey := "blah"
		fn := r.BuildRouteParamIDFetcher(logging.NewNoopLogger(), exampleKey, "thing")
		expected := uint64(0)

		req := testutils.BuildTestRequest(t)
		req = req.WithContext(
			context.WithValue(
				ctx,
				chi.RouteCtxKey,
				&chi.Context{
					URLParams: chi.RouteParams{
						Keys:   []string{exampleKey},
						Values: []string{"expected"},
					},
				},
			),
		)

		actual := fn(req)
		assert.Equal(t, expected, actual)
	})
}

func Test_BuildRouteParamStringIDFetcher(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		r := &chiRouteParamManager{}

		ctx := context.Background()
		exampleKey := "blah"
		fn := r.BuildRouteParamStringIDFetcher(exampleKey)
		expectedInt := uint64(123)
		expected := strconv.FormatUint(expectedInt, 10)
		req := testutils.BuildTestRequest(t).WithContext(
			context.WithValue(
				ctx,
				chi.RouteCtxKey,
				&chi.Context{
					URLParams: chi.RouteParams{
						Keys:   []string{exampleKey},
						Values: []string{strconv.FormatUint(expectedInt, 10)},
					},
				},
			),
		)

		actual := fn(req)
		assert.Equal(t, expected, actual)
	})
}
