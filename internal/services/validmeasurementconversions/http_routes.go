package validmeasurementconversions

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
	// ValidMeasurementConversionIDURIParamKey is a standard string that we'll use to refer to valid measurement conversion IDs with.
	ValidMeasurementConversionIDURIParamKey = "validMeasurementConversionID"
	// ValidMeasurementUnitIDURIParamKey is a standard string that we'll use to refer to valid measurement unit IDs with.
	ValidMeasurementUnitIDURIParamKey = "validMeasurementUnitID"
)

// CreateHandler is our valid measurement conversion creation route.
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
	providedInput := new(types.ValidMeasurementUnitConversionCreationRequestInput)
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

	input := converters.ConvertValidMeasurementConversionCreationRequestInputToValidMeasurementConversionDatabaseCreationInput(providedInput)
	input.ID = identifiers.New()

	tracing.AttachValidMeasurementConversionIDToSpan(span, input.ID)

	validMeasurementConversion, err := s.validMeasurementConversionDataManager.CreateValidMeasurementConversion(ctx, input)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "creating valid measurement conversions")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	dcm := &types.DataChangeMessage{
		DataType:                   types.ValidMeasurementConversionDataType,
		EventType:                  types.ValidMeasurementConversionCreatedCustomerEventType,
		ValidMeasurementConversion: validMeasurementConversion,
		UserID:                     sessionCtxData.Requester.UserID,
	}

	if err = s.dataChangesPublisher.Publish(ctx, dcm); err != nil {
		observability.AcknowledgeError(err, logger, span, "publishing to data changes topic")
	}

	s.encoderDecoder.EncodeResponseWithStatus(ctx, res, validMeasurementConversion, http.StatusCreated)
}

// ReadHandler returns a GET handler that returns a valid measurement conversion.
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

	// determine valid measurement conversion ID.
	validMeasurementConversionID := s.validMeasurementConversionIDFetcher(req)
	tracing.AttachValidMeasurementConversionIDToSpan(span, validMeasurementConversionID)
	logger = logger.WithValue(keys.ValidMeasurementConversionIDKey, validMeasurementConversionID)

	// fetch valid measurement conversion from database.
	x, err := s.validMeasurementConversionDataManager.GetValidMeasurementConversion(ctx, validMeasurementConversionID)
	if errors.Is(err, sql.ErrNoRows) {
		s.encoderDecoder.EncodeNotFoundResponse(ctx, res)
		return
	} else if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving valid measurement conversion")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	// encode our response and peace.
	s.encoderDecoder.RespondWithData(ctx, res, x)
}

// UpdateHandler returns a handler that updates a valid measurement conversion.
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
	input := new(types.ValidMeasurementUnitConversionUpdateRequestInput)
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

	// determine valid measurement conversion ID.
	validMeasurementConversionID := s.validMeasurementConversionIDFetcher(req)
	tracing.AttachValidMeasurementConversionIDToSpan(span, validMeasurementConversionID)
	logger = logger.WithValue(keys.ValidMeasurementConversionIDKey, validMeasurementConversionID)

	// fetch valid measurement conversion from database.
	validMeasurementConversion, err := s.validMeasurementConversionDataManager.GetValidMeasurementConversion(ctx, validMeasurementConversionID)
	if errors.Is(err, sql.ErrNoRows) {
		s.encoderDecoder.EncodeNotFoundResponse(ctx, res)
		return
	} else if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving valid measurement conversion for update")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	// update the valid measurement conversion.
	validMeasurementConversion.Update(input)

	if err = s.validMeasurementConversionDataManager.UpdateValidMeasurementConversion(ctx, validMeasurementConversion); err != nil {
		observability.AcknowledgeError(err, logger, span, "updating valid measurement conversions")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	dcm := &types.DataChangeMessage{
		DataType:                   types.ValidMeasurementConversionDataType,
		EventType:                  types.ValidMeasurementConversionUpdatedCustomerEventType,
		ValidMeasurementConversion: validMeasurementConversion,
		UserID:                     sessionCtxData.Requester.UserID,
	}

	if err = s.dataChangesPublisher.Publish(ctx, dcm); err != nil {
		observability.AcknowledgeError(err, logger, span, "publishing data change message")
	}

	// encode our response and peace.
	s.encoderDecoder.RespondWithData(ctx, res, validMeasurementConversion)
}

// ArchiveHandler returns a handler that archives a valid measurement conversion.
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

	// determine valid measurement conversion ID.
	validMeasurementConversionID := s.validMeasurementConversionIDFetcher(req)
	tracing.AttachValidMeasurementConversionIDToSpan(span, validMeasurementConversionID)
	logger = logger.WithValue(keys.ValidMeasurementConversionIDKey, validMeasurementConversionID)

	exists, existenceCheckErr := s.validMeasurementConversionDataManager.ValidMeasurementConversionExists(ctx, validMeasurementConversionID)
	if existenceCheckErr != nil && !errors.Is(existenceCheckErr, sql.ErrNoRows) {
		observability.AcknowledgeError(existenceCheckErr, logger, span, "checking valid measurement conversion existence")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	} else if !exists || errors.Is(existenceCheckErr, sql.ErrNoRows) {
		s.encoderDecoder.EncodeNotFoundResponse(ctx, res)
		return
	}

	if err = s.validMeasurementConversionDataManager.ArchiveValidMeasurementConversion(ctx, validMeasurementConversionID); err != nil {
		observability.AcknowledgeError(err, logger, span, "archiving valid measurement conversions")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	dcm := &types.DataChangeMessage{
		DataType:  types.ValidMeasurementConversionDataType,
		EventType: types.ValidMeasurementConversionArchivedCustomerEventType,
		UserID:    sessionCtxData.Requester.UserID,
	}

	if err = s.dataChangesPublisher.Publish(ctx, dcm); err != nil {
		observability.AcknowledgeError(err, logger, span, "publishing data change message")
	}

	// encode our response and peace.
	res.WriteHeader(http.StatusNoContent)
}

func (s *service) FromMeasurementUnitHandler(res http.ResponseWriter, req *http.Request) {
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

	// determine valid measurement conversion ID.
	validMeasurementUnitID := s.validMeasurementUnitIDFetcher(req)
	tracing.AttachValidMeasurementUnitIDToSpan(span, validMeasurementUnitID)
	logger = logger.WithValue(keys.ValidMeasurementUnitIDKey, validMeasurementUnitID)

	// fetch valid measurement conversion from database.
	x, err := s.validMeasurementConversionDataManager.GetValidMeasurementConversionsFromUnit(ctx, validMeasurementUnitID)
	if errors.Is(err, sql.ErrNoRows) {
		s.encoderDecoder.EncodeNotFoundResponse(ctx, res)
		return
	} else if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving valid measurement conversion")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	// encode our response and peace.
	s.encoderDecoder.RespondWithData(ctx, res, x)
}

func (s *service) ToMeasurementUnitHandler(res http.ResponseWriter, req *http.Request) {
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

	// determine valid measurement conversion ID.
	validMeasurementUnitID := s.validMeasurementUnitIDFetcher(req)
	tracing.AttachValidMeasurementUnitIDToSpan(span, validMeasurementUnitID)
	logger = logger.WithValue(keys.ValidMeasurementUnitIDKey, validMeasurementUnitID)

	// fetch valid measurement conversion from database.
	x, err := s.validMeasurementConversionDataManager.GetValidMeasurementConversionsToUnit(ctx, validMeasurementUnitID)
	if errors.Is(err, sql.ErrNoRows) {
		s.encoderDecoder.EncodeNotFoundResponse(ctx, res)
		return
	} else if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving valid measurement conversion")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	// encode our response and peace.
	s.encoderDecoder.RespondWithData(ctx, res, x)
}
