package types

import (
	"context"
	"net/http"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const (
	// APIClientCreatedCustomerEventType indicates an API client was created.
	APIClientCreatedCustomerEventType CustomerEventType = "api_client_created"
	// APIClientArchivedCustomerEventType indicates an API client was archived.
	APIClientArchivedCustomerEventType CustomerEventType = "api_client_archived"
)

type (
	// APIClient represents a user-authorized API client.
	APIClient struct {
		_ struct{}

		CreatedAt     time.Time  `json:"createdAt"`
		LastUpdatedAt *time.Time `json:"lastUpdatedAt"`
		ArchivedAt    *time.Time `json:"archivedAt"`
		Name          string     `json:"name"`
		ClientID      string     `json:"clientID"`
		ID            string     `json:"id"`
		BelongsToUser string     `json:"belongsToUser"`
		ClientSecret  []byte     `json:"-"`
	}

	// APIClientCreationRequestInput is a struct for use when creating API clients.
	APIClientCreationRequestInput struct {
		_ struct{}

		UserLoginInput
		Name string `json:"clientName"`
	}

	// APIClientDatabaseCreationInput is a struct for use when creating API clients.
	APIClientDatabaseCreationInput struct {
		_ struct{}

		UserLoginInput
		ID            string `json:"-"`
		Name          string `json:"clientName"`
		ClientID      string `json:"-"`
		BelongsToUser string `json:"-"`
		ClientSecret  []byte `json:"-"`
	}

	// APIClientCreationResponse is a struct for informing users of what their API client's secret key is.
	APIClientCreationResponse struct {
		_ struct{}

		ClientID     string `json:"clientID"`
		ClientSecret string `json:"clientSecret"`
		ID           string `json:"id"`
	}

	// APIClientDataManager handles API clients.
	APIClientDataManager interface {
		GetAPIClientByClientID(ctx context.Context, clientID string) (*APIClient, error)
		GetAPIClientByDatabaseID(ctx context.Context, clientID, ownerUserID string) (*APIClient, error)
		GetAPIClients(ctx context.Context, owneruserID string, filter *QueryFilter) (*QueryFilteredResult[APIClient], error)
		CreateAPIClient(ctx context.Context, input *APIClientDatabaseCreationInput) (*APIClient, error)
		ArchiveAPIClient(ctx context.Context, clientID, ownerUserID string) error
	}

	// APIClientDataService describes a structure capable of serving traffic related to API clients.
	APIClientDataService interface {
		ListHandler(http.ResponseWriter, *http.Request)
		CreateHandler(http.ResponseWriter, *http.Request)
		ReadHandler(http.ResponseWriter, *http.Request)
		ArchiveHandler(http.ResponseWriter, *http.Request)
	}
)

// ValidateWithContext validates an APICreationInput.
func (x *APIClientCreationRequestInput) ValidateWithContext(ctx context.Context, minUsernameLength, minPasswordLength uint8) error {
	if err := x.UserLoginInput.ValidateWithContext(ctx, minUsernameLength, minPasswordLength); err != nil {
		return err
	}

	return validation.ValidateStructWithContext(ctx, x,
		validation.Field(&x.Name, validation.Required),
	)
}
