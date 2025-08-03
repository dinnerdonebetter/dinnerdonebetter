package fakes

import (
	"encoding/base32"
	"fmt"

	"github.com/dinnerdonebetter/backend/internal/authorization"
	types "github.com/dinnerdonebetter/backend/internal/domain/identity"
	"github.com/dinnerdonebetter/backend/internal/platform/database/filtering"
	"github.com/dinnerdonebetter/backend/internal/platform/pointer"

	fake "github.com/brianvoe/gofakeit/v7"
)

// BuildFakeUser builds a faked User.
func BuildFakeUser() *types.User {
	fakeDate := BuildFakeTime()

	return &types.User{
		ID:                        BuildFakeID(),
		FirstName:                 fake.FirstName(),
		LastName:                  fake.LastName(),
		EmailAddress:              fake.Email(),
		Username:                  fmt.Sprintf("%s_%d_%s", fake.Username(), fake.Uint8(), fake.Username()),
		Birthday:                  pointer.To(BuildFakeTime()),
		AccountStatus:             string(types.UnverifiedAccountStatus),
		TwoFactorSecret:           base32.StdEncoding.EncodeToString([]byte(fake.Password(false, true, true, false, false, 32))),
		TwoFactorSecretVerifiedAt: &fakeDate,
		ServiceRole:               authorization.ServiceUserRole.String(),
		CreatedAt:                 BuildFakeTime(),
	}
}

// BuildFakeUsersList builds a faked UserList.
func BuildFakeUsersList() *filtering.QueryFilteredResult[types.User] {
	var examples []*types.User
	for i := 0; i < exampleQuantity; i++ {
		examples = append(examples, BuildFakeUser())
	}

	return &filtering.QueryFilteredResult[types.User]{
		Pagination: filtering.Pagination{
			Page:          1,
			Limit:         50,
			FilteredCount: exampleQuantity / 2,
			TotalCount:    exampleQuantity,
		},
		Data: examples,
	}
}

// BuildFakeUserCreationInput builds a faked UserRegistrationInput.
func BuildFakeUserCreationInput() *types.UserRegistrationInput {
	exampleUser := BuildFakeUser()

	return &types.UserRegistrationInput{
		Username:     exampleUser.Username,
		EmailAddress: fake.Email(),
		FirstName:    exampleUser.FirstName,
		LastName:     exampleUser.LastName,
		Password:     buildFakePassword(),
		Birthday:     exampleUser.Birthday,
	}
}

// BuildFakeUserRegistrationInput builds a faked UserRegistrationInput.
func BuildFakeUserRegistrationInput() *types.UserRegistrationInput {
	user := BuildFakeUser()
	return &types.UserRegistrationInput{
		Username:     user.Username,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		EmailAddress: user.EmailAddress,
		Password:     buildFakePassword(),
		Birthday:     user.Birthday,
	}
}

// BuildFakeUserRegistrationInputFromUser builds a faked UserRegistrationInput.
func BuildFakeUserRegistrationInputFromUser(user *types.User) *types.UserRegistrationInput {
	return &types.UserRegistrationInput{
		Username:     user.Username,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		EmailAddress: user.EmailAddress,
		Password:     buildFakePassword(),
		Birthday:     user.Birthday,
	}
}

// BuildFakeUserRegistrationInputWithInviteFromUser builds a faked UserRegistrationInput.
func BuildFakeUserRegistrationInputWithInviteFromUser(user *types.User) *types.UserRegistrationInput {
	return &types.UserRegistrationInput{
		Username:        user.Username,
		FirstName:       user.FirstName,
		LastName:        user.LastName,
		EmailAddress:    user.EmailAddress,
		Password:        buildFakePassword(),
		Birthday:        user.Birthday,
		InvitationToken: fake.UUID(),
		InvitationID:    BuildFakeID(),
	}
}

// BuildFakeUserCreationResponse builds a faked UserAccountStatusUpdateInput.
func BuildFakeUserCreationResponse() *types.UserCreationResponse {
	user := BuildFakeUser()
	return &types.UserCreationResponse{
		CreatedAt:       user.CreatedAt,
		Birthday:        user.Birthday,
		Username:        user.Username,
		EmailAddress:    user.EmailAddress,
		TwoFactorQRCode: fake.UUID(),
		CreatedUserID:   user.ID,
		AccountStatus:   user.AccountStatus,
		TwoFactorSecret: user.TwoFactorSecret,
		FirstName:       user.FirstName,
		LastName:        user.LastName,
	}
}

// BuildFakeAvatarUpdateInput builds a faked AvatarUpdateInput.
func BuildFakeAvatarUpdateInput() *types.AvatarUpdateInput {
	return &types.AvatarUpdateInput{
		Base64EncodedData: buildUniqueString(),
	}
}

func BuildFakeUserDetailsUpdateRequestInput() *types.UserDetailsUpdateRequestInput {
	return &types.UserDetailsUpdateRequestInput{
		FirstName:       buildUniqueString(),
		LastName:        buildUniqueString(),
		Birthday:        BuildFakeTime(),
		CurrentPassword: fake.Password(true, true, true, false, false, 32),
		TOTPToken:       "123456",
	}
}

func BuildFakeUserDetailsDatabaseUpdateInput() *types.UserDetailsDatabaseUpdateInput {
	return &types.UserDetailsDatabaseUpdateInput{
		FirstName: buildUniqueString(),
		LastName:  buildUniqueString(),
		Birthday:  BuildFakeTime(),
	}
}

// BuildFakeUserPermissionModificationInput builds a faked ModifyUserPermissionsInput.
func BuildFakeUserPermissionModificationInput() *types.ModifyUserPermissionsInput {
	return &types.ModifyUserPermissionsInput{
		Reason:  fake.Sentence(10),
		NewRole: authorization.AccountMemberRole.String(),
	}
}
