package httpserver

import (
	"net/http"
	"strconv"

	models "gitlab.com/prixfixe/prixfixe/models/v1"
	"gitlab.com/prixfixe/prixfixe/services/v1/auth"
	"gitlab.com/prixfixe/prixfixe/services/v1/ingredients"
	"gitlab.com/prixfixe/prixfixe/services/v1/instruments"
	"gitlab.com/prixfixe/prixfixe/services/v1/invitations"
	"gitlab.com/prixfixe/prixfixe/services/v1/iterationmedias"
	"gitlab.com/prixfixe/prixfixe/services/v1/oauth2clients"
	"gitlab.com/prixfixe/prixfixe/services/v1/preparations"
	"gitlab.com/prixfixe/prixfixe/services/v1/recipeiterations"
	"gitlab.com/prixfixe/prixfixe/services/v1/recipes"
	"gitlab.com/prixfixe/prixfixe/services/v1/recipestepevents"
	"gitlab.com/prixfixe/prixfixe/services/v1/recipestepingredients"
	"gitlab.com/prixfixe/prixfixe/services/v1/recipestepinstruments"
	"gitlab.com/prixfixe/prixfixe/services/v1/recipestepproducts"
	"gitlab.com/prixfixe/prixfixe/services/v1/recipesteps"
	"gitlab.com/prixfixe/prixfixe/services/v1/reports"
	"gitlab.com/prixfixe/prixfixe/services/v1/requiredpreparationinstruments"
	"gitlab.com/prixfixe/prixfixe/services/v1/users"
	"gitlab.com/prixfixe/prixfixe/services/v1/webhooks"

	"github.com/go-chi/chi"
	"github.com/google/wire"
	"gitlab.com/verygoodsoftwarenotvirus/logging/v1"
)

var (
	paramFetcherProviders = wire.NewSet(
		ProvideInstrumentServiceUserIDFetcher,
		ProvideIngredientServiceUserIDFetcher,
		ProvidePreparationServiceUserIDFetcher,
		ProvideRequiredPreparationInstrumentServiceUserIDFetcher,
		ProvideRecipeServiceUserIDFetcher,
		ProvideRecipeStepServiceUserIDFetcher,
		ProvideRecipeStepInstrumentServiceUserIDFetcher,
		ProvideRecipeStepIngredientServiceUserIDFetcher,
		ProvideRecipeStepProductServiceUserIDFetcher,
		ProvideRecipeIterationServiceUserIDFetcher,
		ProvideRecipeStepEventServiceUserIDFetcher,
		ProvideIterationMediaServiceUserIDFetcher,
		ProvideInvitationServiceUserIDFetcher,
		ProvideReportServiceUserIDFetcher,
		ProvideUsernameFetcher,
		ProvideOAuth2ServiceClientIDFetcher,
		ProvideAuthUserIDFetcher,
		ProvideInstrumentIDFetcher,
		ProvideIngredientIDFetcher,
		ProvidePreparationIDFetcher,
		ProvideRequiredPreparationInstrumentIDFetcher,
		ProvideRecipeIDFetcher,
		ProvideRecipeStepIDFetcher,
		ProvideRecipeStepInstrumentIDFetcher,
		ProvideRecipeStepIngredientIDFetcher,
		ProvideRecipeStepProductIDFetcher,
		ProvideRecipeIterationIDFetcher,
		ProvideRecipeStepEventIDFetcher,
		ProvideIterationMediaIDFetcher,
		ProvideInvitationIDFetcher,
		ProvideReportIDFetcher,
		ProvideWebhooksUserIDFetcher,
		ProvideWebhookIDFetcher,
	)
)

// ProvideInstrumentServiceUserIDFetcher provides a UserIDFetcher
func ProvideInstrumentServiceUserIDFetcher() instruments.UserIDFetcher {
	return UserIDFetcher
}

// ProvideInstrumentIDFetcher provides an InstrumentIDFetcher
func ProvideInstrumentIDFetcher(logger logging.Logger) instruments.InstrumentIDFetcher {
	return buildChiInstrumentIDFetcher(logger)
}

// ProvideIngredientServiceUserIDFetcher provides a UserIDFetcher
func ProvideIngredientServiceUserIDFetcher() ingredients.UserIDFetcher {
	return UserIDFetcher
}

// ProvideIngredientIDFetcher provides an IngredientIDFetcher
func ProvideIngredientIDFetcher(logger logging.Logger) ingredients.IngredientIDFetcher {
	return buildChiIngredientIDFetcher(logger)
}

// ProvidePreparationServiceUserIDFetcher provides a UserIDFetcher
func ProvidePreparationServiceUserIDFetcher() preparations.UserIDFetcher {
	return UserIDFetcher
}

// ProvidePreparationIDFetcher provides an PreparationIDFetcher
func ProvidePreparationIDFetcher(logger logging.Logger) preparations.PreparationIDFetcher {
	return buildChiPreparationIDFetcher(logger)
}

// ProvideRequiredPreparationInstrumentServiceUserIDFetcher provides a UserIDFetcher
func ProvideRequiredPreparationInstrumentServiceUserIDFetcher() requiredpreparationinstruments.UserIDFetcher {
	return UserIDFetcher
}

// ProvideRequiredPreparationInstrumentIDFetcher provides an RequiredPreparationInstrumentIDFetcher
func ProvideRequiredPreparationInstrumentIDFetcher(logger logging.Logger) requiredpreparationinstruments.RequiredPreparationInstrumentIDFetcher {
	return buildChiRequiredPreparationInstrumentIDFetcher(logger)
}

// ProvideRecipeServiceUserIDFetcher provides a UserIDFetcher
func ProvideRecipeServiceUserIDFetcher() recipes.UserIDFetcher {
	return UserIDFetcher
}

// ProvideRecipeIDFetcher provides an RecipeIDFetcher
func ProvideRecipeIDFetcher(logger logging.Logger) recipes.RecipeIDFetcher {
	return buildChiRecipeIDFetcher(logger)
}

// ProvideRecipeStepServiceUserIDFetcher provides a UserIDFetcher
func ProvideRecipeStepServiceUserIDFetcher() recipesteps.UserIDFetcher {
	return UserIDFetcher
}

// ProvideRecipeStepIDFetcher provides an RecipeStepIDFetcher
func ProvideRecipeStepIDFetcher(logger logging.Logger) recipesteps.RecipeStepIDFetcher {
	return buildChiRecipeStepIDFetcher(logger)
}

// ProvideRecipeStepInstrumentServiceUserIDFetcher provides a UserIDFetcher
func ProvideRecipeStepInstrumentServiceUserIDFetcher() recipestepinstruments.UserIDFetcher {
	return UserIDFetcher
}

// ProvideRecipeStepInstrumentIDFetcher provides an RecipeStepInstrumentIDFetcher
func ProvideRecipeStepInstrumentIDFetcher(logger logging.Logger) recipestepinstruments.RecipeStepInstrumentIDFetcher {
	return buildChiRecipeStepInstrumentIDFetcher(logger)
}

// ProvideRecipeStepIngredientServiceUserIDFetcher provides a UserIDFetcher
func ProvideRecipeStepIngredientServiceUserIDFetcher() recipestepingredients.UserIDFetcher {
	return UserIDFetcher
}

// ProvideRecipeStepIngredientIDFetcher provides an RecipeStepIngredientIDFetcher
func ProvideRecipeStepIngredientIDFetcher(logger logging.Logger) recipestepingredients.RecipeStepIngredientIDFetcher {
	return buildChiRecipeStepIngredientIDFetcher(logger)
}

// ProvideRecipeStepProductServiceUserIDFetcher provides a UserIDFetcher
func ProvideRecipeStepProductServiceUserIDFetcher() recipestepproducts.UserIDFetcher {
	return UserIDFetcher
}

// ProvideRecipeStepProductIDFetcher provides an RecipeStepProductIDFetcher
func ProvideRecipeStepProductIDFetcher(logger logging.Logger) recipestepproducts.RecipeStepProductIDFetcher {
	return buildChiRecipeStepProductIDFetcher(logger)
}

// ProvideRecipeIterationServiceUserIDFetcher provides a UserIDFetcher
func ProvideRecipeIterationServiceUserIDFetcher() recipeiterations.UserIDFetcher {
	return UserIDFetcher
}

// ProvideRecipeIterationIDFetcher provides an RecipeIterationIDFetcher
func ProvideRecipeIterationIDFetcher(logger logging.Logger) recipeiterations.RecipeIterationIDFetcher {
	return buildChiRecipeIterationIDFetcher(logger)
}

// ProvideRecipeStepEventServiceUserIDFetcher provides a UserIDFetcher
func ProvideRecipeStepEventServiceUserIDFetcher() recipestepevents.UserIDFetcher {
	return UserIDFetcher
}

// ProvideRecipeStepEventIDFetcher provides an RecipeStepEventIDFetcher
func ProvideRecipeStepEventIDFetcher(logger logging.Logger) recipestepevents.RecipeStepEventIDFetcher {
	return buildChiRecipeStepEventIDFetcher(logger)
}

// ProvideIterationMediaServiceUserIDFetcher provides a UserIDFetcher
func ProvideIterationMediaServiceUserIDFetcher() iterationmedias.UserIDFetcher {
	return UserIDFetcher
}

// ProvideIterationMediaIDFetcher provides an IterationMediaIDFetcher
func ProvideIterationMediaIDFetcher(logger logging.Logger) iterationmedias.IterationMediaIDFetcher {
	return buildChiIterationMediaIDFetcher(logger)
}

// ProvideInvitationServiceUserIDFetcher provides a UserIDFetcher
func ProvideInvitationServiceUserIDFetcher() invitations.UserIDFetcher {
	return UserIDFetcher
}

// ProvideInvitationIDFetcher provides an InvitationIDFetcher
func ProvideInvitationIDFetcher(logger logging.Logger) invitations.InvitationIDFetcher {
	return buildChiInvitationIDFetcher(logger)
}

// ProvideReportServiceUserIDFetcher provides a UserIDFetcher
func ProvideReportServiceUserIDFetcher() reports.UserIDFetcher {
	return UserIDFetcher
}

// ProvideReportIDFetcher provides an ReportIDFetcher
func ProvideReportIDFetcher(logger logging.Logger) reports.ReportIDFetcher {
	return buildChiReportIDFetcher(logger)
}

// ProvideUsernameFetcher provides a UsernameFetcher
func ProvideUsernameFetcher(logger logging.Logger) users.UserIDFetcher {
	return buildChiUserIDFetcher(logger)
}

// ProvideAuthUserIDFetcher provides a UsernameFetcher
func ProvideAuthUserIDFetcher() auth.UserIDFetcher {
	return UserIDFetcher
}

// ProvideWebhooksUserIDFetcher provides a UserIDFetcher
func ProvideWebhooksUserIDFetcher() webhooks.UserIDFetcher {
	return UserIDFetcher
}

// ProvideWebhookIDFetcher provides an WebhookIDFetcher
func ProvideWebhookIDFetcher(logger logging.Logger) webhooks.WebhookIDFetcher {
	return buildChiWebhookIDFetcher(logger)
}

// ProvideOAuth2ServiceClientIDFetcher provides a ClientIDFetcher
func ProvideOAuth2ServiceClientIDFetcher(logger logging.Logger) oauth2clients.ClientIDFetcher {
	return buildChiOAuth2ClientIDFetcher(logger)
}

// UserIDFetcher fetches a user ID from a request routed by chi.
func UserIDFetcher(req *http.Request) uint64 {
	if userID, ok := req.Context().Value(models.UserIDKey).(uint64); ok {
		return userID
	} else {
		return 0
	}
}

// buildChiUserIDFetcher builds a function that fetches a Username from a request routed by chi.
func buildChiUserIDFetcher(logger logging.Logger) users.UserIDFetcher {
	return func(req *http.Request) uint64 {
		u, err := strconv.ParseUint(chi.URLParam(req, users.URIParamKey), 10, 64)
		if err != nil {
			logger.Error(err, "fetching user ID from request")
		}
		return u
	}
}

// buildChiInstrumentIDFetcher builds a function that fetches a InstrumentID from a request routed by chi.
func buildChiInstrumentIDFetcher(logger logging.Logger) func(req *http.Request) uint64 {
	return func(req *http.Request) uint64 {
		// we can generally disregard this error only because we should be able to validate
		// that the string only contains numbers via chi's regex url param feature.
		u, err := strconv.ParseUint(chi.URLParam(req, instruments.URIParamKey), 10, 64)
		if err != nil {
			logger.Error(err, "fetching InstrumentID from request")
		}
		return u
	}
}

// buildChiIngredientIDFetcher builds a function that fetches a IngredientID from a request routed by chi.
func buildChiIngredientIDFetcher(logger logging.Logger) func(req *http.Request) uint64 {
	return func(req *http.Request) uint64 {
		// we can generally disregard this error only because we should be able to validate
		// that the string only contains numbers via chi's regex url param feature.
		u, err := strconv.ParseUint(chi.URLParam(req, ingredients.URIParamKey), 10, 64)
		if err != nil {
			logger.Error(err, "fetching IngredientID from request")
		}
		return u
	}
}

// buildChiPreparationIDFetcher builds a function that fetches a PreparationID from a request routed by chi.
func buildChiPreparationIDFetcher(logger logging.Logger) func(req *http.Request) uint64 {
	return func(req *http.Request) uint64 {
		// we can generally disregard this error only because we should be able to validate
		// that the string only contains numbers via chi's regex url param feature.
		u, err := strconv.ParseUint(chi.URLParam(req, preparations.URIParamKey), 10, 64)
		if err != nil {
			logger.Error(err, "fetching PreparationID from request")
		}
		return u
	}
}

// buildChiRequiredPreparationInstrumentIDFetcher builds a function that fetches a RequiredPreparationInstrumentID from a request routed by chi.
func buildChiRequiredPreparationInstrumentIDFetcher(logger logging.Logger) func(req *http.Request) uint64 {
	return func(req *http.Request) uint64 {
		// we can generally disregard this error only because we should be able to validate
		// that the string only contains numbers via chi's regex url param feature.
		u, err := strconv.ParseUint(chi.URLParam(req, requiredpreparationinstruments.URIParamKey), 10, 64)
		if err != nil {
			logger.Error(err, "fetching RequiredPreparationInstrumentID from request")
		}
		return u
	}
}

// buildChiRecipeIDFetcher builds a function that fetches a RecipeID from a request routed by chi.
func buildChiRecipeIDFetcher(logger logging.Logger) func(req *http.Request) uint64 {
	return func(req *http.Request) uint64 {
		// we can generally disregard this error only because we should be able to validate
		// that the string only contains numbers via chi's regex url param feature.
		u, err := strconv.ParseUint(chi.URLParam(req, recipes.URIParamKey), 10, 64)
		if err != nil {
			logger.Error(err, "fetching RecipeID from request")
		}
		return u
	}
}

// buildChiRecipeStepIDFetcher builds a function that fetches a RecipeStepID from a request routed by chi.
func buildChiRecipeStepIDFetcher(logger logging.Logger) func(req *http.Request) uint64 {
	return func(req *http.Request) uint64 {
		// we can generally disregard this error only because we should be able to validate
		// that the string only contains numbers via chi's regex url param feature.
		u, err := strconv.ParseUint(chi.URLParam(req, recipesteps.URIParamKey), 10, 64)
		if err != nil {
			logger.Error(err, "fetching RecipeStepID from request")
		}
		return u
	}
}

// buildChiRecipeStepInstrumentIDFetcher builds a function that fetches a RecipeStepInstrumentID from a request routed by chi.
func buildChiRecipeStepInstrumentIDFetcher(logger logging.Logger) func(req *http.Request) uint64 {
	return func(req *http.Request) uint64 {
		// we can generally disregard this error only because we should be able to validate
		// that the string only contains numbers via chi's regex url param feature.
		u, err := strconv.ParseUint(chi.URLParam(req, recipestepinstruments.URIParamKey), 10, 64)
		if err != nil {
			logger.Error(err, "fetching RecipeStepInstrumentID from request")
		}
		return u
	}
}

// buildChiRecipeStepIngredientIDFetcher builds a function that fetches a RecipeStepIngredientID from a request routed by chi.
func buildChiRecipeStepIngredientIDFetcher(logger logging.Logger) func(req *http.Request) uint64 {
	return func(req *http.Request) uint64 {
		// we can generally disregard this error only because we should be able to validate
		// that the string only contains numbers via chi's regex url param feature.
		u, err := strconv.ParseUint(chi.URLParam(req, recipestepingredients.URIParamKey), 10, 64)
		if err != nil {
			logger.Error(err, "fetching RecipeStepIngredientID from request")
		}
		return u
	}
}

// buildChiRecipeStepProductIDFetcher builds a function that fetches a RecipeStepProductID from a request routed by chi.
func buildChiRecipeStepProductIDFetcher(logger logging.Logger) func(req *http.Request) uint64 {
	return func(req *http.Request) uint64 {
		// we can generally disregard this error only because we should be able to validate
		// that the string only contains numbers via chi's regex url param feature.
		u, err := strconv.ParseUint(chi.URLParam(req, recipestepproducts.URIParamKey), 10, 64)
		if err != nil {
			logger.Error(err, "fetching RecipeStepProductID from request")
		}
		return u
	}
}

// buildChiRecipeIterationIDFetcher builds a function that fetches a RecipeIterationID from a request routed by chi.
func buildChiRecipeIterationIDFetcher(logger logging.Logger) func(req *http.Request) uint64 {
	return func(req *http.Request) uint64 {
		// we can generally disregard this error only because we should be able to validate
		// that the string only contains numbers via chi's regex url param feature.
		u, err := strconv.ParseUint(chi.URLParam(req, recipeiterations.URIParamKey), 10, 64)
		if err != nil {
			logger.Error(err, "fetching RecipeIterationID from request")
		}
		return u
	}
}

// buildChiRecipeStepEventIDFetcher builds a function that fetches a RecipeStepEventID from a request routed by chi.
func buildChiRecipeStepEventIDFetcher(logger logging.Logger) func(req *http.Request) uint64 {
	return func(req *http.Request) uint64 {
		// we can generally disregard this error only because we should be able to validate
		// that the string only contains numbers via chi's regex url param feature.
		u, err := strconv.ParseUint(chi.URLParam(req, recipestepevents.URIParamKey), 10, 64)
		if err != nil {
			logger.Error(err, "fetching RecipeStepEventID from request")
		}
		return u
	}
}

// buildChiIterationMediaIDFetcher builds a function that fetches a IterationMediaID from a request routed by chi.
func buildChiIterationMediaIDFetcher(logger logging.Logger) func(req *http.Request) uint64 {
	return func(req *http.Request) uint64 {
		// we can generally disregard this error only because we should be able to validate
		// that the string only contains numbers via chi's regex url param feature.
		u, err := strconv.ParseUint(chi.URLParam(req, iterationmedias.URIParamKey), 10, 64)
		if err != nil {
			logger.Error(err, "fetching IterationMediaID from request")
		}
		return u
	}
}

// buildChiInvitationIDFetcher builds a function that fetches a InvitationID from a request routed by chi.
func buildChiInvitationIDFetcher(logger logging.Logger) func(req *http.Request) uint64 {
	return func(req *http.Request) uint64 {
		// we can generally disregard this error only because we should be able to validate
		// that the string only contains numbers via chi's regex url param feature.
		u, err := strconv.ParseUint(chi.URLParam(req, invitations.URIParamKey), 10, 64)
		if err != nil {
			logger.Error(err, "fetching InvitationID from request")
		}
		return u
	}
}

// buildChiReportIDFetcher builds a function that fetches a ReportID from a request routed by chi.
func buildChiReportIDFetcher(logger logging.Logger) func(req *http.Request) uint64 {
	return func(req *http.Request) uint64 {
		// we can generally disregard this error only because we should be able to validate
		// that the string only contains numbers via chi's regex url param feature.
		u, err := strconv.ParseUint(chi.URLParam(req, reports.URIParamKey), 10, 64)
		if err != nil {
			logger.Error(err, "fetching ReportID from request")
		}
		return u
	}
}

// chiWebhookIDFetcher fetches a WebhookID from a request routed by chi.
func buildChiWebhookIDFetcher(logger logging.Logger) func(req *http.Request) uint64 {
	return func(req *http.Request) uint64 {
		// we can generally disregard this error only because we should be able to validate
		// that the string only contains numbers via chi's regex url param feature.
		u, err := strconv.ParseUint(chi.URLParam(req, webhooks.URIParamKey), 10, 64)
		if err != nil {
			logger.Error(err, "fetching WebhookID from request")
		}
		return u
	}
}

// chiOAuth2ClientIDFetcher fetches a OAuth2ClientID from a request routed by chi.
func buildChiOAuth2ClientIDFetcher(logger logging.Logger) func(req *http.Request) uint64 {
	return func(req *http.Request) uint64 {
		// we can generally disregard this error only because we should be able to validate
		// that the string only contains numbers via chi's regex url param feature.
		u, err := strconv.ParseUint(chi.URLParam(req, oauth2clients.URIParamKey), 10, 64)
		if err != nil {
			logger.Error(err, "fetching OAuth2ClientID from request")
		}
		return u
	}
}
