package validingredients

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
	// ValidIngredientIDURIParamKey is a standard string that we'll use to refer to valid ingredient IDs with.
	ValidIngredientIDURIParamKey = "validIngredientID"
)

// CreateHandler is our valid ingredient creation route.
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
	providedInput := new(types.ValidIngredientCreationRequestInput)
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

	input := converters.ConvertValidIngredientCreationRequestInputToValidIngredientDatabaseCreationInput(providedInput)
	input.ID = identifiers.New()

	tracing.AttachValidIngredientIDToSpan(span, input.ID)

	validIngredient, err := s.validIngredientDataManager.CreateValidIngredient(ctx, input)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "creating valid ingredient")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	if s.dataChangesPublisher != nil {
		dcm := &types.DataChangeMessage{
			DataType:             types.ValidIngredientDataType,
			EventType:            types.ValidIngredientCreatedCustomerEventType,
			ValidIngredient:      validIngredient,
			AttributableToUserID: sessionCtxData.Requester.UserID,
		}

		if err = s.dataChangesPublisher.Publish(ctx, dcm); err != nil {
			observability.AcknowledgeError(err, logger, span, "publishing to data changes topic")
		}
	}

	s.encoderDecoder.EncodeResponseWithStatus(ctx, res, validIngredient, http.StatusCreated)
}

// ReadHandler returns a GET handler that returns a valid ingredient.
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

	// determine valid ingredient ID.
	validIngredientID := s.validIngredientIDFetcher(req)
	tracing.AttachValidIngredientIDToSpan(span, validIngredientID)
	logger = logger.WithValue(keys.ValidIngredientIDKey, validIngredientID)

	// fetch valid ingredient from database.
	x, err := s.validIngredientDataManager.GetValidIngredient(ctx, validIngredientID)
	if errors.Is(err, sql.ErrNoRows) {
		s.encoderDecoder.EncodeNotFoundResponse(ctx, res)
		return
	} else if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving valid ingredient")
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

	validIngredients, err := s.validIngredientDataManager.GetValidIngredients(ctx, filter)
	if errors.Is(err, sql.ErrNoRows) {
		// in the event no rows exist, return an empty list.
		validIngredients = &types.ValidIngredientList{ValidIngredients: []*types.ValidIngredient{}}
	} else if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving valid ingredients")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	// encode our response and peace.
	s.encoderDecoder.RespondWithData(ctx, res, validIngredients)
}

// SearchHandler is our search route.
func (s *service) SearchHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	query := req.URL.Query().Get(types.SearchQueryKey)
	filter := types.ExtractQueryFilterFromRequest(req)
	logger := s.logger.WithRequest(req).
		WithValue(keys.FilterLimitKey, filter.Limit).
		WithValue(keys.FilterPageKey, filter.Page).
		WithValue(keys.FilterSortByKey, filter.SortBy).
		WithValue(keys.SearchQueryKey, query)

	tracing.AttachRequestToSpan(span, req)
	tracing.AttachFilterDataToSpan(span, filter.Page, filter.Limit, filter.SortBy)

	// determine user ID.
	sessionCtxData, err := s.sessionContextDataFetcher(req)
	if err != nil {
		logger.Error(err, "retrieving session context data")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, "unauthenticated", http.StatusUnauthorized)
		return
	}

	tracing.AttachSessionContextDataToSpan(span, sessionCtxData)
	logger = sessionCtxData.AttachToLogger(logger)

	validIngredients, err := s.validIngredientDataManager.SearchForValidIngredients(ctx, query)
	if errors.Is(err, sql.ErrNoRows) {
		// in the event no rows exist, return an empty list.
		validIngredients = []*types.ValidIngredient{}
	} else if err != nil {
		observability.AcknowledgeError(err, logger, span, "searching valid ingredients")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	// encode our response and peace.
	s.encoderDecoder.RespondWithData(ctx, res, validIngredients)
}

// UpdateHandler returns a handler that updates a valid ingredient.
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
	input := new(types.ValidIngredientUpdateRequestInput)
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

	// determine valid ingredient ID.
	validIngredientID := s.validIngredientIDFetcher(req)
	tracing.AttachValidIngredientIDToSpan(span, validIngredientID)
	logger = logger.WithValue(keys.ValidIngredientIDKey, validIngredientID)

	// fetch valid ingredient from database.
	validIngredient, err := s.validIngredientDataManager.GetValidIngredient(ctx, validIngredientID)
	if errors.Is(err, sql.ErrNoRows) {
		s.encoderDecoder.EncodeNotFoundResponse(ctx, res)
		return
	} else if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving valid ingredient for update")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	// update the valid ingredient.
	validIngredient.Update(input)

	if err = s.validIngredientDataManager.UpdateValidIngredient(ctx, validIngredient); err != nil {
		observability.AcknowledgeError(err, logger, span, "updating valid ingredient")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	if s.dataChangesPublisher != nil {
		dcm := &types.DataChangeMessage{
			DataType:             types.ValidIngredientDataType,
			EventType:            types.ValidIngredientUpdatedCustomerEventType,
			ValidIngredient:      validIngredient,
			AttributableToUserID: sessionCtxData.Requester.UserID,
		}

		if err = s.dataChangesPublisher.Publish(ctx, dcm); err != nil {
			observability.AcknowledgeError(err, logger, span, "publishing data change message")
		}
	}

	// encode our response and peace.
	s.encoderDecoder.RespondWithData(ctx, res, validIngredient)
}

// ArchiveHandler returns a handler that archives a valid ingredient.
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

	// determine valid ingredient ID.
	validIngredientID := s.validIngredientIDFetcher(req)
	tracing.AttachValidIngredientIDToSpan(span, validIngredientID)
	logger = logger.WithValue(keys.ValidIngredientIDKey, validIngredientID)

	exists, existenceCheckErr := s.validIngredientDataManager.ValidIngredientExists(ctx, validIngredientID)
	if existenceCheckErr != nil && !errors.Is(existenceCheckErr, sql.ErrNoRows) {
		observability.AcknowledgeError(existenceCheckErr, logger, span, "checking valid ingredient existence")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	} else if !exists || errors.Is(existenceCheckErr, sql.ErrNoRows) {
		s.encoderDecoder.EncodeNotFoundResponse(ctx, res)
		return
	}

	if err = s.validIngredientDataManager.ArchiveValidIngredient(ctx, validIngredientID); err != nil {
		observability.AcknowledgeError(err, logger, span, "updating valid ingredient")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	if s.dataChangesPublisher != nil {
		dcm := &types.DataChangeMessage{
			DataType:             types.ValidIngredientDataType,
			EventType:            types.ValidIngredientArchivedCustomerEventType,
			AttributableToUserID: sessionCtxData.Requester.UserID,
		}

		if err = s.dataChangesPublisher.Publish(ctx, dcm); err != nil {
			observability.AcknowledgeError(err, logger, span, "publishing data change message")
		}
	}

	// encode our response and peace.
	res.WriteHeader(http.StatusNoContent)
}

// RandomHandler returns a GET handler that returns a random valid ingredient.
func (s *service) RandomHandler(res http.ResponseWriter, req *http.Request) {
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

	// fetch valid ingredient from database.
	x, err := s.validIngredientDataManager.GetRandomValidIngredient(ctx)
	if errors.Is(err, sql.ErrNoRows) {
		s.encoderDecoder.EncodeNotFoundResponse(ctx, res)
		return
	} else if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving valid ingredient")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	// encode our response and peace.
	s.encoderDecoder.RespondWithData(ctx, res, x)
}
