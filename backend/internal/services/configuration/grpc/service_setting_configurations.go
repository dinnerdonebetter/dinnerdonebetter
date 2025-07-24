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

func (s *serviceImpl) CreateServiceSettingConfiguration(ctx context.Context, request *configurationsvc.CreateServiceSettingConfigurationRequest) (*configurationsvc.CreateServiceSettingConfigurationResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValue(keys.ServiceSettingIDKey, request.Input.ServiceSettingID)

	created, err := s.serviceSettingsRepository.CreateServiceSettingConfiguration(ctx, converters.ConvertGRPCServiceSettingConfigurationCreationRequestInputToServiceSettingConfigurationDatabaseCreationInput(request.Input))
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed to create service setting")
	}

	x := &configurationsvc.CreateServiceSettingConfigurationResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
		Created: converters.ConvertServiceSettingConfigurationToGRPCServiceSettingConfiguration(created),
	}

	return x, nil
}

func (s *serviceImpl) GetServiceSettingConfigurationByName(ctx context.Context, request *configurationsvc.GetServiceSettingConfigurationByNameRequest) (*configurationsvc.GetServiceSettingConfigurationByNameResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	sessionContextData, err := s.sessionContextDataFetcher(ctx)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Unauthenticated, "failed to fetch session context data")
	}
	logger = logger.WithValue(keys.AccountIDKey, sessionContextData.ActiveAccountID)

	serviceSettingConfig, err := s.serviceSettingsRepository.GetServiceSettingConfigurationForAccountByName(ctx, sessionContextData.ActiveAccountID, request.ServiceSettingConfigurationName)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed to retrieve service setting configuration for account by name")
	}

	x := &configurationsvc.GetServiceSettingConfigurationByNameResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
		Result: converters.ConvertServiceSettingConfigurationToGRPCServiceSettingConfiguration(serviceSettingConfig),
	}

	return x, nil
}

func (s *serviceImpl) GetServiceSettingConfigurationsForAccount(ctx context.Context, request *configurationsvc.GetServiceSettingConfigurationsForAccountRequest) (*configurationsvc.GetServiceSettingConfigurationsForAccountResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	sessionContextData, err := s.sessionContextDataFetcher(ctx)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Unauthenticated, "failed to fetch session context data")
	}
	logger = logger.WithValue(keys.AccountIDKey, sessionContextData.ActiveAccountID)

	serviceSettingConfigs, err := s.serviceSettingsRepository.GetServiceSettingConfigurationsForAccount(ctx, sessionContextData.ActiveAccountID, grpcconverters.ConvertGRPCQueryFilterToQueryFilter(request.Filter))
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed to retrieve service setting configurations for account")
	}

	x := &configurationsvc.GetServiceSettingConfigurationsForAccountResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
	}

	for _, cfg := range serviceSettingConfigs.Data {
		x.Results = append(x.Results, converters.ConvertServiceSettingConfigurationToGRPCServiceSettingConfiguration(cfg))
	}

	return x, nil
}

func (s *serviceImpl) GetServiceSettingConfigurationsForUser(ctx context.Context, request *configurationsvc.GetServiceSettingConfigurationsForUserRequest) (*configurationsvc.GetServiceSettingConfigurationsForUserResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	sessionContextData, err := s.sessionContextDataFetcher(ctx)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Unauthenticated, "failed to fetch session context data")
	}
	logger = logger.WithValue(keys.UserIDKey, sessionContextData.GetUserID())

	serviceSettingConfigs, err := s.serviceSettingsRepository.GetServiceSettingConfigurationsForUser(ctx, sessionContextData.GetUserID(), grpcconverters.ConvertGRPCQueryFilterToQueryFilter(request.Filter))
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed to retrieve service setting configurations for user")
	}

	x := &configurationsvc.GetServiceSettingConfigurationsForUserResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
	}

	for _, cfg := range serviceSettingConfigs.Data {
		x.Results = append(x.Results, converters.ConvertServiceSettingConfigurationToGRPCServiceSettingConfiguration(cfg))
	}

	return x, nil
}

func (s *serviceImpl) UpdateServiceSettingConfiguration(ctx context.Context, request *configurationsvc.UpdateServiceSettingConfigurationRequest) (*configurationsvc.UpdateServiceSettingConfigurationResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValue(keys.ServiceSettingConfigurationIDKey, request.ServiceSettingConfigurationID)

	existing, err := s.serviceSettingsRepository.GetServiceSettingConfiguration(ctx, request.ServiceSettingConfigurationID)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed to retrieve service setting configuration")
	}

	existing.Update(converters.ConvertGRPCServiceSettingConfigurationUpdateRequestInputToServiceSettingConfigurationUpdateRequestInputTo(request.Input))

	if err = s.serviceSettingsRepository.UpdateServiceSettingConfiguration(ctx, existing); err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed to update service setting configuration")
	}

	x := &configurationsvc.UpdateServiceSettingConfigurationResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
		Updated: converters.ConvertServiceSettingConfigurationToGRPCServiceSettingConfiguration(existing),
	}

	return x, nil
}

func (s *serviceImpl) ArchiveServiceSettingConfiguration(ctx context.Context, request *configurationsvc.ArchiveServiceSettingConfigurationRequest) (*configurationsvc.ArchiveServiceSettingConfigurationResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValue(keys.ServiceSettingConfigurationIDKey, request.ServiceSettingConfigurationID)

	if err := s.serviceSettingsRepository.ArchiveServiceSettingConfiguration(ctx, request.ServiceSettingConfigurationID); err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed to archive service setting configuration")
	}

	x := &configurationsvc.ArchiveServiceSettingConfigurationResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
	}

	return x, nil
}
