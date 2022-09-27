package mealplantasks

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/prixfixeco/api_server/internal/observability"
	"github.com/prixfixeco/api_server/internal/observability/keys"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/pkg/types"
)

const (
	// MealPlanTaskIDURIParamKey is a standard string that we'll use to refer to advanced prep step IDs with.
	MealPlanTaskIDURIParamKey = "mealPlanTaskID"
)

// ReadHandler returns a GET handler that returns an advanced prep step.
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

	// determine advanced prep step ID.
	mealPlanTaskID := s.mealPlanTaskIDFetcher(req)
	tracing.AttachMealPlanEventIDToSpan(span, mealPlanTaskID)
	logger = logger.WithValue(keys.MealPlanTaskIDKey, mealPlanTaskID)

	// fetch advanced prep step from database.
	x, err := s.mealPlanTaskDataManager.GetMealPlanTask(ctx, mealPlanTaskID)
	if errors.Is(err, sql.ErrNoRows) {
		s.encoderDecoder.EncodeNotFoundResponse(ctx, res)
		return
	} else if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving advanced prep step")
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

	filter := types.ExtractQueryFilter(req)
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
		observability.AcknowledgeError(err, logger, span, "retrieving advanced prep steps for meal plan")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	// encode our response and peace.
	s.encoderDecoder.RespondWithData(ctx, res, mealPlanTasks)
}

// StatusChangeHandler returns a handler that updates an advanced prep step.
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

	// determine advanced prep step ID.
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

	prepStep, fetchMealPlanTaskErr := s.mealPlanTaskDataManager.GetMealPlanTask(ctx, mealPlanTaskID)
	if fetchMealPlanTaskErr != nil && !errors.Is(fetchMealPlanTaskErr, sql.ErrNoRows) {
		observability.AcknowledgeError(fetchMealPlanTaskErr, logger, span, "checking advanced step existence")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	} else if errors.Is(fetchMealPlanTaskErr, sql.ErrNoRows) {
		s.encoderDecoder.EncodeNotFoundResponse(ctx, res)
		return
	}

	if err := s.mealPlanTaskDataManager.ChangeMealPlanTaskStatus(ctx, providedInput); err != nil {
		observability.AcknowledgeError(err, logger, span, "archiving advanced prep step")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	// NOTE: do what you need here
	prepStep.StatusExplanation = providedInput.StatusExplanation
	prepStep.Status = providedInput.Status
	prepStep.AssignedToUser = providedInput.AssignedToUser

	if s.dataChangesPublisher != nil {
		dcm := &types.DataChangeMessage{
			DataType:                  types.MealPlanTaskDataType,
			EventType:                 types.MealPlanTaskStatusChangedCustomerEventType,
			MealPlanTask:              prepStep,
			MealPlanTaskID:            mealPlanTaskID,
			AttributableToUserID:      sessionCtxData.Requester.UserID,
			AttributableToHouseholdID: sessionCtxData.ActiveHouseholdID,
		}

		if err := s.dataChangesPublisher.Publish(ctx, dcm); err != nil {
			observability.AcknowledgeError(err, logger, span, "publishing data change message")
		}
	}

	s.encoderDecoder.RespondWithData(ctx, res, prepStep)
}
