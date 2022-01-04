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

func TestRecipeStepProducts(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(recipeStepProductsTestSuite))
}

type recipeStepProductsBaseSuite struct {
	suite.Suite
	ctx                      context.Context
	exampleRecipeStepProduct *types.RecipeStepProduct
	exampleRecipeID          string
	exampleRecipeStepID      string
}

var _ suite.SetupTestSuite = (*recipeStepProductsBaseSuite)(nil)

func (s *recipeStepProductsBaseSuite) SetupTest() {
	s.ctx = context.Background()
	s.exampleRecipeID = fakes.BuildFakeID()
	s.exampleRecipeStepID = fakes.BuildFakeID()
	s.exampleRecipeStepProduct = fakes.BuildFakeRecipeStepProduct()
	s.exampleRecipeStepProduct.BelongsToRecipeStep = s.exampleRecipeStepID
}

type recipeStepProductsTestSuite struct {
	suite.Suite

	recipeStepProductsBaseSuite
}

func (s *recipeStepProductsTestSuite) TestClient_GetRecipeStepProduct() {
	const expectedPathFormat = "/api/v1/recipes/%s/recipe_steps/%s/recipe_step_products/%s"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, s.exampleRecipeID, s.exampleRecipeStepID, s.exampleRecipeStepProduct.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleRecipeStepProduct)
		actual, err := c.GetRecipeStepProduct(s.ctx, s.exampleRecipeID, s.exampleRecipeStepID, s.exampleRecipeStepProduct.ID)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, s.exampleRecipeStepProduct, actual)
	})

	s.Run("with invalid recipe ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetRecipeStepProduct(s.ctx, "", s.exampleRecipeStepID, s.exampleRecipeStepProduct.ID)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with invalid recipe step ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetRecipeStepProduct(s.ctx, s.exampleRecipeID, "", s.exampleRecipeStepProduct.ID)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with invalid recipe step product ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetRecipeStepProduct(s.ctx, s.exampleRecipeID, s.exampleRecipeStepID, "")

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetRecipeStepProduct(s.ctx, s.exampleRecipeID, s.exampleRecipeStepID, s.exampleRecipeStepProduct.ID)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, s.exampleRecipeID, s.exampleRecipeStepID, s.exampleRecipeStepProduct.ID)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetRecipeStepProduct(s.ctx, s.exampleRecipeID, s.exampleRecipeStepID, s.exampleRecipeStepProduct.ID)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *recipeStepProductsTestSuite) TestClient_GetRecipeStepProducts() {
	const expectedPath = "/api/v1/recipes/%s/recipe_steps/%s/recipe_step_products"

	s.Run("standard", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)

		exampleRecipeStepProductList := fakes.BuildFakeRecipeStepProductList()

		spec := newRequestSpec(true, http.MethodGet, "includeArchived=false&limit=20&page=1&sortBy=asc", expectedPath, s.exampleRecipeID, s.exampleRecipeStepID)
		c, _ := buildTestClientWithJSONResponse(t, spec, exampleRecipeStepProductList)
		actual, err := c.GetRecipeStepProducts(s.ctx, s.exampleRecipeID, s.exampleRecipeStepID, filter)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, exampleRecipeStepProductList, actual)
	})

	s.Run("with invalid recipe ID", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)

		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetRecipeStepProducts(s.ctx, "", s.exampleRecipeStepID, filter)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with invalid recipe step ID", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)

		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetRecipeStepProducts(s.ctx, s.exampleRecipeID, "", filter)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetRecipeStepProducts(s.ctx, s.exampleRecipeID, s.exampleRecipeStepID, filter)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)

		spec := newRequestSpec(true, http.MethodGet, "includeArchived=false&limit=20&page=1&sortBy=asc", expectedPath, s.exampleRecipeID, s.exampleRecipeStepID)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetRecipeStepProducts(s.ctx, s.exampleRecipeID, s.exampleRecipeStepID, filter)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *recipeStepProductsTestSuite) TestClient_CreateRecipeStepProduct() {
	const expectedPath = "/api/v1/recipes/%s/recipe_steps/%s/recipe_step_products"

	s.Run("standard", func() {
		t := s.T()

		exampleInput := fakes.BuildFakeRecipeStepProductCreationRequestInput()
		exampleInput.BelongsToRecipeStep = s.exampleRecipeStepID

		spec := newRequestSpec(false, http.MethodPost, "", expectedPath, s.exampleRecipeID, s.exampleRecipeStepID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleRecipeStepProduct)

		actual, err := c.CreateRecipeStepProduct(s.ctx, s.exampleRecipeID, exampleInput)
		require.NotEmpty(t, actual)
		assert.NoError(t, err)

		assert.Equal(t, s.exampleRecipeStepProduct, actual)
	})

	s.Run("with invalid recipe ID", func() {
		t := s.T()

		exampleInput := fakes.BuildFakeRecipeStepProductCreationRequestInput()
		exampleInput.BelongsToRecipeStep = s.exampleRecipeStepID

		c, _ := buildSimpleTestClient(t)

		actual, err := c.CreateRecipeStepProduct(s.ctx, "", exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with nil input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		actual, err := c.CreateRecipeStepProduct(s.ctx, s.exampleRecipeID, nil)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with invalid input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		exampleInput := &types.RecipeStepProductCreationRequestInput{}

		actual, err := c.CreateRecipeStepProduct(s.ctx, s.exampleRecipeID, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		exampleInput := fakes.BuildFakeRecipeStepProductCreationRequestInputFromRecipeStepProduct(s.exampleRecipeStepProduct)

		c := buildTestClientWithInvalidURL(t)

		actual, err := c.CreateRecipeStepProduct(s.ctx, s.exampleRecipeID, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		exampleInput := fakes.BuildFakeRecipeStepProductCreationRequestInputFromRecipeStepProduct(s.exampleRecipeStepProduct)
		c, _ := buildTestClientThatWaitsTooLong(t)

		actual, err := c.CreateRecipeStepProduct(s.ctx, s.exampleRecipeID, exampleInput)
		assert.Empty(t, actual)
		assert.Error(t, err)
	})
}

func (s *recipeStepProductsTestSuite) TestClient_UpdateRecipeStepProduct() {
	const expectedPathFormat = "/api/v1/recipes/%s/recipe_steps/%s/recipe_step_products/%s"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, s.exampleRecipeID, s.exampleRecipeStepID, s.exampleRecipeStepProduct.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleRecipeStepProduct)

		err := c.UpdateRecipeStepProduct(s.ctx, s.exampleRecipeID, s.exampleRecipeStepProduct)
		assert.NoError(t, err)
	})

	s.Run("with invalid recipe ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		err := c.UpdateRecipeStepProduct(s.ctx, "", s.exampleRecipeStepProduct)
		assert.Error(t, err)
	})

	s.Run("with nil input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		err := c.UpdateRecipeStepProduct(s.ctx, s.exampleRecipeID, nil)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)

		err := c.UpdateRecipeStepProduct(s.ctx, s.exampleRecipeID, s.exampleRecipeStepProduct)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		c, _ := buildTestClientThatWaitsTooLong(t)

		err := c.UpdateRecipeStepProduct(s.ctx, s.exampleRecipeID, s.exampleRecipeStepProduct)
		assert.Error(t, err)
	})
}

func (s *recipeStepProductsTestSuite) TestClient_ArchiveRecipeStepProduct() {
	const expectedPathFormat = "/api/v1/recipes/%s/recipe_steps/%s/recipe_step_products/%s"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodDelete, "", expectedPathFormat, s.exampleRecipeID, s.exampleRecipeStepID, s.exampleRecipeStepProduct.ID)
		c, _ := buildTestClientWithStatusCodeResponse(t, spec, http.StatusOK)

		err := c.ArchiveRecipeStepProduct(s.ctx, s.exampleRecipeID, s.exampleRecipeStepID, s.exampleRecipeStepProduct.ID)
		assert.NoError(t, err)
	})

	s.Run("with invalid recipe ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		err := c.ArchiveRecipeStepProduct(s.ctx, "", s.exampleRecipeStepID, s.exampleRecipeStepProduct.ID)
		assert.Error(t, err)
	})

	s.Run("with invalid recipe step ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		err := c.ArchiveRecipeStepProduct(s.ctx, s.exampleRecipeID, "", s.exampleRecipeStepProduct.ID)
		assert.Error(t, err)
	})

	s.Run("with invalid recipe step product ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		err := c.ArchiveRecipeStepProduct(s.ctx, s.exampleRecipeID, s.exampleRecipeStepID, "")
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)

		err := c.ArchiveRecipeStepProduct(s.ctx, s.exampleRecipeID, s.exampleRecipeStepID, s.exampleRecipeStepProduct.ID)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		c, _ := buildTestClientThatWaitsTooLong(t)

		err := c.ArchiveRecipeStepProduct(s.ctx, s.exampleRecipeID, s.exampleRecipeStepID, s.exampleRecipeStepProduct.ID)
		assert.Error(t, err)
	})
}
