package mealplangrocerylistitems

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
	// MealPlanGroceryListItemIDURIParamKey is a standard string that we'll use to refer to meal plan grocery list item IDs with.
	MealPlanGroceryListItemIDURIParamKey = "mealPlanGroceryListItemID"
)

// CreateHandler is our meal plan grocery list item creation route.
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
	providedInput := new(types.MealPlanGroceryListItemCreationRequestInput)
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

	input := converters.ConvertMealPlanGroceryListItemCreationRequestInputToMealPlanGroceryListItemDatabaseCreationInput(providedInput)
	input.ID = identifiers.New()
	tracing.AttachMealPlanGroceryListItemIDToSpan(span, input.ID)

	logger = logger.WithValue("input", input)

	// determine meal plan ID.
	mealPlanID := s.mealPlanIDFetcher(req)
	tracing.AttachMealPlanIDToSpan(span, mealPlanID)
	logger = logger.WithValue(keys.MealPlanIDKey, mealPlanID)
	input.BelongsToMealPlan = mealPlanID

	mealPlanGroceryListItem, err := s.mealPlanGroceryListItemDataManager.CreateMealPlanGroceryListItem(ctx, input)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "creating meal plan")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	if s.dataChangesPublisher != nil {
		dcm := &types.DataChangeMessage{
			DataType:                  types.MealPlanDataType,
			EventType:                 types.MealPlanCreatedCustomerEventType,
			MealPlanID:                mealPlanID,
			MealPlanGroceryListItem:   mealPlanGroceryListItem,
			AttributableToUserID:      sessionCtxData.Requester.UserID,
			AttributableToHouseholdID: sessionCtxData.ActiveHouseholdID,
		}

		if err = s.dataChangesPublisher.Publish(ctx, dcm); err != nil {
			observability.AcknowledgeError(err, logger, span, "publishing to data changes topic")
		}
	}

	s.encoderDecoder.EncodeResponseWithStatus(ctx, res, mealPlanGroceryListItem, http.StatusCreated)
}

// ReadHandler returns a GET handler that returns a meal plan grocery list item.
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

	// determine meal plan ID.
	mealPlanID := s.mealPlanIDFetcher(req)
	tracing.AttachMealPlanIDToSpan(span, mealPlanID)
	logger = logger.WithValue(keys.MealPlanIDKey, mealPlanID)

	// determine meal plan grocery list item ID.
	mealPlanGroceryListItemID := s.mealPlanGroceryListItemIDFetcher(req)
	tracing.AttachMealPlanEventIDToSpan(span, mealPlanGroceryListItemID)
	logger = logger.WithValue(keys.MealPlanGroceryListItemIDKey, mealPlanGroceryListItemID)

	// fetch meal plan grocery list item from database.
	x, err := s.mealPlanGroceryListItemDataManager.GetMealPlanGroceryListItem(ctx, mealPlanID, mealPlanGroceryListItemID)
	if errors.Is(err, sql.ErrNoRows) {
		s.encoderDecoder.EncodeNotFoundResponse(ctx, res)
		return
	} else if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving meal plan grocery list item")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	// encode our response and peace.
	s.encoderDecoder.RespondWithData(ctx, res, x)
}

// ListByMealPlanHandler is our list route.
func (s *service) ListByMealPlanHandler(res http.ResponseWriter, req *http.Request) {
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

	// determine meal plan ID.
	mealPlanID := s.mealPlanIDFetcher(req)
	tracing.AttachMealPlanIDToSpan(span, mealPlanID)
	logger = logger.WithValue(keys.MealPlanIDKey, mealPlanID)

	mealPlanGroceryListItems, err := s.mealPlanGroceryListItemDataManager.GetMealPlanGroceryListItemsForMealPlan(ctx, mealPlanID)
	if errors.Is(err, sql.ErrNoRows) {
		// in the event no rows exist, return an empty list.
		mealPlanGroceryListItems = []*types.MealPlanGroceryListItem{}
	} else if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving meal plan grocery list items for meal plan")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	// encode our response and peace.
	s.encoderDecoder.RespondWithData(ctx, res, mealPlanGroceryListItems)
}

// UpdateHandler returns a handler that updates a meal plan grocery list item.
func (s *service) UpdateHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	logger := s.logger.WithRequest(req)
	tracing.AttachRequestToSpan(span, req)

	// determine user ID.
	sessionCtxData, sessionCtxFetchErr := s.sessionContextDataFetcher(req)
	if sessionCtxFetchErr != nil {
		observability.AcknowledgeError(sessionCtxFetchErr, logger, span, "retrieving session context data")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, "unauthenticated", http.StatusUnauthorized)
		return
	}

	tracing.AttachSessionContextDataToSpan(span, sessionCtxData)
	logger = sessionCtxData.AttachToLogger(logger)

	// determine meal plan ID.
	mealPlanID := s.mealPlanIDFetcher(req)
	tracing.AttachMealPlanIDToSpan(span, mealPlanID)
	logger = logger.WithValue(keys.MealPlanIDKey, mealPlanID)

	// determine meal plan grocery list item ID.
	mealPlanGroceryListItemID := s.mealPlanGroceryListItemIDFetcher(req)
	tracing.AttachMealPlanEventIDToSpan(span, mealPlanGroceryListItemID)
	logger = logger.WithValue(keys.MealPlanGroceryListItemIDKey, mealPlanGroceryListItemID)

	// read parsed input struct from request body.
	providedInput := new(types.MealPlanGroceryListItemUpdateRequestInput)
	if err := s.encoderDecoder.DecodeRequest(ctx, req, providedInput); err != nil {
		observability.AcknowledgeError(err, logger, span, "decoding request")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, "invalid request content", http.StatusBadRequest)
		return
	}

	if err := providedInput.ValidateWithContext(ctx); err != nil {
		logger.WithValue(keys.ValidationErrorKey, err).Debug("provided input was invalid")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, err.Error(), http.StatusBadRequest)
		return
	}

	mealPlanGroceryListItem, fetchMealPlanGroceryListItemErr := s.mealPlanGroceryListItemDataManager.GetMealPlanGroceryListItem(ctx, mealPlanID, mealPlanGroceryListItemID)
	if fetchMealPlanGroceryListItemErr != nil && !errors.Is(fetchMealPlanGroceryListItemErr, sql.ErrNoRows) {
		observability.AcknowledgeError(fetchMealPlanGroceryListItemErr, logger, span, "checking meal plan grocery list item existence")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	} else if errors.Is(fetchMealPlanGroceryListItemErr, sql.ErrNoRows) {
		s.encoderDecoder.EncodeNotFoundResponse(ctx, res)
		return
	}

	mealPlanGroceryListItem.Update(providedInput)

	if err := s.mealPlanGroceryListItemDataManager.UpdateMealPlanGroceryListItem(ctx, mealPlanGroceryListItem); err != nil {
		observability.AcknowledgeError(err, logger, span, "archiving meal plan grocery list item")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	if s.dataChangesPublisher != nil {
		dcm := &types.DataChangeMessage{
			DataType:                  types.MealPlanGroceryListItemDataType,
			EventType:                 types.MealPlanGroceryListItemUpdatedCustomerEventType,
			MealPlanGroceryListItem:   mealPlanGroceryListItem,
			MealPlanGroceryListItemID: mealPlanGroceryListItemID,
			AttributableToUserID:      sessionCtxData.Requester.UserID,
			AttributableToHouseholdID: sessionCtxData.ActiveHouseholdID,
		}

		if err := s.dataChangesPublisher.Publish(ctx, dcm); err != nil {
			observability.AcknowledgeError(err, logger, span, "publishing data change message")
		}
	}

	s.encoderDecoder.RespondWithData(ctx, res, mealPlanGroceryListItem)
}

// ArchiveHandler returns a GET handler that returns a meal plan grocery list item.
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

	// determine meal plan ID.
	mealPlanID := s.mealPlanIDFetcher(req)
	tracing.AttachMealPlanIDToSpan(span, mealPlanID)
	logger = logger.WithValue(keys.MealPlanIDKey, mealPlanID)

	// determine meal plan grocery list item ID.
	mealPlanGroceryListItemID := s.mealPlanGroceryListItemIDFetcher(req)
	tracing.AttachMealPlanEventIDToSpan(span, mealPlanGroceryListItemID)
	logger = logger.WithValue(keys.MealPlanGroceryListItemIDKey, mealPlanGroceryListItemID)

	// check that meal plan grocery list item exists in database.
	x, err := s.mealPlanGroceryListItemDataManager.MealPlanGroceryListItemExists(ctx, mealPlanID, mealPlanGroceryListItemID)
	if errors.Is(err, sql.ErrNoRows) || !x {
		s.encoderDecoder.EncodeNotFoundResponse(ctx, res)
		return
	} else if err != nil {
		observability.AcknowledgeError(err, logger, span, "checking meal plan grocery list item existence")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	// fetch meal plan grocery list item from database.
	err = s.mealPlanGroceryListItemDataManager.ArchiveMealPlanGroceryListItem(ctx, mealPlanGroceryListItemID)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "archiving meal plan grocery list item")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	if s.dataChangesPublisher != nil {
		dcm := &types.DataChangeMessage{
			DataType:                  types.MealPlanGroceryListItemDataType,
			EventType:                 types.MealPlanGroceryListItemArchivedCustomerEventType,
			MealPlanGroceryListItemID: mealPlanGroceryListItemID,
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
