package mealplanning

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/dinnerdonebetter/backend/internal/lib/database/filtering"
	"github.com/dinnerdonebetter/backend/internal/lib/identifiers"
	"github.com/dinnerdonebetter/backend/internal/lib/observability"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/lib/pointer"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"

	servertiming "github.com/mitchellh/go-server-timing"
)

// CreateMealPlanGroceryListItemHandler is our meal plan grocery list item creation route.
func (s *service) CreateMealPlanGroceryListItemHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	timing := servertiming.FromContext(ctx)
	logger := s.logger.WithRequest(req).WithSpan(span)
	tracing.AttachRequestToSpan(span, req)

	responseDetails := types.ResponseDetails{
		TraceID: span.SpanContext().TraceID().String(),
	}

	// determine user ID.
	sessionContextTimer := timing.NewMetric("session").WithDesc("fetch session context").Start()
	sessionCtxData, err := s.sessionContextDataFetcher(req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving session context data")
		errRes := types.NewAPIErrorResponse("unauthenticated", types.ErrFetchingSessionContextData, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusUnauthorized)
		return
	}
	sessionContextTimer.Stop()

	tracing.AttachSessionContextDataToSpan(span, sessionCtxData)
	logger = sessionCtxData.AttachToLogger(logger)
	responseDetails.CurrentHouseholdID = sessionCtxData.ActiveHouseholdID

	// read parsed input struct from request body.
	providedInput := new(types.MealPlanGroceryListItemCreationRequestInput)
	if err = s.encoderDecoder.DecodeRequest(ctx, req, providedInput); err != nil {
		observability.AcknowledgeError(err, logger, span, "decoding request")
		errRes := types.NewAPIErrorResponse("invalid request content", types.ErrDecodingRequestInput, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusBadRequest)
		return
	}

	if err = providedInput.ValidateWithContext(ctx); err != nil {
		logger.WithValue(keys.ValidationErrorKey, err).Debug("provided input was invalid")
		errRes := types.NewAPIErrorResponse(err.Error(), types.ErrValidatingRequestInput, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusBadRequest)
		return
	}

	input := converters.ConvertMealPlanGroceryListItemCreationRequestInputToMealPlanGroceryListItemDatabaseCreationInput(providedInput)
	input.ID = identifiers.New()
	tracing.AttachToSpan(span, keys.MealPlanGroceryListItemIDKey, input.ID)

	logger = logger.WithValue("input", input)

	// determine meal plan ID.
	mealPlanID := s.mealPlanIDFetcher(req)
	tracing.AttachToSpan(span, keys.MealPlanIDKey, mealPlanID)
	logger = logger.WithValue(keys.MealPlanIDKey, mealPlanID)
	input.BelongsToMealPlan = mealPlanID

	createTimer := timing.NewMetric("database").WithDesc("create").Start()
	mealPlanGroceryListItem, err := s.mealPlanningDataManager.CreateMealPlanGroceryListItem(ctx, input)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "creating meal plan")
		errRes := types.NewAPIErrorResponse("database error", types.ErrTalkingToDatabase, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusInternalServerError)
		return
	}
	createTimer.Stop()

	dcm := &types.DataChangeMessage{
		EventType:               types.MealPlanGroceryListItemCreatedServiceEventType,
		MealPlanID:              mealPlanID,
		MealPlanGroceryListItem: mealPlanGroceryListItem,
		HouseholdID:             sessionCtxData.ActiveHouseholdID,
		UserID:                  sessionCtxData.Requester.UserID,
	}

	s.dataChangesPublisher.PublishAsync(ctx, dcm)

	responseValue := &types.APIResponse[*types.MealPlanGroceryListItem]{
		Details: responseDetails,
		Data:    mealPlanGroceryListItem,
	}

	s.encoderDecoder.EncodeResponseWithStatus(ctx, res, responseValue, http.StatusCreated)
}

// ReadMealPlanGroceryListItemHandler returns a GET handler that returns a meal plan grocery list item.
func (s *service) ReadMealPlanGroceryListItemHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	timing := servertiming.FromContext(ctx)
	logger := s.logger.WithRequest(req).WithSpan(span)
	tracing.AttachRequestToSpan(span, req)

	responseDetails := types.ResponseDetails{
		TraceID: span.SpanContext().TraceID().String(),
	}

	// determine user ID.
	sessionContextTimer := timing.NewMetric("session").WithDesc("fetch session context").Start()
	sessionCtxData, err := s.sessionContextDataFetcher(req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving session context data")
		errRes := types.NewAPIErrorResponse("unauthenticated", types.ErrFetchingSessionContextData, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusUnauthorized)
		return
	}
	sessionContextTimer.Stop()

	tracing.AttachSessionContextDataToSpan(span, sessionCtxData)
	logger = sessionCtxData.AttachToLogger(logger)
	responseDetails.CurrentHouseholdID = sessionCtxData.ActiveHouseholdID

	// determine meal plan ID.
	mealPlanID := s.mealPlanIDFetcher(req)
	tracing.AttachToSpan(span, keys.MealPlanIDKey, mealPlanID)
	logger = logger.WithValue(keys.MealPlanIDKey, mealPlanID)

	// determine meal plan grocery list item ID.
	mealPlanGroceryListItemID := s.mealPlanGroceryListItemIDFetcher(req)
	tracing.AttachToSpan(span, keys.MealPlanEventIDKey, mealPlanGroceryListItemID)
	logger = logger.WithValue(keys.MealPlanGroceryListItemIDKey, mealPlanGroceryListItemID)

	// fetch meal plan grocery list item from database.
	x, err := s.mealPlanningDataManager.GetMealPlanGroceryListItem(ctx, mealPlanID, mealPlanGroceryListItemID)
	if errors.Is(err, sql.ErrNoRows) {
		errRes := types.NewAPIErrorResponse("not found", types.ErrDataNotFound, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusNotFound)
		return
	} else if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving meal plan grocery list item")
		errRes := types.NewAPIErrorResponse("database error", types.ErrTalkingToDatabase, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusInternalServerError)
		return
	}

	responseValue := &types.APIResponse[*types.MealPlanGroceryListItem]{
		Details: responseDetails,
		Data:    x,
	}

	// encode our response and peace.
	s.encoderDecoder.EncodeResponseWithStatus(ctx, res, responseValue, http.StatusOK)
}

// ListMealPlanGroceryListItemsByMealPlanHandler is our list route.
func (s *service) ListMealPlanGroceryListItemsByMealPlanHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	timing := servertiming.FromContext(ctx)
	logger := s.logger.WithRequest(req).WithSpan(span)
	tracing.AttachRequestToSpan(span, req)
	filter := filtering.ExtractQueryFilterFromRequest(req)
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	responseDetails := types.ResponseDetails{
		TraceID: span.SpanContext().TraceID().String(),
	}

	// determine user ID.
	sessionContextTimer := timing.NewMetric("session").WithDesc("fetch session context").Start()
	sessionCtxData, err := s.sessionContextDataFetcher(req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving session context data")
		errRes := types.NewAPIErrorResponse("unauthenticated", types.ErrFetchingSessionContextData, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusUnauthorized)
		return
	}
	sessionContextTimer.Stop()

	tracing.AttachSessionContextDataToSpan(span, sessionCtxData)
	logger = sessionCtxData.AttachToLogger(logger)
	responseDetails.CurrentHouseholdID = sessionCtxData.ActiveHouseholdID

	// determine meal plan ID.
	mealPlanID := s.mealPlanIDFetcher(req)
	tracing.AttachToSpan(span, keys.MealPlanIDKey, mealPlanID)
	logger = logger.WithValue(keys.MealPlanIDKey, mealPlanID)

	mealPlanGroceryListItems, err := s.mealPlanningDataManager.GetMealPlanGroceryListItemsForMealPlan(ctx, mealPlanID)
	if errors.Is(err, sql.ErrNoRows) {
		// in the event no rows exist, return an empty list.
		mealPlanGroceryListItems = []*types.MealPlanGroceryListItem{}
	} else if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving meal plan grocery list items for meal plan")
		errRes := types.NewAPIErrorResponse("database error", types.ErrTalkingToDatabase, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusInternalServerError)
		return
	}

	responseValue := &types.APIResponse[[]*types.MealPlanGroceryListItem]{
		Details:    responseDetails,
		Data:       mealPlanGroceryListItems,
		Pagination: pointer.To(filter.ToPagination()),
	}

	// encode our response and peace.
	s.encoderDecoder.EncodeResponseWithStatus(ctx, res, responseValue, http.StatusOK)
}

// UpdateMealPlanGroceryListItemHandler returns a handler that updates a meal plan grocery list item.
func (s *service) UpdateMealPlanGroceryListItemHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	timing := servertiming.FromContext(ctx)
	logger := s.logger.WithRequest(req).WithSpan(span)
	tracing.AttachRequestToSpan(span, req)

	responseDetails := types.ResponseDetails{
		TraceID: span.SpanContext().TraceID().String(),
	}

	// determine user ID.
	sessionContextTimer := timing.NewMetric("session").WithDesc("fetch session context").Start()
	sessionCtxData, err := s.sessionContextDataFetcher(req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving session context data")
		errRes := types.NewAPIErrorResponse("unauthenticated", types.ErrFetchingSessionContextData, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusUnauthorized)
		return
	}
	sessionContextTimer.Stop()

	tracing.AttachSessionContextDataToSpan(span, sessionCtxData)
	logger = sessionCtxData.AttachToLogger(logger)
	responseDetails.CurrentHouseholdID = sessionCtxData.ActiveHouseholdID

	// determine meal plan ID.
	mealPlanID := s.mealPlanIDFetcher(req)
	tracing.AttachToSpan(span, keys.MealPlanIDKey, mealPlanID)
	logger = logger.WithValue(keys.MealPlanIDKey, mealPlanID)

	// determine meal plan grocery list item ID.
	mealPlanGroceryListItemID := s.mealPlanGroceryListItemIDFetcher(req)
	tracing.AttachToSpan(span, keys.MealPlanEventIDKey, mealPlanGroceryListItemID)
	logger = logger.WithValue(keys.MealPlanGroceryListItemIDKey, mealPlanGroceryListItemID)

	// read parsed input struct from request body.
	providedInput := new(types.MealPlanGroceryListItemUpdateRequestInput)
	if err = s.encoderDecoder.DecodeRequest(ctx, req, providedInput); err != nil {
		observability.AcknowledgeError(err, logger, span, "decoding request")
		errRes := types.NewAPIErrorResponse("invalid request content", types.ErrDecodingRequestInput, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusBadRequest)
		return
	}

	if err = providedInput.ValidateWithContext(ctx); err != nil {
		logger.WithValue(keys.ValidationErrorKey, err).Debug("provided input was invalid")
		errRes := types.NewAPIErrorResponse(err.Error(), types.ErrValidatingRequestInput, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusBadRequest)
		return
	}

	mealPlanGroceryListItem, err := s.mealPlanningDataManager.GetMealPlanGroceryListItem(ctx, mealPlanID, mealPlanGroceryListItemID)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "checking meal plan grocery list item existence")
		if errors.Is(err, sql.ErrNoRows) {
			errRes := types.NewAPIErrorResponse("not found", types.ErrDataNotFound, responseDetails)
			s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusNotFound)
			return
		}
		errRes := types.NewAPIErrorResponse("database error", types.ErrTalkingToDatabase, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusInternalServerError)
		return
	}

	mealPlanGroceryListItem.Update(providedInput)

	updateTimer := timing.NewMetric("database").WithDesc("update").Start()
	if err = s.mealPlanningDataManager.UpdateMealPlanGroceryListItem(ctx, mealPlanGroceryListItem); err != nil {
		observability.AcknowledgeError(err, logger, span, "archiving meal plan grocery list item")
		errRes := types.NewAPIErrorResponse("database error", types.ErrTalkingToDatabase, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusInternalServerError)
		return
	}
	updateTimer.Stop()

	dcm := &types.DataChangeMessage{
		EventType:                 types.MealPlanGroceryListItemUpdatedServiceEventType,
		MealPlanGroceryListItem:   mealPlanGroceryListItem,
		MealPlanGroceryListItemID: mealPlanGroceryListItemID,
		HouseholdID:               sessionCtxData.ActiveHouseholdID,
		UserID:                    sessionCtxData.Requester.UserID,
	}

	s.dataChangesPublisher.PublishAsync(ctx, dcm)

	responseValue := &types.APIResponse[*types.MealPlanGroceryListItem]{
		Details: responseDetails,
		Data:    mealPlanGroceryListItem,
	}

	s.encoderDecoder.EncodeResponseWithStatus(ctx, res, responseValue, http.StatusOK)
}

// ArchiveMealPlanGroceryListItemHandler returns a GET handler that returns a meal plan grocery list item.
func (s *service) ArchiveMealPlanGroceryListItemHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	timing := servertiming.FromContext(ctx)
	logger := s.logger.WithRequest(req).WithSpan(span)
	tracing.AttachRequestToSpan(span, req)

	responseDetails := types.ResponseDetails{
		TraceID: span.SpanContext().TraceID().String(),
	}

	// determine user ID.
	sessionContextTimer := timing.NewMetric("session").WithDesc("fetch session context").Start()
	sessionCtxData, err := s.sessionContextDataFetcher(req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving session context data")
		errRes := types.NewAPIErrorResponse("unauthenticated", types.ErrFetchingSessionContextData, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusUnauthorized)
		return
	}
	sessionContextTimer.Stop()

	tracing.AttachSessionContextDataToSpan(span, sessionCtxData)
	logger = sessionCtxData.AttachToLogger(logger)
	responseDetails.CurrentHouseholdID = sessionCtxData.ActiveHouseholdID

	// determine meal plan ID.
	mealPlanID := s.mealPlanIDFetcher(req)
	tracing.AttachToSpan(span, keys.MealPlanIDKey, mealPlanID)
	logger = logger.WithValue(keys.MealPlanIDKey, mealPlanID)

	// determine meal plan grocery list item ID.
	mealPlanGroceryListItemID := s.mealPlanGroceryListItemIDFetcher(req)
	tracing.AttachToSpan(span, keys.MealPlanEventIDKey, mealPlanGroceryListItemID)
	logger = logger.WithValue(keys.MealPlanGroceryListItemIDKey, mealPlanGroceryListItemID)

	// check that meal plan grocery list item exists in database.
	existenceTimer := timing.NewMetric("database").WithDesc("existence check").Start()
	x, err := s.mealPlanningDataManager.MealPlanGroceryListItemExists(ctx, mealPlanID, mealPlanGroceryListItemID)
	if errors.Is(err, sql.ErrNoRows) || !x {
		errRes := types.NewAPIErrorResponse("not found", types.ErrDataNotFound, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusNotFound)
		return
	} else if err != nil {
		observability.AcknowledgeError(err, logger, span, "checking meal plan grocery list item existence")
		errRes := types.NewAPIErrorResponse("database error", types.ErrTalkingToDatabase, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusInternalServerError)
		return
	}
	existenceTimer.Stop()

	// fetch meal plan grocery list item from database.
	archiveTimer := timing.NewMetric("database").WithDesc("archive").Start()
	err = s.mealPlanningDataManager.ArchiveMealPlanGroceryListItem(ctx, mealPlanID, mealPlanGroceryListItemID)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "archiving meal plan grocery list item")
		errRes := types.NewAPIErrorResponse("database error", types.ErrTalkingToDatabase, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusInternalServerError)
		return
	}
	archiveTimer.Stop()

	dcm := &types.DataChangeMessage{
		EventType:                 types.MealPlanGroceryListItemArchivedServiceEventType,
		MealPlanGroceryListItemID: mealPlanGroceryListItemID,
		HouseholdID:               sessionCtxData.ActiveHouseholdID,
		UserID:                    sessionCtxData.Requester.UserID,
	}

	s.dataChangesPublisher.PublishAsync(ctx, dcm)

	responseValue := &types.APIResponse[*types.MealPlanGroceryListItem]{
		Details: responseDetails,
	}

	// let everybody go home.
	s.encoderDecoder.EncodeResponseWithStatus(ctx, res, responseValue, http.StatusOK)
}
