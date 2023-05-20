package requests

import (
	"net/http"
	"testing"

	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
)

func TestBuilder_BuildGetValidMeasurementConversionRequest(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/valid_measurement_conversions/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleValidMeasurementConversion := fakes.BuildFakeValidMeasurementConversion()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, exampleValidMeasurementConversion.ID)

		actual, err := helper.builder.BuildGetValidMeasurementConversionRequest(helper.ctx, exampleValidMeasurementConversion.ID)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with invalid valid preparation ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		actual, err := helper.builder.BuildGetValidMeasurementConversionRequest(helper.ctx, "")
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		exampleValidMeasurementConversion := fakes.BuildFakeValidMeasurementConversion()

		actual, err := helper.builder.BuildGetValidMeasurementConversionRequest(helper.ctx, exampleValidMeasurementConversion.ID)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildGetValidMeasurementConversionsFromUnitRequest(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/valid_measurement_conversions/from_unit/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleValidMeasurementUnit := fakes.BuildFakeValidMeasurementUnit()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, exampleValidMeasurementUnit.ID)

		actual, err := helper.builder.BuildGetValidMeasurementConversionsFromUnitRequest(helper.ctx, exampleValidMeasurementUnit.ID)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with invalid valid preparation ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		actual, err := helper.builder.BuildGetValidMeasurementConversionsFromUnitRequest(helper.ctx, "")
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		exampleValidMeasurementUnit := fakes.BuildFakeValidMeasurementUnit()

		actual, err := helper.builder.BuildGetValidMeasurementConversionsFromUnitRequest(helper.ctx, exampleValidMeasurementUnit.ID)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildGetValidMeasurementConversionsToUnitRequest(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/valid_measurement_conversions/to_unit/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleValidMeasurementUnit := fakes.BuildFakeValidMeasurementUnit()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, exampleValidMeasurementUnit.ID)

		actual, err := helper.builder.BuildGetValidMeasurementConversionsToUnitRequest(helper.ctx, exampleValidMeasurementUnit.ID)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with invalid valid preparation ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		actual, err := helper.builder.BuildGetValidMeasurementConversionsToUnitRequest(helper.ctx, "")
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		exampleValidMeasurementUnit := fakes.BuildFakeValidMeasurementUnit()

		actual, err := helper.builder.BuildGetValidMeasurementConversionsToUnitRequest(helper.ctx, exampleValidMeasurementUnit.ID)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildCreateValidMeasurementConversionRequest(T *testing.T) {
	T.Parallel()

	const expectedPath = "/api/v1/valid_measurement_conversions"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleInput := fakes.BuildFakeValidMeasurementConversionCreationRequestInput()

		spec := newRequestSpec(false, http.MethodPost, "", expectedPath)

		actual, err := helper.builder.BuildCreateValidMeasurementConversionRequest(helper.ctx, exampleInput)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		actual, err := helper.builder.BuildCreateValidMeasurementConversionRequest(helper.ctx, nil)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		actual, err := helper.builder.BuildCreateValidMeasurementConversionRequest(helper.ctx, &types.ValidMeasurementUnitConversionCreationRequestInput{})
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		exampleInput := fakes.BuildFakeValidMeasurementConversionCreationRequestInput()

		actual, err := helper.builder.BuildCreateValidMeasurementConversionRequest(helper.ctx, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildUpdateValidMeasurementConversionRequest(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/valid_measurement_conversions/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleValidMeasurementConversion := fakes.BuildFakeValidMeasurementConversion()

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, exampleValidMeasurementConversion.ID)

		actual, err := helper.builder.BuildUpdateValidMeasurementConversionRequest(helper.ctx, exampleValidMeasurementConversion)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		actual, err := helper.builder.BuildUpdateValidMeasurementConversionRequest(helper.ctx, nil)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		exampleValidMeasurementConversion := fakes.BuildFakeValidMeasurementConversion()

		actual, err := helper.builder.BuildUpdateValidMeasurementConversionRequest(helper.ctx, exampleValidMeasurementConversion)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildArchiveValidMeasurementConversionRequest(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/valid_measurement_conversions/%s"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleValidMeasurementConversion := fakes.BuildFakeValidMeasurementConversion()

		spec := newRequestSpec(true, http.MethodDelete, "", expectedPathFormat, exampleValidMeasurementConversion.ID)

		actual, err := helper.builder.BuildArchiveValidMeasurementConversionRequest(helper.ctx, exampleValidMeasurementConversion.ID)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with invalid valid preparation ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		actual, err := helper.builder.BuildArchiveValidMeasurementConversionRequest(helper.ctx, "")
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		exampleValidMeasurementConversion := fakes.BuildFakeValidMeasurementConversion()

		actual, err := helper.builder.BuildArchiveValidMeasurementConversionRequest(helper.ctx, exampleValidMeasurementConversion.ID)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}
