package recipestepingredients

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/prixfixeco/backend/internal/identifiers"
	"github.com/prixfixeco/backend/internal/observability"
	"github.com/prixfixeco/backend/internal/observability/keys"
	"github.com/prixfixeco/backend/internal/observability/tracing"
	"github.com/prixfixeco/backend/pkg/types"
	"github.com/prixfixeco/backend/pkg/types/converters"
)

const (
	// RecipeStepCompletionConditionIDURIParamKey is a standard string that we'll use to refer to recipe step completion condition IDs with.
	RecipeStepCompletionConditionIDURIParamKey = "recipeStepCompletionConditionID"
)

// CreateHandler is our recipe step completion condition creation route.
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
	providedInput := new(types.RecipeStepCompletionConditionForExistingRecipeCreationRequestInput)
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

	input := converters.ConvertRecipeStepCompletionConditionForExistingRecipeCreationRequestInputToRecipeStepCompletionConditionDatabaseCreationInput(providedInput)
	input.ID = identifiers.New()

	// determine recipe ID.
	recipeID := s.recipeIDFetcher(req)
	tracing.AttachRecipeIDToSpan(span, recipeID)
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)

	// determine recipe step ID.
	recipeStepID := s.recipeStepIDFetcher(req)
	tracing.AttachRecipeStepIDToSpan(span, recipeStepID)
	logger = logger.WithValue(keys.RecipeStepIDKey, recipeStepID)

	input.BelongsToRecipeStep = recipeStepID
	tracing.AttachRecipeStepCompletionConditionIDToSpan(span, input.ID)

	recipeStepCompletionCondition, err := s.recipeStepCompletionConditionDataManager.CreateRecipeStepCompletionCondition(ctx, input)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "creating recipe step completion conditions")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	if s.dataChangesPublisher != nil {
		dcm := &types.DataChangeMessage{
			DataType:                      types.RecipeStepCompletionConditionDataType,
			EventType:                     types.RecipeStepCompletionConditionCreatedCustomerEventType,
			RecipeStepCompletionCondition: recipeStepCompletionCondition,
			AttributableToUserID:          sessionCtxData.Requester.UserID,
			AttributableToHouseholdID:     sessionCtxData.ActiveHouseholdID,
		}

		if err = s.dataChangesPublisher.Publish(ctx, dcm); err != nil {
			observability.AcknowledgeError(err, logger, span, "publishing to data changes topic")
		}
	}

	s.encoderDecoder.EncodeResponseWithStatus(ctx, res, recipeStepCompletionCondition, http.StatusCreated)
}

// ReadHandler returns a GET handler that returns a recipe step completion condition.
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

	// determine recipe step ID.
	recipeStepID := s.recipeStepIDFetcher(req)
	tracing.AttachRecipeStepIDToSpan(span, recipeStepID)
	logger = logger.WithValue(keys.RecipeStepIDKey, recipeStepID)

	// determine recipe step completion condition ID.
	recipeStepCompletionConditionID := s.recipeStepCompletionConditionIDFetcher(req)
	tracing.AttachRecipeStepCompletionConditionIDToSpan(span, recipeStepCompletionConditionID)
	logger = logger.WithValue(keys.RecipeStepCompletionConditionIDKey, recipeStepCompletionConditionID)

	// fetch recipe step completion condition from database.
	x, err := s.recipeStepCompletionConditionDataManager.GetRecipeStepCompletionCondition(ctx, recipeID, recipeStepID, recipeStepCompletionConditionID)
	if errors.Is(err, sql.ErrNoRows) {
		s.encoderDecoder.EncodeNotFoundResponse(ctx, res)
		return
	} else if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving recipe step completion condition")
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

	filter := types.ExtractQueryFilterFromRequest(req)
	logger := s.logger.WithRequest(req).
		WithValue(keys.FilterLimitKey, filter.Limit).
		WithValue(keys.FilterPageKey, filter.Page).
		WithValue(keys.FilterSortByKey, filter.SortBy)

	tracing.AttachRequestToSpan(span, req)
	tracing.AttachFilterDataToSpan(span, filter.Page, filter.Limit, filter.SortBy)

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

	// determine recipe step ID.
	recipeStepID := s.recipeStepIDFetcher(req)
	tracing.AttachRecipeStepIDToSpan(span, recipeStepID)
	logger = logger.WithValue(keys.RecipeStepIDKey, recipeStepID)

	recipeStepCompletionConditions, err := s.recipeStepCompletionConditionDataManager.GetRecipeStepCompletionConditions(ctx, recipeID, recipeStepID, filter)
	if errors.Is(err, sql.ErrNoRows) {
		// in the event no rows exist, return an empty list.
		recipeStepCompletionConditions = &types.QueryFilteredResult[types.RecipeStepCompletionCondition]{Data: []*types.RecipeStepCompletionCondition{}}
	} else if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving recipe step completion conditions")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	// encode our response and peace.
	s.encoderDecoder.RespondWithData(ctx, res, recipeStepCompletionConditions)
}

// UpdateHandler returns a handler that updates a recipe step completion condition.
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
	input := new(types.RecipeStepCompletionConditionUpdateRequestInput)
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

	// determine recipe step ID.
	recipeStepID := s.recipeStepIDFetcher(req)
	tracing.AttachRecipeStepIDToSpan(span, recipeStepID)
	logger = logger.WithValue(keys.RecipeStepIDKey, recipeStepID)

	// determine recipe step completion condition ID.
	recipeStepCompletionConditionID := s.recipeStepCompletionConditionIDFetcher(req)
	tracing.AttachRecipeStepCompletionConditionIDToSpan(span, recipeStepCompletionConditionID)
	logger = logger.WithValue(keys.RecipeStepCompletionConditionIDKey, recipeStepCompletionConditionID)

	// fetch recipe step completion condition from database.
	recipeStepCompletionCondition, err := s.recipeStepCompletionConditionDataManager.GetRecipeStepCompletionCondition(ctx, recipeID, recipeStepID, recipeStepCompletionConditionID)
	if errors.Is(err, sql.ErrNoRows) {
		s.encoderDecoder.EncodeNotFoundResponse(ctx, res)
		return
	} else if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving recipe step completion condition for update")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	// update the recipe step completion condition.
	recipeStepCompletionCondition.Update(input)

	if err = s.recipeStepCompletionConditionDataManager.UpdateRecipeStepCompletionCondition(ctx, recipeStepCompletionCondition); err != nil {
		observability.AcknowledgeError(err, logger, span, "updating recipe step completion conditions")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	if s.dataChangesPublisher != nil {
		dcm := &types.DataChangeMessage{
			DataType:                      types.RecipeStepCompletionConditionDataType,
			EventType:                     types.RecipeStepCompletionConditionUpdatedCustomerEventType,
			RecipeStepCompletionCondition: recipeStepCompletionCondition,
			AttributableToUserID:          sessionCtxData.Requester.UserID,
			AttributableToHouseholdID:     sessionCtxData.ActiveHouseholdID,
		}

		if err = s.dataChangesPublisher.Publish(ctx, dcm); err != nil {
			observability.AcknowledgeError(err, logger, span, "publishing data change message")
		}
	}

	// encode our response and peace.
	s.encoderDecoder.RespondWithData(ctx, res, recipeStepCompletionCondition)
}

// ArchiveHandler returns a handler that archives a recipe step completion condition.
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

	// determine recipe step ID.
	recipeStepID := s.recipeStepIDFetcher(req)
	tracing.AttachRecipeStepIDToSpan(span, recipeStepID)
	logger = logger.WithValue(keys.RecipeStepIDKey, recipeStepID)

	// determine recipe step completion condition ID.
	recipeStepCompletionConditionID := s.recipeStepCompletionConditionIDFetcher(req)
	tracing.AttachRecipeStepCompletionConditionIDToSpan(span, recipeStepCompletionConditionID)
	logger = logger.WithValue(keys.RecipeStepCompletionConditionIDKey, recipeStepCompletionConditionID)

	exists, existenceCheckErr := s.recipeStepCompletionConditionDataManager.RecipeStepCompletionConditionExists(ctx, recipeID, recipeStepID, recipeStepCompletionConditionID)
	if existenceCheckErr != nil && !errors.Is(existenceCheckErr, sql.ErrNoRows) {
		observability.AcknowledgeError(existenceCheckErr, logger, span, "checking recipe step completion condition existence")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	} else if !exists || errors.Is(existenceCheckErr, sql.ErrNoRows) {
		s.encoderDecoder.EncodeNotFoundResponse(ctx, res)
		return
	}

	if err = s.recipeStepCompletionConditionDataManager.ArchiveRecipeStepCompletionCondition(ctx, recipeStepID, recipeStepCompletionConditionID); err != nil {
		observability.AcknowledgeError(err, logger, span, "archiving recipe step completion conditions")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	if s.dataChangesPublisher != nil {
		dcm := &types.DataChangeMessage{
			DataType:                  types.RecipeStepCompletionConditionDataType,
			EventType:                 types.RecipeStepCompletionConditionArchivedCustomerEventType,
			AttributableToUserID:      sessionCtxData.Requester.UserID,
			AttributableToHouseholdID: sessionCtxData.ActiveHouseholdID,
		}

		if err = s.dataChangesPublisher.Publish(ctx, dcm); err != nil {
			observability.AcknowledgeError(err, logger, span, "publishing data change message")
		}
	}

	// encode our response and peace.
	res.WriteHeader(http.StatusNoContent)
}
