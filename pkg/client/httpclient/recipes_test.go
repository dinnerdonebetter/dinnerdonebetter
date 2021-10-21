package httpclient

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"gitlab.com/prixfixe/prixfixe/pkg/types"
	"gitlab.com/prixfixe/prixfixe/pkg/types/fakes"
)

func TestRecipes(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(recipesTestSuite))
}

type recipesBaseSuite struct {
	suite.Suite

	ctx           context.Context
	exampleRecipe *types.Recipe
}

var _ suite.SetupTestSuite = (*recipesBaseSuite)(nil)

func (s *recipesBaseSuite) SetupTest() {
	s.ctx = context.Background()
	s.exampleRecipe = fakes.BuildFakeRecipe()
}

type recipesTestSuite struct {
	suite.Suite

	recipesBaseSuite
}

func (s *recipesTestSuite) TestClient_GetRecipe() {
	const expectedPathFormat = "/api/v1/recipes/%s"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, s.exampleRecipe.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleRecipe)
		actual, err := c.GetRecipe(s.ctx, s.exampleRecipe.ID)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, s.exampleRecipe, actual)
	})

	s.Run("with invalid recipe ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetRecipe(s.ctx, "")

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetRecipe(s.ctx, s.exampleRecipe.ID)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, s.exampleRecipe.ID)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetRecipe(s.ctx, s.exampleRecipe.ID)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *recipesTestSuite) TestClient_GetRecipes() {
	const expectedPath = "/api/v1/recipes"

	s.Run("standard", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)

		exampleRecipeList := fakes.BuildFakeRecipeList()

		spec := newRequestSpec(true, http.MethodGet, "includeArchived=false&limit=20&page=1&sortBy=asc", expectedPath)
		c, _ := buildTestClientWithJSONResponse(t, spec, exampleRecipeList)
		actual, err := c.GetRecipes(s.ctx, filter)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, exampleRecipeList, actual)
	})

	s.Run("with error building request", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetRecipes(s.ctx, filter)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)

		spec := newRequestSpec(true, http.MethodGet, "includeArchived=false&limit=20&page=1&sortBy=asc", expectedPath)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetRecipes(s.ctx, filter)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *recipesTestSuite) TestClient_CreateRecipe() {
	const expectedPath = "/api/v1/recipes"

	s.Run("standard", func() {
		t := s.T()

		exampleInput := fakes.BuildFakeRecipeCreationRequestInput()
		exampleInput.BelongsToAccount = ""

		spec := newRequestSpec(false, http.MethodPost, "", expectedPath)
		c, _ := buildTestClientWithJSONResponse(t, spec, &types.PreWriteResponse{ID: s.exampleRecipe.ID})

		actual, err := c.CreateRecipe(s.ctx, exampleInput)
		require.NotEmpty(t, actual)
		assert.NoError(t, err)

		assert.Equal(t, s.exampleRecipe.ID, actual)
	})

	s.Run("with nil input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		actual, err := c.CreateRecipe(s.ctx, nil)
		assert.Empty(t, actual)
		assert.Error(t, err)
	})

	s.Run("with invalid input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		exampleInput := &types.RecipeCreationRequestInput{}

		actual, err := c.CreateRecipe(s.ctx, exampleInput)
		assert.Empty(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		exampleInput := fakes.BuildFakeRecipeCreationRequestInputFromRecipe(s.exampleRecipe)

		c := buildTestClientWithInvalidURL(t)

		actual, err := c.CreateRecipe(s.ctx, exampleInput)
		assert.Empty(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		exampleInput := fakes.BuildFakeRecipeCreationRequestInputFromRecipe(s.exampleRecipe)
		c, _ := buildTestClientThatWaitsTooLong(t)

		actual, err := c.CreateRecipe(s.ctx, exampleInput)
		assert.Empty(t, actual)
		assert.Error(t, err)
	})
}

func (s *recipesTestSuite) TestClient_UpdateRecipe() {
	const expectedPathFormat = "/api/v1/recipes/%s"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, s.exampleRecipe.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleRecipe)

		err := c.UpdateRecipe(s.ctx, s.exampleRecipe)
		assert.NoError(t, err)
	})

	s.Run("with nil input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		err := c.UpdateRecipe(s.ctx, nil)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)

		err := c.UpdateRecipe(s.ctx, s.exampleRecipe)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		c, _ := buildTestClientThatWaitsTooLong(t)

		err := c.UpdateRecipe(s.ctx, s.exampleRecipe)
		assert.Error(t, err)
	})
}

func (s *recipesTestSuite) TestClient_ArchiveRecipe() {
	const expectedPathFormat = "/api/v1/recipes/%s"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodDelete, "", expectedPathFormat, s.exampleRecipe.ID)
		c, _ := buildTestClientWithStatusCodeResponse(t, spec, http.StatusOK)

		err := c.ArchiveRecipe(s.ctx, s.exampleRecipe.ID)
		assert.NoError(t, err)
	})

	s.Run("with invalid recipe ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		err := c.ArchiveRecipe(s.ctx, "")
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)

		err := c.ArchiveRecipe(s.ctx, s.exampleRecipe.ID)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		c, _ := buildTestClientThatWaitsTooLong(t)

		err := c.ArchiveRecipe(s.ctx, s.exampleRecipe.ID)
		assert.Error(t, err)
	})
}
