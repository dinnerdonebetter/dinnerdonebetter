package serversentevents

import (
	"net/http"
)

// StreamSubscriptionHandler subscribes to a stream.
func (s *service) StreamSubscriptionHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	go func() {
		<-ctx.Done()
		return
	}()

	s.eventsServer.HTTPHandler(res, req)
}
