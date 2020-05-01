package httpserver

import (
	"net/http"
	"strconv"

	models "gitlab.com/prixfixe/prixfixe/models/v1"
	authservice "gitlab.com/prixfixe/prixfixe/services/v1/auth"
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
	"github.com/google/wire"
	"gitlab.com/verygoodsoftwarenotvirus/logging/v1"
)

var (
	paramFetcherProviders = wire.NewSet(
		ProvideUsersServiceUserIDFetcher,
		ProvideOAuth2ClientsServiceClientIDFetcher,
		ProvideAuthServiceUserIDFetcher,
		ProvideValidInstrumentsServiceValidInstrumentIDFetcher,
		ProvideValidIngredientsServiceValidIngredientIDFetcher,
		ProvideValidIngredientTagsServiceValidIngredientTagIDFetcher,
		ProvideIngredientTagMappingsServiceValidIngredientIDFetcher,
		ProvideIngredientTagMappingsServiceIngredientTagMappingIDFetcher,
		ProvideValidPreparationsServiceValidPreparationIDFetcher,
		ProvideRequiredPreparationInstrumentsServiceValidPreparationIDFetcher,
		ProvideRequiredPreparationInstrumentsServiceRequiredPreparationInstrumentIDFetcher,
		ProvideValidIngredientPreparationsServiceValidIngredientIDFetcher,
		ProvideValidIngredientPreparationsServiceValidIngredientPreparationIDFetcher,
		ProvideRecipesServiceUserIDFetcher,
		ProvideRecipesServiceRecipeIDFetcher,
		ProvideRecipeTagsServiceUserIDFetcher,
		ProvideRecipeTagsServiceRecipeIDFetcher,
		ProvideRecipeTagsServiceRecipeTagIDFetcher,
		ProvideRecipeStepsServiceUserIDFetcher,
		ProvideRecipeStepsServiceRecipeIDFetcher,
		ProvideRecipeStepsServiceRecipeStepIDFetcher,
		ProvideRecipeStepPreparationsServiceUserIDFetcher,
		ProvideRecipeStepPreparationsServiceRecipeIDFetcher,
		ProvideRecipeStepPreparationsServiceRecipeStepIDFetcher,
		ProvideRecipeStepPreparationsServiceRecipeStepPreparationIDFetcher,
		ProvideRecipeStepIngredientsServiceUserIDFetcher,
		ProvideRecipeStepIngredientsServiceRecipeIDFetcher,
		ProvideRecipeStepIngredientsServiceRecipeStepIDFetcher,
		ProvideRecipeStepIngredientsServiceRecipeStepIngredientIDFetcher,
		ProvideRecipeIterationsServiceUserIDFetcher,
		ProvideRecipeIterationsServiceRecipeIDFetcher,
		ProvideRecipeIterationsServiceRecipeIterationIDFetcher,
		ProvideRecipeIterationStepsServiceUserIDFetcher,
		ProvideRecipeIterationStepsServiceRecipeIDFetcher,
		ProvideRecipeIterationStepsServiceRecipeIterationStepIDFetcher,
		ProvideIterationMediasServiceUserIDFetcher,
		ProvideIterationMediasServiceRecipeIDFetcher,
		ProvideIterationMediasServiceRecipeIterationIDFetcher,
		ProvideIterationMediasServiceIterationMediaIDFetcher,
		ProvideInvitationsServiceUserIDFetcher,
		ProvideInvitationsServiceInvitationIDFetcher,
		ProvideReportsServiceUserIDFetcher,
		ProvideReportsServiceReportIDFetcher,
		ProvideWebhooksServiceUserIDFetcher,
		ProvideWebhooksServiceWebhookIDFetcher,
	)
)

// ProvideValidInstrumentsServiceValidInstrumentIDFetcher provides a ValidInstrumentIDFetcher.
func ProvideValidInstrumentsServiceValidInstrumentIDFetcher(logger logging.Logger) validinstrumentsservice.ValidInstrumentIDFetcher {
	return buildRouteParamValidInstrumentIDFetcher(logger)
}

// ProvideValidIngredientsServiceValidIngredientIDFetcher provides a ValidIngredientIDFetcher.
func ProvideValidIngredientsServiceValidIngredientIDFetcher(logger logging.Logger) validingredientsservice.ValidIngredientIDFetcher {
	return buildRouteParamValidIngredientIDFetcher(logger)
}

// ProvideValidIngredientTagsServiceValidIngredientTagIDFetcher provides a ValidIngredientTagIDFetcher.
func ProvideValidIngredientTagsServiceValidIngredientTagIDFetcher(logger logging.Logger) validingredienttagsservice.ValidIngredientTagIDFetcher {
	return buildRouteParamValidIngredientTagIDFetcher(logger)
}

// ProvideIngredientTagMappingsServiceValidIngredientIDFetcher provides a ValidIngredientIDFetcher.
func ProvideIngredientTagMappingsServiceValidIngredientIDFetcher(logger logging.Logger) ingredienttagmappingsservice.ValidIngredientIDFetcher {
	return buildRouteParamValidIngredientIDFetcher(logger)
}

// ProvideIngredientTagMappingsServiceIngredientTagMappingIDFetcher provides an IngredientTagMappingIDFetcher.
func ProvideIngredientTagMappingsServiceIngredientTagMappingIDFetcher(logger logging.Logger) ingredienttagmappingsservice.IngredientTagMappingIDFetcher {
	return buildRouteParamIngredientTagMappingIDFetcher(logger)
}

// ProvideValidPreparationsServiceValidPreparationIDFetcher provides a ValidPreparationIDFetcher.
func ProvideValidPreparationsServiceValidPreparationIDFetcher(logger logging.Logger) validpreparationsservice.ValidPreparationIDFetcher {
	return buildRouteParamValidPreparationIDFetcher(logger)
}

// ProvideRequiredPreparationInstrumentsServiceValidPreparationIDFetcher provides a ValidPreparationIDFetcher.
func ProvideRequiredPreparationInstrumentsServiceValidPreparationIDFetcher(logger logging.Logger) requiredpreparationinstrumentsservice.ValidPreparationIDFetcher {
	return buildRouteParamValidPreparationIDFetcher(logger)
}

// ProvideRequiredPreparationInstrumentsServiceRequiredPreparationInstrumentIDFetcher provides a RequiredPreparationInstrumentIDFetcher.
func ProvideRequiredPreparationInstrumentsServiceRequiredPreparationInstrumentIDFetcher(logger logging.Logger) requiredpreparationinstrumentsservice.RequiredPreparationInstrumentIDFetcher {
	return buildRouteParamRequiredPreparationInstrumentIDFetcher(logger)
}

// ProvideValidIngredientPreparationsServiceValidIngredientIDFetcher provides a ValidIngredientIDFetcher.
func ProvideValidIngredientPreparationsServiceValidIngredientIDFetcher(logger logging.Logger) validingredientpreparationsservice.ValidIngredientIDFetcher {
	return buildRouteParamValidIngredientIDFetcher(logger)
}

// ProvideValidIngredientPreparationsServiceValidIngredientPreparationIDFetcher provides a ValidIngredientPreparationIDFetcher.
func ProvideValidIngredientPreparationsServiceValidIngredientPreparationIDFetcher(logger logging.Logger) validingredientpreparationsservice.ValidIngredientPreparationIDFetcher {
	return buildRouteParamValidIngredientPreparationIDFetcher(logger)
}

// ProvideRecipesServiceUserIDFetcher provides a UserIDFetcher.
func ProvideRecipesServiceUserIDFetcher() recipesservice.UserIDFetcher {
	return userIDFetcherFromRequestContext
}

// ProvideRecipesServiceRecipeIDFetcher provides a RecipeIDFetcher.
func ProvideRecipesServiceRecipeIDFetcher(logger logging.Logger) recipesservice.RecipeIDFetcher {
	return buildRouteParamRecipeIDFetcher(logger)
}

// ProvideRecipeTagsServiceUserIDFetcher provides a UserIDFetcher.
func ProvideRecipeTagsServiceUserIDFetcher() recipetagsservice.UserIDFetcher {
	return userIDFetcherFromRequestContext
}

// ProvideRecipeTagsServiceRecipeIDFetcher provides a RecipeIDFetcher.
func ProvideRecipeTagsServiceRecipeIDFetcher(logger logging.Logger) recipetagsservice.RecipeIDFetcher {
	return buildRouteParamRecipeIDFetcher(logger)
}

// ProvideRecipeTagsServiceRecipeTagIDFetcher provides a RecipeTagIDFetcher.
func ProvideRecipeTagsServiceRecipeTagIDFetcher(logger logging.Logger) recipetagsservice.RecipeTagIDFetcher {
	return buildRouteParamRecipeTagIDFetcher(logger)
}

// ProvideRecipeStepsServiceUserIDFetcher provides a UserIDFetcher.
func ProvideRecipeStepsServiceUserIDFetcher() recipestepsservice.UserIDFetcher {
	return userIDFetcherFromRequestContext
}

// ProvideRecipeStepsServiceRecipeIDFetcher provides a RecipeIDFetcher.
func ProvideRecipeStepsServiceRecipeIDFetcher(logger logging.Logger) recipestepsservice.RecipeIDFetcher {
	return buildRouteParamRecipeIDFetcher(logger)
}

// ProvideRecipeStepsServiceRecipeStepIDFetcher provides a RecipeStepIDFetcher.
func ProvideRecipeStepsServiceRecipeStepIDFetcher(logger logging.Logger) recipestepsservice.RecipeStepIDFetcher {
	return buildRouteParamRecipeStepIDFetcher(logger)
}

// ProvideRecipeStepPreparationsServiceUserIDFetcher provides a UserIDFetcher.
func ProvideRecipeStepPreparationsServiceUserIDFetcher() recipesteppreparationsservice.UserIDFetcher {
	return userIDFetcherFromRequestContext
}

// ProvideRecipeStepPreparationsServiceRecipeIDFetcher provides a RecipeIDFetcher.
func ProvideRecipeStepPreparationsServiceRecipeIDFetcher(logger logging.Logger) recipesteppreparationsservice.RecipeIDFetcher {
	return buildRouteParamRecipeIDFetcher(logger)
}

// ProvideRecipeStepPreparationsServiceRecipeStepIDFetcher provides a RecipeStepIDFetcher.
func ProvideRecipeStepPreparationsServiceRecipeStepIDFetcher(logger logging.Logger) recipesteppreparationsservice.RecipeStepIDFetcher {
	return buildRouteParamRecipeStepIDFetcher(logger)
}

// ProvideRecipeStepPreparationsServiceRecipeStepPreparationIDFetcher provides a RecipeStepPreparationIDFetcher.
func ProvideRecipeStepPreparationsServiceRecipeStepPreparationIDFetcher(logger logging.Logger) recipesteppreparationsservice.RecipeStepPreparationIDFetcher {
	return buildRouteParamRecipeStepPreparationIDFetcher(logger)
}

// ProvideRecipeStepIngredientsServiceUserIDFetcher provides a UserIDFetcher.
func ProvideRecipeStepIngredientsServiceUserIDFetcher() recipestepingredientsservice.UserIDFetcher {
	return userIDFetcherFromRequestContext
}

// ProvideRecipeStepIngredientsServiceRecipeIDFetcher provides a RecipeIDFetcher.
func ProvideRecipeStepIngredientsServiceRecipeIDFetcher(logger logging.Logger) recipestepingredientsservice.RecipeIDFetcher {
	return buildRouteParamRecipeIDFetcher(logger)
}

// ProvideRecipeStepIngredientsServiceRecipeStepIDFetcher provides a RecipeStepIDFetcher.
func ProvideRecipeStepIngredientsServiceRecipeStepIDFetcher(logger logging.Logger) recipestepingredientsservice.RecipeStepIDFetcher {
	return buildRouteParamRecipeStepIDFetcher(logger)
}

// ProvideRecipeStepIngredientsServiceRecipeStepIngredientIDFetcher provides a RecipeStepIngredientIDFetcher.
func ProvideRecipeStepIngredientsServiceRecipeStepIngredientIDFetcher(logger logging.Logger) recipestepingredientsservice.RecipeStepIngredientIDFetcher {
	return buildRouteParamRecipeStepIngredientIDFetcher(logger)
}

// ProvideRecipeIterationsServiceUserIDFetcher provides a UserIDFetcher.
func ProvideRecipeIterationsServiceUserIDFetcher() recipeiterationsservice.UserIDFetcher {
	return userIDFetcherFromRequestContext
}

// ProvideRecipeIterationsServiceRecipeIDFetcher provides a RecipeIDFetcher.
func ProvideRecipeIterationsServiceRecipeIDFetcher(logger logging.Logger) recipeiterationsservice.RecipeIDFetcher {
	return buildRouteParamRecipeIDFetcher(logger)
}

// ProvideRecipeIterationsServiceRecipeIterationIDFetcher provides a RecipeIterationIDFetcher.
func ProvideRecipeIterationsServiceRecipeIterationIDFetcher(logger logging.Logger) recipeiterationsservice.RecipeIterationIDFetcher {
	return buildRouteParamRecipeIterationIDFetcher(logger)
}

// ProvideRecipeIterationStepsServiceUserIDFetcher provides a UserIDFetcher.
func ProvideRecipeIterationStepsServiceUserIDFetcher() recipeiterationstepsservice.UserIDFetcher {
	return userIDFetcherFromRequestContext
}

// ProvideRecipeIterationStepsServiceRecipeIDFetcher provides a RecipeIDFetcher.
func ProvideRecipeIterationStepsServiceRecipeIDFetcher(logger logging.Logger) recipeiterationstepsservice.RecipeIDFetcher {
	return buildRouteParamRecipeIDFetcher(logger)
}

// ProvideRecipeIterationStepsServiceRecipeIterationStepIDFetcher provides a RecipeIterationStepIDFetcher.
func ProvideRecipeIterationStepsServiceRecipeIterationStepIDFetcher(logger logging.Logger) recipeiterationstepsservice.RecipeIterationStepIDFetcher {
	return buildRouteParamRecipeIterationStepIDFetcher(logger)
}

// ProvideIterationMediasServiceUserIDFetcher provides a UserIDFetcher.
func ProvideIterationMediasServiceUserIDFetcher() iterationmediasservice.UserIDFetcher {
	return userIDFetcherFromRequestContext
}

// ProvideIterationMediasServiceRecipeIDFetcher provides a RecipeIDFetcher.
func ProvideIterationMediasServiceRecipeIDFetcher(logger logging.Logger) iterationmediasservice.RecipeIDFetcher {
	return buildRouteParamRecipeIDFetcher(logger)
}

// ProvideIterationMediasServiceRecipeIterationIDFetcher provides a RecipeIterationIDFetcher.
func ProvideIterationMediasServiceRecipeIterationIDFetcher(logger logging.Logger) iterationmediasservice.RecipeIterationIDFetcher {
	return buildRouteParamRecipeIterationIDFetcher(logger)
}

// ProvideIterationMediasServiceIterationMediaIDFetcher provides an IterationMediaIDFetcher.
func ProvideIterationMediasServiceIterationMediaIDFetcher(logger logging.Logger) iterationmediasservice.IterationMediaIDFetcher {
	return buildRouteParamIterationMediaIDFetcher(logger)
}

// ProvideInvitationsServiceUserIDFetcher provides a UserIDFetcher.
func ProvideInvitationsServiceUserIDFetcher() invitationsservice.UserIDFetcher {
	return userIDFetcherFromRequestContext
}

// ProvideInvitationsServiceInvitationIDFetcher provides an InvitationIDFetcher.
func ProvideInvitationsServiceInvitationIDFetcher(logger logging.Logger) invitationsservice.InvitationIDFetcher {
	return buildRouteParamInvitationIDFetcher(logger)
}

// ProvideReportsServiceUserIDFetcher provides a UserIDFetcher.
func ProvideReportsServiceUserIDFetcher() reportsservice.UserIDFetcher {
	return userIDFetcherFromRequestContext
}

// ProvideReportsServiceReportIDFetcher provides a ReportIDFetcher.
func ProvideReportsServiceReportIDFetcher(logger logging.Logger) reportsservice.ReportIDFetcher {
	return buildRouteParamReportIDFetcher(logger)
}

// ProvideUsersServiceUserIDFetcher provides a UsernameFetcher.
func ProvideUsersServiceUserIDFetcher(logger logging.Logger) usersservice.UserIDFetcher {
	return buildRouteParamUserIDFetcher(logger)
}

// ProvideAuthServiceUserIDFetcher provides a UsernameFetcher.
func ProvideAuthServiceUserIDFetcher() authservice.UserIDFetcher {
	return userIDFetcherFromRequestContext
}

// ProvideWebhooksServiceUserIDFetcher provides a UserIDFetcher.
func ProvideWebhooksServiceUserIDFetcher() webhooksservice.UserIDFetcher {
	return userIDFetcherFromRequestContext
}

// ProvideWebhooksServiceWebhookIDFetcher provides an WebhookIDFetcher.
func ProvideWebhooksServiceWebhookIDFetcher(logger logging.Logger) webhooksservice.WebhookIDFetcher {
	return buildRouteParamWebhookIDFetcher(logger)
}

// ProvideOAuth2ClientsServiceClientIDFetcher provides a ClientIDFetcher.
func ProvideOAuth2ClientsServiceClientIDFetcher(logger logging.Logger) oauth2clientsservice.ClientIDFetcher {
	return buildRouteParamOAuth2ClientIDFetcher(logger)
}

// userIDFetcherFromRequestContext fetches a user ID from a request routed by chi.
func userIDFetcherFromRequestContext(req *http.Request) uint64 {
	if userID, ok := req.Context().Value(models.UserIDKey).(uint64); ok {
		return userID
	}
	return 0
}

// buildRouteParamUserIDFetcher builds a function that fetches a Username from a request routed by chi.
func buildRouteParamUserIDFetcher(logger logging.Logger) usersservice.UserIDFetcher {
	return func(req *http.Request) uint64 {
		u, err := strconv.ParseUint(chi.URLParam(req, usersservice.URIParamKey), 10, 64)
		if err != nil {
			logger.Error(err, "fetching user ID from request")
		}
		return u
	}
}

// buildRouteParamValidInstrumentIDFetcher builds a function that fetches a ValidInstrumentID from a request routed by chi.
func buildRouteParamValidInstrumentIDFetcher(logger logging.Logger) func(req *http.Request) uint64 {
	return func(req *http.Request) uint64 {
		// we can generally disregard this error only because we should be able to validate.
		// that the string only contains numbers via chi's regex url param feature.
		u, err := strconv.ParseUint(chi.URLParam(req, validinstrumentsservice.URIParamKey), 10, 64)
		if err != nil {
			logger.Error(err, "fetching ValidInstrumentID from request")
		}
		return u
	}
}

// buildRouteParamValidIngredientIDFetcher builds a function that fetches a ValidIngredientID from a request routed by chi.
func buildRouteParamValidIngredientIDFetcher(logger logging.Logger) func(req *http.Request) uint64 {
	return func(req *http.Request) uint64 {
		// we can generally disregard this error only because we should be able to validate.
		// that the string only contains numbers via chi's regex url param feature.
		u, err := strconv.ParseUint(chi.URLParam(req, validingredientsservice.URIParamKey), 10, 64)
		if err != nil {
			logger.Error(err, "fetching ValidIngredientID from request")
		}
		return u
	}
}

// buildRouteParamValidIngredientTagIDFetcher builds a function that fetches a ValidIngredientTagID from a request routed by chi.
func buildRouteParamValidIngredientTagIDFetcher(logger logging.Logger) func(req *http.Request) uint64 {
	return func(req *http.Request) uint64 {
		// we can generally disregard this error only because we should be able to validate.
		// that the string only contains numbers via chi's regex url param feature.
		u, err := strconv.ParseUint(chi.URLParam(req, validingredienttagsservice.URIParamKey), 10, 64)
		if err != nil {
			logger.Error(err, "fetching ValidIngredientTagID from request")
		}
		return u
	}
}

// buildRouteParamIngredientTagMappingIDFetcher builds a function that fetches a IngredientTagMappingID from a request routed by chi.
func buildRouteParamIngredientTagMappingIDFetcher(logger logging.Logger) func(req *http.Request) uint64 {
	return func(req *http.Request) uint64 {
		// we can generally disregard this error only because we should be able to validate.
		// that the string only contains numbers via chi's regex url param feature.
		u, err := strconv.ParseUint(chi.URLParam(req, ingredienttagmappingsservice.URIParamKey), 10, 64)
		if err != nil {
			logger.Error(err, "fetching IngredientTagMappingID from request")
		}
		return u
	}
}

// buildRouteParamValidPreparationIDFetcher builds a function that fetches a ValidPreparationID from a request routed by chi.
func buildRouteParamValidPreparationIDFetcher(logger logging.Logger) func(req *http.Request) uint64 {
	return func(req *http.Request) uint64 {
		// we can generally disregard this error only because we should be able to validate.
		// that the string only contains numbers via chi's regex url param feature.
		u, err := strconv.ParseUint(chi.URLParam(req, validpreparationsservice.URIParamKey), 10, 64)
		if err != nil {
			logger.Error(err, "fetching ValidPreparationID from request")
		}
		return u
	}
}

// buildRouteParamRequiredPreparationInstrumentIDFetcher builds a function that fetches a RequiredPreparationInstrumentID from a request routed by chi.
func buildRouteParamRequiredPreparationInstrumentIDFetcher(logger logging.Logger) func(req *http.Request) uint64 {
	return func(req *http.Request) uint64 {
		// we can generally disregard this error only because we should be able to validate.
		// that the string only contains numbers via chi's regex url param feature.
		u, err := strconv.ParseUint(chi.URLParam(req, requiredpreparationinstrumentsservice.URIParamKey), 10, 64)
		if err != nil {
			logger.Error(err, "fetching RequiredPreparationInstrumentID from request")
		}
		return u
	}
}

// buildRouteParamValidIngredientPreparationIDFetcher builds a function that fetches a ValidIngredientPreparationID from a request routed by chi.
func buildRouteParamValidIngredientPreparationIDFetcher(logger logging.Logger) func(req *http.Request) uint64 {
	return func(req *http.Request) uint64 {
		// we can generally disregard this error only because we should be able to validate.
		// that the string only contains numbers via chi's regex url param feature.
		u, err := strconv.ParseUint(chi.URLParam(req, validingredientpreparationsservice.URIParamKey), 10, 64)
		if err != nil {
			logger.Error(err, "fetching ValidIngredientPreparationID from request")
		}
		return u
	}
}

// buildRouteParamRecipeIDFetcher builds a function that fetches a RecipeID from a request routed by chi.
func buildRouteParamRecipeIDFetcher(logger logging.Logger) func(req *http.Request) uint64 {
	return func(req *http.Request) uint64 {
		// we can generally disregard this error only because we should be able to validate.
		// that the string only contains numbers via chi's regex url param feature.
		u, err := strconv.ParseUint(chi.URLParam(req, recipesservice.URIParamKey), 10, 64)
		if err != nil {
			logger.Error(err, "fetching RecipeID from request")
		}
		return u
	}
}

// buildRouteParamRecipeTagIDFetcher builds a function that fetches a RecipeTagID from a request routed by chi.
func buildRouteParamRecipeTagIDFetcher(logger logging.Logger) func(req *http.Request) uint64 {
	return func(req *http.Request) uint64 {
		// we can generally disregard this error only because we should be able to validate.
		// that the string only contains numbers via chi's regex url param feature.
		u, err := strconv.ParseUint(chi.URLParam(req, recipetagsservice.URIParamKey), 10, 64)
		if err != nil {
			logger.Error(err, "fetching RecipeTagID from request")
		}
		return u
	}
}

// buildRouteParamRecipeStepIDFetcher builds a function that fetches a RecipeStepID from a request routed by chi.
func buildRouteParamRecipeStepIDFetcher(logger logging.Logger) func(req *http.Request) uint64 {
	return func(req *http.Request) uint64 {
		// we can generally disregard this error only because we should be able to validate.
		// that the string only contains numbers via chi's regex url param feature.
		u, err := strconv.ParseUint(chi.URLParam(req, recipestepsservice.URIParamKey), 10, 64)
		if err != nil {
			logger.Error(err, "fetching RecipeStepID from request")
		}
		return u
	}
}

// buildRouteParamRecipeStepPreparationIDFetcher builds a function that fetches a RecipeStepPreparationID from a request routed by chi.
func buildRouteParamRecipeStepPreparationIDFetcher(logger logging.Logger) func(req *http.Request) uint64 {
	return func(req *http.Request) uint64 {
		// we can generally disregard this error only because we should be able to validate.
		// that the string only contains numbers via chi's regex url param feature.
		u, err := strconv.ParseUint(chi.URLParam(req, recipesteppreparationsservice.URIParamKey), 10, 64)
		if err != nil {
			logger.Error(err, "fetching RecipeStepPreparationID from request")
		}
		return u
	}
}

// buildRouteParamRecipeStepIngredientIDFetcher builds a function that fetches a RecipeStepIngredientID from a request routed by chi.
func buildRouteParamRecipeStepIngredientIDFetcher(logger logging.Logger) func(req *http.Request) uint64 {
	return func(req *http.Request) uint64 {
		// we can generally disregard this error only because we should be able to validate.
		// that the string only contains numbers via chi's regex url param feature.
		u, err := strconv.ParseUint(chi.URLParam(req, recipestepingredientsservice.URIParamKey), 10, 64)
		if err != nil {
			logger.Error(err, "fetching RecipeStepIngredientID from request")
		}
		return u
	}
}

// buildRouteParamRecipeIterationIDFetcher builds a function that fetches a RecipeIterationID from a request routed by chi.
func buildRouteParamRecipeIterationIDFetcher(logger logging.Logger) func(req *http.Request) uint64 {
	return func(req *http.Request) uint64 {
		// we can generally disregard this error only because we should be able to validate.
		// that the string only contains numbers via chi's regex url param feature.
		u, err := strconv.ParseUint(chi.URLParam(req, recipeiterationsservice.URIParamKey), 10, 64)
		if err != nil {
			logger.Error(err, "fetching RecipeIterationID from request")
		}
		return u
	}
}

// buildRouteParamRecipeIterationStepIDFetcher builds a function that fetches a RecipeIterationStepID from a request routed by chi.
func buildRouteParamRecipeIterationStepIDFetcher(logger logging.Logger) func(req *http.Request) uint64 {
	return func(req *http.Request) uint64 {
		// we can generally disregard this error only because we should be able to validate.
		// that the string only contains numbers via chi's regex url param feature.
		u, err := strconv.ParseUint(chi.URLParam(req, recipeiterationstepsservice.URIParamKey), 10, 64)
		if err != nil {
			logger.Error(err, "fetching RecipeIterationStepID from request")
		}
		return u
	}
}

// buildRouteParamIterationMediaIDFetcher builds a function that fetches a IterationMediaID from a request routed by chi.
func buildRouteParamIterationMediaIDFetcher(logger logging.Logger) func(req *http.Request) uint64 {
	return func(req *http.Request) uint64 {
		// we can generally disregard this error only because we should be able to validate.
		// that the string only contains numbers via chi's regex url param feature.
		u, err := strconv.ParseUint(chi.URLParam(req, iterationmediasservice.URIParamKey), 10, 64)
		if err != nil {
			logger.Error(err, "fetching IterationMediaID from request")
		}
		return u
	}
}

// buildRouteParamInvitationIDFetcher builds a function that fetches a InvitationID from a request routed by chi.
func buildRouteParamInvitationIDFetcher(logger logging.Logger) func(req *http.Request) uint64 {
	return func(req *http.Request) uint64 {
		// we can generally disregard this error only because we should be able to validate.
		// that the string only contains numbers via chi's regex url param feature.
		u, err := strconv.ParseUint(chi.URLParam(req, invitationsservice.URIParamKey), 10, 64)
		if err != nil {
			logger.Error(err, "fetching InvitationID from request")
		}
		return u
	}
}

// buildRouteParamReportIDFetcher builds a function that fetches a ReportID from a request routed by chi.
func buildRouteParamReportIDFetcher(logger logging.Logger) func(req *http.Request) uint64 {
	return func(req *http.Request) uint64 {
		// we can generally disregard this error only because we should be able to validate.
		// that the string only contains numbers via chi's regex url param feature.
		u, err := strconv.ParseUint(chi.URLParam(req, reportsservice.URIParamKey), 10, 64)
		if err != nil {
			logger.Error(err, "fetching ReportID from request")
		}
		return u
	}
}

// buildRouteParamWebhookIDFetcher fetches a WebhookID from a request routed by chi.
func buildRouteParamWebhookIDFetcher(logger logging.Logger) func(req *http.Request) uint64 {
	return func(req *http.Request) uint64 {
		// we can generally disregard this error only because we should be able to validate.
		// that the string only contains numbers via chi's regex url param feature.
		u, err := strconv.ParseUint(chi.URLParam(req, webhooksservice.URIParamKey), 10, 64)
		if err != nil {
			logger.Error(err, "fetching WebhookID from request")
		}
		return u
	}
}

// buildRouteParamOAuth2ClientIDFetcher fetches a OAuth2ClientID from a request routed by chi.
func buildRouteParamOAuth2ClientIDFetcher(logger logging.Logger) func(req *http.Request) uint64 {
	return func(req *http.Request) uint64 {
		// we can generally disregard this error only because we should be able to validate.
		// that the string only contains numbers via chi's regex url param feature.
		u, err := strconv.ParseUint(chi.URLParam(req, oauth2clientsservice.URIParamKey), 10, 64)
		if err != nil {
			logger.Error(err, "fetching OAuth2ClientID from request")
		}
		return u
	}
}
