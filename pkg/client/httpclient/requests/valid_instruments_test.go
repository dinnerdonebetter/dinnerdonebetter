package requests

import (
	"net/http"
	"testing"

	"gitlab.com/prixfixe/prixfixe/pkg/types"
	"gitlab.com/prixfixe/prixfixe/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuilder_BuildValidInstrumentExistsRequest(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/valid_instruments/%d"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleValidInstrument := fakes.BuildFakeValidInstrument()

		actual, err := helper.builder.BuildValidInstrumentExistsRequest(helper.ctx, exampleValidInstrument.ID)
		spec := newRequestSpec(true, http.MethodHead, "", expectedPathFormat, exampleValidInstrument.ID)

		assert.NoError(t, err)
		assertRequestQuality(t, actual, spec)
	})

	T.Run("with invalid valid instrument ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		actual, err := helper.builder.BuildValidInstrumentExistsRequest(helper.ctx, 0)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		exampleValidInstrument := fakes.BuildFakeValidInstrument()

		actual, err := helper.builder.BuildValidInstrumentExistsRequest(helper.ctx, exampleValidInstrument.ID)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildGetValidInstrumentRequest(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/valid_instruments/%d"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleValidInstrument := fakes.BuildFakeValidInstrument()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, exampleValidInstrument.ID)

		actual, err := helper.builder.BuildGetValidInstrumentRequest(helper.ctx, exampleValidInstrument.ID)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with invalid valid instrument ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		actual, err := helper.builder.BuildGetValidInstrumentRequest(helper.ctx, 0)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		exampleValidInstrument := fakes.BuildFakeValidInstrument()

		actual, err := helper.builder.BuildGetValidInstrumentRequest(helper.ctx, exampleValidInstrument.ID)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildGetValidInstrumentsRequest(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/valid_instruments"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		filter := (*types.QueryFilter)(nil)
		spec := newRequestSpec(true, http.MethodGet, "includeArchived=false&limit=20&page=1&sortBy=asc", expectedPathFormat)

		actual, err := helper.builder.BuildGetValidInstrumentsRequest(helper.ctx, filter)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		filter := (*types.QueryFilter)(nil)

		actual, err := helper.builder.BuildGetValidInstrumentsRequest(helper.ctx, filter)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildSearchValidInstrumentsRequest(T *testing.T) {
	T.Parallel()

	const expectedPath = "/api/v1/valid_instruments/search"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		limit := types.DefaultQueryFilter().Limit
		exampleQuery := "whatever"
		spec := newRequestSpec(true, http.MethodGet, "limit=20&q=whatever", expectedPath)

		actual, err := helper.builder.BuildSearchValidInstrumentsRequest(helper.ctx, exampleQuery, limit)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		limit := types.DefaultQueryFilter().Limit
		exampleQuery := "whatever"

		actual, err := helper.builder.BuildSearchValidInstrumentsRequest(helper.ctx, exampleQuery, limit)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildCreateValidInstrumentRequest(T *testing.T) {
	T.Parallel()

	const expectedPath = "/api/v1/valid_instruments"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleInput := fakes.BuildFakeValidInstrumentCreationInput()

		spec := newRequestSpec(false, http.MethodPost, "", expectedPath)

		actual, err := helper.builder.BuildCreateValidInstrumentRequest(helper.ctx, exampleInput)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		actual, err := helper.builder.BuildCreateValidInstrumentRequest(helper.ctx, nil)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		actual, err := helper.builder.BuildCreateValidInstrumentRequest(helper.ctx, &types.ValidInstrumentCreationInput{})
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		exampleInput := fakes.BuildFakeValidInstrumentCreationInput()

		actual, err := helper.builder.BuildCreateValidInstrumentRequest(helper.ctx, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildUpdateValidInstrumentRequest(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/valid_instruments/%d"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleValidInstrument := fakes.BuildFakeValidInstrument()

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, exampleValidInstrument.ID)

		actual, err := helper.builder.BuildUpdateValidInstrumentRequest(helper.ctx, exampleValidInstrument)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		actual, err := helper.builder.BuildUpdateValidInstrumentRequest(helper.ctx, nil)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		exampleValidInstrument := fakes.BuildFakeValidInstrument()

		actual, err := helper.builder.BuildUpdateValidInstrumentRequest(helper.ctx, exampleValidInstrument)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildArchiveValidInstrumentRequest(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/valid_instruments/%d"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleValidInstrument := fakes.BuildFakeValidInstrument()

		spec := newRequestSpec(true, http.MethodDelete, "", expectedPathFormat, exampleValidInstrument.ID)

		actual, err := helper.builder.BuildArchiveValidInstrumentRequest(helper.ctx, exampleValidInstrument.ID)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with invalid valid instrument ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		actual, err := helper.builder.BuildArchiveValidInstrumentRequest(helper.ctx, 0)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		exampleValidInstrument := fakes.BuildFakeValidInstrument()

		actual, err := helper.builder.BuildArchiveValidInstrumentRequest(helper.ctx, exampleValidInstrument.ID)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildGetAuditLogForValidInstrumentRequest(T *testing.T) {
	T.Parallel()

	const expectedPath = "/api/v1/valid_instruments/%d/audit"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleValidInstrument := fakes.BuildFakeValidInstrument()

		actual, err := helper.builder.BuildGetAuditLogForValidInstrumentRequest(helper.ctx, exampleValidInstrument.ID)
		require.NotNil(t, actual)
		assert.NoError(t, err)

		spec := newRequestSpec(true, http.MethodGet, "", expectedPath, exampleValidInstrument.ID)
		assertRequestQuality(t, actual, spec)
	})

	T.Run("with invalid valid instrument ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		actual, err := helper.builder.BuildGetAuditLogForValidInstrumentRequest(helper.ctx, 0)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		exampleValidInstrument := fakes.BuildFakeValidInstrument()

		actual, err := helper.builder.BuildGetAuditLogForValidInstrumentRequest(helper.ctx, exampleValidInstrument.ID)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}
