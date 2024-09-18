package types

import (
	"context"
	"encoding/gob"
	"net/http"
	"slices"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/hashicorp/go-multierror"
)

const (
	// ServiceSettingCreatedCustomerEventType indicates a service setting was created.
	ServiceSettingCreatedCustomerEventType ServiceEventType = "service_setting_created"
	// ServiceSettingArchivedCustomerEventType indicates a service setting was archived.
	ServiceSettingArchivedCustomerEventType ServiceEventType = "service_setting_archived"
)

func init() {
	gob.Register(new(ServiceSetting))
	gob.Register(new(ServiceSettingCreationRequestInput))
}

type (
	// ServiceSetting represents a service setting.
	ServiceSetting struct {
		_ struct{} `json:"-"`

		CreatedAt     time.Time  `json:"createdAt"`
		DefaultValue  *string    `json:"defaultValue"`
		LastUpdatedAt *time.Time `json:"lastUpdatedAt"`
		ArchivedAt    *time.Time `json:"archivedAt"`
		ID            string     `json:"id"`
		Name          string     `json:"name"`
		Type          string     `json:"type"`
		Description   string     `json:"description"`
		Enumeration   []string   `json:"enumeration"`
		AdminsOnly    bool       `json:"adminsOnly"`
	}

	// ServiceSettingCreationRequestInput represents what a user could set as input for creating settings.
	ServiceSettingCreationRequestInput struct {
		_ struct{} `json:"-"`

		DefaultValue *string  `json:"defaultValue"`
		Name         string   `json:"name"`
		Type         string   `json:"type"`
		Description  string   `json:"description"`
		Enumeration  []string `json:"enumeration"`
		AdminsOnly   bool     `json:"adminsOnly"`
	}

	// ServiceSettingDatabaseCreationInput represents what a user could set as input for creating service settings.
	ServiceSettingDatabaseCreationInput struct {
		_ struct{} `json:"-"`

		DefaultValue *string
		ID           string
		Name         string
		Type         string
		Description  string
		Enumeration  []string
		AdminsOnly   bool
	}

	// ServiceSettingDataManager describes a structure capable of storing settings permanently.
	ServiceSettingDataManager interface {
		CreateServiceSetting(ctx context.Context, input *ServiceSettingDatabaseCreationInput) (*ServiceSetting, error)
		ServiceSettingExists(ctx context.Context, serviceSettingID string) (bool, error)
		GetServiceSetting(ctx context.Context, serviceSettingID string) (*ServiceSetting, error)
		GetServiceSettings(ctx context.Context, filter *QueryFilter) (*QueryFilteredResult[ServiceSetting], error)
		SearchForServiceSettings(ctx context.Context, query string) ([]*ServiceSetting, error)
		ArchiveServiceSetting(ctx context.Context, serviceSettingID string) error
	}

	// ServiceSettingDataService describes a structure capable of serving traffic related to service settings.
	ServiceSettingDataService interface {
		CreateHandler(http.ResponseWriter, *http.Request)
		SearchHandler(http.ResponseWriter, *http.Request)
		ListHandler(http.ResponseWriter, *http.Request)
		ReadHandler(http.ResponseWriter, *http.Request)
		ArchiveHandler(http.ResponseWriter, *http.Request)
	}
)

var _ validation.ValidatableWithContext = (*ServiceSettingCreationRequestInput)(nil)

// ValidateWithContext validates a ServiceSettingCreationRequestInput.
func (x *ServiceSettingCreationRequestInput) ValidateWithContext(ctx context.Context) error {
	var result *multierror.Error

	if x.DefaultValue != nil && !slices.Contains(x.Enumeration, *x.DefaultValue) {
		result = multierror.Append(result, errDefaultValueMustBeEnumerationValue)
	}

	if err := validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.Name, validation.Required),
		validation.Field(&x.Type, validation.Required),
		validation.Field(&x.Description, validation.Required),
		validation.Field(&x.Enumeration, validation.When(x.DefaultValue != nil, validation.Required)),
		validation.Field(&x.DefaultValue, validation.When(len(x.Enumeration) != 0, validation.Required)),
	); err != nil {
		result = multierror.Append(result, err)
	}

	return result.ErrorOrNil()
}

var _ validation.ValidatableWithContext = (*ServiceSettingDatabaseCreationInput)(nil)

// ValidateWithContext validates a ServiceSettingDatabaseCreationInput.
func (x *ServiceSettingDatabaseCreationInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.ID, validation.Required),
		validation.Field(&x.Name, validation.Required),
		validation.Field(&x.Description, validation.Required),
		validation.Field(&x.Type, validation.Required),
		validation.Field(&x.Enumeration, validation.When(x.DefaultValue != nil, validation.Required)),
		validation.Field(&x.DefaultValue, validation.When(len(x.Enumeration) != 0, validation.Required)),
	)
}
