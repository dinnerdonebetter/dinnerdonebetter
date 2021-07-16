package validingredientpreparations

import (
	"database/sql"
	"errors"
	"net/http"
	"strings"

	observability "gitlab.com/prixfixe/prixfixe/internal/observability"
	keys "gitlab.com/prixfixe/prixfixe/internal/observability/keys"
	"gitlab.com/prixfixe/prixfixe/internal/observability/tracing"
	"gitlab.com/prixfixe/prixfixe/pkg/types"
)

const (
	// ValidIngredientPreparationIDURIParamKey is a standard string that we'll use to refer to valid ingredient preparation IDs with.
	ValidIngredientPreparationIDURIParamKey = "validIngredientPreparationID"
)

// parseBool differs from strconv.ParseBool in that it returns false by default.
func parseBool(str string) bool {
	switch strings.ToLower(strings.TrimSpace(str)) {
	case "1", "t", "true":
		return true
	default:
		return false
	}
}

// CreateHandler is our valid ingredient preparation creation route.
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

	// check session context data for parsed input struct.
	input := new(types.ValidIngredientPreparationCreationInput)
	if err = s.encoderDecoder.DecodeRequest(ctx, req, input); err != nil {
		observability.AcknowledgeError(err, logger, span, "decoding request")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, "invalid request content", http.StatusBadRequest)
		return
	}

	if err = input.ValidateWithContext(ctx); err != nil {
		logger.WithValue(keys.ValidationErrorKey, err).Debug("invalid input attached to request")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, err.Error(), http.StatusBadRequest)
		return
	}

	// create valid ingredient preparation in database.
	validIngredientPreparation, err := s.validIngredientPreparationDataManager.CreateValidIngredientPreparation(ctx, input, sessionCtxData.Requester.UserID)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "creating valid ingredient preparation")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	tracing.AttachValidIngredientPreparationIDToSpan(span, validIngredientPreparation.ID)

	// notify interested parties.
	s.validIngredientPreparationCounter.Increment(ctx)

	s.encoderDecoder.EncodeResponseWithStatus(ctx, res, validIngredientPreparation, http.StatusCreated)
}

// ReadHandler returns a GET handler that returns a valid ingredient preparation.
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

	// determine valid ingredient preparation ID.
	validIngredientPreparationID := s.validIngredientPreparationIDFetcher(req)
	tracing.AttachValidIngredientPreparationIDToSpan(span, validIngredientPreparationID)
	logger = logger.WithValue(keys.ValidIngredientPreparationIDKey, validIngredientPreparationID)

	// fetch valid ingredient preparation from database.
	x, err := s.validIngredientPreparationDataManager.GetValidIngredientPreparation(ctx, validIngredientPreparationID)
	if errors.Is(err, sql.ErrNoRows) {
		s.encoderDecoder.EncodeNotFoundResponse(ctx, res)
		return
	} else if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving valid ingredient preparation")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	// encode our response and peace.
	s.encoderDecoder.RespondWithData(ctx, res, x)
}

// ExistenceHandler returns a HEAD handler that returns 200 if a valid ingredient preparation exists, 404 otherwise.
func (s *service) ExistenceHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	logger := s.logger.WithRequest(req)
	tracing.AttachRequestToSpan(span, req)

	// determine user ID.
	sessionCtxData, err := s.sessionContextDataFetcher(req)
	if err != nil {
		s.logger.Error(err, "retrieving session context data")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, "unauthenticated", http.StatusUnauthorized)
		return
	}

	tracing.AttachSessionContextDataToSpan(span, sessionCtxData)
	logger = sessionCtxData.AttachToLogger(logger)

	// determine valid ingredient preparation ID.
	validIngredientPreparationID := s.validIngredientPreparationIDFetcher(req)
	tracing.AttachValidIngredientPreparationIDToSpan(span, validIngredientPreparationID)
	logger = logger.WithValue(keys.ValidIngredientPreparationIDKey, validIngredientPreparationID)

	// check the database.
	exists, err := s.validIngredientPreparationDataManager.ValidIngredientPreparationExists(ctx, validIngredientPreparationID)
	if !errors.Is(err, sql.ErrNoRows) {
		observability.AcknowledgeError(err, logger, span, "checking valid ingredient preparation existence")
	}

	if !exists || errors.Is(err, sql.ErrNoRows) {
		s.encoderDecoder.EncodeNotFoundResponse(ctx, res)
	}
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

	validIngredientPreparations, err := s.validIngredientPreparationDataManager.GetValidIngredientPreparations(ctx, filter)
	if errors.Is(err, sql.ErrNoRows) {
		// in the event no rows exist, return an empty list.
		validIngredientPreparations = &types.ValidIngredientPreparationList{ValidIngredientPreparations: []*types.ValidIngredientPreparation{}}
	} else if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving valid ingredient preparations")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	// encode our response and peace.
	s.encoderDecoder.RespondWithData(ctx, res, validIngredientPreparations)
}

// UpdateHandler returns a handler that updates a valid ingredient preparation.
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
	input := new(types.ValidIngredientPreparationUpdateInput)
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

	// determine valid ingredient preparation ID.
	validIngredientPreparationID := s.validIngredientPreparationIDFetcher(req)
	tracing.AttachValidIngredientPreparationIDToSpan(span, validIngredientPreparationID)
	logger = logger.WithValue(keys.ValidIngredientPreparationIDKey, validIngredientPreparationID)

	// fetch valid ingredient preparation from database.
	validIngredientPreparation, err := s.validIngredientPreparationDataManager.GetValidIngredientPreparation(ctx, validIngredientPreparationID)
	if errors.Is(err, sql.ErrNoRows) {
		s.encoderDecoder.EncodeNotFoundResponse(ctx, res)
		return
	} else if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving valid ingredient preparation for update")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	// update the valid ingredient preparation.
	changeReport := validIngredientPreparation.Update(input)
	tracing.AttachChangeSummarySpan(span, "valid_ingredient_preparation", changeReport)

	// update valid ingredient preparation in database.
	if err = s.validIngredientPreparationDataManager.UpdateValidIngredientPreparation(ctx, validIngredientPreparation, sessionCtxData.Requester.UserID, changeReport); err != nil {
		observability.AcknowledgeError(err, logger, span, "updating valid ingredient preparation")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	// notify interested parties.

	// encode our response and peace.
	s.encoderDecoder.RespondWithData(ctx, res, validIngredientPreparation)
}

// ArchiveHandler returns a handler that archives a valid ingredient preparation.
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

	// determine valid ingredient preparation ID.
	validIngredientPreparationID := s.validIngredientPreparationIDFetcher(req)
	tracing.AttachValidIngredientPreparationIDToSpan(span, validIngredientPreparationID)
	logger = logger.WithValue(keys.ValidIngredientPreparationIDKey, validIngredientPreparationID)

	// archive the valid ingredient preparation in the database.
	err = s.validIngredientPreparationDataManager.ArchiveValidIngredientPreparation(ctx, validIngredientPreparationID, sessionCtxData.Requester.UserID)
	if errors.Is(err, sql.ErrNoRows) {
		s.encoderDecoder.EncodeNotFoundResponse(ctx, res)
		return
	} else if err != nil {
		observability.AcknowledgeError(err, logger, span, "archiving valid ingredient preparation")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	// notify interested parties.
	s.validIngredientPreparationCounter.Decrement(ctx)

	// encode our response and peace.
	res.WriteHeader(http.StatusNoContent)
}

// AuditEntryHandler returns a GET handler that returns all audit log entries related to a valid ingredient preparation.
func (s *service) AuditEntryHandler(res http.ResponseWriter, req *http.Request) {
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

	// determine valid ingredient preparation ID.
	validIngredientPreparationID := s.validIngredientPreparationIDFetcher(req)
	tracing.AttachValidIngredientPreparationIDToSpan(span, validIngredientPreparationID)
	logger = logger.WithValue(keys.ValidIngredientPreparationIDKey, validIngredientPreparationID)

	x, err := s.validIngredientPreparationDataManager.GetAuditLogEntriesForValidIngredientPreparation(ctx, validIngredientPreparationID)
	if errors.Is(err, sql.ErrNoRows) {
		s.encoderDecoder.EncodeNotFoundResponse(ctx, res)
		return
	} else if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving audit log entries for valid ingredient preparation")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	// encode our response and peace.
	s.encoderDecoder.RespondWithData(ctx, res, x)
}
