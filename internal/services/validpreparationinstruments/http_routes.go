package validpreparationinstruments

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
	// ValidPreparationInstrumentIDURIParamKey is a standard string that we'll use to refer to valid preparation instrument IDs with.
	ValidPreparationInstrumentIDURIParamKey = "validPreparationInstrumentID"
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

// CreateHandler is our valid preparation instrument creation route.
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
	input := new(types.ValidPreparationInstrumentCreationInput)
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

	// create valid preparation instrument in database.
	validPreparationInstrument, err := s.validPreparationInstrumentDataManager.CreateValidPreparationInstrument(ctx, input, sessionCtxData.Requester.UserID)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "creating valid preparation instrument")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	tracing.AttachValidPreparationInstrumentIDToSpan(span, validPreparationInstrument.ID)

	// notify interested parties.
	s.validPreparationInstrumentCounter.Increment(ctx)

	s.encoderDecoder.EncodeResponseWithStatus(ctx, res, validPreparationInstrument, http.StatusCreated)
}

// ReadHandler returns a GET handler that returns a valid preparation instrument.
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

	// determine valid preparation instrument ID.
	validPreparationInstrumentID := s.validPreparationInstrumentIDFetcher(req)
	tracing.AttachValidPreparationInstrumentIDToSpan(span, validPreparationInstrumentID)
	logger = logger.WithValue(keys.ValidPreparationInstrumentIDKey, validPreparationInstrumentID)

	// fetch valid preparation instrument from database.
	x, err := s.validPreparationInstrumentDataManager.GetValidPreparationInstrument(ctx, validPreparationInstrumentID)
	if errors.Is(err, sql.ErrNoRows) {
		s.encoderDecoder.EncodeNotFoundResponse(ctx, res)
		return
	} else if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving valid preparation instrument")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	// encode our response and peace.
	s.encoderDecoder.RespondWithData(ctx, res, x)
}

// ExistenceHandler returns a HEAD handler that returns 200 if a valid preparation instrument exists, 404 otherwise.
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

	// determine valid preparation instrument ID.
	validPreparationInstrumentID := s.validPreparationInstrumentIDFetcher(req)
	tracing.AttachValidPreparationInstrumentIDToSpan(span, validPreparationInstrumentID)
	logger = logger.WithValue(keys.ValidPreparationInstrumentIDKey, validPreparationInstrumentID)

	// check the database.
	exists, err := s.validPreparationInstrumentDataManager.ValidPreparationInstrumentExists(ctx, validPreparationInstrumentID)
	if !errors.Is(err, sql.ErrNoRows) {
		observability.AcknowledgeError(err, logger, span, "checking valid preparation instrument existence")
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

	validPreparationInstruments, err := s.validPreparationInstrumentDataManager.GetValidPreparationInstruments(ctx, filter)
	if errors.Is(err, sql.ErrNoRows) {
		// in the event no rows exist, return an empty list.
		validPreparationInstruments = &types.ValidPreparationInstrumentList{ValidPreparationInstruments: []*types.ValidPreparationInstrument{}}
	} else if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving valid preparation instruments")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	// encode our response and peace.
	s.encoderDecoder.RespondWithData(ctx, res, validPreparationInstruments)
}

// UpdateHandler returns a handler that updates a valid preparation instrument.
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
	input := new(types.ValidPreparationInstrumentUpdateInput)
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

	// determine valid preparation instrument ID.
	validPreparationInstrumentID := s.validPreparationInstrumentIDFetcher(req)
	tracing.AttachValidPreparationInstrumentIDToSpan(span, validPreparationInstrumentID)
	logger = logger.WithValue(keys.ValidPreparationInstrumentIDKey, validPreparationInstrumentID)

	// fetch valid preparation instrument from database.
	validPreparationInstrument, err := s.validPreparationInstrumentDataManager.GetValidPreparationInstrument(ctx, validPreparationInstrumentID)
	if errors.Is(err, sql.ErrNoRows) {
		s.encoderDecoder.EncodeNotFoundResponse(ctx, res)
		return
	} else if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving valid preparation instrument for update")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	// update the valid preparation instrument.
	changeReport := validPreparationInstrument.Update(input)
	tracing.AttachChangeSummarySpan(span, "valid_preparation_instrument", changeReport)

	// update valid preparation instrument in database.
	if err = s.validPreparationInstrumentDataManager.UpdateValidPreparationInstrument(ctx, validPreparationInstrument, sessionCtxData.Requester.UserID, changeReport); err != nil {
		observability.AcknowledgeError(err, logger, span, "updating valid preparation instrument")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	// notify interested parties.

	// encode our response and peace.
	s.encoderDecoder.RespondWithData(ctx, res, validPreparationInstrument)
}

// ArchiveHandler returns a handler that archives a valid preparation instrument.
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

	// determine valid preparation instrument ID.
	validPreparationInstrumentID := s.validPreparationInstrumentIDFetcher(req)
	tracing.AttachValidPreparationInstrumentIDToSpan(span, validPreparationInstrumentID)
	logger = logger.WithValue(keys.ValidPreparationInstrumentIDKey, validPreparationInstrumentID)

	// archive the valid preparation instrument in the database.
	err = s.validPreparationInstrumentDataManager.ArchiveValidPreparationInstrument(ctx, validPreparationInstrumentID, sessionCtxData.Requester.UserID)
	if errors.Is(err, sql.ErrNoRows) {
		s.encoderDecoder.EncodeNotFoundResponse(ctx, res)
		return
	} else if err != nil {
		observability.AcknowledgeError(err, logger, span, "archiving valid preparation instrument")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	// notify interested parties.
	s.validPreparationInstrumentCounter.Decrement(ctx)

	// encode our response and peace.
	res.WriteHeader(http.StatusNoContent)
}

// AuditEntryHandler returns a GET handler that returns all audit log entries related to a valid preparation instrument.
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

	// determine valid preparation instrument ID.
	validPreparationInstrumentID := s.validPreparationInstrumentIDFetcher(req)
	tracing.AttachValidPreparationInstrumentIDToSpan(span, validPreparationInstrumentID)
	logger = logger.WithValue(keys.ValidPreparationInstrumentIDKey, validPreparationInstrumentID)

	x, err := s.validPreparationInstrumentDataManager.GetAuditLogEntriesForValidPreparationInstrument(ctx, validPreparationInstrumentID)
	if errors.Is(err, sql.ErrNoRows) {
		s.encoderDecoder.EncodeNotFoundResponse(ctx, res)
		return
	} else if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving audit log entries for valid preparation instrument")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	// encode our response and peace.
	s.encoderDecoder.RespondWithData(ctx, res, x)
}
