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

func TestRecipeStepProducts(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(recipeStepProductsTestSuite))
}

type recipeStepProductsBaseSuite struct {
	suite.Suite
	ctx                      context.Context
	exampleRecipeStepProduct *types.RecipeStepProduct
	exampleRecipeID          uint64
	exampleRecipeStepID      uint64
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

func (s *recipeStepProductsTestSuite) TestClient_RecipeStepProductExists() {
	const expectedPathFormat = "/api/v1/recipes/%d/recipe_steps/%d/recipe_step_products/%d"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodHead, "", expectedPathFormat, s.exampleRecipeID, s.exampleRecipeStepID, s.exampleRecipeStepProduct.ID)

		c, _ := buildTestClientWithStatusCodeResponse(t, spec, http.StatusOK)
		actual, err := c.RecipeStepProductExists(s.ctx, s.exampleRecipeID, s.exampleRecipeStepID, s.exampleRecipeStepProduct.ID)

		assert.NoError(t, err)
		assert.True(t, actual)
	})

	s.Run("with invalid recipe ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		actual, err := c.RecipeStepProductExists(s.ctx, 0, s.exampleRecipeStepID, s.exampleRecipeStepProduct.ID)

		assert.Error(t, err)
		assert.False(t, actual)
	})

	s.Run("with invalid recipe step ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		actual, err := c.RecipeStepProductExists(s.ctx, s.exampleRecipeID, 0, s.exampleRecipeStepProduct.ID)

		assert.Error(t, err)
		assert.False(t, actual)
	})

	s.Run("with invalid recipe step product ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		actual, err := c.RecipeStepProductExists(s.ctx, s.exampleRecipeID, s.exampleRecipeStepID, 0)

		assert.Error(t, err)
		assert.False(t, actual)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.RecipeStepProductExists(s.ctx, s.exampleRecipeID, s.exampleRecipeStepID, s.exampleRecipeStepProduct.ID)

		assert.Error(t, err)
		assert.False(t, actual)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		c, _ := buildTestClientThatWaitsTooLong(t)
		actual, err := c.RecipeStepProductExists(s.ctx, s.exampleRecipeID, s.exampleRecipeStepID, s.exampleRecipeStepProduct.ID)

		assert.Error(t, err)
		assert.False(t, actual)
	})
}

func (s *recipeStepProductsTestSuite) TestClient_GetRecipeStepProduct() {
	const expectedPathFormat = "/api/v1/recipes/%d/recipe_steps/%d/recipe_step_products/%d"

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
		actual, err := c.GetRecipeStepProduct(s.ctx, 0, s.exampleRecipeStepID, s.exampleRecipeStepProduct.ID)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with invalid recipe step ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetRecipeStepProduct(s.ctx, s.exampleRecipeID, 0, s.exampleRecipeStepProduct.ID)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with invalid recipe step product ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetRecipeStepProduct(s.ctx, s.exampleRecipeID, s.exampleRecipeStepID, 0)

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
	const expectedPath = "/api/v1/recipes/%d/recipe_steps/%d/recipe_step_products"

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
		actual, err := c.GetRecipeStepProducts(s.ctx, 0, s.exampleRecipeStepID, filter)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with invalid recipe step ID", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)

		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetRecipeStepProducts(s.ctx, s.exampleRecipeID, 0, filter)

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
	const expectedPath = "/api/v1/recipes/%d/recipe_steps/%d/recipe_step_products"

	s.Run("standard", func() {
		t := s.T()

		exampleInput := fakes.BuildFakeRecipeStepProductCreationInput()
		exampleInput.BelongsToRecipeStep = s.exampleRecipeStepID

		spec := newRequestSpec(false, http.MethodPost, "", expectedPath, s.exampleRecipeID, s.exampleRecipeStepID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleRecipeStepProduct)

		actual, err := c.CreateRecipeStepProduct(s.ctx, s.exampleRecipeID, exampleInput)
		require.NotNil(t, actual)
		assert.NoError(t, err)

		assert.Equal(t, s.exampleRecipeStepProduct, actual)
	})

	s.Run("with invalid recipe ID", func() {
		t := s.T()

		exampleInput := fakes.BuildFakeRecipeStepProductCreationInput()
		exampleInput.BelongsToRecipeStep = s.exampleRecipeStepID

		c, _ := buildSimpleTestClient(t)

		actual, err := c.CreateRecipeStepProduct(s.ctx, 0, exampleInput)
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
		exampleInput := &types.RecipeStepProductCreationInput{}

		actual, err := c.CreateRecipeStepProduct(s.ctx, s.exampleRecipeID, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		exampleInput := fakes.BuildFakeRecipeStepProductCreationInputFromRecipeStepProduct(s.exampleRecipeStepProduct)

		c := buildTestClientWithInvalidURL(t)

		actual, err := c.CreateRecipeStepProduct(s.ctx, s.exampleRecipeID, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		exampleInput := fakes.BuildFakeRecipeStepProductCreationInputFromRecipeStepProduct(s.exampleRecipeStepProduct)
		c, _ := buildTestClientThatWaitsTooLong(t)

		actual, err := c.CreateRecipeStepProduct(s.ctx, s.exampleRecipeID, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *recipeStepProductsTestSuite) TestClient_UpdateRecipeStepProduct() {
	const expectedPathFormat = "/api/v1/recipes/%d/recipe_steps/%d/recipe_step_products/%d"

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

		err := c.UpdateRecipeStepProduct(s.ctx, 0, s.exampleRecipeStepProduct)
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
	const expectedPathFormat = "/api/v1/recipes/%d/recipe_steps/%d/recipe_step_products/%d"

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

		err := c.ArchiveRecipeStepProduct(s.ctx, 0, s.exampleRecipeStepID, s.exampleRecipeStepProduct.ID)
		assert.Error(t, err)
	})

	s.Run("with invalid recipe step ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		err := c.ArchiveRecipeStepProduct(s.ctx, s.exampleRecipeID, 0, s.exampleRecipeStepProduct.ID)
		assert.Error(t, err)
	})

	s.Run("with invalid recipe step product ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		err := c.ArchiveRecipeStepProduct(s.ctx, s.exampleRecipeID, s.exampleRecipeStepID, 0)
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

func (s *recipeStepProductsTestSuite) TestClient_GetAuditLogForRecipeStepProduct() {
	const (
		expectedPath   = "/api/v1/recipes/%d/recipe_steps/%d/recipe_step_products/%d/audit"
		expectedMethod = http.MethodGet
	)

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, expectedMethod, "", expectedPath, s.exampleRecipeID, s.exampleRecipeStepID, s.exampleRecipeStepProduct.ID)
		exampleAuditLogEntryList := fakes.BuildFakeAuditLogEntryList().Entries

		c, _ := buildTestClientWithJSONResponse(t, spec, exampleAuditLogEntryList)

		actual, err := c.GetAuditLogForRecipeStepProduct(s.ctx, s.exampleRecipeID, s.exampleRecipeStepID, s.exampleRecipeStepProduct.ID)
		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, exampleAuditLogEntryList, actual)
	})

	s.Run("with invalid recipe ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		actual, err := c.GetAuditLogForRecipeStepProduct(s.ctx, 0, s.exampleRecipeStepID, s.exampleRecipeStepProduct.ID)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with invalid recipe step ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		actual, err := c.GetAuditLogForRecipeStepProduct(s.ctx, s.exampleRecipeID, 0, s.exampleRecipeStepProduct.ID)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with invalid recipe step product ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		actual, err := c.GetAuditLogForRecipeStepProduct(s.ctx, s.exampleRecipeID, s.exampleRecipeStepID, 0)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)

		actual, err := c.GetAuditLogForRecipeStepProduct(s.ctx, s.exampleRecipeID, s.exampleRecipeStepID, s.exampleRecipeStepProduct.ID)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		c, _ := buildTestClientThatWaitsTooLong(t)

		actual, err := c.GetAuditLogForRecipeStepProduct(s.ctx, s.exampleRecipeID, s.exampleRecipeStepID, s.exampleRecipeStepProduct.ID)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}
