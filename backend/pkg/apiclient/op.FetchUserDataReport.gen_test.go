// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

import (
	"context"
	"net/http"
	"testing"

	"github.com/dinnerdonebetter/backend/lib/internal/fake"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_FetchUserDataReport(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/data_privacy/reports/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		userDataAggregationReportID := fake.BuildFakeID()

		data := fake.BuildFakeForTest[*UserDataCollection](t)
		expected := &APIResponse[*UserDataCollection]{
			Data: data,
		}

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, userDataAggregationReportID)
		c, _ := buildTestClientWithJSONResponse(t, spec, expected)
		actual, err := c.FetchUserDataReport(ctx, userDataAggregationReportID)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, data, actual)
	})

	T.Run("with invalid userDataAggregationReport ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildSimpleTestClient(t)
		actual, err := c.FetchUserDataReport(ctx, "")

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with error building request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		userDataAggregationReportID := fake.BuildFakeID()

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.FetchUserDataReport(ctx, userDataAggregationReportID)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with error executing request", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		userDataAggregationReportID := fake.BuildFakeID()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, userDataAggregationReportID)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.FetchUserDataReport(ctx, userDataAggregationReportID)

		require.Nil(t, actual)
		assert.Error(t, err)
	})
}
