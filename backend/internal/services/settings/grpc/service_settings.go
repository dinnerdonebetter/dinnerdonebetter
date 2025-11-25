package grpc

import (
	"context"

	grpcconverters "github.com/dinnerdonebetter/backend/internal/grpc/converters"
	settingssvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/settings"
	"github.com/dinnerdonebetter/backend/internal/grpc/generated/types"
	"github.com/dinnerdonebetter/backend/internal/platform/observability"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/services/settings/grpc/converters"

	"google.golang.org/grpc/codes"
)

func (s *serviceImpl) CreateServiceSetting(ctx context.Context, request *settingssvc.CreateServiceSettingRequest) (*settingssvc.CreateServiceSettingResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	input := converters.ConvertGPRCServiceSettingCreationRequestInputToServiceSettingDatabaseCreationInput(request.Input)
	if err := input.ValidateWithContext(ctx); err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.InvalidArgument, "invalid service setting")
	}

	created, err := s.serviceSettingsRepository.CreateServiceSetting(ctx, input)
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

	filter := grpcconverters.ConvertGRPCQueryFilterToQueryFilter(request.Filter)
	serviceSettings, err := s.serviceSettingsRepository.GetServiceSettings(ctx, filter)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed to retrieve service settings")
	}

	x := &settingssvc.GetServiceSettingsResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
		Pagination: grpcconverters.ConvertPaginationToGRPCPagination(serviceSettings.Pagination, filter),
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

	filter := grpcconverters.ConvertGRPCQueryFilterToQueryFilter(request.Filter)
	tracing.AttachQueryFilterToSpan(span, filter)

	serviceSettings, err := s.serviceSettingsRepository.SearchForServiceSettings(ctx, request.Query, filter)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed to retrieve service settings")
	}

	x := &settingssvc.SearchForServiceSettingsResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
		Pagination: grpcconverters.ConvertPaginationToGRPCPagination(serviceSettings.Pagination, filter),
	}

	for _, serviceSetting := range serviceSettings.Data {
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
