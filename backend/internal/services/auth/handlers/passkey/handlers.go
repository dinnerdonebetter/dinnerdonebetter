package passkey

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/dinnerdonebetter/backend/internal/authentication/cookies"
	authsvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/auth"
	"github.com/dinnerdonebetter/backend/internal/platform/encoding"
	"github.com/dinnerdonebetter/backend/internal/platform/observability"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
	webappauth "github.com/dinnerdonebetter/backend/internal/webapp/auth"
	"github.com/dinnerdonebetter/backend/pkg/client"
)

// AuthOptionsRequest is the request body for passkey authentication options.
type AuthOptionsRequest struct {
	Username string `json:"username"`
}

// AuthOptionsResponse is the response for passkey authentication options.
type AuthOptionsResponse struct {
	Challenge                         string `json:"challenge"`
	PublicKeyCredentialRequestOptions []byte `json:"publicKeyCredentialRequestOptions"`
}

// AuthVerifyRequest is the request body for passkey authentication verify.
// AssertionResponse is the raw JSON bytes of the WebAuthn credential (id, rawId, type, response).
type AuthVerifyRequest struct {
	Challenge         string          `json:"challenge"`
	Username          string          `json:"username"`
	AssertionResponse json.RawMessage `json:"assertionResponse"`
}

// Handlers provides HTTP handlers for passkey authentication.
type Handlers struct {
	tracer              tracing.Tracer
	logger              logging.Logger
	encoder             encoding.ServerEncoderDecoder
	cookieManager       cookies.Manager
	cookieConfig        *cookies.Config
	buildUnauthedClient func(context.Context) (client.Client, error)
}

// NewHandlers creates passkey HTTP handlers.
func NewHandlers(
	tracer tracing.Tracer,
	logger logging.Logger,
	encoder encoding.ServerEncoderDecoder,
	cookieManager cookies.Manager,
	cookieConfig *cookies.Config,
	buildUnauthedClient func(context.Context) (client.Client, error),
) *Handlers {
	return &Handlers{
		tracer:              tracer,
		logger:              logger,
		encoder:             encoder,
		cookieManager:       cookieManager,
		cookieConfig:        cookieConfig,
		buildUnauthedClient: buildUnauthedClient,
	}
}

// AuthOptionsHandler handles POST /auth/passkey/authentication/options.
func (h *Handlers) AuthOptionsHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := h.tracer.StartSpan(req.Context())
	defer span.End()

	logger := h.logger.WithRequest(req)

	var body AuthOptionsRequest
	if err := h.encoder.DecodeRequest(ctx, req, &body); err != nil {
		observability.AcknowledgeError(err, logger, span, "decoding passkey options request")
		h.encoder.EncodeResponseWithStatus(ctx, res, map[string]string{"error": "invalid request"}, http.StatusBadRequest)
		return
	}

	unauthedClient, err := h.buildUnauthedClient(ctx)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "building gRPC client")
		h.encoder.EncodeResponseWithStatus(ctx, res, map[string]string{"error": "service unavailable"}, http.StatusServiceUnavailable)
		return
	}

	optsRes, err := unauthedClient.BeginPasskeyAuthentication(ctx, &authsvc.BeginPasskeyAuthenticationRequest{
		Username: strings.TrimSpace(body.Username),
	})
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "beginning passkey authentication")
		h.encoder.EncodeResponseWithStatus(ctx, res, map[string]string{"error": "failed to get passkey options"}, http.StatusInternalServerError)
		return
	}

	h.encoder.EncodeResponseWithStatus(ctx, res, AuthOptionsResponse{
		PublicKeyCredentialRequestOptions: optsRes.PublicKeyCredentialRequestOptions,
		Challenge:                         optsRes.Challenge,
	}, http.StatusOK)
}

// AuthVerifyHandler handles POST /auth/passkey/authentication/verify.
func (h *Handlers) AuthVerifyHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := h.tracer.StartSpan(req.Context())
	defer span.End()

	logger := h.logger.WithRequest(req)

	var body AuthVerifyRequest
	if err := h.encoder.DecodeRequest(ctx, req, &body); err != nil {
		observability.AcknowledgeError(err, logger, span, "decoding passkey verify request")
		h.encoder.EncodeResponseWithStatus(ctx, res, map[string]string{"error": "invalid request"}, http.StatusBadRequest)
		return
	}

	if len([]byte(body.AssertionResponse)) == 0 || body.Challenge == "" {
		h.encoder.EncodeResponseWithStatus(ctx, res, map[string]string{"error": "assertion_response and challenge are required"}, http.StatusBadRequest)
		return
	}

	unauthedClient, err := h.buildUnauthedClient(ctx)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "building gRPC client")
		h.encoder.EncodeResponseWithStatus(ctx, res, map[string]string{"error": "service unavailable"}, http.StatusServiceUnavailable)
		return
	}

	tokenRes, err := unauthedClient.FinishPasskeyAuthentication(ctx, &authsvc.FinishPasskeyAuthenticationRequest{
		AssertionResponse: []byte(body.AssertionResponse),
		Challenge:         body.Challenge,
		Username:          strings.TrimSpace(body.Username),
	})
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "finishing passkey authentication")
		h.encoder.EncodeResponseWithStatus(ctx, res, map[string]string{"error": "authentication failed"}, http.StatusUnauthorized)
		return
	}

	encodedCookie, err := h.cookieManager.Encode(ctx, h.cookieConfig.CookieName, &webappauth.AuthPayload{AccessToken: tokenRes.Result.AccessToken})
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "encoding auth cookie")
		h.encoder.EncodeResponseWithStatus(ctx, res, map[string]string{"error": "failed to complete login"}, http.StatusInternalServerError)
		return
	}

	http.SetCookie(res, webappauth.BuildCookie(h.cookieConfig, encodedCookie))

	res.Header().Set("HX-Redirect", "/")
	h.encoder.EncodeResponseWithStatus(ctx, res, map[string]string{"redirect": "/"}, http.StatusOK)
}
