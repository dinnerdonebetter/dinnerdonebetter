package httpclient

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/prixfixeco/api_server/pkg/types"
	"github.com/prixfixeco/api_server/pkg/types/fakes"
)

func TestValidMeasurementUnits(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(validMeasurementUnitsTestSuite))
}

type validMeasurementUnitsBaseSuite struct {
	suite.Suite

	ctx                         context.Context
	exampleValidMeasurementUnit *types.ValidMeasurementUnit
}

var _ suite.SetupTestSuite = (*validMeasurementUnitsBaseSuite)(nil)

func (s *validMeasurementUnitsBaseSuite) SetupTest() {
	s.ctx = context.Background()
	s.exampleValidMeasurementUnit = fakes.BuildFakeValidMeasurementUnit()
}

type validMeasurementUnitsTestSuite struct {
	suite.Suite

	validMeasurementUnitsBaseSuite
}

func (s *validMeasurementUnitsTestSuite) TestClient_GetValidMeasurementUnit() {
	const expectedPathFormat = "/api/v1/valid_measurement_units/%s"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, s.exampleValidMeasurementUnit.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleValidMeasurementUnit)
		actual, err := c.GetValidMeasurementUnit(s.ctx, s.exampleValidMeasurementUnit.ID)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, s.exampleValidMeasurementUnit, actual)
	})

	s.Run("with invalid valid ingredient ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetValidMeasurementUnit(s.ctx, "")

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetValidMeasurementUnit(s.ctx, s.exampleValidMeasurementUnit.ID)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, s.exampleValidMeasurementUnit.ID)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetValidMeasurementUnit(s.ctx, s.exampleValidMeasurementUnit.ID)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *validMeasurementUnitsTestSuite) TestClient_GetValidMeasurementUnits() {
	const expectedPath = "/api/v1/valid_measurement_units"

	s.Run("standard", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)

		exampleValidMeasurementUnitList := fakes.BuildFakeValidMeasurementUnitList()

		spec := newRequestSpec(true, http.MethodGet, "includeArchived=false&limit=20&page=1&sortBy=asc", expectedPath)
		c, _ := buildTestClientWithJSONResponse(t, spec, exampleValidMeasurementUnitList)
		actual, err := c.GetValidMeasurementUnits(s.ctx, filter)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, exampleValidMeasurementUnitList, actual)
	})

	s.Run("with error building request", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetValidMeasurementUnits(s.ctx, filter)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)

		spec := newRequestSpec(true, http.MethodGet, "includeArchived=false&limit=20&page=1&sortBy=asc", expectedPath)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetValidMeasurementUnits(s.ctx, filter)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *validMeasurementUnitsTestSuite) TestClient_SearchValidMeasurementUnits() {
	const expectedPath = "/api/v1/valid_measurement_units/search"

	exampleQuery := "whatever"

	s.Run("standard", func() {
		t := s.T()

		exampleValidMeasurementUnitList := fakes.BuildFakeValidMeasurementUnitList()

		spec := newRequestSpec(true, http.MethodGet, "limit=20&q=whatever", expectedPath)
		c, _ := buildTestClientWithJSONResponse(t, spec, exampleValidMeasurementUnitList.ValidMeasurementUnits)
		actual, err := c.SearchValidMeasurementUnits(s.ctx, exampleQuery, 0)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, exampleValidMeasurementUnitList.ValidMeasurementUnits, actual)
	})

	s.Run("with empty query", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		actual, err := c.SearchValidMeasurementUnits(s.ctx, "", 0)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)

		actual, err := c.SearchValidMeasurementUnits(s.ctx, exampleQuery, 0)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with bad response from server", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "limit=20&q=whatever", expectedPath)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.SearchValidMeasurementUnits(s.ctx, exampleQuery, 0)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *validMeasurementUnitsTestSuite) TestClient_CreateValidMeasurementUnit() {
	const expectedPath = "/api/v1/valid_measurement_units"

	s.Run("standard", func() {
		t := s.T()

		exampleInput := fakes.BuildFakeValidMeasurementUnitCreationRequestInput()

		spec := newRequestSpec(false, http.MethodPost, "", expectedPath)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleValidMeasurementUnit)

		actual, err := c.CreateValidMeasurementUnit(s.ctx, exampleInput)
		require.NotEmpty(t, actual)
		assert.NoError(t, err)

		assert.Equal(t, s.exampleValidMeasurementUnit, actual)
	})

	s.Run("with nil input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		actual, err := c.CreateValidMeasurementUnit(s.ctx, nil)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with invalid input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		exampleInput := &types.ValidMeasurementUnitCreationRequestInput{}

		actual, err := c.CreateValidMeasurementUnit(s.ctx, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		exampleInput := fakes.BuildFakeValidMeasurementUnitCreationRequestInputFromValidMeasurementUnit(s.exampleValidMeasurementUnit)

		c := buildTestClientWithInvalidURL(t)

		actual, err := c.CreateValidMeasurementUnit(s.ctx, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		exampleInput := fakes.BuildFakeValidMeasurementUnitCreationRequestInputFromValidMeasurementUnit(s.exampleValidMeasurementUnit)
		c, _ := buildTestClientThatWaitsTooLong(t)

		actual, err := c.CreateValidMeasurementUnit(s.ctx, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *validMeasurementUnitsTestSuite) TestClient_UpdateValidMeasurementUnit() {
	const expectedPathFormat = "/api/v1/valid_measurement_units/%s"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, s.exampleValidMeasurementUnit.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleValidMeasurementUnit)

		err := c.UpdateValidMeasurementUnit(s.ctx, s.exampleValidMeasurementUnit)
		assert.NoError(t, err)
	})

	s.Run("with nil input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		err := c.UpdateValidMeasurementUnit(s.ctx, nil)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)

		err := c.UpdateValidMeasurementUnit(s.ctx, s.exampleValidMeasurementUnit)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		c, _ := buildTestClientThatWaitsTooLong(t)

		err := c.UpdateValidMeasurementUnit(s.ctx, s.exampleValidMeasurementUnit)
		assert.Error(t, err)
	})
}

func (s *validMeasurementUnitsTestSuite) TestClient_ArchiveValidMeasurementUnit() {
	const expectedPathFormat = "/api/v1/valid_measurement_units/%s"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodDelete, "", expectedPathFormat, s.exampleValidMeasurementUnit.ID)
		c, _ := buildTestClientWithStatusCodeResponse(t, spec, http.StatusOK)

		err := c.ArchiveValidMeasurementUnit(s.ctx, s.exampleValidMeasurementUnit.ID)
		assert.NoError(t, err)
	})

	s.Run("with invalid valid ingredient ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		err := c.ArchiveValidMeasurementUnit(s.ctx, "")
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)

		err := c.ArchiveValidMeasurementUnit(s.ctx, s.exampleValidMeasurementUnit.ID)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		c, _ := buildTestClientThatWaitsTooLong(t)

		err := c.ArchiveValidMeasurementUnit(s.ctx, s.exampleValidMeasurementUnit.ID)
		assert.Error(t, err)
	})
}
