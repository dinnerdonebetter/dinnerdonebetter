package grpc

import (
	"context"

	grpcconverters "github.com/dinnerdonebetter/backend/internal/grpc/converters"
	settingssvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/settings"
	"github.com/dinnerdonebetter/backend/internal/grpc/generated/types"
	"github.com/dinnerdonebetter/backend/internal/platform/observability"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/services/settings/grpc/converters"

	"google.golang.org/grpc/codes"
)

func (s *serviceImpl) CreateServiceSetting(ctx context.Context, request *settingssvc.CreateServiceSettingRequest) (*settingssvc.CreateServiceSettingResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	created, err := s.serviceSettingsRepository.CreateServiceSetting(ctx, converters.ConvertGPRCServiceSettingCreationRequestInputToServiceSettingDatabaseCreationInput(request.Input))
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed to create service setting")
	}

	x := &settingssvc.CreateServiceSettingResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
		Created: converters.ConvertServiceSettingToGRPCServiceSetting(created),
	}

	return x, nil
}

func (s *serviceImpl) GetServiceSetting(ctx context.Context, request *settingssvc.GetServiceSettingRequest) (*settingssvc.GetServiceSettingResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValue(keys.ServiceSettingIDKey, request.ServiceSettingID)

	serviceSetting, err := s.serviceSettingsRepository.GetServiceSetting(ctx, request.ServiceSettingID)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed to retrieve service setting")
	}

	x := &settingssvc.GetServiceSettingResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
		Result: converters.ConvertServiceSettingToGRPCServiceSetting(serviceSetting),
	}

	return x, nil
}

func (s *serviceImpl) GetServiceSettings(ctx context.Context, request *settingssvc.GetServiceSettingsRequest) (*settingssvc.GetServiceSettingsResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	serviceSettings, err := s.serviceSettingsRepository.GetServiceSettings(ctx, grpcconverters.ConvertGRPCQueryFilterToQueryFilter(request.Filter))
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed to retrieve service settings")
	}

	x := &settingssvc.GetServiceSettingsResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
	}

	for _, serviceSetting := range serviceSettings.Data {
		x.Results = append(x.Results, converters.ConvertServiceSettingToGRPCServiceSetting(serviceSetting))
	}

	return x, nil
}

func (s *serviceImpl) SearchForServiceSettings(ctx context.Context, request *settingssvc.SearchForServiceSettingsRequest) (*settingssvc.SearchForServiceSettingsResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValue(keys.SearchQueryKey, request.Query)

	serviceSettings, err := s.serviceSettingsRepository.SearchForServiceSettings(ctx, request.Query)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed to retrieve service settings")
	}

	x := &settingssvc.SearchForServiceSettingsResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
	}

	for _, serviceSetting := range serviceSettings {
		x.Results = append(x.Results, converters.ConvertServiceSettingToGRPCServiceSetting(serviceSetting))
	}

	return x, nil
}

func (s *serviceImpl) ArchiveServiceSetting(ctx context.Context, request *settingssvc.ArchiveServiceSettingRequest) (*settingssvc.ArchiveServiceSettingResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValue(keys.ServiceSettingIDKey, request.ServiceSettingID)

	if err := s.serviceSettingsRepository.ArchiveServiceSetting(ctx, request.ServiceSettingID); err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed to archive service setting")
	}

	x := &settingssvc.ArchiveServiceSettingResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
	}

	return x, nil
}
