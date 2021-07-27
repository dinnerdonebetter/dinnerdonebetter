package validpreparations

import (
	"context"
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
	// ValidPreparationIDURIParamKey is a standard string that we'll use to refer to valid preparation IDs with.
	ValidPreparationIDURIParamKey = "validPreparationID"
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

// CreateHandler is our valid preparation creation route.
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
	input := new(types.ValidPreparationCreationInput)
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

	// create valid preparation in database.
	validPreparation, err := s.validPreparationDataManager.CreateValidPreparation(ctx, input, sessionCtxData.Requester.UserID)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "creating valid preparation")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	tracing.AttachValidPreparationIDToSpan(span, validPreparation.ID)
	logger = logger.WithValue(keys.ValidPreparationIDKey, validPreparation.ID)

	// notify interested parties.
	if searchIndexErr := s.search.Index(ctx, validPreparation.ID, validPreparation); searchIndexErr != nil {
		observability.AcknowledgeError(err, logger, span, "adding valid preparation to search index")
	}
	s.validPreparationCounter.Increment(ctx)

	s.encoderDecoder.EncodeResponseWithStatus(ctx, res, validPreparation, http.StatusCreated)
}

// ReadHandler returns a GET handler that returns a valid preparation.
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

	// determine valid preparation ID.
	validPreparationID := s.validPreparationIDFetcher(req)
	tracing.AttachValidPreparationIDToSpan(span, validPreparationID)
	logger = logger.WithValue(keys.ValidPreparationIDKey, validPreparationID)

	// fetch valid preparation from database.
	x, err := s.validPreparationDataManager.GetValidPreparation(ctx, validPreparationID)
	if errors.Is(err, sql.ErrNoRows) {
		s.encoderDecoder.EncodeNotFoundResponse(ctx, res)
		return
	} else if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving valid preparation")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	// encode our response and peace.
	s.encoderDecoder.RespondWithData(ctx, res, x)
}

// ExistenceHandler returns a HEAD handler that returns 200 if a valid preparation exists, 404 otherwise.
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

	// determine valid preparation ID.
	validPreparationID := s.validPreparationIDFetcher(req)
	tracing.AttachValidPreparationIDToSpan(span, validPreparationID)
	logger = logger.WithValue(keys.ValidPreparationIDKey, validPreparationID)

	// check the database.
	exists, err := s.validPreparationDataManager.ValidPreparationExists(ctx, validPreparationID)
	if !errors.Is(err, sql.ErrNoRows) {
		observability.AcknowledgeError(err, logger, span, "checking valid preparation existence")
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

	validPreparations, err := s.validPreparationDataManager.GetValidPreparations(ctx, filter)
	if errors.Is(err, sql.ErrNoRows) {
		// in the event no rows exist, return an empty list.
		validPreparations = &types.ValidPreparationList{ValidPreparations: []*types.ValidPreparation{}}
	} else if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving valid preparations")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	// encode our response and peace.
	s.encoderDecoder.RespondWithData(ctx, res, validPreparations)
}

// SearchForValidPreparations handles searching for and retrieving ValidPreparations.
func (s *service) SearchForValidPreparations(ctx context.Context, sessionCtxData *types.SessionContextData, query string, filter *types.QueryFilter) ([]*types.ValidPreparation, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithValue(keys.SearchQueryKey, query)
	tracing.AttachSearchQueryToSpan(span, query)

	relevantIDs, err := s.search.Search(ctx, query, sessionCtxData.ActiveAccountID)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "executing valid ingredient search query")
	}

	// fetch valid ingredients from database.
	return s.validPreparationDataManager.GetValidPreparationsWithIDs(ctx, filter.Limit, relevantIDs)
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

	relevantIDs, err := s.search.Search(ctx, query, sessionCtxData.ActiveAccountID)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "executing valid preparation search query")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	// fetch valid preparations from database.
	validPreparations, err := s.validPreparationDataManager.GetValidPreparationsWithIDs(ctx, filter.Limit, relevantIDs)
	if errors.Is(err, sql.ErrNoRows) {
		// in the event no rows exist, return an empty list.
		validPreparations = []*types.ValidPreparation{}
	} else if err != nil {
		observability.AcknowledgeError(err, logger, span, "searching valid preparations")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	// encode our response and peace.
	s.encoderDecoder.RespondWithData(ctx, res, validPreparations)
}

// UpdateHandler returns a handler that updates a valid preparation.
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
	input := new(types.ValidPreparationUpdateInput)
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

	// determine valid preparation ID.
	validPreparationID := s.validPreparationIDFetcher(req)
	tracing.AttachValidPreparationIDToSpan(span, validPreparationID)
	logger = logger.WithValue(keys.ValidPreparationIDKey, validPreparationID)

	// fetch valid preparation from database.
	validPreparation, err := s.validPreparationDataManager.GetValidPreparation(ctx, validPreparationID)
	if errors.Is(err, sql.ErrNoRows) {
		s.encoderDecoder.EncodeNotFoundResponse(ctx, res)
		return
	} else if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving valid preparation for update")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	// update the valid preparation.
	changeReport := validPreparation.Update(input)
	tracing.AttachChangeSummarySpan(span, "valid_preparation", changeReport)

	// update valid preparation in database.
	if err = s.validPreparationDataManager.UpdateValidPreparation(ctx, validPreparation, sessionCtxData.Requester.UserID, changeReport); err != nil {
		observability.AcknowledgeError(err, logger, span, "updating valid preparation")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	// notify interested parties.
	if searchIndexErr := s.search.Index(ctx, validPreparation.ID, validPreparation); searchIndexErr != nil {
		observability.AcknowledgeError(err, logger, span, "updating valid preparation in search index")
	}

	// encode our response and peace.
	s.encoderDecoder.RespondWithData(ctx, res, validPreparation)
}

// ArchiveHandler returns a handler that archives a valid preparation.
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

	// determine valid preparation ID.
	validPreparationID := s.validPreparationIDFetcher(req)
	tracing.AttachValidPreparationIDToSpan(span, validPreparationID)
	logger = logger.WithValue(keys.ValidPreparationIDKey, validPreparationID)

	// archive the valid preparation in the database.
	err = s.validPreparationDataManager.ArchiveValidPreparation(ctx, validPreparationID, sessionCtxData.Requester.UserID)
	if errors.Is(err, sql.ErrNoRows) {
		s.encoderDecoder.EncodeNotFoundResponse(ctx, res)
		return
	} else if err != nil {
		observability.AcknowledgeError(err, logger, span, "archiving valid preparation")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	// notify interested parties.
	s.validPreparationCounter.Decrement(ctx)

	if indexDeleteErr := s.search.Delete(ctx, validPreparationID); indexDeleteErr != nil {
		observability.AcknowledgeError(err, logger, span, "removing from search index")
	}

	// encode our response and peace.
	res.WriteHeader(http.StatusNoContent)
}

// AuditEntryHandler returns a GET handler that returns all audit log entries related to a valid preparation.
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

	// determine valid preparation ID.
	validPreparationID := s.validPreparationIDFetcher(req)
	tracing.AttachValidPreparationIDToSpan(span, validPreparationID)
	logger = logger.WithValue(keys.ValidPreparationIDKey, validPreparationID)

	x, err := s.validPreparationDataManager.GetAuditLogEntriesForValidPreparation(ctx, validPreparationID)
	if errors.Is(err, sql.ErrNoRows) {
		s.encoderDecoder.EncodeNotFoundResponse(ctx, res)
		return
	} else if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving audit log entries for valid preparation")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	// encode our response and peace.
	s.encoderDecoder.RespondWithData(ctx, res, x)
}
