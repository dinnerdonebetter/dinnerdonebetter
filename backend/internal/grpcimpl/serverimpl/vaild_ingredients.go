package serverimpl

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/grpc/messages"
	"github.com/dinnerdonebetter/backend/internal/grpcimpl/converters"
	"github.com/dinnerdonebetter/backend/internal/lib/observability"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/tracing"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Server) CreateValidIngredient(ctx context.Context, input *messages.ValidIngredientCreationRequestInput) (*messages.ValidIngredient, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.Clone()

	if err := validation.ValidateStructWithContext(
		ctx,
		input,
		validation.Field(&input.Name, validation.Required),
	); err != nil {
		return nil, observability.PrepareError(err, span, "validating input")
	}

	created, err := s.dataManager.CreateValidIngredient(ctx, converters.ConvertValidIngredientCreationRequestInputToValidIngredientDatabaseCreationInput(input))
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "creating valid ingredient")
	}

	tracing.AttachToSpan(span, keys.ValidIngredientIDKey, created.ID)
	output := converters.ConvertValidIngredientToProtobuf(created)

	return output, nil
}

func (s *Server) GetValidIngredient(ctx context.Context, request *messages.GetValidIngredientRequest) (*messages.ValidIngredient, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithValue(keys.ValidIngredientIDKey, request.ValidIngredientID)
	tracing.AttachToSpan(span, keys.ValidIngredientIDKey, request.ValidIngredientID)

	validIngredient, err := s.dataManager.GetValidIngredient(ctx, request.ValidIngredientID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "getting valid ingredient")
	}

	output := converters.ConvertValidIngredientToProtobuf(validIngredient)

	return output, nil
}

func (s *Server) GetRandomValidIngredient(ctx context.Context, _ *emptypb.Empty) (*messages.ValidIngredient, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	ingredient, err := s.dataManager.GetRandomValidIngredient(ctx)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, s.logger, span, "getting random valid ingredient")
	}

	tracing.AttachToSpan(span, keys.ValidIngredientIDKey, ingredient.ID)
	output := converters.ConvertValidIngredientToProtobuf(ingredient)

	return output, nil
}

func (s *Server) GetValidIngredients(ctx context.Context, request *messages.GetValidIngredientsRequest) (*messages.ValidIngredient, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) SearchForValidIngredients(ctx context.Context, request *messages.SearchForValidIngredientsRequest) (*messages.ValidIngredient, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) SearchValidIngredientsByPreparation(ctx context.Context, request *messages.SearchValidIngredientsByPreparationRequest) (*messages.ValidIngredient, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) UpdateValidIngredient(ctx context.Context, request *messages.UpdateValidIngredientRequest) (*messages.ValidIngredient, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithValue(keys.ValidIngredientIDKey, request.ValidIngredientID)
	tracing.AttachToSpan(span, keys.ValidIngredientIDKey, request.ValidIngredientID)

	updated := converters.ConvertUpdateValidIngredientRequestToValidIngredient(request)
	if err := s.dataManager.UpdateValidIngredient(ctx, updated); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "updating valid ingredient")
	}

	return nil, nil
}

func (s *Server) ArchiveValidIngredient(ctx context.Context, request *messages.ArchiveValidIngredientRequest) (*messages.ValidIngredient, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithValue(keys.ValidIngredientIDKey, request.ValidIngredientID)
	tracing.AttachToSpan(span, keys.ValidIngredientIDKey, request.ValidIngredientID)

	validIngredient, err := s.dataManager.GetValidIngredient(ctx, request.ValidIngredientID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "getting valid ingredient")
	}

	if err = s.dataManager.ArchiveValidIngredient(ctx, request.ValidIngredientID); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "archiving valid ingredient")
	}

	return converters.ConvertValidIngredientToProtobuf(validIngredient), nil
}
