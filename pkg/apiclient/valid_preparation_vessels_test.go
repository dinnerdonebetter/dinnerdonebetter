package apiclient

import (
	"context"
	"net/http"
	"testing"

	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

func TestValidPreparationVessels(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(validPreparationVesselsTestSuite))
}

type validPreparationVesselsBaseSuite struct {
	suite.Suite

	ctx                           context.Context
	exampleValidPreparationVessel *types.ValidPreparationVessel
}

var _ suite.SetupTestSuite = (*validPreparationVesselsBaseSuite)(nil)

func (s *validPreparationVesselsBaseSuite) SetupTest() {
	s.ctx = context.Background()
	s.exampleValidPreparationVessel = fakes.BuildFakeValidPreparationVessel()
}

type validPreparationVesselsTestSuite struct {
	suite.Suite

	validPreparationVesselsBaseSuite
}

func (s *validPreparationVesselsTestSuite) TestClient_GetValidPreparationVessel() {
	const expectedPathFormat = "/api/v1/valid_preparation_vessels/%s"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, s.exampleValidPreparationVessel.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleValidPreparationVessel)
		actual, err := c.GetValidPreparationVessel(s.ctx, s.exampleValidPreparationVessel.ID)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, s.exampleValidPreparationVessel, actual)
	})

	s.Run("with invalid valid ingredient preparation ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetValidPreparationVessel(s.ctx, "")

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetValidPreparationVessel(s.ctx, s.exampleValidPreparationVessel.ID)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, s.exampleValidPreparationVessel.ID)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetValidPreparationVessel(s.ctx, s.exampleValidPreparationVessel.ID)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *validPreparationVesselsTestSuite) TestClient_GetValidPreparationVessels() {
	const expectedPath = "/api/v1/valid_preparation_vessels"

	s.Run("standard", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)

		exampleValidPreparationVesselList := fakes.BuildFakeValidPreparationVesselList()

		spec := newRequestSpec(true, http.MethodGet, "limit=50&page=1&sortBy=asc", expectedPath)
		c, _ := buildTestClientWithJSONResponse(t, spec, exampleValidPreparationVesselList)
		actual, err := c.GetValidPreparationVessels(s.ctx, filter)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, exampleValidPreparationVesselList, actual)
	})

	s.Run("with error building request", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetValidPreparationVessels(s.ctx, filter)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)

		spec := newRequestSpec(true, http.MethodGet, "limit=50&page=1&sortBy=asc", expectedPath)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetValidPreparationVessels(s.ctx, filter)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *validIngredientMeasurementUnitsTestSuite) TestClient_GetValidPreparationVesselsForPreparation() {
	const expectedPath = "/api/v1/valid_preparation_vessels/by_preparation/%s"

	exampleValidPreparation := fakes.BuildFakeValidPreparation()

	s.Run("standard", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)

		exampleValidIngredientMeasurementUnitList := fakes.BuildFakeValidPreparationVesselList()

		spec := newRequestSpec(true, http.MethodGet, "limit=50&page=1&sortBy=asc", expectedPath, exampleValidPreparation.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, exampleValidIngredientMeasurementUnitList)
		actual, err := c.GetValidPreparationVesselsForPreparation(s.ctx, exampleValidPreparation.ID, filter)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, exampleValidIngredientMeasurementUnitList, actual)
	})

	s.Run("with invalid ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetValidPreparationVesselsForPreparation(s.ctx, "", nil)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetValidPreparationVesselsForPreparation(s.ctx, exampleValidPreparation.ID, filter)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)

		spec := newRequestSpec(true, http.MethodGet, "limit=50&page=1&sortBy=asc", expectedPath, exampleValidPreparation.ID)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetValidPreparationVesselsForPreparation(s.ctx, exampleValidPreparation.ID, filter)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *validIngredientMeasurementUnitsTestSuite) TestClient_GetValidPreparationVesselsForVessel() {
	const expectedPath = "/api/v1/valid_preparation_vessels/by_vessel/%s"

	exampleValidInstrument := fakes.BuildFakeValidInstrument()

	s.Run("standard", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)

		exampleValidIngredientMeasurementUnitList := fakes.BuildFakeValidPreparationVesselList()

		spec := newRequestSpec(true, http.MethodGet, "limit=50&page=1&sortBy=asc", expectedPath, exampleValidInstrument.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, exampleValidIngredientMeasurementUnitList)
		actual, err := c.GetValidPreparationVesselsForVessel(s.ctx, exampleValidInstrument.ID, filter)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, exampleValidIngredientMeasurementUnitList, actual)
	})

	s.Run("with invalid ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetValidPreparationVesselsForVessel(s.ctx, "", nil)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetValidPreparationVesselsForVessel(s.ctx, exampleValidInstrument.ID, filter)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)

		spec := newRequestSpec(true, http.MethodGet, "limit=50&page=1&sortBy=asc", expectedPath, exampleValidInstrument.ID)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetValidPreparationVesselsForVessel(s.ctx, exampleValidInstrument.ID, filter)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *validPreparationVesselsTestSuite) TestClient_CreateValidPreparationVessel() {
	const expectedPath = "/api/v1/valid_preparation_vessels"

	s.Run("standard", func() {
		t := s.T()

		exampleInput := fakes.BuildFakeValidPreparationVesselCreationRequestInput()

		spec := newRequestSpec(false, http.MethodPost, "", expectedPath)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleValidPreparationVessel)

		actual, err := c.CreateValidPreparationVessel(s.ctx, exampleInput)
		assert.NoError(t, err)
		assert.Equal(t, s.exampleValidPreparationVessel, actual)
	})

	s.Run("with nil input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		actual, err := c.CreateValidPreparationVessel(s.ctx, nil)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with invalid input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		exampleInput := &types.ValidPreparationVesselCreationRequestInput{}

		actual, err := c.CreateValidPreparationVessel(s.ctx, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		exampleInput := converters.ConvertValidPreparationVesselToValidPreparationVesselCreationRequestInput(s.exampleValidPreparationVessel)

		c := buildTestClientWithInvalidURL(t)

		actual, err := c.CreateValidPreparationVessel(s.ctx, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		exampleInput := converters.ConvertValidPreparationVesselToValidPreparationVesselCreationRequestInput(s.exampleValidPreparationVessel)
		c, _ := buildTestClientThatWaitsTooLong(t)

		actual, err := c.CreateValidPreparationVessel(s.ctx, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *validPreparationVesselsTestSuite) TestClient_UpdateValidPreparationVessel() {
	const expectedPathFormat = "/api/v1/valid_preparation_vessels/%s"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, s.exampleValidPreparationVessel.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleValidPreparationVessel)

		err := c.UpdateValidPreparationVessel(s.ctx, s.exampleValidPreparationVessel)
		assert.NoError(t, err)
	})

	s.Run("with nil input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		err := c.UpdateValidPreparationVessel(s.ctx, nil)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)

		err := c.UpdateValidPreparationVessel(s.ctx, s.exampleValidPreparationVessel)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		c, _ := buildTestClientThatWaitsTooLong(t)

		err := c.UpdateValidPreparationVessel(s.ctx, s.exampleValidPreparationVessel)
		assert.Error(t, err)
	})
}

func (s *validPreparationVesselsTestSuite) TestClient_ArchiveValidPreparationVessel() {
	const expectedPathFormat = "/api/v1/valid_preparation_vessels/%s"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodDelete, "", expectedPathFormat, s.exampleValidPreparationVessel.ID)
		c, _ := buildTestClientWithStatusCodeResponse(t, spec, http.StatusOK)

		err := c.ArchiveValidPreparationVessel(s.ctx, s.exampleValidPreparationVessel.ID)
		assert.NoError(t, err)
	})

	s.Run("with invalid valid ingredient preparation ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		err := c.ArchiveValidPreparationVessel(s.ctx, "")
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)

		err := c.ArchiveValidPreparationVessel(s.ctx, s.exampleValidPreparationVessel.ID)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		c, _ := buildTestClientThatWaitsTooLong(t)

		err := c.ArchiveValidPreparationVessel(s.ctx, s.exampleValidPreparationVessel.ID)
		assert.Error(t, err)
	})
}
