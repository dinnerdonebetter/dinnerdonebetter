package websockets

import (
	"net/http"

	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
)

// WebsocketSubscriptionHandler is our subscription route.
func (s *service) WebsocketSubscriptionHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	logger := s.logger.WithRequest(req)
	tracing.AttachRequestToSpan(span, req)

	logger.Debug("notification subscription route hit")

	// determine user ID.
	sessionCtxData, err := s.sessionContextDataFetcher(req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving session context data")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, "unauthenticated", http.StatusUnauthorized)
		return
	}

	tracing.AttachSessionContextDataToSpan(span, sessionCtxData)
	logger = sessionCtxData.AttachToLogger(logger)

	cookie, err := req.Cookie(s.authConfig.Cookies.Name)
	if err != nil {
		logger.Error(err, "checking websocket subscription request for cookies")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, "unauthenticated", http.StatusUnauthorized)
		return
	}

	// we have to amend the cookie here, lest it get set to an invalid path and domain.
	cookie.Domain = s.authConfig.Cookies.Domain
	cookie.Path = "/"

	wsHeader := http.Header{}
	wsHeader.Add("Set-Cookie", cookie.String())
	wsHeader.Add("Authorization", req.Header.Get("Authorization"))

	conn, err := s.websocketConnectionUpgrader.Upgrade(res, req, wsHeader)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving session context data")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	s.connectionsHat.Lock()
	defer s.connectionsHat.Unlock()

	_, ok := s.connections[sessionCtxData.Requester.UserID]
	if ok {
		s.connections[sessionCtxData.Requester.UserID] = append(s.connections[sessionCtxData.Requester.UserID], conn)
	} else {
		s.connections[sessionCtxData.Requester.UserID] = []websocketConnection{conn}
	}
}
