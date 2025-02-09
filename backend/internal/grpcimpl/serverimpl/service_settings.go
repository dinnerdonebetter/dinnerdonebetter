package serverimpl

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/grpc/messages"
)

func (s *Server) ArchiveServiceSetting(ctx context.Context, request *messages.ArchiveServiceSettingRequest) (*messages.ArchiveServiceSettingResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) ArchiveServiceSettingConfiguration(ctx context.Context, request *messages.ArchiveServiceSettingConfigurationRequest) (*messages.ArchiveServiceSettingConfigurationResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) CreateServiceSetting(ctx context.Context, request *messages.CreateServiceSettingRequest) (*messages.CreateServiceSettingResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) CreateServiceSettingConfiguration(ctx context.Context, request *messages.CreateServiceSettingConfigurationRequest) (*messages.CreateServiceSettingConfigurationResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetServiceSetting(ctx context.Context, request *messages.GetServiceSettingRequest) (*messages.GetServiceSettingResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetServiceSettingConfigurationByName(ctx context.Context, request *messages.GetServiceSettingConfigurationByNameRequest) (*messages.GetServiceSettingConfigurationByNameResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetServiceSettingConfigurationsForHousehold(ctx context.Context, request *messages.GetServiceSettingConfigurationsForHouseholdRequest) (*messages.GetServiceSettingConfigurationsForHouseholdResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetServiceSettingConfigurationsForUser(ctx context.Context, request *messages.GetServiceSettingConfigurationsForUserRequest) (*messages.GetServiceSettingConfigurationsForUserResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetServiceSettings(ctx context.Context, request *messages.GetServiceSettingsRequest) (*messages.GetServiceSettingsResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) UpdateServiceSettingConfiguration(ctx context.Context, request *messages.UpdateServiceSettingConfigurationRequest) (*messages.UpdateServiceSettingConfigurationResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) SearchForServiceSettings(ctx context.Context, request *messages.SearchForServiceSettingsRequest) (*messages.SearchForServiceSettingsResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}
