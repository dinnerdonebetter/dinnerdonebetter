package mealplanoptionvotes

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

	// determine meal plan event ID.
	mealPlanEventID := s.mealPlanEventIDFetcher(req)
	tracing.AttachMealPlanEventIDToSpan(span, mealPlanEventID)
	logger = logger.WithValue(keys.MealPlanEventIDKey, mealPlanEventID)

	// read parsed input struct from request body.
	providedInput := new(types.MealPlanOptionVoteCreationRequestInput)
	if err = s.encoderDecoder.DecodeRequest(ctx, req, providedInput); err != nil {
		observability.AcknowledgeError(err, logger, span, "decoding request")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, "invalid request content", http.StatusBadRequest)
		return
	}

	// note, this is where you would call providedInput.ValidateWithContext, if that currently had any effect.

	input := converters.ConvertMealPlanOptionVoteCreationRequestInputToMealPlanOptionVoteDatabaseCreationInput(providedInput)
	for i := range input.Votes {
		input.Votes[i].ID = identifiers.New()
		input.Votes[i].ByUser = sessionCtxData.Requester.UserID
		tracing.AttachMealPlanOptionVoteIDToSpan(span, input.Votes[i].ID)
	}
	input.ByUser = sessionCtxData.Requester.UserID

	mealPlanOptionVotes, err := s.dataManager.CreateMealPlanOptionVote(ctx, input)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "creating meal plan option vote")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	for _, vote := range mealPlanOptionVotes {
		if s.dataChangesPublisher != nil {
			dcm := &types.DataChangeMessage{
				DataType:             types.MealPlanOptionVoteDataType,
				EventType:            types.MealPlanOptionVoteCreatedCustomerEventType,
				MealPlanID:           mealPlanID,
				MealPlanOptionID:     vote.BelongsToMealPlanOption,
				MealPlanOptionVote:   vote,
				MealPlanOptionVoteID: vote.ID,
				HouseholdID:          sessionCtxData.ActiveHouseholdID,
				UserID:               sessionCtxData.Requester.UserID,
			}

			if err = s.dataChangesPublisher.Publish(ctx, dcm); err != nil {
				observability.AcknowledgeError(err, logger, span, "publishing data change message about meal plan option vote")
			}
		}
	}

	if len(mealPlanOptionVotes) > 0 {
		lastVote := mealPlanOptionVotes[len(mealPlanOptionVotes)-1]

		// have all votes been received for an option? if so, finalize it
		mealPlanOptionFinalized, optionFinalizationErr := s.dataManager.FinalizeMealPlanOption(ctx, mealPlanID, mealPlanEventID, lastVote.BelongsToMealPlanOption, sessionCtxData.ActiveHouseholdID)
		if optionFinalizationErr != nil {
			observability.AcknowledgeError(optionFinalizationErr, logger, span, "finalizing meal plan option vote")
			s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
			return
		}

		// have all options for the meal plan been selected? if so, finalize the meal plan and fire event
		if mealPlanOptionFinalized {
			logger.Debug("meal plan option finalized")
			// fire event
			dcm := &types.DataChangeMessage{
				DataType:             types.MealPlanOptionDataType,
				EventType:            types.MealPlanOptionFinalizedCreatedCustomerEventType,
				MealPlanID:           mealPlanID,
				MealPlanOptionID:     lastVote.BelongsToMealPlanOption,
				MealPlanOptionVote:   lastVote,
				MealPlanOptionVoteID: lastVote.ID,
				HouseholdID:          sessionCtxData.ActiveHouseholdID,
				UserID:               sessionCtxData.Requester.UserID,
			}

			if dataChangePublishErr := s.dataChangesPublisher.Publish(ctx, dcm); dataChangePublishErr != nil {
				observability.AcknowledgeError(dataChangePublishErr, logger, span, "publishing data change message about meal plan option finalization")
			}

			mealPlanFinalized, finalizationErr := s.dataManager.AttemptToFinalizeMealPlan(ctx, mealPlanID, sessionCtxData.ActiveHouseholdID)
			if finalizationErr != nil {
				observability.AcknowledgeError(err, logger, span, "finalizing meal plan option vote")
				s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
				return
			}

			if mealPlanFinalized {
				logger.Debug("meal plan finalized")
				// fire event
				dcm = &types.DataChangeMessage{
					DataType:             types.MealPlanDataType,
					EventType:            types.MealPlanFinalizedCustomerEventType,
					MealPlanID:           mealPlanID,
					MealPlanOptionID:     lastVote.BelongsToMealPlanOption,
					MealPlanOptionVote:   lastVote,
					MealPlanOptionVoteID: lastVote.ID,
					HouseholdID:          sessionCtxData.ActiveHouseholdID,
					UserID:               sessionCtxData.Requester.UserID,
				}
				if dataChangePublishErr := s.dataChangesPublisher.Publish(ctx, dcm); dataChangePublishErr != nil {
					observability.AcknowledgeError(dataChangePublishErr, logger, span, "publishing data change message about meal plan finalization")
				}
			}
		}
	}

	s.encoderDecoder.EncodeResponseWithStatus(ctx, res, mealPlanOptionVotes, http.StatusCreated)
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

	// determine meal plan event ID.
	mealPlanEventID := s.mealPlanEventIDFetcher(req)
	tracing.AttachMealPlanEventIDToSpan(span, mealPlanEventID)
	logger = logger.WithValue(keys.MealPlanEventIDKey, mealPlanEventID)

	// determine meal plan option ID.
	mealPlanOptionID := s.mealPlanOptionIDFetcher(req)
	tracing.AttachMealPlanOptionIDToSpan(span, mealPlanOptionID)
	logger = logger.WithValue(keys.MealPlanOptionIDKey, mealPlanOptionID)

	// determine meal plan option vote ID.
	mealPlanOptionVoteID := s.mealPlanOptionVoteIDFetcher(req)
	tracing.AttachMealPlanOptionVoteIDToSpan(span, mealPlanOptionVoteID)
	logger = logger.WithValue(keys.MealPlanOptionVoteIDKey, mealPlanOptionVoteID)

	// fetch meal plan option vote from database.
	x, err := s.dataManager.GetMealPlanOptionVote(ctx, mealPlanID, mealPlanEventID, mealPlanOptionID, mealPlanOptionVoteID)
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

	// determine meal plan event ID.
	mealPlanEventID := s.mealPlanEventIDFetcher(req)
	tracing.AttachMealPlanEventIDToSpan(span, mealPlanEventID)
	logger = logger.WithValue(keys.MealPlanEventIDKey, mealPlanEventID)

	// determine meal plan option ID.
	mealPlanOptionID := s.mealPlanOptionIDFetcher(req)
	tracing.AttachMealPlanOptionIDToSpan(span, mealPlanOptionID)
	logger = logger.WithValue(keys.MealPlanOptionIDKey, mealPlanOptionID)

	mealPlanOptionVotes, err := s.dataManager.GetMealPlanOptionVotes(ctx, mealPlanID, mealPlanEventID, mealPlanOptionID, filter)
	if errors.Is(err, sql.ErrNoRows) {
		// in the event no rows exist, return an empty list.
		mealPlanOptionVotes = &types.QueryFilteredResult[types.MealPlanOptionVote]{Data: []*types.MealPlanOptionVote{}}
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

	// determine meal plan event ID.
	mealPlanEventID := s.mealPlanEventIDFetcher(req)
	tracing.AttachMealPlanEventIDToSpan(span, mealPlanEventID)
	logger = logger.WithValue(keys.MealPlanEventIDKey, mealPlanEventID)

	// determine meal plan option ID.
	mealPlanOptionID := s.mealPlanOptionIDFetcher(req)
	tracing.AttachMealPlanOptionIDToSpan(span, mealPlanOptionID)
	logger = logger.WithValue(keys.MealPlanOptionIDKey, mealPlanOptionID)

	// determine meal plan option vote ID.
	mealPlanOptionVoteID := s.mealPlanOptionVoteIDFetcher(req)
	tracing.AttachMealPlanOptionVoteIDToSpan(span, mealPlanOptionVoteID)
	logger = logger.WithValue(keys.MealPlanOptionVoteIDKey, mealPlanOptionVoteID)

	// fetch meal plan option vote from database.
	mealPlanOptionVote, err := s.dataManager.GetMealPlanOptionVote(ctx, mealPlanID, mealPlanEventID, mealPlanOptionID, mealPlanOptionVoteID)
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

	if err = s.dataManager.UpdateMealPlanOptionVote(ctx, mealPlanOptionVote); err != nil {
		observability.AcknowledgeError(err, logger, span, "updating meal plan option vote")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	if s.dataChangesPublisher != nil {
		dcm := &types.DataChangeMessage{
			DataType:             types.MealPlanOptionVoteDataType,
			EventType:            types.MealPlanOptionVoteUpdatedCustomerEventType,
			MealPlanID:           mealPlanID,
			MealPlanOptionID:     mealPlanOptionID,
			MealPlanOptionVote:   mealPlanOptionVote,
			MealPlanOptionVoteID: mealPlanOptionVote.ID,
			HouseholdID:          sessionCtxData.ActiveHouseholdID,
			UserID:               sessionCtxData.Requester.UserID,
		}

		if err = s.dataChangesPublisher.Publish(ctx, dcm); err != nil {
			observability.AcknowledgeError(err, logger, span, "publishing data change message")
		}
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

	// determine meal plan event ID.
	mealPlanEventID := s.mealPlanEventIDFetcher(req)
	tracing.AttachMealPlanEventIDToSpan(span, mealPlanEventID)
	logger = logger.WithValue(keys.MealPlanEventIDKey, mealPlanEventID)

	// determine meal plan option ID.
	mealPlanOptionID := s.mealPlanOptionIDFetcher(req)
	tracing.AttachMealPlanOptionIDToSpan(span, mealPlanOptionID)
	logger = logger.WithValue(keys.MealPlanOptionIDKey, mealPlanOptionID)

	// determine meal plan option vote ID.
	mealPlanOptionVoteID := s.mealPlanOptionVoteIDFetcher(req)
	tracing.AttachMealPlanOptionVoteIDToSpan(span, mealPlanOptionVoteID)
	logger = logger.WithValue(keys.MealPlanOptionVoteIDKey, mealPlanOptionVoteID)

	exists, existenceCheckErr := s.dataManager.MealPlanOptionVoteExists(ctx, mealPlanID, mealPlanEventID, mealPlanOptionID, mealPlanOptionVoteID)
	if existenceCheckErr != nil && !errors.Is(existenceCheckErr, sql.ErrNoRows) {
		observability.AcknowledgeError(existenceCheckErr, logger, span, "checking meal plan option vote existence")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	} else if !exists || errors.Is(existenceCheckErr, sql.ErrNoRows) {
		s.encoderDecoder.EncodeNotFoundResponse(ctx, res)
		return
	}

	if err = s.dataManager.ArchiveMealPlanOptionVote(ctx, mealPlanID, mealPlanEventID, mealPlanOptionID, mealPlanOptionVoteID); err != nil {
		observability.AcknowledgeError(err, logger, span, "archiving meal plan option vote")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	if s.dataChangesPublisher != nil {
		dcm := &types.DataChangeMessage{
			DataType:             types.MealPlanOptionVoteDataType,
			EventType:            types.MealPlanOptionVoteArchivedCustomerEventType,
			MealPlanID:           mealPlanID,
			MealPlanOptionID:     mealPlanOptionID,
			MealPlanOptionVoteID: mealPlanOptionVoteID,
			HouseholdID:          sessionCtxData.ActiveHouseholdID,
			UserID:               sessionCtxData.Requester.UserID,
		}

		if err = s.dataChangesPublisher.Publish(ctx, dcm); err != nil {
			observability.AcknowledgeError(err, logger, span, "publishing data change message")
		}
	}

	// encode our response and peace.
	res.WriteHeader(http.StatusNoContent)
}
