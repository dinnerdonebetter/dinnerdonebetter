package recipes

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/segmentio/ksuid"

	"github.com/prixfixeco/api_server/internal/observability"
	"github.com/prixfixeco/api_server/internal/observability/keys"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/pkg/types"
)

const (
	// RecipeIDURIParamKey is a standard string that we'll use to refer to recipe IDs with.
	RecipeIDURIParamKey = "recipeID"
)

// CreateHandler is our recipe creation route.
func (s *service) CreateHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	logger := s.logger.WithRequest(req)
	tracing.AttachRequestToSpan(span, req)

	// determine user ID.
	sessionCtxData, err := s.sessionContextDataFetcher(req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving session context data")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, "unauthenticated", http.StatusUnauthorized)
		return
	}

	tracing.AttachSessionContextDataToSpan(span, sessionCtxData)
	logger = sessionCtxData.AttachToLogger(logger)

	// read parsed input struct from request body.
	providedInput := new(types.RecipeCreationRequestInput)
	if err = s.encoderDecoder.DecodeRequest(ctx, req, providedInput); err != nil {
		observability.AcknowledgeError(err, logger, span, "decoding request")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, "invalid request content", http.StatusBadRequest)
		return
	}

	if err = providedInput.ValidateWithContext(ctx); err != nil {
		logger.WithValue(keys.ValidationErrorKey, err).Debug("provided input was invalid")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, err.Error(), http.StatusBadRequest)
		return
	}

	input := types.RecipeDatabaseCreationInputFromRecipeCreationInput(providedInput)
	input.ID = ksuid.New().String()

	for i, step := range input.Steps {
		input.Steps[i].ID = ksuid.New().String()
		for j := range step.Ingredients {
			input.Steps[i].Ingredients[j].ID = ksuid.New().String()
		}
	}

	input.CreatedByUser = sessionCtxData.Requester.UserID
	tracing.AttachRecipeIDToSpan(span, input.ID)

	recipe, err := s.recipeDataManager.CreateRecipe(ctx, input)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "creating recipe")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	if s.dataChangesPublisher != nil {
		dcm := &types.DataChangeMessage{
			DataType:                  types.RecipeDataType,
			MessageType:               "recipe_created",
			Recipe:                    recipe,
			AttributableToUserID:      sessionCtxData.Requester.UserID,
			AttributableToHouseholdID: sessionCtxData.ActiveHouseholdID,
		}

		logger.Debug("publishing to data change")
		if err = s.dataChangesPublisher.Publish(ctx, dcm); err != nil {
			observability.AcknowledgeError(err, logger, span, "publishing to data changes topic")
		}
		logger.Debug("published to data change")
	}

	if err = s.customerDataCollector.EventOccurred(ctx, "recipe_created", sessionCtxData.Requester.UserID, map[string]interface{}{
		keys.HouseholdIDKey: sessionCtxData.ActiveHouseholdID,
		keys.RecipeIDKey:    input.ID,
	}); err != nil {
		logger.Error(err, "notifying customer data platform")
	}

	s.encoderDecoder.EncodeResponseWithStatus(ctx, res, recipe, http.StatusCreated)
}

// ReadHandler returns a GET handler that returns a recipe.
func (s *service) ReadHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	logger := s.logger.WithRequest(req)
	tracing.AttachRequestToSpan(span, req)

	// determine user ID.
	sessionCtxData, err := s.sessionContextDataFetcher(req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving session context data")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, "unauthenticated", http.StatusUnauthorized)
		return
	}

	tracing.AttachSessionContextDataToSpan(span, sessionCtxData)
	logger = sessionCtxData.AttachToLogger(logger)

	// determine recipe ID.
	recipeID := s.recipeIDFetcher(req)
	tracing.AttachRecipeIDToSpan(span, recipeID)
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)

	// fetch recipe from database.
	x, err := s.recipeDataManager.GetRecipe(ctx, recipeID)
	if errors.Is(err, sql.ErrNoRows) {
		s.encoderDecoder.EncodeNotFoundResponse(ctx, res)
		return
	} else if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving recipe")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	// encode our response and peace.
	s.encoderDecoder.RespondWithData(ctx, res, x)
}

// ListHandler is our list route.
func (s *service) ListHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	filter := types.ExtractQueryFilter(req)
	logger := s.logger.WithRequest(req).
		WithValue(keys.FilterLimitKey, filter.Limit).
		WithValue(keys.FilterPageKey, filter.Page).
		WithValue(keys.FilterSortByKey, string(filter.SortBy))

	tracing.AttachRequestToSpan(span, req)
	tracing.AttachFilterToSpan(span, filter.Page, filter.Limit, string(filter.SortBy))

	// determine user ID.
	sessionCtxData, err := s.sessionContextDataFetcher(req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving session context data")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, "unauthenticated", http.StatusUnauthorized)
		return
	}

	tracing.AttachSessionContextDataToSpan(span, sessionCtxData)
	logger = sessionCtxData.AttachToLogger(logger)

	recipes, err := s.recipeDataManager.GetRecipes(ctx, filter)
	if errors.Is(err, sql.ErrNoRows) {
		// in the event no rows exist, return an empty list.
		recipes = &types.RecipeList{Recipes: []*types.Recipe{}}
	} else if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving recipes")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	// encode our response and peace.
	s.encoderDecoder.RespondWithData(ctx, res, recipes)
}

// UpdateHandler returns a handler that updates a recipe.
func (s *service) UpdateHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	logger := s.logger.WithRequest(req)
	tracing.AttachRequestToSpan(span, req)

	// determine user ID.
	sessionCtxData, err := s.sessionContextDataFetcher(req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving session context data")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, "unauthenticated", http.StatusUnauthorized)
		return
	}

	tracing.AttachSessionContextDataToSpan(span, sessionCtxData)
	logger = sessionCtxData.AttachToLogger(logger)

	// check for parsed input attached to session context data.
	input := new(types.RecipeUpdateRequestInput)
	if err = s.encoderDecoder.DecodeRequest(ctx, req, input); err != nil {
		logger.Error(err, "error encountered decoding request body")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, "invalid request content", http.StatusBadRequest)
		return
	}

	if err = input.ValidateWithContext(ctx); err != nil {
		logger.Error(err, "provided input was invalid")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, err.Error(), http.StatusBadRequest)
		return
	}

	// determine recipe ID.
	recipeID := s.recipeIDFetcher(req)
	tracing.AttachRecipeIDToSpan(span, recipeID)
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)

	// fetch recipe from database.
	recipe, err := s.recipeDataManager.GetRecipeByIDAndUser(ctx, recipeID, sessionCtxData.Requester.UserID)
	if errors.Is(err, sql.ErrNoRows) {
		s.encoderDecoder.EncodeNotFoundResponse(ctx, res)
		return
	} else if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving recipe for update")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	// update the recipe.
	recipe.Update(input)

	if err = s.recipeDataManager.UpdateRecipe(ctx, recipe); err != nil {
		observability.AcknowledgeError(err, logger, span, "updating recipe")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	if s.dataChangesPublisher != nil {
		dcm := &types.DataChangeMessage{
			DataType:                  types.RecipeDataType,
			MessageType:               "recipe_updated",
			Recipe:                    recipe,
			AttributableToUserID:      sessionCtxData.Requester.UserID,
			AttributableToHouseholdID: sessionCtxData.ActiveHouseholdID,
		}

		if err = s.dataChangesPublisher.Publish(ctx, dcm); err != nil {
			observability.AcknowledgeError(err, logger, span, "publishing data change message")
		}
	}

	if err = s.customerDataCollector.EventOccurred(ctx, "recipe_updated", sessionCtxData.Requester.UserID, map[string]interface{}{
		keys.HouseholdIDKey: sessionCtxData.ActiveHouseholdID,
		keys.RecipeIDKey:    recipeID,
	}); err != nil {
		logger.Error(err, "notifying customer data platform")
	}

	// encode our response and peace.
	s.encoderDecoder.RespondWithData(ctx, res, recipe)
}

// ArchiveHandler returns a handler that archives a recipe.
func (s *service) ArchiveHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	logger := s.logger.WithRequest(req)
	tracing.AttachRequestToSpan(span, req)

	// determine user ID.
	sessionCtxData, err := s.sessionContextDataFetcher(req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving session context data")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, "unauthenticated", http.StatusUnauthorized)
		return
	}

	tracing.AttachSessionContextDataToSpan(span, sessionCtxData)
	logger = sessionCtxData.AttachToLogger(logger)

	// determine recipe ID.
	recipeID := s.recipeIDFetcher(req)
	tracing.AttachRecipeIDToSpan(span, recipeID)
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)

	exists, existenceCheckErr := s.recipeDataManager.RecipeExists(ctx, recipeID)
	if existenceCheckErr != nil && !errors.Is(existenceCheckErr, sql.ErrNoRows) {
		observability.AcknowledgeError(existenceCheckErr, logger, span, "checking recipe existence")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	} else if !exists || errors.Is(existenceCheckErr, sql.ErrNoRows) {
		s.encoderDecoder.EncodeNotFoundResponse(ctx, res)
		return
	}

	if err = s.recipeDataManager.ArchiveRecipe(ctx, recipeID, sessionCtxData.Requester.UserID); err != nil {
		observability.AcknowledgeError(err, logger, span, "archiving recipe")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	if s.dataChangesPublisher != nil {
		dcm := &types.DataChangeMessage{
			DataType:                  types.RecipeDataType,
			MessageType:               "recipe_archived",
			AttributableToUserID:      sessionCtxData.Requester.UserID,
			AttributableToHouseholdID: sessionCtxData.ActiveHouseholdID,
		}

		if err = s.dataChangesPublisher.Publish(ctx, dcm); err != nil {
			observability.AcknowledgeError(err, logger, span, "publishing data change message")
		}
	}

	if err = s.customerDataCollector.EventOccurred(ctx, "recipe_archived", sessionCtxData.Requester.UserID, map[string]interface{}{
		keys.HouseholdIDKey: sessionCtxData.ActiveHouseholdID,
		keys.RecipeIDKey:    recipeID,
	}); err != nil {
		logger.Error(err, "notifying customer data platform")
	}

	// encode our response and peace.
	res.WriteHeader(http.StatusNoContent)
}
