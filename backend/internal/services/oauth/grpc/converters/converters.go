package converters

import (
	"github.com/dinnerdonebetter/backend/internal/domain/oauth"
	grpcconverters "github.com/dinnerdonebetter/backend/internal/grpc/converters"
	oauthsvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/oauth"
)

func ConvertGRPCCreateOAuth2ClientRequestToOAuth2ClientCreationRequestInput(input *oauthsvc.CreateOAuth2ClientRequest) *oauth.OAuth2ClientCreationRequestInput {
	return &oauth.OAuth2ClientCreationRequestInput{
		Name:        input.Name,
		Description: input.Description,
	}
}

func ConvertOAuth2ClientToGRPCOAuth2Client(client *oauth.OAuth2Client) *oauthsvc.OAuth2Client {
	return &oauthsvc.OAuth2Client{
		CreatedAt:    grpcconverters.ConvertTimeToPBTimestamp(client.CreatedAt),
		ArchivedAt:   grpcconverters.ConvertTimePointerToPBTimestamp(client.ArchivedAt),
		Name:         client.Name,
		Description:  client.Description,
		ClientID:     client.ClientID,
		ID:           client.ID,
		ClientSecret: client.ClientSecret,
	}
}
