package mealplanoptionvotes

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
	// MealPlanOptionVoteIDURIParamKey is a standard string that we'll use to refer to meal plan option vote IDs with.
	MealPlanOptionVoteIDURIParamKey = "mealPlanOptionVoteID"
)

// CreateHandler is our meal plan option vote creation route.
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

	// determine meal plan ID.
	mealPlanID := s.mealPlanIDFetcher(req)
	tracing.AttachMealPlanIDToSpan(span, mealPlanID)
	logger = logger.WithValue(keys.MealPlanIDKey, mealPlanID)

	// determine meal plan option ID.
	mealPlanOptionID := s.mealPlanOptionIDFetcher(req)
	tracing.AttachMealPlanOptionIDToSpan(span, mealPlanOptionID)
	logger = logger.WithValue(keys.MealPlanOptionIDKey, mealPlanOptionID)

	// read parsed input struct from request body.
	providedInput := new(types.MealPlanOptionVoteCreationRequestInput)
	if err = s.encoderDecoder.DecodeRequest(ctx, req, providedInput); err != nil {
		observability.AcknowledgeError(err, logger, span, "decoding request")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, "invalid request content", http.StatusBadRequest)
		return
	}
	providedInput.BelongsToMealPlanOption = mealPlanOptionID
	providedInput.ByUser = sessionCtxData.Requester.UserID

	// note, this is where you would call providedInput.ValidateWithContext, if that currently had any effect.

	input := types.MealPlanOptionVoteDatabaseCreationInputFromMealPlanOptionVoteCreationInput(providedInput)
	input.ID = ksuid.New().String()
	input.ByUser = sessionCtxData.Requester.UserID
	tracing.AttachMealPlanOptionVoteIDToSpan(span, input.ID)

	// create meal plan option vote in database.
	preWrite := &types.PreWriteMessage{
		DataType:                  types.MealPlanOptionVoteDataType,
		MealPlanID:                mealPlanID,
		MealPlanOptionID:          mealPlanOptionID,
		MealPlanOptionVote:        input,
		AttributableToUserID:      sessionCtxData.Requester.UserID,
		AttributableToHouseholdID: sessionCtxData.ActiveHouseholdID,
	}
	if err = s.preWritesPublisher.Publish(ctx, preWrite); err != nil {
		observability.AcknowledgeError(err, logger, span, "publishing meal plan option vote write message")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	pwr := types.PreWriteResponse{ID: input.ID}

	if err = s.customerDataCollector.EventOccurred(ctx, "meal_plan_option_vote_created", sessionCtxData.Requester.UserID, map[string]interface{}{
		keys.MealPlanIDKey:       mealPlanID,
		keys.MealPlanOptionIDKey: mealPlanOptionID,
		keys.HouseholdIDKey:      sessionCtxData.ActiveHouseholdID,
	}); err != nil {
		logger.Error(err, "notifying customer data platform")
	}

	s.encoderDecoder.EncodeResponseWithStatus(ctx, res, pwr, http.StatusAccepted)
}

// ReadHandler returns a GET handler that returns a meal plan option vote.
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

	// determine meal plan option ID.
	mealPlanOptionID := s.mealPlanOptionIDFetcher(req)
	tracing.AttachMealPlanOptionIDToSpan(span, mealPlanOptionID)
	logger = logger.WithValue(keys.MealPlanOptionIDKey, mealPlanOptionID)

	// determine meal plan option vote ID.
	mealPlanOptionVoteID := s.mealPlanOptionVoteIDFetcher(req)
	tracing.AttachMealPlanOptionVoteIDToSpan(span, mealPlanOptionVoteID)
	logger = logger.WithValue(keys.MealPlanOptionVoteIDKey, mealPlanOptionVoteID)

	// fetch meal plan option vote from database.
	x, err := s.mealPlanOptionVoteDataManager.GetMealPlanOptionVote(ctx, mealPlanID, mealPlanOptionID, mealPlanOptionVoteID)
	if errors.Is(err, sql.ErrNoRows) {
		s.encoderDecoder.EncodeNotFoundResponse(ctx, res)
		return
	} else if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving meal plan option vote")
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

	// determine meal plan ID.
	mealPlanID := s.mealPlanIDFetcher(req)
	tracing.AttachMealPlanIDToSpan(span, mealPlanID)
	logger = logger.WithValue(keys.MealPlanIDKey, mealPlanID)

	// determine meal plan option ID.
	mealPlanOptionID := s.mealPlanOptionIDFetcher(req)
	tracing.AttachMealPlanOptionIDToSpan(span, mealPlanOptionID)
	logger = logger.WithValue(keys.MealPlanOptionIDKey, mealPlanOptionID)

	mealPlanOptionVotes, err := s.mealPlanOptionVoteDataManager.GetMealPlanOptionVotes(ctx, mealPlanID, mealPlanOptionID, filter)
	if errors.Is(err, sql.ErrNoRows) {
		// in the event no rows exist, return an empty list.
		mealPlanOptionVotes = &types.MealPlanOptionVoteList{MealPlanOptionVotes: []*types.MealPlanOptionVote{}}
	} else if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving meal plan option votes")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	// encode our response and peace.
	s.encoderDecoder.RespondWithData(ctx, res, mealPlanOptionVotes)
}

// UpdateHandler returns a handler that updates a meal plan option vote.
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
	input := new(types.MealPlanOptionVoteUpdateRequestInput)
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

	// determine meal plan ID.
	mealPlanID := s.mealPlanIDFetcher(req)
	tracing.AttachMealPlanIDToSpan(span, mealPlanID)
	logger = logger.WithValue(keys.MealPlanIDKey, mealPlanID)

	// determine meal plan option ID.
	mealPlanOptionID := s.mealPlanOptionIDFetcher(req)
	tracing.AttachMealPlanOptionIDToSpan(span, mealPlanOptionID)
	logger = logger.WithValue(keys.MealPlanOptionIDKey, mealPlanOptionID)

	// determine meal plan option vote ID.
	mealPlanOptionVoteID := s.mealPlanOptionVoteIDFetcher(req)
	tracing.AttachMealPlanOptionVoteIDToSpan(span, mealPlanOptionVoteID)
	logger = logger.WithValue(keys.MealPlanOptionVoteIDKey, mealPlanOptionVoteID)

	// fetch meal plan option vote from database.
	mealPlanOptionVote, err := s.mealPlanOptionVoteDataManager.GetMealPlanOptionVote(ctx, mealPlanID, mealPlanOptionID, mealPlanOptionVoteID)
	if errors.Is(err, sql.ErrNoRows) {
		s.encoderDecoder.EncodeNotFoundResponse(ctx, res)
		return
	} else if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving meal plan option vote for update")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	// update the meal plan option vote.
	mealPlanOptionVote.Update(input)

	pum := &types.PreUpdateMessage{
		DataType:                  types.MealPlanOptionVoteDataType,
		MealPlanID:                mealPlanID,
		MealPlanOptionID:          mealPlanOptionID,
		MealPlanOptionVote:        mealPlanOptionVote,
		AttributableToUserID:      sessionCtxData.Requester.UserID,
		AttributableToHouseholdID: sessionCtxData.ActiveHouseholdID,
	}
	if err = s.preUpdatesPublisher.Publish(ctx, pum); err != nil {
		observability.AcknowledgeError(err, logger, span, "publishing meal plan option vote update message")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	if err = s.customerDataCollector.EventOccurred(ctx, "meal_plan_option_vote_updated", sessionCtxData.Requester.UserID, map[string]interface{}{
		keys.MealPlanIDKey:       mealPlanID,
		keys.MealPlanOptionIDKey: mealPlanOptionID,
		keys.HouseholdIDKey:      sessionCtxData.ActiveHouseholdID,
	}); err != nil {
		logger.Error(err, "notifying customer data platform")
	}

	// encode our response and peace.
	s.encoderDecoder.RespondWithData(ctx, res, mealPlanOptionVote)
}

// ArchiveHandler returns a handler that archives a meal plan option vote.
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

	// determine meal plan option ID.
	mealPlanOptionID := s.mealPlanOptionIDFetcher(req)
	tracing.AttachMealPlanOptionIDToSpan(span, mealPlanOptionID)
	logger = logger.WithValue(keys.MealPlanOptionIDKey, mealPlanOptionID)

	// determine meal plan option vote ID.
	mealPlanOptionVoteID := s.mealPlanOptionVoteIDFetcher(req)
	tracing.AttachMealPlanOptionVoteIDToSpan(span, mealPlanOptionVoteID)
	logger = logger.WithValue(keys.MealPlanOptionVoteIDKey, mealPlanOptionVoteID)

	exists, existenceCheckErr := s.mealPlanOptionVoteDataManager.MealPlanOptionVoteExists(ctx, mealPlanID, mealPlanOptionID, mealPlanOptionVoteID)
	if existenceCheckErr != nil && !errors.Is(existenceCheckErr, sql.ErrNoRows) {
		observability.AcknowledgeError(existenceCheckErr, logger, span, "checking meal plan option vote existence")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	} else if !exists || errors.Is(existenceCheckErr, sql.ErrNoRows) {
		s.encoderDecoder.EncodeNotFoundResponse(ctx, res)
		return
	}

	pam := &types.PreArchiveMessage{
		DataType:                  types.MealPlanOptionVoteDataType,
		MealPlanID:                mealPlanID,
		MealPlanOptionID:          mealPlanOptionID,
		MealPlanOptionVoteID:      mealPlanOptionVoteID,
		AttributableToUserID:      sessionCtxData.Requester.UserID,
		AttributableToHouseholdID: sessionCtxData.ActiveHouseholdID,
	}
	if err = s.preArchivesPublisher.Publish(ctx, pam); err != nil {
		observability.AcknowledgeError(err, logger, span, "publishing meal plan option vote archive message")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	if err = s.customerDataCollector.EventOccurred(ctx, "meal_plan_option_vote_archived", sessionCtxData.Requester.UserID, map[string]interface{}{
		keys.MealPlanIDKey:       mealPlanID,
		keys.MealPlanOptionIDKey: mealPlanOptionID,
		keys.HouseholdIDKey:      sessionCtxData.ActiveHouseholdID,
	}); err != nil {
		logger.Error(err, "notifying customer data platform")
	}

	// encode our response and peace.
	res.WriteHeader(http.StatusNoContent)
}
