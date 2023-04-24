package types

import (
	"context"
	"encoding/gob"
	"net/http"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const (
	// ServiceSettingDataType indicates an event is related to a service setting.
	ServiceSettingDataType dataType = "service_setting"

	// ServiceSettingCreatedCustomerEventType indicates a service setting was created.
	ServiceSettingCreatedCustomerEventType CustomerEventType = "service_setting_created"
	// ServiceSettingUpdatedCustomerEventType indicates a service setting was updated.
	ServiceSettingUpdatedCustomerEventType CustomerEventType = "service_setting_updated"
	// ServiceSettingArchivedCustomerEventType indicates a service setting was archived.
	ServiceSettingArchivedCustomerEventType CustomerEventType = "service_setting_archived"
)

func init() {
	gob.Register(new(ServiceSetting))
	gob.Register(new(ServiceSettingCreationRequestInput))
	gob.Register(new(ServiceSettingUpdateRequestInput))
}

type (
	// ServiceSetting represents a service setting.
	ServiceSetting struct {
		_ struct{}

		CreatedAt     time.Time  `json:"createdAt"`
		DefaultValue  *string    `json:"defaultValue"`
		LastUpdatedAt *time.Time `json:"lastUpdatedAt"`
		ArchivedAt    *time.Time `json:"archivedAt"`
		ID            string     `json:"id"`
		Name          string     `json:"name"`
		Type          string     `json:"type"`
		Description   string     `json:"description"`
		AdminsOnly    bool       `json:"adminsOnly"`
	}

	// ConfiguredServiceSetting represents a configured service setting.
	ConfiguredServiceSetting struct {
		_ struct{}

		CreatedAt        time.Time  `json:"createdAt"`
		LastUpdatedAt    *time.Time `json:"lastUpdatedAt"`
		ArchivedAt       *time.Time `json:"archivedAt"`
		ID               string     `json:"id"`
		Value            string     `json:"value"`
		Notes            string     `json:"notes"`
		ServiceSettingID string     `json:"serviceSettingID"`
		UserID           string     `json:"userID"`
		HouseholdID      string     `json:"householdID"`
	}

	// ServiceSettingCreationRequestInput represents what a user could set as input for creating settings.
	ServiceSettingCreationRequestInput struct {
		_ struct{}

		CreatedAt     time.Time  `json:"createdAt"`
		DefaultValue  *string    `json:"defaultValue"`
		LastUpdatedAt *time.Time `json:"lastUpdatedAt"`
		ArchivedAt    *time.Time `json:"archivedAt"`
		ID            string     `json:"id"`
		Name          string     `json:"name"`
		Type          string     `json:"type"`
		Description   string     `json:"description"`
		AdminsOnly    bool       `json:"adminsOnly"`
	}

	// ConfiguredServiceSettingCreationRequestInput represents what a user could set as input for creating settings.
	ConfiguredServiceSettingCreationRequestInput struct {
		_ struct{}

		Value            string `json:"value"`
		Notes            string `json:"notes"`
		ServiceSettingID string `json:"serviceSettingID"`
		UserID           string `json:"userID"`
		HouseholdID      string `json:"householdID"`
	}

	// ServiceSettingDatabaseCreationInput represents what a user could set as input for creating service settings.
	ServiceSettingDatabaseCreationInput struct {
		_ struct{}

		DefaultValue *string
		ID           string
		Name         string
		Type         string
		Description  string
		AdminsOnly   bool
	}

	// ConfiguredServiceSettingDatabaseCreationInput represents what a user could set as input for creating service settings.
	ConfiguredServiceSettingDatabaseCreationInput struct {
		_ struct{}

		ID               string
		Value            string
		Notes            string
		ServiceSettingID string
		UserID           string
		HouseholdID      string
	}

	// ServiceSettingUpdateRequestInput represents what a user could set as input for updating service settings.
	ServiceSettingUpdateRequestInput struct {
		_ struct{}

		Name         *string `json:"name"`
		Type         *string `json:"type"`
		Description  *string `json:"description"`
		DefaultValue *string `json:"defaultValue"`
		AdminsOnly   *bool   `json:"adminsOnly"`
	}

	// ConfiguredServiceSettingUpdateRequestInput represents what a user could set as input for updating service settings.
	ConfiguredServiceSettingUpdateRequestInput struct {
		_ struct{}

		Value            *string `json:"value"`
		Notes            *string `json:"notes"`
		ServiceSettingID *string `json:"serviceSettingID"`
		UserID           *string `json:"userID"`
		HouseholdID      *string `json:"householdID"`
	}

	// ServiceSettingDataManager describes a structure capable of storing settings permanently.
	ServiceSettingDataManager interface {
		ServiceSettingExists(ctx context.Context, validPreparationID string) (bool, error)
		GetServiceSetting(ctx context.Context, validPreparationID string) (*ServiceSetting, error)
		GetRandomServiceSetting(ctx context.Context) (*ServiceSetting, error)
		GetServiceSettings(ctx context.Context, filter *QueryFilter) (*QueryFilteredResult[ServiceSetting], error)
		SearchForServiceSettings(ctx context.Context, query string) ([]*ServiceSetting, error)
		CreateServiceSetting(ctx context.Context, input *ServiceSettingDatabaseCreationInput) (*ServiceSetting, error)
		UpdateServiceSetting(ctx context.Context, updated *ServiceSetting) error
		ArchiveServiceSetting(ctx context.Context, validPreparationID string) error
	}

	// ServiceSettingDataService describes a structure capable of serving traffic related to service settings.
	ServiceSettingDataService interface {
		SearchHandler(res http.ResponseWriter, req *http.Request)
		ListHandler(res http.ResponseWriter, req *http.Request)
		CreateHandler(res http.ResponseWriter, req *http.Request)
		ReadHandler(res http.ResponseWriter, req *http.Request)
		RandomHandler(res http.ResponseWriter, req *http.Request)
		UpdateHandler(res http.ResponseWriter, req *http.Request)
		ArchiveHandler(res http.ResponseWriter, req *http.Request)
	}
)

// Update merges an ServiceSettingUpdateRequestInput with a service setting.
func (x *ServiceSetting) Update(input *ServiceSettingUpdateRequestInput) {
	if input.Name != nil && *input.Name != x.Name {
		x.Name = *input.Name
	}

	if input.Type != nil && *input.Type != x.Type {
		x.Type = *input.Type
	}

	if input.Description != nil && *input.Description != x.Description {
		x.Description = *input.Description
	}

	if input.DefaultValue != nil && input.DefaultValue != x.DefaultValue {
		x.DefaultValue = input.DefaultValue
	}

	if input.AdminsOnly != nil && *input.AdminsOnly != x.AdminsOnly {
		x.AdminsOnly = *input.AdminsOnly
	}
}

// Update merges an ConfiguredServiceSettingUpdateRequestInput with a service setting.
func (x *ConfiguredServiceSetting) Update(input *ConfiguredServiceSettingUpdateRequestInput) {
	if input.Value != nil && *input.Value != x.Value {
		x.Value = *input.Value
	}
	if input.Notes != nil && *input.Notes != x.Notes {
		x.Notes = *input.Notes
	}
	if input.ServiceSettingID != nil && *input.ServiceSettingID != x.ServiceSettingID {
		x.ServiceSettingID = *input.ServiceSettingID
	}
	if input.UserID != nil && *input.UserID != x.UserID {
		x.UserID = *input.UserID
	}
	if input.HouseholdID != nil && *input.HouseholdID != x.HouseholdID {
		x.HouseholdID = *input.HouseholdID
	}
}

var _ validation.ValidatableWithContext = (*ServiceSettingCreationRequestInput)(nil)

// ValidateWithContext validates a ServiceSettingCreationRequestInput.
func (x *ServiceSettingCreationRequestInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.Name, validation.Required),
		validation.Field(&x.Type, validation.Required),
	)
}

var _ validation.ValidatableWithContext = (*ConfiguredServiceSettingCreationRequestInput)(nil)

// ValidateWithContext validates a ServiceSettingCreationRequestInput.
func (x *ConfiguredServiceSettingCreationRequestInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.Value, validation.Required),
		validation.Field(&x.Notes, validation.Required),
		validation.Field(&x.UserID, validation.Required),
	)
}

var _ validation.ValidatableWithContext = (*ServiceSettingDatabaseCreationInput)(nil)

// ValidateWithContext validates a ServiceSettingDatabaseCreationInput.
func (x *ServiceSettingDatabaseCreationInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.ID, validation.Required),
		validation.Field(&x.Type, validation.Required),
		validation.Field(&x.Name, validation.Required),
	)
}

var _ validation.ValidatableWithContext = (*ConfiguredServiceSettingDatabaseCreationInput)(nil)

// ValidateWithContext validates a ConfiguredServiceSettingDatabaseCreationInput.
func (x *ConfiguredServiceSettingDatabaseCreationInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.ID, validation.Required),
		validation.Field(&x.Value, validation.Required),
		validation.Field(&x.UserID, validation.Required),
	)
}

var _ validation.ValidatableWithContext = (*ServiceSettingUpdateRequestInput)(nil)

// ValidateWithContext validates a ServiceSettingUpdateRequestInput.
func (x *ServiceSettingUpdateRequestInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.Name, validation.Required),
		validation.Field(&x.Description, validation.Required),
	)
}

var _ validation.ValidatableWithContext = (*ConfiguredServiceSettingUpdateRequestInput)(nil)

// ValidateWithContext validates a ConfiguredServiceSettingUpdateRequestInput.
func (x *ConfiguredServiceSettingUpdateRequestInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.Value, validation.Required),
		validation.Field(&x.UserID, validation.Required),
	)
}
