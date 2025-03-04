package types

import (
	"context"
	"encoding/gob"
	"net/http"
	"time"

	"github.com/dinnerdonebetter/backend/internal/lib/database/filtering"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const (
	// ServiceSettingConfigurationCreatedServiceEventType indicates a service setting was created.
	ServiceSettingConfigurationCreatedServiceEventType = "service_setting_configuration_created"
	// ServiceSettingConfigurationUpdatedServiceEventType indicates a service setting was updated.
	ServiceSettingConfigurationUpdatedServiceEventType = "service_setting_configuration_updated"
	// ServiceSettingConfigurationArchivedServiceEventType indicates a service setting was archived.
	ServiceSettingConfigurationArchivedServiceEventType = "service_setting_configuration_archived"
)

func init() {
	gob.Register(new(ServiceSetting))
	gob.Register(new(ServiceSettingCreationRequestInput))
}

type (
	// ServiceSettingConfiguration represents a configured service setting configurations.
	ServiceSettingConfiguration struct {
		_ struct{} `json:"-"`

		CreatedAt          time.Time      `json:"createdAt"`
		LastUpdatedAt      *time.Time     `json:"lastUpdatedAt"`
		ArchivedAt         *time.Time     `json:"archivedAt"`
		ID                 string         `json:"id"`
		Value              string         `json:"value"`
		Notes              string         `json:"notes"`
		BelongsToUser      string         `json:"belongsToUser"`
		BelongsToHousehold string         `json:"belongsToHousehold"`
		ServiceSetting     ServiceSetting `json:"serviceSetting"`
	}

	// ServiceSettingConfigurationCreationRequestInput represents what a user could set as input for creating settings configurations.
	ServiceSettingConfigurationCreationRequestInput struct {
		_ struct{} `json:"-"`

		Value              string `json:"value"`
		Notes              string `json:"notes"`
		ServiceSettingID   string `json:"serviceSettingID"`
		BelongsToUser      string `json:"belongsToUser"`
		BelongsToHousehold string `json:"belongsToHousehold"`
	}

	// ServiceSettingConfigurationDatabaseCreationInput represents what a user could set as input for creating service settings configurations.
	ServiceSettingConfigurationDatabaseCreationInput struct {
		_ struct{} `json:"-"`

		ID                 string `json:"-"`
		Value              string `json:"-"`
		Notes              string `json:"-"`
		ServiceSettingID   string `json:"-"`
		BelongsToUser      string `json:"-"`
		BelongsToHousehold string `json:"-"`
	}

	// ServiceSettingConfigurationUpdateRequestInput represents what a user could set as input for updating service settings configurations.
	ServiceSettingConfigurationUpdateRequestInput struct {
		_ struct{} `json:"-"`

		Value              *string `json:"value"`
		Notes              *string `json:"notes"`
		ServiceSettingID   *string `json:"serviceSettingID"`
		BelongsToUser      *string `json:"belongsToUser"`
		BelongsToHousehold *string `json:"belongsToHousehold"`
	}

	// ServiceSettingConfigurationDataManager describes a structure capable of storing settings permanently.
	ServiceSettingConfigurationDataManager interface {
		ServiceSettingConfigurationExists(ctx context.Context, serviceSettingConfigurationID string) (bool, error)
		GetServiceSettingConfiguration(ctx context.Context, serviceSettingConfigurationID string) (*ServiceSettingConfiguration, error)
		GetServiceSettingConfigurationForUserByName(ctx context.Context, userID, serviceSettingConfigurationName string) (*ServiceSettingConfiguration, error)
		GetServiceSettingConfigurationForHouseholdByName(ctx context.Context, householdID, serviceSettingConfigurationName string) (*ServiceSettingConfiguration, error)
		GetServiceSettingConfigurationsForUser(ctx context.Context, userID string, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[ServiceSettingConfiguration], error)
		GetServiceSettingConfigurationsForHousehold(ctx context.Context, householdID string, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[ServiceSettingConfiguration], error)
		CreateServiceSettingConfiguration(ctx context.Context, input *ServiceSettingConfigurationDatabaseCreationInput) (*ServiceSettingConfiguration, error)
		UpdateServiceSettingConfiguration(ctx context.Context, updated *ServiceSettingConfiguration) error
		ArchiveServiceSettingConfiguration(ctx context.Context, serviceSettingConfigurationID string) error
	}

	// ServiceSettingConfigurationDataService describes a structure capable of serving traffic related to service settings.
	ServiceSettingConfigurationDataService interface {
		CreateServiceSettingConfigurationHandler(http.ResponseWriter, *http.Request)
		GetServiceSettingConfigurationsForUserHandler(http.ResponseWriter, *http.Request)
		GetServiceSettingConfigurationsForHouseholdHandler(http.ResponseWriter, *http.Request)
		GetServiceSettingConfigurationsForUserByNameHandler(http.ResponseWriter, *http.Request)
		UpdateServiceSettingConfigurationHandler(http.ResponseWriter, *http.Request)
		ArchiveServiceSettingConfigurationHandler(http.ResponseWriter, *http.Request)
	}
)

// Update merges an ServiceSettingConfigurationUpdateRequestInput with a service setting.
func (x *ServiceSettingConfiguration) Update(input *ServiceSettingConfigurationUpdateRequestInput) {
	if input.Value != nil && *input.Value != x.Value {
		x.Value = *input.Value
	}

	if input.Notes != nil && *input.Notes != x.Notes {
		x.Notes = *input.Notes
	}

	if input.ServiceSettingID != nil && *input.ServiceSettingID != x.ServiceSetting.ID {
		x.ServiceSetting.ID = *input.ServiceSettingID
	}

	if input.BelongsToUser != nil && *input.BelongsToUser != x.BelongsToUser {
		x.BelongsToUser = *input.BelongsToUser
	}

	if input.BelongsToHousehold != nil && *input.BelongsToHousehold != x.BelongsToHousehold {
		x.BelongsToHousehold = *input.BelongsToHousehold
	}
}

var _ validation.ValidatableWithContext = (*ServiceSettingConfigurationCreationRequestInput)(nil)

// ValidateWithContext validates a ServiceSettingCreationRequestInput.
func (x *ServiceSettingConfigurationCreationRequestInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.Value, validation.Required),
		validation.Field(&x.BelongsToUser, validation.Required),
	)
}

var _ validation.ValidatableWithContext = (*ServiceSettingConfigurationDatabaseCreationInput)(nil)

// ValidateWithContext validates a ServiceSettingConfigurationDatabaseCreationInput.
func (x *ServiceSettingConfigurationDatabaseCreationInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.ID, validation.Required),
		validation.Field(&x.Value, validation.Required),
		validation.Field(&x.BelongsToUser, validation.Required),
	)
}

var _ validation.ValidatableWithContext = (*ServiceSettingConfigurationUpdateRequestInput)(nil)

// ValidateWithContext validates a ServiceSettingConfigurationUpdateRequestInput.
func (x *ServiceSettingConfigurationUpdateRequestInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.Value, validation.Required),
		validation.Field(&x.ServiceSettingID, validation.Required),
	)
}
