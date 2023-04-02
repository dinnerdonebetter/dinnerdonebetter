package apiclient

import (
	"context"
	"net/http"
	"testing"

	"github.com/prixfixeco/backend/pkg/types"
	"github.com/prixfixeco/backend/pkg/types/converters"
	"github.com/prixfixeco/backend/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

func TestValidMeasurementConversions(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(validMeasurementConversionsTestSuite))
}

type validMeasurementConversionsBaseSuite struct {
	suite.Suite

	ctx                                   context.Context
	exampleValidMeasurementUnit           *types.ValidMeasurementUnit
	exampleValidMeasurementConversion     *types.ValidMeasurementUnitConversion
	exampleValidMeasurementConversionList []*types.ValidMeasurementUnitConversion
}

var _ suite.SetupTestSuite = (*validMeasurementConversionsBaseSuite)(nil)

func (s *validMeasurementConversionsBaseSuite) SetupTest() {
	s.ctx = context.Background()
	s.exampleValidMeasurementUnit = fakes.BuildFakeValidMeasurementUnit()
	s.exampleValidMeasurementConversion = fakes.BuildFakeValidMeasurementConversion()
	s.exampleValidMeasurementConversionList = []*types.ValidMeasurementUnitConversion{
		s.exampleValidMeasurementConversion,
	}
}

type validMeasurementConversionsTestSuite struct {
	suite.Suite

	validMeasurementConversionsBaseSuite
}

func (s *validMeasurementConversionsTestSuite) TestClient_GetValidMeasurementConversion() {
	const expectedPathFormat = "/api/v1/valid_measurement_conversions/%s"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, s.exampleValidMeasurementConversion.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleValidMeasurementConversion)
		actual, err := c.GetValidMeasurementConversion(s.ctx, s.exampleValidMeasurementConversion.ID)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, s.exampleValidMeasurementConversion, actual)
	})

	s.Run("with invalid valid preparation ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetValidMeasurementConversion(s.ctx, "")

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetValidMeasurementConversion(s.ctx, s.exampleValidMeasurementConversion.ID)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, s.exampleValidMeasurementConversion.ID)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetValidMeasurementConversion(s.ctx, s.exampleValidMeasurementConversion.ID)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *validMeasurementConversionsTestSuite) TestClient_GetValidMeasurementConversionsFromUnit() {
	const expectedPathFormat = "/api/v1/valid_measurement_conversions/from_unit/%s"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, s.exampleValidMeasurementUnit.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleValidMeasurementConversionList)
		actual, err := c.GetValidMeasurementConversionsFromUnit(s.ctx, s.exampleValidMeasurementUnit.ID)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, s.exampleValidMeasurementConversionList, actual)
	})

	s.Run("with invalid valid preparation ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetValidMeasurementConversionsFromUnit(s.ctx, "")

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetValidMeasurementConversionsFromUnit(s.ctx, s.exampleValidMeasurementUnit.ID)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, s.exampleValidMeasurementUnit.ID)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetValidMeasurementConversionsFromUnit(s.ctx, s.exampleValidMeasurementUnit.ID)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *validMeasurementConversionsTestSuite) TestClient_CreateValidMeasurementConversion() {
	const expectedPath = "/api/v1/valid_measurement_conversions"

	s.Run("standard", func() {
		t := s.T()

		exampleInput := fakes.BuildFakeValidMeasurementConversionCreationRequestInput()

		spec := newRequestSpec(false, http.MethodPost, "", expectedPath)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleValidMeasurementConversion)

		actual, err := c.CreateValidMeasurementConversion(s.ctx, exampleInput)
		require.NotEmpty(t, actual)
		assert.NoError(t, err)

		assert.Equal(t, s.exampleValidMeasurementConversion, actual)
	})

	s.Run("with nil input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		actual, err := c.CreateValidMeasurementConversion(s.ctx, nil)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with invalid input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		exampleInput := &types.ValidMeasurementUnitConversionCreationRequestInput{}

		actual, err := c.CreateValidMeasurementConversion(s.ctx, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		exampleInput := converters.ConvertValidMeasurementConversionToValidMeasurementConversionCreationRequestInput(s.exampleValidMeasurementConversion)

		c := buildTestClientWithInvalidURL(t)

		actual, err := c.CreateValidMeasurementConversion(s.ctx, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		exampleInput := converters.ConvertValidMeasurementConversionToValidMeasurementConversionCreationRequestInput(s.exampleValidMeasurementConversion)
		c, _ := buildTestClientThatWaitsTooLong(t)

		actual, err := c.CreateValidMeasurementConversion(s.ctx, exampleInput)
		assert.Empty(t, actual)
		assert.Error(t, err)
	})
}

func (s *validMeasurementConversionsTestSuite) TestClient_UpdateValidMeasurementConversion() {
	const expectedPathFormat = "/api/v1/valid_measurement_conversions/%s"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, s.exampleValidMeasurementConversion.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleValidMeasurementConversion)

		err := c.UpdateValidMeasurementConversion(s.ctx, s.exampleValidMeasurementConversion)
		assert.NoError(t, err)
	})

	s.Run("with nil input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		err := c.UpdateValidMeasurementConversion(s.ctx, nil)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)

		err := c.UpdateValidMeasurementConversion(s.ctx, s.exampleValidMeasurementConversion)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		c, _ := buildTestClientThatWaitsTooLong(t)

		err := c.UpdateValidMeasurementConversion(s.ctx, s.exampleValidMeasurementConversion)
		assert.Error(t, err)
	})
}

func (s *validMeasurementConversionsTestSuite) TestClient_ArchiveValidMeasurementConversion() {
	const expectedPathFormat = "/api/v1/valid_measurement_conversions/%s"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodDelete, "", expectedPathFormat, s.exampleValidMeasurementConversion.ID)
		c, _ := buildTestClientWithStatusCodeResponse(t, spec, http.StatusOK)

		err := c.ArchiveValidMeasurementConversion(s.ctx, s.exampleValidMeasurementConversion.ID)
		assert.NoError(t, err)
	})

	s.Run("with invalid valid preparation ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		err := c.ArchiveValidMeasurementConversion(s.ctx, "")
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)

		err := c.ArchiveValidMeasurementConversion(s.ctx, s.exampleValidMeasurementConversion.ID)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		c, _ := buildTestClientThatWaitsTooLong(t)

		err := c.ArchiveValidMeasurementConversion(s.ctx, s.exampleValidMeasurementConversion.ID)
		assert.Error(t, err)
	})
}
