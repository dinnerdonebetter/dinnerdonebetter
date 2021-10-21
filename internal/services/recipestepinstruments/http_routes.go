package recipestepinstruments

import (
	"database/sql"
	"errors"
	"net/http"
	"strings"

	"github.com/segmentio/ksuid"

	observability "gitlab.com/prixfixe/prixfixe/internal/observability"
	keys "gitlab.com/prixfixe/prixfixe/internal/observability/keys"
	"gitlab.com/prixfixe/prixfixe/internal/observability/tracing"
	"gitlab.com/prixfixe/prixfixe/pkg/types"
)

const (
	// RecipeStepInstrumentIDURIParamKey is a standard string that we'll use to refer to recipe step instrument IDs with.
	RecipeStepInstrumentIDURIParamKey = "recipeStepInstrumentID"
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

// CreateHandler is our recipe step instrument creation route.
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
	providedInput := new(types.RecipeStepInstrumentCreationRequestInput)
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

	input := types.RecipeStepInstrumentDatabaseCreationInputFromRecipeStepInstrumentCreationInput(providedInput)
	input.ID = ksuid.New().String()

	// determine recipe step ID.
	recipeStepID := s.recipeStepIDFetcher(req)
	tracing.AttachRecipeStepIDToSpan(span, recipeStepID)
	logger = logger.WithValue(keys.RecipeStepIDKey, recipeStepID)

	input.BelongsToRecipeStep = recipeStepID
	tracing.AttachRecipeStepInstrumentIDToSpan(span, input.ID)

	// create recipe step instrument in database.
	preWrite := &types.PreWriteMessage{
		DataType:                types.RecipeStepInstrumentDataType,
		RecipeStepID:            recipeStepID,
		RecipeStepInstrument:    input,
		AttributableToUserID:    sessionCtxData.Requester.UserID,
		AttributableToAccountID: sessionCtxData.ActiveAccountID,
	}
	if err = s.preWritesPublisher.Publish(ctx, preWrite); err != nil {
		observability.AcknowledgeError(err, logger, span, "publishing recipe step instrument write message")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	pwr := types.PreWriteResponse{ID: input.ID}

	s.encoderDecoder.EncodeResponseWithStatus(ctx, res, pwr, http.StatusAccepted)
}

// ReadHandler returns a GET handler that returns a recipe step instrument.
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

	// determine recipe ID.
	recipeID := s.recipeIDFetcher(req)
	tracing.AttachRecipeIDToSpan(span, recipeID)
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)

	// determine recipe step ID.
	recipeStepID := s.recipeStepIDFetcher(req)
	tracing.AttachRecipeStepIDToSpan(span, recipeStepID)
	logger = logger.WithValue(keys.RecipeStepIDKey, recipeStepID)

	// determine recipe step instrument ID.
	recipeStepInstrumentID := s.recipeStepInstrumentIDFetcher(req)
	tracing.AttachRecipeStepInstrumentIDToSpan(span, recipeStepInstrumentID)
	logger = logger.WithValue(keys.RecipeStepInstrumentIDKey, recipeStepInstrumentID)

	// fetch recipe step instrument from database.
	x, err := s.recipeStepInstrumentDataManager.GetRecipeStepInstrument(ctx, recipeID, recipeStepID, recipeStepInstrumentID)
	if errors.Is(err, sql.ErrNoRows) {
		s.encoderDecoder.EncodeNotFoundResponse(ctx, res)
		return
	} else if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving recipe step instrument")
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

	// determine recipe ID.
	recipeID := s.recipeIDFetcher(req)
	tracing.AttachRecipeIDToSpan(span, recipeID)
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)

	// determine recipe step ID.
	recipeStepID := s.recipeStepIDFetcher(req)
	tracing.AttachRecipeStepIDToSpan(span, recipeStepID)
	logger = logger.WithValue(keys.RecipeStepIDKey, recipeStepID)

	recipeStepInstruments, err := s.recipeStepInstrumentDataManager.GetRecipeStepInstruments(ctx, recipeID, recipeStepID, filter)
	if errors.Is(err, sql.ErrNoRows) {
		// in the event no rows exist, return an empty list.
		recipeStepInstruments = &types.RecipeStepInstrumentList{RecipeStepInstruments: []*types.RecipeStepInstrument{}}
	} else if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving recipe step instruments")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	// encode our response and peace.
	s.encoderDecoder.RespondWithData(ctx, res, recipeStepInstruments)
}

// UpdateHandler returns a handler that updates a recipe step instrument.
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
	input := new(types.RecipeStepInstrumentUpdateRequestInput)
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

	// determine recipe ID.
	recipeID := s.recipeIDFetcher(req)
	tracing.AttachRecipeIDToSpan(span, recipeID)
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)

	// determine recipe step ID.
	recipeStepID := s.recipeStepIDFetcher(req)
	tracing.AttachRecipeStepIDToSpan(span, recipeStepID)
	logger = logger.WithValue(keys.RecipeStepIDKey, recipeStepID)

	// determine recipe step instrument ID.
	recipeStepInstrumentID := s.recipeStepInstrumentIDFetcher(req)
	tracing.AttachRecipeStepInstrumentIDToSpan(span, recipeStepInstrumentID)
	logger = logger.WithValue(keys.RecipeStepInstrumentIDKey, recipeStepInstrumentID)

	// fetch recipe step instrument from database.
	recipeStepInstrument, err := s.recipeStepInstrumentDataManager.GetRecipeStepInstrument(ctx, recipeID, recipeStepID, recipeStepInstrumentID)
	if errors.Is(err, sql.ErrNoRows) {
		s.encoderDecoder.EncodeNotFoundResponse(ctx, res)
		return
	} else if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving recipe step instrument for update")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	// update the recipe step instrument.
	recipeStepInstrument.Update(input)

	pum := &types.PreUpdateMessage{
		DataType:                types.RecipeStepInstrumentDataType,
		RecipeStepID:            recipeStepID,
		RecipeStepInstrument:    recipeStepInstrument,
		AttributableToUserID:    sessionCtxData.Requester.UserID,
		AttributableToAccountID: sessionCtxData.ActiveAccountID,
	}
	if err = s.preUpdatesPublisher.Publish(ctx, pum); err != nil {
		observability.AcknowledgeError(err, logger, span, "publishing recipe step instrument update message")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	// encode our response and peace.
	s.encoderDecoder.RespondWithData(ctx, res, recipeStepInstrument)
}

// ArchiveHandler returns a handler that archives a recipe step instrument.
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

	// determine recipe ID.
	recipeID := s.recipeIDFetcher(req)
	tracing.AttachRecipeIDToSpan(span, recipeID)
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)

	// determine recipe step ID.
	recipeStepID := s.recipeStepIDFetcher(req)
	tracing.AttachRecipeStepIDToSpan(span, recipeStepID)
	logger = logger.WithValue(keys.RecipeStepIDKey, recipeStepID)

	// determine recipe step instrument ID.
	recipeStepInstrumentID := s.recipeStepInstrumentIDFetcher(req)
	tracing.AttachRecipeStepInstrumentIDToSpan(span, recipeStepInstrumentID)
	logger = logger.WithValue(keys.RecipeStepInstrumentIDKey, recipeStepInstrumentID)

	exists, existenceCheckErr := s.recipeStepInstrumentDataManager.RecipeStepInstrumentExists(ctx, recipeID, recipeStepID, recipeStepInstrumentID)
	if existenceCheckErr != nil && !errors.Is(existenceCheckErr, sql.ErrNoRows) {
		observability.AcknowledgeError(existenceCheckErr, logger, span, "checking recipe step instrument existence")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	} else if !exists || errors.Is(existenceCheckErr, sql.ErrNoRows) {
		s.encoderDecoder.EncodeNotFoundResponse(ctx, res)
		return
	}

	pam := &types.PreArchiveMessage{
		DataType:                types.RecipeStepInstrumentDataType,
		RecipeStepID:            recipeStepID,
		RecipeStepInstrumentID:  recipeStepInstrumentID,
		AttributableToUserID:    sessionCtxData.Requester.UserID,
		AttributableToAccountID: sessionCtxData.ActiveAccountID,
	}
	if err = s.preArchivesPublisher.Publish(ctx, pam); err != nil {
		observability.AcknowledgeError(err, logger, span, "publishing recipe step instrument archive message")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	// encode our response and peace.
	res.WriteHeader(http.StatusNoContent)
}
