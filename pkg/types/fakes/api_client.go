package fakes

import (
	"fmt"

	fake "github.com/brianvoe/gofakeit/v5"
	"github.com/segmentio/ksuid"

	"github.com/prixfixeco/api_server/pkg/types"
)

// BuildFakeAPIClient builds a faked APIClient.
func BuildFakeAPIClient() *types.APIClient {
	return &types.APIClient{
		ID:            ksuid.New().String(),
		Name:          fake.Password(true, true, true, false, false, 32),
		ClientID:      ksuid.New().String(),
		ClientSecret:  []byte(fake.Password(true, true, true, true, false, 32)),
		BelongsToUser: fake.UUID(),
		CreatedAt:     fake.Date(),
	}
}

// BuildFakeAPIClientCreationResponseFromClient builds a faked APIClientCreationResponse.
func BuildFakeAPIClientCreationResponseFromClient(client *types.APIClient) *types.APIClientCreationResponse {
	return &types.APIClientCreationResponse{
		ID:           client.ID,
		ClientID:     client.ClientID,
		ClientSecret: string(client.ClientSecret),
	}
}

// BuildFakeAPIClientList builds a faked APIClientList.
func BuildFakeAPIClientList() *types.APIClientList {
	var examples []*types.APIClient
	for i := 0; i < exampleQuantity; i++ {
		examples = append(examples, BuildFakeAPIClient())
	}

	return &types.APIClientList{
		Pagination: types.Pagination{
			Page:          1,
			Limit:         20,
			FilteredCount: exampleQuantity / 2,
			TotalCount:    exampleQuantity,
		},
		Clients: examples,
	}
}

// BuildFakeAPIClientCreationInput builds a faked APIClientCreationRequestInput.
func BuildFakeAPIClientCreationInput() *types.APIClientCreationRequestInput {
	client := BuildFakeAPIClient()

	return &types.APIClientCreationRequestInput{
		UserLoginInput: types.UserLoginInput{
			Username:  fake.Username(),
			Password:  fake.Password(true, true, true, true, false, 32),
			TOTPToken: fmt.Sprintf("0%s", fake.Zip()),
		},
		Name:          client.Name,
		ClientID:      client.ClientID,
		BelongsToUser: client.BelongsToUser,
	}
}

// BuildFakeAPIClientCreationInputFromClient builds a faked APIClientCreationRequestInput.
func BuildFakeAPIClientCreationInputFromClient(client *types.APIClient) *types.APIClientCreationRequestInput {
	return &types.APIClientCreationRequestInput{
		ID: client.ID,
		UserLoginInput: types.UserLoginInput{
			Username:  fake.Username(),
			Password:  fake.Password(true, true, true, true, false, 32),
			TOTPToken: fmt.Sprintf("0%s", fake.Zip()),
		},
		Name:          client.Name,
		ClientID:      client.ClientID,
		ClientSecret:  client.ClientSecret,
		BelongsToUser: client.BelongsToUser,
	}
}
