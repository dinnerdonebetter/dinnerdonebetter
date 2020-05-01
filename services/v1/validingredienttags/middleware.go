package validingredienttags

import (
	"context"
	"net/http"

	"gitlab.com/prixfixe/prixfixe/internal/v1/tracing"
	models "gitlab.com/prixfixe/prixfixe/models/v1"
)

// CreationInputMiddleware is a middleware for fetching, parsing, and attaching an ValidIngredientTagInput struct from a request.
func (s *Service) CreationInputMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		x := new(models.ValidIngredientTagCreationInput)
		ctx, span := tracing.StartSpan(req.Context(), "CreationInputMiddleware")
		defer span.End()

		logger := s.logger.WithRequest(req)

		if err := s.encoderDecoder.DecodeRequest(req, x); err != nil {
			logger.Error(err, "error encountered decoding request body")
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		ctx = context.WithValue(ctx, CreateMiddlewareCtxKey, x)
		next.ServeHTTP(res, req.WithContext(ctx))
	})
}

// UpdateInputMiddleware is a middleware for fetching, parsing, and attaching an ValidIngredientTagInput struct from a request.
// This is the same as the creation one, but that won't always be the case.
func (s *Service) UpdateInputMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		x := new(models.ValidIngredientTagUpdateInput)
		ctx, span := tracing.StartSpan(req.Context(), "UpdateInputMiddleware")
		defer span.End()

		logger := s.logger.WithRequest(req)

		if err := s.encoderDecoder.DecodeRequest(req, x); err != nil {
			logger.Error(err, "error encountered decoding request body")
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		ctx = context.WithValue(ctx, UpdateMiddlewareCtxKey, x)
		next.ServeHTTP(res, req.WithContext(ctx))
	})
}
