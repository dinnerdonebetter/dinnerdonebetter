package validinstruments

import (
	"database/sql"
	"errors"
	"net/http"
	"strings"

	"github.com/segmentio/ksuid"

	observability "github.com/prixfixeco/api_server/internal/observability"
	keys "github.com/prixfixeco/api_server/internal/observability/keys"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/pkg/types"
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

	// read parsed input struct from request body.
	providedInput := new(types.ValidInstrumentCreationRequestInput)
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

	input := types.ValidInstrumentDatabaseCreationInputFromValidInstrumentCreationInput(providedInput)
	input.ID = ksuid.New().String()

	tracing.AttachValidInstrumentIDToSpan(span, input.ID)

	// create valid instrument in database.
	preWrite := &types.PreWriteMessage{
		DataType:                  types.ValidInstrumentDataType,
		ValidInstrument:           input,
		AttributableToUserID:      sessionCtxData.Requester.UserID,
		AttributableToHouseholdID: sessionCtxData.ActiveHouseholdID,
	}
	if err = s.preWritesPublisher.Publish(ctx, preWrite); err != nil {
		observability.AcknowledgeError(err, logger, span, "publishing valid instrument write message")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	pwr := types.PreWriteResponse{ID: input.ID}

	s.encoderDecoder.EncodeResponseWithStatus(ctx, res, pwr, http.StatusAccepted)
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
		logger.Error(err, "retrieving session context data")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, "unauthenticated", http.StatusUnauthorized)
		return
	}

	tracing.AttachSessionContextDataToSpan(span, sessionCtxData)
	logger = sessionCtxData.AttachToLogger(logger)

	relevantIDs, err := s.search.Search(ctx, "name", query, "")
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "executing valid instrument search query")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	// fetch valid instruments from database.
	validInstruments, err := s.validInstrumentDataManager.GetValidInstrumentsWithIDs(ctx, filter.Limit, relevantIDs)
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
	input := new(types.ValidInstrumentUpdateRequestInput)
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
	validInstrument.Update(input)

	pum := &types.PreUpdateMessage{
		DataType:                  types.ValidInstrumentDataType,
		ValidInstrument:           validInstrument,
		AttributableToUserID:      sessionCtxData.Requester.UserID,
		AttributableToHouseholdID: sessionCtxData.ActiveHouseholdID,
	}
	if err = s.preUpdatesPublisher.Publish(ctx, pum); err != nil {
		observability.AcknowledgeError(err, logger, span, "publishing valid instrument update message")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
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

	exists, existenceCheckErr := s.validInstrumentDataManager.ValidInstrumentExists(ctx, validInstrumentID)
	if existenceCheckErr != nil && !errors.Is(existenceCheckErr, sql.ErrNoRows) {
		observability.AcknowledgeError(existenceCheckErr, logger, span, "checking valid instrument existence")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	} else if !exists || errors.Is(existenceCheckErr, sql.ErrNoRows) {
		s.encoderDecoder.EncodeNotFoundResponse(ctx, res)
		return
	}

	pam := &types.PreArchiveMessage{
		DataType:                  types.ValidInstrumentDataType,
		ValidInstrumentID:         validInstrumentID,
		AttributableToUserID:      sessionCtxData.Requester.UserID,
		AttributableToHouseholdID: sessionCtxData.ActiveHouseholdID,
	}
	if err = s.preArchivesPublisher.Publish(ctx, pam); err != nil {
		observability.AcknowledgeError(err, logger, span, "publishing valid instrument archive message")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	// encode our response and peace.
	res.WriteHeader(http.StatusNoContent)
}
