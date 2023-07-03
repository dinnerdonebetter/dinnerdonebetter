package wasm

import (
	_ "embed"
	"net/http"

	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
)

//go:embed assets/wasm_exec.js
var wasmExecJS []byte

// WASMHandler is our valid ingredient creation route.
func (s *Service) WASMHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	logger := s.logger.WithRequest(req)
	tracing.AttachRequestToSpan(span, req)

	logger.Info("WASM route")

	if _, err := res.Write(wasmExecJS); err != nil {
		observability.AcknowledgeError(err, logger, span, "failed to write wasm_exec.js")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, "failed to write wasm_exec.js", http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusOK)
}
