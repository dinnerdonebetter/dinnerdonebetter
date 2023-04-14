package vendorproxy

import (
	"net/http"

	"github.com/prixfixeco/backend/internal/observability"
	"github.com/prixfixeco/backend/internal/observability/tracing"
	"github.com/prixfixeco/backend/pkg/types"
)

const (
	// FeatureFlagURIParamKey is a standard string that we'll use to refer to feature flag identifiers with.
	FeatureFlagURIParamKey = "featureFlag"
)

// FeatureFlagStatusHandler is our feature flag management route.
func (s *service) FeatureFlagStatusHandler(res http.ResponseWriter, req *http.Request) {
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

	featureFlag := s.featureFlagURLFetcher(req)

	result, err := s.featureFlagManager.CanUseFeature(ctx, sessionCtxData.Requester.UserID, featureFlag)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "checking feature flag")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, "internal error", http.StatusInternalServerError)
		return
	}

	responseCode := http.StatusOK
	if !result {
		responseCode = http.StatusForbidden
	}

	res.WriteHeader(responseCode)
}

type analyticsTrackRequest struct {
	Properties map[string]any `json:"properties"`
	Event      string         `json:"event"`
}

// AnalyticsTrackHandler returns a handler that proxies requests to the analytics service.
func (s *service) AnalyticsTrackHandler(res http.ResponseWriter, req *http.Request) {
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

	var x *analyticsTrackRequest
	if err = s.encoderDecoder.DecodeRequest(ctx, req, &x); err != nil {
		observability.AcknowledgeError(err, logger, span, "decoding request")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, "", http.StatusBadRequest)
		return
	}

	if err = s.eventReporter.EventOccurred(ctx, types.CustomerEventType(x.Event), sessionCtxData.Requester.UserID, x.Properties); err != nil {
		observability.AcknowledgeError(err, logger, span, "reporting event")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, "", http.StatusInternalServerError)
		return
	}
}

type analyticsIdentifyRequest struct {
	Properties map[string]any `json:"properties"`
}

// AnalyticsIdentifyHandler returns a handler that proxies requests to the analytics service.
func (s *service) AnalyticsIdentifyHandler(res http.ResponseWriter, req *http.Request) {
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

	var x *analyticsIdentifyRequest
	if err = s.encoderDecoder.DecodeRequest(ctx, req, &x); err != nil {
		observability.AcknowledgeError(err, logger, span, "decoding request")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, "", http.StatusBadRequest)
		return
	}

	if err = s.eventReporter.AddUser(ctx, sessionCtxData.Requester.UserID, x.Properties); err != nil {
		observability.AcknowledgeError(err, logger, span, "reporting event")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, "", http.StatusInternalServerError)
		return
	}
}
