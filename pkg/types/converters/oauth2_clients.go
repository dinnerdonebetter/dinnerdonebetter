package converters

import (
	"fmt"

	"github.com/dinnerdonebetter/backend/internal/identifiers"
	"github.com/dinnerdonebetter/backend/pkg/types"

	fake "github.com/brianvoe/gofakeit/v5"
)

func ConvertOAuth2ClientCreationRequestInputToOAuth2ClientDatabaseCreationInput(x *types.OAuth2ClientCreationRequestInput) *types.OAuth2ClientDatabaseCreationInput {
	return &types.OAuth2ClientDatabaseCreationInput{
		ID:           identifiers.New(),
		Name:         x.Name,
		Description:  x.Description,
		ClientID:     "",
		ClientSecret: "",
	}
}

// ConvertOAuth2ClientToOAuth2ClientDatabaseCreationInput builds a faked OAuth2ClientDatabaseCreationInput.
func ConvertOAuth2ClientToOAuth2ClientDatabaseCreationInput(client *types.OAuth2Client) *types.OAuth2ClientDatabaseCreationInput {
	return &types.OAuth2ClientDatabaseCreationInput{
		ID:           client.ID,
		Name:         client.Name,
		Description:  client.Description,
		ClientID:     client.ClientID,
		ClientSecret: client.ClientSecret,
	}
}

// ConvertOAuth2ClientToOAuth2ClientCreationInput builds a faked OAuth2ClientCreationRequestInput.
func ConvertOAuth2ClientToOAuth2ClientCreationInput(client *types.OAuth2Client) *types.OAuth2ClientCreationRequestInput {
	return &types.OAuth2ClientCreationRequestInput{
		UserLoginInput: types.UserLoginInput{
			Username:  fake.Username(),
			Password:  fake.UUID(),
			TOTPToken: fmt.Sprintf("0%s", fake.Zip()),
		},
		Name:        client.Name,
		Description: client.Description,
	}
}

// ConvertOAuth2ClientToOAuth2ClientCreationResponse builds a faked OAuth2ClientCreationRequestInput.
func ConvertOAuth2ClientToOAuth2ClientCreationResponse(client *types.OAuth2Client) *types.OAuth2ClientCreationResponse {
	return &types.OAuth2ClientCreationResponse{
		Name:        client.Name,
		Description: client.Description,
	}
}
