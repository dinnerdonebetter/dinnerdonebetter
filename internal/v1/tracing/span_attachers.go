package tracing

import (
	"strconv"

	models "gitlab.com/prixfixe/prixfixe/models/v1"

	"go.opencensus.io/trace"
)

const (
	validInstrumentIDSpanAttachmentKey               = "valid_instrument_id"
	validIngredientIDSpanAttachmentKey               = "valid_ingredient_id"
	validPreparationIDSpanAttachmentKey              = "valid_preparation_id"
	validIngredientPreparationIDSpanAttachmentKey    = "valid_ingredient_preparation_id"
	requiredPreparationInstrumentIDSpanAttachmentKey = "required_preparation_instrument_id"
	recipeIDSpanAttachmentKey                        = "recipe_id"
	recipeStepIDSpanAttachmentKey                    = "recipe_step_id"
	recipeStepInstrumentIDSpanAttachmentKey          = "recipe_step_instrument_id"
	recipeStepIngredientIDSpanAttachmentKey          = "recipe_step_ingredient_id"
	recipeStepProductIDSpanAttachmentKey             = "recipe_step_product_id"
	recipeIterationIDSpanAttachmentKey               = "recipe_iteration_id"
	recipeStepEventIDSpanAttachmentKey               = "recipe_step_event_id"
	iterationMediaIDSpanAttachmentKey                = "iteration_media_id"
	invitationIDSpanAttachmentKey                    = "invitation_id"
	reportIDSpanAttachmentKey                        = "report_id"
	userIDSpanAttachmentKey                          = "user_id"
	usernameSpanAttachmentKey                        = "username"
	filterPageSpanAttachmentKey                      = "filter_page"
	filterLimitSpanAttachmentKey                     = "filter_limit"
	oauth2ClientDatabaseIDSpanAttachmentKey          = "oauth2client_id"
	oauth2ClientIDSpanAttachmentKey                  = "client_id"
	webhookIDSpanAttachmentKey                       = "webhook_id"
	requestURISpanAttachmentKey                      = "request_uri"
)

func attachUint64ToSpan(span *trace.Span, attachmentKey string, id uint64) {
	if span != nil {
		span.AddAttributes(trace.StringAttribute(attachmentKey, strconv.FormatUint(id, 10)))
	}
}

func attachStringToSpan(span *trace.Span, key, str string) {
	if span != nil {
		span.AddAttributes(trace.StringAttribute(key, str))
	}
}

// AttachFilterToSpan provides a consistent way to attach a filter's info to a span.
func AttachFilterToSpan(span *trace.Span, filter *models.QueryFilter) {
	if filter != nil && span != nil {
		span.AddAttributes(
			trace.StringAttribute(filterPageSpanAttachmentKey, strconv.FormatUint(filter.QueryPage(), 10)),
			trace.StringAttribute(filterLimitSpanAttachmentKey, strconv.FormatUint(uint64(filter.Limit), 10)),
		)
	}
}

// AttachValidInstrumentIDToSpan attaches a valid instrument ID to a given span.
func AttachValidInstrumentIDToSpan(span *trace.Span, validInstrumentID uint64) {
	attachUint64ToSpan(span, validInstrumentIDSpanAttachmentKey, validInstrumentID)
}

// AttachValidIngredientIDToSpan attaches a valid ingredient ID to a given span.
func AttachValidIngredientIDToSpan(span *trace.Span, validIngredientID uint64) {
	attachUint64ToSpan(span, validIngredientIDSpanAttachmentKey, validIngredientID)
}

// AttachValidPreparationIDToSpan attaches a valid preparation ID to a given span.
func AttachValidPreparationIDToSpan(span *trace.Span, validPreparationID uint64) {
	attachUint64ToSpan(span, validPreparationIDSpanAttachmentKey, validPreparationID)
}

// AttachValidIngredientPreparationIDToSpan attaches a valid ingredient preparation ID to a given span.
func AttachValidIngredientPreparationIDToSpan(span *trace.Span, validIngredientPreparationID uint64) {
	attachUint64ToSpan(span, validIngredientPreparationIDSpanAttachmentKey, validIngredientPreparationID)
}

// AttachRequiredPreparationInstrumentIDToSpan attaches a required preparation instrument ID to a given span.
func AttachRequiredPreparationInstrumentIDToSpan(span *trace.Span, requiredPreparationInstrumentID uint64) {
	attachUint64ToSpan(span, requiredPreparationInstrumentIDSpanAttachmentKey, requiredPreparationInstrumentID)
}

// AttachRecipeIDToSpan attaches a recipe ID to a given span.
func AttachRecipeIDToSpan(span *trace.Span, recipeID uint64) {
	attachUint64ToSpan(span, recipeIDSpanAttachmentKey, recipeID)
}

// AttachRecipeStepIDToSpan attaches a recipe step ID to a given span.
func AttachRecipeStepIDToSpan(span *trace.Span, recipeStepID uint64) {
	attachUint64ToSpan(span, recipeStepIDSpanAttachmentKey, recipeStepID)
}

// AttachRecipeStepInstrumentIDToSpan attaches a recipe step instrument ID to a given span.
func AttachRecipeStepInstrumentIDToSpan(span *trace.Span, recipeStepInstrumentID uint64) {
	attachUint64ToSpan(span, recipeStepInstrumentIDSpanAttachmentKey, recipeStepInstrumentID)
}

// AttachRecipeStepIngredientIDToSpan attaches a recipe step ingredient ID to a given span.
func AttachRecipeStepIngredientIDToSpan(span *trace.Span, recipeStepIngredientID uint64) {
	attachUint64ToSpan(span, recipeStepIngredientIDSpanAttachmentKey, recipeStepIngredientID)
}

// AttachRecipeStepProductIDToSpan attaches a recipe step product ID to a given span.
func AttachRecipeStepProductIDToSpan(span *trace.Span, recipeStepProductID uint64) {
	attachUint64ToSpan(span, recipeStepProductIDSpanAttachmentKey, recipeStepProductID)
}

// AttachRecipeIterationIDToSpan attaches a recipe iteration ID to a given span.
func AttachRecipeIterationIDToSpan(span *trace.Span, recipeIterationID uint64) {
	attachUint64ToSpan(span, recipeIterationIDSpanAttachmentKey, recipeIterationID)
}

// AttachRecipeStepEventIDToSpan attaches a recipe step event ID to a given span.
func AttachRecipeStepEventIDToSpan(span *trace.Span, recipeStepEventID uint64) {
	attachUint64ToSpan(span, recipeStepEventIDSpanAttachmentKey, recipeStepEventID)
}

// AttachIterationMediaIDToSpan attaches an iteration media ID to a given span.
func AttachIterationMediaIDToSpan(span *trace.Span, iterationMediaID uint64) {
	attachUint64ToSpan(span, iterationMediaIDSpanAttachmentKey, iterationMediaID)
}

// AttachInvitationIDToSpan attaches an invitation ID to a given span.
func AttachInvitationIDToSpan(span *trace.Span, invitationID uint64) {
	attachUint64ToSpan(span, invitationIDSpanAttachmentKey, invitationID)
}

// AttachReportIDToSpan attaches a report ID to a given span.
func AttachReportIDToSpan(span *trace.Span, reportID uint64) {
	attachUint64ToSpan(span, reportIDSpanAttachmentKey, reportID)
}

// AttachUserIDToSpan provides a consistent way to attach a user's ID to a span.
func AttachUserIDToSpan(span *trace.Span, userID uint64) {
	attachUint64ToSpan(span, userIDSpanAttachmentKey, userID)
}

// AttachOAuth2ClientDatabaseIDToSpan is a consistent way to attach an oauth2 client's ID to a span.
func AttachOAuth2ClientDatabaseIDToSpan(span *trace.Span, oauth2ClientID uint64) {
	attachUint64ToSpan(span, oauth2ClientDatabaseIDSpanAttachmentKey, oauth2ClientID)
}

// AttachOAuth2ClientIDToSpan is a consistent way to attach an oauth2 client's Client ID to a span.
func AttachOAuth2ClientIDToSpan(span *trace.Span, clientID string) {
	attachStringToSpan(span, oauth2ClientIDSpanAttachmentKey, clientID)
}

// AttachUsernameToSpan provides a consistent way to attach a user's username to a span.
func AttachUsernameToSpan(span *trace.Span, username string) {
	attachStringToSpan(span, usernameSpanAttachmentKey, username)
}

// AttachWebhookIDToSpan provides a consistent way to attach a webhook's ID to a span.
func AttachWebhookIDToSpan(span *trace.Span, webhookID uint64) {
	attachUint64ToSpan(span, webhookIDSpanAttachmentKey, webhookID)
}

// AttachRequestURIToSpan attaches a given URI to a span.
func AttachRequestURIToSpan(span *trace.Span, uri string) {
	attachStringToSpan(span, requestURISpanAttachmentKey, uri)
}
