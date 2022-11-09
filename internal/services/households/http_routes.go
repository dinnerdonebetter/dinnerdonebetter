package households

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/prixfixeco/api_server/internal/authorization"
	"github.com/prixfixeco/api_server/internal/identifiers"
	"github.com/prixfixeco/api_server/internal/observability"
	"github.com/prixfixeco/api_server/internal/observability/keys"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/pkg/types"
	"github.com/prixfixeco/api_server/pkg/types/converters"
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
	logger := s.logger.WithRequest(req).
		WithValue(keys.FilterLimitKey, filter.Limit).
		WithValue(keys.FilterPageKey, filter.Page).
		WithValue(keys.FilterSortByKey, filter.SortBy)

	tracing.AttachRequestToSpan(span, req)
	tracing.AttachFilterDataToSpan(span, filter.Page, filter.Limit, filter.SortBy)

	// fetch session context data
	sessionCtxData, err := s.sessionContextDataFetcher(req)
	if err != nil {
		s.encoderDecoder.EncodeErrorResponse(ctx, res, "unauthenticated", http.StatusUnauthorized)
		return
	}

	requester := sessionCtxData.Requester.UserID
	logger = logger.WithValue(keys.RequesterIDKey, requester)
	tracing.AttachSessionContextDataToSpan(span, sessionCtxData)

	filter.AttachToLogger(logger)
	logger.Debug("filtering?!")

	fetchFunc := s.householdDataManager.GetHouseholds
	if sessionCtxData.ServiceRolePermissionChecker().IsServiceAdmin() {
		fetchFunc = s.householdDataManager.GetHouseholdsForAdmin
	}

	households, err := fetchFunc(ctx, requester, filter)
	if errors.Is(err, sql.ErrNoRows) {
		// in the event no rows exist, return an empty list.
		households = &types.HouseholdList{Households: []*types.Household{}}
	} else if err != nil {
		observability.AcknowledgeError(err, logger, span, "fetching households")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	// encode our response and say farewell.
	s.encoderDecoder.RespondWithData(ctx, res, households)
}

// CreateHandler is our household creation route.
func (s *service) CreateHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	logger := s.logger.WithRequest(req)
	tracing.AttachRequestToSpan(span, req)

	// retrieve session context data.
	sessionCtxData, err := s.sessionContextDataFetcher(req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving session context data")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, "unauthenticated", http.StatusUnauthorized)
		return
	}

	tracing.AttachSessionContextDataToSpan(span, sessionCtxData)
	logger = sessionCtxData.AttachToLogger(logger)

	// read parsed input struct from request body.
	providedInput := new(types.HouseholdCreationRequestInput)
	if err = s.encoderDecoder.DecodeRequest(ctx, req, providedInput); err != nil {
		observability.AcknowledgeError(err, logger, span, "decoding request body")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, "invalid request content", http.StatusBadRequest)
		return
	}

	if providedInput.TimeZone == "" {
		providedInput.TimeZone = types.DefaultHouseholdTimeZone
	}

	if err = providedInput.ValidateWithContext(ctx); err != nil {
		logger.WithValue(keys.ValidationErrorKey, err).Debug("invalid input attached to request")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, err.Error(), http.StatusBadRequest)
		return
	}

	tracing.AttachSessionContextDataToSpan(span, sessionCtxData)
	requester := sessionCtxData.Requester.UserID
	logger = logger.WithValue(keys.RequesterIDKey, requester)

	input := converters.ConvertHouseholdCreationInputToHouseholdDatabaseCreationInput(providedInput)
	input.ID = identifiers.New()
	input.BelongsToUser = requester

	logger = logger.WithValue(keys.NameKey, input.Name)

	// create household in database.
	household, err := s.householdDataManager.CreateHousehold(ctx, input)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "creating household")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	logger = logger.WithValue(keys.HouseholdIDKey, household.ID)
	tracing.AttachHouseholdIDToSpan(span, household.ID)

	// notify relevant parties.
	logger.Debug("created household")
	s.householdCounter.Increment(ctx)

	if s.dataChangesPublisher != nil {
		dcm := &types.DataChangeMessage{
			DataType:                  types.HouseholdDataType,
			EventType:                 types.HouseholdCreatedCustomerEventType,
			Household:                 household,
			AttributableToUserID:      requester,
			AttributableToHouseholdID: household.ID,
		}
		if err = s.dataChangesPublisher.Publish(ctx, dcm); err != nil {
			observability.AcknowledgeError(err, logger, span, "publishing data change message for created household")
		}
	}

	s.encoderDecoder.EncodeResponseWithStatus(ctx, res, household, http.StatusCreated)
}

// CurrentInfoHandler returns a handler that returns the current household.
func (s *service) CurrentInfoHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	logger := s.logger.WithRequest(req)
	tracing.AttachRequestToSpan(span, req)

	logger.Debug("current household info requested")

	// determine user ID.
	sessionCtxData, err := s.sessionContextDataFetcher(req)
	if err != nil {
		logger.Info("session context data missing from request")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, "unauthenticated", http.StatusUnauthorized)
		return
	}

	requester := sessionCtxData.Requester.UserID
	logger = logger.WithValue(keys.RequesterIDKey, requester)
	tracing.AttachSessionContextDataToSpan(span, sessionCtxData)

	// determine household ID.
	householdID := sessionCtxData.ActiveHouseholdID
	logger = logger.WithValue(keys.HouseholdIDKey, householdID)
	tracing.AttachHouseholdIDToSpan(span, householdID)

	// fetch household from database.
	household, err := s.householdDataManager.GetHousehold(ctx, householdID, requester)
	if errors.Is(err, sql.ErrNoRows) {
		logger.Info("household ID is invalid")
		s.encoderDecoder.EncodeNotFoundResponse(ctx, res)
		return
	} else if err != nil {
		logger.Info("something is fucked!")
		observability.AcknowledgeError(err, logger, span, "fetching household from database")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	logger.Debug("responding with current household info")

	// encode our response and peace.
	s.encoderDecoder.RespondWithData(ctx, res, household)
}

// ReadHandler returns a GET handler that returns a household.
func (s *service) ReadHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	logger := s.logger.WithRequest(req)
	tracing.AttachRequestToSpan(span, req)

	// determine user ID.
	sessionCtxData, err := s.sessionContextDataFetcher(req)
	if err != nil {
		s.encoderDecoder.EncodeErrorResponse(ctx, res, "unauthenticated", http.StatusUnauthorized)
		return
	}

	requester := sessionCtxData.Requester.UserID
	logger = logger.WithValue(keys.RequesterIDKey, requester)
	tracing.AttachSessionContextDataToSpan(span, sessionCtxData)

	// determine household ID.
	householdID := s.householdIDFetcher(req)
	logger = logger.WithValue(keys.HouseholdIDKey, householdID)
	tracing.AttachHouseholdIDToSpan(span, householdID)

	// fetch household from database.
	var household *types.Household
	if sessionCtxData.ServiceRolePermissionChecker().IsServiceAdmin() {
		household, err = s.householdDataManager.GetHouseholdByID(ctx, householdID)
	} else {
		household, err = s.householdDataManager.GetHousehold(ctx, householdID, requester)
	}

	if errors.Is(err, sql.ErrNoRows) {
		s.encoderDecoder.EncodeNotFoundResponse(ctx, res)
		return
	} else if err != nil {
		observability.AcknowledgeError(err, logger, span, "fetching household from database")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
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

	household.Members = append(admins, plainUsers...)

	// encode our response and peace.
	s.encoderDecoder.RespondWithData(ctx, res, household)
}

// UpdateHandler returns a handler that updates a household.
func (s *service) UpdateHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	logger := s.logger.WithRequest(req)
	tracing.AttachRequestToSpan(span, req)

	// determine user ID.
	sessionCtxData, err := s.sessionContextDataFetcher(req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "fetching session context data")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, "unauthenticated", http.StatusUnauthorized)
		return
	}

	tracing.AttachSessionContextDataToSpan(span, sessionCtxData)
	logger = sessionCtxData.AttachToLogger(logger)

	input := new(types.HouseholdUpdateRequestInput)
	if err = s.encoderDecoder.DecodeRequest(ctx, req, input); err != nil {
		observability.AcknowledgeError(err, logger, span, "decoding request body")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, "invalid request content", http.StatusBadRequest)
		return
	}

	if err = input.ValidateWithContext(ctx); err != nil {
		logger.WithValue(keys.ValidationErrorKey, err).Debug("invalid input attached to request")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, err.Error(), http.StatusBadRequest)
		return
	}

	tracing.AttachSessionContextDataToSpan(span, sessionCtxData)
	requester := sessionCtxData.Requester.UserID
	logger = logger.WithValue(keys.RequesterIDKey, requester)
	input.BelongsToUser = requester

	// determine household ID.
	householdID := s.householdIDFetcher(req)
	logger = logger.WithValue(keys.HouseholdIDKey, householdID)
	tracing.AttachHouseholdIDToSpan(span, householdID)

	// fetch household from database.
	var household *types.Household
	if sessionCtxData.ServiceRolePermissionChecker().IsServiceAdmin() {
		household, err = s.householdDataManager.GetHouseholdByID(ctx, householdID)
	} else {
		household, err = s.householdDataManager.GetHousehold(ctx, householdID, requester)
	}

	if errors.Is(err, sql.ErrNoRows) {
		s.encoderDecoder.EncodeNotFoundResponse(ctx, res)
		return
	} else if err != nil {
		observability.AcknowledgeError(err, logger, span, "fetching household from database")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	// update the data structure.
	household.Update(input)

	// update household in database.
	if err = s.householdDataManager.UpdateHousehold(ctx, household); err != nil {
		observability.AcknowledgeError(err, logger, span, "updating household")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	if s.dataChangesPublisher != nil {
		dcm := &types.DataChangeMessage{
			DataType:                  types.HouseholdDataType,
			EventType:                 types.HouseholdUpdatedCustomerEventType,
			Household:                 household,
			AttributableToUserID:      requester,
			AttributableToHouseholdID: household.ID,
		}
		if err = s.dataChangesPublisher.Publish(ctx, dcm); err != nil {
			observability.AcknowledgeError(err, logger, span, "publishing data change message for created household")
		}
	}

	// encode our response and peace.
	s.encoderDecoder.RespondWithData(ctx, res, household)
}

// ArchiveHandler returns a handler that archives a household.
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

	requester := sessionCtxData.Requester.UserID
	logger = logger.WithValue(keys.RequesterIDKey, requester)
	tracing.AttachSessionContextDataToSpan(span, sessionCtxData)

	// determine household ID.
	householdID := s.householdIDFetcher(req)
	tracing.AttachHouseholdIDToSpan(span, householdID)
	logger = logger.WithValue(keys.HouseholdIDKey, householdID)

	// archive the household in the database.
	err = s.householdDataManager.ArchiveHousehold(ctx, householdID, requester)
	if errors.Is(err, sql.ErrNoRows) {
		s.encoderDecoder.EncodeNotFoundResponse(ctx, res)
		return
	} else if err != nil {
		observability.AcknowledgeError(err, logger, span, "archiving household")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	if s.dataChangesPublisher != nil {
		dcm := &types.DataChangeMessage{
			DataType:                  types.HouseholdDataType,
			EventType:                 types.HouseholdArchivedCustomerEventType,
			AttributableToUserID:      requester,
			AttributableToHouseholdID: householdID,
		}
		if err = s.dataChangesPublisher.Publish(ctx, dcm); err != nil {
			observability.AcknowledgeError(err, logger, span, "publishing data change message for created household")
		}
	}

	s.householdCounter.Decrement(ctx)

	// encode our response and peace.
	res.WriteHeader(http.StatusNoContent)
}

// ModifyMemberPermissionsHandler is our household creation route.
func (s *service) ModifyMemberPermissionsHandler(res http.ResponseWriter, req *http.Request) {
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
	input := new(types.ModifyUserPermissionsInput)
	if err = s.encoderDecoder.DecodeRequest(ctx, req, input); err != nil {
		observability.AcknowledgeError(err, logger, span, "decoding request body")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, "invalid request content", http.StatusBadRequest)
		return
	}

	if err = input.ValidateWithContext(ctx); err != nil {
		logger.WithValue(keys.ValidationErrorKey, err).Debug("invalid input attached to request")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, err.Error(), http.StatusBadRequest)
		return
	}

	requester := sessionCtxData.Requester.UserID
	logger = logger.WithValue(keys.RequesterIDKey, requester)
	tracing.AttachSessionContextDataToSpan(span, sessionCtxData)

	householdID := s.householdIDFetcher(req)
	logger = logger.WithValue(keys.HouseholdIDKey, householdID)
	tracing.AttachHouseholdIDToSpan(span, householdID)

	userID := s.userIDFetcher(req)
	logger = logger.WithValue(keys.UserIDKey, userID)
	tracing.AttachUserIDToSpan(span, userID)

	// create household in database.
	if err = s.householdMembershipDataManager.ModifyUserPermissions(ctx, householdID, userID, input); err != nil {
		observability.AcknowledgeError(err, logger, span, "modifying user permissions")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	if s.dataChangesPublisher != nil {
		dcm := &types.DataChangeMessage{
			DataType:                  types.UserMembershipDataType,
			EventType:                 types.HouseholdMembershipPermissionsUpdatedCustomerEventType,
			AttributableToUserID:      requester,
			HouseholdID:               householdID,
			AttributableToHouseholdID: householdID,
		}
		if err = s.dataChangesPublisher.Publish(ctx, dcm); err != nil {
			observability.AcknowledgeError(err, logger, span, "publishing data change message for created household")
		}
	}

	res.WriteHeader(http.StatusAccepted)
}

// TransferHouseholdOwnershipHandler is our household creation route.
func (s *service) TransferHouseholdOwnershipHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	logger := s.logger.WithRequest(req)
	tracing.AttachRequestToSpan(span, req)

	// determine user ID.
	sessionCtxData, err := s.sessionContextDataFetcher(req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "transferring household ownership")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, "unauthenticated", http.StatusUnauthorized)
		return
	}

	tracing.AttachSessionContextDataToSpan(span, sessionCtxData)
	logger = sessionCtxData.AttachToLogger(logger)

	// read parsed input struct from request body.
	input := new(types.HouseholdOwnershipTransferInput)
	if err = s.encoderDecoder.DecodeRequest(ctx, req, input); err != nil {
		observability.AcknowledgeError(err, logger, span, "decoding request body")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, "invalid request content", http.StatusBadRequest)
		return
	}

	if err = input.ValidateWithContext(ctx); err != nil {
		logger.WithValue(keys.ValidationErrorKey, err).Debug("invalid input attached to request")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, err.Error(), http.StatusBadRequest)
		return
	}

	householdID := s.householdIDFetcher(req)
	tracing.AttachHouseholdIDToSpan(span, householdID)
	logger = logger.WithValue(keys.HouseholdIDKey, householdID)

	requester := sessionCtxData.Requester.UserID
	tracing.AttachSessionContextDataToSpan(span, sessionCtxData)
	logger = logger.WithValue(keys.RequesterIDKey, requester)

	// transfer ownership of household in database.
	if err = s.householdMembershipDataManager.TransferHouseholdOwnership(ctx, householdID, input); err != nil {
		observability.AcknowledgeError(err, logger, span, "transferring household ownership")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	if s.dataChangesPublisher != nil {
		dcm := &types.DataChangeMessage{
			DataType:                  types.HouseholdDataType,
			EventType:                 types.HouseholdOwnershipTransferredCustomerEventType,
			AttributableToUserID:      requester,
			AttributableToHouseholdID: householdID,
		}
		if err = s.dataChangesPublisher.Publish(ctx, dcm); err != nil {
			observability.AcknowledgeError(err, logger, span, "publishing data change message for created household")
		}
	}

	res.WriteHeader(http.StatusAccepted)
}

// RemoveMemberHandler is our household creation route.
func (s *service) RemoveMemberHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	logger := s.logger.WithRequest(req)
	tracing.AttachRequestToSpan(span, req)

	// read parsed input struct from request body.
	reason := req.URL.Query().Get("reason")
	logger = logger.WithValue(keys.ReasonKey, reason)

	// determine user ID.
	sessionCtxData, err := s.sessionContextDataFetcher(req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving session context data")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, "unauthenticated", http.StatusUnauthorized)
		return
	}

	requester := sessionCtxData.Requester.UserID
	logger = logger.WithValue(keys.RequesterIDKey, requester)
	tracing.AttachSessionContextDataToSpan(span, sessionCtxData)

	householdID := s.householdIDFetcher(req)
	logger = logger.WithValue(keys.HouseholdIDKey, householdID)
	tracing.AttachHouseholdIDToSpan(span, householdID)

	userID := s.userIDFetcher(req)
	logger = logger.WithValue(keys.UserIDKey, userID)
	tracing.AttachUserIDToSpan(span, userID)

	// remove user from household in database.
	if err = s.householdMembershipDataManager.RemoveUserFromHousehold(ctx, userID, householdID); err != nil {
		observability.AcknowledgeError(err, logger, span, "removing user from household")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	logger.Info("user removed from household")

	if s.dataChangesPublisher != nil {
		dcm := &types.DataChangeMessage{
			DataType:                  types.UserMembershipDataType,
			EventType:                 types.HouseholdMemberRemovedCustomerEventType,
			AttributableToUserID:      userID,
			AttributableToHouseholdID: householdID,
		}
		if err = s.dataChangesPublisher.Publish(ctx, dcm); err != nil {
			observability.AcknowledgeError(err, logger, span, "publishing data change message for created household")
		}
	}

	res.WriteHeader(http.StatusAccepted)
}

// MarkAsDefaultHouseholdHandler is our household creation route.
func (s *service) MarkAsDefaultHouseholdHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	logger := s.logger.WithRequest(req)
	tracing.AttachRequestToSpan(span, req)

	householdID := s.householdIDFetcher(req)
	logger = logger.WithValue(keys.HouseholdIDKey, householdID)
	tracing.AttachHouseholdIDToSpan(span, householdID)

	// determine user ID.
	sessionCtxData, err := s.sessionContextDataFetcher(req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving session context data")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, "unauthenticated", http.StatusUnauthorized)
		return
	}

	requester := sessionCtxData.Requester.UserID
	logger = logger.WithValue(keys.RequesterIDKey, requester)
	tracing.AttachSessionContextDataToSpan(span, sessionCtxData)

	// mark household as default in database.
	if err = s.householdMembershipDataManager.MarkHouseholdAsUserDefault(ctx, requester, householdID); err != nil {
		observability.AcknowledgeError(err, logger, span, "marking household as default")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	if s.dataChangesPublisher != nil {
		dcm := &types.DataChangeMessage{
			DataType:                  types.UserMembershipDataType,
			EventType:                 types.HouseholdMemberRemovedCustomerEventType,
			AttributableToUserID:      requester,
			AttributableToHouseholdID: householdID,
		}
		if err = s.dataChangesPublisher.Publish(ctx, dcm); err != nil {
			observability.AcknowledgeError(err, logger, span, "publishing data change message for created household")
		}
	}

	res.WriteHeader(http.StatusAccepted)
}
