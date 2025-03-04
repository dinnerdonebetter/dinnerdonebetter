package fakes

import (
	"github.com/dinnerdonebetter/backend/internal/lib/database/filtering"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"

	fake "github.com/brianvoe/gofakeit/v7"
)

// BuildFakeHouseholdInvitation builds a faked HouseholdInvitation.
func BuildFakeHouseholdInvitation() *types.HouseholdInvitation {
	return &types.HouseholdInvitation{
		FromUser:             *BuildFakeUser(),
		ToEmail:              fake.Email(),
		ToName:               buildUniqueString(),
		ToUser:               func(s string) *string { return &s }(buildUniqueString()),
		Note:                 buildUniqueString(),
		StatusNote:           buildUniqueString(),
		Token:                fake.UUID(),
		DestinationHousehold: *BuildFakeHousehold(),
		ID:                   BuildFakeID(),
		ExpiresAt:            BuildFakeTime(),
		Status:               string(types.PendingHouseholdInvitationStatus),
		CreatedAt:            BuildFakeTime(),
	}
}

// BuildFakeHouseholdInvitationsList builds a faked HouseholdInvitationList.
func BuildFakeHouseholdInvitationsList() *filtering.QueryFilteredResult[types.HouseholdInvitation] {
	var examples []*types.HouseholdInvitation
	for range exampleQuantity {
		examples = append(examples, BuildFakeHouseholdInvitation())
	}

	return &filtering.QueryFilteredResult[types.HouseholdInvitation]{
		Pagination: filtering.Pagination{
			Page:          1,
			Limit:         50,
			FilteredCount: exampleQuantity / 2,
			TotalCount:    exampleQuantity,
		},
		Data: examples,
	}
}

func BuildFakeHouseholdInvitationUpdateRequestInput() *types.HouseholdInvitationUpdateRequestInput {
	return &types.HouseholdInvitationUpdateRequestInput{
		Token: BuildFakeID(),
		Note:  fake.Sentence(3),
	}
}

// BuildFakeHouseholdInvitationCreationRequestInput builds a faked HouseholdInvitationCreationRequestInput from a webhook.
func BuildFakeHouseholdInvitationCreationRequestInput() *types.HouseholdInvitationCreationRequestInput {
	invitation := BuildFakeHouseholdInvitation()
	return converters.ConvertHouseholdInvitationToHouseholdInvitationCreationInput(invitation)
}
