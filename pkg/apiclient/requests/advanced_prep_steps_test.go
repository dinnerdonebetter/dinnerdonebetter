package requests

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/prixfixeco/api_server/pkg/types"
	"github.com/prixfixeco/api_server/pkg/types/fakes"
)

func TestBuilder_BuildGetAdvancedPrepStepRequest(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/advanced_prep_steps/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleAdvancedPrepStep := fakes.BuildFakeAdvancedPrepStep()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, exampleAdvancedPrepStep.ID)

		actual, err := helper.builder.BuildGetAdvancedPrepStepRequest(helper.ctx, exampleAdvancedPrepStep.ID)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with invalid advanced prep step ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		actual, err := helper.builder.BuildGetAdvancedPrepStepRequest(helper.ctx, "")
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		exampleAdvancedPrepStep := fakes.BuildFakeAdvancedPrepStep()

		actual, err := helper.builder.BuildGetAdvancedPrepStepRequest(helper.ctx, exampleAdvancedPrepStep.ID)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildGetAdvancedPrepStepsRequest(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/advanced_prep_steps"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		filter := (*types.QueryFilter)(nil)
		spec := newRequestSpec(true, http.MethodGet, "limit=20&page=1&sortBy=asc", expectedPathFormat)

		actual, err := helper.builder.BuildGetAdvancedPrepStepsRequest(helper.ctx, filter)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		filter := (*types.QueryFilter)(nil)

		actual, err := helper.builder.BuildGetAdvancedPrepStepsRequest(helper.ctx, filter)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildChangeAdvancedPrepStepStatusRequest(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/advanced_prep_steps/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleInput := fakes.BuildFakeAdvancedPrepStepStatusChangeRequestInput()
		spec := newRequestSpec(false, http.MethodPost, "", expectedPathFormat, exampleInput.ID)

		actual, err := helper.builder.BuildChangeAdvancedPrepStepStatusRequest(helper.ctx, exampleInput)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with invalid advanced prep step ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		actual, err := helper.builder.BuildChangeAdvancedPrepStepStatusRequest(helper.ctx, nil)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		exampleInput := fakes.BuildFakeAdvancedPrepStepStatusChangeRequestInput()

		actual, err := helper.builder.BuildChangeAdvancedPrepStepStatusRequest(helper.ctx, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}
