package advancedprepsteps

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
	// AdvancedPrepStepIDURIParamKey is a standard string that we'll use to refer to advanced prep step IDs with.
	AdvancedPrepStepIDURIParamKey = "advancedPrepStepID"
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
	advancedPrepStepID := s.advancedPrepStepIDFetcher(req)
	tracing.AttachMealPlanEventIDToSpan(span, advancedPrepStepID)
	logger = logger.WithValue(keys.AdvancedPrepStepIDKey, advancedPrepStepID)

	// fetch advanced prep step from database.
	x, err := s.advancedPrepStepDataManager.GetAdvancedPrepStep(ctx, advancedPrepStepID)
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

	advancedPrepSteps, err := s.advancedPrepStepDataManager.GetAdvancedPrepStepsForMealPlan(ctx, mealPlanID)
	if errors.Is(err, sql.ErrNoRows) {
		// in the event no rows exist, return an empty list.
		advancedPrepSteps = []*types.AdvancedPrepStep{}
	} else if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving advanced prep steps for meal plan")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	// encode our response and peace.
	s.encoderDecoder.RespondWithData(ctx, res, advancedPrepSteps)
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
	advancedPrepStepID := s.advancedPrepStepIDFetcher(req)
	tracing.AttachMealPlanEventIDToSpan(span, advancedPrepStepID)
	logger = logger.WithValue(keys.AdvancedPrepStepIDKey, advancedPrepStepID)

	// read parsed input struct from request body.
	providedInput := new(types.AdvancedPrepStepStatusChangeRequestInput)
	if err := s.encoderDecoder.DecodeRequest(ctx, req, providedInput); err != nil {
		observability.AcknowledgeError(err, logger, span, "decoding request")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, "invalid request content", http.StatusBadRequest)
		return
	}
	providedInput.ID = advancedPrepStepID

	if err := providedInput.ValidateWithContext(ctx); err != nil {
		logger.WithValue(keys.ValidationErrorKey, err).Debug("provided input was invalid")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, err.Error(), http.StatusBadRequest)
		return
	}

	exists, existenceCheckErr := s.advancedPrepStepDataManager.AdvancedPrepStepExists(ctx, mealPlanID, advancedPrepStepID)
	if existenceCheckErr != nil && !errors.Is(existenceCheckErr, sql.ErrNoRows) {
		observability.AcknowledgeError(existenceCheckErr, logger, span, "checking recipe existence")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	} else if !exists || errors.Is(existenceCheckErr, sql.ErrNoRows) {
		s.encoderDecoder.EncodeNotFoundResponse(ctx, res)
		return
	}

	if err := s.advancedPrepStepDataManager.ChangeAdvancedPrepStepStatus(ctx, providedInput); err != nil {
		observability.AcknowledgeError(err, logger, span, "archiving advanced prep step")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	if s.dataChangesPublisher != nil {
		dcm := &types.DataChangeMessage{
			DataType:                  types.AdvancedPrepStepDataType,
			EventType:                 types.AdvancedPrepStepStatusChangedCustomerEventType,
			MealPlanEventID:           advancedPrepStepID,
			AttributableToUserID:      sessionCtxData.Requester.UserID,
			AttributableToHouseholdID: sessionCtxData.ActiveHouseholdID,
		}

		if err := s.dataChangesPublisher.Publish(ctx, dcm); err != nil {
			observability.AcknowledgeError(err, logger, span, "publishing data change message")
		}
	}

	// encode our response and peace.
	res.WriteHeader(http.StatusNoContent)
}
