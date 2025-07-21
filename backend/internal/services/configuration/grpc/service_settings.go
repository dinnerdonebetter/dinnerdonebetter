package grpc

import (
	"context"

	grpcconverters "github.com/dinnerdonebetter/backend/internal/grpc/converters"
	configurationsvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/configuration"
	"github.com/dinnerdonebetter/backend/internal/grpc/generated/types"
	"github.com/dinnerdonebetter/backend/internal/platform/observability"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/services/configuration/grpc/converters"

	"google.golang.org/grpc/codes"
)

func (s *ServiceImpl) CreateServiceSetting(ctx context.Context, request *configurationsvc.CreateServiceSettingRequest) (*configurationsvc.CreateServiceSettingResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	created, err := s.serviceSettingsRepository.CreateServiceSetting(ctx, converters.ConvertGPRCServiceSettingCreationRequestInputToServiceSettingDatabaseCreationInput(request.Input))
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed to create service setting")
	}

	x := &configurationsvc.CreateServiceSettingResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
		Created: converters.ConvertServiceSettingToGRPCServiceSetting(created),
	}

	return x, nil
}

func (s *ServiceImpl) GetServiceSetting(ctx context.Context, request *configurationsvc.GetServiceSettingRequest) (*configurationsvc.GetServiceSettingResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValue(keys.ServiceSettingIDKey, request.ServiceSettingID)

	serviceSetting, err := s.serviceSettingsRepository.GetServiceSetting(ctx, request.ServiceSettingID)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed to retrieve service setting")
	}

	x := &configurationsvc.GetServiceSettingResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
		Result: converters.ConvertServiceSettingToGRPCServiceSetting(serviceSetting),
	}

	return x, nil
}

func (s *ServiceImpl) GetServiceSettings(ctx context.Context, request *configurationsvc.GetServiceSettingsRequest) (*configurationsvc.GetServiceSettingsResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	serviceSettings, err := s.serviceSettingsRepository.GetServiceSettings(ctx, grpcconverters.ConvertGRPCQueryFilterToQueryFilter(request.Filter))
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed to retrieve service settings")
	}

	x := &configurationsvc.GetServiceSettingsResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
	}

	for _, serviceSetting := range serviceSettings.Data {
		x.Results = append(x.Results, converters.ConvertServiceSettingToGRPCServiceSetting(serviceSetting))
	}

	return x, nil
}

func (s *ServiceImpl) SearchForServiceSettings(ctx context.Context, request *configurationsvc.SearchForServiceSettingsRequest) (*configurationsvc.SearchForServiceSettingsResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValue(keys.SearchQueryKey, request.Query)

	serviceSettings, err := s.serviceSettingsRepository.SearchForServiceSettings(ctx, request.Query)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed to retrieve service settings")
	}

	x := &configurationsvc.SearchForServiceSettingsResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
	}

	for _, serviceSetting := range serviceSettings {
		x.Results = append(x.Results, converters.ConvertServiceSettingToGRPCServiceSetting(serviceSetting))
	}

	return x, nil
}

func (s *ServiceImpl) ArchiveServiceSetting(ctx context.Context, request *configurationsvc.ArchiveServiceSettingRequest) (*configurationsvc.ArchiveServiceSettingResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValue(keys.ServiceSettingIDKey, request.ServiceSettingID)

	if err := s.serviceSettingsRepository.ArchiveServiceSetting(ctx, request.ServiceSettingID); err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed to archive service setting")
	}

	x := &configurationsvc.ArchiveServiceSettingResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
	}

	return x, nil
}
