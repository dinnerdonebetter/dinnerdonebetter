package validinstruments

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"strings"

	"gitlab.com/prixfixe/prixfixe/internal/database"

	observability "gitlab.com/prixfixe/prixfixe/internal/observability"
	keys "gitlab.com/prixfixe/prixfixe/internal/observability/keys"
	"gitlab.com/prixfixe/prixfixe/internal/observability/tracing"
	"gitlab.com/prixfixe/prixfixe/pkg/types"
)

const (
	// ValidInstrumentIDURIParamKey is a standard string that we'll use to refer to valid instrument IDs with.
	ValidInstrumentIDURIParamKey = "validInstrumentID"
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

// CreateHandler is our valid instrument creation route.
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
	input := new(types.ValidInstrumentCreationInput)
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

	// create valid instrument in database.
	validInstrument, err := s.validInstrumentDataManager.CreateValidInstrument(ctx, input, sessionCtxData.Requester.UserID)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "creating valid instrument")

		if errors.Is(err, database.ErrUniqueConstraintViolation) {
			s.encoderDecoder.EncodeRejectedDuplicateResponse(ctx, res)
		} else {
			s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		}

		return
	}

	tracing.AttachValidInstrumentIDToSpan(span, validInstrument.ID)
	logger = logger.WithValue(keys.ValidInstrumentIDKey, validInstrument.ID)

	// notify interested parties.
	if searchIndexErr := s.search.Index(ctx, validInstrument.ID, validInstrument); searchIndexErr != nil {
		observability.AcknowledgeError(err, logger, span, "adding valid instrument to search index")
	}
	s.validInstrumentCounter.Increment(ctx)

	s.encoderDecoder.EncodeResponseWithStatus(ctx, res, validInstrument, http.StatusCreated)
}

// ReadHandler returns a GET handler that returns a valid instrument.
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

	// determine valid instrument ID.
	validInstrumentID := s.validInstrumentIDFetcher(req)
	tracing.AttachValidInstrumentIDToSpan(span, validInstrumentID)
	logger = logger.WithValue(keys.ValidInstrumentIDKey, validInstrumentID)

	// fetch valid instrument from database.
	x, err := s.validInstrumentDataManager.GetValidInstrument(ctx, validInstrumentID)
	if errors.Is(err, sql.ErrNoRows) {
		s.encoderDecoder.EncodeNotFoundResponse(ctx, res)
		return
	} else if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving valid instrument")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	// encode our response and peace.
	s.encoderDecoder.RespondWithData(ctx, res, x)
}

// ExistenceHandler returns a HEAD handler that returns 200 if a valid instrument exists, 404 otherwise.
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

	// determine valid instrument ID.
	validInstrumentID := s.validInstrumentIDFetcher(req)
	tracing.AttachValidInstrumentIDToSpan(span, validInstrumentID)
	logger = logger.WithValue(keys.ValidInstrumentIDKey, validInstrumentID)

	// check the database.
	exists, err := s.validInstrumentDataManager.ValidInstrumentExists(ctx, validInstrumentID)
	if !errors.Is(err, sql.ErrNoRows) {
		observability.AcknowledgeError(err, logger, span, "checking valid instrument existence")
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

	validInstruments, err := s.validInstrumentDataManager.GetValidInstruments(ctx, filter)
	if errors.Is(err, sql.ErrNoRows) {
		// in the event no rows exist, return an empty list.
		validInstruments = &types.ValidInstrumentList{ValidInstruments: []*types.ValidInstrument{}}
	} else if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving valid instruments")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	// encode our response and peace.
	s.encoderDecoder.RespondWithData(ctx, res, validInstruments)
}

// SearchForValidInstruments handles searching for and retrieving ValidInstruments.
func (s *service) SearchForValidInstruments(ctx context.Context, sessionCtxData *types.SessionContextData, query string, filter *types.QueryFilter) ([]*types.ValidInstrument, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithValue(keys.SearchQueryKey, query)
	tracing.AttachSearchQueryToSpan(span, query)

	relevantIDs, err := s.search.Search(ctx, query, sessionCtxData.ActiveAccountID)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "executing valid ingredient search query")
	}

	// fetch valid ingredients from database.
	return s.validInstrumentDataManager.GetValidInstrumentsWithIDs(ctx, filter.Limit, relevantIDs)
}

// SearchHandler is our search route.
func (s *service) SearchHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	query := req.URL.Query().Get(types.SearchQueryKey)
	filter := types.ExtractQueryFilter(req)
	logger := s.logger.WithRequest(req).
		WithValue(keys.FilterLimitKey, filter.Limit).
		WithValue(keys.FilterPageKey, filter.Page).
		WithValue(keys.FilterSortByKey, string(filter.SortBy)).
		WithValue(keys.SearchQueryKey, query)

	tracing.AttachRequestToSpan(span, req)
	tracing.AttachFilterToSpan(span, filter.Page, filter.Limit, string(filter.SortBy))

	// determine user ID.
	sessionCtxData, err := s.sessionContextDataFetcher(req)
	if err != nil {
		s.logger.Error(err, "retrieving session context data")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, "unauthenticated", http.StatusUnauthorized)
		return
	}

	tracing.AttachSessionContextDataToSpan(span, sessionCtxData)
	logger = sessionCtxData.AttachToLogger(logger)

	// fetch valid instruments from database.
	validInstruments, err := s.SearchForValidInstruments(ctx, sessionCtxData, query, filter)
	if errors.Is(err, sql.ErrNoRows) {
		// in the event no rows exist, return an empty list.
		validInstruments = []*types.ValidInstrument{}
	} else if err != nil {
		observability.AcknowledgeError(err, logger, span, "searching valid instruments")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	// encode our response and peace.
	s.encoderDecoder.RespondWithData(ctx, res, validInstruments)
}

// UpdateHandler returns a handler that updates a valid instrument.
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
	input := new(types.ValidInstrumentUpdateInput)
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

	// determine valid instrument ID.
	validInstrumentID := s.validInstrumentIDFetcher(req)
	tracing.AttachValidInstrumentIDToSpan(span, validInstrumentID)
	logger = logger.WithValue(keys.ValidInstrumentIDKey, validInstrumentID)

	// fetch valid instrument from database.
	validInstrument, err := s.validInstrumentDataManager.GetValidInstrument(ctx, validInstrumentID)
	if errors.Is(err, sql.ErrNoRows) {
		s.encoderDecoder.EncodeNotFoundResponse(ctx, res)
		return
	} else if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving valid instrument for update")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	// update the valid instrument.
	changeReport := validInstrument.Update(input)
	tracing.AttachChangeSummarySpan(span, "valid_instrument", changeReport)

	// update valid instrument in database.
	if err = s.validInstrumentDataManager.UpdateValidInstrument(ctx, validInstrument, sessionCtxData.Requester.UserID, changeReport); err != nil {
		observability.AcknowledgeError(err, logger, span, "updating valid instrument")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	// notify interested parties.
	if searchIndexErr := s.search.Index(ctx, validInstrument.ID, validInstrument); searchIndexErr != nil {
		observability.AcknowledgeError(err, logger, span, "updating valid instrument in search index")
	}

	// encode our response and peace.
	s.encoderDecoder.RespondWithData(ctx, res, validInstrument)
}

// ArchiveHandler returns a handler that archives a valid instrument.
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

	// determine valid instrument ID.
	validInstrumentID := s.validInstrumentIDFetcher(req)
	tracing.AttachValidInstrumentIDToSpan(span, validInstrumentID)
	logger = logger.WithValue(keys.ValidInstrumentIDKey, validInstrumentID)

	// archive the valid instrument in the database.
	err = s.validInstrumentDataManager.ArchiveValidInstrument(ctx, validInstrumentID, sessionCtxData.Requester.UserID)
	if errors.Is(err, sql.ErrNoRows) {
		s.encoderDecoder.EncodeNotFoundResponse(ctx, res)
		return
	} else if err != nil {
		observability.AcknowledgeError(err, logger, span, "archiving valid instrument")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	// notify interested parties.
	s.validInstrumentCounter.Decrement(ctx)

	if indexDeleteErr := s.search.Delete(ctx, validInstrumentID); indexDeleteErr != nil {
		observability.AcknowledgeError(err, logger, span, "removing from search index")
	}

	// encode our response and peace.
	res.WriteHeader(http.StatusNoContent)
}

// AuditEntryHandler returns a GET handler that returns all audit log entries related to a valid instrument.
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

	// determine valid instrument ID.
	validInstrumentID := s.validInstrumentIDFetcher(req)
	tracing.AttachValidInstrumentIDToSpan(span, validInstrumentID)
	logger = logger.WithValue(keys.ValidInstrumentIDKey, validInstrumentID)

	x, err := s.validInstrumentDataManager.GetAuditLogEntriesForValidInstrument(ctx, validInstrumentID)
	if errors.Is(err, sql.ErrNoRows) {
		s.encoderDecoder.EncodeNotFoundResponse(ctx, res)
		return
	} else if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving audit log entries for valid instrument")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	// encode our response and peace.
	s.encoderDecoder.RespondWithData(ctx, res, x)
}
