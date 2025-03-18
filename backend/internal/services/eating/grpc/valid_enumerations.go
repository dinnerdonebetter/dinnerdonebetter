package grpc

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/grpc/messages"
	"github.com/dinnerdonebetter/backend/internal/lib/observability"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/tracing"
	grpcconverters "github.com/dinnerdonebetter/backend/internal/services/eating/grpc/converters"

	"google.golang.org/grpc/codes"
)

func (s *serviceImpl) ArchiveValidIngredient(ctx context.Context, request *messages.ArchiveValidIngredientRequest) (*messages.ArchiveValidIngredientResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	if err := s.validEnumerationsManager.ArchiveValidIngredient(ctx, request.ValidIngredientID); err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "archiving valid ingredient")
	}

	return &messages.ArchiveValidIngredientResponse{}, nil
}

func (s *serviceImpl) ArchiveValidIngredientGroup(ctx context.Context, request *messages.ArchiveValidIngredientGroupRequest) (*messages.ArchiveValidIngredientGroupResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	if err := s.validEnumerationsManager.ArchiveValidIngredientGroup(ctx, request.ValidIngredientGroupID); err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "archiving valid ingredient group")
	}

	return &messages.ArchiveValidIngredientGroupResponse{}, nil
}

func (s *serviceImpl) ArchiveValidIngredientMeasurementUnit(ctx context.Context, request *messages.ArchiveValidIngredientMeasurementUnitRequest) (*messages.ArchiveValidIngredientMeasurementUnitResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	if err := s.validEnumerationsManager.ArchiveValidIngredientMeasurementUnit(ctx, request.ValidIngredientMeasurementUnitID); err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "archiving valid ingredient measurement unit")
	}

	return &messages.ArchiveValidIngredientMeasurementUnitResponse{}, nil
}

func (s *serviceImpl) ArchiveValidIngredientPreparation(ctx context.Context, request *messages.ArchiveValidIngredientPreparationRequest) (*messages.ArchiveValidIngredientPreparationResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	if err := s.validEnumerationsManager.ArchiveValidIngredientPreparation(ctx, request.ValidIngredientPreparationID); err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "archiving valid ingredient preparation")
	}

	return &messages.ArchiveValidIngredientPreparationResponse{}, nil
}

func (s *serviceImpl) ArchiveValidIngredientState(ctx context.Context, request *messages.ArchiveValidIngredientStateRequest) (*messages.ArchiveValidIngredientStateResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	if err := s.validEnumerationsManager.ArchiveValidIngredientState(ctx, request.ValidIngredientStateID); err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "archiving valid ingredient state")
	}

	return &messages.ArchiveValidIngredientStateResponse{}, nil
}

func (s *serviceImpl) ArchiveValidIngredientStateIngredient(ctx context.Context, request *messages.ArchiveValidIngredientStateIngredientRequest) (*messages.ArchiveValidIngredientStateIngredientResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	if err := s.validEnumerationsManager.ArchiveValidIngredientStateIngredient(ctx, request.ValidIngredientStateIngredientID); err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "archiving valid ingredient state ingredient")
	}

	return &messages.ArchiveValidIngredientStateIngredientResponse{}, nil
}

func (s *serviceImpl) ArchiveValidInstrument(ctx context.Context, request *messages.ArchiveValidInstrumentRequest) (*messages.ArchiveValidInstrumentResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	if err := s.validEnumerationsManager.ArchiveValidInstrument(ctx, request.ValidInstrumentID); err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "archiving valid instrument")
	}

	return &messages.ArchiveValidInstrumentResponse{}, nil
}

func (s *serviceImpl) ArchiveValidMeasurementUnit(ctx context.Context, request *messages.ArchiveValidMeasurementUnitRequest) (*messages.ArchiveValidMeasurementUnitResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	if err := s.validEnumerationsManager.ArchiveValidMeasurementUnit(ctx, request.ValidMeasurementUnitID); err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "archiving valid measurement unit")
	}

	return &messages.ArchiveValidMeasurementUnitResponse{}, nil
}

func (s *serviceImpl) ArchiveValidMeasurementUnitConversion(ctx context.Context, request *messages.ArchiveValidMeasurementUnitConversionRequest) (*messages.ArchiveValidMeasurementUnitConversionResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	if err := s.validEnumerationsManager.ArchiveValidMeasurementUnitConversion(ctx, request.ValidMeasurementUnitConversionID); err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "archiving valid measurement unit conversion")
	}

	return &messages.ArchiveValidMeasurementUnitConversionResponse{}, nil
}

func (s *serviceImpl) ArchiveValidPreparation(ctx context.Context, request *messages.ArchiveValidPreparationRequest) (*messages.ArchiveValidPreparationResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	if err := s.validEnumerationsManager.ArchiveValidPreparation(ctx, request.ValidPreparationID); err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "archiving valid preparation")
	}

	return &messages.ArchiveValidPreparationResponse{}, nil
}

func (s *serviceImpl) ArchiveValidPreparationInstrument(ctx context.Context, request *messages.ArchiveValidPreparationInstrumentRequest) (*messages.ArchiveValidPreparationInstrumentResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	if err := s.validEnumerationsManager.ArchiveValidPreparationInstrument(ctx, request.ValidPreparationInstrumentID); err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "archiving valid preparation instrument")
	}

	return &messages.ArchiveValidPreparationInstrumentResponse{}, nil
}

func (s *serviceImpl) ArchiveValidPreparationVessel(ctx context.Context, request *messages.ArchiveValidPreparationVesselRequest) (*messages.ArchiveValidPreparationVesselResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	if err := s.validEnumerationsManager.ArchiveValidPreparationVessel(ctx, request.ValidPreparationVesselID); err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "archiving valid preparation vessel")
	}

	return &messages.ArchiveValidPreparationVesselResponse{}, nil
}

func (s *serviceImpl) ArchiveValidVessel(ctx context.Context, request *messages.ArchiveValidVesselRequest) (*messages.ArchiveValidVesselResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	if err := s.validEnumerationsManager.ArchiveValidVessel(ctx, request.ValidVesselID); err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "archiving valid vessel")
	}

	return &messages.ArchiveValidVesselResponse{}, nil
}

func (s *serviceImpl) CreateValidIngredient(ctx context.Context, request *messages.CreateValidIngredientRequest) (*messages.CreateValidIngredientResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	created, err := s.validEnumerationsManager.CreateValidIngredient(ctx, grpcconverters.ConvertGRPCCreateValidIngredientRequestToValidIngredientCreationRequestInput(request))
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "creating valid ingredient")
	}

	result := &messages.CreateValidIngredientResponse{
		Result: grpcconverters.ConvertValidIngredientToGRPCValidIngredient(created),
	}

	return result, nil
}

func (s *serviceImpl) CreateValidIngredientGroup(ctx context.Context, request *messages.CreateValidIngredientGroupRequest) (*messages.CreateValidIngredientGroupResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	created, err := s.validEnumerationsManager.CreateValidIngredientGroup(ctx, grpcconverters.ConvertGRPCCreateValidIngredientGroupRequestToValidIngredientGroupCreationRequestInput(request))
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "creating valid ingredient")
	}

	result := &messages.CreateValidIngredientGroupResponse{
		Result: grpcconverters.ConvertValidIngredientGroupToGRPCValidIngredientGroup(created),
	}

	return result, nil
}

func (s *serviceImpl) CreateValidIngredientMeasurementUnit(ctx context.Context, request *messages.CreateValidIngredientMeasurementUnitRequest) (*messages.CreateValidIngredientMeasurementUnitResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	created, err := s.validEnumerationsManager.CreateValidIngredientMeasurementUnit(ctx, grpcconverters.ConvertGRPCCreateValidIngredientMeasurementUnitRequestToValidIngredientMeasurementUnitCreationRequestInput(request))
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "creating valid ingredient measurement unit")
	}

	result := &messages.CreateValidIngredientMeasurementUnitResponse{
		Result: grpcconverters.ConvertValidIngredientMeasurementUnitToGRPCValidIngredientMeasurementUnit(created),
	}

	return result, nil
}

func (s *serviceImpl) CreateValidIngredientPreparation(ctx context.Context, request *messages.CreateValidIngredientPreparationRequest) (*messages.CreateValidIngredientPreparationResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	created, err := s.validEnumerationsManager.CreateValidIngredientPreparation(ctx, grpcconverters.ConvertGRPCCreateValidIngredientPreparationRequestToValidIngredientPreparationCreationRequestInput(request))
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "creating valid ingredient preparation")
	}

	result := &messages.CreateValidIngredientPreparationResponse{
		Result: grpcconverters.ConvertValidIngredientPreparationToGRPCValidIngredientPreparation(created),
	}

	return result, nil
}

func (s *serviceImpl) CreateValidIngredientState(ctx context.Context, request *messages.CreateValidIngredientStateRequest) (*messages.CreateValidIngredientStateResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	created, err := s.validEnumerationsManager.CreateValidIngredientState(ctx, grpcconverters.ConvertGRPCCreateValidIngredientStateRequestToValidIngredientStateCreationRequestInput(request))
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "creating valid ingredient state")
	}

	result := &messages.CreateValidIngredientStateResponse{
		Result: grpcconverters.ConvertValidIngredientStateToGRPCValidIngredientState(created),
	}

	return result, nil
}

func (s *serviceImpl) CreateValidIngredientStateIngredient(ctx context.Context, request *messages.CreateValidIngredientStateIngredientRequest) (*messages.CreateValidIngredientStateIngredientResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	created, err := s.validEnumerationsManager.CreateValidIngredientStateIngredient(ctx, grpcconverters.ConvertGRPCCreateValidIngredientStateIngredientRequestToValidIngredientStateIngredientCreationRequestInput(request))
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "creating valid ingredient state ingredient")
	}

	result := &messages.CreateValidIngredientStateIngredientResponse{
		Result: grpcconverters.ConvertValidIngredientStateIngredientToGRPCValidIngredientStateIngredient(created),
	}

	return result, nil
}

func (s *serviceImpl) CreateValidInstrument(ctx context.Context, request *messages.CreateValidInstrumentRequest) (*messages.CreateValidInstrumentResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	created, err := s.validEnumerationsManager.CreateValidInstrument(ctx, grpcconverters.ConvertGRPCCreateValidInstrumentRequestToValidInstrumentCreationRequestInput(request))
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "creating valid instrument")
	}

	result := &messages.CreateValidInstrumentResponse{
		Result: grpcconverters.ConvertValidInstrumentToGRPCValidInstrument(created),
	}

	return result, nil
}

func (s *serviceImpl) CreateValidMeasurementUnit(ctx context.Context, request *messages.CreateValidMeasurementUnitRequest) (*messages.CreateValidMeasurementUnitResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	created, err := s.validEnumerationsManager.CreateValidMeasurementUnit(ctx, grpcconverters.ConvertGRPCCreateValidMeasurementUnitRequestToValidMeasurementUnitCreationRequestInput(request))
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "creating valid measurement unit")
	}

	result := &messages.CreateValidMeasurementUnitResponse{
		Result: grpcconverters.ConvertValidMeasurementUnitToGRPCValidMeasurementUnit(created),
	}

	return result, nil
}

func (s *serviceImpl) CreateValidMeasurementUnitConversion(ctx context.Context, request *messages.CreateValidMeasurementUnitConversionRequest) (*messages.CreateValidMeasurementUnitConversionResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	created, err := s.validEnumerationsManager.CreateValidMeasurementUnitConversion(ctx, grpcconverters.ConvertGRPCCreateValidMeasurementUnitConversionRequestToValidMeasurementUnitConversionCreationRequestInput(request))
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "creating valid measurement unit conversion")
	}

	result := &messages.CreateValidMeasurementUnitConversionResponse{
		Result: grpcconverters.ConvertValidMeasurementUnitConversionToGRPCValidMeasurementUnitConversion(created),
	}

	return result, nil
}

func (s *serviceImpl) CreateValidPreparation(ctx context.Context, request *messages.CreateValidPreparationRequest) (*messages.CreateValidPreparationResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	created, err := s.validEnumerationsManager.CreateValidPreparation(ctx, grpcconverters.ConvertGRPCCreateValidPreparationRequestToValidPreparationCreationRequestInput(request))
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "creating valid preparation")
	}

	result := &messages.CreateValidPreparationResponse{
		Result: grpcconverters.ConvertValidPreparationToGRPCValidPreparation(created),
	}

	return result, nil
}

func (s *serviceImpl) CreateValidPreparationInstrument(ctx context.Context, request *messages.CreateValidPreparationInstrumentRequest) (*messages.CreateValidPreparationInstrumentResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	created, err := s.validEnumerationsManager.CreateValidPreparationInstrument(ctx, grpcconverters.ConvertGRPCCreateValidPreparationInstrumentRequestToValidPreparationInstrumentCreationRequestInput(request))
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "creating valid preparation instrument")
	}

	result := &messages.CreateValidPreparationInstrumentResponse{
		Result: grpcconverters.ConvertValidPreparationInstrumentToGRPCValidPreparationInstrument(created),
	}

	return result, nil
}

func (s *serviceImpl) CreateValidPreparationVessel(ctx context.Context, request *messages.CreateValidPreparationVesselRequest) (*messages.CreateValidPreparationVesselResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	created, err := s.validEnumerationsManager.CreateValidPreparationVessel(ctx, grpcconverters.ConvertGRPCCreateValidPreparationVesselRequestToValidPreparationVesselCreationRequestInput(request))
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "creating valid preparation vessel")
	}

	result := &messages.CreateValidPreparationVesselResponse{
		Result: grpcconverters.ConvertValidPreparationVesselToGRPCValidPreparationVessel(created),
	}

	return result, nil
}

func (s *serviceImpl) CreateValidVessel(ctx context.Context, request *messages.CreateValidVesselRequest) (*messages.CreateValidVesselResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	created, err := s.validEnumerationsManager.CreateValidVessel(ctx, grpcconverters.ConvertGRPCCreateValidVesselRequestToValidVesselCreationRequestInput(request))
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "creating valid vessel")
	}

	result := &messages.CreateValidVesselResponse{
		Result: grpcconverters.ConvertValidVesselToGRPCValidVessel(created),
	}

	return result, nil
}

func (s *serviceImpl) GetRandomValidIngredient(ctx context.Context, _ *messages.GetRandomValidIngredientRequest) (*messages.GetRandomValidIngredientResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	selected, err := s.validEnumerationsManager.RandomValidIngredient(ctx)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "fetching random valid ingredient")
	}

	res := &messages.GetRandomValidIngredientResponse{
		Result: grpcconverters.ConvertValidIngredientToGRPCValidIngredient(selected),
	}

	return res, nil
}

func (s *serviceImpl) GetRandomValidInstrument(ctx context.Context, _ *messages.GetRandomValidInstrumentRequest) (*messages.GetRandomValidInstrumentResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	selected, err := s.validEnumerationsManager.RandomValidInstrument(ctx)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "fetching random valid instrument")
	}

	res := &messages.GetRandomValidInstrumentResponse{
		Result: grpcconverters.ConvertValidInstrumentToGRPCValidInstrument(selected),
	}

	return res, nil
}

func (s *serviceImpl) GetRandomValidPreparation(ctx context.Context, _ *messages.GetRandomValidPreparationRequest) (*messages.GetRandomValidPreparationResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	selected, err := s.validEnumerationsManager.RandomValidPreparation(ctx)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "fetching random valid preparation")
	}

	res := &messages.GetRandomValidPreparationResponse{
		Result: grpcconverters.ConvertValidPreparationToGRPCValidPreparation(selected),
	}

	return res, nil
}

func (s *serviceImpl) GetRandomValidVessel(ctx context.Context, _ *messages.GetRandomValidVesselRequest) (*messages.GetRandomValidVesselResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	selected, err := s.validEnumerationsManager.RandomValidVessel(ctx)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "fetching random valid vessel")
	}

	res := &messages.GetRandomValidVesselResponse{
		Result: grpcconverters.ConvertValidVesselToGRPCValidVessel(selected),
	}

	return res, nil
}

func (s *serviceImpl) GetValidIngredient(ctx context.Context, request *messages.GetValidIngredientRequest) (*messages.GetValidIngredientResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	x, err := s.validEnumerationsManager.ReadValidIngredient(ctx, request.ValidIngredientID)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "fetching valid ingredient")
	}

	res := &messages.GetValidIngredientResponse{
		Result: grpcconverters.ConvertValidIngredientToGRPCValidIngredient(x),
	}

	return res, nil
}

func (s *serviceImpl) GetValidIngredientGroup(ctx context.Context, request *messages.GetValidIngredientGroupRequest) (*messages.GetValidIngredientGroupResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	x, err := s.validEnumerationsManager.ReadValidIngredientGroup(ctx, request.ValidIngredientGroupID)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "fetching valid ingredient group")
	}

	res := &messages.GetValidIngredientGroupResponse{
		Result: grpcconverters.ConvertValidIngredientGroupToGRPCValidIngredientGroup(x),
	}

	return res, nil
}

func (s *serviceImpl) GetValidIngredientGroups(ctx context.Context, request *messages.GetValidIngredientGroupsRequest) (*messages.GetValidIngredientGroupsResponse, error) {
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

func (s *serviceImpl) GetValidIngredientMeasurementUnit(ctx context.Context, request *messages.GetValidIngredientMeasurementUnitRequest) (*messages.GetValidIngredientMeasurementUnitResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	x, err := s.validEnumerationsManager.ReadValidIngredientMeasurementUnit(ctx, request.ValidIngredientMeasurementUnitID)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "fetching valid ingredient measurement unit")
	}

	res := &messages.GetValidIngredientMeasurementUnitResponse{
		Result: grpcconverters.ConvertValidIngredientMeasurementUnitToGRPCValidIngredientMeasurementUnit(x),
	}

	return res, nil
}

func (s *serviceImpl) GetValidIngredientMeasurementUnits(ctx context.Context, request *messages.GetValidIngredientMeasurementUnitsRequest) (*messages.GetValidIngredientMeasurementUnitsResponse, error) {
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

func (s *serviceImpl) GetValidIngredientMeasurementUnitsByIngredient(ctx context.Context, request *messages.GetValidIngredientMeasurementUnitsByIngredientRequest) (*messages.GetValidIngredientMeasurementUnitsByIngredientResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)
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

func (s *serviceImpl) GetValidIngredientMeasurementUnitsByMeasurementUnit(ctx context.Context, request *messages.GetValidIngredientMeasurementUnitsByMeasurementUnitRequest) (*messages.GetValidIngredientMeasurementUnitsByMeasurementUnitResponse, error) {
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

func (s *serviceImpl) GetValidIngredientPreparation(ctx context.Context, request *messages.GetValidIngredientPreparationRequest) (*messages.GetValidIngredientPreparationResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	x, err := s.validEnumerationsManager.ReadValidIngredientPreparation(ctx, request.ValidIngredientPreparationID)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "fetching valid ingredient preparation")
	}

	res := &messages.GetValidIngredientPreparationResponse{
		Result: grpcconverters.ConvertValidIngredientPreparationToGRPCValidIngredientPreparation(x),
	}

	return res, nil
}

func (s *serviceImpl) GetValidIngredientPreparations(ctx context.Context, request *messages.GetValidIngredientPreparationsRequest) (*messages.GetValidIngredientPreparationsResponse, error) {
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

func (s *serviceImpl) GetValidIngredientPreparationsByIngredient(ctx context.Context, request *messages.GetValidIngredientPreparationsByIngredientRequest) (*messages.GetValidIngredientPreparationsByIngredientResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)
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

func (s *serviceImpl) GetValidIngredientPreparationsByPreparation(ctx context.Context, request *messages.GetValidIngredientPreparationsByPreparationRequest) (*messages.GetValidIngredientPreparationsByPreparationResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)
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

func (s *serviceImpl) GetValidIngredientState(ctx context.Context, request *messages.GetValidIngredientStateRequest) (*messages.GetValidIngredientStateResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	x, err := s.validEnumerationsManager.ReadValidIngredientState(ctx, request.ValidIngredientStateID)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "fetching valid ingredient state")
	}

	res := &messages.GetValidIngredientStateResponse{
		Result: grpcconverters.ConvertValidIngredientStateToGRPCValidIngredientState(x),
	}

	return res, nil
}

func (s *serviceImpl) GetValidIngredientStateIngredient(ctx context.Context, request *messages.GetValidIngredientStateIngredientRequest) (*messages.GetValidIngredientStateIngredientResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	x, err := s.validEnumerationsManager.ReadValidIngredientStateIngredient(ctx, request.ValidIngredientStateIngredientID)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "fetching valid ingredient state ingredient")
	}

	res := &messages.GetValidIngredientStateIngredientResponse{
		Result: grpcconverters.ConvertValidIngredientStateIngredientToGRPCValidIngredientStateIngredient(x),
	}

	return res, nil
}

func (s *serviceImpl) GetValidIngredientStateIngredients(ctx context.Context, request *messages.GetValidIngredientStateIngredientsRequest) (*messages.GetValidIngredientStateIngredientsResponse, error) {
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
		res.Result = append(res.Result, grpcconverters.ConvertValidIngredientStateIngredientToGRPCValidIngredientStateIngredient(y))
	}

	return res, nil
}

func (s *serviceImpl) GetValidIngredientStateIngredientsByIngredient(ctx context.Context, request *messages.GetValidIngredientStateIngredientsByIngredientRequest) (*messages.GetValidIngredientStateIngredientsByIngredientResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)
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

func (s *serviceImpl) GetValidIngredientStateIngredientsByIngredientState(ctx context.Context, request *messages.GetValidIngredientStateIngredientsByIngredientStateRequest) (*messages.GetValidIngredientStateIngredientsByIngredientStateResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)
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

func (s *serviceImpl) GetValidIngredientStates(ctx context.Context, request *messages.GetValidIngredientStatesRequest) (*messages.GetValidIngredientStatesResponse, error) {
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

func (s *serviceImpl) GetValidIngredients(ctx context.Context, request *messages.GetValidIngredientsRequest) (*messages.GetValidIngredientsResponse, error) {
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

func (s *serviceImpl) GetValidInstrument(ctx context.Context, request *messages.GetValidInstrumentRequest) (*messages.GetValidInstrumentResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	x, err := s.validEnumerationsManager.ReadValidInstrument(ctx, request.ValidInstrumentID)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "fetching valid instrument")
	}

	res := &messages.GetValidInstrumentResponse{
		Result: grpcconverters.ConvertValidInstrumentToGRPCValidInstrument(x),
	}

	return res, nil
}

func (s *serviceImpl) GetValidInstruments(ctx context.Context, request *messages.GetValidInstrumentsRequest) (*messages.GetValidInstrumentsResponse, error) {
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

func (s *serviceImpl) GetValidMeasurementUnit(ctx context.Context, request *messages.GetValidMeasurementUnitRequest) (*messages.GetValidMeasurementUnitResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	x, err := s.validEnumerationsManager.ReadValidMeasurementUnit(ctx, request.ValidMeasurementUnitID)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "fetching valid measurement unit")
	}

	res := &messages.GetValidMeasurementUnitResponse{
		Result: grpcconverters.ConvertValidMeasurementUnitToGRPCValidMeasurementUnit(x),
	}

	return res, nil
}

func (s *serviceImpl) GetValidMeasurementUnitConversion(ctx context.Context, request *messages.GetValidMeasurementUnitConversionRequest) (*messages.GetValidMeasurementUnitConversionResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	x, err := s.validEnumerationsManager.ReadValidMeasurementUnitConversion(ctx, request.ValidMeasurementUnitConversionID)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "fetching valid measurement unit conversion")
	}

	res := &messages.GetValidMeasurementUnitConversionResponse{
		Result: grpcconverters.ConvertValidMeasurementUnitConversionToGRPCValidMeasurementUnitConversion(x),
	}

	return res, nil
}

func (s *serviceImpl) GetValidMeasurementUnitConversionsFromUnit(ctx context.Context, request *messages.GetValidMeasurementUnitConversionsFromUnitRequest) (*messages.GetValidMeasurementUnitConversionsFromUnitResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

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

func (s *serviceImpl) GetValidMeasurementUnitConversionsToUnit(ctx context.Context, request *messages.GetValidMeasurementUnitConversionsToUnitRequest) (*messages.GetValidMeasurementUnitConversionsToUnitResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

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

func (s *serviceImpl) GetValidMeasurementUnits(ctx context.Context, request *messages.GetValidMeasurementUnitsRequest) (*messages.GetValidMeasurementUnitsResponse, error) {
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

func (s *serviceImpl) GetValidPreparation(ctx context.Context, request *messages.GetValidPreparationRequest) (*messages.GetValidPreparationResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	x, err := s.validEnumerationsManager.ReadValidPreparation(ctx, request.ValidPreparationID)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "fetching valid preparation")
	}

	res := &messages.GetValidPreparationResponse{
		Result: grpcconverters.ConvertValidPreparationToGRPCValidPreparation(x),
	}

	return res, nil
}

func (s *serviceImpl) GetValidPreparationInstrument(ctx context.Context, request *messages.GetValidPreparationInstrumentRequest) (*messages.GetValidPreparationInstrumentResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	x, err := s.validEnumerationsManager.ReadValidPreparationInstrument(ctx, request.ValidPreparationInstrumentID)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "fetching valid preparation instrument")
	}

	res := &messages.GetValidPreparationInstrumentResponse{
		Result: grpcconverters.ConvertValidPreparationInstrumentToGRPCValidPreparationInstrument(x),
	}

	return res, nil
}

func (s *serviceImpl) GetValidPreparationInstruments(ctx context.Context, request *messages.GetValidPreparationInstrumentsRequest) (*messages.GetValidPreparationInstrumentsResponse, error) {
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

func (s *serviceImpl) GetValidPreparationInstrumentsByInstrument(ctx context.Context, request *messages.GetValidPreparationInstrumentsByInstrumentRequest) (*messages.GetValidPreparationInstrumentsByInstrumentResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)
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

func (s *serviceImpl) GetValidPreparationInstrumentsByPreparation(ctx context.Context, request *messages.GetValidPreparationInstrumentsByPreparationRequest) (*messages.GetValidPreparationInstrumentsByPreparationResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)
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
		res.Result = append(res.Result, grpcconverters.ConvertValidPreparationInstrumentToGRPCValidPreparationInstrument(y))
	}

	return res, nil
}

func (s *serviceImpl) GetValidPreparationVessel(ctx context.Context, request *messages.GetValidPreparationVesselRequest) (*messages.GetValidPreparationVesselResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	x, err := s.validEnumerationsManager.ReadValidPreparationVessel(ctx, request.ValidPreparationVesselID)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "fetching valid preparation vessel")
	}

	res := &messages.GetValidPreparationVesselResponse{
		Result: grpcconverters.ConvertValidPreparationVesselToGRPCValidPreparationVessel(x),
	}

	return res, nil
}

func (s *serviceImpl) GetValidPreparationVessels(ctx context.Context, request *messages.GetValidPreparationVesselsRequest) (*messages.GetValidPreparationVesselsResponse, error) {
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
		res.Result = append(res.Result, grpcconverters.ConvertValidPreparationVesselToGRPCValidPreparationVessel(y))
	}

	return res, nil
}

func (s *serviceImpl) GetValidPreparationVesselsByPreparation(ctx context.Context, request *messages.GetValidPreparationVesselsByPreparationRequest) (*messages.GetValidPreparationVesselsByPreparationResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)
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

func (s *serviceImpl) GetValidPreparationVesselsByVessel(ctx context.Context, request *messages.GetValidPreparationVesselsByVesselRequest) (*messages.GetValidPreparationVesselsByVesselResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)
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

func (s *serviceImpl) GetValidPreparations(ctx context.Context, request *messages.GetValidPreparationsRequest) (*messages.GetValidPreparationsResponse, error) {
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

func (s *serviceImpl) GetValidVessel(ctx context.Context, request *messages.GetValidVesselRequest) (*messages.GetValidVesselResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	x, err := s.validEnumerationsManager.ReadValidVessel(ctx, request.ValidVesselID)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "fetching valid vessel")
	}

	res := &messages.GetValidVesselResponse{
		Result: grpcconverters.ConvertValidVesselToGRPCValidVessel(x),
	}

	return res, nil
}

func (s *serviceImpl) GetValidVessels(ctx context.Context, request *messages.GetValidVesselsRequest) (*messages.GetValidVesselsResponse, error) {
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

func (s *serviceImpl) SearchForValidIngredientGroups(ctx context.Context, request *messages.SearchForValidIngredientGroupsRequest) (*messages.SearchForValidIngredientGroupsResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return &messages.SearchForValidIngredientGroupsResponse{}, observability.PrepareAndLogError(nil, logger, span, "")
}

func (s *serviceImpl) SearchForValidIngredientStates(ctx context.Context, request *messages.SearchForValidIngredientStatesRequest) (*messages.SearchForValidIngredientStatesResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return &messages.SearchForValidIngredientStatesResponse{}, observability.PrepareAndLogError(nil, logger, span, "")
}

func (s *serviceImpl) SearchForValidIngredients(ctx context.Context, request *messages.SearchForValidIngredientsRequest) (*messages.SearchForValidIngredientsResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return &messages.SearchForValidIngredientsResponse{}, observability.PrepareAndLogError(nil, logger, span, "")
}

func (s *serviceImpl) SearchForValidInstruments(ctx context.Context, request *messages.SearchForValidInstrumentsRequest) (*messages.SearchForValidInstrumentsResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return &messages.SearchForValidInstrumentsResponse{}, observability.PrepareAndLogError(nil, logger, span, "")
}

func (s *serviceImpl) SearchForValidMeasurementUnits(ctx context.Context, request *messages.SearchForValidMeasurementUnitsRequest) (*messages.SearchForValidMeasurementUnitsResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return &messages.SearchForValidMeasurementUnitsResponse{}, observability.PrepareAndLogError(nil, logger, span, "")
}

func (s *serviceImpl) SearchForValidPreparations(ctx context.Context, request *messages.SearchForValidPreparationsRequest) (*messages.SearchForValidPreparationsResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return &messages.SearchForValidPreparationsResponse{}, observability.PrepareAndLogError(nil, logger, span, "")
}

func (s *serviceImpl) SearchForValidVessels(ctx context.Context, request *messages.SearchForValidVesselsRequest) (*messages.SearchForValidVesselsResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return &messages.SearchForValidVesselsResponse{}, observability.PrepareAndLogError(nil, logger, span, "")
}

func (s *serviceImpl) SearchValidIngredientsByPreparation(ctx context.Context, request *messages.SearchValidIngredientsByPreparationRequest) (*messages.SearchValidIngredientsByPreparationResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return &messages.SearchValidIngredientsByPreparationResponse{}, observability.PrepareAndLogError(nil, logger, span, "")
}

func (s *serviceImpl) SearchValidMeasurementUnitsByIngredient(ctx context.Context, request *messages.SearchValidMeasurementUnitsByIngredientRequest) (*messages.SearchValidMeasurementUnitsByIngredientResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return &messages.SearchValidMeasurementUnitsByIngredientResponse{}, observability.PrepareAndLogError(nil, logger, span, "")
}

func (s *serviceImpl) UpdateValidIngredient(ctx context.Context, request *messages.UpdateValidIngredientRequest) (*messages.UpdateValidIngredientResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return &messages.UpdateValidIngredientResponse{}, observability.PrepareAndLogError(nil, logger, span, "")
}

func (s *serviceImpl) UpdateValidIngredientGroup(ctx context.Context, request *messages.UpdateValidIngredientGroupRequest) (*messages.UpdateValidIngredientGroupResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return &messages.UpdateValidIngredientGroupResponse{}, observability.PrepareAndLogError(nil, logger, span, "")
}

func (s *serviceImpl) UpdateValidIngredientMeasurementUnit(ctx context.Context, request *messages.UpdateValidIngredientMeasurementUnitRequest) (*messages.UpdateValidIngredientMeasurementUnitResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return &messages.UpdateValidIngredientMeasurementUnitResponse{}, observability.PrepareAndLogError(nil, logger, span, "")
}

func (s *serviceImpl) UpdateValidIngredientPreparation(ctx context.Context, request *messages.UpdateValidIngredientPreparationRequest) (*messages.UpdateValidIngredientPreparationResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return &messages.UpdateValidIngredientPreparationResponse{}, observability.PrepareAndLogError(nil, logger, span, "")
}

func (s *serviceImpl) UpdateValidIngredientState(ctx context.Context, request *messages.UpdateValidIngredientStateRequest) (*messages.UpdateValidIngredientStateResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return &messages.UpdateValidIngredientStateResponse{}, observability.PrepareAndLogError(nil, logger, span, "")
}

func (s *serviceImpl) UpdateValidIngredientStateIngredient(ctx context.Context, request *messages.UpdateValidIngredientStateIngredientRequest) (*messages.UpdateValidIngredientStateIngredientResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return &messages.UpdateValidIngredientStateIngredientResponse{}, observability.PrepareAndLogError(nil, logger, span, "")
}

func (s *serviceImpl) UpdateValidInstrument(ctx context.Context, request *messages.UpdateValidInstrumentRequest) (*messages.UpdateValidInstrumentResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return &messages.UpdateValidInstrumentResponse{}, observability.PrepareAndLogError(nil, logger, span, "")
}

func (s *serviceImpl) UpdateValidMeasurementUnit(ctx context.Context, request *messages.UpdateValidMeasurementUnitRequest) (*messages.UpdateValidMeasurementUnitResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return &messages.UpdateValidMeasurementUnitResponse{}, observability.PrepareAndLogError(nil, logger, span, "")
}

func (s *serviceImpl) UpdateValidMeasurementUnitConversion(ctx context.Context, request *messages.UpdateValidMeasurementUnitConversionRequest) (*messages.UpdateValidMeasurementUnitConversionResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return &messages.UpdateValidMeasurementUnitConversionResponse{}, observability.PrepareAndLogError(nil, logger, span, "")
}

func (s *serviceImpl) UpdateValidPreparation(ctx context.Context, request *messages.UpdateValidPreparationRequest) (*messages.UpdateValidPreparationResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return &messages.UpdateValidPreparationResponse{}, observability.PrepareAndLogError(nil, logger, span, "")
}

func (s *serviceImpl) UpdateValidPreparationInstrument(ctx context.Context, request *messages.UpdateValidPreparationInstrumentRequest) (*messages.UpdateValidPreparationInstrumentResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return &messages.UpdateValidPreparationInstrumentResponse{}, observability.PrepareAndLogError(nil, logger, span, "")
}

func (s *serviceImpl) UpdateValidPreparationVessel(ctx context.Context, request *messages.UpdateValidPreparationVesselRequest) (*messages.UpdateValidPreparationVesselResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return &messages.UpdateValidPreparationVesselResponse{}, observability.PrepareAndLogError(nil, logger, span, "")
}

func (s *serviceImpl) UpdateValidVessel(ctx context.Context, request *messages.UpdateValidVesselRequest) (*messages.UpdateValidVesselResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return &messages.UpdateValidVesselResponse{}, observability.PrepareAndLogError(nil, logger, span, "")
}
