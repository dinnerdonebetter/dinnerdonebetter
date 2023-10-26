package households

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/dinnerdonebetter/backend/internal/authorization"
	"github.com/dinnerdonebetter/backend/internal/identifiers"
	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"
)

const (
	// HouseholdIDURIParamKey is a standard string that we'll use to refer to household IDs with.
	HouseholdIDURIParamKey = "householdID"
	// UserIDURIParamKey is a standard string that we'll use to refer to user IDs with.
	UserIDURIParamKey = "userID"
)

var _ types.HouseholdDataService = (*service)(nil)

// ListHandler is our list route.
func (s *service) ListHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	filter := types.ExtractQueryFilterFromRequest(req)
	logger := s.logger.WithRequest(req)
	logger = filter.AttachToLogger(logger)

	responseDetails := types.ResponseDetails{
		TraceID: span.SpanContext().TraceID().String(),
	}

	tracing.AttachRequestToSpan(span, req)
	tracing.AttachFilterDataToSpan(span, filter.Page, filter.Limit, filter.SortBy)

	// fetch session context data
	sessionCtxData, err := s.sessionContextDataFetcher(req)
	if err != nil {
		errRes := types.NewAPIErrorResponse("unauthenticated", types.ErrFetchingSessionContextData, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusUnauthorized)
		return
	}

	requester := sessionCtxData.Requester.UserID
	logger = logger.WithValue(keys.RequesterIDKey, requester)
	tracing.AttachSessionContextDataToSpan(span, sessionCtxData)

	filter.AttachToLogger(logger)

	households, err := s.householdDataManager.GetHouseholds(ctx, requester, filter)
	if errors.Is(err, sql.ErrNoRows) {
		// in the event no rows exist, return an empty list.
		households = &types.QueryFilteredResult[types.Household]{Data: []*types.Household{}}
	} else if err != nil {
		observability.AcknowledgeError(err, logger, span, "fetching households")
		errRes := types.NewAPIErrorResponse("database error", types.ErrTalkingToDatabase, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusInternalServerError)
		return
	}

	responseValue := &types.APIResponse[[]*types.Household]{
		Details:    responseDetails,
		Data:       households.Data,
		Pagination: &households.Pagination,
	}

	// encode our response and say farewell.
	s.encoderDecoder.RespondWithData(ctx, res, responseValue)
}

// CreateHandler is our household creation route.
func (s *service) CreateHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	logger := s.logger.WithRequest(req)
	tracing.AttachRequestToSpan(span, req)

	responseDetails := types.ResponseDetails{
		TraceID: span.SpanContext().TraceID().String(),
	}

	// retrieve session context data.
	sessionCtxData, err := s.sessionContextDataFetcher(req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving session context data")
		errRes := types.NewAPIErrorResponse("unauthenticated", types.ErrFetchingSessionContextData, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusUnauthorized)
		return
	}

	tracing.AttachSessionContextDataToSpan(span, sessionCtxData)
	logger = sessionCtxData.AttachToLogger(logger)

	// read parsed input struct from request body.
	providedInput := new(types.HouseholdCreationRequestInput)
	if err = s.encoderDecoder.DecodeRequest(ctx, req, providedInput); err != nil {
		observability.AcknowledgeError(err, logger, span, "decoding request body")
		errRes := types.NewAPIErrorResponse("invalid request content", types.ErrDecodingRequestInput, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusBadRequest)
		return
	}

	if err = providedInput.ValidateWithContext(ctx); err != nil {
		logger.WithValue(keys.ValidationErrorKey, err).Debug("invalid input attached to request")
		errRes := types.NewAPIErrorResponse(err.Error(), types.ErrValidatingRequestInput, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusBadRequest)
		return
	}

	tracing.AttachSessionContextDataToSpan(span, sessionCtxData)
	requester := sessionCtxData.Requester.UserID
	logger = logger.WithValue(keys.RequesterIDKey, requester)

	input := converters.ConvertHouseholdCreationInputToHouseholdDatabaseCreationInput(providedInput)
	input.ID = identifiers.New()
	input.BelongsToUser = requester

	input.WebhookEncryptionKey, err = s.secretGenerator.GenerateHexEncodedString(ctx, 128)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "generating webhook encryption key")
		errRes := types.NewAPIErrorResponse("encryption key error", types.ErrSecretGeneration, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusInternalServerError)
		return
	}

	logger = logger.WithValue(keys.NameKey, input.Name)

	// create household in database.
	household, err := s.householdDataManager.CreateHousehold(ctx, input)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "creating household")
		errRes := types.NewAPIErrorResponse("database error", types.ErrTalkingToDatabase, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusInternalServerError)
		return
	}

	logger = logger.WithValue(keys.HouseholdIDKey, household.ID)
	tracing.AttachToSpan(span, keys.HouseholdIDKey, household.ID)

	// notify relevant parties.
	logger.Debug("created household")

	dcm := &types.DataChangeMessage{
		EventType:   types.HouseholdCreatedCustomerEventType,
		Household:   household,
		HouseholdID: household.ID,
		UserID:      sessionCtxData.Requester.UserID,
	}
	if err = s.dataChangesPublisher.Publish(ctx, dcm); err != nil {
		observability.AcknowledgeError(err, logger, span, "publishing data change message for created household")
	}

	responseValue := &types.APIResponse[*types.Household]{
		Details: responseDetails,
		Data:    household,
	}

	s.encoderDecoder.EncodeResponseWithStatus(ctx, res, responseValue, http.StatusCreated)
}

// CurrentInfoHandler returns a handler that returns the current household.
func (s *service) CurrentInfoHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	logger := s.logger.WithRequest(req)
	tracing.AttachRequestToSpan(span, req)

	responseDetails := types.ResponseDetails{
		TraceID: span.SpanContext().TraceID().String(),
	}

	logger.Debug("current household info requested")

	// determine user ID.
	sessionCtxData, err := s.sessionContextDataFetcher(req)
	if err != nil {
		logger.Info("session context data missing from request")
		errRes := types.NewAPIErrorResponse("unauthenticated", types.ErrFetchingSessionContextData, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusUnauthorized)
		return
	}

	requester := sessionCtxData.Requester.UserID
	logger = logger.WithValue(keys.RequesterIDKey, requester)
	tracing.AttachSessionContextDataToSpan(span, sessionCtxData)

	// determine household ID.
	householdID := sessionCtxData.ActiveHouseholdID
	logger = logger.WithValue(keys.HouseholdIDKey, householdID)
	tracing.AttachToSpan(span, keys.HouseholdIDKey, householdID)

	// fetch household from database.
	household, err := s.householdDataManager.GetHousehold(ctx, householdID)
	if errors.Is(err, sql.ErrNoRows) {
		errRes := types.NewAPIErrorResponse("not found", types.ErrDataNotFound, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusNotFound)
		return
	} else if err != nil {
		observability.AcknowledgeError(err, logger, span, "fetching household from database")
		errRes := types.NewAPIErrorResponse("database error", types.ErrTalkingToDatabase, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusInternalServerError)
		return
	}

	logger.Debug("responding with current household info")

	responseValue := &types.APIResponse[*types.Household]{
		Details: responseDetails,
		Data:    household,
	}

	// encode our response and peace.
	s.encoderDecoder.RespondWithData(ctx, res, responseValue)
}

// ReadHandler returns a GET handler that returns a household.
func (s *service) ReadHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	logger := s.logger.WithRequest(req)
	tracing.AttachRequestToSpan(span, req)

	responseDetails := types.ResponseDetails{
		TraceID: span.SpanContext().TraceID().String(),
	}

	// determine user ID.
	sessionCtxData, err := s.sessionContextDataFetcher(req)
	if err != nil {
		errRes := types.NewAPIErrorResponse("unauthenticated", types.ErrFetchingSessionContextData, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusUnauthorized)
		return
	}

	requester := sessionCtxData.Requester.UserID
	logger = logger.WithValue(keys.RequesterIDKey, requester)
	tracing.AttachSessionContextDataToSpan(span, sessionCtxData)

	// determine household ID.
	householdID := s.householdIDFetcher(req)
	logger = logger.WithValue(keys.HouseholdIDKey, householdID)
	tracing.AttachToSpan(span, keys.HouseholdIDKey, householdID)

	// fetch household from database.
	household, err := s.householdDataManager.GetHousehold(ctx, householdID)
	if errors.Is(err, sql.ErrNoRows) {
		errRes := types.NewAPIErrorResponse("not found", types.ErrDataNotFound, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusNotFound)
		return
	} else if err != nil {
		observability.AcknowledgeError(err, logger, span, "fetching household from database")
		errRes := types.NewAPIErrorResponse("database error", types.ErrTalkingToDatabase, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusInternalServerError)
		return
	}

	admins := []*types.HouseholdUserMembershipWithUser{}
	plainUsers := []*types.HouseholdUserMembershipWithUser{}
	for _, member := range household.Members {
		if member.HouseholdRole == authorization.HouseholdAdminRole.String() {
			admins = append(admins, member)
		} else {
			plainUsers = append(plainUsers, member)
		}
	}

	household.Members = []*types.HouseholdUserMembershipWithUser{}
	household.Members = append(household.Members, admins...)
	household.Members = append(household.Members, plainUsers...)

	responseValue := &types.APIResponse[*types.Household]{
		Details: responseDetails,
		Data:    household,
	}

	// encode our response and peace.
	s.encoderDecoder.RespondWithData(ctx, res, responseValue)
}

// UpdateHandler returns a handler that updates a household.
func (s *service) UpdateHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	logger := s.logger.WithRequest(req)
	tracing.AttachRequestToSpan(span, req)

	responseDetails := types.ResponseDetails{
		TraceID: span.SpanContext().TraceID().String(),
	}

	// determine user ID.
	sessionCtxData, err := s.sessionContextDataFetcher(req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "fetching session context data")
		errRes := types.NewAPIErrorResponse("unauthenticated", types.ErrFetchingSessionContextData, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusUnauthorized)
		return
	}

	tracing.AttachSessionContextDataToSpan(span, sessionCtxData)
	logger = sessionCtxData.AttachToLogger(logger)

	input := new(types.HouseholdUpdateRequestInput)
	if err = s.encoderDecoder.DecodeRequest(ctx, req, input); err != nil {
		observability.AcknowledgeError(err, logger, span, "decoding request body")
		errRes := types.NewAPIErrorResponse("invalid request content", types.ErrDecodingRequestInput, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusBadRequest)
		return
	}

	if err = input.ValidateWithContext(ctx); err != nil {
		logger.WithValue(keys.ValidationErrorKey, err).Debug("invalid input attached to request")
		errRes := types.NewAPIErrorResponse(err.Error(), types.ErrValidatingRequestInput, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusBadRequest)
		return
	}

	tracing.AttachSessionContextDataToSpan(span, sessionCtxData)
	requester := sessionCtxData.Requester.UserID
	logger = logger.WithValue(keys.RequesterIDKey, requester)
	input.BelongsToUser = requester

	// determine household ID.
	householdID := s.householdIDFetcher(req)
	logger = logger.WithValue(keys.HouseholdIDKey, householdID)
	tracing.AttachToSpan(span, keys.HouseholdIDKey, householdID)

	// fetch household from database.
	household, err := s.householdDataManager.GetHousehold(ctx, householdID)
	if errors.Is(err, sql.ErrNoRows) {
		errRes := types.NewAPIErrorResponse("not found", types.ErrDataNotFound, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusNotFound)
		return
	} else if err != nil {
		observability.AcknowledgeError(err, logger, span, "fetching household from database")
		errRes := types.NewAPIErrorResponse("database error", types.ErrTalkingToDatabase, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusInternalServerError)
		return
	}

	// update the data structure.
	household.Update(input)

	// update household in database.
	if err = s.householdDataManager.UpdateHousehold(ctx, household); err != nil {
		observability.AcknowledgeError(err, logger, span, "updating household")
		errRes := types.NewAPIErrorResponse("database error", types.ErrTalkingToDatabase, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusInternalServerError)
		return
	}

	dcm := &types.DataChangeMessage{
		EventType:   types.HouseholdUpdatedCustomerEventType,
		Household:   household,
		HouseholdID: household.ID,
		UserID:      sessionCtxData.Requester.UserID,
	}
	if err = s.dataChangesPublisher.Publish(ctx, dcm); err != nil {
		observability.AcknowledgeError(err, logger, span, "publishing data change message for created household")
	}

	responseValue := &types.APIResponse[*types.Household]{
		Details: responseDetails,
		Data:    household,
	}

	// encode our response and peace.
	s.encoderDecoder.RespondWithData(ctx, res, responseValue)
}

// ArchiveHandler returns a handler that archives a household.
func (s *service) ArchiveHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	logger := s.logger.WithRequest(req)
	tracing.AttachRequestToSpan(span, req)

	responseDetails := types.ResponseDetails{
		TraceID: span.SpanContext().TraceID().String(),
	}

	// determine user ID.
	sessionCtxData, err := s.sessionContextDataFetcher(req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving session context data")
		errRes := types.NewAPIErrorResponse("unauthenticated", types.ErrFetchingSessionContextData, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusUnauthorized)
		return
	}

	requester := sessionCtxData.Requester.UserID
	logger = logger.WithValue(keys.RequesterIDKey, requester)
	tracing.AttachSessionContextDataToSpan(span, sessionCtxData)

	// determine household ID.
	householdID := s.householdIDFetcher(req)
	tracing.AttachToSpan(span, keys.HouseholdIDKey, householdID)
	logger = logger.WithValue(keys.HouseholdIDKey, householdID)

	// archive the household in the database.
	err = s.householdDataManager.ArchiveHousehold(ctx, householdID, requester)
	if errors.Is(err, sql.ErrNoRows) {
		errRes := types.NewAPIErrorResponse("not found", types.ErrDataNotFound, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusNotFound)
		return
	} else if err != nil {
		observability.AcknowledgeError(err, logger, span, "archiving household")
		errRes := types.NewAPIErrorResponse("database error", types.ErrTalkingToDatabase, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusInternalServerError)
		return
	}

	dcm := &types.DataChangeMessage{
		EventType:   types.HouseholdArchivedCustomerEventType,
		HouseholdID: householdID,
		UserID:      sessionCtxData.Requester.UserID,
	}
	if err = s.dataChangesPublisher.Publish(ctx, dcm); err != nil {
		observability.AcknowledgeError(err, logger, span, "publishing data change message for created household")
	}

	responseValue := &types.APIResponse[*types.Household]{
		Details: responseDetails,
	}

	// let everybody go home.
	s.encoderDecoder.RespondWithData(ctx, res, responseValue)
}

// ModifyMemberPermissionsHandler is our household creation route.
func (s *service) ModifyMemberPermissionsHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	logger := s.logger.WithRequest(req)
	tracing.AttachRequestToSpan(span, req)

	responseDetails := types.ResponseDetails{
		TraceID: span.SpanContext().TraceID().String(),
	}

	// determine user ID.
	sessionCtxData, err := s.sessionContextDataFetcher(req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving session context data")
		errRes := types.NewAPIErrorResponse("unauthenticated", types.ErrFetchingSessionContextData, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusUnauthorized)
		return
	}

	tracing.AttachSessionContextDataToSpan(span, sessionCtxData)
	logger = sessionCtxData.AttachToLogger(logger)

	// read parsed input struct from request body.
	input := new(types.ModifyUserPermissionsInput)
	if err = s.encoderDecoder.DecodeRequest(ctx, req, input); err != nil {
		observability.AcknowledgeError(err, logger, span, "decoding request body")
		errRes := types.NewAPIErrorResponse("invalid request content", types.ErrDecodingRequestInput, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusBadRequest)
		return
	}

	if err = input.ValidateWithContext(ctx); err != nil {
		logger.WithValue(keys.ValidationErrorKey, err).Debug("invalid input attached to request")
		errRes := types.NewAPIErrorResponse(err.Error(), types.ErrValidatingRequestInput, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusBadRequest)
		return
	}

	requester := sessionCtxData.Requester.UserID
	logger = logger.WithValue(keys.RequesterIDKey, requester)
	tracing.AttachSessionContextDataToSpan(span, sessionCtxData)

	householdID := s.householdIDFetcher(req)
	logger = logger.WithValue(keys.HouseholdIDKey, householdID)
	tracing.AttachToSpan(span, keys.HouseholdIDKey, householdID)

	userID := s.userIDFetcher(req)
	logger = logger.WithValue(keys.UserIDKey, userID)
	tracing.AttachToSpan(span, keys.UserIDKey, userID)

	// create household in database.
	if err = s.householdMembershipDataManager.ModifyUserPermissions(ctx, householdID, userID, input); err != nil {
		observability.AcknowledgeError(err, logger, span, "modifying user permissions")
		errRes := types.NewAPIErrorResponse("database error", types.ErrTalkingToDatabase, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusInternalServerError)
		return
	}

	dcm := &types.DataChangeMessage{
		EventType:   types.HouseholdMembershipPermissionsUpdatedCustomerEventType,
		HouseholdID: householdID,
		UserID:      sessionCtxData.Requester.UserID,
	}
	if err = s.dataChangesPublisher.Publish(ctx, dcm); err != nil {
		observability.AcknowledgeError(err, logger, span, "publishing data change message for created household")
	}

	responseValue := &types.APIResponse[*types.Webhook]{
		Details: responseDetails,
	}

	s.encoderDecoder.EncodeResponseWithStatus(ctx, res, responseValue, http.StatusAccepted)
}

// TransferHouseholdOwnershipHandler is our household creation route.
func (s *service) TransferHouseholdOwnershipHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	logger := s.logger.WithRequest(req)
	tracing.AttachRequestToSpan(span, req)

	responseDetails := types.ResponseDetails{
		TraceID: span.SpanContext().TraceID().String(),
	}

	// determine user ID.
	sessionCtxData, err := s.sessionContextDataFetcher(req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "transferring household ownership")
		errRes := types.NewAPIErrorResponse("unauthenticated", types.ErrFetchingSessionContextData, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusUnauthorized)
		return
	}

	tracing.AttachSessionContextDataToSpan(span, sessionCtxData)
	logger = sessionCtxData.AttachToLogger(logger)

	// read parsed input struct from request body.
	input := new(types.HouseholdOwnershipTransferInput)
	if err = s.encoderDecoder.DecodeRequest(ctx, req, input); err != nil {
		observability.AcknowledgeError(err, logger, span, "decoding request body")
		errRes := types.NewAPIErrorResponse("invalid request content", types.ErrDecodingRequestInput, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusBadRequest)
		return
	}

	if err = input.ValidateWithContext(ctx); err != nil {
		logger.WithValue(keys.ValidationErrorKey, err).Debug("invalid input attached to request")
		errRes := types.NewAPIErrorResponse(err.Error(), types.ErrValidatingRequestInput, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusBadRequest)
		return
	}

	householdID := s.householdIDFetcher(req)
	tracing.AttachToSpan(span, keys.HouseholdIDKey, householdID)
	logger = logger.WithValue(keys.HouseholdIDKey, householdID)

	requester := sessionCtxData.Requester.UserID
	tracing.AttachSessionContextDataToSpan(span, sessionCtxData)
	logger = logger.WithValue(keys.RequesterIDKey, requester)

	// transfer ownership of household in database.
	if err = s.householdMembershipDataManager.TransferHouseholdOwnership(ctx, householdID, input); err != nil {
		observability.AcknowledgeError(err, logger, span, "transferring household ownership")
		errRes := types.NewAPIErrorResponse("database error", types.ErrTalkingToDatabase, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusInternalServerError)
		return
	}

	dcm := &types.DataChangeMessage{
		EventType:   types.HouseholdOwnershipTransferredCustomerEventType,
		HouseholdID: householdID,
		UserID:      sessionCtxData.Requester.UserID,
	}
	if err = s.dataChangesPublisher.Publish(ctx, dcm); err != nil {
		observability.AcknowledgeError(err, logger, span, "publishing data change message for created household")
	}

	responseValue := &types.APIResponse[*types.Webhook]{
		Details: responseDetails,
	}

	s.encoderDecoder.EncodeResponseWithStatus(ctx, res, responseValue, http.StatusAccepted)
}

// RemoveMemberHandler is our household creation route.
func (s *service) RemoveMemberHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	logger := s.logger.WithRequest(req)
	tracing.AttachRequestToSpan(span, req)

	responseDetails := types.ResponseDetails{
		TraceID: span.SpanContext().TraceID().String(),
	}

	// read parsed input struct from request body.
	reason := req.URL.Query().Get("reason")
	logger = logger.WithValue(keys.ReasonKey, reason)

	// determine user ID.
	sessionCtxData, err := s.sessionContextDataFetcher(req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving session context data")
		errRes := types.NewAPIErrorResponse("unauthenticated", types.ErrFetchingSessionContextData, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusUnauthorized)
		return
	}

	requester := sessionCtxData.Requester.UserID
	logger = logger.WithValue(keys.RequesterIDKey, requester)
	tracing.AttachSessionContextDataToSpan(span, sessionCtxData)

	householdID := s.householdIDFetcher(req)
	logger = logger.WithValue(keys.HouseholdIDKey, householdID)
	tracing.AttachToSpan(span, keys.HouseholdIDKey, householdID)

	userID := s.userIDFetcher(req)
	logger = logger.WithValue(keys.UserIDKey, userID)
	tracing.AttachToSpan(span, keys.UserIDKey, userID)

	// remove user from household in database.
	if err = s.householdMembershipDataManager.RemoveUserFromHousehold(ctx, userID, householdID); err != nil {
		observability.AcknowledgeError(err, logger, span, "removing user from household")
		errRes := types.NewAPIErrorResponse("database error", types.ErrTalkingToDatabase, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusInternalServerError)
		return
	}

	logger.Info("user removed from household")

	dcm := &types.DataChangeMessage{
		EventType:   types.HouseholdMemberRemovedCustomerEventType,
		HouseholdID: householdID,
		UserID:      sessionCtxData.Requester.UserID,
	}
	if err = s.dataChangesPublisher.Publish(ctx, dcm); err != nil {
		observability.AcknowledgeError(err, logger, span, "publishing data change message for created household")
	}

	responseValue := &types.APIResponse[*types.Webhook]{
		Details: responseDetails,
	}

	s.encoderDecoder.EncodeResponseWithStatus(ctx, res, responseValue, http.StatusAccepted)
}

// MarkAsDefaultHouseholdHandler is our household creation route.
func (s *service) MarkAsDefaultHouseholdHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	logger := s.logger.WithRequest(req)
	tracing.AttachRequestToSpan(span, req)

	responseDetails := types.ResponseDetails{
		TraceID: span.SpanContext().TraceID().String(),
	}

	householdID := s.householdIDFetcher(req)
	logger = logger.WithValue(keys.HouseholdIDKey, householdID)
	tracing.AttachToSpan(span, keys.HouseholdIDKey, householdID)

	// determine user ID.
	sessionCtxData, err := s.sessionContextDataFetcher(req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving session context data")
		errRes := types.NewAPIErrorResponse("unauthenticated", types.ErrFetchingSessionContextData, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusUnauthorized)
		return
	}

	requester := sessionCtxData.Requester.UserID
	logger = logger.WithValue(keys.RequesterIDKey, requester)
	tracing.AttachSessionContextDataToSpan(span, sessionCtxData)

	// mark household as default in database.
	if err = s.householdMembershipDataManager.MarkHouseholdAsUserDefault(ctx, requester, householdID); err != nil {
		observability.AcknowledgeError(err, logger, span, "marking household as default")
		errRes := types.NewAPIErrorResponse("database error", types.ErrTalkingToDatabase, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusInternalServerError)
		return
	}

	dcm := &types.DataChangeMessage{
		EventType:   types.HouseholdMemberRemovedCustomerEventType,
		HouseholdID: householdID,
		UserID:      sessionCtxData.Requester.UserID,
	}
	if err = s.dataChangesPublisher.Publish(ctx, dcm); err != nil {
		observability.AcknowledgeError(err, logger, span, "publishing data change message for created household")
	}

	responseValue := &types.APIResponse[*types.Webhook]{
		Details: responseDetails,
	}

	s.encoderDecoder.EncodeResponseWithStatus(ctx, res, responseValue, http.StatusAccepted)
}
