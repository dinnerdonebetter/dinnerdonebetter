package httpclient

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"gitlab.com/prixfixe/prixfixe/pkg/types"
	"gitlab.com/prixfixe/prixfixe/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

func TestValidIngredients(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(validIngredientsTestSuite))
}

type validIngredientsBaseSuite struct {
	suite.Suite

	ctx                    context.Context
	exampleValidIngredient *types.ValidIngredient
}

var _ suite.SetupTestSuite = (*validIngredientsBaseSuite)(nil)

func (s *validIngredientsBaseSuite) SetupTest() {
	s.ctx = context.Background()
	s.exampleValidIngredient = fakes.BuildFakeValidIngredient()
}

type validIngredientsTestSuite struct {
	suite.Suite

	validIngredientsBaseSuite
}

func (s *validIngredientsTestSuite) TestClient_ValidIngredientExists() {
	const expectedPathFormat = "/api/v1/valid_ingredients/%d"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodHead, "", expectedPathFormat, s.exampleValidIngredient.ID)

		c, _ := buildTestClientWithStatusCodeResponse(t, spec, http.StatusOK)
		actual, err := c.ValidIngredientExists(s.ctx, s.exampleValidIngredient.ID)

		assert.NoError(t, err)
		assert.True(t, actual)
	})

	s.Run("with invalid valid ingredient ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		actual, err := c.ValidIngredientExists(s.ctx, 0)

		assert.Error(t, err)
		assert.False(t, actual)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.ValidIngredientExists(s.ctx, s.exampleValidIngredient.ID)

		assert.Error(t, err)
		assert.False(t, actual)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		c, _ := buildTestClientThatWaitsTooLong(t)
		actual, err := c.ValidIngredientExists(s.ctx, s.exampleValidIngredient.ID)

		assert.Error(t, err)
		assert.False(t, actual)
	})
}

func (s *validIngredientsTestSuite) TestClient_GetValidIngredient() {
	const expectedPathFormat = "/api/v1/valid_ingredients/%d"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, s.exampleValidIngredient.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleValidIngredient)
		actual, err := c.GetValidIngredient(s.ctx, s.exampleValidIngredient.ID)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, s.exampleValidIngredient, actual)
	})

	s.Run("with invalid valid ingredient ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetValidIngredient(s.ctx, 0)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetValidIngredient(s.ctx, s.exampleValidIngredient.ID)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, s.exampleValidIngredient.ID)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetValidIngredient(s.ctx, s.exampleValidIngredient.ID)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *validIngredientsTestSuite) TestClient_GetValidIngredients() {
	const expectedPath = "/api/v1/valid_ingredients"

	s.Run("standard", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)

		exampleValidIngredientList := fakes.BuildFakeValidIngredientList()

		spec := newRequestSpec(true, http.MethodGet, "includeArchived=false&limit=20&page=1&sortBy=asc", expectedPath)
		c, _ := buildTestClientWithJSONResponse(t, spec, exampleValidIngredientList)
		actual, err := c.GetValidIngredients(s.ctx, filter)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, exampleValidIngredientList, actual)
	})

	s.Run("with error building request", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetValidIngredients(s.ctx, filter)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)

		spec := newRequestSpec(true, http.MethodGet, "includeArchived=false&limit=20&page=1&sortBy=asc", expectedPath)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetValidIngredients(s.ctx, filter)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *validIngredientsTestSuite) TestClient_SearchValidIngredients() {
	const expectedPath = "/api/v1/valid_ingredients/search"

	exampleQuery := "whatever"

	s.Run("standard", func() {
		t := s.T()

		exampleValidPreparation := fakes.BuildFakeValidPreparation()
		exampleValidIngredientList := fakes.BuildFakeValidIngredientList()

		spec := newRequestSpec(true, http.MethodGet, fmt.Sprintf("limit=20&pid=%d&q=whatever", exampleValidPreparation.ID), expectedPath)
		c, _ := buildTestClientWithJSONResponse(t, spec, exampleValidIngredientList.ValidIngredients)
		actual, err := c.SearchValidIngredients(s.ctx, exampleQuery, exampleValidPreparation.ID, 0)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, exampleValidIngredientList.ValidIngredients, actual)
	})

	s.Run("with invalid valid preparation ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		actual, err := c.SearchValidIngredients(s.ctx, "", 0, 0)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with empty query", func() {
		t := s.T()

		exampleValidPreparation := fakes.BuildFakeValidPreparation()
		c, _ := buildSimpleTestClient(t)
		actual, err := c.SearchValidIngredients(s.ctx, "", exampleValidPreparation.ID, 0)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		exampleValidPreparation := fakes.BuildFakeValidPreparation()
		c := buildTestClientWithInvalidURL(t)

		actual, err := c.SearchValidIngredients(s.ctx, exampleQuery, exampleValidPreparation.ID, 0)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with bad response from server", func() {
		t := s.T()

		exampleValidPreparation := fakes.BuildFakeValidPreparation()
		spec := newRequestSpec(true, http.MethodGet, fmt.Sprintf("limit=20&pid=%d&q=whatever", exampleValidPreparation.ID), expectedPath)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.SearchValidIngredients(s.ctx, exampleQuery, exampleValidPreparation.ID, 0)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *validIngredientsTestSuite) TestClient_CreateValidIngredient() {
	const expectedPath = "/api/v1/valid_ingredients"

	s.Run("standard", func() {
		t := s.T()

		exampleInput := fakes.BuildFakeValidIngredientCreationInput()

		spec := newRequestSpec(false, http.MethodPost, "", expectedPath)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleValidIngredient)

		actual, err := c.CreateValidIngredient(s.ctx, exampleInput)
		require.NotNil(t, actual)
		assert.NoError(t, err)

		assert.Equal(t, s.exampleValidIngredient, actual)
	})

	s.Run("with nil input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		actual, err := c.CreateValidIngredient(s.ctx, nil)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with invalid input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		exampleInput := &types.ValidIngredientCreationInput{}

		actual, err := c.CreateValidIngredient(s.ctx, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		exampleInput := fakes.BuildFakeValidIngredientCreationInputFromValidIngredient(s.exampleValidIngredient)

		c := buildTestClientWithInvalidURL(t)

		actual, err := c.CreateValidIngredient(s.ctx, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		exampleInput := fakes.BuildFakeValidIngredientCreationInputFromValidIngredient(s.exampleValidIngredient)
		c, _ := buildTestClientThatWaitsTooLong(t)

		actual, err := c.CreateValidIngredient(s.ctx, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *validIngredientsTestSuite) TestClient_UpdateValidIngredient() {
	const expectedPathFormat = "/api/v1/valid_ingredients/%d"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, s.exampleValidIngredient.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleValidIngredient)

		err := c.UpdateValidIngredient(s.ctx, s.exampleValidIngredient)
		assert.NoError(t, err)
	})

	s.Run("with nil input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		err := c.UpdateValidIngredient(s.ctx, nil)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)

		err := c.UpdateValidIngredient(s.ctx, s.exampleValidIngredient)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		c, _ := buildTestClientThatWaitsTooLong(t)

		err := c.UpdateValidIngredient(s.ctx, s.exampleValidIngredient)
		assert.Error(t, err)
	})
}

func (s *validIngredientsTestSuite) TestClient_ArchiveValidIngredient() {
	const expectedPathFormat = "/api/v1/valid_ingredients/%d"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodDelete, "", expectedPathFormat, s.exampleValidIngredient.ID)
		c, _ := buildTestClientWithStatusCodeResponse(t, spec, http.StatusOK)

		err := c.ArchiveValidIngredient(s.ctx, s.exampleValidIngredient.ID)
		assert.NoError(t, err)
	})

	s.Run("with invalid valid ingredient ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		err := c.ArchiveValidIngredient(s.ctx, 0)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)

		err := c.ArchiveValidIngredient(s.ctx, s.exampleValidIngredient.ID)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		c, _ := buildTestClientThatWaitsTooLong(t)

		err := c.ArchiveValidIngredient(s.ctx, s.exampleValidIngredient.ID)
		assert.Error(t, err)
	})
}

func (s *validIngredientsTestSuite) TestClient_GetAuditLogForValidIngredient() {
	const (
		expectedPath   = "/api/v1/valid_ingredients/%d/audit"
		expectedMethod = http.MethodGet
	)

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, expectedMethod, "", expectedPath, s.exampleValidIngredient.ID)
		exampleAuditLogEntryList := fakes.BuildFakeAuditLogEntryList().Entries

		c, _ := buildTestClientWithJSONResponse(t, spec, exampleAuditLogEntryList)

		actual, err := c.GetAuditLogForValidIngredient(s.ctx, s.exampleValidIngredient.ID)
		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, exampleAuditLogEntryList, actual)
	})

	s.Run("with invalid valid ingredient ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		actual, err := c.GetAuditLogForValidIngredient(s.ctx, 0)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)

		actual, err := c.GetAuditLogForValidIngredient(s.ctx, s.exampleValidIngredient.ID)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		c, _ := buildTestClientThatWaitsTooLong(t)

		actual, err := c.GetAuditLogForValidIngredient(s.ctx, s.exampleValidIngredient.ID)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}
