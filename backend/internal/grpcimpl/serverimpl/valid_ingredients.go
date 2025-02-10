package serverimpl

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/grpc/messages"
	"github.com/dinnerdonebetter/backend/internal/grpcimpl/converters"
	"github.com/dinnerdonebetter/backend/internal/lib/observability"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/types"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func (s *Server) CreateValidIngredient(ctx context.Context, request *messages.CreateValidIngredientRequest) (*messages.CreateValidIngredientResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	if err := validation.ValidateStructWithContext(
		ctx,
		request,
		validation.Field(&request.Name, validation.Required),
	); err != nil {
		return nil, observability.PrepareError(err, span, "validating input")
	}

	created, err := s.dataManager.CreateValidIngredient(ctx, converters.ConvertCreateValidIngredientRequestToValidIngredientDatabaseCreationInput(request))
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "creating valid ingredient")
	}

	tracing.AttachToSpan(span, keys.ValidIngredientIDKey, created.ID)
	output := converters.ConvertValidIngredientToProtobuf(created)

	return &messages.CreateValidIngredientResponse{Result: output}, nil
}

func (s *Server) GetValidIngredient(ctx context.Context, request *messages.GetValidIngredientRequest) (*messages.GetValidIngredientResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	logger = logger.WithValue(keys.ValidIngredientIDKey, request.ValidIngredientID)
	tracing.AttachToSpan(span, keys.ValidIngredientIDKey, request.ValidIngredientID)

	validIngredient, err := s.dataManager.GetValidIngredient(ctx, request.ValidIngredientID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "getting valid ingredient")
	}

	output := converters.ConvertValidIngredientToProtobuf(validIngredient)

	return &messages.GetValidIngredientResponse{Result: output}, nil
}

func (s *Server) GetValidIngredients(ctx context.Context, request *messages.GetValidIngredientsRequest) (*messages.GetValidIngredientsResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	validIngredients, err := s.dataManager.GetValidIngredients(ctx, nil)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "getting valid ingredients")
	}

	outbound := &messages.GetValidIngredientsResponse{
		Meta: &messages.ResponseMeta{
			TraceID:            "",
			CurrentHouseholdID: "",
		},
		Results: []*messages.ValidIngredient{},
	}

	for _, validIngredient := range validIngredients.Data {
		outbound.Results = append(outbound.Results, converters.ConvertValidIngredientToProtobuf(validIngredient))
	}

	return outbound, nil
}

func (s *Server) GetRandomValidIngredient(ctx context.Context, _ *messages.GetRandomValidIngredientRequest) (*messages.GetRandomValidIngredientResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	ingredient, err := s.dataManager.GetRandomValidIngredient(ctx)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "getting random valid ingredient")
	}

	tracing.AttachToSpan(span, keys.ValidIngredientIDKey, ingredient.ID)
	output := converters.ConvertValidIngredientToProtobuf(ingredient)

	return &messages.GetRandomValidIngredientResponse{Result: output}, nil
}

func (s *Server) SearchForValidIngredients(ctx context.Context, request *messages.SearchForValidIngredientsRequest) (*messages.SearchForValidIngredientsResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	query := request.Query
	logger = logger.WithValue(keys.SearchQueryKey, query)
	tracing.AttachToSpan(span, keys.SearchQueryKey, query)

	result := &messages.SearchForValidIngredientsResponse{
		Meta: &messages.ResponseMeta{
			TraceID:            "",
			CurrentHouseholdID: "",
		},
		Filter:  request.Filter,
		Results: []*messages.ValidIngredient{},
	}

	var (
		validIngredients []*types.ValidIngredient
	)
	if s.config.Services.ValidEnumerations.UseSearchService {
		validIngredientSubsets, err := s.validIngredientSearchIndex.Search(ctx, query)
		if err != nil {
			return nil, observability.PrepareAndLogError(err, logger, span, "searching for valid ingredients")
		}

		ids := []string{}
		for _, validIngredientSubset := range validIngredientSubsets {
			ids = append(ids, validIngredientSubset.ID)
		}

		validIngredients, err = s.dataManager.GetValidIngredientsWithIDs(ctx, ids)
		if err != nil {
			return nil, observability.PrepareAndLogError(err, logger, span, "reading valid ingredients by IDs")
		}
	} else {
		validIngredientsResult, err := s.dataManager.SearchForValidIngredients(ctx, query, nil)
		if err != nil {
			return nil, observability.PrepareAndLogError(err, logger, span, "reading valid ingredients by IDs")
		}

		validIngredients = validIngredientsResult.Data
	}

	for _, validIngredient := range validIngredients {
		result.Results = append(result.Results, converters.ConvertValidIngredientToProtobuf(validIngredient))
	}

	return result, nil
}

func (s *Server) SearchValidIngredientsByPreparation(ctx context.Context, request *messages.SearchValidIngredientsByPreparationRequest) (*messages.SearchValidIngredientsByPreparationResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	logger = logger.WithValue(keys.SearchQueryKey, request.Query)
	tracing.AttachToSpan(span, keys.SearchQueryKey, request.Query)

	validIngredients, err := s.dataManager.SearchForValidIngredientsForPreparation(ctx, request.ValidPreparationID, request.Query, nil)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "reading valid ingredients by IDs")
	}

	output := &messages.SearchValidIngredientsByPreparationResponse{
		Meta:    s.buildResponseMeta(span),
		Results: []*messages.ValidIngredient{},
	}

	for _, validIngredient := range validIngredients.Data {
		output.Results = append(output.Results, converters.ConvertValidIngredientToProtobuf(validIngredient))
	}

	return output, nil
}

func (s *Server) UpdateValidIngredient(ctx context.Context, request *messages.UpdateValidIngredientRequest) (*messages.UpdateValidIngredientResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	logger = logger.WithValue(keys.ValidIngredientIDKey, request.ValidIngredientID)
	tracing.AttachToSpan(span, keys.ValidIngredientIDKey, request.ValidIngredientID)

	updated := converters.ConvertUpdateValidIngredientRequestToValidIngredient(request)
	if err := s.dataManager.UpdateValidIngredient(ctx, updated); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "updating valid ingredient")
	}

	return &messages.UpdateValidIngredientResponse{}, nil
}

func (s *Server) ArchiveValidIngredient(ctx context.Context, request *messages.ArchiveValidIngredientRequest) (*messages.ArchiveValidIngredientResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	logger = logger.WithValue(keys.ValidIngredientIDKey, request.ValidIngredientID)
	tracing.AttachToSpan(span, keys.ValidIngredientIDKey, request.ValidIngredientID)

	if err := s.dataManager.ArchiveValidIngredient(ctx, request.ValidIngredientID); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "archiving valid ingredient")
	}

	return &messages.ArchiveValidIngredientResponse{}, nil
}
