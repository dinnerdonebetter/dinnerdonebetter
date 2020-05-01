package httpserver

import (
	"context"
	"fmt"
	"testing"

	models "gitlab.com/prixfixe/prixfixe/models/v1"
	ingredienttagmappingsservice "gitlab.com/prixfixe/prixfixe/services/v1/ingredienttagmappings"
	invitationsservice "gitlab.com/prixfixe/prixfixe/services/v1/invitations"
	iterationmediasservice "gitlab.com/prixfixe/prixfixe/services/v1/iterationmedias"
	oauth2clientsservice "gitlab.com/prixfixe/prixfixe/services/v1/oauth2clients"
	recipeiterationsservice "gitlab.com/prixfixe/prixfixe/services/v1/recipeiterations"
	recipeiterationstepsservice "gitlab.com/prixfixe/prixfixe/services/v1/recipeiterationsteps"
	recipesservice "gitlab.com/prixfixe/prixfixe/services/v1/recipes"
	recipestepingredientsservice "gitlab.com/prixfixe/prixfixe/services/v1/recipestepingredients"
	recipesteppreparationsservice "gitlab.com/prixfixe/prixfixe/services/v1/recipesteppreparations"
	recipestepsservice "gitlab.com/prixfixe/prixfixe/services/v1/recipesteps"
	recipetagsservice "gitlab.com/prixfixe/prixfixe/services/v1/recipetags"
	reportsservice "gitlab.com/prixfixe/prixfixe/services/v1/reports"
	requiredpreparationinstrumentsservice "gitlab.com/prixfixe/prixfixe/services/v1/requiredpreparationinstruments"
	usersservice "gitlab.com/prixfixe/prixfixe/services/v1/users"
	validingredientpreparationsservice "gitlab.com/prixfixe/prixfixe/services/v1/validingredientpreparations"
	validingredientsservice "gitlab.com/prixfixe/prixfixe/services/v1/validingredients"
	validingredienttagsservice "gitlab.com/prixfixe/prixfixe/services/v1/validingredienttags"
	validinstrumentsservice "gitlab.com/prixfixe/prixfixe/services/v1/validinstruments"
	validpreparationsservice "gitlab.com/prixfixe/prixfixe/services/v1/validpreparations"
	webhooksservice "gitlab.com/prixfixe/prixfixe/services/v1/webhooks"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
	"gitlab.com/verygoodsoftwarenotvirus/logging/v1/noop"
)

func TestProvideValidInstrumentsServiceValidInstrumentIDFetcher(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		_ = ProvideValidInstrumentsServiceValidInstrumentIDFetcher(noop.ProvideNoopLogger())
	})
}

func TestProvideValidIngredientsServiceValidIngredientIDFetcher(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		_ = ProvideValidIngredientsServiceValidIngredientIDFetcher(noop.ProvideNoopLogger())
	})
}

func TestProvideValidIngredientTagsServiceValidIngredientTagIDFetcher(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		_ = ProvideValidIngredientTagsServiceValidIngredientTagIDFetcher(noop.ProvideNoopLogger())
	})
}

func TestProvideIngredientTagMappingsServiceValidIngredientIDFetcher(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		_ = ProvideIngredientTagMappingsServiceValidIngredientIDFetcher(noop.ProvideNoopLogger())
	})
}

func TestProvideIngredientTagMappingsServiceIngredientTagMappingIDFetcher(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		_ = ProvideIngredientTagMappingsServiceIngredientTagMappingIDFetcher(noop.ProvideNoopLogger())
	})
}

func TestProvideValidPreparationsServiceValidPreparationIDFetcher(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		_ = ProvideValidPreparationsServiceValidPreparationIDFetcher(noop.ProvideNoopLogger())
	})
}

func TestProvideRequiredPreparationInstrumentsServiceValidPreparationIDFetcher(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		_ = ProvideRequiredPreparationInstrumentsServiceValidPreparationIDFetcher(noop.ProvideNoopLogger())
	})
}

func TestProvideRequiredPreparationInstrumentsServiceRequiredPreparationInstrumentIDFetcher(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		_ = ProvideRequiredPreparationInstrumentsServiceRequiredPreparationInstrumentIDFetcher(noop.ProvideNoopLogger())
	})
}

func TestProvideValidIngredientPreparationsServiceValidIngredientIDFetcher(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		_ = ProvideValidIngredientPreparationsServiceValidIngredientIDFetcher(noop.ProvideNoopLogger())
	})
}

func TestProvideValidIngredientPreparationsServiceValidIngredientPreparationIDFetcher(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		_ = ProvideValidIngredientPreparationsServiceValidIngredientPreparationIDFetcher(noop.ProvideNoopLogger())
	})
}

func TestProvideRecipesServiceUserIDFetcher(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		_ = ProvideRecipesServiceUserIDFetcher()
	})
}

func TestProvideRecipesServiceRecipeIDFetcher(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		_ = ProvideRecipesServiceRecipeIDFetcher(noop.ProvideNoopLogger())
	})
}

func TestProvideRecipeTagsServiceUserIDFetcher(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		_ = ProvideRecipeTagsServiceUserIDFetcher()
	})
}

func TestProvideRecipeTagsServiceRecipeIDFetcher(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		_ = ProvideRecipeTagsServiceRecipeIDFetcher(noop.ProvideNoopLogger())
	})
}

func TestProvideRecipeTagsServiceRecipeTagIDFetcher(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		_ = ProvideRecipeTagsServiceRecipeTagIDFetcher(noop.ProvideNoopLogger())
	})
}

func TestProvideRecipeStepsServiceUserIDFetcher(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		_ = ProvideRecipeStepsServiceUserIDFetcher()
	})
}

func TestProvideRecipeStepsServiceRecipeIDFetcher(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		_ = ProvideRecipeStepsServiceRecipeIDFetcher(noop.ProvideNoopLogger())
	})
}

func TestProvideRecipeStepsServiceRecipeStepIDFetcher(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		_ = ProvideRecipeStepsServiceRecipeStepIDFetcher(noop.ProvideNoopLogger())
	})
}

func TestProvideRecipeStepPreparationsServiceUserIDFetcher(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		_ = ProvideRecipeStepPreparationsServiceUserIDFetcher()
	})
}

func TestProvideRecipeStepPreparationsServiceRecipeIDFetcher(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		_ = ProvideRecipeStepPreparationsServiceRecipeIDFetcher(noop.ProvideNoopLogger())
	})
}

func TestProvideRecipeStepPreparationsServiceRecipeStepIDFetcher(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		_ = ProvideRecipeStepPreparationsServiceRecipeStepIDFetcher(noop.ProvideNoopLogger())
	})
}

func TestProvideRecipeStepPreparationsServiceRecipeStepPreparationIDFetcher(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		_ = ProvideRecipeStepPreparationsServiceRecipeStepPreparationIDFetcher(noop.ProvideNoopLogger())
	})
}

func TestProvideRecipeStepIngredientsServiceUserIDFetcher(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		_ = ProvideRecipeStepIngredientsServiceUserIDFetcher()
	})
}

func TestProvideRecipeStepIngredientsServiceRecipeIDFetcher(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		_ = ProvideRecipeStepIngredientsServiceRecipeIDFetcher(noop.ProvideNoopLogger())
	})
}

func TestProvideRecipeStepIngredientsServiceRecipeStepIDFetcher(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		_ = ProvideRecipeStepIngredientsServiceRecipeStepIDFetcher(noop.ProvideNoopLogger())
	})
}

func TestProvideRecipeStepIngredientsServiceRecipeStepIngredientIDFetcher(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		_ = ProvideRecipeStepIngredientsServiceRecipeStepIngredientIDFetcher(noop.ProvideNoopLogger())
	})
}

func TestProvideRecipeIterationsServiceUserIDFetcher(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		_ = ProvideRecipeIterationsServiceUserIDFetcher()
	})
}

func TestProvideRecipeIterationsServiceRecipeIDFetcher(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		_ = ProvideRecipeIterationsServiceRecipeIDFetcher(noop.ProvideNoopLogger())
	})
}

func TestProvideRecipeIterationsServiceRecipeIterationIDFetcher(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		_ = ProvideRecipeIterationsServiceRecipeIterationIDFetcher(noop.ProvideNoopLogger())
	})
}

func TestProvideRecipeIterationStepsServiceUserIDFetcher(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		_ = ProvideRecipeIterationStepsServiceUserIDFetcher()
	})
}

func TestProvideRecipeIterationStepsServiceRecipeIDFetcher(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		_ = ProvideRecipeIterationStepsServiceRecipeIDFetcher(noop.ProvideNoopLogger())
	})
}

func TestProvideRecipeIterationStepsServiceRecipeIterationStepIDFetcher(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		_ = ProvideRecipeIterationStepsServiceRecipeIterationStepIDFetcher(noop.ProvideNoopLogger())
	})
}

func TestProvideIterationMediasServiceUserIDFetcher(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		_ = ProvideIterationMediasServiceUserIDFetcher()
	})
}

func TestProvideIterationMediasServiceRecipeIDFetcher(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		_ = ProvideIterationMediasServiceRecipeIDFetcher(noop.ProvideNoopLogger())
	})
}

func TestProvideIterationMediasServiceRecipeIterationIDFetcher(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		_ = ProvideIterationMediasServiceRecipeIterationIDFetcher(noop.ProvideNoopLogger())
	})
}

func TestProvideIterationMediasServiceIterationMediaIDFetcher(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		_ = ProvideIterationMediasServiceIterationMediaIDFetcher(noop.ProvideNoopLogger())
	})
}

func TestProvideInvitationsServiceUserIDFetcher(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		_ = ProvideInvitationsServiceUserIDFetcher()
	})
}

func TestProvideInvitationsServiceInvitationIDFetcher(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		_ = ProvideInvitationsServiceInvitationIDFetcher(noop.ProvideNoopLogger())
	})
}

func TestProvideReportsServiceUserIDFetcher(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		_ = ProvideReportsServiceUserIDFetcher()
	})
}

func TestProvideReportsServiceReportIDFetcher(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		_ = ProvideReportsServiceReportIDFetcher(noop.ProvideNoopLogger())
	})
}

func TestProvideUsersServiceUserIDFetcher(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		_ = ProvideUsersServiceUserIDFetcher(noop.ProvideNoopLogger())
	})
}

func TestProvideAuthServiceUserIDFetcher(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		_ = ProvideAuthServiceUserIDFetcher()
	})
}

func TestProvideWebhooksServiceUserIDFetcher(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		_ = ProvideWebhooksServiceUserIDFetcher()
	})
}

func TestProvideWebhooksServiceWebhookIDFetcher(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		_ = ProvideWebhooksServiceWebhookIDFetcher(noop.ProvideNoopLogger())
	})
}

func TestProvideOAuth2ClientsServiceClientIDFetcher(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		_ = ProvideOAuth2ClientsServiceClientIDFetcher(noop.ProvideNoopLogger())
	})
}

func Test_userIDFetcherFromRequestContext(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		expected := uint64(123)

		req := buildRequest(t)
		req = req.WithContext(
			context.WithValue(
				req.Context(),
				models.UserIDKey,
				expected,
			),
		)

		actual := userIDFetcherFromRequestContext(req)
		assert.Equal(t, expected, actual)
	})

	T.Run("without attached value", func(t *testing.T) {
		req := buildRequest(t)
		actual := userIDFetcherFromRequestContext(req)

		assert.Zero(t, actual)
	})
}

func Test_buildRouteParamUserIDFetcher(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		fn := buildRouteParamUserIDFetcher(noop.ProvideNoopLogger())
		expected := uint64(123)

		req := buildRequest(t)
		req = req.WithContext(
			context.WithValue(
				req.Context(),
				chi.RouteCtxKey,
				&chi.Context{
					URLParams: chi.RouteParams{
						Keys:   []string{usersservice.URIParamKey},
						Values: []string{fmt.Sprintf("%d", expected)},
					},
				},
			),
		)

		actual := fn(req)
		assert.Equal(t, expected, actual)
	})

	T.Run("with invalid value somehow", func(t *testing.T) {
		// NOTE: This will probably never happen in dev or production
		fn := buildRouteParamUserIDFetcher(noop.ProvideNoopLogger())
		expected := uint64(0)

		req := buildRequest(t)
		req = req.WithContext(
			context.WithValue(
				req.Context(),
				chi.RouteCtxKey,
				&chi.Context{
					URLParams: chi.RouteParams{
						Keys:   []string{usersservice.URIParamKey},
						Values: []string{"expected"},
					},
				},
			),
		)

		actual := fn(req)
		assert.Equal(t, expected, actual)
	})
}

func Test_buildRouteParamValidInstrumentIDFetcher(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		fn := buildRouteParamValidInstrumentIDFetcher(noop.ProvideNoopLogger())
		expected := uint64(123)

		req := buildRequest(t)
		req = req.WithContext(
			context.WithValue(
				req.Context(),
				chi.RouteCtxKey,
				&chi.Context{
					URLParams: chi.RouteParams{
						Keys:   []string{validinstrumentsservice.URIParamKey},
						Values: []string{fmt.Sprintf("%d", expected)},
					},
				},
			),
		)

		actual := fn(req)
		assert.Equal(t, expected, actual)
	})

	T.Run("with invalid value somehow", func(t *testing.T) {
		// NOTE: This will probably never happen in dev or production
		fn := buildRouteParamValidInstrumentIDFetcher(noop.ProvideNoopLogger())
		expected := uint64(0)

		req := buildRequest(t)
		req = req.WithContext(
			context.WithValue(
				req.Context(),
				chi.RouteCtxKey,
				&chi.Context{
					URLParams: chi.RouteParams{
						Keys:   []string{validinstrumentsservice.URIParamKey},
						Values: []string{"expected"},
					},
				},
			),
		)

		actual := fn(req)
		assert.Equal(t, expected, actual)
	})
}

func Test_buildRouteParamValidIngredientIDFetcher(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		fn := buildRouteParamValidIngredientIDFetcher(noop.ProvideNoopLogger())
		expected := uint64(123)

		req := buildRequest(t)
		req = req.WithContext(
			context.WithValue(
				req.Context(),
				chi.RouteCtxKey,
				&chi.Context{
					URLParams: chi.RouteParams{
						Keys:   []string{validingredientsservice.URIParamKey},
						Values: []string{fmt.Sprintf("%d", expected)},
					},
				},
			),
		)

		actual := fn(req)
		assert.Equal(t, expected, actual)
	})

	T.Run("with invalid value somehow", func(t *testing.T) {
		// NOTE: This will probably never happen in dev or production
		fn := buildRouteParamValidIngredientIDFetcher(noop.ProvideNoopLogger())
		expected := uint64(0)

		req := buildRequest(t)
		req = req.WithContext(
			context.WithValue(
				req.Context(),
				chi.RouteCtxKey,
				&chi.Context{
					URLParams: chi.RouteParams{
						Keys:   []string{validingredientsservice.URIParamKey},
						Values: []string{"expected"},
					},
				},
			),
		)

		actual := fn(req)
		assert.Equal(t, expected, actual)
	})
}

func Test_buildRouteParamValidIngredientTagIDFetcher(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		fn := buildRouteParamValidIngredientTagIDFetcher(noop.ProvideNoopLogger())
		expected := uint64(123)

		req := buildRequest(t)
		req = req.WithContext(
			context.WithValue(
				req.Context(),
				chi.RouteCtxKey,
				&chi.Context{
					URLParams: chi.RouteParams{
						Keys:   []string{validingredienttagsservice.URIParamKey},
						Values: []string{fmt.Sprintf("%d", expected)},
					},
				},
			),
		)

		actual := fn(req)
		assert.Equal(t, expected, actual)
	})

	T.Run("with invalid value somehow", func(t *testing.T) {
		// NOTE: This will probably never happen in dev or production
		fn := buildRouteParamValidIngredientTagIDFetcher(noop.ProvideNoopLogger())
		expected := uint64(0)

		req := buildRequest(t)
		req = req.WithContext(
			context.WithValue(
				req.Context(),
				chi.RouteCtxKey,
				&chi.Context{
					URLParams: chi.RouteParams{
						Keys:   []string{validingredienttagsservice.URIParamKey},
						Values: []string{"expected"},
					},
				},
			),
		)

		actual := fn(req)
		assert.Equal(t, expected, actual)
	})
}

func Test_buildRouteParamIngredientTagMappingIDFetcher(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		fn := buildRouteParamIngredientTagMappingIDFetcher(noop.ProvideNoopLogger())
		expected := uint64(123)

		req := buildRequest(t)
		req = req.WithContext(
			context.WithValue(
				req.Context(),
				chi.RouteCtxKey,
				&chi.Context{
					URLParams: chi.RouteParams{
						Keys:   []string{ingredienttagmappingsservice.URIParamKey},
						Values: []string{fmt.Sprintf("%d", expected)},
					},
				},
			),
		)

		actual := fn(req)
		assert.Equal(t, expected, actual)
	})

	T.Run("with invalid value somehow", func(t *testing.T) {
		// NOTE: This will probably never happen in dev or production
		fn := buildRouteParamIngredientTagMappingIDFetcher(noop.ProvideNoopLogger())
		expected := uint64(0)

		req := buildRequest(t)
		req = req.WithContext(
			context.WithValue(
				req.Context(),
				chi.RouteCtxKey,
				&chi.Context{
					URLParams: chi.RouteParams{
						Keys:   []string{ingredienttagmappingsservice.URIParamKey},
						Values: []string{"expected"},
					},
				},
			),
		)

		actual := fn(req)
		assert.Equal(t, expected, actual)
	})
}

func Test_buildRouteParamValidPreparationIDFetcher(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		fn := buildRouteParamValidPreparationIDFetcher(noop.ProvideNoopLogger())
		expected := uint64(123)

		req := buildRequest(t)
		req = req.WithContext(
			context.WithValue(
				req.Context(),
				chi.RouteCtxKey,
				&chi.Context{
					URLParams: chi.RouteParams{
						Keys:   []string{validpreparationsservice.URIParamKey},
						Values: []string{fmt.Sprintf("%d", expected)},
					},
				},
			),
		)

		actual := fn(req)
		assert.Equal(t, expected, actual)
	})

	T.Run("with invalid value somehow", func(t *testing.T) {
		// NOTE: This will probably never happen in dev or production
		fn := buildRouteParamValidPreparationIDFetcher(noop.ProvideNoopLogger())
		expected := uint64(0)

		req := buildRequest(t)
		req = req.WithContext(
			context.WithValue(
				req.Context(),
				chi.RouteCtxKey,
				&chi.Context{
					URLParams: chi.RouteParams{
						Keys:   []string{validpreparationsservice.URIParamKey},
						Values: []string{"expected"},
					},
				},
			),
		)

		actual := fn(req)
		assert.Equal(t, expected, actual)
	})
}

func Test_buildRouteParamRequiredPreparationInstrumentIDFetcher(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		fn := buildRouteParamRequiredPreparationInstrumentIDFetcher(noop.ProvideNoopLogger())
		expected := uint64(123)

		req := buildRequest(t)
		req = req.WithContext(
			context.WithValue(
				req.Context(),
				chi.RouteCtxKey,
				&chi.Context{
					URLParams: chi.RouteParams{
						Keys:   []string{requiredpreparationinstrumentsservice.URIParamKey},
						Values: []string{fmt.Sprintf("%d", expected)},
					},
				},
			),
		)

		actual := fn(req)
		assert.Equal(t, expected, actual)
	})

	T.Run("with invalid value somehow", func(t *testing.T) {
		// NOTE: This will probably never happen in dev or production
		fn := buildRouteParamRequiredPreparationInstrumentIDFetcher(noop.ProvideNoopLogger())
		expected := uint64(0)

		req := buildRequest(t)
		req = req.WithContext(
			context.WithValue(
				req.Context(),
				chi.RouteCtxKey,
				&chi.Context{
					URLParams: chi.RouteParams{
						Keys:   []string{requiredpreparationinstrumentsservice.URIParamKey},
						Values: []string{"expected"},
					},
				},
			),
		)

		actual := fn(req)
		assert.Equal(t, expected, actual)
	})
}

func Test_buildRouteParamValidIngredientPreparationIDFetcher(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		fn := buildRouteParamValidIngredientPreparationIDFetcher(noop.ProvideNoopLogger())
		expected := uint64(123)

		req := buildRequest(t)
		req = req.WithContext(
			context.WithValue(
				req.Context(),
				chi.RouteCtxKey,
				&chi.Context{
					URLParams: chi.RouteParams{
						Keys:   []string{validingredientpreparationsservice.URIParamKey},
						Values: []string{fmt.Sprintf("%d", expected)},
					},
				},
			),
		)

		actual := fn(req)
		assert.Equal(t, expected, actual)
	})

	T.Run("with invalid value somehow", func(t *testing.T) {
		// NOTE: This will probably never happen in dev or production
		fn := buildRouteParamValidIngredientPreparationIDFetcher(noop.ProvideNoopLogger())
		expected := uint64(0)

		req := buildRequest(t)
		req = req.WithContext(
			context.WithValue(
				req.Context(),
				chi.RouteCtxKey,
				&chi.Context{
					URLParams: chi.RouteParams{
						Keys:   []string{validingredientpreparationsservice.URIParamKey},
						Values: []string{"expected"},
					},
				},
			),
		)

		actual := fn(req)
		assert.Equal(t, expected, actual)
	})
}

func Test_buildRouteParamRecipeIDFetcher(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		fn := buildRouteParamRecipeIDFetcher(noop.ProvideNoopLogger())
		expected := uint64(123)

		req := buildRequest(t)
		req = req.WithContext(
			context.WithValue(
				req.Context(),
				chi.RouteCtxKey,
				&chi.Context{
					URLParams: chi.RouteParams{
						Keys:   []string{recipesservice.URIParamKey},
						Values: []string{fmt.Sprintf("%d", expected)},
					},
				},
			),
		)

		actual := fn(req)
		assert.Equal(t, expected, actual)
	})

	T.Run("with invalid value somehow", func(t *testing.T) {
		// NOTE: This will probably never happen in dev or production
		fn := buildRouteParamRecipeIDFetcher(noop.ProvideNoopLogger())
		expected := uint64(0)

		req := buildRequest(t)
		req = req.WithContext(
			context.WithValue(
				req.Context(),
				chi.RouteCtxKey,
				&chi.Context{
					URLParams: chi.RouteParams{
						Keys:   []string{recipesservice.URIParamKey},
						Values: []string{"expected"},
					},
				},
			),
		)

		actual := fn(req)
		assert.Equal(t, expected, actual)
	})
}

func Test_buildRouteParamRecipeTagIDFetcher(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		fn := buildRouteParamRecipeTagIDFetcher(noop.ProvideNoopLogger())
		expected := uint64(123)

		req := buildRequest(t)
		req = req.WithContext(
			context.WithValue(
				req.Context(),
				chi.RouteCtxKey,
				&chi.Context{
					URLParams: chi.RouteParams{
						Keys:   []string{recipetagsservice.URIParamKey},
						Values: []string{fmt.Sprintf("%d", expected)},
					},
				},
			),
		)

		actual := fn(req)
		assert.Equal(t, expected, actual)
	})

	T.Run("with invalid value somehow", func(t *testing.T) {
		// NOTE: This will probably never happen in dev or production
		fn := buildRouteParamRecipeTagIDFetcher(noop.ProvideNoopLogger())
		expected := uint64(0)

		req := buildRequest(t)
		req = req.WithContext(
			context.WithValue(
				req.Context(),
				chi.RouteCtxKey,
				&chi.Context{
					URLParams: chi.RouteParams{
						Keys:   []string{recipetagsservice.URIParamKey},
						Values: []string{"expected"},
					},
				},
			),
		)

		actual := fn(req)
		assert.Equal(t, expected, actual)
	})
}

func Test_buildRouteParamRecipeStepIDFetcher(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		fn := buildRouteParamRecipeStepIDFetcher(noop.ProvideNoopLogger())
		expected := uint64(123)

		req := buildRequest(t)
		req = req.WithContext(
			context.WithValue(
				req.Context(),
				chi.RouteCtxKey,
				&chi.Context{
					URLParams: chi.RouteParams{
						Keys:   []string{recipestepsservice.URIParamKey},
						Values: []string{fmt.Sprintf("%d", expected)},
					},
				},
			),
		)

		actual := fn(req)
		assert.Equal(t, expected, actual)
	})

	T.Run("with invalid value somehow", func(t *testing.T) {
		// NOTE: This will probably never happen in dev or production
		fn := buildRouteParamRecipeStepIDFetcher(noop.ProvideNoopLogger())
		expected := uint64(0)

		req := buildRequest(t)
		req = req.WithContext(
			context.WithValue(
				req.Context(),
				chi.RouteCtxKey,
				&chi.Context{
					URLParams: chi.RouteParams{
						Keys:   []string{recipestepsservice.URIParamKey},
						Values: []string{"expected"},
					},
				},
			),
		)

		actual := fn(req)
		assert.Equal(t, expected, actual)
	})
}

func Test_buildRouteParamRecipeStepPreparationIDFetcher(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		fn := buildRouteParamRecipeStepPreparationIDFetcher(noop.ProvideNoopLogger())
		expected := uint64(123)

		req := buildRequest(t)
		req = req.WithContext(
			context.WithValue(
				req.Context(),
				chi.RouteCtxKey,
				&chi.Context{
					URLParams: chi.RouteParams{
						Keys:   []string{recipesteppreparationsservice.URIParamKey},
						Values: []string{fmt.Sprintf("%d", expected)},
					},
				},
			),
		)

		actual := fn(req)
		assert.Equal(t, expected, actual)
	})

	T.Run("with invalid value somehow", func(t *testing.T) {
		// NOTE: This will probably never happen in dev or production
		fn := buildRouteParamRecipeStepPreparationIDFetcher(noop.ProvideNoopLogger())
		expected := uint64(0)

		req := buildRequest(t)
		req = req.WithContext(
			context.WithValue(
				req.Context(),
				chi.RouteCtxKey,
				&chi.Context{
					URLParams: chi.RouteParams{
						Keys:   []string{recipesteppreparationsservice.URIParamKey},
						Values: []string{"expected"},
					},
				},
			),
		)

		actual := fn(req)
		assert.Equal(t, expected, actual)
	})
}

func Test_buildRouteParamRecipeStepIngredientIDFetcher(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		fn := buildRouteParamRecipeStepIngredientIDFetcher(noop.ProvideNoopLogger())
		expected := uint64(123)

		req := buildRequest(t)
		req = req.WithContext(
			context.WithValue(
				req.Context(),
				chi.RouteCtxKey,
				&chi.Context{
					URLParams: chi.RouteParams{
						Keys:   []string{recipestepingredientsservice.URIParamKey},
						Values: []string{fmt.Sprintf("%d", expected)},
					},
				},
			),
		)

		actual := fn(req)
		assert.Equal(t, expected, actual)
	})

	T.Run("with invalid value somehow", func(t *testing.T) {
		// NOTE: This will probably never happen in dev or production
		fn := buildRouteParamRecipeStepIngredientIDFetcher(noop.ProvideNoopLogger())
		expected := uint64(0)

		req := buildRequest(t)
		req = req.WithContext(
			context.WithValue(
				req.Context(),
				chi.RouteCtxKey,
				&chi.Context{
					URLParams: chi.RouteParams{
						Keys:   []string{recipestepingredientsservice.URIParamKey},
						Values: []string{"expected"},
					},
				},
			),
		)

		actual := fn(req)
		assert.Equal(t, expected, actual)
	})
}

func Test_buildRouteParamRecipeIterationIDFetcher(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		fn := buildRouteParamRecipeIterationIDFetcher(noop.ProvideNoopLogger())
		expected := uint64(123)

		req := buildRequest(t)
		req = req.WithContext(
			context.WithValue(
				req.Context(),
				chi.RouteCtxKey,
				&chi.Context{
					URLParams: chi.RouteParams{
						Keys:   []string{recipeiterationsservice.URIParamKey},
						Values: []string{fmt.Sprintf("%d", expected)},
					},
				},
			),
		)

		actual := fn(req)
		assert.Equal(t, expected, actual)
	})

	T.Run("with invalid value somehow", func(t *testing.T) {
		// NOTE: This will probably never happen in dev or production
		fn := buildRouteParamRecipeIterationIDFetcher(noop.ProvideNoopLogger())
		expected := uint64(0)

		req := buildRequest(t)
		req = req.WithContext(
			context.WithValue(
				req.Context(),
				chi.RouteCtxKey,
				&chi.Context{
					URLParams: chi.RouteParams{
						Keys:   []string{recipeiterationsservice.URIParamKey},
						Values: []string{"expected"},
					},
				},
			),
		)

		actual := fn(req)
		assert.Equal(t, expected, actual)
	})
}

func Test_buildRouteParamRecipeIterationStepIDFetcher(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		fn := buildRouteParamRecipeIterationStepIDFetcher(noop.ProvideNoopLogger())
		expected := uint64(123)

		req := buildRequest(t)
		req = req.WithContext(
			context.WithValue(
				req.Context(),
				chi.RouteCtxKey,
				&chi.Context{
					URLParams: chi.RouteParams{
						Keys:   []string{recipeiterationstepsservice.URIParamKey},
						Values: []string{fmt.Sprintf("%d", expected)},
					},
				},
			),
		)

		actual := fn(req)
		assert.Equal(t, expected, actual)
	})

	T.Run("with invalid value somehow", func(t *testing.T) {
		// NOTE: This will probably never happen in dev or production
		fn := buildRouteParamRecipeIterationStepIDFetcher(noop.ProvideNoopLogger())
		expected := uint64(0)

		req := buildRequest(t)
		req = req.WithContext(
			context.WithValue(
				req.Context(),
				chi.RouteCtxKey,
				&chi.Context{
					URLParams: chi.RouteParams{
						Keys:   []string{recipeiterationstepsservice.URIParamKey},
						Values: []string{"expected"},
					},
				},
			),
		)

		actual := fn(req)
		assert.Equal(t, expected, actual)
	})
}

func Test_buildRouteParamIterationMediaIDFetcher(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		fn := buildRouteParamIterationMediaIDFetcher(noop.ProvideNoopLogger())
		expected := uint64(123)

		req := buildRequest(t)
		req = req.WithContext(
			context.WithValue(
				req.Context(),
				chi.RouteCtxKey,
				&chi.Context{
					URLParams: chi.RouteParams{
						Keys:   []string{iterationmediasservice.URIParamKey},
						Values: []string{fmt.Sprintf("%d", expected)},
					},
				},
			),
		)

		actual := fn(req)
		assert.Equal(t, expected, actual)
	})

	T.Run("with invalid value somehow", func(t *testing.T) {
		// NOTE: This will probably never happen in dev or production
		fn := buildRouteParamIterationMediaIDFetcher(noop.ProvideNoopLogger())
		expected := uint64(0)

		req := buildRequest(t)
		req = req.WithContext(
			context.WithValue(
				req.Context(),
				chi.RouteCtxKey,
				&chi.Context{
					URLParams: chi.RouteParams{
						Keys:   []string{iterationmediasservice.URIParamKey},
						Values: []string{"expected"},
					},
				},
			),
		)

		actual := fn(req)
		assert.Equal(t, expected, actual)
	})
}

func Test_buildRouteParamInvitationIDFetcher(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		fn := buildRouteParamInvitationIDFetcher(noop.ProvideNoopLogger())
		expected := uint64(123)

		req := buildRequest(t)
		req = req.WithContext(
			context.WithValue(
				req.Context(),
				chi.RouteCtxKey,
				&chi.Context{
					URLParams: chi.RouteParams{
						Keys:   []string{invitationsservice.URIParamKey},
						Values: []string{fmt.Sprintf("%d", expected)},
					},
				},
			),
		)

		actual := fn(req)
		assert.Equal(t, expected, actual)
	})

	T.Run("with invalid value somehow", func(t *testing.T) {
		// NOTE: This will probably never happen in dev or production
		fn := buildRouteParamInvitationIDFetcher(noop.ProvideNoopLogger())
		expected := uint64(0)

		req := buildRequest(t)
		req = req.WithContext(
			context.WithValue(
				req.Context(),
				chi.RouteCtxKey,
				&chi.Context{
					URLParams: chi.RouteParams{
						Keys:   []string{invitationsservice.URIParamKey},
						Values: []string{"expected"},
					},
				},
			),
		)

		actual := fn(req)
		assert.Equal(t, expected, actual)
	})
}

func Test_buildRouteParamReportIDFetcher(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		fn := buildRouteParamReportIDFetcher(noop.ProvideNoopLogger())
		expected := uint64(123)

		req := buildRequest(t)
		req = req.WithContext(
			context.WithValue(
				req.Context(),
				chi.RouteCtxKey,
				&chi.Context{
					URLParams: chi.RouteParams{
						Keys:   []string{reportsservice.URIParamKey},
						Values: []string{fmt.Sprintf("%d", expected)},
					},
				},
			),
		)

		actual := fn(req)
		assert.Equal(t, expected, actual)
	})

	T.Run("with invalid value somehow", func(t *testing.T) {
		// NOTE: This will probably never happen in dev or production
		fn := buildRouteParamReportIDFetcher(noop.ProvideNoopLogger())
		expected := uint64(0)

		req := buildRequest(t)
		req = req.WithContext(
			context.WithValue(
				req.Context(),
				chi.RouteCtxKey,
				&chi.Context{
					URLParams: chi.RouteParams{
						Keys:   []string{reportsservice.URIParamKey},
						Values: []string{"expected"},
					},
				},
			),
		)

		actual := fn(req)
		assert.Equal(t, expected, actual)
	})
}

func Test_buildRouteParamWebhookIDFetcher(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		fn := buildRouteParamWebhookIDFetcher(noop.ProvideNoopLogger())
		expected := uint64(123)

		req := buildRequest(t)
		req = req.WithContext(
			context.WithValue(
				req.Context(),
				chi.RouteCtxKey,
				&chi.Context{
					URLParams: chi.RouteParams{
						Keys:   []string{webhooksservice.URIParamKey},
						Values: []string{fmt.Sprintf("%d", expected)},
					},
				},
			),
		)

		actual := fn(req)
		assert.Equal(t, expected, actual)
	})

	T.Run("with invalid value somehow", func(t *testing.T) {
		// NOTE: This will probably never happen in dev or production
		fn := buildRouteParamWebhookIDFetcher(noop.ProvideNoopLogger())
		expected := uint64(0)

		req := buildRequest(t)
		req = req.WithContext(
			context.WithValue(
				req.Context(),
				chi.RouteCtxKey,
				&chi.Context{
					URLParams: chi.RouteParams{
						Keys:   []string{webhooksservice.URIParamKey},
						Values: []string{"expected"},
					},
				},
			),
		)

		actual := fn(req)
		assert.Equal(t, expected, actual)
	})
}

func Test_buildRouteParamOAuth2ClientIDFetcher(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		fn := buildRouteParamOAuth2ClientIDFetcher(noop.ProvideNoopLogger())
		expected := uint64(123)

		req := buildRequest(t)
		req = req.WithContext(
			context.WithValue(
				req.Context(),
				chi.RouteCtxKey,
				&chi.Context{
					URLParams: chi.RouteParams{
						Keys:   []string{oauth2clientsservice.URIParamKey},
						Values: []string{fmt.Sprintf("%d", expected)},
					},
				},
			),
		)

		actual := fn(req)
		assert.Equal(t, expected, actual)
	})

	T.Run("with invalid value somehow", func(t *testing.T) {
		// NOTE: This will probably never happen in dev or production
		fn := buildRouteParamOAuth2ClientIDFetcher(noop.ProvideNoopLogger())
		expected := uint64(0)

		req := buildRequest(t)
		req = req.WithContext(
			context.WithValue(
				req.Context(),
				chi.RouteCtxKey,
				&chi.Context{
					URLParams: chi.RouteParams{
						Keys:   []string{oauth2clientsservice.URIParamKey},
						Values: []string{"expected"},
					},
				},
			),
		)

		actual := fn(req)
		assert.Equal(t, expected, actual)
	})
}
