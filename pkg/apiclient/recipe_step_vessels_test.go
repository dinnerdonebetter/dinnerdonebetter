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

func TestRecipeStepVessels(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(recipeStepVesselsTestSuite))
}

type recipeStepVesselsBaseSuite struct {
	suite.Suite
	ctx                     context.Context
	exampleRecipeStepVessel *types.RecipeStepVessel
	exampleRecipeID         string
	exampleRecipeStepID     string
}

var _ suite.SetupTestSuite = (*recipeStepVesselsBaseSuite)(nil)

func (s *recipeStepVesselsBaseSuite) SetupTest() {
	s.ctx = context.Background()
	s.exampleRecipeID = fakes.BuildFakeID()
	s.exampleRecipeStepID = fakes.BuildFakeID()
	s.exampleRecipeStepVessel = fakes.BuildFakeRecipeStepVessel()
	s.exampleRecipeStepVessel.BelongsToRecipeStep = s.exampleRecipeStepID
}

type recipeStepVesselsTestSuite struct {
	suite.Suite

	recipeStepVesselsBaseSuite
}

func (s *recipeStepVesselsTestSuite) TestClient_GetRecipeStepVessel() {
	const expectedPathFormat = "/api/v1/recipes/%s/steps/%s/vessels/%s"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, s.exampleRecipeID, s.exampleRecipeStepID, s.exampleRecipeStepVessel.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleRecipeStepVessel)
		actual, err := c.GetRecipeStepVessel(s.ctx, s.exampleRecipeID, s.exampleRecipeStepID, s.exampleRecipeStepVessel.ID)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, s.exampleRecipeStepVessel, actual)
	})

	s.Run("with invalid recipe ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetRecipeStepVessel(s.ctx, "", s.exampleRecipeStepID, s.exampleRecipeStepVessel.ID)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with invalid recipe step ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetRecipeStepVessel(s.ctx, s.exampleRecipeID, "", s.exampleRecipeStepVessel.ID)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with invalid recipe step vessel ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetRecipeStepVessel(s.ctx, s.exampleRecipeID, s.exampleRecipeStepID, "")

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetRecipeStepVessel(s.ctx, s.exampleRecipeID, s.exampleRecipeStepID, s.exampleRecipeStepVessel.ID)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, s.exampleRecipeID, s.exampleRecipeStepID, s.exampleRecipeStepVessel.ID)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetRecipeStepVessel(s.ctx, s.exampleRecipeID, s.exampleRecipeStepID, s.exampleRecipeStepVessel.ID)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *recipeStepVesselsTestSuite) TestClient_GetRecipeStepVessels() {
	const expectedPath = "/api/v1/recipes/%s/steps/%s/vessels"

	s.Run("standard", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)

		exampleRecipeStepVesselList := fakes.BuildFakeRecipeStepVesselList()

		spec := newRequestSpec(true, http.MethodGet, "limit=20&page=1&sortBy=asc", expectedPath, s.exampleRecipeID, s.exampleRecipeStepID)
		c, _ := buildTestClientWithJSONResponse(t, spec, exampleRecipeStepVesselList)
		actual, err := c.GetRecipeStepVessels(s.ctx, s.exampleRecipeID, s.exampleRecipeStepID, filter)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, exampleRecipeStepVesselList, actual)
	})

	s.Run("with invalid recipe ID", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)

		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetRecipeStepVessels(s.ctx, "", s.exampleRecipeStepID, filter)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with invalid recipe step ID", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)

		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetRecipeStepVessels(s.ctx, s.exampleRecipeID, "", filter)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetRecipeStepVessels(s.ctx, s.exampleRecipeID, s.exampleRecipeStepID, filter)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)

		spec := newRequestSpec(true, http.MethodGet, "limit=20&page=1&sortBy=asc", expectedPath, s.exampleRecipeID, s.exampleRecipeStepID)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetRecipeStepVessels(s.ctx, s.exampleRecipeID, s.exampleRecipeStepID, filter)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *recipeStepVesselsTestSuite) TestClient_CreateRecipeStepVessel() {
	const expectedPath = "/api/v1/recipes/%s/steps/%s/vessels"

	s.Run("standard", func() {
		t := s.T()

		exampleInput := fakes.BuildFakeRecipeStepVesselCreationRequestInput()

		spec := newRequestSpec(false, http.MethodPost, "", expectedPath, s.exampleRecipeID, s.exampleRecipeStepID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleRecipeStepVessel)

		actual, err := c.CreateRecipeStepVessel(s.ctx, s.exampleRecipeID, s.exampleRecipeStepID, exampleInput)
		assert.NoError(t, err)
		assert.Equal(t, s.exampleRecipeStepVessel, actual)
	})

	s.Run("with invalid recipe ID", func() {
		t := s.T()

		exampleInput := fakes.BuildFakeRecipeStepVesselCreationRequestInput()

		c, _ := buildSimpleTestClient(t)

		actual, err := c.CreateRecipeStepVessel(s.ctx, "", s.exampleRecipeStepID, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with nil input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		actual, err := c.CreateRecipeStepVessel(s.ctx, s.exampleRecipeID, s.exampleRecipeStepID, nil)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with invalid input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		exampleInput := &types.RecipeStepVesselCreationRequestInput{}

		actual, err := c.CreateRecipeStepVessel(s.ctx, s.exampleRecipeID, s.exampleRecipeStepID, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		exampleInput := converters.ConvertRecipeStepVesselToRecipeStepVesselCreationRequestInput(s.exampleRecipeStepVessel)

		c := buildTestClientWithInvalidURL(t)

		actual, err := c.CreateRecipeStepVessel(s.ctx, s.exampleRecipeID, s.exampleRecipeStepID, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		exampleInput := converters.ConvertRecipeStepVesselToRecipeStepVesselCreationRequestInput(s.exampleRecipeStepVessel)
		c, _ := buildTestClientThatWaitsTooLong(t)

		actual, err := c.CreateRecipeStepVessel(s.ctx, s.exampleRecipeID, s.exampleRecipeStepID, exampleInput)
		assert.Empty(t, actual)
		assert.Error(t, err)
	})
}

func (s *recipeStepVesselsTestSuite) TestClient_UpdateRecipeStepVessel() {
	const expectedPathFormat = "/api/v1/recipes/%s/steps/%s/vessels/%s"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, s.exampleRecipeID, s.exampleRecipeStepID, s.exampleRecipeStepVessel.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleRecipeStepVessel)

		err := c.UpdateRecipeStepVessel(s.ctx, s.exampleRecipeID, s.exampleRecipeStepVessel)
		assert.NoError(t, err)
	})

	s.Run("with invalid recipe ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		err := c.UpdateRecipeStepVessel(s.ctx, "", s.exampleRecipeStepVessel)
		assert.Error(t, err)
	})

	s.Run("with nil input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		err := c.UpdateRecipeStepVessel(s.ctx, s.exampleRecipeID, nil)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)

		err := c.UpdateRecipeStepVessel(s.ctx, s.exampleRecipeID, s.exampleRecipeStepVessel)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		c, _ := buildTestClientThatWaitsTooLong(t)

		err := c.UpdateRecipeStepVessel(s.ctx, s.exampleRecipeID, s.exampleRecipeStepVessel)
		assert.Error(t, err)
	})
}

func (s *recipeStepVesselsTestSuite) TestClient_ArchiveRecipeStepVessel() {
	const expectedPathFormat = "/api/v1/recipes/%s/steps/%s/vessels/%s"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodDelete, "", expectedPathFormat, s.exampleRecipeID, s.exampleRecipeStepID, s.exampleRecipeStepVessel.ID)
		c, _ := buildTestClientWithStatusCodeResponse(t, spec, http.StatusOK)

		err := c.ArchiveRecipeStepVessel(s.ctx, s.exampleRecipeID, s.exampleRecipeStepID, s.exampleRecipeStepVessel.ID)
		assert.NoError(t, err)
	})

	s.Run("with invalid recipe ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		err := c.ArchiveRecipeStepVessel(s.ctx, "", s.exampleRecipeStepID, s.exampleRecipeStepVessel.ID)
		assert.Error(t, err)
	})

	s.Run("with invalid recipe step ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		err := c.ArchiveRecipeStepVessel(s.ctx, s.exampleRecipeID, "", s.exampleRecipeStepVessel.ID)
		assert.Error(t, err)
	})

	s.Run("with invalid recipe step vessel ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		err := c.ArchiveRecipeStepVessel(s.ctx, s.exampleRecipeID, s.exampleRecipeStepID, "")
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)

		err := c.ArchiveRecipeStepVessel(s.ctx, s.exampleRecipeID, s.exampleRecipeStepID, s.exampleRecipeStepVessel.ID)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		c, _ := buildTestClientThatWaitsTooLong(t)

		err := c.ArchiveRecipeStepVessel(s.ctx, s.exampleRecipeID, s.exampleRecipeStepID, s.exampleRecipeStepVessel.ID)
		assert.Error(t, err)
	})
}
