package mealplantasks

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
	// MealPlanTaskIDURIParamKey is a standard string that we'll use to refer to meal plan task IDs with.
	MealPlanTaskIDURIParamKey = "mealPlanTaskID"
)

// CreateHandler is our meal plan task creation route.
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
	providedInput := new(types.MealPlanTaskCreationRequestInput)
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

	input := converters.ConvertMealPlanTaskCreationRequestInputToMealPlanTaskDatabaseCreationInput(providedInput)
	input.ID = identifiers.New()
	tracing.AttachMealPlanTaskIDToSpan(span, input.ID)

	logger = logger.WithValue("input", input)

	// determine meal plan ID.
	mealPlanID := s.mealPlanIDFetcher(req)
	tracing.AttachMealPlanIDToSpan(span, mealPlanID)
	logger = logger.WithValue(keys.MealPlanIDKey, mealPlanID)

	mealPlanTask, err := s.mealPlanTaskDataManager.CreateMealPlanTask(ctx, input)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "creating meal plan")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	if s.dataChangesPublisher != nil {
		dcm := &types.DataChangeMessage{
			DataType:     types.MealPlanDataType,
			EventType:    types.MealPlanCreatedCustomerEventType,
			MealPlanID:   mealPlanID,
			MealPlanTask: mealPlanTask,
			HouseholdID:  sessionCtxData.ActiveHouseholdID,
			UserID:       sessionCtxData.Requester.UserID,
		}

		if err = s.dataChangesPublisher.Publish(ctx, dcm); err != nil {
			observability.AcknowledgeError(err, logger, span, "publishing to data changes topic")
		}
	}

	s.encoderDecoder.EncodeResponseWithStatus(ctx, res, mealPlanTask, http.StatusCreated)
}

// ReadHandler returns a GET handler that returns a meal plan task.
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

	// determine meal plan task ID.
	mealPlanTaskID := s.mealPlanTaskIDFetcher(req)
	tracing.AttachMealPlanEventIDToSpan(span, mealPlanTaskID)
	logger = logger.WithValue(keys.MealPlanTaskIDKey, mealPlanTaskID)

	// fetch meal plan task from database.
	x, err := s.mealPlanTaskDataManager.GetMealPlanTask(ctx, mealPlanTaskID)
	if errors.Is(err, sql.ErrNoRows) {
		s.encoderDecoder.EncodeNotFoundResponse(ctx, res)
		return
	} else if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving meal plan task")
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

	mealPlanTasks, err := s.mealPlanTaskDataManager.GetMealPlanTasksForMealPlan(ctx, mealPlanID)
	if errors.Is(err, sql.ErrNoRows) {
		// in the event no rows exist, return an empty list.
		mealPlanTasks = []*types.MealPlanTask{}
	} else if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving meal plan tasks for meal plan")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	// encode our response and peace.
	s.encoderDecoder.RespondWithData(ctx, res, mealPlanTasks)
}

// StatusChangeHandler returns a handler that updates a meal plan task.
func (s *service) StatusChangeHandler(res http.ResponseWriter, req *http.Request) {
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

	// determine meal plan task ID.
	mealPlanTaskID := s.mealPlanTaskIDFetcher(req)
	tracing.AttachMealPlanEventIDToSpan(span, mealPlanTaskID)
	logger = logger.WithValue(keys.MealPlanTaskIDKey, mealPlanTaskID)

	// read parsed input struct from request body.
	providedInput := new(types.MealPlanTaskStatusChangeRequestInput)
	if err := s.encoderDecoder.DecodeRequest(ctx, req, providedInput); err != nil {
		observability.AcknowledgeError(err, logger, span, "decoding request")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, "invalid request content", http.StatusBadRequest)
		return
	}
	providedInput.ID = mealPlanTaskID

	if err := providedInput.ValidateWithContext(ctx); err != nil {
		logger.WithValue(keys.ValidationErrorKey, err).Debug("provided input was invalid")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, err.Error(), http.StatusBadRequest)
		return
	}

	mealPlanTask, fetchMealPlanTaskErr := s.mealPlanTaskDataManager.GetMealPlanTask(ctx, mealPlanTaskID)
	if fetchMealPlanTaskErr != nil && !errors.Is(fetchMealPlanTaskErr, sql.ErrNoRows) {
		observability.AcknowledgeError(fetchMealPlanTaskErr, logger, span, "checking meal plan task existence")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	} else if errors.Is(fetchMealPlanTaskErr, sql.ErrNoRows) {
		s.encoderDecoder.EncodeNotFoundResponse(ctx, res)
		return
	}

	mealPlanTask.Update(providedInput)

	if err := s.mealPlanTaskDataManager.ChangeMealPlanTaskStatus(ctx, providedInput); err != nil {
		observability.AcknowledgeError(err, logger, span, "archiving meal plan task")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	if s.dataChangesPublisher != nil {
		dcm := &types.DataChangeMessage{
			DataType:       types.MealPlanTaskDataType,
			EventType:      types.MealPlanTaskStatusChangedCustomerEventType,
			MealPlanTask:   mealPlanTask,
			MealPlanTaskID: mealPlanTaskID,
			HouseholdID:    sessionCtxData.ActiveHouseholdID,
			UserID:         sessionCtxData.Requester.UserID,
		}

		if err := s.dataChangesPublisher.Publish(ctx, dcm); err != nil {
			observability.AcknowledgeError(err, logger, span, "publishing data change message")
		}
	}

	s.encoderDecoder.RespondWithData(ctx, res, mealPlanTask)
}
