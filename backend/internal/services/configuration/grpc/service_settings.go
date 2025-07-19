package grpc

import (
	"context"
	configurationsvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/configuration"
	"github.com/dinnerdonebetter/backend/internal/grpc/generated/types"
)

func (s *ServiceImpl) GetServiceSettingConfigurationsForAccount(ctx context.Context, request *configurationsvc.GetServiceSettingConfigurationsForAccountRequest) (*configurationsvc.GetServiceSettingConfigurationsForAccountResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	x := &configurationsvc.GetServiceSettingConfigurationsForAccountResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
	}

	return x, nil
}

func (s *ServiceImpl) GetServiceSettingConfigurationsForUser(ctx context.Context, request *configurationsvc.GetServiceSettingConfigurationsForUserRequest) (*configurationsvc.GetServiceSettingConfigurationsForUserResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	x := &configurationsvc.GetServiceSettingConfigurationsForUserResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
	}

	return x, nil
}

func (s *ServiceImpl) SearchForServiceSettings(ctx context.Context, request *configurationsvc.SearchForServiceSettingsRequest) (*configurationsvc.SearchForServiceSettingsResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	x := &configurationsvc.SearchForServiceSettingsResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
	}

	return x, nil
}

func (s *ServiceImpl) ArchiveServiceSetting(ctx context.Context, request *configurationsvc.ArchiveServiceSettingRequest) (*configurationsvc.ArchiveServiceSettingResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	x := &configurationsvc.ArchiveServiceSettingResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
	}

	return x, nil
}

func (s *ServiceImpl) ArchiveServiceSettingConfiguration(ctx context.Context, request *configurationsvc.ArchiveServiceSettingConfigurationRequest) (*configurationsvc.ArchiveServiceSettingConfigurationResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	x := &configurationsvc.ArchiveServiceSettingConfigurationResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
	}

	return x, nil
}

func (s *ServiceImpl) CreateServiceSetting(ctx context.Context, request *configurationsvc.CreateServiceSettingRequest) (*configurationsvc.CreateServiceSettingResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	x := &configurationsvc.CreateServiceSettingResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
	}

	return x, nil
}

func (s *ServiceImpl) CreateServiceSettingConfiguration(ctx context.Context, request *configurationsvc.CreateServiceSettingConfigurationRequest) (*configurationsvc.CreateServiceSettingConfigurationResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	x := &configurationsvc.CreateServiceSettingConfigurationResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
	}

	return x, nil
}

func (s *ServiceImpl) GetServiceSetting(ctx context.Context, request *configurationsvc.GetServiceSettingRequest) (*configurationsvc.GetServiceSettingResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	x := &configurationsvc.GetServiceSettingResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
	}

	return x, nil
}

func (s *ServiceImpl) GetServiceSettingConfigurationByName(ctx context.Context, request *configurationsvc.GetServiceSettingConfigurationByNameRequest) (*configurationsvc.GetServiceSettingConfigurationByNameResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	x := &configurationsvc.GetServiceSettingConfigurationByNameResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
	}

	return x, nil
}

func (s *ServiceImpl) GetServiceSettings(ctx context.Context, request *configurationsvc.GetServiceSettingsRequest) (*configurationsvc.GetServiceSettingsResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	x := &configurationsvc.GetServiceSettingsResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
	}

	return x, nil
}

func (s *ServiceImpl) UpdateServiceSettingConfiguration(ctx context.Context, request *configurationsvc.UpdateServiceSettingConfigurationRequest) (*configurationsvc.UpdateServiceSettingConfigurationResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	x := &configurationsvc.UpdateServiceSettingConfigurationResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
	}

	return x, nil
}
