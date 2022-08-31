package types

import (
	"context"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const (
	// APIClientDataType reference API client events.
	APIClientDataType dataType = "api_client"

	// APIClientCreatedCustomerEventType indicates an API client was created.
	APIClientCreatedCustomerEventType CustomerEventType = "api_client_created"
	// APIClientUpdatedCustomerEventType indicates an API client was updated.
	APIClientUpdatedCustomerEventType CustomerEventType = "api_client_updated"
	// APIClientArchivedCustomerEventType indicates an API client was archived.
	APIClientArchivedCustomerEventType CustomerEventType = "api_client_archived"
)

type (
	// APIClient represents a user-authorized API client.
	APIClient struct {
		_ struct{}

		LastUpdatedAt *uint64 `json:"lastUpdatedAt"`
		ArchivedAt    *uint64 `json:"archivedAt"`
		Name          string  `json:"name"`
		ClientID      string  `json:"clientID"`
		ID            string  `json:"id"`
		BelongsToUser string  `json:"belongsToUser"`
		ClientSecret  []byte  `json:"-"`
		CreatedAt     uint64  `json:"createdAt"`
	}

	// APIClientList is a response struct containing a list of API clients.
	APIClientList struct {
		_ struct{}

		Clients []*APIClient `json:"data"`
		Pagination
	}

	// APIClientCreationRequestInput is a struct for use when creating API clients.
	APIClientCreationRequestInput struct {
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
		GetAPIClients(ctx context.Context, owneruserID string, filter *QueryFilter) (*APIClientList, error)
		CreateAPIClient(ctx context.Context, input *APIClientCreationRequestInput) (*APIClient, error)
		ArchiveAPIClient(ctx context.Context, clientID, ownerUserID string) error
	}

	// APIClientDataService describes a structure capable of serving traffic related to API clients.
	APIClientDataService interface {
		ListHandler(res http.ResponseWriter, req *http.Request)
		CreateHandler(res http.ResponseWriter, req *http.Request)
		ReadHandler(res http.ResponseWriter, req *http.Request)
		ArchiveHandler(res http.ResponseWriter, req *http.Request)
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
