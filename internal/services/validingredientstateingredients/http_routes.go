package validingredientstateingredients

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
	// ValidIngredientStateIngredientIDURIParamKey is a standard string that we'll use to refer to valid ingredient state ingredient IDs with.
	ValidIngredientStateIngredientIDURIParamKey = "validIngredientStateIngredientID"
	// ValidIngredientStateIDURIParamKey is a standard string that we'll use to refer to valid preparation IDs with.
	ValidIngredientStateIDURIParamKey = "validIngredientStateID"
	// ValidIngredientIDURIParamKey is a standard string that we'll use to refer to valid ingredient IDs with.
	ValidIngredientIDURIParamKey = "validIngredientID"
)

// CreateHandler is our valid ingredient state ingredient creation route.
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
	providedInput := new(types.ValidIngredientStateIngredientCreationRequestInput)
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

	input := converters.ConvertValidIngredientStateIngredientCreationRequestInputToValidIngredientStateIngredientDatabaseCreationInput(providedInput)
	input.ID = identifiers.New()

	tracing.AttachValidIngredientStateIngredientIDToSpan(span, input.ID)

	validIngredientStateIngredient, err := s.validIngredientStateIngredientDataManager.CreateValidIngredientStateIngredient(ctx, input)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "creating valid ingredient state ingredient")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	dcm := &types.DataChangeMessage{
		EventType:                      types.ValidIngredientStateIngredientCreatedCustomerEventType,
		ValidIngredientStateIngredient: validIngredientStateIngredient,
		UserID:                         sessionCtxData.Requester.UserID,
	}

	if err = s.dataChangesPublisher.Publish(ctx, dcm); err != nil {
		observability.AcknowledgeError(err, logger, span, "publishing to data changes topic")
	}

	s.encoderDecoder.EncodeResponseWithStatus(ctx, res, validIngredientStateIngredient, http.StatusCreated)
}

// ReadHandler returns a GET handler that returns a valid ingredient state ingredient.
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

	// determine valid ingredient state ingredient ID.
	validIngredientStateIngredientID := s.validIngredientStateIngredientIDFetcher(req)
	tracing.AttachValidIngredientStateIngredientIDToSpan(span, validIngredientStateIngredientID)
	logger = logger.WithValue(keys.ValidIngredientStateIngredientIDKey, validIngredientStateIngredientID)

	// fetch valid ingredient state ingredient from database.
	x, err := s.validIngredientStateIngredientDataManager.GetValidIngredientStateIngredient(ctx, validIngredientStateIngredientID)
	if errors.Is(err, sql.ErrNoRows) {
		s.encoderDecoder.EncodeNotFoundResponse(ctx, res)
		return
	} else if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving valid ingredient state ingredient")
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

	validIngredientStateIngredients, err := s.validIngredientStateIngredientDataManager.GetValidIngredientStateIngredients(ctx, filter)
	if errors.Is(err, sql.ErrNoRows) {
		// in the event no rows exist, return an empty list.
		validIngredientStateIngredients = &types.QueryFilteredResult[types.ValidIngredientStateIngredient]{Data: []*types.ValidIngredientStateIngredient{}}
	} else if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving valid ingredient state ingredients")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	// encode our response and peace.
	s.encoderDecoder.RespondWithData(ctx, res, validIngredientStateIngredients)
}

// UpdateHandler returns a handler that updates a valid ingredient state ingredient.
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
	input := new(types.ValidIngredientStateIngredientUpdateRequestInput)
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

	// determine valid ingredient state ingredient ID.
	validIngredientStateIngredientID := s.validIngredientStateIngredientIDFetcher(req)
	tracing.AttachValidIngredientStateIngredientIDToSpan(span, validIngredientStateIngredientID)
	logger = logger.WithValue(keys.ValidIngredientStateIngredientIDKey, validIngredientStateIngredientID)

	// fetch valid ingredient state ingredient from database.
	validIngredientStateIngredient, err := s.validIngredientStateIngredientDataManager.GetValidIngredientStateIngredient(ctx, validIngredientStateIngredientID)
	if errors.Is(err, sql.ErrNoRows) {
		s.encoderDecoder.EncodeNotFoundResponse(ctx, res)
		return
	} else if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving valid ingredient state ingredient for update")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	// update the valid ingredient state ingredient.
	validIngredientStateIngredient.Update(input)

	if err = s.validIngredientStateIngredientDataManager.UpdateValidIngredientStateIngredient(ctx, validIngredientStateIngredient); err != nil {
		observability.AcknowledgeError(err, logger, span, "updating valid ingredient state ingredient")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	dcm := &types.DataChangeMessage{
		EventType:                      types.ValidIngredientStateIngredientUpdatedCustomerEventType,
		ValidIngredientStateIngredient: validIngredientStateIngredient,
		UserID:                         sessionCtxData.Requester.UserID,
	}

	if err = s.dataChangesPublisher.Publish(ctx, dcm); err != nil {
		observability.AcknowledgeError(err, logger, span, "publishing data change message")
	}

	// encode our response and peace.
	s.encoderDecoder.RespondWithData(ctx, res, validIngredientStateIngredient)
}

// ArchiveHandler returns a handler that archives a valid ingredient state ingredient.
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

	// determine valid ingredient state ingredient ID.
	validIngredientStateIngredientID := s.validIngredientStateIngredientIDFetcher(req)
	tracing.AttachValidIngredientStateIngredientIDToSpan(span, validIngredientStateIngredientID)
	logger = logger.WithValue(keys.ValidIngredientStateIngredientIDKey, validIngredientStateIngredientID)

	exists, existenceCheckErr := s.validIngredientStateIngredientDataManager.ValidIngredientStateIngredientExists(ctx, validIngredientStateIngredientID)
	if existenceCheckErr != nil && !errors.Is(existenceCheckErr, sql.ErrNoRows) {
		observability.AcknowledgeError(existenceCheckErr, logger, span, "checking valid ingredient state ingredient existence")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	} else if !exists || errors.Is(existenceCheckErr, sql.ErrNoRows) {
		s.encoderDecoder.EncodeNotFoundResponse(ctx, res)
		return
	}

	if err = s.validIngredientStateIngredientDataManager.ArchiveValidIngredientStateIngredient(ctx, validIngredientStateIngredientID); err != nil {
		observability.AcknowledgeError(err, logger, span, "archiving valid ingredient state ingredient")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	dcm := &types.DataChangeMessage{
		EventType: types.ValidIngredientStateIngredientArchivedCustomerEventType,
		UserID:    sessionCtxData.Requester.UserID,
	}

	if err = s.dataChangesPublisher.Publish(ctx, dcm); err != nil {
		observability.AcknowledgeError(err, logger, span, "publishing data change message")
	}

	// encode our response and peace.
	res.WriteHeader(http.StatusNoContent)
}

// SearchByIngredientHandler is our valid ingredient state ingredient search route.
func (s *service) SearchByIngredientHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	tracing.AttachRequestToSpan(span, req)

	filter := types.ExtractQueryFilterFromRequest(req)
	tracing.AttachFilterDataToSpan(span, filter.Page, filter.Limit, filter.SortBy)

	logger := s.logger.WithRequest(req).
		WithValue(keys.FilterLimitKey, filter.Limit).
		WithValue(keys.FilterPageKey, filter.Page).
		WithValue(keys.FilterSortByKey, filter.SortBy)

	validIngredientID := s.validIngredientIDFetcher(req)
	logger = logger.WithValue(keys.ValidIngredientIDKey, validIngredientID)

	// determine user ID.
	sessionCtxData, err := s.sessionContextDataFetcher(req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving session context data")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, "unauthenticated", http.StatusUnauthorized)
		return
	}

	tracing.AttachSessionContextDataToSpan(span, sessionCtxData)
	logger = sessionCtxData.AttachToLogger(logger)

	validIngredientStateIngredients, err := s.validIngredientStateIngredientDataManager.GetValidIngredientStateIngredientsForIngredient(ctx, validIngredientID, filter)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "searching for valid ingredient state ingredients")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	s.encoderDecoder.EncodeResponseWithStatus(ctx, res, validIngredientStateIngredients, http.StatusOK)
}

// SearchByIngredientStateHandler is our valid ingredient state ingredient search route.
func (s *service) SearchByIngredientStateHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	tracing.AttachRequestToSpan(span, req)

	filter := types.ExtractQueryFilterFromRequest(req)
	tracing.AttachFilterDataToSpan(span, filter.Page, filter.Limit, filter.SortBy)

	logger := s.logger.WithRequest(req).
		WithValue(keys.FilterLimitKey, filter.Limit).
		WithValue(keys.FilterPageKey, filter.Page).
		WithValue(keys.FilterSortByKey, filter.SortBy)

	validIngredientStateID := s.validIngredientStateIDFetcher(req)
	logger = logger.WithValue(keys.ValidIngredientStateIDKey, validIngredientStateID)

	// determine user ID.
	sessionCtxData, err := s.sessionContextDataFetcher(req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving session context data")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, "unauthenticated", http.StatusUnauthorized)
		return
	}

	tracing.AttachSessionContextDataToSpan(span, sessionCtxData)
	logger = sessionCtxData.AttachToLogger(logger)

	validIngredientStateIngredients, err := s.validIngredientStateIngredientDataManager.GetValidIngredientStateIngredientsForIngredientState(ctx, validIngredientStateID, filter)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "searching for valid ingredient state ingredients")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	s.encoderDecoder.EncodeResponseWithStatus(ctx, res, validIngredientStateIngredients, http.StatusOK)
}
