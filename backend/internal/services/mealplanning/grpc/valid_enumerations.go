package eatinggrpc

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/grpc/messages"
	"github.com/dinnerdonebetter/backend/internal/platform/observability"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
	grpcconverters "github.com/dinnerdonebetter/backend/internal/services/mealplanning/grpc/converters"

	"google.golang.org/grpc/codes"
)

func (s *ServiceImpl) ArchiveValidIngredient(ctx context.Context, request *messages.ArchiveValidIngredientRequest) (*messages.ArchiveValidIngredientResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValue(keys.ValidIngredientIDKey, request.ValidIngredientID)
	tracing.AttachToSpan(span, keys.ValidIngredientIDKey, request.ValidIngredientID)

	if err := s.validEnumerationsManager.ArchiveValidIngredient(ctx, request.ValidIngredientID); err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "archiving valid ingredient")
	}

	res := &messages.ArchiveValidIngredientResponse{}

	return res, nil
}

func (s *ServiceImpl) ArchiveValidIngredientGroup(ctx context.Context, request *messages.ArchiveValidIngredientGroupRequest) (*messages.ArchiveValidIngredientGroupResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValue(keys.ValidIngredientGroupIDKey, request.ValidIngredientGroupID)
	tracing.AttachToSpan(span, keys.ValidIngredientGroupIDKey, request.ValidIngredientGroupID)

	if err := s.validEnumerationsManager.ArchiveValidIngredientGroup(ctx, request.ValidIngredientGroupID); err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "archiving valid ingredient group")
	}

	res := &messages.ArchiveValidIngredientGroupResponse{}

	return res, nil
}

func (s *ServiceImpl) ArchiveValidIngredientMeasurementUnit(ctx context.Context, request *messages.ArchiveValidIngredientMeasurementUnitRequest) (*messages.ArchiveValidIngredientMeasurementUnitResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValue(keys.ValidIngredientMeasurementUnitIDKey, request.ValidIngredientMeasurementUnitID)
	tracing.AttachToSpan(span, keys.ValidIngredientMeasurementUnitIDKey, request.ValidIngredientMeasurementUnitID)

	if err := s.validEnumerationsManager.ArchiveValidIngredientMeasurementUnit(ctx, request.ValidIngredientMeasurementUnitID); err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "archiving valid ingredient measurement unit")
	}

	res := &messages.ArchiveValidIngredientMeasurementUnitResponse{}

	return res, nil
}

func (s *ServiceImpl) ArchiveValidIngredientPreparation(ctx context.Context, request *messages.ArchiveValidIngredientPreparationRequest) (*messages.ArchiveValidIngredientPreparationResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValue(keys.ValidIngredientPreparationIDKey, request.ValidIngredientPreparationID)
	tracing.AttachToSpan(span, keys.ValidIngredientPreparationIDKey, request.ValidIngredientPreparationID)

	if err := s.validEnumerationsManager.ArchiveValidIngredientPreparation(ctx, request.ValidIngredientPreparationID); err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "archiving valid ingredient preparation")
	}

	res := &messages.ArchiveValidIngredientPreparationResponse{}

	return res, nil
}

func (s *ServiceImpl) ArchiveValidIngredientState(ctx context.Context, request *messages.ArchiveValidIngredientStateRequest) (*messages.ArchiveValidIngredientStateResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValue(keys.ValidIngredientStateIDKey, request.ValidIngredientStateID)
	tracing.AttachToSpan(span, keys.ValidIngredientStateIDKey, request.ValidIngredientStateID)

	if err := s.validEnumerationsManager.ArchiveValidIngredientState(ctx, request.ValidIngredientStateID); err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "archiving valid ingredient state")
	}

	res := &messages.ArchiveValidIngredientStateResponse{}

	return res, nil
}

func (s *ServiceImpl) ArchiveValidIngredientStateIngredient(ctx context.Context, request *messages.ArchiveValidIngredientStateIngredientRequest) (*messages.ArchiveValidIngredientStateIngredientResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValue(keys.ValidIngredientStateIngredientIDKey, request.ValidIngredientStateIngredientID)
	tracing.AttachToSpan(span, keys.ValidIngredientStateIngredientIDKey, request.ValidIngredientStateIngredientID)

	if err := s.validEnumerationsManager.ArchiveValidIngredientStateIngredient(ctx, request.ValidIngredientStateIngredientID); err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "archiving valid ingredient state ingredient")
	}

	res := &messages.ArchiveValidIngredientStateIngredientResponse{}

	return res, nil
}

func (s *ServiceImpl) ArchiveValidInstrument(ctx context.Context, request *messages.ArchiveValidInstrumentRequest) (*messages.ArchiveValidInstrumentResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValue(keys.ValidInstrumentIDKey, request.ValidInstrumentID)
	tracing.AttachToSpan(span, keys.ValidInstrumentIDKey, request.ValidInstrumentID)

	if err := s.validEnumerationsManager.ArchiveValidInstrument(ctx, request.ValidInstrumentID); err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "archiving valid instrument")
	}

	res := &messages.ArchiveValidInstrumentResponse{}

	return res, nil
}

func (s *ServiceImpl) ArchiveValidMeasurementUnit(ctx context.Context, request *messages.ArchiveValidMeasurementUnitRequest) (*messages.ArchiveValidMeasurementUnitResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValue(keys.ValidMeasurementUnitIDKey, request.ValidMeasurementUnitID)
	tracing.AttachToSpan(span, keys.ValidMeasurementUnitIDKey, request.ValidMeasurementUnitID)

	if err := s.validEnumerationsManager.ArchiveValidMeasurementUnit(ctx, request.ValidMeasurementUnitID); err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "archiving valid measurement unit")
	}

	res := &messages.ArchiveValidMeasurementUnitResponse{}

	return res, nil
}

func (s *ServiceImpl) ArchiveValidMeasurementUnitConversion(ctx context.Context, request *messages.ArchiveValidMeasurementUnitConversionRequest) (*messages.ArchiveValidMeasurementUnitConversionResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValue(keys.ValidMeasurementUnitConversionIDKey, request.ValidMeasurementUnitConversionID)
	tracing.AttachToSpan(span, keys.ValidMeasurementUnitConversionIDKey, request.ValidMeasurementUnitConversionID)

	if err := s.validEnumerationsManager.ArchiveValidMeasurementUnitConversion(ctx, request.ValidMeasurementUnitConversionID); err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "archiving valid measurement unit conversion")
	}

	res := &messages.ArchiveValidMeasurementUnitConversionResponse{}

	return res, nil
}

func (s *ServiceImpl) ArchiveValidPreparation(ctx context.Context, request *messages.ArchiveValidPreparationRequest) (*messages.ArchiveValidPreparationResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValue(keys.ValidPreparationIDKey, request.ValidPreparationID)
	tracing.AttachToSpan(span, keys.ValidPreparationIDKey, request.ValidPreparationID)

	if err := s.validEnumerationsManager.ArchiveValidPreparation(ctx, request.ValidPreparationID); err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "archiving valid preparation")
	}

	res := &messages.ArchiveValidPreparationResponse{}

	return res, nil
}

func (s *ServiceImpl) ArchiveValidPreparationInstrument(ctx context.Context, request *messages.ArchiveValidPreparationInstrumentRequest) (*messages.ArchiveValidPreparationInstrumentResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValue(keys.ValidPreparationInstrumentIDKey, request.ValidPreparationInstrumentID)
	tracing.AttachToSpan(span, keys.ValidPreparationInstrumentIDKey, request.ValidPreparationInstrumentID)

	if err := s.validEnumerationsManager.ArchiveValidPreparationInstrument(ctx, request.ValidPreparationInstrumentID); err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "archiving valid preparation instrument")
	}

	res := &messages.ArchiveValidPreparationInstrumentResponse{}

	return res, nil
}

func (s *ServiceImpl) ArchiveValidPreparationVessel(ctx context.Context, request *messages.ArchiveValidPreparationVesselRequest) (*messages.ArchiveValidPreparationVesselResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValue(keys.ValidPreparationVesselIDKey, request.ValidPreparationVesselID)
	tracing.AttachToSpan(span, keys.ValidPreparationVesselIDKey, request.ValidPreparationVesselID)

	if err := s.validEnumerationsManager.ArchiveValidPreparationVessel(ctx, request.ValidPreparationVesselID); err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "archiving valid preparation vessel")
	}

	res := &messages.ArchiveValidPreparationVesselResponse{}

	return res, nil
}

func (s *ServiceImpl) ArchiveValidVessel(ctx context.Context, request *messages.ArchiveValidVesselRequest) (*messages.ArchiveValidVesselResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValue(keys.ValidVesselIDKey, request.ValidVesselID)
	tracing.AttachToSpan(span, keys.ValidVesselIDKey, request.ValidVesselID)

	if err := s.validEnumerationsManager.ArchiveValidVessel(ctx, request.ValidVesselID); err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "archiving valid vessel")
	}

	res := &messages.ArchiveValidVesselResponse{}

	return res, nil
}

func (s *ServiceImpl) CreateValidIngredient(ctx context.Context, request *messages.CreateValidIngredientRequest) (*messages.CreateValidIngredientResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	created, err := s.validEnumerationsManager.CreateValidIngredient(ctx, grpcconverters.ConvertGRPCCreateValidIngredientRequestToValidIngredientCreationRequestInput(request))
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "creating valid ingredient")
	}
	tracing.AttachToSpan(span, keys.ValidIngredientIDKey, created.ID)

	result := &messages.CreateValidIngredientResponse{
		Result: grpcconverters.ConvertValidIngredientToGRPCValidIngredient(created),
	}

	return result, nil
}

func (s *ServiceImpl) CreateValidIngredientGroup(ctx context.Context, request *messages.CreateValidIngredientGroupRequest) (*messages.CreateValidIngredientGroupResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	created, err := s.validEnumerationsManager.CreateValidIngredientGroup(ctx, grpcconverters.ConvertGRPCCreateValidIngredientGroupRequestToValidIngredientGroupCreationRequestInput(request))
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "creating valid ingredient")
	}
	tracing.AttachToSpan(span, keys.ValidIngredientGroupIDKey, created.ID)

	result := &messages.CreateValidIngredientGroupResponse{
		Result: grpcconverters.ConvertValidIngredientGroupToGRPCValidIngredientGroup(created),
	}

	return result, nil
}

func (s *ServiceImpl) CreateValidIngredientMeasurementUnit(ctx context.Context, request *messages.CreateValidIngredientMeasurementUnitRequest) (*messages.CreateValidIngredientMeasurementUnitResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	created, err := s.validEnumerationsManager.CreateValidIngredientMeasurementUnit(ctx, grpcconverters.ConvertGRPCCreateValidIngredientMeasurementUnitRequestToValidIngredientMeasurementUnitCreationRequestInput(request))
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "creating valid ingredient measurement unit")
	}
	tracing.AttachToSpan(span, keys.ValidIngredientMeasurementUnitIDKey, created.ID)

	result := &messages.CreateValidIngredientMeasurementUnitResponse{
		Result: grpcconverters.ConvertValidIngredientMeasurementUnitToGRPCValidIngredientMeasurementUnit(created),
	}

	return result, nil
}

func (s *ServiceImpl) CreateValidIngredientPreparation(ctx context.Context, request *messages.CreateValidIngredientPreparationRequest) (*messages.CreateValidIngredientPreparationResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	created, err := s.validEnumerationsManager.CreateValidIngredientPreparation(ctx, grpcconverters.ConvertGRPCCreateValidIngredientPreparationRequestToValidIngredientPreparationCreationRequestInput(request))
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "creating valid ingredient preparation")
	}
	tracing.AttachToSpan(span, keys.ValidIngredientPreparationIDKey, created.ID)

	result := &messages.CreateValidIngredientPreparationResponse{
		Result: grpcconverters.ConvertValidIngredientPreparationToGRPCValidIngredientPreparation(created),
	}

	return result, nil
}

func (s *ServiceImpl) CreateValidIngredientState(ctx context.Context, request *messages.CreateValidIngredientStateRequest) (*messages.CreateValidIngredientStateResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	created, err := s.validEnumerationsManager.CreateValidIngredientState(ctx, grpcconverters.ConvertGRPCCreateValidIngredientStateRequestToValidIngredientStateCreationRequestInput(request))
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "creating valid ingredient state")
	}
	tracing.AttachToSpan(span, keys.ValidIngredientStateIDKey, created.ID)

	result := &messages.CreateValidIngredientStateResponse{
		Result: grpcconverters.ConvertValidIngredientStateToGRPCValidIngredientState(created),
	}

	return result, nil
}

func (s *ServiceImpl) CreateValidIngredientStateIngredient(ctx context.Context, request *messages.CreateValidIngredientStateIngredientRequest) (*messages.CreateValidIngredientStateIngredientResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	created, err := s.validEnumerationsManager.CreateValidIngredientStateIngredient(ctx, grpcconverters.ConvertGRPCCreateValidIngredientStateIngredientRequestToValidIngredientStateIngredientCreationRequestInput(request))
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "creating valid ingredient state ingredient")
	}
	tracing.AttachToSpan(span, keys.ValidIngredientStateIngredientIDKey, created.ID)

	result := &messages.CreateValidIngredientStateIngredientResponse{
		Result: grpcconverters.ConvertValidIngredientStateIngredientToGRPCValidIngredientStateIngredient(created),
	}

	return result, nil
}

func (s *ServiceImpl) CreateValidInstrument(ctx context.Context, request *messages.CreateValidInstrumentRequest) (*messages.CreateValidInstrumentResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	created, err := s.validEnumerationsManager.CreateValidInstrument(ctx, grpcconverters.ConvertGRPCCreateValidInstrumentRequestToValidInstrumentCreationRequestInput(request))
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "creating valid instrument")
	}
	tracing.AttachToSpan(span, keys.ValidInstrumentIDKey, created.ID)

	result := &messages.CreateValidInstrumentResponse{
		Result: grpcconverters.ConvertValidInstrumentToGRPCValidInstrument(created),
	}

	return result, nil
}

func (s *ServiceImpl) CreateValidMeasurementUnit(ctx context.Context, request *messages.CreateValidMeasurementUnitRequest) (*messages.CreateValidMeasurementUnitResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	created, err := s.validEnumerationsManager.CreateValidMeasurementUnit(ctx, grpcconverters.ConvertGRPCCreateValidMeasurementUnitRequestToValidMeasurementUnitCreationRequestInput(request))
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "creating valid measurement unit")
	}
	tracing.AttachToSpan(span, keys.ValidMeasurementUnitIDKey, created.ID)

	result := &messages.CreateValidMeasurementUnitResponse{
		Result: grpcconverters.ConvertValidMeasurementUnitToGRPCValidMeasurementUnit(created),
	}

	return result, nil
}

func (s *ServiceImpl) CreateValidMeasurementUnitConversion(ctx context.Context, request *messages.CreateValidMeasurementUnitConversionRequest) (*messages.CreateValidMeasurementUnitConversionResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	created, err := s.validEnumerationsManager.CreateValidMeasurementUnitConversion(ctx, grpcconverters.ConvertGRPCCreateValidMeasurementUnitConversionRequestToValidMeasurementUnitConversionCreationRequestInput(request))
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "creating valid measurement unit conversion")
	}
	tracing.AttachToSpan(span, keys.ValidMeasurementUnitConversionIDKey, created.ID)

	result := &messages.CreateValidMeasurementUnitConversionResponse{
		Result: grpcconverters.ConvertValidMeasurementUnitConversionToGRPCValidMeasurementUnitConversion(created),
	}

	return result, nil
}

func (s *ServiceImpl) CreateValidPreparation(ctx context.Context, request *messages.CreateValidPreparationRequest) (*messages.CreateValidPreparationResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	created, err := s.validEnumerationsManager.CreateValidPreparation(ctx, grpcconverters.ConvertGRPCCreateValidPreparationRequestToValidPreparationCreationRequestInput(request))
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "creating valid preparation")
	}
	tracing.AttachToSpan(span, keys.ValidPreparationIDKey, created.ID)

	result := &messages.CreateValidPreparationResponse{
		Result: grpcconverters.ConvertValidPreparationToGRPCValidPreparation(created),
	}

	return result, nil
}

func (s *ServiceImpl) CreateValidPreparationInstrument(ctx context.Context, request *messages.CreateValidPreparationInstrumentRequest) (*messages.CreateValidPreparationInstrumentResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	created, err := s.validEnumerationsManager.CreateValidPreparationInstrument(ctx, grpcconverters.ConvertGRPCCreateValidPreparationInstrumentRequestToValidPreparationInstrumentCreationRequestInput(request))
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "creating valid preparation instrument")
	}
	tracing.AttachToSpan(span, keys.ValidPreparationInstrumentIDKey, created.ID)

	result := &messages.CreateValidPreparationInstrumentResponse{
		Result: grpcconverters.ConvertValidPreparationInstrumentToGRPCValidPreparationInstrument(created),
	}

	return result, nil
}

func (s *ServiceImpl) CreateValidPreparationVessel(ctx context.Context, request *messages.CreateValidPreparationVesselRequest) (*messages.CreateValidPreparationVesselResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	created, err := s.validEnumerationsManager.CreateValidPreparationVessel(ctx, grpcconverters.ConvertGRPCCreateValidPreparationVesselRequestToValidPreparationVesselCreationRequestInput(request))
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "creating valid preparation vessel")
	}
	tracing.AttachToSpan(span, keys.ValidPreparationVesselIDKey, created.ID)

	result := &messages.CreateValidPreparationVesselResponse{
		Result: grpcconverters.ConvertValidPreparationVesselToGRPCValidPreparationVessel(created),
	}

	return result, nil
}

func (s *ServiceImpl) CreateValidVessel(ctx context.Context, request *messages.CreateValidVesselRequest) (*messages.CreateValidVesselResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	created, err := s.validEnumerationsManager.CreateValidVessel(ctx, grpcconverters.ConvertGRPCCreateValidVesselRequestToValidVesselCreationRequestInput(request))
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "creating valid vessel")
	}
	tracing.AttachToSpan(span, keys.ValidVesselIDKey, created.ID)

	result := &messages.CreateValidVesselResponse{
		Result: grpcconverters.ConvertValidVesselToGRPCValidVessel(created),
	}

	return result, nil
}

func (s *ServiceImpl) GetRandomValidIngredient(ctx context.Context, _ *messages.GetRandomValidIngredientRequest) (*messages.GetRandomValidIngredientResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	selected, err := s.validEnumerationsManager.RandomValidIngredient(ctx)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "fetching random valid ingredient")
	}
	tracing.AttachToSpan(span, keys.ValidIngredientIDKey, selected.ID)

	res := &messages.GetRandomValidIngredientResponse{
		Result: grpcconverters.ConvertValidIngredientToGRPCValidIngredient(selected),
	}

	return res, nil
}

func (s *ServiceImpl) GetRandomValidInstrument(ctx context.Context, _ *messages.GetRandomValidInstrumentRequest) (*messages.GetRandomValidInstrumentResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	selected, err := s.validEnumerationsManager.RandomValidInstrument(ctx)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "fetching random valid instrument")
	}
	tracing.AttachToSpan(span, keys.ValidInstrumentIDKey, selected.ID)

	res := &messages.GetRandomValidInstrumentResponse{
		Result: grpcconverters.ConvertValidInstrumentToGRPCValidInstrument(selected),
	}

	return res, nil
}

func (s *ServiceImpl) GetRandomValidPreparation(ctx context.Context, _ *messages.GetRandomValidPreparationRequest) (*messages.GetRandomValidPreparationResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	selected, err := s.validEnumerationsManager.RandomValidPreparation(ctx)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "fetching random valid preparation")
	}
	tracing.AttachToSpan(span, keys.ValidPreparationIDKey, selected.ID)

	res := &messages.GetRandomValidPreparationResponse{
		Result: grpcconverters.ConvertValidPreparationToGRPCValidPreparation(selected),
	}

	return res, nil
}

func (s *ServiceImpl) GetRandomValidVessel(ctx context.Context, _ *messages.GetRandomValidVesselRequest) (*messages.GetRandomValidVesselResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	selected, err := s.validEnumerationsManager.RandomValidVessel(ctx)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "fetching random valid vessel")
	}
	tracing.AttachToSpan(span, keys.ValidVesselIDKey, selected.ID)

	res := &messages.GetRandomValidVesselResponse{
		Result: grpcconverters.ConvertValidVesselToGRPCValidVessel(selected),
	}

	return res, nil
}

func (s *ServiceImpl) GetValidIngredient(ctx context.Context, request *messages.GetValidIngredientRequest) (*messages.GetValidIngredientResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValue(keys.ValidIngredientIDKey, request.ValidIngredientID)
	tracing.AttachToSpan(span, keys.ValidIngredientIDKey, request.ValidIngredientID)

	x, err := s.validEnumerationsManager.ReadValidIngredient(ctx, request.ValidIngredientID)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "fetching valid ingredient")
	}

	res := &messages.GetValidIngredientResponse{
		Result: grpcconverters.ConvertValidIngredientToGRPCValidIngredient(x),
	}

	return res, nil
}

func (s *ServiceImpl) GetValidIngredientGroup(ctx context.Context, request *messages.GetValidIngredientGroupRequest) (*messages.GetValidIngredientGroupResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValue(keys.ValidIngredientGroupIDKey, request.ValidIngredientGroupID)
	tracing.AttachToSpan(span, keys.ValidIngredientGroupIDKey, request.ValidIngredientGroupID)

	x, err := s.validEnumerationsManager.ReadValidIngredientGroup(ctx, request.ValidIngredientGroupID)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "fetching valid ingredient group")
	}

	res := &messages.GetValidIngredientGroupResponse{
		Result: grpcconverters.ConvertValidIngredientGroupToGRPCValidIngredientGroup(x),
	}

	return res, nil
}

func (s *ServiceImpl) GetValidIngredientGroups(ctx context.Context, request *messages.GetValidIngredientGroupsRequest) (*messages.GetValidIngredientGroupsResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)
	filter := grpcconverters.ConvertGRPCQueryFilterToQueryFilter(request.Filter)
	tracing.AttachQueryFilterToSpan(span, filter)

	x, _, err := s.validEnumerationsManager.ListValidIngredientGroups(ctx, filter)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "fetching valid ingredient groups")
	}

	res := &messages.GetValidIngredientGroupsResponse{
		Filter: request.Filter,
	}
	for _, y := range x {
		res.Results = append(res.Results, grpcconverters.ConvertValidIngredientGroupToGRPCValidIngredientGroup(y))
	}

	return res, nil
}

func (s *ServiceImpl) GetValidIngredientMeasurementUnit(ctx context.Context, request *messages.GetValidIngredientMeasurementUnitRequest) (*messages.GetValidIngredientMeasurementUnitResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValue(keys.ValidIngredientMeasurementUnitIDKey, request.ValidIngredientMeasurementUnitID)
	tracing.AttachToSpan(span, keys.ValidIngredientMeasurementUnitIDKey, request.ValidIngredientMeasurementUnitID)

	x, err := s.validEnumerationsManager.ReadValidIngredientMeasurementUnit(ctx, request.ValidIngredientMeasurementUnitID)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "fetching valid ingredient measurement unit")
	}

	res := &messages.GetValidIngredientMeasurementUnitResponse{
		Result: grpcconverters.ConvertValidIngredientMeasurementUnitToGRPCValidIngredientMeasurementUnit(x),
	}

	return res, nil
}

func (s *ServiceImpl) GetValidIngredientMeasurementUnits(ctx context.Context, request *messages.GetValidIngredientMeasurementUnitsRequest) (*messages.GetValidIngredientMeasurementUnitsResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)
	filter := grpcconverters.ConvertGRPCQueryFilterToQueryFilter(request.Filter)
	tracing.AttachQueryFilterToSpan(span, filter)

	x, _, err := s.validEnumerationsManager.ListValidIngredientMeasurementUnits(ctx, filter)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "fetching valid ingredient measurement units")
	}

	res := &messages.GetValidIngredientMeasurementUnitsResponse{
		Filter: request.Filter,
	}
	for _, y := range x {
		res.Results = append(res.Results, grpcconverters.ConvertValidIngredientMeasurementUnitToGRPCValidIngredientMeasurementUnit(y))
	}

	return res, nil
}

func (s *ServiceImpl) GetValidIngredientMeasurementUnitsByIngredient(ctx context.Context, request *messages.GetValidIngredientMeasurementUnitsByIngredientRequest) (*messages.GetValidIngredientMeasurementUnitsByIngredientResponse, error) {
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

	res := &messages.GetValidIngredientMeasurementUnitsByIngredientResponse{
		Filter: request.Filter,
	}

	for _, y := range x {
		res.Results = append(res.Results, grpcconverters.ConvertValidIngredientMeasurementUnitToGRPCValidIngredientMeasurementUnit(y))
	}

	return res, nil
}

func (s *ServiceImpl) GetValidIngredientMeasurementUnitsByMeasurementUnit(ctx context.Context, request *messages.GetValidIngredientMeasurementUnitsByMeasurementUnitRequest) (*messages.GetValidIngredientMeasurementUnitsByMeasurementUnitResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)
	filter := grpcconverters.ConvertGRPCQueryFilterToQueryFilter(request.Filter)
	tracing.AttachQueryFilterToSpan(span, filter)

	x, err := s.validEnumerationsManager.SearchValidIngredientMeasurementUnitsByMeasurementUnit(ctx, request.ValidMeasurementUnitID, filter)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "fetching valid ingredient measurement units by measurement unit")
	}

	res := &messages.GetValidIngredientMeasurementUnitsByMeasurementUnitResponse{
		Filter: request.Filter,
	}

	for _, y := range x {
		res.Results = append(res.Results, grpcconverters.ConvertValidIngredientMeasurementUnitToGRPCValidIngredientMeasurementUnit(y))
	}

	return res, nil
}

func (s *ServiceImpl) GetValidIngredientPreparation(ctx context.Context, request *messages.GetValidIngredientPreparationRequest) (*messages.GetValidIngredientPreparationResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValue(keys.ValidIngredientPreparationIDKey, request.ValidIngredientPreparationID)
	tracing.AttachToSpan(span, keys.ValidIngredientPreparationIDKey, request.ValidIngredientPreparationID)

	x, err := s.validEnumerationsManager.ReadValidIngredientPreparation(ctx, request.ValidIngredientPreparationID)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "fetching valid ingredient preparation")
	}

	res := &messages.GetValidIngredientPreparationResponse{
		Result: grpcconverters.ConvertValidIngredientPreparationToGRPCValidIngredientPreparation(x),
	}

	return res, nil
}

func (s *ServiceImpl) GetValidIngredientPreparations(ctx context.Context, request *messages.GetValidIngredientPreparationsRequest) (*messages.GetValidIngredientPreparationsResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)
	filter := grpcconverters.ConvertGRPCQueryFilterToQueryFilter(request.Filter)
	tracing.AttachQueryFilterToSpan(span, filter)

	x, _, err := s.validEnumerationsManager.ListValidIngredientPreparations(ctx, filter)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "fetching valid ingredient preparations")
	}

	res := &messages.GetValidIngredientPreparationsResponse{
		Filter: request.Filter,
	}

	for _, y := range x {
		res.Results = append(res.Results, grpcconverters.ConvertValidIngredientPreparationToGRPCValidIngredientPreparation(y))
	}

	return res, nil
}

func (s *ServiceImpl) GetValidIngredientPreparationsByIngredient(ctx context.Context, request *messages.GetValidIngredientPreparationsByIngredientRequest) (*messages.GetValidIngredientPreparationsByIngredientResponse, error) {
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

	res := &messages.GetValidIngredientPreparationsByIngredientResponse{
		Filter: request.Filter,
	}

	for _, y := range x {
		res.Results = append(res.Results, grpcconverters.ConvertValidIngredientPreparationToGRPCValidIngredientPreparation(y))
	}

	return res, nil
}

func (s *ServiceImpl) GetValidIngredientPreparationsByPreparation(ctx context.Context, request *messages.GetValidIngredientPreparationsByPreparationRequest) (*messages.GetValidIngredientPreparationsByPreparationResponse, error) {
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

	res := &messages.GetValidIngredientPreparationsByPreparationResponse{
		Filter: request.Filter,
	}

	for _, y := range x {
		res.Results = append(res.Results, grpcconverters.ConvertValidIngredientPreparationToGRPCValidIngredientPreparation(y))
	}

	return res, nil
}

func (s *ServiceImpl) GetValidIngredientState(ctx context.Context, request *messages.GetValidIngredientStateRequest) (*messages.GetValidIngredientStateResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValue(keys.ValidIngredientStateIDKey, request.ValidIngredientStateID)
	tracing.AttachToSpan(span, keys.ValidIngredientStateIDKey, request.ValidIngredientStateID)

	x, err := s.validEnumerationsManager.ReadValidIngredientState(ctx, request.ValidIngredientStateID)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "fetching valid ingredient state")
	}

	res := &messages.GetValidIngredientStateResponse{
		Result: grpcconverters.ConvertValidIngredientStateToGRPCValidIngredientState(x),
	}

	return res, nil
}

func (s *ServiceImpl) GetValidIngredientStateIngredient(ctx context.Context, request *messages.GetValidIngredientStateIngredientRequest) (*messages.GetValidIngredientStateIngredientResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValue(keys.ValidIngredientStateIngredientIDKey, request.ValidIngredientStateIngredientID)
	tracing.AttachToSpan(span, keys.ValidIngredientStateIngredientIDKey, request.ValidIngredientStateIngredientID)

	x, err := s.validEnumerationsManager.ReadValidIngredientStateIngredient(ctx, request.ValidIngredientStateIngredientID)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "fetching valid ingredient state ingredient")
	}

	res := &messages.GetValidIngredientStateIngredientResponse{
		Result: grpcconverters.ConvertValidIngredientStateIngredientToGRPCValidIngredientStateIngredient(x),
	}

	return res, nil
}

func (s *ServiceImpl) GetValidIngredientStateIngredients(ctx context.Context, request *messages.GetValidIngredientStateIngredientsRequest) (*messages.GetValidIngredientStateIngredientsResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)
	filter := grpcconverters.ConvertGRPCQueryFilterToQueryFilter(request.Filter)
	tracing.AttachQueryFilterToSpan(span, filter)

	x, _, err := s.validEnumerationsManager.ListValidIngredientStateIngredients(ctx, filter)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "fetching valid ingredient state ingredients")
	}

	res := &messages.GetValidIngredientStateIngredientsResponse{
		Filter: request.Filter,
	}

	for _, y := range x {
		res.Results = append(res.Results, grpcconverters.ConvertValidIngredientStateIngredientToGRPCValidIngredientStateIngredient(y))
	}

	return res, nil
}

func (s *ServiceImpl) GetValidIngredientStateIngredientsByIngredient(ctx context.Context, request *messages.GetValidIngredientStateIngredientsByIngredientRequest) (*messages.GetValidIngredientStateIngredientsByIngredientResponse, error) {
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

	res := &messages.GetValidIngredientStateIngredientsByIngredientResponse{
		Filter: request.Filter,
	}

	for _, y := range x {
		res.Results = append(res.Results, grpcconverters.ConvertValidIngredientStateIngredientToGRPCValidIngredientStateIngredient(y))
	}

	return res, nil
}

func (s *ServiceImpl) GetValidIngredientStateIngredientsByIngredientState(ctx context.Context, request *messages.GetValidIngredientStateIngredientsByIngredientStateRequest) (*messages.GetValidIngredientStateIngredientsByIngredientStateResponse, error) {
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

	res := &messages.GetValidIngredientStateIngredientsByIngredientStateResponse{
		Filter: request.Filter,
	}

	for _, y := range x {
		res.Results = append(res.Results, grpcconverters.ConvertValidIngredientStateIngredientToGRPCValidIngredientStateIngredient(y))
	}

	return res, nil
}

func (s *ServiceImpl) GetValidIngredientStates(ctx context.Context, request *messages.GetValidIngredientStatesRequest) (*messages.GetValidIngredientStatesResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)
	filter := grpcconverters.ConvertGRPCQueryFilterToQueryFilter(request.Filter)
	tracing.AttachQueryFilterToSpan(span, filter)

	x, _, err := s.validEnumerationsManager.ListValidIngredientStates(ctx, filter)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "fetching valid ingredient states")
	}

	res := &messages.GetValidIngredientStatesResponse{
		Filter: request.Filter,
	}

	for _, y := range x {
		res.Results = append(res.Results, grpcconverters.ConvertValidIngredientStateToGRPCValidIngredientState(y))
	}

	return res, nil
}

func (s *ServiceImpl) GetValidIngredients(ctx context.Context, request *messages.GetValidIngredientsRequest) (*messages.GetValidIngredientsResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)
	filter := grpcconverters.ConvertGRPCQueryFilterToQueryFilter(request.Filter)
	tracing.AttachQueryFilterToSpan(span, filter)

	x, _, err := s.validEnumerationsManager.ListValidIngredients(ctx, filter)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "fetching valid ingredients")
	}

	res := &messages.GetValidIngredientsResponse{
		Filter: request.Filter,
	}

	for _, y := range x {
		res.Results = append(res.Results, grpcconverters.ConvertValidIngredientToGRPCValidIngredient(y))
	}

	return res, nil
}

func (s *ServiceImpl) GetValidInstrument(ctx context.Context, request *messages.GetValidInstrumentRequest) (*messages.GetValidInstrumentResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValue(keys.ValidInstrumentIDKey, request.ValidInstrumentID)
	tracing.AttachToSpan(span, keys.ValidInstrumentIDKey, request.ValidInstrumentID)

	x, err := s.validEnumerationsManager.ReadValidInstrument(ctx, request.ValidInstrumentID)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "fetching valid instrument")
	}

	res := &messages.GetValidInstrumentResponse{
		Result: grpcconverters.ConvertValidInstrumentToGRPCValidInstrument(x),
	}

	return res, nil
}

func (s *ServiceImpl) GetValidInstruments(ctx context.Context, request *messages.GetValidInstrumentsRequest) (*messages.GetValidInstrumentsResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)
	filter := grpcconverters.ConvertGRPCQueryFilterToQueryFilter(request.Filter)
	tracing.AttachQueryFilterToSpan(span, filter)

	x, _, err := s.validEnumerationsManager.ListValidInstruments(ctx, filter)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "fetching valid instruments")
	}

	res := &messages.GetValidInstrumentsResponse{
		Filter: request.Filter,
	}

	for _, y := range x {
		res.Results = append(res.Results, grpcconverters.ConvertValidInstrumentToGRPCValidInstrument(y))
	}

	return res, nil
}

func (s *ServiceImpl) GetValidMeasurementUnit(ctx context.Context, request *messages.GetValidMeasurementUnitRequest) (*messages.GetValidMeasurementUnitResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValue(keys.ValidMeasurementUnitIDKey, request.ValidMeasurementUnitID)
	tracing.AttachToSpan(span, keys.ValidMeasurementUnitIDKey, request.ValidMeasurementUnitID)

	x, err := s.validEnumerationsManager.ReadValidMeasurementUnit(ctx, request.ValidMeasurementUnitID)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "fetching valid measurement unit")
	}

	res := &messages.GetValidMeasurementUnitResponse{
		Result: grpcconverters.ConvertValidMeasurementUnitToGRPCValidMeasurementUnit(x),
	}

	return res, nil
}

func (s *ServiceImpl) GetValidMeasurementUnitConversion(ctx context.Context, request *messages.GetValidMeasurementUnitConversionRequest) (*messages.GetValidMeasurementUnitConversionResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValue(keys.ValidMeasurementUnitConversionIDKey, request.ValidMeasurementUnitConversionID)
	tracing.AttachToSpan(span, keys.ValidMeasurementUnitConversionIDKey, request.ValidMeasurementUnitConversionID)

	x, err := s.validEnumerationsManager.ReadValidMeasurementUnitConversion(ctx, request.ValidMeasurementUnitConversionID)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "fetching valid measurement unit conversion")
	}

	res := &messages.GetValidMeasurementUnitConversionResponse{
		Result: grpcconverters.ConvertValidMeasurementUnitConversionToGRPCValidMeasurementUnitConversion(x),
	}

	return res, nil
}

func (s *ServiceImpl) GetValidMeasurementUnitConversionsFromUnit(ctx context.Context, request *messages.GetValidMeasurementUnitConversionsFromUnitRequest) (*messages.GetValidMeasurementUnitConversionsFromUnitResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValue(keys.ValidMeasurementUnitIDKey, request.ValidMeasurementUnitID)
	tracing.AttachToSpan(span, keys.ValidMeasurementUnitIDKey, request.ValidMeasurementUnitID)

	x, err := s.validEnumerationsManager.ValidMeasurementUnitConversionsFromMeasurementUnit(ctx, request.ValidMeasurementUnitID)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "fetching valid measurement unit conversions from unit")
	}

	res := &messages.GetValidMeasurementUnitConversionsFromUnitResponse{}

	for _, y := range x {
		res.Results = append(res.Results, grpcconverters.ConvertValidMeasurementUnitConversionToGRPCValidMeasurementUnitConversion(y))
	}

	return res, nil
}

func (s *ServiceImpl) GetValidMeasurementUnitConversionsToUnit(ctx context.Context, request *messages.GetValidMeasurementUnitConversionsToUnitRequest) (*messages.GetValidMeasurementUnitConversionsToUnitResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValue(keys.ValidMeasurementUnitIDKey, request.ValidMeasurementUnitID)
	tracing.AttachToSpan(span, keys.ValidMeasurementUnitIDKey, request.ValidMeasurementUnitID)

	x, err := s.validEnumerationsManager.ValidMeasurementUnitConversionsToMeasurementUnit(ctx, request.ValidMeasurementUnitID)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "fetching valid measurement unit conversions to unit")
	}

	res := &messages.GetValidMeasurementUnitConversionsToUnitResponse{}
	for _, y := range x {
		res.Results = append(res.Results, grpcconverters.ConvertValidMeasurementUnitConversionToGRPCValidMeasurementUnitConversion(y))
	}

	return res, nil
}

func (s *ServiceImpl) GetValidMeasurementUnits(ctx context.Context, request *messages.GetValidMeasurementUnitsRequest) (*messages.GetValidMeasurementUnitsResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)
	filter := grpcconverters.ConvertGRPCQueryFilterToQueryFilter(request.Filter)
	tracing.AttachQueryFilterToSpan(span, filter)

	x, _, err := s.validEnumerationsManager.ListValidMeasurementUnits(ctx, filter)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "fetching valid measurement units")
	}

	res := &messages.GetValidMeasurementUnitsResponse{
		Filter: request.Filter,
	}

	for _, y := range x {
		res.Results = append(res.Results, grpcconverters.ConvertValidMeasurementUnitToGRPCValidMeasurementUnit(y))
	}

	return res, nil
}

func (s *ServiceImpl) GetValidPreparation(ctx context.Context, request *messages.GetValidPreparationRequest) (*messages.GetValidPreparationResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValue(keys.ValidPreparationIDKey, request.ValidPreparationID)
	tracing.AttachToSpan(span, keys.ValidPreparationIDKey, request.ValidPreparationID)

	x, err := s.validEnumerationsManager.ReadValidPreparation(ctx, request.ValidPreparationID)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "fetching valid preparation")
	}

	res := &messages.GetValidPreparationResponse{
		Result: grpcconverters.ConvertValidPreparationToGRPCValidPreparation(x),
	}

	return res, nil
}

func (s *ServiceImpl) GetValidPreparationInstrument(ctx context.Context, request *messages.GetValidPreparationInstrumentRequest) (*messages.GetValidPreparationInstrumentResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValue(keys.ValidPreparationInstrumentIDKey, request.ValidPreparationInstrumentID)
	tracing.AttachToSpan(span, keys.ValidPreparationInstrumentIDKey, request.ValidPreparationInstrumentID)

	x, err := s.validEnumerationsManager.ReadValidPreparationInstrument(ctx, request.ValidPreparationInstrumentID)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "fetching valid preparation instrument")
	}

	res := &messages.GetValidPreparationInstrumentResponse{
		Result: grpcconverters.ConvertValidPreparationInstrumentToGRPCValidPreparationInstrument(x),
	}

	return res, nil
}

func (s *ServiceImpl) GetValidPreparationInstruments(ctx context.Context, request *messages.GetValidPreparationInstrumentsRequest) (*messages.GetValidPreparationInstrumentsResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)
	filter := grpcconverters.ConvertGRPCQueryFilterToQueryFilter(request.Filter)
	tracing.AttachQueryFilterToSpan(span, filter)

	x, _, err := s.validEnumerationsManager.ListValidPreparationInstruments(ctx, filter)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "fetching valid preparation instruments")
	}

	res := &messages.GetValidPreparationInstrumentsResponse{
		Filter: request.Filter,
	}

	for _, y := range x {
		res.Results = append(res.Results, grpcconverters.ConvertValidPreparationInstrumentToGRPCValidPreparationInstrument(y))
	}

	return res, nil
}

func (s *ServiceImpl) GetValidPreparationInstrumentsByInstrument(ctx context.Context, request *messages.GetValidPreparationInstrumentsByInstrumentRequest) (*messages.GetValidPreparationInstrumentsByInstrumentResponse, error) {
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

	res := &messages.GetValidPreparationInstrumentsByInstrumentResponse{
		Filter: request.Filter,
	}

	for _, y := range x {
		res.Results = append(res.Results, grpcconverters.ConvertValidPreparationInstrumentToGRPCValidPreparationInstrument(y))
	}

	return res, nil
}

func (s *ServiceImpl) GetValidPreparationInstrumentsByPreparation(ctx context.Context, request *messages.GetValidPreparationInstrumentsByPreparationRequest) (*messages.GetValidPreparationInstrumentsByPreparationResponse, error) {
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

	res := &messages.GetValidPreparationInstrumentsByPreparationResponse{
		Filter: request.Filter,
	}

	for _, y := range x {
		res.Results = append(res.Results, grpcconverters.ConvertValidPreparationInstrumentToGRPCValidPreparationInstrument(y))
	}

	return res, nil
}

func (s *ServiceImpl) GetValidPreparationVessel(ctx context.Context, request *messages.GetValidPreparationVesselRequest) (*messages.GetValidPreparationVesselResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValue(keys.ValidPreparationVesselIDKey, request.ValidPreparationVesselID)
	tracing.AttachToSpan(span, keys.ValidPreparationVesselIDKey, request.ValidPreparationVesselID)

	x, err := s.validEnumerationsManager.ReadValidPreparationVessel(ctx, request.ValidPreparationVesselID)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "fetching valid preparation vessel")
	}

	res := &messages.GetValidPreparationVesselResponse{
		Result: grpcconverters.ConvertValidPreparationVesselToGRPCValidPreparationVessel(x),
	}

	return res, nil
}

func (s *ServiceImpl) GetValidPreparationVessels(ctx context.Context, request *messages.GetValidPreparationVesselsRequest) (*messages.GetValidPreparationVesselsResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)
	filter := grpcconverters.ConvertGRPCQueryFilterToQueryFilter(request.Filter)
	tracing.AttachQueryFilterToSpan(span, filter)

	x, _, err := s.validEnumerationsManager.ListValidPreparationVessels(ctx, filter)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "fetching valid preparation vessels")
	}

	res := &messages.GetValidPreparationVesselsResponse{
		Filter: request.Filter,
	}

	for _, y := range x {
		res.Results = append(res.Results, grpcconverters.ConvertValidPreparationVesselToGRPCValidPreparationVessel(y))
	}

	return res, nil
}

func (s *ServiceImpl) GetValidPreparationVesselsByPreparation(ctx context.Context, request *messages.GetValidPreparationVesselsByPreparationRequest) (*messages.GetValidPreparationVesselsByPreparationResponse, error) {
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

	res := &messages.GetValidPreparationVesselsByPreparationResponse{
		Filter: request.Filter,
	}

	for _, y := range x {
		res.Results = append(res.Results, grpcconverters.ConvertValidPreparationVesselToGRPCValidPreparationVessel(y))
	}

	return res, nil
}

func (s *ServiceImpl) GetValidPreparationVesselsByVessel(ctx context.Context, request *messages.GetValidPreparationVesselsByVesselRequest) (*messages.GetValidPreparationVesselsByVesselResponse, error) {
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

	res := &messages.GetValidPreparationVesselsByVesselResponse{
		Filter: request.Filter,
	}

	for _, y := range x {
		res.Results = append(res.Results, grpcconverters.ConvertValidPreparationVesselToGRPCValidPreparationVessel(y))
	}

	return res, nil
}

func (s *ServiceImpl) GetValidPreparations(ctx context.Context, request *messages.GetValidPreparationsRequest) (*messages.GetValidPreparationsResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)
	filter := grpcconverters.ConvertGRPCQueryFilterToQueryFilter(request.Filter)
	tracing.AttachQueryFilterToSpan(span, filter)

	x, _, err := s.validEnumerationsManager.ListValidPreparations(ctx, filter)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "fetching valid preparations")
	}

	res := &messages.GetValidPreparationsResponse{
		Filter: request.Filter,
	}

	for _, y := range x {
		res.Results = append(res.Results, grpcconverters.ConvertValidPreparationToGRPCValidPreparation(y))
	}

	return res, nil
}

func (s *ServiceImpl) GetValidVessel(ctx context.Context, request *messages.GetValidVesselRequest) (*messages.GetValidVesselResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValue(keys.ValidVesselIDKey, request.ValidVesselID)
	tracing.AttachToSpan(span, keys.ValidVesselIDKey, request.ValidVesselID)

	x, err := s.validEnumerationsManager.ReadValidVessel(ctx, request.ValidVesselID)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "fetching valid vessel")
	}

	res := &messages.GetValidVesselResponse{
		Result: grpcconverters.ConvertValidVesselToGRPCValidVessel(x),
	}

	return res, nil
}

func (s *ServiceImpl) GetValidVessels(ctx context.Context, request *messages.GetValidVesselsRequest) (*messages.GetValidVesselsResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)
	filter := grpcconverters.ConvertGRPCQueryFilterToQueryFilter(request.Filter)
	tracing.AttachQueryFilterToSpan(span, filter)

	x, _, err := s.validEnumerationsManager.ListValidVessels(ctx, filter)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "fetching valid vessels")
	}

	res := &messages.GetValidVesselsResponse{
		Filter: request.Filter,
	}

	for _, y := range x {
		res.Results = append(res.Results, grpcconverters.ConvertValidVesselToGRPCValidVessel(y))
	}

	return res, nil
}

func (s *ServiceImpl) SearchForValidIngredientGroups(ctx context.Context, request *messages.SearchForValidIngredientGroupsRequest) (*messages.SearchForValidIngredientGroupsResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)
	filter := grpcconverters.ConvertGRPCQueryFilterToQueryFilter(request.Filter)
	tracing.AttachQueryFilterToSpan(span, filter)

	x, err := s.validEnumerationsManager.SearchValidIngredientGroups(ctx, request.Query, request.UseDatabase, filter)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "searching for valid ingredient groups")
	}

	res := &messages.SearchForValidIngredientGroupsResponse{
		Filter: request.Filter,
	}

	for _, y := range x {
		res.Results = append(res.Results, grpcconverters.ConvertValidIngredientGroupToGRPCValidIngredientGroup(y))
	}

	return res, nil
}

func (s *ServiceImpl) SearchForValidIngredientStates(ctx context.Context, request *messages.SearchForValidIngredientStatesRequest) (*messages.SearchForValidIngredientStatesResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)
	filter := grpcconverters.ConvertGRPCQueryFilterToQueryFilter(request.Filter)
	tracing.AttachQueryFilterToSpan(span, filter)

	x, err := s.validEnumerationsManager.SearchValidIngredientStates(ctx, request.Query, request.UseDatabase, filter)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "searching for valid ingredient states")
	}

	res := &messages.SearchForValidIngredientStatesResponse{
		Filter: request.Filter,
	}

	for _, y := range x {
		res.Results = append(res.Results, grpcconverters.ConvertValidIngredientStateToGRPCValidIngredientState(y))
	}

	return res, nil
}

func (s *ServiceImpl) SearchForValidIngredients(ctx context.Context, request *messages.SearchForValidIngredientsRequest) (*messages.SearchForValidIngredientsResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)
	filter := grpcconverters.ConvertGRPCQueryFilterToQueryFilter(request.Filter)
	tracing.AttachQueryFilterToSpan(span, filter)

	x, err := s.validEnumerationsManager.SearchValidIngredients(ctx, request.Query, request.UseDatabase, filter)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "searching for valid ingredients")
	}

	res := &messages.SearchForValidIngredientsResponse{
		Filter: request.Filter,
	}

	for _, y := range x {
		res.Results = append(res.Results, grpcconverters.ConvertValidIngredientToGRPCValidIngredient(y))
	}

	return res, nil
}

func (s *ServiceImpl) SearchForValidInstruments(ctx context.Context, request *messages.SearchForValidInstrumentsRequest) (*messages.SearchForValidInstrumentsResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)
	filter := grpcconverters.ConvertGRPCQueryFilterToQueryFilter(request.Filter)
	tracing.AttachQueryFilterToSpan(span, filter)

	x, err := s.validEnumerationsManager.SearchValidInstruments(ctx, request.Query, request.UseDatabase, filter)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "searching for valid instruments")
	}

	res := &messages.SearchForValidInstrumentsResponse{
		Filter: request.Filter,
	}

	for _, y := range x {
		res.Results = append(res.Results, grpcconverters.ConvertValidInstrumentToGRPCValidInstrument(y))
	}

	return res, nil
}

func (s *ServiceImpl) SearchForValidMeasurementUnits(ctx context.Context, request *messages.SearchForValidMeasurementUnitsRequest) (*messages.SearchForValidMeasurementUnitsResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)
	filter := grpcconverters.ConvertGRPCQueryFilterToQueryFilter(request.Filter)
	tracing.AttachQueryFilterToSpan(span, filter)

	x, err := s.validEnumerationsManager.SearchValidMeasurementUnits(ctx, request.Query, request.UseDatabase, filter)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "searching for valid measurement units")
	}

	res := &messages.SearchForValidMeasurementUnitsResponse{
		Filter: request.Filter,
	}

	for _, y := range x {
		res.Results = append(res.Results, grpcconverters.ConvertValidMeasurementUnitToGRPCValidMeasurementUnit(y))
	}

	return res, nil
}

func (s *ServiceImpl) SearchForValidPreparations(ctx context.Context, request *messages.SearchForValidPreparationsRequest) (*messages.SearchForValidPreparationsResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)
	filter := grpcconverters.ConvertGRPCQueryFilterToQueryFilter(request.Filter)
	tracing.AttachQueryFilterToSpan(span, filter)

	x, err := s.validEnumerationsManager.SearchValidPreparations(ctx, request.Query, request.UseDatabase, filter)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "searching for valid preparations")
	}

	res := &messages.SearchForValidPreparationsResponse{
		Filter: request.Filter,
	}

	for _, y := range x {
		res.Results = append(res.Results, grpcconverters.ConvertValidPreparationToGRPCValidPreparation(y))
	}

	return res, nil
}

func (s *ServiceImpl) SearchForValidVessels(ctx context.Context, request *messages.SearchForValidVesselsRequest) (*messages.SearchForValidVesselsResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)
	filter := grpcconverters.ConvertGRPCQueryFilterToQueryFilter(request.Filter)
	tracing.AttachQueryFilterToSpan(span, filter)

	x, err := s.validEnumerationsManager.SearchValidVessels(ctx, request.Query, request.UseDatabase, filter)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "searching for valid vessels")
	}

	res := &messages.SearchForValidVesselsResponse{
		Filter: request.Filter,
	}

	for _, y := range x {
		res.Results = append(res.Results, grpcconverters.ConvertValidVesselToGRPCValidVessel(y))
	}

	return res, nil
}

func (s *ServiceImpl) SearchValidIngredientsByPreparation(ctx context.Context, request *messages.SearchValidIngredientsByPreparationRequest) (*messages.SearchValidIngredientsByPreparationResponse, error) {
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

	res := &messages.SearchValidIngredientsByPreparationResponse{
		Filter: request.Filter,
	}

	for _, y := range x {
		res.Results = append(res.Results, grpcconverters.ConvertValidIngredientToGRPCValidIngredient(y))
	}

	return res, nil
}

func (s *ServiceImpl) SearchValidMeasurementUnitsByIngredient(ctx context.Context, request *messages.SearchValidMeasurementUnitsByIngredientRequest) (*messages.SearchValidMeasurementUnitsByIngredientResponse, error) {
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

	res := &messages.SearchValidMeasurementUnitsByIngredientResponse{
		Filter: request.Filter,
	}

	for _, y := range x {
		res.Results = append(res.Results, grpcconverters.ConvertValidMeasurementUnitToGRPCValidMeasurementUnit(y))
	}

	return res, nil
}

func (s *ServiceImpl) UpdateValidIngredient(ctx context.Context, request *messages.UpdateValidIngredientRequest) (*messages.UpdateValidIngredientResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValue(keys.ValidIngredientIDKey, request.ValidIngredientID)
	tracing.AttachToSpan(span, keys.ValidIngredientIDKey, request.ValidIngredientID)

	input := grpcconverters.ConvertGRPCValidIngredientUpdateRequestInputToValidIngredientUpdateRequestInput(request.Input)
	updated, err := s.validEnumerationsManager.UpdateValidIngredient(ctx, request.ValidIngredientID, input)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "updating valid ingredient")
	}

	res := &messages.UpdateValidIngredientResponse{
		Result: grpcconverters.ConvertValidIngredientToGRPCValidIngredient(updated),
	}

	return res, nil
}

func (s *ServiceImpl) UpdateValidIngredientGroup(ctx context.Context, request *messages.UpdateValidIngredientGroupRequest) (*messages.UpdateValidIngredientGroupResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValue(keys.ValidIngredientGroupIDKey, request.ValidIngredientGroupID)
	tracing.AttachToSpan(span, keys.ValidIngredientGroupIDKey, request.ValidIngredientGroupID)

	input := grpcconverters.ConvertGRPCValidIngredientGroupUpdateRequestInputToValidIngredientGroupUpdateRequestInput(request.Input)

	updated, err := s.validEnumerationsManager.UpdateValidIngredientGroup(ctx, request.ValidIngredientGroupID, input)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "updating valid ingredient group")
	}

	res := &messages.UpdateValidIngredientGroupResponse{
		Result: grpcconverters.ConvertValidIngredientGroupToGRPCValidIngredientGroup(updated),
	}

	return res, nil
}

func (s *ServiceImpl) UpdateValidIngredientMeasurementUnit(ctx context.Context, request *messages.UpdateValidIngredientMeasurementUnitRequest) (*messages.UpdateValidIngredientMeasurementUnitResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValue(keys.ValidIngredientMeasurementUnitIDKey, request.ValidIngredientMeasurementUnitID)
	tracing.AttachToSpan(span, keys.ValidIngredientMeasurementUnitIDKey, request.ValidIngredientMeasurementUnitID)

	input := grpcconverters.ConvertGRPCValidIngredientMeasurementUnitUpdateRequestInputToValidIngredientMeasurementUnitUpdateRequestInput(request.Input)

	updated, err := s.validEnumerationsManager.UpdateValidIngredientMeasurementUnit(ctx, request.ValidIngredientMeasurementUnitID, input)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "updating valid ingredient measurement unit")
	}

	res := &messages.UpdateValidIngredientMeasurementUnitResponse{
		Result: grpcconverters.ConvertValidIngredientMeasurementUnitToGRPCValidIngredientMeasurementUnit(updated),
	}

	return res, nil
}

func (s *ServiceImpl) UpdateValidIngredientPreparation(ctx context.Context, request *messages.UpdateValidIngredientPreparationRequest) (*messages.UpdateValidIngredientPreparationResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValue(keys.ValidIngredientPreparationIDKey, request.ValidIngredientPreparationID)
	tracing.AttachToSpan(span, keys.ValidIngredientPreparationIDKey, request.ValidIngredientPreparationID)

	input := grpcconverters.ConvertGRPCValidIngredientPreparationUpdateRequestInputToValidIngredientPreparationUpdateRequestInput(request.Input)

	updated, err := s.validEnumerationsManager.UpdateValidIngredientPreparation(ctx, request.ValidIngredientPreparationID, input)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "updating valid ingredient preparation")
	}

	res := &messages.UpdateValidIngredientPreparationResponse{
		Result: grpcconverters.ConvertValidIngredientPreparationToGRPCValidIngredientPreparation(updated),
	}

	return res, nil
}

func (s *ServiceImpl) UpdateValidIngredientState(ctx context.Context, request *messages.UpdateValidIngredientStateRequest) (*messages.UpdateValidIngredientStateResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValue(keys.ValidIngredientStateIDKey, request.ValidIngredientStateID)
	tracing.AttachToSpan(span, keys.ValidIngredientStateIDKey, request.ValidIngredientStateID)

	input := grpcconverters.ConvertGRPCValidIngredientStateUpdateRequestInputToValidIngredientStateUpdateRequestInput(request.Input)

	updated, err := s.validEnumerationsManager.UpdateValidIngredientState(ctx, request.ValidIngredientStateID, input)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "updating valid ingredient state")
	}

	res := &messages.UpdateValidIngredientStateResponse{
		Result: grpcconverters.ConvertValidIngredientStateToGRPCValidIngredientState(updated),
	}

	return res, nil
}

func (s *ServiceImpl) UpdateValidIngredientStateIngredient(ctx context.Context, request *messages.UpdateValidIngredientStateIngredientRequest) (*messages.UpdateValidIngredientStateIngredientResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValue(keys.ValidIngredientStateIngredientIDKey, request.ValidIngredientStateIngredientID)
	tracing.AttachToSpan(span, keys.ValidIngredientStateIngredientIDKey, request.ValidIngredientStateIngredientID)

	input := grpcconverters.ConvertGRPCValidIngredientStateIngredientUpdateRequestInputToValidIngredientStateIngredientUpdateRequestInput(request.Input)

	updated, err := s.validEnumerationsManager.UpdateValidIngredientStateIngredient(ctx, request.ValidIngredientStateIngredientID, input)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "updating valid ingredient state ingredient")
	}

	res := &messages.UpdateValidIngredientStateIngredientResponse{
		Result: grpcconverters.ConvertValidIngredientStateIngredientToGRPCValidIngredientStateIngredient(updated),
	}

	return res, nil
}

func (s *ServiceImpl) UpdateValidInstrument(ctx context.Context, request *messages.UpdateValidInstrumentRequest) (*messages.UpdateValidInstrumentResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValue(keys.ValidInstrumentIDKey, request.ValidInstrumentID)
	tracing.AttachToSpan(span, keys.ValidInstrumentIDKey, request.ValidInstrumentID)

	input := grpcconverters.ConvertGRPCValidInstrumentUpdateRequestInputToValidInstrumentUpdateRequestInput(request.Input)

	updated, err := s.validEnumerationsManager.UpdateValidInstrument(ctx, request.ValidInstrumentID, input)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "updating valid instrument")
	}

	res := &messages.UpdateValidInstrumentResponse{
		Result: grpcconverters.ConvertValidInstrumentToGRPCValidInstrument(updated),
	}

	return res, nil
}

func (s *ServiceImpl) UpdateValidMeasurementUnit(ctx context.Context, request *messages.UpdateValidMeasurementUnitRequest) (*messages.UpdateValidMeasurementUnitResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValue(keys.ValidMeasurementUnitIDKey, request.ValidMeasurementUnitID)
	tracing.AttachToSpan(span, keys.ValidMeasurementUnitIDKey, request.ValidMeasurementUnitID)

	input := grpcconverters.ConvertGRPCValidMeasurementUnitUpdateRequestInputToValidMeasurementUnitUpdateRequestInput(request.Input)

	updated, err := s.validEnumerationsManager.UpdateValidMeasurementUnit(ctx, request.ValidMeasurementUnitID, input)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "updating valid measurement unit")
	}

	res := &messages.UpdateValidMeasurementUnitResponse{
		Result: grpcconverters.ConvertValidMeasurementUnitToGRPCValidMeasurementUnit(updated),
	}

	return res, nil
}

func (s *ServiceImpl) UpdateValidMeasurementUnitConversion(ctx context.Context, request *messages.UpdateValidMeasurementUnitConversionRequest) (*messages.UpdateValidMeasurementUnitConversionResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValue(keys.ValidMeasurementUnitConversionIDKey, request.ValidMeasurementUnitConversionID)
	tracing.AttachToSpan(span, keys.ValidMeasurementUnitConversionIDKey, request.ValidMeasurementUnitConversionID)

	input := grpcconverters.ConvertGRPCValidMeasurementUnitConversionUpdateRequestInputToValidMeasurementUnitConversionUpdateRequestInput(request.Input)

	updated, err := s.validEnumerationsManager.UpdateValidMeasurementUnitConversion(ctx, request.ValidMeasurementUnitConversionID, input)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "updating valid measurement unit conversion")
	}

	res := &messages.UpdateValidMeasurementUnitConversionResponse{
		Result: grpcconverters.ConvertValidMeasurementUnitConversionToGRPCValidMeasurementUnitConversion(updated),
	}

	return res, nil
}

func (s *ServiceImpl) UpdateValidPreparation(ctx context.Context, request *messages.UpdateValidPreparationRequest) (*messages.UpdateValidPreparationResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValue(keys.ValidPreparationIDKey, request.ValidPreparationID)
	tracing.AttachToSpan(span, keys.ValidPreparationIDKey, request.ValidPreparationID)

	input := grpcconverters.ConvertGRPCValidPreparationUpdateRequestInputToValidPreparationUpdateRequestInput(request.Input)

	updated, err := s.validEnumerationsManager.UpdateValidPreparation(ctx, request.ValidPreparationID, input)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "updating valid preparation")
	}

	res := &messages.UpdateValidPreparationResponse{
		Result: grpcconverters.ConvertValidPreparationToGRPCValidPreparation(updated),
	}

	return res, nil
}

func (s *ServiceImpl) UpdateValidPreparationInstrument(ctx context.Context, request *messages.UpdateValidPreparationInstrumentRequest) (*messages.UpdateValidPreparationInstrumentResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValue(keys.ValidPreparationInstrumentIDKey, request.ValidPreparationInstrumentID)
	tracing.AttachToSpan(span, keys.ValidPreparationInstrumentIDKey, request.ValidPreparationInstrumentID)

	input := grpcconverters.ConvertGRPCValidPreparationInstrumentUpdateRequestInputToValidPreparationInstrumentUpdateRequestInput(request.Input)

	updated, err := s.validEnumerationsManager.UpdateValidPreparationInstrument(ctx, request.ValidPreparationInstrumentID, input)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "updating valid preparation instrument")
	}

	res := &messages.UpdateValidPreparationInstrumentResponse{
		Result: grpcconverters.ConvertValidPreparationInstrumentToGRPCValidPreparationInstrument(updated),
	}

	return res, nil
}

func (s *ServiceImpl) UpdateValidPreparationVessel(ctx context.Context, request *messages.UpdateValidPreparationVesselRequest) (*messages.UpdateValidPreparationVesselResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValue(keys.ValidPreparationVesselIDKey, request.ValidPreparationVesselID)
	tracing.AttachToSpan(span, keys.ValidPreparationVesselIDKey, request.ValidPreparationVesselID)

	input := grpcconverters.ConvertGRPCValidPreparationVesselUpdateRequestInputToValidPreparationVesselUpdateRequestInput(request.Input)

	updated, err := s.validEnumerationsManager.UpdateValidPreparationVessel(ctx, request.ValidPreparationVesselID, input)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "updating valid preparation vessel")
	}

	res := &messages.UpdateValidPreparationVesselResponse{
		Result: grpcconverters.ConvertValidPreparationVesselToGRPCValidPreparationVessel(updated),
	}

	return res, nil
}

func (s *ServiceImpl) UpdateValidVessel(ctx context.Context, request *messages.UpdateValidVesselRequest) (*messages.UpdateValidVesselResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValue(keys.ValidVesselIDKey, request.ValidVesselID)
	tracing.AttachToSpan(span, keys.ValidVesselIDKey, request.ValidVesselID)

	input := grpcconverters.ConvertGRPCValidVesselUpdateRequestInputToValidVesselUpdateRequestInput(request.Input)

	updated, err := s.validEnumerationsManager.UpdateValidVessel(ctx, request.ValidVesselID, input)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "updating valid vessel")
	}

	res := &messages.UpdateValidVesselResponse{
		Result: grpcconverters.ConvertValidVesselToGRPCValidVessel(updated),
	}

	return res, nil
}
