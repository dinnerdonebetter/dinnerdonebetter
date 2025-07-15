package grpc

import (
	"context"

	grpcconverters "github.com/dinnerdonebetter/backend/internal/grpc/converters"
	"github.com/dinnerdonebetter/backend/internal/grpc/generated/services/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/platform/observability"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
	mpgrpcconverters "github.com/dinnerdonebetter/backend/internal/services/mealplanning/grpc/converters"

	"google.golang.org/grpc/codes"
)

func (s *ServiceImpl) ArchiveValidIngredient(ctx context.Context, request *mealplanning.ArchiveValidIngredientRequest) (*mealplanning.ArchiveValidIngredientResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValue(keys.ValidIngredientIDKey, request.ValidIngredientID)
	tracing.AttachToSpan(span, keys.ValidIngredientIDKey, request.ValidIngredientID)

	if err := s.validEnumerationsManager.ArchiveValidIngredient(ctx, request.ValidIngredientID); err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "archiving valid ingredient")
	}

	res := &mealplanning.ArchiveValidIngredientResponse{}

	return res, nil
}

func (s *ServiceImpl) ArchiveValidIngredientGroup(ctx context.Context, request *mealplanning.ArchiveValidIngredientGroupRequest) (*mealplanning.ArchiveValidIngredientGroupResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValue(keys.ValidIngredientGroupIDKey, request.ValidIngredientGroupID)
	tracing.AttachToSpan(span, keys.ValidIngredientGroupIDKey, request.ValidIngredientGroupID)

	if err := s.validEnumerationsManager.ArchiveValidIngredientGroup(ctx, request.ValidIngredientGroupID); err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "archiving valid ingredient group")
	}

	res := &mealplanning.ArchiveValidIngredientGroupResponse{}

	return res, nil
}

func (s *ServiceImpl) ArchiveValidIngredientMeasurementUnit(ctx context.Context, request *mealplanning.ArchiveValidIngredientMeasurementUnitRequest) (*mealplanning.ArchiveValidIngredientMeasurementUnitResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValue(keys.ValidIngredientMeasurementUnitIDKey, request.ValidIngredientMeasurementUnitID)
	tracing.AttachToSpan(span, keys.ValidIngredientMeasurementUnitIDKey, request.ValidIngredientMeasurementUnitID)

	if err := s.validEnumerationsManager.ArchiveValidIngredientMeasurementUnit(ctx, request.ValidIngredientMeasurementUnitID); err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "archiving valid ingredient measurement unit")
	}

	res := &mealplanning.ArchiveValidIngredientMeasurementUnitResponse{}

	return res, nil
}

func (s *ServiceImpl) ArchiveValidIngredientPreparation(ctx context.Context, request *mealplanning.ArchiveValidIngredientPreparationRequest) (*mealplanning.ArchiveValidIngredientPreparationResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValue(keys.ValidIngredientPreparationIDKey, request.ValidIngredientPreparationID)
	tracing.AttachToSpan(span, keys.ValidIngredientPreparationIDKey, request.ValidIngredientPreparationID)

	if err := s.validEnumerationsManager.ArchiveValidIngredientPreparation(ctx, request.ValidIngredientPreparationID); err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "archiving valid ingredient preparation")
	}

	res := &mealplanning.ArchiveValidIngredientPreparationResponse{}

	return res, nil
}

func (s *ServiceImpl) ArchiveValidIngredientState(ctx context.Context, request *mealplanning.ArchiveValidIngredientStateRequest) (*mealplanning.ArchiveValidIngredientStateResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValue(keys.ValidIngredientStateIDKey, request.ValidIngredientStateID)
	tracing.AttachToSpan(span, keys.ValidIngredientStateIDKey, request.ValidIngredientStateID)

	if err := s.validEnumerationsManager.ArchiveValidIngredientState(ctx, request.ValidIngredientStateID); err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "archiving valid ingredient state")
	}

	res := &mealplanning.ArchiveValidIngredientStateResponse{}

	return res, nil
}

func (s *ServiceImpl) ArchiveValidIngredientStateIngredient(ctx context.Context, request *mealplanning.ArchiveValidIngredientStateIngredientRequest) (*mealplanning.ArchiveValidIngredientStateIngredientResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValue(keys.ValidIngredientStateIngredientIDKey, request.ValidIngredientStateIngredientID)
	tracing.AttachToSpan(span, keys.ValidIngredientStateIngredientIDKey, request.ValidIngredientStateIngredientID)

	if err := s.validEnumerationsManager.ArchiveValidIngredientStateIngredient(ctx, request.ValidIngredientStateIngredientID); err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "archiving valid ingredient state ingredient")
	}

	res := &mealplanning.ArchiveValidIngredientStateIngredientResponse{}

	return res, nil
}

func (s *ServiceImpl) ArchiveValidInstrument(ctx context.Context, request *mealplanning.ArchiveValidInstrumentRequest) (*mealplanning.ArchiveValidInstrumentResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValue(keys.ValidInstrumentIDKey, request.ValidInstrumentID)
	tracing.AttachToSpan(span, keys.ValidInstrumentIDKey, request.ValidInstrumentID)

	if err := s.validEnumerationsManager.ArchiveValidInstrument(ctx, request.ValidInstrumentID); err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "archiving valid instrument")
	}

	res := &mealplanning.ArchiveValidInstrumentResponse{}

	return res, nil
}

func (s *ServiceImpl) ArchiveValidMeasurementUnit(ctx context.Context, request *mealplanning.ArchiveValidMeasurementUnitRequest) (*mealplanning.ArchiveValidMeasurementUnitResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValue(keys.ValidMeasurementUnitIDKey, request.ValidMeasurementUnitID)
	tracing.AttachToSpan(span, keys.ValidMeasurementUnitIDKey, request.ValidMeasurementUnitID)

	if err := s.validEnumerationsManager.ArchiveValidMeasurementUnit(ctx, request.ValidMeasurementUnitID); err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "archiving valid measurement unit")
	}

	res := &mealplanning.ArchiveValidMeasurementUnitResponse{}

	return res, nil
}

func (s *ServiceImpl) ArchiveValidMeasurementUnitConversion(ctx context.Context, request *mealplanning.ArchiveValidMeasurementUnitConversionRequest) (*mealplanning.ArchiveValidMeasurementUnitConversionResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValue(keys.ValidMeasurementUnitConversionIDKey, request.ValidMeasurementUnitConversionID)
	tracing.AttachToSpan(span, keys.ValidMeasurementUnitConversionIDKey, request.ValidMeasurementUnitConversionID)

	if err := s.validEnumerationsManager.ArchiveValidMeasurementUnitConversion(ctx, request.ValidMeasurementUnitConversionID); err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "archiving valid measurement unit conversion")
	}

	res := &mealplanning.ArchiveValidMeasurementUnitConversionResponse{}

	return res, nil
}

func (s *ServiceImpl) ArchiveValidPreparation(ctx context.Context, request *mealplanning.ArchiveValidPreparationRequest) (*mealplanning.ArchiveValidPreparationResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValue(keys.ValidPreparationIDKey, request.ValidPreparationID)
	tracing.AttachToSpan(span, keys.ValidPreparationIDKey, request.ValidPreparationID)

	if err := s.validEnumerationsManager.ArchiveValidPreparation(ctx, request.ValidPreparationID); err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "archiving valid preparation")
	}

	res := &mealplanning.ArchiveValidPreparationResponse{}

	return res, nil
}

func (s *ServiceImpl) ArchiveValidPreparationInstrument(ctx context.Context, request *mealplanning.ArchiveValidPreparationInstrumentRequest) (*mealplanning.ArchiveValidPreparationInstrumentResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValue(keys.ValidPreparationInstrumentIDKey, request.ValidPreparationInstrumentID)
	tracing.AttachToSpan(span, keys.ValidPreparationInstrumentIDKey, request.ValidPreparationInstrumentID)

	if err := s.validEnumerationsManager.ArchiveValidPreparationInstrument(ctx, request.ValidPreparationInstrumentID); err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "archiving valid preparation instrument")
	}

	res := &mealplanning.ArchiveValidPreparationInstrumentResponse{}

	return res, nil
}

func (s *ServiceImpl) ArchiveValidPreparationVessel(ctx context.Context, request *mealplanning.ArchiveValidPreparationVesselRequest) (*mealplanning.ArchiveValidPreparationVesselResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValue(keys.ValidPreparationVesselIDKey, request.ValidPreparationVesselID)
	tracing.AttachToSpan(span, keys.ValidPreparationVesselIDKey, request.ValidPreparationVesselID)

	if err := s.validEnumerationsManager.ArchiveValidPreparationVessel(ctx, request.ValidPreparationVesselID); err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "archiving valid preparation vessel")
	}

	res := &mealplanning.ArchiveValidPreparationVesselResponse{}

	return res, nil
}

func (s *ServiceImpl) ArchiveValidVessel(ctx context.Context, request *mealplanning.ArchiveValidVesselRequest) (*mealplanning.ArchiveValidVesselResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValue(keys.ValidVesselIDKey, request.ValidVesselID)
	tracing.AttachToSpan(span, keys.ValidVesselIDKey, request.ValidVesselID)

	if err := s.validEnumerationsManager.ArchiveValidVessel(ctx, request.ValidVesselID); err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "archiving valid vessel")
	}

	res := &mealplanning.ArchiveValidVesselResponse{}

	return res, nil
}

func (s *ServiceImpl) CreateValidIngredient(ctx context.Context, request *mealplanning.CreateValidIngredientRequest) (*mealplanning.CreateValidIngredientResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	created, err := s.validEnumerationsManager.CreateValidIngredient(ctx, mpgrpcconverters.ConvertGRPCCreateValidIngredientRequestToValidIngredientCreationRequestInput(request))
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "creating valid ingredient")
	}
	tracing.AttachToSpan(span, keys.ValidIngredientIDKey, created.ID)

	result := &mealplanning.CreateValidIngredientResponse{
		Result: mpgrpcconverters.ConvertValidIngredientToGRPCValidIngredient(created),
	}

	return result, nil
}

func (s *ServiceImpl) CreateValidIngredientGroup(ctx context.Context, request *mealplanning.CreateValidIngredientGroupRequest) (*mealplanning.CreateValidIngredientGroupResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	created, err := s.validEnumerationsManager.CreateValidIngredientGroup(ctx, mpgrpcconverters.ConvertGRPCCreateValidIngredientGroupRequestToValidIngredientGroupCreationRequestInput(request))
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "creating valid ingredient")
	}
	tracing.AttachToSpan(span, keys.ValidIngredientGroupIDKey, created.ID)

	result := &mealplanning.CreateValidIngredientGroupResponse{
		Result: mpgrpcconverters.ConvertValidIngredientGroupToGRPCValidIngredientGroup(created),
	}

	return result, nil
}

func (s *ServiceImpl) CreateValidIngredientMeasurementUnit(ctx context.Context, request *mealplanning.CreateValidIngredientMeasurementUnitRequest) (*mealplanning.CreateValidIngredientMeasurementUnitResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	created, err := s.validEnumerationsManager.CreateValidIngredientMeasurementUnit(ctx, mpgrpcconverters.ConvertGRPCCreateValidIngredientMeasurementUnitRequestToValidIngredientMeasurementUnitCreationRequestInput(request))
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "creating valid ingredient measurement unit")
	}
	tracing.AttachToSpan(span, keys.ValidIngredientMeasurementUnitIDKey, created.ID)

	result := &mealplanning.CreateValidIngredientMeasurementUnitResponse{
		Result: mpgrpcconverters.ConvertValidIngredientMeasurementUnitToGRPCValidIngredientMeasurementUnit(created),
	}

	return result, nil
}

func (s *ServiceImpl) CreateValidIngredientPreparation(ctx context.Context, request *mealplanning.CreateValidIngredientPreparationRequest) (*mealplanning.CreateValidIngredientPreparationResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	created, err := s.validEnumerationsManager.CreateValidIngredientPreparation(ctx, mpgrpcconverters.ConvertGRPCCreateValidIngredientPreparationRequestToValidIngredientPreparationCreationRequestInput(request))
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "creating valid ingredient preparation")
	}
	tracing.AttachToSpan(span, keys.ValidIngredientPreparationIDKey, created.ID)

	result := &mealplanning.CreateValidIngredientPreparationResponse{
		Result: mpgrpcconverters.ConvertValidIngredientPreparationToGRPCValidIngredientPreparation(created),
	}

	return result, nil
}

func (s *ServiceImpl) CreateValidIngredientState(ctx context.Context, request *mealplanning.CreateValidIngredientStateRequest) (*mealplanning.CreateValidIngredientStateResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	created, err := s.validEnumerationsManager.CreateValidIngredientState(ctx, mpgrpcconverters.ConvertGRPCCreateValidIngredientStateRequestToValidIngredientStateCreationRequestInput(request))
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "creating valid ingredient state")
	}
	tracing.AttachToSpan(span, keys.ValidIngredientStateIDKey, created.ID)

	result := &mealplanning.CreateValidIngredientStateResponse{
		Result: mpgrpcconverters.ConvertValidIngredientStateToGRPCValidIngredientState(created),
	}

	return result, nil
}

func (s *ServiceImpl) CreateValidIngredientStateIngredient(ctx context.Context, request *mealplanning.CreateValidIngredientStateIngredientRequest) (*mealplanning.CreateValidIngredientStateIngredientResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	created, err := s.validEnumerationsManager.CreateValidIngredientStateIngredient(ctx, mpgrpcconverters.ConvertGRPCCreateValidIngredientStateIngredientRequestToValidIngredientStateIngredientCreationRequestInput(request))
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "creating valid ingredient state ingredient")
	}
	tracing.AttachToSpan(span, keys.ValidIngredientStateIngredientIDKey, created.ID)

	result := &mealplanning.CreateValidIngredientStateIngredientResponse{
		Result: mpgrpcconverters.ConvertValidIngredientStateIngredientToGRPCValidIngredientStateIngredient(created),
	}

	return result, nil
}

func (s *ServiceImpl) CreateValidInstrument(ctx context.Context, request *mealplanning.CreateValidInstrumentRequest) (*mealplanning.CreateValidInstrumentResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	created, err := s.validEnumerationsManager.CreateValidInstrument(ctx, mpgrpcconverters.ConvertGRPCCreateValidInstrumentRequestToValidInstrumentCreationRequestInput(request))
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "creating valid instrument")
	}
	tracing.AttachToSpan(span, keys.ValidInstrumentIDKey, created.ID)

	result := &mealplanning.CreateValidInstrumentResponse{
		Result: mpgrpcconverters.ConvertValidInstrumentToGRPCValidInstrument(created),
	}

	return result, nil
}

func (s *ServiceImpl) CreateValidMeasurementUnit(ctx context.Context, request *mealplanning.CreateValidMeasurementUnitRequest) (*mealplanning.CreateValidMeasurementUnitResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	created, err := s.validEnumerationsManager.CreateValidMeasurementUnit(ctx, mpgrpcconverters.ConvertGRPCCreateValidMeasurementUnitRequestToValidMeasurementUnitCreationRequestInput(request))
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "creating valid measurement unit")
	}
	tracing.AttachToSpan(span, keys.ValidMeasurementUnitIDKey, created.ID)

	result := &mealplanning.CreateValidMeasurementUnitResponse{
		Result: mpgrpcconverters.ConvertValidMeasurementUnitToGRPCValidMeasurementUnit(created),
	}

	return result, nil
}

func (s *ServiceImpl) CreateValidMeasurementUnitConversion(ctx context.Context, request *mealplanning.CreateValidMeasurementUnitConversionRequest) (*mealplanning.CreateValidMeasurementUnitConversionResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	created, err := s.validEnumerationsManager.CreateValidMeasurementUnitConversion(ctx, mpgrpcconverters.ConvertGRPCCreateValidMeasurementUnitConversionRequestToValidMeasurementUnitConversionCreationRequestInput(request))
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "creating valid measurement unit conversion")
	}
	tracing.AttachToSpan(span, keys.ValidMeasurementUnitConversionIDKey, created.ID)

	result := &mealplanning.CreateValidMeasurementUnitConversionResponse{
		Result: mpgrpcconverters.ConvertValidMeasurementUnitConversionToGRPCValidMeasurementUnitConversion(created),
	}

	return result, nil
}

func (s *ServiceImpl) CreateValidPreparation(ctx context.Context, request *mealplanning.CreateValidPreparationRequest) (*mealplanning.CreateValidPreparationResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	created, err := s.validEnumerationsManager.CreateValidPreparation(ctx, mpgrpcconverters.ConvertGRPCCreateValidPreparationRequestToValidPreparationCreationRequestInput(request))
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "creating valid preparation")
	}
	tracing.AttachToSpan(span, keys.ValidPreparationIDKey, created.ID)

	result := &mealplanning.CreateValidPreparationResponse{
		Result: mpgrpcconverters.ConvertValidPreparationToGRPCValidPreparation(created),
	}

	return result, nil
}

func (s *ServiceImpl) CreateValidPreparationInstrument(ctx context.Context, request *mealplanning.CreateValidPreparationInstrumentRequest) (*mealplanning.CreateValidPreparationInstrumentResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	created, err := s.validEnumerationsManager.CreateValidPreparationInstrument(ctx, mpgrpcconverters.ConvertGRPCCreateValidPreparationInstrumentRequestToValidPreparationInstrumentCreationRequestInput(request))
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "creating valid preparation instrument")
	}
	tracing.AttachToSpan(span, keys.ValidPreparationInstrumentIDKey, created.ID)

	result := &mealplanning.CreateValidPreparationInstrumentResponse{
		Result: mpgrpcconverters.ConvertValidPreparationInstrumentToGRPCValidPreparationInstrument(created),
	}

	return result, nil
}

func (s *ServiceImpl) CreateValidPreparationVessel(ctx context.Context, request *mealplanning.CreateValidPreparationVesselRequest) (*mealplanning.CreateValidPreparationVesselResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	created, err := s.validEnumerationsManager.CreateValidPreparationVessel(ctx, mpgrpcconverters.ConvertGRPCCreateValidPreparationVesselRequestToValidPreparationVesselCreationRequestInput(request))
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "creating valid preparation vessel")
	}
	tracing.AttachToSpan(span, keys.ValidPreparationVesselIDKey, created.ID)

	result := &mealplanning.CreateValidPreparationVesselResponse{
		Result: mpgrpcconverters.ConvertValidPreparationVesselToGRPCValidPreparationVessel(created),
	}

	return result, nil
}

func (s *ServiceImpl) CreateValidVessel(ctx context.Context, request *mealplanning.CreateValidVesselRequest) (*mealplanning.CreateValidVesselResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	created, err := s.validEnumerationsManager.CreateValidVessel(ctx, mpgrpcconverters.ConvertGRPCCreateValidVesselRequestToValidVesselCreationRequestInput(request))
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "creating valid vessel")
	}
	tracing.AttachToSpan(span, keys.ValidVesselIDKey, created.ID)

	result := &mealplanning.CreateValidVesselResponse{
		Result: mpgrpcconverters.ConvertValidVesselToGRPCValidVessel(created),
	}

	return result, nil
}

func (s *ServiceImpl) GetRandomValidIngredient(ctx context.Context, _ *mealplanning.GetRandomValidIngredientRequest) (*mealplanning.GetRandomValidIngredientResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	selected, err := s.validEnumerationsManager.RandomValidIngredient(ctx)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "fetching random valid ingredient")
	}
	tracing.AttachToSpan(span, keys.ValidIngredientIDKey, selected.ID)

	res := &mealplanning.GetRandomValidIngredientResponse{
		Result: mpgrpcconverters.ConvertValidIngredientToGRPCValidIngredient(selected),
	}

	return res, nil
}

func (s *ServiceImpl) GetRandomValidInstrument(ctx context.Context, _ *mealplanning.GetRandomValidInstrumentRequest) (*mealplanning.GetRandomValidInstrumentResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	selected, err := s.validEnumerationsManager.RandomValidInstrument(ctx)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "fetching random valid instrument")
	}
	tracing.AttachToSpan(span, keys.ValidInstrumentIDKey, selected.ID)

	res := &mealplanning.GetRandomValidInstrumentResponse{
		Result: mpgrpcconverters.ConvertValidInstrumentToGRPCValidInstrument(selected),
	}

	return res, nil
}

func (s *ServiceImpl) GetRandomValidPreparation(ctx context.Context, _ *mealplanning.GetRandomValidPreparationRequest) (*mealplanning.GetRandomValidPreparationResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	selected, err := s.validEnumerationsManager.RandomValidPreparation(ctx)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "fetching random valid preparation")
	}
	tracing.AttachToSpan(span, keys.ValidPreparationIDKey, selected.ID)

	res := &mealplanning.GetRandomValidPreparationResponse{
		Result: mpgrpcconverters.ConvertValidPreparationToGRPCValidPreparation(selected),
	}

	return res, nil
}

func (s *ServiceImpl) GetRandomValidVessel(ctx context.Context, _ *mealplanning.GetRandomValidVesselRequest) (*mealplanning.GetRandomValidVesselResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	selected, err := s.validEnumerationsManager.RandomValidVessel(ctx)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "fetching random valid vessel")
	}
	tracing.AttachToSpan(span, keys.ValidVesselIDKey, selected.ID)

	res := &mealplanning.GetRandomValidVesselResponse{
		Result: mpgrpcconverters.ConvertValidVesselToGRPCValidVessel(selected),
	}

	return res, nil
}

func (s *ServiceImpl) GetValidIngredient(ctx context.Context, request *mealplanning.GetValidIngredientRequest) (*mealplanning.GetValidIngredientResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValue(keys.ValidIngredientIDKey, request.ValidIngredientID)
	tracing.AttachToSpan(span, keys.ValidIngredientIDKey, request.ValidIngredientID)

	x, err := s.validEnumerationsManager.ReadValidIngredient(ctx, request.ValidIngredientID)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "fetching valid ingredient")
	}

	res := &mealplanning.GetValidIngredientResponse{
		Result: mpgrpcconverters.ConvertValidIngredientToGRPCValidIngredient(x),
	}

	return res, nil
}

func (s *ServiceImpl) GetValidIngredientGroup(ctx context.Context, request *mealplanning.GetValidIngredientGroupRequest) (*mealplanning.GetValidIngredientGroupResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValue(keys.ValidIngredientGroupIDKey, request.ValidIngredientGroupID)
	tracing.AttachToSpan(span, keys.ValidIngredientGroupIDKey, request.ValidIngredientGroupID)

	x, err := s.validEnumerationsManager.ReadValidIngredientGroup(ctx, request.ValidIngredientGroupID)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "fetching valid ingredient group")
	}

	res := &mealplanning.GetValidIngredientGroupResponse{
		Result: mpgrpcconverters.ConvertValidIngredientGroupToGRPCValidIngredientGroup(x),
	}

	return res, nil
}

func (s *ServiceImpl) GetValidIngredientGroups(ctx context.Context, request *mealplanning.GetValidIngredientGroupsRequest) (*mealplanning.GetValidIngredientGroupsResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)
	filter := grpcconverters.ConvertGRPCQueryFilterToQueryFilter(request.Filter)
	tracing.AttachQueryFilterToSpan(span, filter)

	x, _, err := s.validEnumerationsManager.ListValidIngredientGroups(ctx, filter)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "fetching valid ingredient groups")
	}

	res := &mealplanning.GetValidIngredientGroupsResponse{
		Filter: request.Filter,
	}
	for _, y := range x {
		res.Results = append(res.Results, mpgrpcconverters.ConvertValidIngredientGroupToGRPCValidIngredientGroup(y))
	}

	return res, nil
}

func (s *ServiceImpl) GetValidIngredientMeasurementUnit(ctx context.Context, request *mealplanning.GetValidIngredientMeasurementUnitRequest) (*mealplanning.GetValidIngredientMeasurementUnitResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValue(keys.ValidIngredientMeasurementUnitIDKey, request.ValidIngredientMeasurementUnitID)
	tracing.AttachToSpan(span, keys.ValidIngredientMeasurementUnitIDKey, request.ValidIngredientMeasurementUnitID)

	x, err := s.validEnumerationsManager.ReadValidIngredientMeasurementUnit(ctx, request.ValidIngredientMeasurementUnitID)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "fetching valid ingredient measurement unit")
	}

	res := &mealplanning.GetValidIngredientMeasurementUnitResponse{
		Result: mpgrpcconverters.ConvertValidIngredientMeasurementUnitToGRPCValidIngredientMeasurementUnit(x),
	}

	return res, nil
}

func (s *ServiceImpl) GetValidIngredientMeasurementUnits(ctx context.Context, request *mealplanning.GetValidIngredientMeasurementUnitsRequest) (*mealplanning.GetValidIngredientMeasurementUnitsResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)
	filter := grpcconverters.ConvertGRPCQueryFilterToQueryFilter(request.Filter)
	tracing.AttachQueryFilterToSpan(span, filter)

	x, _, err := s.validEnumerationsManager.ListValidIngredientMeasurementUnits(ctx, filter)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "fetching valid ingredient measurement units")
	}

	res := &mealplanning.GetValidIngredientMeasurementUnitsResponse{
		Filter: request.Filter,
	}
	for _, y := range x {
		res.Results = append(res.Results, mpgrpcconverters.ConvertValidIngredientMeasurementUnitToGRPCValidIngredientMeasurementUnit(y))
	}

	return res, nil
}

func (s *ServiceImpl) GetValidIngredientMeasurementUnitsByIngredient(ctx context.Context, request *mealplanning.GetValidIngredientMeasurementUnitsByIngredientRequest) (*mealplanning.GetValidIngredientMeasurementUnitsByIngredientResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValue(keys.ValidIngredientIDKey, request.ValidIngredientID)
	tracing.AttachToSpan(span, keys.ValidIngredientIDKey, request.ValidIngredientID)
	filter := grpcconverters.ConvertGRPCQueryFilterToQueryFilter(request.Filter)
	tracing.AttachQueryFilterToSpan(span, filter)

	x, err := s.validEnumerationsManager.SearchValidIngredientMeasurementUnitsByIngredient(ctx, request.ValidIngredientID, filter)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "fetching valid ingredient measurement units by ingredient")
	}

	res := &mealplanning.GetValidIngredientMeasurementUnitsByIngredientResponse{
		Filter: request.Filter,
	}

	for _, y := range x {
		res.Results = append(res.Results, mpgrpcconverters.ConvertValidIngredientMeasurementUnitToGRPCValidIngredientMeasurementUnit(y))
	}

	return res, nil
}

func (s *ServiceImpl) GetValidIngredientMeasurementUnitsByMeasurementUnit(ctx context.Context, request *mealplanning.GetValidIngredientMeasurementUnitsByMeasurementUnitRequest) (*mealplanning.GetValidIngredientMeasurementUnitsByMeasurementUnitResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)
	filter := grpcconverters.ConvertGRPCQueryFilterToQueryFilter(request.Filter)
	tracing.AttachQueryFilterToSpan(span, filter)

	x, err := s.validEnumerationsManager.SearchValidIngredientMeasurementUnitsByMeasurementUnit(ctx, request.ValidMeasurementUnitID, filter)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "fetching valid ingredient measurement units by measurement unit")
	}

	res := &mealplanning.GetValidIngredientMeasurementUnitsByMeasurementUnitResponse{
		Filter: request.Filter,
	}

	for _, y := range x {
		res.Results = append(res.Results, mpgrpcconverters.ConvertValidIngredientMeasurementUnitToGRPCValidIngredientMeasurementUnit(y))
	}

	return res, nil
}

func (s *ServiceImpl) GetValidIngredientPreparation(ctx context.Context, request *mealplanning.GetValidIngredientPreparationRequest) (*mealplanning.GetValidIngredientPreparationResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValue(keys.ValidIngredientPreparationIDKey, request.ValidIngredientPreparationID)
	tracing.AttachToSpan(span, keys.ValidIngredientPreparationIDKey, request.ValidIngredientPreparationID)

	x, err := s.validEnumerationsManager.ReadValidIngredientPreparation(ctx, request.ValidIngredientPreparationID)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "fetching valid ingredient preparation")
	}

	res := &mealplanning.GetValidIngredientPreparationResponse{
		Result: mpgrpcconverters.ConvertValidIngredientPreparationToGRPCValidIngredientPreparation(x),
	}

	return res, nil
}

func (s *ServiceImpl) GetValidIngredientPreparations(ctx context.Context, request *mealplanning.GetValidIngredientPreparationsRequest) (*mealplanning.GetValidIngredientPreparationsResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)
	filter := grpcconverters.ConvertGRPCQueryFilterToQueryFilter(request.Filter)
	tracing.AttachQueryFilterToSpan(span, filter)

	x, _, err := s.validEnumerationsManager.ListValidIngredientPreparations(ctx, filter)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "fetching valid ingredient preparations")
	}

	res := &mealplanning.GetValidIngredientPreparationsResponse{
		Filter: request.Filter,
	}

	for _, y := range x {
		res.Results = append(res.Results, mpgrpcconverters.ConvertValidIngredientPreparationToGRPCValidIngredientPreparation(y))
	}

	return res, nil
}

func (s *ServiceImpl) GetValidIngredientPreparationsByIngredient(ctx context.Context, request *mealplanning.GetValidIngredientPreparationsByIngredientRequest) (*mealplanning.GetValidIngredientPreparationsByIngredientResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValue(keys.ValidIngredientIDKey, request.ValidIngredientID)
	tracing.AttachToSpan(span, keys.ValidIngredientIDKey, request.ValidIngredientID)
	filter := grpcconverters.ConvertGRPCQueryFilterToQueryFilter(request.Filter)
	tracing.AttachQueryFilterToSpan(span, filter)

	x, err := s.validEnumerationsManager.SearchValidIngredientPreparationsByIngredient(ctx, request.ValidIngredientID, filter)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "fetching valid ingredient preparations by ingredient")
	}

	res := &mealplanning.GetValidIngredientPreparationsByIngredientResponse{
		Filter: request.Filter,
	}

	for _, y := range x {
		res.Results = append(res.Results, mpgrpcconverters.ConvertValidIngredientPreparationToGRPCValidIngredientPreparation(y))
	}

	return res, nil
}

func (s *ServiceImpl) GetValidIngredientPreparationsByPreparation(ctx context.Context, request *mealplanning.GetValidIngredientPreparationsByPreparationRequest) (*mealplanning.GetValidIngredientPreparationsByPreparationResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValue(keys.ValidPreparationIDKey, request.ValidPreparationID)
	tracing.AttachToSpan(span, keys.ValidPreparationIDKey, request.ValidPreparationID)
	filter := grpcconverters.ConvertGRPCQueryFilterToQueryFilter(request.Filter)
	tracing.AttachQueryFilterToSpan(span, filter)

	x, err := s.validEnumerationsManager.SearchValidIngredientPreparationsByPreparation(ctx, request.ValidPreparationID, filter)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "fetching valid ingredient preparations by preparation")
	}

	res := &mealplanning.GetValidIngredientPreparationsByPreparationResponse{
		Filter: request.Filter,
	}

	for _, y := range x {
		res.Results = append(res.Results, mpgrpcconverters.ConvertValidIngredientPreparationToGRPCValidIngredientPreparation(y))
	}

	return res, nil
}

func (s *ServiceImpl) GetValidIngredientState(ctx context.Context, request *mealplanning.GetValidIngredientStateRequest) (*mealplanning.GetValidIngredientStateResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValue(keys.ValidIngredientStateIDKey, request.ValidIngredientStateID)
	tracing.AttachToSpan(span, keys.ValidIngredientStateIDKey, request.ValidIngredientStateID)

	x, err := s.validEnumerationsManager.ReadValidIngredientState(ctx, request.ValidIngredientStateID)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "fetching valid ingredient state")
	}

	res := &mealplanning.GetValidIngredientStateResponse{
		Result: mpgrpcconverters.ConvertValidIngredientStateToGRPCValidIngredientState(x),
	}

	return res, nil
}

func (s *ServiceImpl) GetValidIngredientStateIngredient(ctx context.Context, request *mealplanning.GetValidIngredientStateIngredientRequest) (*mealplanning.GetValidIngredientStateIngredientResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValue(keys.ValidIngredientStateIngredientIDKey, request.ValidIngredientStateIngredientID)
	tracing.AttachToSpan(span, keys.ValidIngredientStateIngredientIDKey, request.ValidIngredientStateIngredientID)

	x, err := s.validEnumerationsManager.ReadValidIngredientStateIngredient(ctx, request.ValidIngredientStateIngredientID)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "fetching valid ingredient state ingredient")
	}

	res := &mealplanning.GetValidIngredientStateIngredientResponse{
		Result: mpgrpcconverters.ConvertValidIngredientStateIngredientToGRPCValidIngredientStateIngredient(x),
	}

	return res, nil
}

func (s *ServiceImpl) GetValidIngredientStateIngredients(ctx context.Context, request *mealplanning.GetValidIngredientStateIngredientsRequest) (*mealplanning.GetValidIngredientStateIngredientsResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)
	filter := grpcconverters.ConvertGRPCQueryFilterToQueryFilter(request.Filter)
	tracing.AttachQueryFilterToSpan(span, filter)

	x, _, err := s.validEnumerationsManager.ListValidIngredientStateIngredients(ctx, filter)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "fetching valid ingredient state ingredients")
	}

	res := &mealplanning.GetValidIngredientStateIngredientsResponse{
		Filter: request.Filter,
	}

	for _, y := range x {
		res.Results = append(res.Results, mpgrpcconverters.ConvertValidIngredientStateIngredientToGRPCValidIngredientStateIngredient(y))
	}

	return res, nil
}

func (s *ServiceImpl) GetValidIngredientStateIngredientsByIngredient(ctx context.Context, request *mealplanning.GetValidIngredientStateIngredientsByIngredientRequest) (*mealplanning.GetValidIngredientStateIngredientsByIngredientResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValue(keys.ValidIngredientIDKey, request.ValidIngredientID)
	tracing.AttachToSpan(span, keys.ValidIngredientIDKey, request.ValidIngredientID)
	filter := grpcconverters.ConvertGRPCQueryFilterToQueryFilter(request.Filter)
	tracing.AttachQueryFilterToSpan(span, filter)

	x, err := s.validEnumerationsManager.SearchValidIngredientStateIngredientsByIngredient(ctx, request.ValidIngredientID, filter)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "fetching valid ingredient state ingredients by ingredient")
	}

	res := &mealplanning.GetValidIngredientStateIngredientsByIngredientResponse{
		Filter: request.Filter,
	}

	for _, y := range x {
		res.Results = append(res.Results, mpgrpcconverters.ConvertValidIngredientStateIngredientToGRPCValidIngredientStateIngredient(y))
	}

	return res, nil
}

func (s *ServiceImpl) GetValidIngredientStateIngredientsByIngredientState(ctx context.Context, request *mealplanning.GetValidIngredientStateIngredientsByIngredientStateRequest) (*mealplanning.GetValidIngredientStateIngredientsByIngredientStateResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValue(keys.ValidIngredientStateIDKey, request.ValidIngredientStateID)
	tracing.AttachToSpan(span, keys.ValidIngredientStateIDKey, request.ValidIngredientStateID)
	filter := grpcconverters.ConvertGRPCQueryFilterToQueryFilter(request.Filter)
	tracing.AttachQueryFilterToSpan(span, filter)

	x, err := s.validEnumerationsManager.SearchValidIngredientStateIngredientsByIngredientState(ctx, request.ValidIngredientStateID, filter)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "fetching valid ingredient state ingredients by ingredient state")
	}

	res := &mealplanning.GetValidIngredientStateIngredientsByIngredientStateResponse{
		Filter: request.Filter,
	}

	for _, y := range x {
		res.Results = append(res.Results, mpgrpcconverters.ConvertValidIngredientStateIngredientToGRPCValidIngredientStateIngredient(y))
	}

	return res, nil
}

func (s *ServiceImpl) GetValidIngredientStates(ctx context.Context, request *mealplanning.GetValidIngredientStatesRequest) (*mealplanning.GetValidIngredientStatesResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)
	filter := grpcconverters.ConvertGRPCQueryFilterToQueryFilter(request.Filter)
	tracing.AttachQueryFilterToSpan(span, filter)

	x, _, err := s.validEnumerationsManager.ListValidIngredientStates(ctx, filter)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "fetching valid ingredient states")
	}

	res := &mealplanning.GetValidIngredientStatesResponse{
		Filter: request.Filter,
	}

	for _, y := range x {
		res.Results = append(res.Results, mpgrpcconverters.ConvertValidIngredientStateToGRPCValidIngredientState(y))
	}

	return res, nil
}

func (s *ServiceImpl) GetValidIngredients(ctx context.Context, request *mealplanning.GetValidIngredientsRequest) (*mealplanning.GetValidIngredientsResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)
	filter := grpcconverters.ConvertGRPCQueryFilterToQueryFilter(request.Filter)
	tracing.AttachQueryFilterToSpan(span, filter)

	x, _, err := s.validEnumerationsManager.ListValidIngredients(ctx, filter)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "fetching valid ingredients")
	}

	res := &mealplanning.GetValidIngredientsResponse{
		Filter: request.Filter,
	}

	for _, y := range x {
		res.Results = append(res.Results, mpgrpcconverters.ConvertValidIngredientToGRPCValidIngredient(y))
	}

	return res, nil
}

func (s *ServiceImpl) GetValidInstrument(ctx context.Context, request *mealplanning.GetValidInstrumentRequest) (*mealplanning.GetValidInstrumentResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValue(keys.ValidInstrumentIDKey, request.ValidInstrumentID)
	tracing.AttachToSpan(span, keys.ValidInstrumentIDKey, request.ValidInstrumentID)

	x, err := s.validEnumerationsManager.ReadValidInstrument(ctx, request.ValidInstrumentID)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "fetching valid instrument")
	}

	res := &mealplanning.GetValidInstrumentResponse{
		Result: mpgrpcconverters.ConvertValidInstrumentToGRPCValidInstrument(x),
	}

	return res, nil
}

func (s *ServiceImpl) GetValidInstruments(ctx context.Context, request *mealplanning.GetValidInstrumentsRequest) (*mealplanning.GetValidInstrumentsResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)
	filter := grpcconverters.ConvertGRPCQueryFilterToQueryFilter(request.Filter)
	tracing.AttachQueryFilterToSpan(span, filter)

	x, _, err := s.validEnumerationsManager.ListValidInstruments(ctx, filter)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "fetching valid instruments")
	}

	res := &mealplanning.GetValidInstrumentsResponse{
		Filter: request.Filter,
	}

	for _, y := range x {
		res.Results = append(res.Results, mpgrpcconverters.ConvertValidInstrumentToGRPCValidInstrument(y))
	}

	return res, nil
}

func (s *ServiceImpl) GetValidMeasurementUnit(ctx context.Context, request *mealplanning.GetValidMeasurementUnitRequest) (*mealplanning.GetValidMeasurementUnitResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValue(keys.ValidMeasurementUnitIDKey, request.ValidMeasurementUnitID)
	tracing.AttachToSpan(span, keys.ValidMeasurementUnitIDKey, request.ValidMeasurementUnitID)

	x, err := s.validEnumerationsManager.ReadValidMeasurementUnit(ctx, request.ValidMeasurementUnitID)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "fetching valid measurement unit")
	}

	res := &mealplanning.GetValidMeasurementUnitResponse{
		Result: mpgrpcconverters.ConvertValidMeasurementUnitToGRPCValidMeasurementUnit(x),
	}

	return res, nil
}

func (s *ServiceImpl) GetValidMeasurementUnitConversion(ctx context.Context, request *mealplanning.GetValidMeasurementUnitConversionRequest) (*mealplanning.GetValidMeasurementUnitConversionResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValue(keys.ValidMeasurementUnitConversionIDKey, request.ValidMeasurementUnitConversionID)
	tracing.AttachToSpan(span, keys.ValidMeasurementUnitConversionIDKey, request.ValidMeasurementUnitConversionID)

	x, err := s.validEnumerationsManager.ReadValidMeasurementUnitConversion(ctx, request.ValidMeasurementUnitConversionID)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "fetching valid measurement unit conversion")
	}

	res := &mealplanning.GetValidMeasurementUnitConversionResponse{
		Result: mpgrpcconverters.ConvertValidMeasurementUnitConversionToGRPCValidMeasurementUnitConversion(x),
	}

	return res, nil
}

func (s *ServiceImpl) GetValidMeasurementUnitConversionsFromUnit(ctx context.Context, request *mealplanning.GetValidMeasurementUnitConversionsFromUnitRequest) (*mealplanning.GetValidMeasurementUnitConversionsFromUnitResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValue(keys.ValidMeasurementUnitIDKey, request.ValidMeasurementUnitID)
	tracing.AttachToSpan(span, keys.ValidMeasurementUnitIDKey, request.ValidMeasurementUnitID)

	x, err := s.validEnumerationsManager.ValidMeasurementUnitConversionsFromMeasurementUnit(ctx, request.ValidMeasurementUnitID)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "fetching valid measurement unit conversions from unit")
	}

	res := &mealplanning.GetValidMeasurementUnitConversionsFromUnitResponse{}

	for _, y := range x {
		res.Results = append(res.Results, mpgrpcconverters.ConvertValidMeasurementUnitConversionToGRPCValidMeasurementUnitConversion(y))
	}

	return res, nil
}

func (s *ServiceImpl) GetValidMeasurementUnitConversionsToUnit(ctx context.Context, request *mealplanning.GetValidMeasurementUnitConversionsToUnitRequest) (*mealplanning.GetValidMeasurementUnitConversionsToUnitResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValue(keys.ValidMeasurementUnitIDKey, request.ValidMeasurementUnitID)
	tracing.AttachToSpan(span, keys.ValidMeasurementUnitIDKey, request.ValidMeasurementUnitID)

	x, err := s.validEnumerationsManager.ValidMeasurementUnitConversionsToMeasurementUnit(ctx, request.ValidMeasurementUnitID)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "fetching valid measurement unit conversions to unit")
	}

	res := &mealplanning.GetValidMeasurementUnitConversionsToUnitResponse{}
	for _, y := range x {
		res.Results = append(res.Results, mpgrpcconverters.ConvertValidMeasurementUnitConversionToGRPCValidMeasurementUnitConversion(y))
	}

	return res, nil
}

func (s *ServiceImpl) GetValidMeasurementUnits(ctx context.Context, request *mealplanning.GetValidMeasurementUnitsRequest) (*mealplanning.GetValidMeasurementUnitsResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)
	filter := grpcconverters.ConvertGRPCQueryFilterToQueryFilter(request.Filter)
	tracing.AttachQueryFilterToSpan(span, filter)

	x, _, err := s.validEnumerationsManager.ListValidMeasurementUnits(ctx, filter)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "fetching valid measurement units")
	}

	res := &mealplanning.GetValidMeasurementUnitsResponse{
		Filter: request.Filter,
	}

	for _, y := range x {
		res.Results = append(res.Results, mpgrpcconverters.ConvertValidMeasurementUnitToGRPCValidMeasurementUnit(y))
	}

	return res, nil
}

func (s *ServiceImpl) GetValidPreparation(ctx context.Context, request *mealplanning.GetValidPreparationRequest) (*mealplanning.GetValidPreparationResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValue(keys.ValidPreparationIDKey, request.ValidPreparationID)
	tracing.AttachToSpan(span, keys.ValidPreparationIDKey, request.ValidPreparationID)

	x, err := s.validEnumerationsManager.ReadValidPreparation(ctx, request.ValidPreparationID)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "fetching valid preparation")
	}

	res := &mealplanning.GetValidPreparationResponse{
		Result: mpgrpcconverters.ConvertValidPreparationToGRPCValidPreparation(x),
	}

	return res, nil
}

func (s *ServiceImpl) GetValidPreparationInstrument(ctx context.Context, request *mealplanning.GetValidPreparationInstrumentRequest) (*mealplanning.GetValidPreparationInstrumentResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValue(keys.ValidPreparationInstrumentIDKey, request.ValidPreparationInstrumentID)
	tracing.AttachToSpan(span, keys.ValidPreparationInstrumentIDKey, request.ValidPreparationInstrumentID)

	x, err := s.validEnumerationsManager.ReadValidPreparationInstrument(ctx, request.ValidPreparationInstrumentID)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "fetching valid preparation instrument")
	}

	res := &mealplanning.GetValidPreparationInstrumentResponse{
		Result: mpgrpcconverters.ConvertValidPreparationInstrumentToGRPCValidPreparationInstrument(x),
	}

	return res, nil
}

func (s *ServiceImpl) GetValidPreparationInstruments(ctx context.Context, request *mealplanning.GetValidPreparationInstrumentsRequest) (*mealplanning.GetValidPreparationInstrumentsResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)
	filter := grpcconverters.ConvertGRPCQueryFilterToQueryFilter(request.Filter)
	tracing.AttachQueryFilterToSpan(span, filter)

	x, _, err := s.validEnumerationsManager.ListValidPreparationInstruments(ctx, filter)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "fetching valid preparation instruments")
	}

	res := &mealplanning.GetValidPreparationInstrumentsResponse{
		Filter: request.Filter,
	}

	for _, y := range x {
		res.Results = append(res.Results, mpgrpcconverters.ConvertValidPreparationInstrumentToGRPCValidPreparationInstrument(y))
	}

	return res, nil
}

func (s *ServiceImpl) GetValidPreparationInstrumentsByInstrument(ctx context.Context, request *mealplanning.GetValidPreparationInstrumentsByInstrumentRequest) (*mealplanning.GetValidPreparationInstrumentsByInstrumentResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValue(keys.ValidInstrumentIDKey, request.ValidInstrumentID)
	tracing.AttachToSpan(span, keys.ValidInstrumentIDKey, request.ValidInstrumentID)
	filter := grpcconverters.ConvertGRPCQueryFilterToQueryFilter(request.Filter)
	tracing.AttachQueryFilterToSpan(span, filter)

	x, err := s.validEnumerationsManager.SearchValidPreparationInstrumentsByInstrument(ctx, request.ValidInstrumentID, filter)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "fetching valid preparation instruments by instrument")
	}

	res := &mealplanning.GetValidPreparationInstrumentsByInstrumentResponse{
		Filter: request.Filter,
	}

	for _, y := range x {
		res.Results = append(res.Results, mpgrpcconverters.ConvertValidPreparationInstrumentToGRPCValidPreparationInstrument(y))
	}

	return res, nil
}

func (s *ServiceImpl) GetValidPreparationInstrumentsByPreparation(ctx context.Context, request *mealplanning.GetValidPreparationInstrumentsByPreparationRequest) (*mealplanning.GetValidPreparationInstrumentsByPreparationResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValue(keys.ValidPreparationIDKey, request.ValidPreparationID)
	tracing.AttachToSpan(span, keys.ValidPreparationIDKey, request.ValidPreparationID)
	filter := grpcconverters.ConvertGRPCQueryFilterToQueryFilter(request.Filter)
	tracing.AttachQueryFilterToSpan(span, filter)

	x, err := s.validEnumerationsManager.SearchValidPreparationInstrumentsByPreparation(ctx, request.ValidPreparationID, filter)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "fetching valid preparation instruments by preparation")
	}

	res := &mealplanning.GetValidPreparationInstrumentsByPreparationResponse{
		Filter: request.Filter,
	}

	for _, y := range x {
		res.Results = append(res.Results, mpgrpcconverters.ConvertValidPreparationInstrumentToGRPCValidPreparationInstrument(y))
	}

	return res, nil
}

func (s *ServiceImpl) GetValidPreparationVessel(ctx context.Context, request *mealplanning.GetValidPreparationVesselRequest) (*mealplanning.GetValidPreparationVesselResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValue(keys.ValidPreparationVesselIDKey, request.ValidPreparationVesselID)
	tracing.AttachToSpan(span, keys.ValidPreparationVesselIDKey, request.ValidPreparationVesselID)

	x, err := s.validEnumerationsManager.ReadValidPreparationVessel(ctx, request.ValidPreparationVesselID)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "fetching valid preparation vessel")
	}

	res := &mealplanning.GetValidPreparationVesselResponse{
		Result: mpgrpcconverters.ConvertValidPreparationVesselToGRPCValidPreparationVessel(x),
	}

	return res, nil
}

func (s *ServiceImpl) GetValidPreparationVessels(ctx context.Context, request *mealplanning.GetValidPreparationVesselsRequest) (*mealplanning.GetValidPreparationVesselsResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)
	filter := grpcconverters.ConvertGRPCQueryFilterToQueryFilter(request.Filter)
	tracing.AttachQueryFilterToSpan(span, filter)

	x, _, err := s.validEnumerationsManager.ListValidPreparationVessels(ctx, filter)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "fetching valid preparation vessels")
	}

	res := &mealplanning.GetValidPreparationVesselsResponse{
		Filter: request.Filter,
	}

	for _, y := range x {
		res.Results = append(res.Results, mpgrpcconverters.ConvertValidPreparationVesselToGRPCValidPreparationVessel(y))
	}

	return res, nil
}

func (s *ServiceImpl) GetValidPreparationVesselsByPreparation(ctx context.Context, request *mealplanning.GetValidPreparationVesselsByPreparationRequest) (*mealplanning.GetValidPreparationVesselsByPreparationResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValue(keys.ValidPreparationIDKey, request.ValidPreparationID)
	tracing.AttachToSpan(span, keys.ValidPreparationIDKey, request.ValidPreparationID)
	filter := grpcconverters.ConvertGRPCQueryFilterToQueryFilter(request.Filter)
	tracing.AttachQueryFilterToSpan(span, filter)

	x, err := s.validEnumerationsManager.SearchValidPreparationVesselsByPreparation(ctx, request.ValidPreparationID, filter)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "fetching valid preparation vessels by preparation")
	}

	res := &mealplanning.GetValidPreparationVesselsByPreparationResponse{
		Filter: request.Filter,
	}

	for _, y := range x {
		res.Results = append(res.Results, mpgrpcconverters.ConvertValidPreparationVesselToGRPCValidPreparationVessel(y))
	}

	return res, nil
}

func (s *ServiceImpl) GetValidPreparationVesselsByVessel(ctx context.Context, request *mealplanning.GetValidPreparationVesselsByVesselRequest) (*mealplanning.GetValidPreparationVesselsByVesselResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValue(keys.ValidVesselIDKey, request.ValidVesselID)
	tracing.AttachToSpan(span, keys.ValidVesselIDKey, request.ValidVesselID)
	filter := grpcconverters.ConvertGRPCQueryFilterToQueryFilter(request.Filter)
	tracing.AttachQueryFilterToSpan(span, filter)

	x, err := s.validEnumerationsManager.SearchValidPreparationVesselsByVessel(ctx, request.ValidVesselID, filter)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "fetching valid preparation vessels by vessel")
	}

	res := &mealplanning.GetValidPreparationVesselsByVesselResponse{
		Filter: request.Filter,
	}

	for _, y := range x {
		res.Results = append(res.Results, mpgrpcconverters.ConvertValidPreparationVesselToGRPCValidPreparationVessel(y))
	}

	return res, nil
}

func (s *ServiceImpl) GetValidPreparations(ctx context.Context, request *mealplanning.GetValidPreparationsRequest) (*mealplanning.GetValidPreparationsResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)
	filter := grpcconverters.ConvertGRPCQueryFilterToQueryFilter(request.Filter)
	tracing.AttachQueryFilterToSpan(span, filter)

	x, _, err := s.validEnumerationsManager.ListValidPreparations(ctx, filter)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "fetching valid preparations")
	}

	res := &mealplanning.GetValidPreparationsResponse{
		Filter: request.Filter,
	}

	for _, y := range x {
		res.Results = append(res.Results, mpgrpcconverters.ConvertValidPreparationToGRPCValidPreparation(y))
	}

	return res, nil
}

func (s *ServiceImpl) GetValidVessel(ctx context.Context, request *mealplanning.GetValidVesselRequest) (*mealplanning.GetValidVesselResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValue(keys.ValidVesselIDKey, request.ValidVesselID)
	tracing.AttachToSpan(span, keys.ValidVesselIDKey, request.ValidVesselID)

	x, err := s.validEnumerationsManager.ReadValidVessel(ctx, request.ValidVesselID)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "fetching valid vessel")
	}

	res := &mealplanning.GetValidVesselResponse{
		Result: mpgrpcconverters.ConvertValidVesselToGRPCValidVessel(x),
	}

	return res, nil
}

func (s *ServiceImpl) GetValidVessels(ctx context.Context, request *mealplanning.GetValidVesselsRequest) (*mealplanning.GetValidVesselsResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)
	filter := grpcconverters.ConvertGRPCQueryFilterToQueryFilter(request.Filter)
	tracing.AttachQueryFilterToSpan(span, filter)

	x, _, err := s.validEnumerationsManager.ListValidVessels(ctx, filter)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "fetching valid vessels")
	}

	res := &mealplanning.GetValidVesselsResponse{
		Filter: request.Filter,
	}

	for _, y := range x {
		res.Results = append(res.Results, mpgrpcconverters.ConvertValidVesselToGRPCValidVessel(y))
	}

	return res, nil
}

func (s *ServiceImpl) SearchForValidIngredientGroups(ctx context.Context, request *mealplanning.SearchForValidIngredientGroupsRequest) (*mealplanning.SearchForValidIngredientGroupsResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)
	filter := grpcconverters.ConvertGRPCQueryFilterToQueryFilter(request.Filter)
	tracing.AttachQueryFilterToSpan(span, filter)

	x, err := s.validEnumerationsManager.SearchValidIngredientGroups(ctx, request.Query, request.UseDatabase, filter)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "searching for valid ingredient groups")
	}

	res := &mealplanning.SearchForValidIngredientGroupsResponse{
		Filter: request.Filter,
	}

	for _, y := range x {
		res.Results = append(res.Results, mpgrpcconverters.ConvertValidIngredientGroupToGRPCValidIngredientGroup(y))
	}

	return res, nil
}

func (s *ServiceImpl) SearchForValidIngredientStates(ctx context.Context, request *mealplanning.SearchForValidIngredientStatesRequest) (*mealplanning.SearchForValidIngredientStatesResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)
	filter := grpcconverters.ConvertGRPCQueryFilterToQueryFilter(request.Filter)
	tracing.AttachQueryFilterToSpan(span, filter)

	x, err := s.validEnumerationsManager.SearchValidIngredientStates(ctx, request.Query, request.UseDatabase, filter)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "searching for valid ingredient states")
	}

	res := &mealplanning.SearchForValidIngredientStatesResponse{
		Filter: request.Filter,
	}

	for _, y := range x {
		res.Results = append(res.Results, mpgrpcconverters.ConvertValidIngredientStateToGRPCValidIngredientState(y))
	}

	return res, nil
}

func (s *ServiceImpl) SearchForValidIngredients(ctx context.Context, request *mealplanning.SearchForValidIngredientsRequest) (*mealplanning.SearchForValidIngredientsResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)
	filter := grpcconverters.ConvertGRPCQueryFilterToQueryFilter(request.Filter)
	tracing.AttachQueryFilterToSpan(span, filter)

	x, err := s.validEnumerationsManager.SearchValidIngredients(ctx, request.Query, request.UseDatabase, filter)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "searching for valid ingredients")
	}

	res := &mealplanning.SearchForValidIngredientsResponse{
		Filter: request.Filter,
	}

	for _, y := range x {
		res.Results = append(res.Results, mpgrpcconverters.ConvertValidIngredientToGRPCValidIngredient(y))
	}

	return res, nil
}

func (s *ServiceImpl) SearchForValidInstruments(ctx context.Context, request *mealplanning.SearchForValidInstrumentsRequest) (*mealplanning.SearchForValidInstrumentsResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)
	filter := grpcconverters.ConvertGRPCQueryFilterToQueryFilter(request.Filter)
	tracing.AttachQueryFilterToSpan(span, filter)

	x, err := s.validEnumerationsManager.SearchValidInstruments(ctx, request.Query, request.UseDatabase, filter)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "searching for valid instruments")
	}

	res := &mealplanning.SearchForValidInstrumentsResponse{
		Filter: request.Filter,
	}

	for _, y := range x {
		res.Results = append(res.Results, mpgrpcconverters.ConvertValidInstrumentToGRPCValidInstrument(y))
	}

	return res, nil
}

func (s *ServiceImpl) SearchForValidMeasurementUnits(ctx context.Context, request *mealplanning.SearchForValidMeasurementUnitsRequest) (*mealplanning.SearchForValidMeasurementUnitsResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)
	filter := grpcconverters.ConvertGRPCQueryFilterToQueryFilter(request.Filter)
	tracing.AttachQueryFilterToSpan(span, filter)

	x, err := s.validEnumerationsManager.SearchValidMeasurementUnits(ctx, request.Query, request.UseDatabase, filter)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "searching for valid measurement units")
	}

	res := &mealplanning.SearchForValidMeasurementUnitsResponse{
		Filter: request.Filter,
	}

	for _, y := range x {
		res.Results = append(res.Results, mpgrpcconverters.ConvertValidMeasurementUnitToGRPCValidMeasurementUnit(y))
	}

	return res, nil
}

func (s *ServiceImpl) SearchForValidPreparations(ctx context.Context, request *mealplanning.SearchForValidPreparationsRequest) (*mealplanning.SearchForValidPreparationsResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)
	filter := grpcconverters.ConvertGRPCQueryFilterToQueryFilter(request.Filter)
	tracing.AttachQueryFilterToSpan(span, filter)

	x, err := s.validEnumerationsManager.SearchValidPreparations(ctx, request.Query, request.UseDatabase, filter)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "searching for valid preparations")
	}

	res := &mealplanning.SearchForValidPreparationsResponse{
		Filter: request.Filter,
	}

	for _, y := range x {
		res.Results = append(res.Results, mpgrpcconverters.ConvertValidPreparationToGRPCValidPreparation(y))
	}

	return res, nil
}

func (s *ServiceImpl) SearchForValidVessels(ctx context.Context, request *mealplanning.SearchForValidVesselsRequest) (*mealplanning.SearchForValidVesselsResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)
	filter := grpcconverters.ConvertGRPCQueryFilterToQueryFilter(request.Filter)
	tracing.AttachQueryFilterToSpan(span, filter)

	x, err := s.validEnumerationsManager.SearchValidVessels(ctx, request.Query, request.UseDatabase, filter)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "searching for valid vessels")
	}

	res := &mealplanning.SearchForValidVesselsResponse{
		Filter: request.Filter,
	}

	for _, y := range x {
		res.Results = append(res.Results, mpgrpcconverters.ConvertValidVesselToGRPCValidVessel(y))
	}

	return res, nil
}

func (s *ServiceImpl) SearchValidIngredientsByPreparation(ctx context.Context, request *mealplanning.SearchValidIngredientsByPreparationRequest) (*mealplanning.SearchValidIngredientsByPreparationResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValue(keys.ValidPreparationIDKey, request.ValidPreparationID)
	tracing.AttachToSpan(span, keys.ValidPreparationIDKey, request.ValidPreparationID)
	filter := grpcconverters.ConvertGRPCQueryFilterToQueryFilter(request.Filter)
	tracing.AttachQueryFilterToSpan(span, filter)

	x, err := s.validEnumerationsManager.SearchValidIngredientsByPreparationAndIngredientName(ctx, request.ValidPreparationID, request.Query, filter)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "searching for valid ingredients by preparation")
	}

	res := &mealplanning.SearchValidIngredientsByPreparationResponse{
		Filter: request.Filter,
	}

	for _, y := range x {
		res.Results = append(res.Results, mpgrpcconverters.ConvertValidIngredientToGRPCValidIngredient(y))
	}

	return res, nil
}

func (s *ServiceImpl) SearchValidMeasurementUnitsByIngredient(ctx context.Context, request *mealplanning.SearchValidMeasurementUnitsByIngredientRequest) (*mealplanning.SearchValidMeasurementUnitsByIngredientResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValue(keys.ValidIngredientIDKey, request.ValidIngredientID)
	tracing.AttachToSpan(span, keys.ValidIngredientIDKey, request.ValidIngredientID)
	filter := grpcconverters.ConvertGRPCQueryFilterToQueryFilter(request.Filter)
	tracing.AttachQueryFilterToSpan(span, filter)

	x, err := s.validEnumerationsManager.SearchValidMeasurementUnitsByIngredientID(ctx, request.ValidIngredientID, filter)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "searching for valid measurement units by ingredient")
	}

	res := &mealplanning.SearchValidMeasurementUnitsByIngredientResponse{
		Filter: request.Filter,
	}

	for _, y := range x {
		res.Results = append(res.Results, mpgrpcconverters.ConvertValidMeasurementUnitToGRPCValidMeasurementUnit(y))
	}

	return res, nil
}

func (s *ServiceImpl) UpdateValidIngredient(ctx context.Context, request *mealplanning.UpdateValidIngredientRequest) (*mealplanning.UpdateValidIngredientResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValue(keys.ValidIngredientIDKey, request.ValidIngredientID)
	tracing.AttachToSpan(span, keys.ValidIngredientIDKey, request.ValidIngredientID)

	input := mpgrpcconverters.ConvertGRPCValidIngredientUpdateRequestInputToValidIngredientUpdateRequestInput(request.Input)
	updated, err := s.validEnumerationsManager.UpdateValidIngredient(ctx, request.ValidIngredientID, input)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "updating valid ingredient")
	}

	res := &mealplanning.UpdateValidIngredientResponse{
		Result: mpgrpcconverters.ConvertValidIngredientToGRPCValidIngredient(updated),
	}

	return res, nil
}

func (s *ServiceImpl) UpdateValidIngredientGroup(ctx context.Context, request *mealplanning.UpdateValidIngredientGroupRequest) (*mealplanning.UpdateValidIngredientGroupResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValue(keys.ValidIngredientGroupIDKey, request.ValidIngredientGroupID)
	tracing.AttachToSpan(span, keys.ValidIngredientGroupIDKey, request.ValidIngredientGroupID)

	input := mpgrpcconverters.ConvertGRPCValidIngredientGroupUpdateRequestInputToValidIngredientGroupUpdateRequestInput(request.Input)

	updated, err := s.validEnumerationsManager.UpdateValidIngredientGroup(ctx, request.ValidIngredientGroupID, input)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "updating valid ingredient group")
	}

	res := &mealplanning.UpdateValidIngredientGroupResponse{
		Result: mpgrpcconverters.ConvertValidIngredientGroupToGRPCValidIngredientGroup(updated),
	}

	return res, nil
}

func (s *ServiceImpl) UpdateValidIngredientMeasurementUnit(ctx context.Context, request *mealplanning.UpdateValidIngredientMeasurementUnitRequest) (*mealplanning.UpdateValidIngredientMeasurementUnitResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValue(keys.ValidIngredientMeasurementUnitIDKey, request.ValidIngredientMeasurementUnitID)
	tracing.AttachToSpan(span, keys.ValidIngredientMeasurementUnitIDKey, request.ValidIngredientMeasurementUnitID)

	input := mpgrpcconverters.ConvertGRPCValidIngredientMeasurementUnitUpdateRequestInputToValidIngredientMeasurementUnitUpdateRequestInput(request.Input)

	updated, err := s.validEnumerationsManager.UpdateValidIngredientMeasurementUnit(ctx, request.ValidIngredientMeasurementUnitID, input)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "updating valid ingredient measurement unit")
	}

	res := &mealplanning.UpdateValidIngredientMeasurementUnitResponse{
		Result: mpgrpcconverters.ConvertValidIngredientMeasurementUnitToGRPCValidIngredientMeasurementUnit(updated),
	}

	return res, nil
}

func (s *ServiceImpl) UpdateValidIngredientPreparation(ctx context.Context, request *mealplanning.UpdateValidIngredientPreparationRequest) (*mealplanning.UpdateValidIngredientPreparationResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValue(keys.ValidIngredientPreparationIDKey, request.ValidIngredientPreparationID)
	tracing.AttachToSpan(span, keys.ValidIngredientPreparationIDKey, request.ValidIngredientPreparationID)

	input := mpgrpcconverters.ConvertGRPCValidIngredientPreparationUpdateRequestInputToValidIngredientPreparationUpdateRequestInput(request.Input)

	updated, err := s.validEnumerationsManager.UpdateValidIngredientPreparation(ctx, request.ValidIngredientPreparationID, input)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "updating valid ingredient preparation")
	}

	res := &mealplanning.UpdateValidIngredientPreparationResponse{
		Result: mpgrpcconverters.ConvertValidIngredientPreparationToGRPCValidIngredientPreparation(updated),
	}

	return res, nil
}

func (s *ServiceImpl) UpdateValidIngredientState(ctx context.Context, request *mealplanning.UpdateValidIngredientStateRequest) (*mealplanning.UpdateValidIngredientStateResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValue(keys.ValidIngredientStateIDKey, request.ValidIngredientStateID)
	tracing.AttachToSpan(span, keys.ValidIngredientStateIDKey, request.ValidIngredientStateID)

	input := mpgrpcconverters.ConvertGRPCValidIngredientStateUpdateRequestInputToValidIngredientStateUpdateRequestInput(request.Input)

	updated, err := s.validEnumerationsManager.UpdateValidIngredientState(ctx, request.ValidIngredientStateID, input)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "updating valid ingredient state")
	}

	res := &mealplanning.UpdateValidIngredientStateResponse{
		Result: mpgrpcconverters.ConvertValidIngredientStateToGRPCValidIngredientState(updated),
	}

	return res, nil
}

func (s *ServiceImpl) UpdateValidIngredientStateIngredient(ctx context.Context, request *mealplanning.UpdateValidIngredientStateIngredientRequest) (*mealplanning.UpdateValidIngredientStateIngredientResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValue(keys.ValidIngredientStateIngredientIDKey, request.ValidIngredientStateIngredientID)
	tracing.AttachToSpan(span, keys.ValidIngredientStateIngredientIDKey, request.ValidIngredientStateIngredientID)

	input := mpgrpcconverters.ConvertGRPCValidIngredientStateIngredientUpdateRequestInputToValidIngredientStateIngredientUpdateRequestInput(request.Input)

	updated, err := s.validEnumerationsManager.UpdateValidIngredientStateIngredient(ctx, request.ValidIngredientStateIngredientID, input)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "updating valid ingredient state ingredient")
	}

	res := &mealplanning.UpdateValidIngredientStateIngredientResponse{
		Result: mpgrpcconverters.ConvertValidIngredientStateIngredientToGRPCValidIngredientStateIngredient(updated),
	}

	return res, nil
}

func (s *ServiceImpl) UpdateValidInstrument(ctx context.Context, request *mealplanning.UpdateValidInstrumentRequest) (*mealplanning.UpdateValidInstrumentResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValue(keys.ValidInstrumentIDKey, request.ValidInstrumentID)
	tracing.AttachToSpan(span, keys.ValidInstrumentIDKey, request.ValidInstrumentID)

	input := mpgrpcconverters.ConvertGRPCValidInstrumentUpdateRequestInputToValidInstrumentUpdateRequestInput(request.Input)

	updated, err := s.validEnumerationsManager.UpdateValidInstrument(ctx, request.ValidInstrumentID, input)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "updating valid instrument")
	}

	res := &mealplanning.UpdateValidInstrumentResponse{
		Result: mpgrpcconverters.ConvertValidInstrumentToGRPCValidInstrument(updated),
	}

	return res, nil
}

func (s *ServiceImpl) UpdateValidMeasurementUnit(ctx context.Context, request *mealplanning.UpdateValidMeasurementUnitRequest) (*mealplanning.UpdateValidMeasurementUnitResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValue(keys.ValidMeasurementUnitIDKey, request.ValidMeasurementUnitID)
	tracing.AttachToSpan(span, keys.ValidMeasurementUnitIDKey, request.ValidMeasurementUnitID)

	input := mpgrpcconverters.ConvertGRPCValidMeasurementUnitUpdateRequestInputToValidMeasurementUnitUpdateRequestInput(request.Input)

	updated, err := s.validEnumerationsManager.UpdateValidMeasurementUnit(ctx, request.ValidMeasurementUnitID, input)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "updating valid measurement unit")
	}

	res := &mealplanning.UpdateValidMeasurementUnitResponse{
		Result: mpgrpcconverters.ConvertValidMeasurementUnitToGRPCValidMeasurementUnit(updated),
	}

	return res, nil
}

func (s *ServiceImpl) UpdateValidMeasurementUnitConversion(ctx context.Context, request *mealplanning.UpdateValidMeasurementUnitConversionRequest) (*mealplanning.UpdateValidMeasurementUnitConversionResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValue(keys.ValidMeasurementUnitConversionIDKey, request.ValidMeasurementUnitConversionID)
	tracing.AttachToSpan(span, keys.ValidMeasurementUnitConversionIDKey, request.ValidMeasurementUnitConversionID)

	input := mpgrpcconverters.ConvertGRPCValidMeasurementUnitConversionUpdateRequestInputToValidMeasurementUnitConversionUpdateRequestInput(request.Input)

	updated, err := s.validEnumerationsManager.UpdateValidMeasurementUnitConversion(ctx, request.ValidMeasurementUnitConversionID, input)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "updating valid measurement unit conversion")
	}

	res := &mealplanning.UpdateValidMeasurementUnitConversionResponse{
		Result: mpgrpcconverters.ConvertValidMeasurementUnitConversionToGRPCValidMeasurementUnitConversion(updated),
	}

	return res, nil
}

func (s *ServiceImpl) UpdateValidPreparation(ctx context.Context, request *mealplanning.UpdateValidPreparationRequest) (*mealplanning.UpdateValidPreparationResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValue(keys.ValidPreparationIDKey, request.ValidPreparationID)
	tracing.AttachToSpan(span, keys.ValidPreparationIDKey, request.ValidPreparationID)

	input := mpgrpcconverters.ConvertGRPCValidPreparationUpdateRequestInputToValidPreparationUpdateRequestInput(request.Input)

	updated, err := s.validEnumerationsManager.UpdateValidPreparation(ctx, request.ValidPreparationID, input)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "updating valid preparation")
	}

	res := &mealplanning.UpdateValidPreparationResponse{
		Result: mpgrpcconverters.ConvertValidPreparationToGRPCValidPreparation(updated),
	}

	return res, nil
}

func (s *ServiceImpl) UpdateValidPreparationInstrument(ctx context.Context, request *mealplanning.UpdateValidPreparationInstrumentRequest) (*mealplanning.UpdateValidPreparationInstrumentResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValue(keys.ValidPreparationInstrumentIDKey, request.ValidPreparationInstrumentID)
	tracing.AttachToSpan(span, keys.ValidPreparationInstrumentIDKey, request.ValidPreparationInstrumentID)

	input := mpgrpcconverters.ConvertGRPCValidPreparationInstrumentUpdateRequestInputToValidPreparationInstrumentUpdateRequestInput(request.Input)

	updated, err := s.validEnumerationsManager.UpdateValidPreparationInstrument(ctx, request.ValidPreparationInstrumentID, input)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "updating valid preparation instrument")
	}

	res := &mealplanning.UpdateValidPreparationInstrumentResponse{
		Result: mpgrpcconverters.ConvertValidPreparationInstrumentToGRPCValidPreparationInstrument(updated),
	}

	return res, nil
}

func (s *ServiceImpl) UpdateValidPreparationVessel(ctx context.Context, request *mealplanning.UpdateValidPreparationVesselRequest) (*mealplanning.UpdateValidPreparationVesselResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValue(keys.ValidPreparationVesselIDKey, request.ValidPreparationVesselID)
	tracing.AttachToSpan(span, keys.ValidPreparationVesselIDKey, request.ValidPreparationVesselID)

	input := mpgrpcconverters.ConvertGRPCValidPreparationVesselUpdateRequestInputToValidPreparationVesselUpdateRequestInput(request.Input)

	updated, err := s.validEnumerationsManager.UpdateValidPreparationVessel(ctx, request.ValidPreparationVesselID, input)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "updating valid preparation vessel")
	}

	res := &mealplanning.UpdateValidPreparationVesselResponse{
		Result: mpgrpcconverters.ConvertValidPreparationVesselToGRPCValidPreparationVessel(updated),
	}

	return res, nil
}

func (s *ServiceImpl) UpdateValidVessel(ctx context.Context, request *mealplanning.UpdateValidVesselRequest) (*mealplanning.UpdateValidVesselResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValue(keys.ValidVesselIDKey, request.ValidVesselID)
	tracing.AttachToSpan(span, keys.ValidVesselIDKey, request.ValidVesselID)

	input := mpgrpcconverters.ConvertGRPCValidVesselUpdateRequestInputToValidVesselUpdateRequestInput(request.Input)

	updated, err := s.validEnumerationsManager.UpdateValidVessel(ctx, request.ValidVesselID, input)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "updating valid vessel")
	}

	res := &mealplanning.UpdateValidVesselResponse{
		Result: mpgrpcconverters.ConvertValidVesselToGRPCValidVessel(updated),
	}

	return res, nil
}
