package mealplanning

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/pkg/identifiers"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"

	servertiming "github.com/mitchellh/go-server-timing"
)

// CreateMealPlanEventHandler is our meal plan creation route.
func (s *service) CreateMealPlanEventHandler(res http.ResponseWriter, req *http.Request) {
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
	responseDetails.CurrentHouseholdID = sessionCtxData.ActiveHouseholdID // determine meal plan ID.
	mealPlanID := s.mealPlanIDFetcher(req)
	tracing.AttachToSpan(span, keys.MealPlanIDKey, mealPlanID)
	logger = logger.WithValue(keys.MealPlanIDKey, mealPlanID)

	// read parsed input struct from request body.
	providedInput := new(types.MealPlanEventCreationRequestInput)
	if err = s.encoderDecoder.DecodeRequest(ctx, req, providedInput); err != nil {
		observability.AcknowledgeError(err, logger, span, "decoding request")
		errRes := types.NewAPIErrorResponse("invalid request content", types.ErrDecodingRequestInput, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusBadRequest)
		return
	}

	if err = providedInput.ValidateWithContext(ctx); err != nil {
		logger.Error("provided input was invalid", err)
		errRes := types.NewAPIErrorResponse(err.Error(), types.ErrValidatingRequestInput, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusBadRequest)
		return
	}

	input := converters.ConvertMealPlanEventCreationRequestInputToMealPlanEventDatabaseCreationInput(providedInput)
	input.ID = identifiers.New()
	input.BelongsToMealPlan = mealPlanID
	tracing.AttachToSpan(span, keys.MealPlanEventIDKey, input.ID)

	for i := range input.Options {
		input.Options[i].ID = identifiers.New()
		input.Options[i].BelongsToMealPlanEvent = input.ID
	}

	createTimer := timing.NewMetric("database").WithDesc("create").Start()
	mealPlanEvent, err := s.mealPlanningDataManager.CreateMealPlanEvent(ctx, input)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "creating meal plan")
		errRes := types.NewAPIErrorResponse("database error", types.ErrTalkingToDatabase, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusInternalServerError)
		return
	}
	createTimer.Stop()

	dcm := &types.DataChangeMessage{
		EventType:     types.MealPlanEventCreatedServiceEventType,
		MealPlanEvent: mealPlanEvent,
		HouseholdID:   sessionCtxData.ActiveHouseholdID,
		UserID:        sessionCtxData.Requester.UserID,
	}

	go s.dataChangesPublisher.PublishAsync(ctx, dcm)

	responseValue := &types.APIResponse[*types.MealPlanEvent]{
		Details: responseDetails,
		Data:    mealPlanEvent,
	}

	s.encoderDecoder.EncodeResponseWithStatus(ctx, res, responseValue, http.StatusCreated)
}

// ReadMealPlanEventHandler returns a GET handler that returns a meal plan.
func (s *service) ReadMealPlanEventHandler(res http.ResponseWriter, req *http.Request) {
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
	responseDetails.CurrentHouseholdID = sessionCtxData.ActiveHouseholdID // determine meal plan ID.
	mealPlanID := s.mealPlanIDFetcher(req)
	tracing.AttachToSpan(span, keys.MealPlanIDKey, mealPlanID)
	logger = logger.WithValue(keys.MealPlanIDKey, mealPlanID)

	// determine meal plan ID.
	mealPlanEventID := s.mealPlanEventIDFetcher(req)
	tracing.AttachToSpan(span, keys.MealPlanEventIDKey, mealPlanEventID)
	logger = logger.WithValue(keys.MealPlanEventIDKey, mealPlanEventID)

	// fetch meal plan from database.
	readTimer := timing.NewMetric("database").WithDesc("fetch").Start()
	x, err := s.mealPlanningDataManager.GetMealPlanEvent(ctx, mealPlanID, mealPlanEventID)
	if errors.Is(err, sql.ErrNoRows) {
		errRes := types.NewAPIErrorResponse("not found", types.ErrDataNotFound, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusNotFound)
		return
	} else if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving meal plan")
		errRes := types.NewAPIErrorResponse("database error", types.ErrTalkingToDatabase, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusInternalServerError)
		return
	}
	readTimer.Stop()

	responseValue := &types.APIResponse[*types.MealPlanEvent]{
		Details: responseDetails,
		Data:    x,
	}

	// encode our response and peace.
	s.encoderDecoder.RespondWithData(ctx, res, responseValue)
}

// ListMealPlanEventHandler is our list route.
func (s *service) ListMealPlanEventHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	timing := servertiming.FromContext(ctx)
	filter := types.ExtractQueryFilterFromRequest(req)
	logger := s.logger.WithRequest(req).WithSpan(span)
	logger = filter.AttachToLogger(logger)

	responseDetails := types.ResponseDetails{
		TraceID: span.SpanContext().TraceID().String(),
	}

	tracing.AttachRequestToSpan(span, req)
	tracing.AttachQueryFilterToSpan(span, filter)

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
	responseDetails.CurrentHouseholdID = sessionCtxData.ActiveHouseholdID // determine meal plan ID.
	mealPlanID := s.mealPlanIDFetcher(req)
	tracing.AttachToSpan(span, keys.MealPlanIDKey, mealPlanID)
	logger = logger.WithValue(keys.MealPlanIDKey, mealPlanID)

	readTimer := timing.NewMetric("database").WithDesc("fetch").Start()
	mealPlanEvents, err := s.mealPlanningDataManager.GetMealPlanEvents(ctx, mealPlanID, filter)
	if errors.Is(err, sql.ErrNoRows) {
		// in the event no rows exist, return an empty list.
		mealPlanEvents = &types.QueryFilteredResult[types.MealPlanEvent]{Data: []*types.MealPlanEvent{}}
	} else if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving meal plans")
		errRes := types.NewAPIErrorResponse("database error", types.ErrTalkingToDatabase, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusInternalServerError)
		return
	}
	readTimer.Stop()

	responseValue := &types.APIResponse[[]*types.MealPlanEvent]{
		Details:    responseDetails,
		Data:       mealPlanEvents.Data,
		Pagination: &mealPlanEvents.Pagination,
	}

	// encode our response and peace.
	s.encoderDecoder.RespondWithData(ctx, res, responseValue)
}

// UpdateMealPlanEventHandler returns a handler that updates a meal plan.
func (s *service) UpdateMealPlanEventHandler(res http.ResponseWriter, req *http.Request) {
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
	responseDetails.CurrentHouseholdID = sessionCtxData.ActiveHouseholdID // check for parsed input attached to session context data.
	input := new(types.MealPlanEventUpdateRequestInput)
	if err = s.encoderDecoder.DecodeRequest(ctx, req, input); err != nil {
		logger.Error("error encountered decoding request body", err)
		errRes := types.NewAPIErrorResponse("invalid request content", types.ErrDecodingRequestInput, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusBadRequest)
		return
	}

	if err = input.ValidateWithContext(ctx); err != nil {
		logger.Error("provided input was invalid", err)
		errRes := types.NewAPIErrorResponse(err.Error(), types.ErrValidatingRequestInput, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusBadRequest)
		return
	}

	// determine meal plan ID.
	mealPlanID := s.mealPlanIDFetcher(req)
	tracing.AttachToSpan(span, keys.MealPlanIDKey, mealPlanID)
	logger = logger.WithValue(keys.MealPlanIDKey, mealPlanID)

	// determine meal plan ID.
	mealPlanEventID := s.mealPlanEventIDFetcher(req)
	tracing.AttachToSpan(span, keys.MealPlanEventIDKey, mealPlanEventID)
	logger = logger.WithValue(keys.MealPlanEventIDKey, mealPlanEventID)

	// fetch meal plan from database.
	readTimer := timing.NewMetric("database").WithDesc("fetch").Start()
	mealPlanEvent, err := s.mealPlanningDataManager.GetMealPlanEvent(ctx, mealPlanID, mealPlanEventID)
	if errors.Is(err, sql.ErrNoRows) {
		errRes := types.NewAPIErrorResponse("not found", types.ErrDataNotFound, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusNotFound)
		return
	} else if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving meal plan for update")
		errRes := types.NewAPIErrorResponse("database error", types.ErrTalkingToDatabase, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusInternalServerError)
		return
	}
	readTimer.Stop()

	// update the meal plan.
	mealPlanEvent.Update(input)

	updateTimer := timing.NewMetric("database").WithDesc("update").Start()
	if err = s.mealPlanningDataManager.UpdateMealPlanEvent(ctx, mealPlanEvent); err != nil {
		observability.AcknowledgeError(err, logger, span, "updating meal plan")
		errRes := types.NewAPIErrorResponse("database error", types.ErrTalkingToDatabase, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusInternalServerError)
		return
	}
	updateTimer.Stop()

	dcm := &types.DataChangeMessage{
		EventType:     types.MealPlanEventUpdatedServiceEventType,
		MealPlanEvent: mealPlanEvent,
		HouseholdID:   sessionCtxData.ActiveHouseholdID,
		UserID:        sessionCtxData.Requester.UserID,
	}

	go s.dataChangesPublisher.PublishAsync(ctx, dcm)

	responseValue := &types.APIResponse[*types.MealPlanEvent]{
		Details: responseDetails,
		Data:    mealPlanEvent,
	}

	// encode our response and peace.
	s.encoderDecoder.RespondWithData(ctx, res, responseValue)
}

// ArchiveMealPlanEventHandler returns a handler that archives a meal plan.
func (s *service) ArchiveMealPlanEventHandler(res http.ResponseWriter, req *http.Request) {
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

	// determine meal plan ID.
	mealPlanEventID := s.mealPlanEventIDFetcher(req)
	tracing.AttachToSpan(span, keys.MealPlanEventIDKey, mealPlanEventID)
	logger = logger.WithValue(keys.MealPlanEventIDKey, mealPlanEventID)

	existenceTimer := timing.NewMetric("database").WithDesc("existence check").Start()
	exists, err := s.mealPlanningDataManager.MealPlanEventExists(ctx, mealPlanID, mealPlanEventID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		observability.AcknowledgeError(err, logger, span, "checking meal plan existence")
		errRes := types.NewAPIErrorResponse("database error", types.ErrTalkingToDatabase, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusInternalServerError)
		return
	} else if !exists || errors.Is(err, sql.ErrNoRows) {
		errRes := types.NewAPIErrorResponse("not found", types.ErrDataNotFound, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusNotFound)
		return
	}
	existenceTimer.Stop()

	archiveTimer := timing.NewMetric("database").WithDesc("archive").Start()
	if err = s.mealPlanningDataManager.ArchiveMealPlanEvent(ctx, mealPlanID, mealPlanEventID); err != nil {
		observability.AcknowledgeError(err, logger, span, "archiving meal plan")
		errRes := types.NewAPIErrorResponse("database error", types.ErrTalkingToDatabase, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusInternalServerError)
		return
	}
	archiveTimer.Stop()

	dcm := &types.DataChangeMessage{
		EventType:       types.MealPlanEventArchivedServiceEventType,
		MealPlanEventID: mealPlanEventID,
		HouseholdID:     sessionCtxData.ActiveHouseholdID,
		UserID:          sessionCtxData.Requester.UserID,
	}

	go s.dataChangesPublisher.PublishAsync(ctx, dcm)

	responseValue := &types.APIResponse[*types.MealPlanEvent]{
		Details: responseDetails,
	}

	// let everybody go home.
	s.encoderDecoder.RespondWithData(ctx, res, responseValue)
}
