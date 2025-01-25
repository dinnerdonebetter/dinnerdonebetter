package mealplanning

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/dinnerdonebetter/backend/internal/lib/identifiers"
	"github.com/dinnerdonebetter/backend/internal/lib/observability"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"

	servertiming "github.com/mitchellh/go-server-timing"
)

// CreateMealPlanOptionVoteHandler is our meal plan option vote creation route.
func (s *service) CreateMealPlanOptionVoteHandler(res http.ResponseWriter, req *http.Request) {
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

	// determine meal plan event ID.
	mealPlanEventID := s.mealPlanEventIDFetcher(req)
	tracing.AttachToSpan(span, keys.MealPlanEventIDKey, mealPlanEventID)
	logger = logger.WithValue(keys.MealPlanEventIDKey, mealPlanEventID)

	// read parsed input struct from request body.
	providedInput := new(types.MealPlanOptionVoteCreationRequestInput)
	if err = s.encoderDecoder.DecodeRequest(ctx, req, providedInput); err != nil {
		observability.AcknowledgeError(err, logger, span, "decoding request")
		errRes := types.NewAPIErrorResponse("invalid request content", types.ErrDecodingRequestInput, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusBadRequest)
		return
	}

	// note, this is where you would call providedInput.ValidateWithContext, if that currently had any effect.
	if err = providedInput.ValidateWithContext(ctx); err != nil {
		observability.AcknowledgeError(err, logger, span, "decoding request")
		errRes := types.NewAPIErrorResponse("invalid request content", types.ErrDecodingRequestInput, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusBadRequest)
		return
	}

	eligible, err := s.mealPlanningDataManager.MealPlanEventIsEligibleForVoting(ctx, mealPlanID, mealPlanEventID)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "checking event vote eligibility")
		errRes := types.NewAPIErrorResponse("database error", types.ErrTalkingToDatabase, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusInternalServerError)
		return
	}
	if !eligible {
		errRes := types.NewAPIErrorResponse("meal plan event is not eligible for voting", types.ErrNothingSpecific, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusBadRequest)
		return
	}

	input := converters.ConvertMealPlanOptionVoteCreationRequestInputToMealPlanOptionVoteDatabaseCreationInput(providedInput)
	for i := range input.Votes {
		input.Votes[i].ID = identifiers.New()
		input.Votes[i].ByUser = sessionCtxData.Requester.UserID
		tracing.AttachToSpan(span, keys.MealPlanOptionVoteIDKey, input.Votes[i].ID)
	}
	input.ByUser = sessionCtxData.Requester.UserID

	mealPlanOptionVotes, err := s.mealPlanningDataManager.CreateMealPlanOptionVote(ctx, input)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "creating meal plan option vote")
		errRes := types.NewAPIErrorResponse("database error", types.ErrTalkingToDatabase, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusInternalServerError)
		return
	}

	for _, vote := range mealPlanOptionVotes {
		dcm := &types.DataChangeMessage{
			EventType:            types.MealPlanOptionVoteCreatedServiceEventType,
			MealPlanID:           mealPlanID,
			MealPlanOptionID:     vote.BelongsToMealPlanOption,
			MealPlanOptionVote:   vote,
			MealPlanOptionVoteID: vote.ID,
			HouseholdID:          sessionCtxData.ActiveHouseholdID,
			UserID:               sessionCtxData.Requester.UserID,
		}

		go s.dataChangesPublisher.PublishAsync(ctx, dcm)
	}

	if len(mealPlanOptionVotes) > 0 {
		lastVote := mealPlanOptionVotes[len(mealPlanOptionVotes)-1]

		// have all votes been received for an option? if so, finalize it
		mealPlanOptionFinalized, optionFinalizationErr := s.mealPlanningDataManager.FinalizeMealPlanOption(ctx, mealPlanID, mealPlanEventID, lastVote.BelongsToMealPlanOption, sessionCtxData.ActiveHouseholdID)
		if optionFinalizationErr != nil {
			observability.AcknowledgeError(optionFinalizationErr, logger, span, "finalizing meal plan option vote")
			errRes := types.NewAPIErrorResponse("database error", types.ErrTalkingToDatabase, responseDetails)
			s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusInternalServerError)
			return
		}

		// have all options for the meal plan been selected? if so, finalize the meal plan and fire event
		if mealPlanOptionFinalized {
			logger.Debug("meal plan option finalized")
			// fire event
			dcm := &types.DataChangeMessage{
				EventType:            types.MealPlanOptionFinalizedCreatedServiceEventType,
				MealPlanID:           mealPlanID,
				MealPlanOptionID:     lastVote.BelongsToMealPlanOption,
				MealPlanOptionVote:   lastVote,
				MealPlanOptionVoteID: lastVote.ID,
				HouseholdID:          sessionCtxData.ActiveHouseholdID,
				UserID:               sessionCtxData.Requester.UserID,
			}

			go s.dataChangesPublisher.PublishAsync(ctx, dcm)

			mealPlanFinalized, finalizationErr := s.mealPlanningDataManager.AttemptToFinalizeMealPlan(ctx, mealPlanID, sessionCtxData.ActiveHouseholdID)
			if finalizationErr != nil {
				observability.AcknowledgeError(err, logger, span, "finalizing meal plan option vote")
				errRes := types.NewAPIErrorResponse("database error", types.ErrTalkingToDatabase, responseDetails)
				s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusInternalServerError)
				return
			}

			if mealPlanFinalized {
				logger.Debug("meal plan finalized")
				// fire event
				dcm = &types.DataChangeMessage{
					EventType:            types.MealPlanFinalizedServiceEventType,
					MealPlanID:           mealPlanID,
					MealPlanOptionID:     lastVote.BelongsToMealPlanOption,
					MealPlanOptionVote:   lastVote,
					MealPlanOptionVoteID: lastVote.ID,
					HouseholdID:          sessionCtxData.ActiveHouseholdID,
					UserID:               sessionCtxData.Requester.UserID,
				}

				go s.dataChangesPublisher.PublishAsync(ctx, dcm)
			}
		}
	}

	responseValue := &types.APIResponse[[]*types.MealPlanOptionVote]{
		Details: responseDetails,
		Data:    mealPlanOptionVotes,
	}

	s.encoderDecoder.EncodeResponseWithStatus(ctx, res, responseValue, http.StatusCreated)
}

// ReadMealPlanOptionVoteHandler returns a GET handler that returns a meal plan option vote.
func (s *service) ReadMealPlanOptionVoteHandler(res http.ResponseWriter, req *http.Request) {
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

	// determine meal plan event ID.
	mealPlanEventID := s.mealPlanEventIDFetcher(req)
	tracing.AttachToSpan(span, keys.MealPlanEventIDKey, mealPlanEventID)
	logger = logger.WithValue(keys.MealPlanEventIDKey, mealPlanEventID)

	// determine meal plan option ID.
	mealPlanOptionID := s.mealPlanOptionIDFetcher(req)
	tracing.AttachToSpan(span, keys.MealPlanOptionIDKey, mealPlanOptionID)
	logger = logger.WithValue(keys.MealPlanOptionIDKey, mealPlanOptionID)

	// determine meal plan option vote ID.
	mealPlanOptionVoteID := s.mealPlanOptionVoteIDFetcher(req)
	tracing.AttachToSpan(span, keys.MealPlanOptionVoteIDKey, mealPlanOptionVoteID)
	logger = logger.WithValue(keys.MealPlanOptionVoteIDKey, mealPlanOptionVoteID)

	// fetch meal plan option vote from database.
	x, err := s.mealPlanningDataManager.GetMealPlanOptionVote(ctx, mealPlanID, mealPlanEventID, mealPlanOptionID, mealPlanOptionVoteID)
	if errors.Is(err, sql.ErrNoRows) {
		errRes := types.NewAPIErrorResponse("not found", types.ErrDataNotFound, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusNotFound)
		return
	} else if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving meal plan option vote")
		errRes := types.NewAPIErrorResponse("database error", types.ErrTalkingToDatabase, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusInternalServerError)
		return
	}

	responseValue := &types.APIResponse[*types.MealPlanOptionVote]{
		Details: responseDetails,
		Data:    x,
	}

	// encode our response and peace.
	s.encoderDecoder.RespondWithData(ctx, res, responseValue)
}

// ListMealPlanOptionVoteHandler is our list route.
func (s *service) ListMealPlanOptionVoteHandler(res http.ResponseWriter, req *http.Request) {
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
	responseDetails.CurrentHouseholdID = sessionCtxData.ActiveHouseholdID

	// determine meal plan ID.
	mealPlanID := s.mealPlanIDFetcher(req)
	tracing.AttachToSpan(span, keys.MealPlanIDKey, mealPlanID)
	logger = logger.WithValue(keys.MealPlanIDKey, mealPlanID)

	// determine meal plan event ID.
	mealPlanEventID := s.mealPlanEventIDFetcher(req)
	tracing.AttachToSpan(span, keys.MealPlanEventIDKey, mealPlanEventID)
	logger = logger.WithValue(keys.MealPlanEventIDKey, mealPlanEventID)

	// determine meal plan option ID.
	mealPlanOptionID := s.mealPlanOptionIDFetcher(req)
	tracing.AttachToSpan(span, keys.MealPlanOptionIDKey, mealPlanOptionID)
	logger = logger.WithValue(keys.MealPlanOptionIDKey, mealPlanOptionID)

	mealPlanOptionVotes, err := s.mealPlanningDataManager.GetMealPlanOptionVotes(ctx, mealPlanID, mealPlanEventID, mealPlanOptionID, filter)
	if errors.Is(err, sql.ErrNoRows) {
		// in the event no rows exist, return an empty list.
		mealPlanOptionVotes = &types.QueryFilteredResult[types.MealPlanOptionVote]{Data: []*types.MealPlanOptionVote{}}
	} else if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving meal plan option votes")
		errRes := types.NewAPIErrorResponse("database error", types.ErrTalkingToDatabase, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusInternalServerError)
		return
	}

	responseValue := &types.APIResponse[[]*types.MealPlanOptionVote]{
		Details:    responseDetails,
		Data:       mealPlanOptionVotes.Data,
		Pagination: &mealPlanOptionVotes.Pagination,
	}

	// encode our response and peace.
	s.encoderDecoder.RespondWithData(ctx, res, responseValue)
}

// UpdateMealPlanOptionVoteHandler returns a handler that updates a meal plan option vote.
func (s *service) UpdateMealPlanOptionVoteHandler(res http.ResponseWriter, req *http.Request) {
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

	// check for parsed input attached to session context data.
	input := new(types.MealPlanOptionVoteUpdateRequestInput)
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

	// determine meal plan event ID.
	mealPlanEventID := s.mealPlanEventIDFetcher(req)
	tracing.AttachToSpan(span, keys.MealPlanEventIDKey, mealPlanEventID)
	logger = logger.WithValue(keys.MealPlanEventIDKey, mealPlanEventID)

	// determine meal plan option ID.
	mealPlanOptionID := s.mealPlanOptionIDFetcher(req)
	tracing.AttachToSpan(span, keys.MealPlanOptionIDKey, mealPlanOptionID)
	logger = logger.WithValue(keys.MealPlanOptionIDKey, mealPlanOptionID)

	// determine meal plan option vote ID.
	mealPlanOptionVoteID := s.mealPlanOptionVoteIDFetcher(req)
	tracing.AttachToSpan(span, keys.MealPlanOptionVoteIDKey, mealPlanOptionVoteID)
	logger = logger.WithValue(keys.MealPlanOptionVoteIDKey, mealPlanOptionVoteID)

	// fetch meal plan option vote from database.
	mealPlanOptionVote, err := s.mealPlanningDataManager.GetMealPlanOptionVote(ctx, mealPlanID, mealPlanEventID, mealPlanOptionID, mealPlanOptionVoteID)
	if errors.Is(err, sql.ErrNoRows) {
		errRes := types.NewAPIErrorResponse("not found", types.ErrDataNotFound, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusNotFound)
		return
	} else if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving meal plan option vote for update")
		errRes := types.NewAPIErrorResponse("database error", types.ErrTalkingToDatabase, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusInternalServerError)
		return
	}

	// update the meal plan option vote.
	mealPlanOptionVote.Update(input)

	if err = s.mealPlanningDataManager.UpdateMealPlanOptionVote(ctx, mealPlanOptionVote); err != nil {
		observability.AcknowledgeError(err, logger, span, "updating meal plan option vote")
		errRes := types.NewAPIErrorResponse("database error", types.ErrTalkingToDatabase, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusInternalServerError)
		return
	}

	dcm := &types.DataChangeMessage{
		EventType:            types.MealPlanOptionVoteUpdatedServiceEventType,
		MealPlanID:           mealPlanID,
		MealPlanOptionID:     mealPlanOptionID,
		MealPlanOptionVote:   mealPlanOptionVote,
		MealPlanOptionVoteID: mealPlanOptionVote.ID,
		HouseholdID:          sessionCtxData.ActiveHouseholdID,
		UserID:               sessionCtxData.Requester.UserID,
	}

	go s.dataChangesPublisher.PublishAsync(ctx, dcm)

	responseValue := &types.APIResponse[*types.MealPlanOptionVote]{
		Details: responseDetails,
		Data:    mealPlanOptionVote,
	}

	// encode our response and peace.
	s.encoderDecoder.RespondWithData(ctx, res, responseValue)
}

// ArchiveMealPlanOptionVoteHandler returns a handler that archives a meal plan option vote.
func (s *service) ArchiveMealPlanOptionVoteHandler(res http.ResponseWriter, req *http.Request) {
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

	// determine meal plan event ID.
	mealPlanEventID := s.mealPlanEventIDFetcher(req)
	tracing.AttachToSpan(span, keys.MealPlanEventIDKey, mealPlanEventID)
	logger = logger.WithValue(keys.MealPlanEventIDKey, mealPlanEventID)

	// determine meal plan option ID.
	mealPlanOptionID := s.mealPlanOptionIDFetcher(req)
	tracing.AttachToSpan(span, keys.MealPlanOptionIDKey, mealPlanOptionID)
	logger = logger.WithValue(keys.MealPlanOptionIDKey, mealPlanOptionID)

	// determine meal plan option vote ID.
	mealPlanOptionVoteID := s.mealPlanOptionVoteIDFetcher(req)
	tracing.AttachToSpan(span, keys.MealPlanOptionVoteIDKey, mealPlanOptionVoteID)
	logger = logger.WithValue(keys.MealPlanOptionVoteIDKey, mealPlanOptionVoteID)

	exists, err := s.mealPlanningDataManager.MealPlanOptionVoteExists(ctx, mealPlanID, mealPlanEventID, mealPlanOptionID, mealPlanOptionVoteID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		observability.AcknowledgeError(err, logger, span, "checking meal plan option vote existence")
		errRes := types.NewAPIErrorResponse("database error", types.ErrTalkingToDatabase, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusInternalServerError)
		return
	} else if !exists || errors.Is(err, sql.ErrNoRows) {
		errRes := types.NewAPIErrorResponse("not found", types.ErrDataNotFound, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusNotFound)
		return
	}

	if err = s.mealPlanningDataManager.ArchiveMealPlanOptionVote(ctx, mealPlanID, mealPlanEventID, mealPlanOptionID, mealPlanOptionVoteID); err != nil {
		observability.AcknowledgeError(err, logger, span, "archiving meal plan option vote")
		errRes := types.NewAPIErrorResponse("database error", types.ErrTalkingToDatabase, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusInternalServerError)
		return
	}

	dcm := &types.DataChangeMessage{
		EventType:            types.MealPlanOptionVoteArchivedServiceEventType,
		MealPlanID:           mealPlanID,
		MealPlanOptionID:     mealPlanOptionID,
		MealPlanOptionVoteID: mealPlanOptionVoteID,
		HouseholdID:          sessionCtxData.ActiveHouseholdID,
		UserID:               sessionCtxData.Requester.UserID,
	}

	go s.dataChangesPublisher.PublishAsync(ctx, dcm)

	responseValue := &types.APIResponse[*types.MealPlanOptionVote]{
		Details: responseDetails,
	}

	// let everybody go home.
	s.encoderDecoder.RespondWithData(ctx, res, responseValue)
}
