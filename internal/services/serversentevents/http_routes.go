package serversentevents

import (
	"net/http"

	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
)

// StreamSubscriptionHandler subscribes to a stream.
func (s *service) StreamSubscriptionHandler(res http.ResponseWriter, req *http.Request) {
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

	s.eventsServer.CreateStream(sessionCtxData.Requester.UserID)
	req.URL.Query().Set("stream", sessionCtxData.Requester.UserID)

	go func() {
		<-ctx.Done()
		return
	}()

	s.eventsServer.HTTPHandler(res, req)
}
