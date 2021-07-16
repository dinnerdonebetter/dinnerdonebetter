package requests

import (
	"net/http"
	"testing"

	"gitlab.com/prixfixe/prixfixe/pkg/types"
	"gitlab.com/prixfixe/prixfixe/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuilder_BuildValidPreparationExistsRequest(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/valid_preparations/%d"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleValidPreparation := fakes.BuildFakeValidPreparation()

		actual, err := helper.builder.BuildValidPreparationExistsRequest(helper.ctx, exampleValidPreparation.ID)
		spec := newRequestSpec(true, http.MethodHead, "", expectedPathFormat, exampleValidPreparation.ID)

		assert.NoError(t, err)
		assertRequestQuality(t, actual, spec)
	})

	T.Run("with invalid valid preparation ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		actual, err := helper.builder.BuildValidPreparationExistsRequest(helper.ctx, 0)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		exampleValidPreparation := fakes.BuildFakeValidPreparation()

		actual, err := helper.builder.BuildValidPreparationExistsRequest(helper.ctx, exampleValidPreparation.ID)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildGetValidPreparationRequest(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/valid_preparations/%d"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleValidPreparation := fakes.BuildFakeValidPreparation()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, exampleValidPreparation.ID)

		actual, err := helper.builder.BuildGetValidPreparationRequest(helper.ctx, exampleValidPreparation.ID)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with invalid valid preparation ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		actual, err := helper.builder.BuildGetValidPreparationRequest(helper.ctx, 0)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		exampleValidPreparation := fakes.BuildFakeValidPreparation()

		actual, err := helper.builder.BuildGetValidPreparationRequest(helper.ctx, exampleValidPreparation.ID)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildGetValidPreparationsRequest(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/valid_preparations"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		filter := (*types.QueryFilter)(nil)
		spec := newRequestSpec(true, http.MethodGet, "includeArchived=false&limit=20&page=1&sortBy=asc", expectedPathFormat)

		actual, err := helper.builder.BuildGetValidPreparationsRequest(helper.ctx, filter)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		filter := (*types.QueryFilter)(nil)

		actual, err := helper.builder.BuildGetValidPreparationsRequest(helper.ctx, filter)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildSearchValidPreparationsRequest(T *testing.T) {
	T.Parallel()

	const expectedPath = "/api/v1/valid_preparations/search"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		limit := types.DefaultQueryFilter().Limit
		exampleQuery := "whatever"
		spec := newRequestSpec(true, http.MethodGet, "limit=20&q=whatever", expectedPath)

		actual, err := helper.builder.BuildSearchValidPreparationsRequest(helper.ctx, exampleQuery, limit)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		limit := types.DefaultQueryFilter().Limit
		exampleQuery := "whatever"

		actual, err := helper.builder.BuildSearchValidPreparationsRequest(helper.ctx, exampleQuery, limit)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildCreateValidPreparationRequest(T *testing.T) {
	T.Parallel()

	const expectedPath = "/api/v1/valid_preparations"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleInput := fakes.BuildFakeValidPreparationCreationInput()

		spec := newRequestSpec(false, http.MethodPost, "", expectedPath)

		actual, err := helper.builder.BuildCreateValidPreparationRequest(helper.ctx, exampleInput)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		actual, err := helper.builder.BuildCreateValidPreparationRequest(helper.ctx, nil)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		actual, err := helper.builder.BuildCreateValidPreparationRequest(helper.ctx, &types.ValidPreparationCreationInput{})
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		exampleInput := fakes.BuildFakeValidPreparationCreationInput()

		actual, err := helper.builder.BuildCreateValidPreparationRequest(helper.ctx, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildUpdateValidPreparationRequest(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/valid_preparations/%d"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleValidPreparation := fakes.BuildFakeValidPreparation()

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, exampleValidPreparation.ID)

		actual, err := helper.builder.BuildUpdateValidPreparationRequest(helper.ctx, exampleValidPreparation)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		actual, err := helper.builder.BuildUpdateValidPreparationRequest(helper.ctx, nil)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		exampleValidPreparation := fakes.BuildFakeValidPreparation()

		actual, err := helper.builder.BuildUpdateValidPreparationRequest(helper.ctx, exampleValidPreparation)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildArchiveValidPreparationRequest(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/valid_preparations/%d"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleValidPreparation := fakes.BuildFakeValidPreparation()

		spec := newRequestSpec(true, http.MethodDelete, "", expectedPathFormat, exampleValidPreparation.ID)

		actual, err := helper.builder.BuildArchiveValidPreparationRequest(helper.ctx, exampleValidPreparation.ID)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with invalid valid preparation ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		actual, err := helper.builder.BuildArchiveValidPreparationRequest(helper.ctx, 0)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		exampleValidPreparation := fakes.BuildFakeValidPreparation()

		actual, err := helper.builder.BuildArchiveValidPreparationRequest(helper.ctx, exampleValidPreparation.ID)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildGetAuditLogForValidPreparationRequest(T *testing.T) {
	T.Parallel()

	const expectedPath = "/api/v1/valid_preparations/%d/audit"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleValidPreparation := fakes.BuildFakeValidPreparation()

		actual, err := helper.builder.BuildGetAuditLogForValidPreparationRequest(helper.ctx, exampleValidPreparation.ID)
		require.NotNil(t, actual)
		assert.NoError(t, err)

		spec := newRequestSpec(true, http.MethodGet, "", expectedPath, exampleValidPreparation.ID)
		assertRequestQuality(t, actual, spec)
	})

	T.Run("with invalid valid preparation ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		actual, err := helper.builder.BuildGetAuditLogForValidPreparationRequest(helper.ctx, 0)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		exampleValidPreparation := fakes.BuildFakeValidPreparation()

		actual, err := helper.builder.BuildGetAuditLogForValidPreparationRequest(helper.ctx, exampleValidPreparation.ID)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}
