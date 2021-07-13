package httpclient

import (
	"context"
	"net/http"
	"testing"

	"gitlab.com/prixfixe/prixfixe/pkg/types"
	"gitlab.com/prixfixe/prixfixe/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

func TestValidIngredientPreparations(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(validIngredientPreparationsTestSuite))
}

type validIngredientPreparationsBaseSuite struct {
	suite.Suite

	ctx                               context.Context
	exampleValidIngredientPreparation *types.ValidIngredientPreparation
}

var _ suite.SetupTestSuite = (*validIngredientPreparationsBaseSuite)(nil)

func (s *validIngredientPreparationsBaseSuite) SetupTest() {
	s.ctx = context.Background()
	s.exampleValidIngredientPreparation = fakes.BuildFakeValidIngredientPreparation()
}

type validIngredientPreparationsTestSuite struct {
	suite.Suite

	validIngredientPreparationsBaseSuite
}

func (s *validIngredientPreparationsTestSuite) TestClient_ValidIngredientPreparationExists() {
	const expectedPathFormat = "/api/v1/valid_ingredient_preparations/%d"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodHead, "", expectedPathFormat, s.exampleValidIngredientPreparation.ID)

		c, _ := buildTestClientWithStatusCodeResponse(t, spec, http.StatusOK)
		actual, err := c.ValidIngredientPreparationExists(s.ctx, s.exampleValidIngredientPreparation.ID)

		assert.NoError(t, err)
		assert.True(t, actual)
	})

	s.Run("with invalid valid ingredient preparation ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		actual, err := c.ValidIngredientPreparationExists(s.ctx, 0)

		assert.Error(t, err)
		assert.False(t, actual)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.ValidIngredientPreparationExists(s.ctx, s.exampleValidIngredientPreparation.ID)

		assert.Error(t, err)
		assert.False(t, actual)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		c, _ := buildTestClientThatWaitsTooLong(t)
		actual, err := c.ValidIngredientPreparationExists(s.ctx, s.exampleValidIngredientPreparation.ID)

		assert.Error(t, err)
		assert.False(t, actual)
	})
}

func (s *validIngredientPreparationsTestSuite) TestClient_GetValidIngredientPreparation() {
	const expectedPathFormat = "/api/v1/valid_ingredient_preparations/%d"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, s.exampleValidIngredientPreparation.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleValidIngredientPreparation)
		actual, err := c.GetValidIngredientPreparation(s.ctx, s.exampleValidIngredientPreparation.ID)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, s.exampleValidIngredientPreparation, actual)
	})

	s.Run("with invalid valid ingredient preparation ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetValidIngredientPreparation(s.ctx, 0)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetValidIngredientPreparation(s.ctx, s.exampleValidIngredientPreparation.ID)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, s.exampleValidIngredientPreparation.ID)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetValidIngredientPreparation(s.ctx, s.exampleValidIngredientPreparation.ID)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *validIngredientPreparationsTestSuite) TestClient_GetValidIngredientPreparations() {
	const expectedPath = "/api/v1/valid_ingredient_preparations"

	s.Run("standard", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)

		exampleValidIngredientPreparationList := fakes.BuildFakeValidIngredientPreparationList()

		spec := newRequestSpec(true, http.MethodGet, "includeArchived=false&limit=20&page=1&sortBy=asc", expectedPath)
		c, _ := buildTestClientWithJSONResponse(t, spec, exampleValidIngredientPreparationList)
		actual, err := c.GetValidIngredientPreparations(s.ctx, filter)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, exampleValidIngredientPreparationList, actual)
	})

	s.Run("with error building request", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetValidIngredientPreparations(s.ctx, filter)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)

		spec := newRequestSpec(true, http.MethodGet, "includeArchived=false&limit=20&page=1&sortBy=asc", expectedPath)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetValidIngredientPreparations(s.ctx, filter)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *validIngredientPreparationsTestSuite) TestClient_CreateValidIngredientPreparation() {
	const expectedPath = "/api/v1/valid_ingredient_preparations"

	s.Run("standard", func() {
		t := s.T()

		exampleInput := fakes.BuildFakeValidIngredientPreparationCreationInput()

		spec := newRequestSpec(false, http.MethodPost, "", expectedPath)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleValidIngredientPreparation)

		actual, err := c.CreateValidIngredientPreparation(s.ctx, exampleInput)
		require.NotNil(t, actual)
		assert.NoError(t, err)

		assert.Equal(t, s.exampleValidIngredientPreparation, actual)
	})

	s.Run("with nil input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		actual, err := c.CreateValidIngredientPreparation(s.ctx, nil)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with invalid input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		exampleInput := &types.ValidIngredientPreparationCreationInput{}

		actual, err := c.CreateValidIngredientPreparation(s.ctx, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		exampleInput := fakes.BuildFakeValidIngredientPreparationCreationInputFromValidIngredientPreparation(s.exampleValidIngredientPreparation)

		c := buildTestClientWithInvalidURL(t)

		actual, err := c.CreateValidIngredientPreparation(s.ctx, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		exampleInput := fakes.BuildFakeValidIngredientPreparationCreationInputFromValidIngredientPreparation(s.exampleValidIngredientPreparation)
		c, _ := buildTestClientThatWaitsTooLong(t)

		actual, err := c.CreateValidIngredientPreparation(s.ctx, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *validIngredientPreparationsTestSuite) TestClient_UpdateValidIngredientPreparation() {
	const expectedPathFormat = "/api/v1/valid_ingredient_preparations/%d"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, s.exampleValidIngredientPreparation.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleValidIngredientPreparation)

		err := c.UpdateValidIngredientPreparation(s.ctx, s.exampleValidIngredientPreparation)
		assert.NoError(t, err)
	})

	s.Run("with nil input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		err := c.UpdateValidIngredientPreparation(s.ctx, nil)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)

		err := c.UpdateValidIngredientPreparation(s.ctx, s.exampleValidIngredientPreparation)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		c, _ := buildTestClientThatWaitsTooLong(t)

		err := c.UpdateValidIngredientPreparation(s.ctx, s.exampleValidIngredientPreparation)
		assert.Error(t, err)
	})
}

func (s *validIngredientPreparationsTestSuite) TestClient_ArchiveValidIngredientPreparation() {
	const expectedPathFormat = "/api/v1/valid_ingredient_preparations/%d"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodDelete, "", expectedPathFormat, s.exampleValidIngredientPreparation.ID)
		c, _ := buildTestClientWithStatusCodeResponse(t, spec, http.StatusOK)

		err := c.ArchiveValidIngredientPreparation(s.ctx, s.exampleValidIngredientPreparation.ID)
		assert.NoError(t, err)
	})

	s.Run("with invalid valid ingredient preparation ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		err := c.ArchiveValidIngredientPreparation(s.ctx, 0)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)

		err := c.ArchiveValidIngredientPreparation(s.ctx, s.exampleValidIngredientPreparation.ID)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		c, _ := buildTestClientThatWaitsTooLong(t)

		err := c.ArchiveValidIngredientPreparation(s.ctx, s.exampleValidIngredientPreparation.ID)
		assert.Error(t, err)
	})
}

func (s *validIngredientPreparationsTestSuite) TestClient_GetAuditLogForValidIngredientPreparation() {
	const (
		expectedPath   = "/api/v1/valid_ingredient_preparations/%d/audit"
		expectedMethod = http.MethodGet
	)

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, expectedMethod, "", expectedPath, s.exampleValidIngredientPreparation.ID)
		exampleAuditLogEntryList := fakes.BuildFakeAuditLogEntryList().Entries

		c, _ := buildTestClientWithJSONResponse(t, spec, exampleAuditLogEntryList)

		actual, err := c.GetAuditLogForValidIngredientPreparation(s.ctx, s.exampleValidIngredientPreparation.ID)
		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, exampleAuditLogEntryList, actual)
	})

	s.Run("with invalid valid ingredient preparation ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		actual, err := c.GetAuditLogForValidIngredientPreparation(s.ctx, 0)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)

		actual, err := c.GetAuditLogForValidIngredientPreparation(s.ctx, s.exampleValidIngredientPreparation.ID)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		c, _ := buildTestClientThatWaitsTooLong(t)

		actual, err := c.GetAuditLogForValidIngredientPreparation(s.ctx, s.exampleValidIngredientPreparation.ID)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}
