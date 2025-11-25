package grpc

import (
	"context"

	grpcconverters "github.com/dinnerdonebetter/backend/internal/grpc/converters"
	"github.com/dinnerdonebetter/backend/internal/grpc/generated/services/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/platform/observability"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
	mealplanningconverters "github.com/dinnerdonebetter/backend/internal/services/mealplanning/grpc/converters"

	"google.golang.org/grpc/codes"
)

func (s *serviceImpl) ArchiveValidIngredient(ctx context.Context, request *mealplanning.ArchiveValidIngredientRequest) (*mealplanning.ArchiveValidIngredientResponse, error) {
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

func (s *serviceImpl) ArchiveValidIngredientGroup(ctx context.Context, request *mealplanning.ArchiveValidIngredientGroupRequest) (*mealplanning.ArchiveValidIngredientGroupResponse, error) {
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

func (s *serviceImpl) ArchiveValidIngredientMeasurementUnit(ctx context.Context, request *mealplanning.ArchiveValidIngredientMeasurementUnitRequest) (*mealplanning.ArchiveValidIngredientMeasurementUnitResponse, error) {
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

func (s *serviceImpl) ArchiveValidIngredientPreparation(ctx context.Context, request *mealplanning.ArchiveValidIngredientPreparationRequest) (*mealplanning.ArchiveValidIngredientPreparationResponse, error) {
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

func (s *serviceImpl) ArchiveValidIngredientState(ctx context.Context, request *mealplanning.ArchiveValidIngredientStateRequest) (*mealplanning.ArchiveValidIngredientStateResponse, error) {
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

func (s *serviceImpl) ArchiveValidIngredientStateIngredient(ctx context.Context, request *mealplanning.ArchiveValidIngredientStateIngredientRequest) (*mealplanning.ArchiveValidIngredientStateIngredientResponse, error) {
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

func (s *serviceImpl) ArchiveValidInstrument(ctx context.Context, request *mealplanning.ArchiveValidInstrumentRequest) (*mealplanning.ArchiveValidInstrumentResponse, error) {
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

func (s *serviceImpl) ArchiveValidMeasurementUnit(ctx context.Context, request *mealplanning.ArchiveValidMeasurementUnitRequest) (*mealplanning.ArchiveValidMeasurementUnitResponse, error) {
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

func (s *serviceImpl) ArchiveValidMeasurementUnitConversion(ctx context.Context, request *mealplanning.ArchiveValidMeasurementUnitConversionRequest) (*mealplanning.ArchiveValidMeasurementUnitConversionResponse, error) {
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

func (s *serviceImpl) ArchiveValidPreparation(ctx context.Context, request *mealplanning.ArchiveValidPreparationRequest) (*mealplanning.ArchiveValidPreparationResponse, error) {
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

func (s *serviceImpl) ArchiveValidPreparationInstrument(ctx context.Context, request *mealplanning.ArchiveValidPreparationInstrumentRequest) (*mealplanning.ArchiveValidPreparationInstrumentResponse, error) {
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

func (s *serviceImpl) ArchiveValidPreparationVessel(ctx context.Context, request *mealplanning.ArchiveValidPreparationVesselRequest) (*mealplanning.ArchiveValidPreparationVesselResponse, error) {
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

func (s *serviceImpl) ArchiveValidVessel(ctx context.Context, request *mealplanning.ArchiveValidVesselRequest) (*mealplanning.ArchiveValidVesselResponse, error) {
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

func (s *serviceImpl) CreateValidIngredient(ctx context.Context, request *mealplanning.CreateValidIngredientRequest) (*mealplanning.CreateValidIngredientResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	created, err := s.validEnumerationsManager.CreateValidIngredient(ctx, mealplanningconverters.ConvertGRPCCreateValidIngredientRequestToValidIngredientCreationRequestInput(request.Input))
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "creating valid ingredient")
	}
	tracing.AttachToSpan(span, keys.ValidIngredientIDKey, created.ID)

	result := &mealplanning.CreateValidIngredientResponse{
		Result: mealplanningconverters.ConvertValidIngredientToGRPCValidIngredient(created),
	}

	return result, nil
}

func (s *serviceImpl) CreateValidIngredientGroup(ctx context.Context, request *mealplanning.CreateValidIngredientGroupRequest) (*mealplanning.CreateValidIngredientGroupResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	created, err := s.validEnumerationsManager.CreateValidIngredientGroup(ctx, mealplanningconverters.ConvertGRPCValidIngredientGroupCreationRequestInputToValidIngredientGroupCreationRequestInput(request.Input))
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "creating valid ingredient")
	}
	tracing.AttachToSpan(span, keys.ValidIngredientGroupIDKey, created.ID)

	result := &mealplanning.CreateValidIngredientGroupResponse{
		Result: mealplanningconverters.ConvertValidIngredientGroupToGRPCValidIngredientGroup(created),
	}

	return result, nil
}

func (s *serviceImpl) CreateValidIngredientMeasurementUnit(ctx context.Context, request *mealplanning.CreateValidIngredientMeasurementUnitRequest) (*mealplanning.CreateValidIngredientMeasurementUnitResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	created, err := s.validEnumerationsManager.CreateValidIngredientMeasurementUnit(ctx, mealplanningconverters.ConvertGRPCCreateValidIngredientMeasurementUnitRequestToValidIngredientMeasurementUnitCreationRequestInput(request.Input))
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "creating valid ingredient measurement unit")
	}
	tracing.AttachToSpan(span, keys.ValidIngredientMeasurementUnitIDKey, created.ID)

	result := &mealplanning.CreateValidIngredientMeasurementUnitResponse{
		Result: mealplanningconverters.ConvertValidIngredientMeasurementUnitToGRPCValidIngredientMeasurementUnit(created),
	}

	return result, nil
}

func (s *serviceImpl) CreateValidIngredientPreparation(ctx context.Context, request *mealplanning.CreateValidIngredientPreparationRequest) (*mealplanning.CreateValidIngredientPreparationResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	created, err := s.validEnumerationsManager.CreateValidIngredientPreparation(ctx, mealplanningconverters.ConvertGRPCCreateValidIngredientPreparationRequestToValidIngredientPreparationCreationRequestInput(request.Input))
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "creating valid ingredient preparation")
	}
	tracing.AttachToSpan(span, keys.ValidIngredientPreparationIDKey, created.ID)

	result := &mealplanning.CreateValidIngredientPreparationResponse{
		Result: mealplanningconverters.ConvertValidIngredientPreparationToGRPCValidIngredientPreparation(created),
	}

	return result, nil
}

func (s *serviceImpl) CreateValidIngredientState(ctx context.Context, request *mealplanning.CreateValidIngredientStateRequest) (*mealplanning.CreateValidIngredientStateResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	created, err := s.validEnumerationsManager.CreateValidIngredientState(ctx, mealplanningconverters.ConvertGRPCCreateValidIngredientStateRequestToValidIngredientStateCreationRequestInput(request.Input))
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "creating valid ingredient state")
	}
	tracing.AttachToSpan(span, keys.ValidIngredientStateIDKey, created.ID)

	result := &mealplanning.CreateValidIngredientStateResponse{
		Result: mealplanningconverters.ConvertValidIngredientStateToGRPCValidIngredientState(created),
	}

	return result, nil
}

func (s *serviceImpl) CreateValidIngredientStateIngredient(ctx context.Context, request *mealplanning.CreateValidIngredientStateIngredientRequest) (*mealplanning.CreateValidIngredientStateIngredientResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	created, err := s.validEnumerationsManager.CreateValidIngredientStateIngredient(ctx, mealplanningconverters.ConvertGRPCCreateValidIngredientStateIngredientRequestToValidIngredientStateIngredientCreationRequestInput(request.Input))
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "creating valid ingredient state ingredient")
	}
	tracing.AttachToSpan(span, keys.ValidIngredientStateIngredientIDKey, created.ID)

	result := &mealplanning.CreateValidIngredientStateIngredientResponse{
		Result: mealplanningconverters.ConvertValidIngredientStateIngredientToGRPCValidIngredientStateIngredient(created),
	}

	return result, nil
}

func (s *serviceImpl) CreateValidInstrument(ctx context.Context, request *mealplanning.CreateValidInstrumentRequest) (*mealplanning.CreateValidInstrumentResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	created, err := s.validEnumerationsManager.CreateValidInstrument(ctx, mealplanningconverters.ConvertGRPCCreateValidInstrumentRequestToValidInstrumentCreationRequestInput(request.Input))
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "creating valid instrument")
	}
	tracing.AttachToSpan(span, keys.ValidInstrumentIDKey, created.ID)

	result := &mealplanning.CreateValidInstrumentResponse{
		Result: mealplanningconverters.ConvertValidInstrumentToGRPCValidInstrument(created),
	}

	return result, nil
}

func (s *serviceImpl) CreateValidMeasurementUnit(ctx context.Context, request *mealplanning.CreateValidMeasurementUnitRequest) (*mealplanning.CreateValidMeasurementUnitResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	created, err := s.validEnumerationsManager.CreateValidMeasurementUnit(ctx, mealplanningconverters.ConvertGRPCValidMeasurementUnitCreationRequestInputToValidMeasurementUnitCreationRequestInput(request.Input))
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "creating valid measurement unit")
	}
	tracing.AttachToSpan(span, keys.ValidMeasurementUnitIDKey, created.ID)

	result := &mealplanning.CreateValidMeasurementUnitResponse{
		Result: mealplanningconverters.ConvertValidMeasurementUnitToGRPCValidMeasurementUnit(created),
	}

	return result, nil
}

func (s *serviceImpl) CreateValidMeasurementUnitConversion(ctx context.Context, request *mealplanning.CreateValidMeasurementUnitConversionRequest) (*mealplanning.CreateValidMeasurementUnitConversionResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	created, err := s.validEnumerationsManager.CreateValidMeasurementUnitConversion(ctx, mealplanningconverters.ConvertGRPCCreateValidMeasurementUnitConversionRequestToValidMeasurementUnitConversionCreationRequestInput(request.Input))
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "creating valid measurement unit conversion")
	}
	tracing.AttachToSpan(span, keys.ValidMeasurementUnitConversionIDKey, created.ID)

	result := &mealplanning.CreateValidMeasurementUnitConversionResponse{
		Result: mealplanningconverters.ConvertValidMeasurementUnitConversionToGRPCValidMeasurementUnitConversion(created),
	}

	return result, nil
}

func (s *serviceImpl) CreateValidPreparation(ctx context.Context, request *mealplanning.CreateValidPreparationRequest) (*mealplanning.CreateValidPreparationResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	created, err := s.validEnumerationsManager.CreateValidPreparation(ctx, mealplanningconverters.ConvertGRPCValidPreparationCreationRequestInputToValidPreparationCreationRequestInput(request.Input))
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "creating valid preparation")
	}
	tracing.AttachToSpan(span, keys.ValidPreparationIDKey, created.ID)

	result := &mealplanning.CreateValidPreparationResponse{
		Result: mealplanningconverters.ConvertValidPreparationToGRPCValidPreparation(created),
	}

	return result, nil
}

func (s *serviceImpl) CreateValidPreparationInstrument(ctx context.Context, request *mealplanning.CreateValidPreparationInstrumentRequest) (*mealplanning.CreateValidPreparationInstrumentResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	created, err := s.validEnumerationsManager.CreateValidPreparationInstrument(ctx, mealplanningconverters.ConvertGRPCCreateValidPreparationInstrumentRequestToValidPreparationInstrumentCreationRequestInput(request.Input))
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "creating valid preparation instrument")
	}
	tracing.AttachToSpan(span, keys.ValidPreparationInstrumentIDKey, created.ID)

	result := &mealplanning.CreateValidPreparationInstrumentResponse{
		Result: mealplanningconverters.ConvertValidPreparationInstrumentToGRPCValidPreparationInstrument(created),
	}

	return result, nil
}

func (s *serviceImpl) CreateValidPreparationVessel(ctx context.Context, request *mealplanning.CreateValidPreparationVesselRequest) (*mealplanning.CreateValidPreparationVesselResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	created, err := s.validEnumerationsManager.CreateValidPreparationVessel(ctx, mealplanningconverters.ConvertGRPCCreateValidPreparationVesselRequestToValidPreparationVesselCreationRequestInput(request.Input))
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "creating valid preparation vessel")
	}
	tracing.AttachToSpan(span, keys.ValidPreparationVesselIDKey, created.ID)

	result := &mealplanning.CreateValidPreparationVesselResponse{
		Result: mealplanningconverters.ConvertValidPreparationVesselToGRPCValidPreparationVessel(created),
	}

	return result, nil
}

func (s *serviceImpl) CreateValidVessel(ctx context.Context, request *mealplanning.CreateValidVesselRequest) (*mealplanning.CreateValidVesselResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	created, err := s.validEnumerationsManager.CreateValidVessel(ctx, mealplanningconverters.ConvertGRPCValidVesselCreationRequestInputToValidVesselCreationRequestInput(request.Input))
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "creating valid vessel")
	}
	tracing.AttachToSpan(span, keys.ValidVesselIDKey, created.ID)

	result := &mealplanning.CreateValidVesselResponse{
		Result: mealplanningconverters.ConvertValidVesselToGRPCValidVessel(created),
	}

	return result, nil
}

func (s *serviceImpl) GetRandomValidIngredient(ctx context.Context, _ *mealplanning.GetRandomValidIngredientRequest) (*mealplanning.GetRandomValidIngredientResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	selected, err := s.validEnumerationsManager.RandomValidIngredient(ctx)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "fetching random valid ingredient")
	}
	tracing.AttachToSpan(span, keys.ValidIngredientIDKey, selected.ID)

	res := &mealplanning.GetRandomValidIngredientResponse{
		Result: mealplanningconverters.ConvertValidIngredientToGRPCValidIngredient(selected),
	}

	return res, nil
}

func (s *serviceImpl) GetRandomValidInstrument(ctx context.Context, _ *mealplanning.GetRandomValidInstrumentRequest) (*mealplanning.GetRandomValidInstrumentResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	selected, err := s.validEnumerationsManager.RandomValidInstrument(ctx)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "fetching random valid instrument")
	}
	tracing.AttachToSpan(span, keys.ValidInstrumentIDKey, selected.ID)

	res := &mealplanning.GetRandomValidInstrumentResponse{
		Result: mealplanningconverters.ConvertValidInstrumentToGRPCValidInstrument(selected),
	}

	return res, nil
}

func (s *serviceImpl) GetRandomValidPreparation(ctx context.Context, _ *mealplanning.GetRandomValidPreparationRequest) (*mealplanning.GetRandomValidPreparationResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	selected, err := s.validEnumerationsManager.RandomValidPreparation(ctx)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "fetching random valid preparation")
	}
	tracing.AttachToSpan(span, keys.ValidPreparationIDKey, selected.ID)

	res := &mealplanning.GetRandomValidPreparationResponse{
		Result: mealplanningconverters.ConvertValidPreparationToGRPCValidPreparation(selected),
	}

	return res, nil
}

func (s *serviceImpl) GetRandomValidVessel(ctx context.Context, _ *mealplanning.GetRandomValidVesselRequest) (*mealplanning.GetRandomValidVesselResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	selected, err := s.validEnumerationsManager.RandomValidVessel(ctx)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "fetching random valid vessel")
	}
	tracing.AttachToSpan(span, keys.ValidVesselIDKey, selected.ID)

	res := &mealplanning.GetRandomValidVesselResponse{
		Result: mealplanningconverters.ConvertValidVesselToGRPCValidVessel(selected),
	}

	return res, nil
}

func (s *serviceImpl) GetValidIngredient(ctx context.Context, request *mealplanning.GetValidIngredientRequest) (*mealplanning.GetValidIngredientResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValue(keys.ValidIngredientIDKey, request.ValidIngredientID)
	tracing.AttachToSpan(span, keys.ValidIngredientIDKey, request.ValidIngredientID)

	x, err := s.validEnumerationsManager.ReadValidIngredient(ctx, request.ValidIngredientID)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "fetching valid ingredient")
	}

	res := &mealplanning.GetValidIngredientResponse{
		Result: mealplanningconverters.ConvertValidIngredientToGRPCValidIngredient(x),
	}

	return res, nil
}

func (s *serviceImpl) GetValidIngredientGroup(ctx context.Context, request *mealplanning.GetValidIngredientGroupRequest) (*mealplanning.GetValidIngredientGroupResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValue(keys.ValidIngredientGroupIDKey, request.ValidIngredientGroupID)
	tracing.AttachToSpan(span, keys.ValidIngredientGroupIDKey, request.ValidIngredientGroupID)

	x, err := s.validEnumerationsManager.ReadValidIngredientGroup(ctx, request.ValidIngredientGroupID)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "fetching valid ingredient group")
	}

	res := &mealplanning.GetValidIngredientGroupResponse{
		Result: mealplanningconverters.ConvertValidIngredientGroupToGRPCValidIngredientGroup(x),
	}

	return res, nil
}

func (s *serviceImpl) GetValidIngredientGroups(ctx context.Context, request *mealplanning.GetValidIngredientGroupsRequest) (*mealplanning.GetValidIngredientGroupsResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)
	filter := grpcconverters.ConvertGRPCQueryFilterToQueryFilter(request.Filter)
	tracing.AttachQueryFilterToSpan(span, filter)

	x, err := s.validEnumerationsManager.ListValidIngredientGroups(ctx, filter)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "fetching valid ingredient groups")
	}

	res := &mealplanning.GetValidIngredientGroupsResponse{
		Filter: grpcconverters.ConvertQueryFilterToGRPCQueryFilter(filter, x.Pagination),
	}
	for _, y := range x.Data {
		res.Results = append(res.Results, mealplanningconverters.ConvertValidIngredientGroupToGRPCValidIngredientGroup(y))
	}

	return res, nil
}

func (s *serviceImpl) GetValidIngredientMeasurementUnit(ctx context.Context, request *mealplanning.GetValidIngredientMeasurementUnitRequest) (*mealplanning.GetValidIngredientMeasurementUnitResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValue(keys.ValidIngredientMeasurementUnitIDKey, request.ValidIngredientMeasurementUnitID)
	tracing.AttachToSpan(span, keys.ValidIngredientMeasurementUnitIDKey, request.ValidIngredientMeasurementUnitID)

	x, err := s.validEnumerationsManager.ReadValidIngredientMeasurementUnit(ctx, request.ValidIngredientMeasurementUnitID)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "fetching valid ingredient measurement unit")
	}

	res := &mealplanning.GetValidIngredientMeasurementUnitResponse{
		Result: mealplanningconverters.ConvertValidIngredientMeasurementUnitToGRPCValidIngredientMeasurementUnit(x),
	}

	return res, nil
}

func (s *serviceImpl) GetValidIngredientMeasurementUnits(ctx context.Context, request *mealplanning.GetValidIngredientMeasurementUnitsRequest) (*mealplanning.GetValidIngredientMeasurementUnitsResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)
	filter := grpcconverters.ConvertGRPCQueryFilterToQueryFilter(request.Filter)
	tracing.AttachQueryFilterToSpan(span, filter)

	x, err := s.validEnumerationsManager.ListValidIngredientMeasurementUnits(ctx, filter)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "fetching valid ingredient measurement units")
	}

	res := &mealplanning.GetValidIngredientMeasurementUnitsResponse{
		Filter: grpcconverters.ConvertQueryFilterToGRPCQueryFilter(filter, x.Pagination),
	}
	for _, y := range x.Data {
		res.Results = append(res.Results, mealplanningconverters.ConvertValidIngredientMeasurementUnitToGRPCValidIngredientMeasurementUnit(y))
	}

	return res, nil
}

func (s *serviceImpl) GetValidIngredientMeasurementUnitsByIngredient(ctx context.Context, request *mealplanning.GetValidIngredientMeasurementUnitsByIngredientRequest) (*mealplanning.GetValidIngredientMeasurementUnitsByIngredientResponse, error) {
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
		Filter: request.Filter, // TODO: grpcconverters.ConvertQueryFilterToGRPCQueryFilter(filter, x.Pagination),
	}

	for _, y := range x {
		res.Results = append(res.Results, mealplanningconverters.ConvertValidIngredientMeasurementUnitToGRPCValidIngredientMeasurementUnit(y))
	}

	return res, nil
}

func (s *serviceImpl) GetValidIngredientMeasurementUnitsByMeasurementUnit(ctx context.Context, request *mealplanning.GetValidIngredientMeasurementUnitsByMeasurementUnitRequest) (*mealplanning.GetValidIngredientMeasurementUnitsByMeasurementUnitResponse, error) {
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
		Filter: request.Filter, // TODO: grpcconverters.ConvertQueryFilterToGRPCQueryFilter(filter, x.Pagination),
	}

	for _, y := range x {
		res.Results = append(res.Results, mealplanningconverters.ConvertValidIngredientMeasurementUnitToGRPCValidIngredientMeasurementUnit(y))
	}

	return res, nil
}

func (s *serviceImpl) GetValidIngredientPreparation(ctx context.Context, request *mealplanning.GetValidIngredientPreparationRequest) (*mealplanning.GetValidIngredientPreparationResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValue(keys.ValidIngredientPreparationIDKey, request.ValidIngredientPreparationID)
	tracing.AttachToSpan(span, keys.ValidIngredientPreparationIDKey, request.ValidIngredientPreparationID)

	x, err := s.validEnumerationsManager.ReadValidIngredientPreparation(ctx, request.ValidIngredientPreparationID)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "fetching valid ingredient preparation")
	}

	res := &mealplanning.GetValidIngredientPreparationResponse{
		Result: mealplanningconverters.ConvertValidIngredientPreparationToGRPCValidIngredientPreparation(x),
	}

	return res, nil
}

func (s *serviceImpl) GetValidIngredientPreparations(ctx context.Context, request *mealplanning.GetValidIngredientPreparationsRequest) (*mealplanning.GetValidIngredientPreparationsResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)
	filter := grpcconverters.ConvertGRPCQueryFilterToQueryFilter(request.Filter)
	tracing.AttachQueryFilterToSpan(span, filter)

	x, err := s.validEnumerationsManager.ListValidIngredientPreparations(ctx, filter)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "fetching valid ingredient preparations")
	}

	res := &mealplanning.GetValidIngredientPreparationsResponse{
		Filter: grpcconverters.ConvertQueryFilterToGRPCQueryFilter(filter, x.Pagination),
	}

	for _, y := range x.Data {
		res.Results = append(res.Results, mealplanningconverters.ConvertValidIngredientPreparationToGRPCValidIngredientPreparation(y))
	}

	return res, nil
}

func (s *serviceImpl) GetValidIngredientPreparationsByIngredient(ctx context.Context, request *mealplanning.GetValidIngredientPreparationsByIngredientRequest) (*mealplanning.GetValidIngredientPreparationsByIngredientResponse, error) {
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
		Filter: grpcconverters.ConvertQueryFilterToGRPCQueryFilter(filter, x.Pagination),
	}

	for _, y := range x.Data {
		res.Results = append(res.Results, mealplanningconverters.ConvertValidIngredientPreparationToGRPCValidIngredientPreparation(y))
	}

	return res, nil
}

func (s *serviceImpl) GetValidIngredientPreparationsByPreparation(ctx context.Context, request *mealplanning.GetValidIngredientPreparationsByPreparationRequest) (*mealplanning.GetValidIngredientPreparationsByPreparationResponse, error) {
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
		Filter: grpcconverters.ConvertQueryFilterToGRPCQueryFilter(filter, x.Pagination),
	}

	for _, y := range x.Data {
		res.Results = append(res.Results, mealplanningconverters.ConvertValidIngredientPreparationToGRPCValidIngredientPreparation(y))
	}

	return res, nil
}

func (s *serviceImpl) GetValidIngredientState(ctx context.Context, request *mealplanning.GetValidIngredientStateRequest) (*mealplanning.GetValidIngredientStateResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValue(keys.ValidIngredientStateIDKey, request.ValidIngredientStateID)
	tracing.AttachToSpan(span, keys.ValidIngredientStateIDKey, request.ValidIngredientStateID)

	x, err := s.validEnumerationsManager.ReadValidIngredientState(ctx, request.ValidIngredientStateID)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "fetching valid ingredient state")
	}

	res := &mealplanning.GetValidIngredientStateResponse{
		Result: mealplanningconverters.ConvertValidIngredientStateToGRPCValidIngredientState(x),
	}

	return res, nil
}

func (s *serviceImpl) GetValidIngredientStateIngredient(ctx context.Context, request *mealplanning.GetValidIngredientStateIngredientRequest) (*mealplanning.GetValidIngredientStateIngredientResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValue(keys.ValidIngredientStateIngredientIDKey, request.ValidIngredientStateIngredientID)
	tracing.AttachToSpan(span, keys.ValidIngredientStateIngredientIDKey, request.ValidIngredientStateIngredientID)

	x, err := s.validEnumerationsManager.ReadValidIngredientStateIngredient(ctx, request.ValidIngredientStateIngredientID)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "fetching valid ingredient state ingredient")
	}

	res := &mealplanning.GetValidIngredientStateIngredientResponse{
		Result: mealplanningconverters.ConvertValidIngredientStateIngredientToGRPCValidIngredientStateIngredient(x),
	}

	return res, nil
}

func (s *serviceImpl) GetValidIngredientStateIngredients(ctx context.Context, request *mealplanning.GetValidIngredientStateIngredientsRequest) (*mealplanning.GetValidIngredientStateIngredientsResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)
	filter := grpcconverters.ConvertGRPCQueryFilterToQueryFilter(request.Filter)
	tracing.AttachQueryFilterToSpan(span, filter)

	x, err := s.validEnumerationsManager.ListValidIngredientStateIngredients(ctx, filter)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "fetching valid ingredient state ingredients")
	}

	res := &mealplanning.GetValidIngredientStateIngredientsResponse{
		Filter: grpcconverters.ConvertQueryFilterToGRPCQueryFilter(filter, x.Pagination),
	}

	for _, y := range x.Data {
		res.Results = append(res.Results, mealplanningconverters.ConvertValidIngredientStateIngredientToGRPCValidIngredientStateIngredient(y))
	}

	return res, nil
}

func (s *serviceImpl) GetValidIngredientStateIngredientsByIngredient(ctx context.Context, request *mealplanning.GetValidIngredientStateIngredientsByIngredientRequest) (*mealplanning.GetValidIngredientStateIngredientsByIngredientResponse, error) {
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
		Filter: request.Filter, // TODO: grpcconverters.ConvertQueryFilterToGRPCQueryFilter(filter, x.Pagination),
	}

	for _, y := range x {
		res.Results = append(res.Results, mealplanningconverters.ConvertValidIngredientStateIngredientToGRPCValidIngredientStateIngredient(y))
	}

	return res, nil
}

func (s *serviceImpl) GetValidIngredientStateIngredientsByIngredientState(ctx context.Context, request *mealplanning.GetValidIngredientStateIngredientsByIngredientStateRequest) (*mealplanning.GetValidIngredientStateIngredientsByIngredientStateResponse, error) {
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
		Filter: request.Filter, // TODO: grpcconverters.ConvertQueryFilterToGRPCQueryFilter(filter, x.Pagination),
	}

	for _, y := range x {
		res.Results = append(res.Results, mealplanningconverters.ConvertValidIngredientStateIngredientToGRPCValidIngredientStateIngredient(y))
	}

	return res, nil
}

func (s *serviceImpl) GetValidIngredientStates(ctx context.Context, request *mealplanning.GetValidIngredientStatesRequest) (*mealplanning.GetValidIngredientStatesResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)
	filter := grpcconverters.ConvertGRPCQueryFilterToQueryFilter(request.Filter)
	tracing.AttachQueryFilterToSpan(span, filter)

	x, err := s.validEnumerationsManager.ListValidIngredientStates(ctx, filter)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "fetching valid ingredient states")
	}

	res := &mealplanning.GetValidIngredientStatesResponse{
		Filter: grpcconverters.ConvertQueryFilterToGRPCQueryFilter(filter, x.Pagination),
	}

	for _, y := range x.Data {
		res.Results = append(res.Results, mealplanningconverters.ConvertValidIngredientStateToGRPCValidIngredientState(y))
	}

	return res, nil
}

func (s *serviceImpl) GetValidIngredients(ctx context.Context, request *mealplanning.GetValidIngredientsRequest) (*mealplanning.GetValidIngredientsResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)
	filter := grpcconverters.ConvertGRPCQueryFilterToQueryFilter(request.Filter)
	tracing.AttachQueryFilterToSpan(span, filter)

	x, err := s.validEnumerationsManager.ListValidIngredients(ctx, filter)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "fetching valid ingredients")
	}

	logger.WithValue("pagination", x.Pagination).Info("Valid ingredients retrieved")

	res := &mealplanning.GetValidIngredientsResponse{
		Filter: grpcconverters.ConvertQueryFilterToGRPCQueryFilter(filter, x.Pagination),
	}

	for _, y := range x.Data {
		res.Results = append(res.Results, mealplanningconverters.ConvertValidIngredientToGRPCValidIngredient(y))
	}

	return res, nil
}

func (s *serviceImpl) GetValidInstrument(ctx context.Context, request *mealplanning.GetValidInstrumentRequest) (*mealplanning.GetValidInstrumentResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValue(keys.ValidInstrumentIDKey, request.ValidInstrumentID)
	tracing.AttachToSpan(span, keys.ValidInstrumentIDKey, request.ValidInstrumentID)

	x, err := s.validEnumerationsManager.ReadValidInstrument(ctx, request.ValidInstrumentID)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "fetching valid instrument")
	}

	res := &mealplanning.GetValidInstrumentResponse{
		Result: mealplanningconverters.ConvertValidInstrumentToGRPCValidInstrument(x),
	}

	return res, nil
}

func (s *serviceImpl) GetValidInstruments(ctx context.Context, request *mealplanning.GetValidInstrumentsRequest) (*mealplanning.GetValidInstrumentsResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)
	filter := grpcconverters.ConvertGRPCQueryFilterToQueryFilter(request.Filter)
	tracing.AttachQueryFilterToSpan(span, filter)

	x, err := s.validEnumerationsManager.ListValidInstruments(ctx, filter)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "fetching valid instruments")
	}

	res := &mealplanning.GetValidInstrumentsResponse{
		Filter: grpcconverters.ConvertQueryFilterToGRPCQueryFilter(filter, x.Pagination),
	}

	for _, y := range x.Data {
		res.Results = append(res.Results, mealplanningconverters.ConvertValidInstrumentToGRPCValidInstrument(y))
	}

	return res, nil
}

func (s *serviceImpl) GetValidMeasurementUnit(ctx context.Context, request *mealplanning.GetValidMeasurementUnitRequest) (*mealplanning.GetValidMeasurementUnitResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValue(keys.ValidMeasurementUnitIDKey, request.ValidMeasurementUnitID)
	tracing.AttachToSpan(span, keys.ValidMeasurementUnitIDKey, request.ValidMeasurementUnitID)

	x, err := s.validEnumerationsManager.ReadValidMeasurementUnit(ctx, request.ValidMeasurementUnitID)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "fetching valid measurement unit")
	}

	res := &mealplanning.GetValidMeasurementUnitResponse{
		Result: mealplanningconverters.ConvertValidMeasurementUnitToGRPCValidMeasurementUnit(x),
	}

	return res, nil
}

func (s *serviceImpl) GetValidMeasurementUnitConversion(ctx context.Context, request *mealplanning.GetValidMeasurementUnitConversionRequest) (*mealplanning.GetValidMeasurementUnitConversionResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValue(keys.ValidMeasurementUnitConversionIDKey, request.ValidMeasurementUnitConversionID)
	tracing.AttachToSpan(span, keys.ValidMeasurementUnitConversionIDKey, request.ValidMeasurementUnitConversionID)

	x, err := s.validEnumerationsManager.ReadValidMeasurementUnitConversion(ctx, request.ValidMeasurementUnitConversionID)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "fetching valid measurement unit conversion")
	}

	res := &mealplanning.GetValidMeasurementUnitConversionResponse{
		Result: mealplanningconverters.ConvertValidMeasurementUnitConversionToGRPCValidMeasurementUnitConversion(x),
	}

	return res, nil
}

func (s *serviceImpl) GetValidMeasurementUnitConversionsForUnit(ctx context.Context, request *mealplanning.GetValidMeasurementUnitConversionsForUnitRequest) (*mealplanning.GetValidMeasurementUnitConversionsForUnitResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValue(keys.ValidMeasurementUnitIDKey, request.ValidMeasurementUnitID)
	tracing.AttachToSpan(span, keys.ValidMeasurementUnitIDKey, request.ValidMeasurementUnitID)

	filter := grpcconverters.ConvertGRPCQueryFilterToQueryFilter(request.Filter)
	tracing.AttachQueryFilterToSpan(span, filter)

	x, err := s.validEnumerationsManager.ValidMeasurementUnitConversionsForMeasurementUnit(ctx, request.ValidMeasurementUnitID, filter)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "fetching valid measurement unit conversions to unit")
	}

	// TODO: add filter to response
	res := &mealplanning.GetValidMeasurementUnitConversionsForUnitResponse{}
	for _, y := range x.Data {
		res.Results = append(res.Results, mealplanningconverters.ConvertValidMeasurementUnitConversionToGRPCValidMeasurementUnitConversion(y))
	}

	return res, nil
}

func (s *serviceImpl) GetValidMeasurementUnits(ctx context.Context, request *mealplanning.GetValidMeasurementUnitsRequest) (*mealplanning.GetValidMeasurementUnitsResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)
	filter := grpcconverters.ConvertGRPCQueryFilterToQueryFilter(request.Filter)
	tracing.AttachQueryFilterToSpan(span, filter)

	x, err := s.validEnumerationsManager.ListValidMeasurementUnits(ctx, filter)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "fetching valid measurement units")
	}

	res := &mealplanning.GetValidMeasurementUnitsResponse{
		Filter: grpcconverters.ConvertQueryFilterToGRPCQueryFilter(filter, x.Pagination),
	}

	for _, y := range x.Data {
		res.Results = append(res.Results, mealplanningconverters.ConvertValidMeasurementUnitToGRPCValidMeasurementUnit(y))
	}

	return res, nil
}

func (s *serviceImpl) GetValidPreparation(ctx context.Context, request *mealplanning.GetValidPreparationRequest) (*mealplanning.GetValidPreparationResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValue(keys.ValidPreparationIDKey, request.ValidPreparationID)
	tracing.AttachToSpan(span, keys.ValidPreparationIDKey, request.ValidPreparationID)

	x, err := s.validEnumerationsManager.ReadValidPreparation(ctx, request.ValidPreparationID)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "fetching valid preparation")
	}

	res := &mealplanning.GetValidPreparationResponse{
		Result: mealplanningconverters.ConvertValidPreparationToGRPCValidPreparation(x),
	}

	return res, nil
}

func (s *serviceImpl) GetValidPreparationInstrument(ctx context.Context, request *mealplanning.GetValidPreparationInstrumentRequest) (*mealplanning.GetValidPreparationInstrumentResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValue(keys.ValidPreparationInstrumentIDKey, request.ValidPreparationInstrumentID)
	tracing.AttachToSpan(span, keys.ValidPreparationInstrumentIDKey, request.ValidPreparationInstrumentID)

	x, err := s.validEnumerationsManager.ReadValidPreparationInstrument(ctx, request.ValidPreparationInstrumentID)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "fetching valid preparation instrument")
	}

	res := &mealplanning.GetValidPreparationInstrumentResponse{
		Result: mealplanningconverters.ConvertValidPreparationInstrumentToGRPCValidPreparationInstrument(x),
	}

	return res, nil
}

func (s *serviceImpl) GetValidPreparationInstruments(ctx context.Context, request *mealplanning.GetValidPreparationInstrumentsRequest) (*mealplanning.GetValidPreparationInstrumentsResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)
	filter := grpcconverters.ConvertGRPCQueryFilterToQueryFilter(request.Filter)
	tracing.AttachQueryFilterToSpan(span, filter)

	x, err := s.validEnumerationsManager.ListValidPreparationInstruments(ctx, filter)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "fetching valid preparation instruments")
	}

	res := &mealplanning.GetValidPreparationInstrumentsResponse{
		Filter: grpcconverters.ConvertQueryFilterToGRPCQueryFilter(filter, x.Pagination),
	}

	for _, y := range x.Data {
		res.Results = append(res.Results, mealplanningconverters.ConvertValidPreparationInstrumentToGRPCValidPreparationInstrument(y))
	}

	return res, nil
}

func (s *serviceImpl) GetValidPreparationInstrumentsByInstrument(ctx context.Context, request *mealplanning.GetValidPreparationInstrumentsByInstrumentRequest) (*mealplanning.GetValidPreparationInstrumentsByInstrumentResponse, error) {
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
		Filter: request.Filter, // TODO: grpcconverters.ConvertQueryFilterToGRPCQueryFilter(filter, x.Pagination),
	}

	for _, y := range x {
		res.Results = append(res.Results, mealplanningconverters.ConvertValidPreparationInstrumentToGRPCValidPreparationInstrument(y))
	}

	return res, nil
}

func (s *serviceImpl) GetValidPreparationInstrumentsByPreparation(ctx context.Context, request *mealplanning.GetValidPreparationInstrumentsByPreparationRequest) (*mealplanning.GetValidPreparationInstrumentsByPreparationResponse, error) {
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
		Filter: request.Filter, // TODO: grpcconverters.ConvertQueryFilterToGRPCQueryFilter(filter, x.Pagination),
	}

	for _, y := range x {
		res.Results = append(res.Results, mealplanningconverters.ConvertValidPreparationInstrumentToGRPCValidPreparationInstrument(y))
	}

	return res, nil
}

func (s *serviceImpl) GetValidPreparationVessel(ctx context.Context, request *mealplanning.GetValidPreparationVesselRequest) (*mealplanning.GetValidPreparationVesselResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValue(keys.ValidPreparationVesselIDKey, request.ValidPreparationVesselID)
	tracing.AttachToSpan(span, keys.ValidPreparationVesselIDKey, request.ValidPreparationVesselID)

	x, err := s.validEnumerationsManager.ReadValidPreparationVessel(ctx, request.ValidPreparationVesselID)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "fetching valid preparation vessel")
	}

	res := &mealplanning.GetValidPreparationVesselResponse{
		Result: mealplanningconverters.ConvertValidPreparationVesselToGRPCValidPreparationVessel(x),
	}

	return res, nil
}

func (s *serviceImpl) GetValidPreparationVessels(ctx context.Context, request *mealplanning.GetValidPreparationVesselsRequest) (*mealplanning.GetValidPreparationVesselsResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)
	filter := grpcconverters.ConvertGRPCQueryFilterToQueryFilter(request.Filter)
	tracing.AttachQueryFilterToSpan(span, filter)

	x, err := s.validEnumerationsManager.ListValidPreparationVessels(ctx, filter)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "fetching valid preparation vessels")
	}

	res := &mealplanning.GetValidPreparationVesselsResponse{
		Filter: grpcconverters.ConvertQueryFilterToGRPCQueryFilter(filter, x.Pagination),
	}

	for _, y := range x.Data {
		res.Results = append(res.Results, mealplanningconverters.ConvertValidPreparationVesselToGRPCValidPreparationVessel(y))
	}

	return res, nil
}

func (s *serviceImpl) GetValidPreparationVesselsByPreparation(ctx context.Context, request *mealplanning.GetValidPreparationVesselsByPreparationRequest) (*mealplanning.GetValidPreparationVesselsByPreparationResponse, error) {
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
		Filter: request.Filter, // TODO: grpcconverters.ConvertQueryFilterToGRPCQueryFilter(filter, x.Pagination),
	}

	for _, y := range x {
		res.Results = append(res.Results, mealplanningconverters.ConvertValidPreparationVesselToGRPCValidPreparationVessel(y))
	}

	return res, nil
}

func (s *serviceImpl) GetValidPreparationVesselsByVessel(ctx context.Context, request *mealplanning.GetValidPreparationVesselsByVesselRequest) (*mealplanning.GetValidPreparationVesselsByVesselResponse, error) {
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
		Filter: request.Filter, // TODO: grpcconverters.ConvertQueryFilterToGRPCQueryFilter(filter, x.Pagination),
	}

	for _, y := range x {
		res.Results = append(res.Results, mealplanningconverters.ConvertValidPreparationVesselToGRPCValidPreparationVessel(y))
	}

	return res, nil
}

func (s *serviceImpl) GetValidPreparations(ctx context.Context, request *mealplanning.GetValidPreparationsRequest) (*mealplanning.GetValidPreparationsResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)
	filter := grpcconverters.ConvertGRPCQueryFilterToQueryFilter(request.Filter)
	tracing.AttachQueryFilterToSpan(span, filter)

	x, err := s.validEnumerationsManager.ListValidPreparations(ctx, filter)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "fetching valid preparations")
	}

	res := &mealplanning.GetValidPreparationsResponse{
		Filter: grpcconverters.ConvertQueryFilterToGRPCQueryFilter(filter, x.Pagination),
	}

	for _, y := range x.Data {
		res.Results = append(res.Results, mealplanningconverters.ConvertValidPreparationToGRPCValidPreparation(y))
	}

	return res, nil
}

func (s *serviceImpl) GetValidVessel(ctx context.Context, request *mealplanning.GetValidVesselRequest) (*mealplanning.GetValidVesselResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValue(keys.ValidVesselIDKey, request.ValidVesselID)
	tracing.AttachToSpan(span, keys.ValidVesselIDKey, request.ValidVesselID)

	x, err := s.validEnumerationsManager.ReadValidVessel(ctx, request.ValidVesselID)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "fetching valid vessel")
	}

	res := &mealplanning.GetValidVesselResponse{
		Result: mealplanningconverters.ConvertValidVesselToGRPCValidVessel(x),
	}

	return res, nil
}

func (s *serviceImpl) GetValidVessels(ctx context.Context, request *mealplanning.GetValidVesselsRequest) (*mealplanning.GetValidVesselsResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)
	filter := grpcconverters.ConvertGRPCQueryFilterToQueryFilter(request.Filter)
	tracing.AttachQueryFilterToSpan(span, filter)

	x, err := s.validEnumerationsManager.ListValidVessels(ctx, filter)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "fetching valid vessels")
	}

	res := &mealplanning.GetValidVesselsResponse{
		Filter: grpcconverters.ConvertQueryFilterToGRPCQueryFilter(filter, x.Pagination),
	}

	for _, y := range x.Data {
		res.Results = append(res.Results, mealplanningconverters.ConvertValidVesselToGRPCValidVessel(y))
	}

	return res, nil
}

func (s *serviceImpl) SearchForValidIngredientGroups(ctx context.Context, request *mealplanning.SearchForValidIngredientGroupsRequest) (*mealplanning.SearchForValidIngredientGroupsResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)
	filter := grpcconverters.ConvertGRPCQueryFilterToQueryFilter(request.Filter)
	tracing.AttachQueryFilterToSpan(span, filter)

	x, err := s.validEnumerationsManager.SearchValidIngredientGroups(ctx, request.Query, request.UseSearchService, filter)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "searching for valid ingredient groups")
	}

	res := &mealplanning.SearchForValidIngredientGroupsResponse{
		Filter: request.Filter, // TODO: grpcconverters.ConvertQueryFilterToGRPCQueryFilter(filter, x.Pagination),
	}

	for _, y := range x.Data {
		res.Results = append(res.Results, mealplanningconverters.ConvertValidIngredientGroupToGRPCValidIngredientGroup(y))
	}

	return res, nil
}

func (s *serviceImpl) SearchForValidIngredientStates(ctx context.Context, request *mealplanning.SearchForValidIngredientStatesRequest) (*mealplanning.SearchForValidIngredientStatesResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)
	filter := grpcconverters.ConvertGRPCQueryFilterToQueryFilter(request.Filter)
	tracing.AttachQueryFilterToSpan(span, filter)

	x, err := s.validEnumerationsManager.SearchValidIngredientStates(ctx, request.Query, request.UseSearchService, filter)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "searching for valid ingredient states")
	}

	res := &mealplanning.SearchForValidIngredientStatesResponse{
		Filter: request.Filter, // TODO: grpcconverters.ConvertQueryFilterToGRPCQueryFilter(filter, x.Pagination),
	}

	for _, y := range x {
		res.Results = append(res.Results, mealplanningconverters.ConvertValidIngredientStateToGRPCValidIngredientState(y))
	}

	return res, nil
}

func (s *serviceImpl) SearchForValidIngredients(ctx context.Context, request *mealplanning.SearchForValidIngredientsRequest) (*mealplanning.SearchForValidIngredientsResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)
	filter := grpcconverters.ConvertGRPCQueryFilterToQueryFilter(request.Filter)
	tracing.AttachQueryFilterToSpan(span, filter)

	x, err := s.validEnumerationsManager.SearchValidIngredients(ctx, request.Query, request.UseSearchService, filter)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "searching for valid ingredients")
	}

	res := &mealplanning.SearchForValidIngredientsResponse{
		Filter: request.Filter, // TODO: grpcconverters.ConvertQueryFilterToGRPCQueryFilter(filter, x.Pagination),
	}

	for _, y := range x {
		res.Results = append(res.Results, mealplanningconverters.ConvertValidIngredientToGRPCValidIngredient(y))
	}

	return res, nil
}

func (s *serviceImpl) SearchForValidInstruments(ctx context.Context, request *mealplanning.SearchForValidInstrumentsRequest) (*mealplanning.SearchForValidInstrumentsResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)
	filter := grpcconverters.ConvertGRPCQueryFilterToQueryFilter(request.Filter)
	tracing.AttachQueryFilterToSpan(span, filter)

	x, err := s.validEnumerationsManager.SearchValidInstruments(ctx, request.Query, request.UseSearchService, filter)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "searching for valid instruments")
	}

	res := &mealplanning.SearchForValidInstrumentsResponse{
		Filter: request.Filter, // TODO: grpcconverters.ConvertQueryFilterToGRPCQueryFilter(filter, x.Pagination),
	}

	for _, y := range x {
		res.Results = append(res.Results, mealplanningconverters.ConvertValidInstrumentToGRPCValidInstrument(y))
	}

	return res, nil
}

func (s *serviceImpl) SearchForValidMeasurementUnits(ctx context.Context, request *mealplanning.SearchForValidMeasurementUnitsRequest) (*mealplanning.SearchForValidMeasurementUnitsResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)
	filter := grpcconverters.ConvertGRPCQueryFilterToQueryFilter(request.Filter)
	tracing.AttachQueryFilterToSpan(span, filter)

	x, err := s.validEnumerationsManager.SearchValidMeasurementUnits(ctx, request.Query, request.UseSearchService, filter)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "searching for valid measurement units")
	}

	res := &mealplanning.SearchForValidMeasurementUnitsResponse{
		Filter: request.Filter, // TODO: grpcconverters.ConvertQueryFilterToGRPCQueryFilter(filter, x.Pagination),
	}

	for _, y := range x {
		res.Results = append(res.Results, mealplanningconverters.ConvertValidMeasurementUnitToGRPCValidMeasurementUnit(y))
	}

	return res, nil
}

func (s *serviceImpl) SearchForValidPreparations(ctx context.Context, request *mealplanning.SearchForValidPreparationsRequest) (*mealplanning.SearchForValidPreparationsResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)
	filter := grpcconverters.ConvertGRPCQueryFilterToQueryFilter(request.Filter)
	tracing.AttachQueryFilterToSpan(span, filter)

	x, err := s.validEnumerationsManager.SearchValidPreparations(ctx, request.Query, request.UseSearchService, filter)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "searching for valid preparations")
	}

	res := &mealplanning.SearchForValidPreparationsResponse{
		Filter: request.Filter, // TODO: grpcconverters.ConvertQueryFilterToGRPCQueryFilter(filter, x.Pagination),
	}

	for _, y := range x {
		res.Results = append(res.Results, mealplanningconverters.ConvertValidPreparationToGRPCValidPreparation(y))
	}

	return res, nil
}

func (s *serviceImpl) SearchForValidVessels(ctx context.Context, request *mealplanning.SearchForValidVesselsRequest) (*mealplanning.SearchForValidVesselsResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)
	filter := grpcconverters.ConvertGRPCQueryFilterToQueryFilter(request.Filter)
	tracing.AttachQueryFilterToSpan(span, filter)

	x, err := s.validEnumerationsManager.SearchValidVessels(ctx, request.Query, request.UseSearchService, filter)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "searching for valid vessels")
	}

	res := &mealplanning.SearchForValidVesselsResponse{
		Filter: request.Filter, // TODO: grpcconverters.ConvertQueryFilterToGRPCQueryFilter(filter, x.Pagination),
	}

	for _, y := range x {
		res.Results = append(res.Results, mealplanningconverters.ConvertValidVesselToGRPCValidVessel(y))
	}

	return res, nil
}

func (s *serviceImpl) SearchValidIngredientsByPreparation(ctx context.Context, request *mealplanning.SearchValidIngredientsByPreparationRequest) (*mealplanning.SearchValidIngredientsByPreparationResponse, error) {
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
		Filter: request.Filter, // TODO: grpcconverters.ConvertQueryFilterToGRPCQueryFilter(filter, x.Pagination),
	}

	for _, y := range x {
		res.Results = append(res.Results, mealplanningconverters.ConvertValidIngredientToGRPCValidIngredient(y))
	}

	return res, nil
}

func (s *serviceImpl) SearchValidMeasurementUnitsByIngredient(ctx context.Context, request *mealplanning.SearchValidMeasurementUnitsByIngredientRequest) (*mealplanning.SearchValidMeasurementUnitsByIngredientResponse, error) {
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
		Filter: request.Filter, // TODO: grpcconverters.ConvertQueryFilterToGRPCQueryFilter(filter, x.Pagination),
	}

	for _, y := range x {
		res.Results = append(res.Results, mealplanningconverters.ConvertValidMeasurementUnitToGRPCValidMeasurementUnit(y))
	}

	return res, nil
}

func (s *serviceImpl) UpdateValidIngredient(ctx context.Context, request *mealplanning.UpdateValidIngredientRequest) (*mealplanning.UpdateValidIngredientResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValue(keys.ValidIngredientIDKey, request.ValidIngredientID)
	tracing.AttachToSpan(span, keys.ValidIngredientIDKey, request.ValidIngredientID)

	input := mealplanningconverters.ConvertGRPCValidIngredientUpdateRequestInputToValidIngredientUpdateRequestInput(request.Input)
	updated, err := s.validEnumerationsManager.UpdateValidIngredient(ctx, request.ValidIngredientID, input)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "updating valid ingredient")
	}

	res := &mealplanning.UpdateValidIngredientResponse{
		Result: mealplanningconverters.ConvertValidIngredientToGRPCValidIngredient(updated),
	}

	return res, nil
}

func (s *serviceImpl) UpdateValidIngredientGroup(ctx context.Context, request *mealplanning.UpdateValidIngredientGroupRequest) (*mealplanning.UpdateValidIngredientGroupResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValue(keys.ValidIngredientGroupIDKey, request.ValidIngredientGroupID)
	tracing.AttachToSpan(span, keys.ValidIngredientGroupIDKey, request.ValidIngredientGroupID)

	input := mealplanningconverters.ConvertGRPCValidIngredientGroupUpdateRequestInputToValidIngredientGroupUpdateRequestInput(request.Input)

	updated, err := s.validEnumerationsManager.UpdateValidIngredientGroup(ctx, request.ValidIngredientGroupID, input)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "updating valid ingredient group")
	}

	res := &mealplanning.UpdateValidIngredientGroupResponse{
		Result: mealplanningconverters.ConvertValidIngredientGroupToGRPCValidIngredientGroup(updated),
	}

	return res, nil
}

func (s *serviceImpl) UpdateValidIngredientMeasurementUnit(ctx context.Context, request *mealplanning.UpdateValidIngredientMeasurementUnitRequest) (*mealplanning.UpdateValidIngredientMeasurementUnitResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValue(keys.ValidIngredientMeasurementUnitIDKey, request.ValidIngredientMeasurementUnitID)
	tracing.AttachToSpan(span, keys.ValidIngredientMeasurementUnitIDKey, request.ValidIngredientMeasurementUnitID)

	input := mealplanningconverters.ConvertGRPCValidIngredientMeasurementUnitUpdateRequestInputToValidIngredientMeasurementUnitUpdateRequestInput(request.Input)

	updated, err := s.validEnumerationsManager.UpdateValidIngredientMeasurementUnit(ctx, request.ValidIngredientMeasurementUnitID, input)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "updating valid ingredient measurement unit")
	}

	res := &mealplanning.UpdateValidIngredientMeasurementUnitResponse{
		Result: mealplanningconverters.ConvertValidIngredientMeasurementUnitToGRPCValidIngredientMeasurementUnit(updated),
	}

	return res, nil
}

func (s *serviceImpl) UpdateValidIngredientPreparation(ctx context.Context, request *mealplanning.UpdateValidIngredientPreparationRequest) (*mealplanning.UpdateValidIngredientPreparationResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValue(keys.ValidIngredientPreparationIDKey, request.ValidIngredientPreparationID)
	tracing.AttachToSpan(span, keys.ValidIngredientPreparationIDKey, request.ValidIngredientPreparationID)

	input := mealplanningconverters.ConvertGRPCValidIngredientPreparationUpdateRequestInputToValidIngredientPreparationUpdateRequestInput(request.Input)

	updated, err := s.validEnumerationsManager.UpdateValidIngredientPreparation(ctx, request.ValidIngredientPreparationID, input)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "updating valid ingredient preparation")
	}

	res := &mealplanning.UpdateValidIngredientPreparationResponse{
		Result: mealplanningconverters.ConvertValidIngredientPreparationToGRPCValidIngredientPreparation(updated),
	}

	return res, nil
}

func (s *serviceImpl) UpdateValidIngredientState(ctx context.Context, request *mealplanning.UpdateValidIngredientStateRequest) (*mealplanning.UpdateValidIngredientStateResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValue(keys.ValidIngredientStateIDKey, request.ValidIngredientStateID)
	tracing.AttachToSpan(span, keys.ValidIngredientStateIDKey, request.ValidIngredientStateID)

	input := mealplanningconverters.ConvertGRPCValidIngredientStateUpdateRequestInputToValidIngredientStateUpdateRequestInput(request.Input)

	updated, err := s.validEnumerationsManager.UpdateValidIngredientState(ctx, request.ValidIngredientStateID, input)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "updating valid ingredient state")
	}

	res := &mealplanning.UpdateValidIngredientStateResponse{
		Result: mealplanningconverters.ConvertValidIngredientStateToGRPCValidIngredientState(updated),
	}

	return res, nil
}

func (s *serviceImpl) UpdateValidIngredientStateIngredient(ctx context.Context, request *mealplanning.UpdateValidIngredientStateIngredientRequest) (*mealplanning.UpdateValidIngredientStateIngredientResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValue(keys.ValidIngredientStateIngredientIDKey, request.ValidIngredientStateIngredientID)
	tracing.AttachToSpan(span, keys.ValidIngredientStateIngredientIDKey, request.ValidIngredientStateIngredientID)

	input := mealplanningconverters.ConvertGRPCValidIngredientStateIngredientUpdateRequestInputToValidIngredientStateIngredientUpdateRequestInput(request.Input)

	updated, err := s.validEnumerationsManager.UpdateValidIngredientStateIngredient(ctx, request.ValidIngredientStateIngredientID, input)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "updating valid ingredient state ingredient")
	}

	res := &mealplanning.UpdateValidIngredientStateIngredientResponse{
		Result: mealplanningconverters.ConvertValidIngredientStateIngredientToGRPCValidIngredientStateIngredient(updated),
	}

	return res, nil
}

func (s *serviceImpl) UpdateValidInstrument(ctx context.Context, request *mealplanning.UpdateValidInstrumentRequest) (*mealplanning.UpdateValidInstrumentResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValue(keys.ValidInstrumentIDKey, request.ValidInstrumentID)
	tracing.AttachToSpan(span, keys.ValidInstrumentIDKey, request.ValidInstrumentID)

	input := mealplanningconverters.ConvertGRPCValidInstrumentUpdateRequestInputToValidInstrumentUpdateRequestInput(request.Input)

	updated, err := s.validEnumerationsManager.UpdateValidInstrument(ctx, request.ValidInstrumentID, input)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "updating valid instrument")
	}

	res := &mealplanning.UpdateValidInstrumentResponse{
		Result: mealplanningconverters.ConvertValidInstrumentToGRPCValidInstrument(updated),
	}

	return res, nil
}

func (s *serviceImpl) UpdateValidMeasurementUnit(ctx context.Context, request *mealplanning.UpdateValidMeasurementUnitRequest) (*mealplanning.UpdateValidMeasurementUnitResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValue(keys.ValidMeasurementUnitIDKey, request.ValidMeasurementUnitID)
	tracing.AttachToSpan(span, keys.ValidMeasurementUnitIDKey, request.ValidMeasurementUnitID)

	input := mealplanningconverters.ConvertGRPCValidMeasurementUnitUpdateRequestInputToValidMeasurementUnitUpdateRequestInput(request.Input)

	updated, err := s.validEnumerationsManager.UpdateValidMeasurementUnit(ctx, request.ValidMeasurementUnitID, input)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "updating valid measurement unit")
	}

	res := &mealplanning.UpdateValidMeasurementUnitResponse{
		Result: mealplanningconverters.ConvertValidMeasurementUnitToGRPCValidMeasurementUnit(updated),
	}

	return res, nil
}

func (s *serviceImpl) UpdateValidMeasurementUnitConversion(ctx context.Context, request *mealplanning.UpdateValidMeasurementUnitConversionRequest) (*mealplanning.UpdateValidMeasurementUnitConversionResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValue(keys.ValidMeasurementUnitConversionIDKey, request.ValidMeasurementUnitConversionID)
	tracing.AttachToSpan(span, keys.ValidMeasurementUnitConversionIDKey, request.ValidMeasurementUnitConversionID)

	input := mealplanningconverters.ConvertGRPCValidMeasurementUnitConversionUpdateRequestInputToValidMeasurementUnitConversionUpdateRequestInput(request.Input)

	updated, err := s.validEnumerationsManager.UpdateValidMeasurementUnitConversion(ctx, request.ValidMeasurementUnitConversionID, input)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "updating valid measurement unit conversion")
	}

	res := &mealplanning.UpdateValidMeasurementUnitConversionResponse{
		Result: mealplanningconverters.ConvertValidMeasurementUnitConversionToGRPCValidMeasurementUnitConversion(updated),
	}

	return res, nil
}

func (s *serviceImpl) UpdateValidPreparation(ctx context.Context, request *mealplanning.UpdateValidPreparationRequest) (*mealplanning.UpdateValidPreparationResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValue(keys.ValidPreparationIDKey, request.ValidPreparationID)
	tracing.AttachToSpan(span, keys.ValidPreparationIDKey, request.ValidPreparationID)

	input := mealplanningconverters.ConvertGRPCValidPreparationUpdateRequestInputToValidPreparationUpdateRequestInput(request.Input)

	updated, err := s.validEnumerationsManager.UpdateValidPreparation(ctx, request.ValidPreparationID, input)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "updating valid preparation")
	}

	res := &mealplanning.UpdateValidPreparationResponse{
		Result: mealplanningconverters.ConvertValidPreparationToGRPCValidPreparation(updated),
	}

	return res, nil
}

func (s *serviceImpl) UpdateValidPreparationInstrument(ctx context.Context, request *mealplanning.UpdateValidPreparationInstrumentRequest) (*mealplanning.UpdateValidPreparationInstrumentResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValue(keys.ValidPreparationInstrumentIDKey, request.ValidPreparationInstrumentID)
	tracing.AttachToSpan(span, keys.ValidPreparationInstrumentIDKey, request.ValidPreparationInstrumentID)

	input := mealplanningconverters.ConvertGRPCValidPreparationInstrumentUpdateRequestInputToValidPreparationInstrumentUpdateRequestInput(request.Input)

	updated, err := s.validEnumerationsManager.UpdateValidPreparationInstrument(ctx, request.ValidPreparationInstrumentID, input)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "updating valid preparation instrument")
	}

	res := &mealplanning.UpdateValidPreparationInstrumentResponse{
		Result: mealplanningconverters.ConvertValidPreparationInstrumentToGRPCValidPreparationInstrument(updated),
	}

	return res, nil
}

func (s *serviceImpl) UpdateValidPreparationVessel(ctx context.Context, request *mealplanning.UpdateValidPreparationVesselRequest) (*mealplanning.UpdateValidPreparationVesselResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValue(keys.ValidPreparationVesselIDKey, request.ValidPreparationVesselID)
	tracing.AttachToSpan(span, keys.ValidPreparationVesselIDKey, request.ValidPreparationVesselID)

	input := mealplanningconverters.ConvertGRPCValidPreparationVesselUpdateRequestInputToValidPreparationVesselUpdateRequestInput(request.Input)

	updated, err := s.validEnumerationsManager.UpdateValidPreparationVessel(ctx, request.ValidPreparationVesselID, input)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "updating valid preparation vessel")
	}

	res := &mealplanning.UpdateValidPreparationVesselResponse{
		Result: mealplanningconverters.ConvertValidPreparationVesselToGRPCValidPreparationVessel(updated),
	}

	return res, nil
}

func (s *serviceImpl) UpdateValidVessel(ctx context.Context, request *mealplanning.UpdateValidVesselRequest) (*mealplanning.UpdateValidVesselResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValue(keys.ValidVesselIDKey, request.ValidVesselID)
	tracing.AttachToSpan(span, keys.ValidVesselIDKey, request.ValidVesselID)

	input := mealplanningconverters.ConvertGRPCValidVesselUpdateRequestInputToValidVesselUpdateRequestInput(request.Input)

	updated, err := s.validEnumerationsManager.UpdateValidVessel(ctx, request.ValidVesselID, input)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "updating valid vessel")
	}

	res := &mealplanning.UpdateValidVesselResponse{
		Result: mealplanningconverters.ConvertValidVesselToGRPCValidVessel(updated),
	}

	return res, nil
}
