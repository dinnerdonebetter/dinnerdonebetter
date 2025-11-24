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

func (s *serviceImpl) CreateServiceSettingConfiguration(ctx context.Context, request *settingssvc.CreateServiceSettingConfigurationRequest) (*settingssvc.CreateServiceSettingConfigurationResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValue(keys.ServiceSettingIDKey, request.Input.ServiceSettingID)

	sessionContextData, err := s.sessionContextDataFetcher(ctx)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Unauthenticated, "failed to fetch session context data")
	}

	input := converters.ConvertGRPCServiceSettingConfigurationCreationRequestInputToServiceSettingConfigurationDatabaseCreationInput(request.Input, sessionContextData.GetUserID(), sessionContextData.GetActiveAccountID())
	input.BelongsToUser = sessionContextData.GetUserID()
	input.BelongsToAccount = sessionContextData.GetActiveAccountID()

	if err = input.ValidateWithContext(ctx); err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.InvalidArgument, "invalid service setting configuration")
	}

	created, err := s.serviceSettingsRepository.CreateServiceSettingConfiguration(ctx, input)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed to create service setting")
	}

	x := &settingssvc.CreateServiceSettingConfigurationResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
		Created: converters.ConvertServiceSettingConfigurationToGRPCServiceSettingConfiguration(created),
	}

	return x, nil
}

func (s *serviceImpl) GetServiceSettingConfigurationByName(ctx context.Context, request *settingssvc.GetServiceSettingConfigurationByNameRequest) (*settingssvc.GetServiceSettingConfigurationByNameResponse, error) {
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

	x := &settingssvc.GetServiceSettingConfigurationByNameResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
		Result: converters.ConvertServiceSettingConfigurationToGRPCServiceSettingConfiguration(serviceSettingConfig),
	}

	return x, nil
}

func (s *serviceImpl) GetServiceSettingConfigurationsForAccount(ctx context.Context, request *settingssvc.GetServiceSettingConfigurationsForAccountRequest) (*settingssvc.GetServiceSettingConfigurationsForAccountResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	sessionContextData, err := s.sessionContextDataFetcher(ctx)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Unauthenticated, "failed to fetch session context data")
	}
	logger = logger.WithValue(keys.AccountIDKey, sessionContextData.ActiveAccountID)

	filter := grpcconverters.ConvertGRPCQueryFilterToQueryFilter(request.Filter)

	serviceSettingConfigs, err := s.serviceSettingsRepository.GetServiceSettingConfigurationsForAccount(ctx, sessionContextData.ActiveAccountID, filter)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed to retrieve service setting configurations for account")
	}

	x := &settingssvc.GetServiceSettingConfigurationsForAccountResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
	}

	for _, cfg := range serviceSettingConfigs.Data {
		x.Results = append(x.Results, converters.ConvertServiceSettingConfigurationToGRPCServiceSettingConfiguration(cfg))
	}

	return x, nil
}

func (s *serviceImpl) GetServiceSettingConfigurationsForUser(ctx context.Context, request *settingssvc.GetServiceSettingConfigurationsForUserRequest) (*settingssvc.GetServiceSettingConfigurationsForUserResponse, error) {
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

	x := &settingssvc.GetServiceSettingConfigurationsForUserResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
	}

	for _, cfg := range serviceSettingConfigs.Data {
		x.Results = append(x.Results, converters.ConvertServiceSettingConfigurationToGRPCServiceSettingConfiguration(cfg))
	}

	return x, nil
}

func (s *serviceImpl) UpdateServiceSettingConfiguration(ctx context.Context, request *settingssvc.UpdateServiceSettingConfigurationRequest) (*settingssvc.UpdateServiceSettingConfigurationResponse, error) {
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

	x := &settingssvc.UpdateServiceSettingConfigurationResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
		Updated: converters.ConvertServiceSettingConfigurationToGRPCServiceSettingConfiguration(existing),
	}

	return x, nil
}

func (s *serviceImpl) ArchiveServiceSettingConfiguration(ctx context.Context, request *settingssvc.ArchiveServiceSettingConfigurationRequest) (*settingssvc.ArchiveServiceSettingConfigurationResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValue(keys.ServiceSettingConfigurationIDKey, request.ServiceSettingConfigurationID)

	if err := s.serviceSettingsRepository.ArchiveServiceSettingConfiguration(ctx, request.ServiceSettingConfigurationID); err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed to archive service setting configuration")
	}

	x := &settingssvc.ArchiveServiceSettingConfigurationResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
	}

	return x, nil
}
