package admin

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/prixfixeco/api_server/internal/observability"
	"github.com/prixfixeco/api_server/internal/observability/keys"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/pkg/types"
)

const (
	// UserIDURIParamKey is used to refer to user IDs in router params.
	UserIDURIParamKey = "userID"
)

// UserAccountStatusChangeHandler changes a user's status.
func (s *service) UserAccountStatusChangeHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	logger := s.logger.WithRequest(req)
	tracing.AttachRequestToSpan(span, req)

	sessionCtxData, err := s.sessionContextDataFetcher(req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving session context data")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	tracing.AttachSessionContextDataToSpan(span, sessionCtxData)
	logger = sessionCtxData.AttachToLogger(logger)

	// read parsed input struct from request body.
	input := new(types.UserAccountStatusUpdateInput)
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

	logger = logger.WithValue("new_status", input.NewStatus)

	if !sessionCtxData.Requester.ServicePermissions.CanUpdateUserAccountStatuses() {
		// this should never happen in production
		s.encoderDecoder.EncodeErrorResponse(ctx, res, "inadequate permissions for route", http.StatusForbidden)
		return
	}

	requester := sessionCtxData.Requester.UserID
	logger = logger.WithValue("ban_giver", requester)
	logger = logger.WithValue("status_change_recipient", input.TargetUserID)

	if err = s.userDB.UpdateUserAccountStatus(ctx, input.TargetUserID, input); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.encoderDecoder.EncodeNotFoundResponse(ctx, res)
		} else {
			observability.AcknowledgeError(err, logger, span, "retrieving session context data")
			s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		}

		return
	}

	s.encoderDecoder.EncodeResponseWithStatus(ctx, res, nil, http.StatusAccepted)
}
