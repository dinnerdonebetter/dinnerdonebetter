package householdinstrumentownerships

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/dinnerdonebetter/backend/internal/identifiers"
	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"
)

const (
	// HouseholdInstrumentOwnershipIDURIParamKey is a standard string that we'll use to refer to household instrument ownership IDs with.
	HouseholdInstrumentOwnershipIDURIParamKey = "householdInstrumentOwnershipID"
)

// CreateHandler is our household instrument ownership creation route.
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
	providedInput := new(types.HouseholdInstrumentOwnershipCreationRequestInput)
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

	input := converters.ConvertHouseholdInstrumentOwnershipCreationRequestInputToHouseholdInstrumentOwnershipDatabaseCreationInput(providedInput)
	input.ID = identifiers.New()

	tracing.AttachHouseholdInstrumentOwnershipIDToSpan(span, input.ID)

	householdInstrumentOwnership, err := s.householdInstrumentOwnershipDataManager.CreateHouseholdInstrumentOwnership(ctx, input)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "creating household instrument ownership")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	dcm := &types.DataChangeMessage{
		EventType:                    types.HouseholdInstrumentOwnershipCreatedCustomerEventType,
		HouseholdInstrumentOwnership: householdInstrumentOwnership,
		UserID:                       sessionCtxData.Requester.UserID,
	}

	if err = s.dataChangesPublisher.Publish(ctx, dcm); err != nil {
		observability.AcknowledgeError(err, logger, span, "publishing to data changes topic")
	}

	s.encoderDecoder.EncodeResponseWithStatus(ctx, res, householdInstrumentOwnership, http.StatusCreated)
}

// ReadHandler returns a GET handler that returns a household instrument ownership.
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

	householdID := sessionCtxData.ActiveHouseholdID
	tracing.AttachSessionContextDataToSpan(span, sessionCtxData)
	logger = sessionCtxData.AttachToLogger(logger)

	// determine household instrument ownership ID.
	householdInstrumentOwnershipID := s.householdInstrumentOwnershipIDFetcher(req)
	tracing.AttachHouseholdInstrumentOwnershipIDToSpan(span, householdInstrumentOwnershipID)
	logger = logger.WithValue(keys.HouseholdInstrumentOwnershipIDKey, householdInstrumentOwnershipID)

	// fetch household instrument ownership from database.
	x, err := s.householdInstrumentOwnershipDataManager.GetHouseholdInstrumentOwnership(ctx, householdInstrumentOwnershipID, householdID)
	if errors.Is(err, sql.ErrNoRows) {
		s.encoderDecoder.EncodeNotFoundResponse(ctx, res)
		return
	} else if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving household instrument ownership")
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

	householdID := sessionCtxData.ActiveHouseholdID
	tracing.AttachSessionContextDataToSpan(span, sessionCtxData)
	logger = sessionCtxData.AttachToLogger(logger)

	householdInstrumentOwnerships, err := s.householdInstrumentOwnershipDataManager.GetHouseholdInstrumentOwnerships(ctx, householdID, filter)
	if errors.Is(err, sql.ErrNoRows) {
		// in the event no rows exist, return an empty list.
		householdInstrumentOwnerships = &types.QueryFilteredResult[types.HouseholdInstrumentOwnership]{Data: []*types.HouseholdInstrumentOwnership{}}
	} else if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving household instrument ownerships")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	// encode our response and peace.
	s.encoderDecoder.RespondWithData(ctx, res, householdInstrumentOwnerships)
}

// UpdateHandler returns a handler that updates a household instrument ownership.
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

	householdID := sessionCtxData.ActiveHouseholdID
	tracing.AttachSessionContextDataToSpan(span, sessionCtxData)
	logger = sessionCtxData.AttachToLogger(logger)

	// check for parsed input attached to session context data.
	input := new(types.HouseholdInstrumentOwnershipUpdateRequestInput)
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

	// determine household instrument ownership ID.
	householdInstrumentOwnershipID := s.householdInstrumentOwnershipIDFetcher(req)
	tracing.AttachHouseholdInstrumentOwnershipIDToSpan(span, householdInstrumentOwnershipID)
	logger = logger.WithValue(keys.HouseholdInstrumentOwnershipIDKey, householdInstrumentOwnershipID)

	// fetch household instrument ownership from database.
	householdInstrumentOwnership, err := s.householdInstrumentOwnershipDataManager.GetHouseholdInstrumentOwnership(ctx, householdInstrumentOwnershipID, householdID)
	if errors.Is(err, sql.ErrNoRows) {
		s.encoderDecoder.EncodeNotFoundResponse(ctx, res)
		return
	} else if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving household instrument ownership for update")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	// update the household instrument ownership.
	householdInstrumentOwnership.Update(input)

	if err = s.householdInstrumentOwnershipDataManager.UpdateHouseholdInstrumentOwnership(ctx, householdInstrumentOwnership); err != nil {
		observability.AcknowledgeError(err, logger, span, "updating household instrument ownership")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	dcm := &types.DataChangeMessage{
		EventType:                    types.HouseholdInstrumentOwnershipUpdatedCustomerEventType,
		HouseholdInstrumentOwnership: householdInstrumentOwnership,
		UserID:                       sessionCtxData.Requester.UserID,
	}

	if err = s.dataChangesPublisher.Publish(ctx, dcm); err != nil {
		observability.AcknowledgeError(err, logger, span, "publishing data change message")
	}

	// encode our response and peace.
	s.encoderDecoder.RespondWithData(ctx, res, householdInstrumentOwnership)
}

// ArchiveHandler returns a handler that archives a household instrument ownership.
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

	householdID := sessionCtxData.ActiveHouseholdID
	tracing.AttachSessionContextDataToSpan(span, sessionCtxData)
	logger = sessionCtxData.AttachToLogger(logger)

	// determine household instrument ownership ID.
	householdInstrumentOwnershipID := s.householdInstrumentOwnershipIDFetcher(req)
	tracing.AttachHouseholdInstrumentOwnershipIDToSpan(span, householdInstrumentOwnershipID)
	logger = logger.WithValue(keys.HouseholdInstrumentOwnershipIDKey, householdInstrumentOwnershipID)

	exists, existenceCheckErr := s.householdInstrumentOwnershipDataManager.HouseholdInstrumentOwnershipExists(ctx, householdInstrumentOwnershipID, householdID)
	if existenceCheckErr != nil && !errors.Is(existenceCheckErr, sql.ErrNoRows) {
		observability.AcknowledgeError(existenceCheckErr, logger, span, "checking household instrument ownership existence")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	} else if !exists || errors.Is(existenceCheckErr, sql.ErrNoRows) {
		s.encoderDecoder.EncodeNotFoundResponse(ctx, res)
		return
	}

	if err = s.householdInstrumentOwnershipDataManager.ArchiveHouseholdInstrumentOwnership(ctx, householdInstrumentOwnershipID, householdID); err != nil {
		observability.AcknowledgeError(err, logger, span, "archiving household instrument ownership")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	dcm := &types.DataChangeMessage{
		EventType: types.HouseholdInstrumentOwnershipArchivedCustomerEventType,
		UserID:    sessionCtxData.Requester.UserID,
	}

	if err = s.dataChangesPublisher.Publish(ctx, dcm); err != nil {
		observability.AcknowledgeError(err, logger, span, "publishing data change message")
	}

	// encode our response and peace.
	res.WriteHeader(http.StatusNoContent)
}
