package httpserver

import (
	"net/http"
	"strconv"

	models "gitlab.com/prixfixe/prixfixe/models/v1"
	invitationsservice "gitlab.com/prixfixe/prixfixe/services/v1/invitations"
	iterationmediasservice "gitlab.com/prixfixe/prixfixe/services/v1/iterationmedias"
	oauth2clientsservice "gitlab.com/prixfixe/prixfixe/services/v1/oauth2clients"
	recipeiterationsservice "gitlab.com/prixfixe/prixfixe/services/v1/recipeiterations"
	recipesservice "gitlab.com/prixfixe/prixfixe/services/v1/recipes"
	recipestepeventsservice "gitlab.com/prixfixe/prixfixe/services/v1/recipestepevents"
	recipestepingredientsservice "gitlab.com/prixfixe/prixfixe/services/v1/recipestepingredients"
	recipestepinstrumentsservice "gitlab.com/prixfixe/prixfixe/services/v1/recipestepinstruments"
	recipestepproductsservice "gitlab.com/prixfixe/prixfixe/services/v1/recipestepproducts"
	recipestepsservice "gitlab.com/prixfixe/prixfixe/services/v1/recipesteps"
	reportsservice "gitlab.com/prixfixe/prixfixe/services/v1/reports"
	requiredpreparationinstrumentsservice "gitlab.com/prixfixe/prixfixe/services/v1/requiredpreparationinstruments"
	usersservice "gitlab.com/prixfixe/prixfixe/services/v1/users"
	validingredientpreparationsservice "gitlab.com/prixfixe/prixfixe/services/v1/validingredientpreparations"
	validingredientsservice "gitlab.com/prixfixe/prixfixe/services/v1/validingredients"
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
		ProvideWebhooksServiceWebhookIDFetcher,
		ProvideWebhooksServiceUserIDFetcher,
		ProvideValidInstrumentsServiceValidInstrumentIDFetcher,
		ProvideValidIngredientsServiceValidIngredientIDFetcher,
		ProvideValidPreparationsServiceValidPreparationIDFetcher,
		ProvideValidIngredientPreparationsServiceValidIngredientPreparationIDFetcher,
		ProvideRequiredPreparationInstrumentsServiceRequiredPreparationInstrumentIDFetcher,
		ProvideRecipesServiceRecipeIDFetcher,
		ProvideRecipesServiceUserIDFetcher,
		ProvideRecipeStepsServiceRecipeIDFetcher,
		ProvideRecipeStepsServiceRecipeStepIDFetcher,
		ProvideRecipeStepsServiceUserIDFetcher,
		ProvideRecipeStepInstrumentsServiceRecipeIDFetcher,
		ProvideRecipeStepInstrumentsServiceRecipeStepIDFetcher,
		ProvideRecipeStepInstrumentsServiceRecipeStepInstrumentIDFetcher,
		ProvideRecipeStepInstrumentsServiceUserIDFetcher,
		ProvideRecipeStepIngredientsServiceRecipeIDFetcher,
		ProvideRecipeStepIngredientsServiceRecipeStepIDFetcher,
		ProvideRecipeStepIngredientsServiceRecipeStepIngredientIDFetcher,
		ProvideRecipeStepIngredientsServiceUserIDFetcher,
		ProvideRecipeStepProductsServiceRecipeIDFetcher,
		ProvideRecipeStepProductsServiceRecipeStepIDFetcher,
		ProvideRecipeStepProductsServiceRecipeStepProductIDFetcher,
		ProvideRecipeStepProductsServiceUserIDFetcher,
		ProvideRecipeIterationsServiceRecipeIDFetcher,
		ProvideRecipeIterationsServiceRecipeIterationIDFetcher,
		ProvideRecipeIterationsServiceUserIDFetcher,
		ProvideRecipeStepEventsServiceRecipeIDFetcher,
		ProvideRecipeStepEventsServiceRecipeStepIDFetcher,
		ProvideRecipeStepEventsServiceRecipeStepEventIDFetcher,
		ProvideRecipeStepEventsServiceUserIDFetcher,
		ProvideIterationMediasServiceRecipeIDFetcher,
		ProvideIterationMediasServiceRecipeIterationIDFetcher,
		ProvideIterationMediasServiceIterationMediaIDFetcher,
		ProvideIterationMediasServiceUserIDFetcher,
		ProvideInvitationsServiceInvitationIDFetcher,
		ProvideInvitationsServiceUserIDFetcher,
		ProvideReportsServiceReportIDFetcher,
		ProvideReportsServiceUserIDFetcher,
	)
)

// ProvideUsersServiceUserIDFetcher provides a UsernameFetcher.
func ProvideUsersServiceUserIDFetcher(logger logging.Logger) usersservice.UserIDFetcher {
	return buildRouteParamUserIDFetcher(logger)
}

// ProvideOAuth2ClientsServiceClientIDFetcher provides a ClientIDFetcher.
func ProvideOAuth2ClientsServiceClientIDFetcher(logger logging.Logger) oauth2clientsservice.ClientIDFetcher {
	return buildRouteParamOAuth2ClientIDFetcher(logger)
}

// ProvideWebhooksServiceWebhookIDFetcher provides an WebhookIDFetcher.
func ProvideWebhooksServiceWebhookIDFetcher(logger logging.Logger) webhooksservice.WebhookIDFetcher {
	return buildRouteParamWebhookIDFetcher(logger)
}

// ProvideWebhooksServiceUserIDFetcher provides a UserIDFetcher.
func ProvideWebhooksServiceUserIDFetcher() webhooksservice.UserIDFetcher {
	return userIDFetcherFromRequestContext
}

// ProvideValidInstrumentsServiceValidInstrumentIDFetcher provides a ValidInstrumentIDFetcher.
func ProvideValidInstrumentsServiceValidInstrumentIDFetcher(logger logging.Logger) validinstrumentsservice.ValidInstrumentIDFetcher {
	return buildRouteParamValidInstrumentIDFetcher(logger)
}

// ProvideValidIngredientsServiceValidIngredientIDFetcher provides a ValidIngredientIDFetcher.
func ProvideValidIngredientsServiceValidIngredientIDFetcher(logger logging.Logger) validingredientsservice.ValidIngredientIDFetcher {
	return buildRouteParamValidIngredientIDFetcher(logger)
}

// ProvideValidPreparationsServiceValidPreparationIDFetcher provides a ValidPreparationIDFetcher.
func ProvideValidPreparationsServiceValidPreparationIDFetcher(logger logging.Logger) validpreparationsservice.ValidPreparationIDFetcher {
	return buildRouteParamValidPreparationIDFetcher(logger)
}

// ProvideValidIngredientPreparationsServiceValidIngredientPreparationIDFetcher provides a ValidIngredientPreparationIDFetcher.
func ProvideValidIngredientPreparationsServiceValidIngredientPreparationIDFetcher(logger logging.Logger) validingredientpreparationsservice.ValidIngredientPreparationIDFetcher {
	return buildRouteParamValidIngredientPreparationIDFetcher(logger)
}

// ProvideRequiredPreparationInstrumentsServiceRequiredPreparationInstrumentIDFetcher provides a RequiredPreparationInstrumentIDFetcher.
func ProvideRequiredPreparationInstrumentsServiceRequiredPreparationInstrumentIDFetcher(logger logging.Logger) requiredpreparationinstrumentsservice.RequiredPreparationInstrumentIDFetcher {
	return buildRouteParamRequiredPreparationInstrumentIDFetcher(logger)
}

// ProvideRecipesServiceRecipeIDFetcher provides a RecipeIDFetcher.
func ProvideRecipesServiceRecipeIDFetcher(logger logging.Logger) recipesservice.RecipeIDFetcher {
	return buildRouteParamRecipeIDFetcher(logger)
}

// ProvideRecipesServiceUserIDFetcher provides a UserIDFetcher.
func ProvideRecipesServiceUserIDFetcher() recipesservice.UserIDFetcher {
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

// ProvideRecipeStepsServiceUserIDFetcher provides a UserIDFetcher.
func ProvideRecipeStepsServiceUserIDFetcher() recipestepsservice.UserIDFetcher {
	return userIDFetcherFromRequestContext
}

// ProvideRecipeStepInstrumentsServiceRecipeIDFetcher provides a RecipeIDFetcher.
func ProvideRecipeStepInstrumentsServiceRecipeIDFetcher(logger logging.Logger) recipestepinstrumentsservice.RecipeIDFetcher {
	return buildRouteParamRecipeIDFetcher(logger)
}

// ProvideRecipeStepInstrumentsServiceRecipeStepIDFetcher provides a RecipeStepIDFetcher.
func ProvideRecipeStepInstrumentsServiceRecipeStepIDFetcher(logger logging.Logger) recipestepinstrumentsservice.RecipeStepIDFetcher {
	return buildRouteParamRecipeStepIDFetcher(logger)
}

// ProvideRecipeStepInstrumentsServiceRecipeStepInstrumentIDFetcher provides a RecipeStepInstrumentIDFetcher.
func ProvideRecipeStepInstrumentsServiceRecipeStepInstrumentIDFetcher(logger logging.Logger) recipestepinstrumentsservice.RecipeStepInstrumentIDFetcher {
	return buildRouteParamRecipeStepInstrumentIDFetcher(logger)
}

// ProvideRecipeStepInstrumentsServiceUserIDFetcher provides a UserIDFetcher.
func ProvideRecipeStepInstrumentsServiceUserIDFetcher() recipestepinstrumentsservice.UserIDFetcher {
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

// ProvideRecipeStepIngredientsServiceUserIDFetcher provides a UserIDFetcher.
func ProvideRecipeStepIngredientsServiceUserIDFetcher() recipestepingredientsservice.UserIDFetcher {
	return userIDFetcherFromRequestContext
}

// ProvideRecipeStepProductsServiceRecipeIDFetcher provides a RecipeIDFetcher.
func ProvideRecipeStepProductsServiceRecipeIDFetcher(logger logging.Logger) recipestepproductsservice.RecipeIDFetcher {
	return buildRouteParamRecipeIDFetcher(logger)
}

// ProvideRecipeStepProductsServiceRecipeStepIDFetcher provides a RecipeStepIDFetcher.
func ProvideRecipeStepProductsServiceRecipeStepIDFetcher(logger logging.Logger) recipestepproductsservice.RecipeStepIDFetcher {
	return buildRouteParamRecipeStepIDFetcher(logger)
}

// ProvideRecipeStepProductsServiceRecipeStepProductIDFetcher provides a RecipeStepProductIDFetcher.
func ProvideRecipeStepProductsServiceRecipeStepProductIDFetcher(logger logging.Logger) recipestepproductsservice.RecipeStepProductIDFetcher {
	return buildRouteParamRecipeStepProductIDFetcher(logger)
}

// ProvideRecipeStepProductsServiceUserIDFetcher provides a UserIDFetcher.
func ProvideRecipeStepProductsServiceUserIDFetcher() recipestepproductsservice.UserIDFetcher {
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

// ProvideRecipeIterationsServiceUserIDFetcher provides a UserIDFetcher.
func ProvideRecipeIterationsServiceUserIDFetcher() recipeiterationsservice.UserIDFetcher {
	return userIDFetcherFromRequestContext
}

// ProvideRecipeStepEventsServiceRecipeIDFetcher provides a RecipeIDFetcher.
func ProvideRecipeStepEventsServiceRecipeIDFetcher(logger logging.Logger) recipestepeventsservice.RecipeIDFetcher {
	return buildRouteParamRecipeIDFetcher(logger)
}

// ProvideRecipeStepEventsServiceRecipeStepIDFetcher provides a RecipeStepIDFetcher.
func ProvideRecipeStepEventsServiceRecipeStepIDFetcher(logger logging.Logger) recipestepeventsservice.RecipeStepIDFetcher {
	return buildRouteParamRecipeStepIDFetcher(logger)
}

// ProvideRecipeStepEventsServiceRecipeStepEventIDFetcher provides a RecipeStepEventIDFetcher.
func ProvideRecipeStepEventsServiceRecipeStepEventIDFetcher(logger logging.Logger) recipestepeventsservice.RecipeStepEventIDFetcher {
	return buildRouteParamRecipeStepEventIDFetcher(logger)
}

// ProvideRecipeStepEventsServiceUserIDFetcher provides a UserIDFetcher.
func ProvideRecipeStepEventsServiceUserIDFetcher() recipestepeventsservice.UserIDFetcher {
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

// ProvideIterationMediasServiceUserIDFetcher provides a UserIDFetcher.
func ProvideIterationMediasServiceUserIDFetcher() iterationmediasservice.UserIDFetcher {
	return userIDFetcherFromRequestContext
}

// ProvideInvitationsServiceInvitationIDFetcher provides an InvitationIDFetcher.
func ProvideInvitationsServiceInvitationIDFetcher(logger logging.Logger) invitationsservice.InvitationIDFetcher {
	return buildRouteParamInvitationIDFetcher(logger)
}

// ProvideInvitationsServiceUserIDFetcher provides a UserIDFetcher.
func ProvideInvitationsServiceUserIDFetcher() invitationsservice.UserIDFetcher {
	return userIDFetcherFromRequestContext
}

// ProvideReportsServiceReportIDFetcher provides a ReportIDFetcher.
func ProvideReportsServiceReportIDFetcher(logger logging.Logger) reportsservice.ReportIDFetcher {
	return buildRouteParamReportIDFetcher(logger)
}

// ProvideReportsServiceUserIDFetcher provides a UserIDFetcher.
func ProvideReportsServiceUserIDFetcher() reportsservice.UserIDFetcher {
	return userIDFetcherFromRequestContext
}

// userIDFetcherFromRequestContext fetches a user ID from a request routed by chi.
func userIDFetcherFromRequestContext(req *http.Request) uint64 {
	if si, ok := req.Context().Value(models.SessionInfoKey).(*models.SessionInfo); ok && si != nil {
		return si.UserID
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

// buildRouteParamRecipeStepInstrumentIDFetcher builds a function that fetches a RecipeStepInstrumentID from a request routed by chi.
func buildRouteParamRecipeStepInstrumentIDFetcher(logger logging.Logger) func(req *http.Request) uint64 {
	return func(req *http.Request) uint64 {
		// we can generally disregard this error only because we should be able to validate.
		// that the string only contains numbers via chi's regex url param feature.
		u, err := strconv.ParseUint(chi.URLParam(req, recipestepinstrumentsservice.URIParamKey), 10, 64)
		if err != nil {
			logger.Error(err, "fetching RecipeStepInstrumentID from request")
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

// buildRouteParamRecipeStepProductIDFetcher builds a function that fetches a RecipeStepProductID from a request routed by chi.
func buildRouteParamRecipeStepProductIDFetcher(logger logging.Logger) func(req *http.Request) uint64 {
	return func(req *http.Request) uint64 {
		// we can generally disregard this error only because we should be able to validate.
		// that the string only contains numbers via chi's regex url param feature.
		u, err := strconv.ParseUint(chi.URLParam(req, recipestepproductsservice.URIParamKey), 10, 64)
		if err != nil {
			logger.Error(err, "fetching RecipeStepProductID from request")
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

// buildRouteParamRecipeStepEventIDFetcher builds a function that fetches a RecipeStepEventID from a request routed by chi.
func buildRouteParamRecipeStepEventIDFetcher(logger logging.Logger) func(req *http.Request) uint64 {
	return func(req *http.Request) uint64 {
		// we can generally disregard this error only because we should be able to validate.
		// that the string only contains numbers via chi's regex url param feature.
		u, err := strconv.ParseUint(chi.URLParam(req, recipestepeventsservice.URIParamKey), 10, 64)
		if err != nil {
			logger.Error(err, "fetching RecipeStepEventID from request")
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
