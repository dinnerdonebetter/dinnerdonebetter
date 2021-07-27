package validingredients

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
	// ValidIngredientIDURIParamKey is a standard string that we'll use to refer to valid ingredient IDs with.
	ValidIngredientIDURIParamKey = "validIngredientID"
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

	// check session context data for parsed input struct.
	input := new(types.ValidIngredientCreationInput)
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

	// create valid ingredient in database.
	validIngredient, err := s.validIngredientDataManager.CreateValidIngredient(ctx, input, sessionCtxData.Requester.UserID)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "creating valid ingredient")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	tracing.AttachValidIngredientIDToSpan(span, validIngredient.ID)
	logger = logger.WithValue(keys.ValidIngredientIDKey, validIngredient.ID)

	// notify interested parties.
	if searchIndexErr := s.search.Index(ctx, validIngredient.ID, validIngredient); searchIndexErr != nil {
		observability.AcknowledgeError(err, logger, span, "adding valid ingredient to search index")
	}
	s.validIngredientCounter.Increment(ctx)

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

// ExistenceHandler returns a HEAD handler that returns 200 if a valid ingredient exists, 404 otherwise.
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

	// determine valid ingredient ID.
	validIngredientID := s.validIngredientIDFetcher(req)
	tracing.AttachValidIngredientIDToSpan(span, validIngredientID)
	logger = logger.WithValue(keys.ValidIngredientIDKey, validIngredientID)

	// check the database.
	exists, err := s.validIngredientDataManager.ValidIngredientExists(ctx, validIngredientID)
	if !errors.Is(err, sql.ErrNoRows) {
		observability.AcknowledgeError(err, logger, span, "checking valid ingredient existence")
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

func (s *service) SearchForValidIngredients(ctx context.Context, sessionCtxData *types.SessionContextData, query string, filter *types.QueryFilter) ([]*types.ValidIngredient, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithValue(keys.SearchQueryKey, query)
	tracing.AttachSearchQueryToSpan(span, query)

	relevantIDs, err := s.search.Search(ctx, query, sessionCtxData.ActiveAccountID)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "executing valid ingredient search query")
	}

	// fetch valid ingredients from database.
	return s.validIngredientDataManager.GetValidIngredientsWithIDs(ctx, filter.Limit, relevantIDs)
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
	tracing.AttachSearchQueryToSpan(span, query)
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

	// fetch valid ingredients from database.
	validIngredients, err := s.SearchForValidIngredients(ctx, sessionCtxData, query, filter)
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
	input := new(types.ValidIngredientUpdateInput)
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
	changeReport := validIngredient.Update(input)
	tracing.AttachChangeSummarySpan(span, "valid_ingredient", changeReport)

	// update valid ingredient in database.
	if err = s.validIngredientDataManager.UpdateValidIngredient(ctx, validIngredient, sessionCtxData.Requester.UserID, changeReport); err != nil {
		observability.AcknowledgeError(err, logger, span, "updating valid ingredient")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	// notify interested parties.
	if searchIndexErr := s.search.Index(ctx, validIngredient.ID, validIngredient); searchIndexErr != nil {
		observability.AcknowledgeError(err, logger, span, "updating valid ingredient in search index")
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

	// archive the valid ingredient in the database.
	err = s.validIngredientDataManager.ArchiveValidIngredient(ctx, validIngredientID, sessionCtxData.Requester.UserID)
	if errors.Is(err, sql.ErrNoRows) {
		s.encoderDecoder.EncodeNotFoundResponse(ctx, res)
		return
	} else if err != nil {
		observability.AcknowledgeError(err, logger, span, "archiving valid ingredient")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	// notify interested parties.
	s.validIngredientCounter.Decrement(ctx)

	if indexDeleteErr := s.search.Delete(ctx, validIngredientID); indexDeleteErr != nil {
		observability.AcknowledgeError(err, logger, span, "removing from search index")
	}

	// encode our response and peace.
	res.WriteHeader(http.StatusNoContent)
}

// AuditEntryHandler returns a GET handler that returns all audit log entries related to a valid ingredient.
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

	// determine valid ingredient ID.
	validIngredientID := s.validIngredientIDFetcher(req)
	tracing.AttachValidIngredientIDToSpan(span, validIngredientID)
	logger = logger.WithValue(keys.ValidIngredientIDKey, validIngredientID)

	x, err := s.validIngredientDataManager.GetAuditLogEntriesForValidIngredient(ctx, validIngredientID)
	if errors.Is(err, sql.ErrNoRows) {
		s.encoderDecoder.EncodeNotFoundResponse(ctx, res)
		return
	} else if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving audit log entries for valid ingredient")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	// encode our response and peace.
	s.encoderDecoder.RespondWithData(ctx, res, x)
}
