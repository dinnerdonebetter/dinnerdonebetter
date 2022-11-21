package converters

import (
	"github.com/prixfixeco/backend/internal/identifiers"
	"github.com/prixfixeco/backend/pkg/types"
)

func ConvertAPIClientCreationRequestInputToAPIClientDatabaseCreationInput(x *types.APIClientCreationRequestInput) *types.APIClientDatabaseCreationInput {
	return &types.APIClientDatabaseCreationInput{
		ID: identifiers.New(),
		UserLoginInput: types.UserLoginInput{
			Username:  x.Username,
			Password:  x.Password,
			TOTPToken: x.TOTPToken,
		},
		Name:          x.Name,
		ClientID:      "",
		BelongsToUser: "",
		ClientSecret:  nil,
	}
}
