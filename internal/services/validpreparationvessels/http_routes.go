package validpreparationvessels

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
	// ValidPreparationVesselIDURIParamKey is a standard string that we'll use to refer to valid preparation vessel IDs with.
	ValidPreparationVesselIDURIParamKey = "validPreparationVesselID"
	// ValidPreparationIDURIParamKey is a standard string that we'll use to refer to valid preparation IDs with.
	ValidPreparationIDURIParamKey = "validPreparationID"
	// ValidVesselIDURIParamKey is a standard string that we'll use to refer to valid preparation IDs with.
	ValidVesselIDURIParamKey = "ValidVesselID"
)

// CreateHandler is our valid preparation vessel creation route.
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
	providedInput := new(types.ValidPreparationVesselCreationRequestInput)
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

	input := converters.ConvertValidPreparationVesselCreationRequestInputToValidPreparationVesselDatabaseCreationInput(providedInput)
	input.ID = identifiers.New()

	tracing.AttachToSpan(span, keys.ValidVesselIDKey, input.ID)

	validPreparationVessel, err := s.validPreparationVesselDataManager.CreateValidPreparationVessel(ctx, input)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "creating valid preparation vessel")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	dcm := &types.DataChangeMessage{
		EventType:              types.ValidPreparationVesselCreatedCustomerEventType,
		ValidPreparationVessel: validPreparationVessel,
		UserID:                 sessionCtxData.Requester.UserID,
	}

	if err = s.dataChangesPublisher.Publish(ctx, dcm); err != nil {
		observability.AcknowledgeError(err, logger, span, "publishing to data changes topic")
	}

	s.encoderDecoder.EncodeResponseWithStatus(ctx, res, validPreparationVessel, http.StatusCreated)
}

// ReadHandler returns a GET handler that returns a valid preparation vessel.
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

	// determine valid preparation vessel ID.
	validPreparationVesselID := s.validPreparationVesselIDFetcher(req)
	tracing.AttachToSpan(span, keys.ValidVesselIDKey, validPreparationVesselID)
	logger = logger.WithValue(keys.ValidPreparationVesselIDKey, validPreparationVesselID)

	// fetch valid preparation vessel from database.
	x, err := s.validPreparationVesselDataManager.GetValidPreparationVessel(ctx, validPreparationVesselID)
	if errors.Is(err, sql.ErrNoRows) {
		s.encoderDecoder.EncodeNotFoundResponse(ctx, res)
		return
	} else if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving valid preparation vessel")
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

	validPreparationVessels, err := s.validPreparationVesselDataManager.GetValidPreparationVessels(ctx, filter)
	if errors.Is(err, sql.ErrNoRows) {
		// in the event no rows exist, return an empty list.
		validPreparationVessels = &types.QueryFilteredResult[types.ValidPreparationVessel]{Data: []*types.ValidPreparationVessel{}}
	} else if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving valid preparation vessels")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	// encode our response and peace.
	s.encoderDecoder.RespondWithData(ctx, res, validPreparationVessels)
}

// UpdateHandler returns a handler that updates a valid preparation vessel.
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
	input := new(types.ValidPreparationVesselUpdateRequestInput)
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

	// determine valid preparation vessel ID.
	validPreparationVesselID := s.validPreparationVesselIDFetcher(req)
	tracing.AttachToSpan(span, keys.ValidVesselIDKey, validPreparationVesselID)
	logger = logger.WithValue(keys.ValidPreparationVesselIDKey, validPreparationVesselID)

	// fetch valid preparation vessel from database.
	validPreparationVessel, err := s.validPreparationVesselDataManager.GetValidPreparationVessel(ctx, validPreparationVesselID)
	if errors.Is(err, sql.ErrNoRows) {
		s.encoderDecoder.EncodeNotFoundResponse(ctx, res)
		return
	} else if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving valid preparation vessel for update")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	// update the valid preparation vessel.
	validPreparationVessel.Update(input)

	if err = s.validPreparationVesselDataManager.UpdateValidPreparationVessel(ctx, validPreparationVessel); err != nil {
		observability.AcknowledgeError(err, logger, span, "updating valid preparation vessel")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	dcm := &types.DataChangeMessage{
		EventType:              types.ValidPreparationVesselUpdatedCustomerEventType,
		ValidPreparationVessel: validPreparationVessel,
		UserID:                 sessionCtxData.Requester.UserID,
	}

	if err = s.dataChangesPublisher.Publish(ctx, dcm); err != nil {
		observability.AcknowledgeError(err, logger, span, "publishing data change message")
	}

	// encode our response and peace.
	s.encoderDecoder.RespondWithData(ctx, res, validPreparationVessel)
}

// ArchiveHandler returns a handler that archives a valid preparation vessel.
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

	// determine valid preparation vessel ID.
	validPreparationVesselID := s.validPreparationVesselIDFetcher(req)
	tracing.AttachToSpan(span, keys.ValidVesselIDKey, validPreparationVesselID)
	logger = logger.WithValue(keys.ValidPreparationVesselIDKey, validPreparationVesselID)

	exists, existenceCheckErr := s.validPreparationVesselDataManager.ValidPreparationVesselExists(ctx, validPreparationVesselID)
	if existenceCheckErr != nil && !errors.Is(existenceCheckErr, sql.ErrNoRows) {
		observability.AcknowledgeError(existenceCheckErr, logger, span, "checking valid preparation vessel existence")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	} else if !exists || errors.Is(existenceCheckErr, sql.ErrNoRows) {
		s.encoderDecoder.EncodeNotFoundResponse(ctx, res)
		return
	}

	if err = s.validPreparationVesselDataManager.ArchiveValidPreparationVessel(ctx, validPreparationVesselID); err != nil {
		observability.AcknowledgeError(err, logger, span, "archiving valid preparation vessel")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	dcm := &types.DataChangeMessage{
		EventType: types.ValidPreparationVesselArchivedCustomerEventType,
		UserID:    sessionCtxData.Requester.UserID,
	}

	if err = s.dataChangesPublisher.Publish(ctx, dcm); err != nil {
		observability.AcknowledgeError(err, logger, span, "publishing data change message")
	}

	// encode our response and peace.
	res.WriteHeader(http.StatusNoContent)
}

// SearchByPreparationHandler is our valid preparation vessel search route for preparations.
func (s *service) SearchByPreparationHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	tracing.AttachRequestToSpan(span, req)

	filter := types.ExtractQueryFilterFromRequest(req)
	tracing.AttachFilterDataToSpan(span, filter.Page, filter.Limit, filter.SortBy)

	logger := s.logger.WithRequest(req).
		WithValue(keys.FilterLimitKey, filter.Limit).
		WithValue(keys.FilterPageKey, filter.Page).
		WithValue(keys.FilterSortByKey, filter.SortBy)

	validPreparationID := s.validPreparationIDFetcher(req)

	// determine user ID.
	sessionCtxData, err := s.sessionContextDataFetcher(req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving session context data")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, "unauthenticated", http.StatusUnauthorized)
		return
	}

	tracing.AttachSessionContextDataToSpan(span, sessionCtxData)
	logger = sessionCtxData.AttachToLogger(logger)

	validPreparationVessels, err := s.validPreparationVesselDataManager.GetValidPreparationVesselsForPreparation(ctx, validPreparationID, filter)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "searching for valid preparation vessels")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	s.encoderDecoder.EncodeResponseWithStatus(ctx, res, validPreparationVessels, http.StatusOK)
}

// SearchByVesselHandler is our valid preparation vessel search route by vessel.
func (s *service) SearchByVesselHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	tracing.AttachRequestToSpan(span, req)

	filter := types.ExtractQueryFilterFromRequest(req)
	tracing.AttachFilterDataToSpan(span, filter.Page, filter.Limit, filter.SortBy)

	logger := s.logger.WithRequest(req).
		WithValue(keys.FilterLimitKey, filter.Limit).
		WithValue(keys.FilterPageKey, filter.Page).
		WithValue(keys.FilterSortByKey, filter.SortBy)

	validVesselID := s.validVesselIDFetcher(req)
	logger = logger.WithValue(keys.ValidVesselIDKey, validVesselID)

	// determine user ID.
	sessionCtxData, err := s.sessionContextDataFetcher(req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving session context data")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, "unauthenticated", http.StatusUnauthorized)
		return
	}

	tracing.AttachSessionContextDataToSpan(span, sessionCtxData)
	logger = sessionCtxData.AttachToLogger(logger)

	validPreparationVessels, err := s.validPreparationVesselDataManager.GetValidPreparationVesselsForVessel(ctx, validVesselID, filter)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "searching for valid preparation vessels")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	s.encoderDecoder.EncodeResponseWithStatus(ctx, res, validPreparationVessels, http.StatusOK)
}
