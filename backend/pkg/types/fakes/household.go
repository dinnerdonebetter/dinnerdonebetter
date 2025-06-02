package fakes

import (
	"github.com/dinnerdonebetter/backend/internal/lib/database/filtering"
	"github.com/dinnerdonebetter/backend/internal/lib/pointer"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"

	fake "github.com/brianvoe/gofakeit/v7"
)

// BuildFakeHousehold builds a faked household.
func BuildFakeHousehold() *types.Household {
	householdID := BuildFakeID()

	var memberships []*types.HouseholdUserMembershipWithUser
	for i := 0; i < exampleQuantity; i++ {
		membership := BuildFakeHouseholdUserMembershipWithUser()
		membership.BelongsToHousehold = householdID
		memberships = append(memberships, membership)
	}

	fakeAddress := fake.Address()
	key := fake.BitcoinPrivateKey()

	return &types.Household{
		ID:                         householdID,
		Name:                       fake.UUID(),
		BillingStatus:              types.UnpaidHouseholdBillingStatus,
		ContactPhone:               fake.PhoneFormatted(),
		PaymentProcessorCustomerID: fake.UUID(),
		CreatedAt:                  BuildFakeTime(),
		BelongsToUser:              fake.UUID(),
		Members:                    memberships,
		AddressLine1:               fakeAddress.Address,
		AddressLine2:               "",
		City:                       fakeAddress.City,
		State:                      fakeAddress.State,
		ZipCode:                    fakeAddress.Zip,
		Country:                    fakeAddress.Country,
		Latitude:                   pointer.To(buildFakeNumber()),
		Longitude:                  pointer.To(buildFakeNumber()),
		WebhookEncryptionKey:       key,
	}
}

// BuildFakeHouseholdsList builds a faked HouseholdList.
func BuildFakeHouseholdsList() *filtering.QueryFilteredResult[types.Household] {
	var examples []*types.Household
	for i := 0; i < exampleQuantity; i++ {
		examples = append(examples, BuildFakeHousehold())
	}

	return &filtering.QueryFilteredResult[types.Household]{
		Pagination: filtering.Pagination{
			Page:          1,
			Limit:         50,
			FilteredCount: exampleQuantity / 2,
			TotalCount:    exampleQuantity,
		},
		Data: examples,
	}
}

func BuildFakeHouseholdOwnershipTransferInput() *types.HouseholdOwnershipTransferInput {
	return &types.HouseholdOwnershipTransferInput{
		Reason:       fake.Sentence(5),
		CurrentOwner: BuildFakeID(),
		NewOwner:     BuildFakeID(),
	}
}

// BuildFakeHouseholdUpdateRequestInput builds a faked HouseholdUpdateRequestInput from a household.
func BuildFakeHouseholdUpdateRequestInput() *types.HouseholdUpdateRequestInput {
	household := BuildFakeHousehold()
	return converters.ConvertHouseholdToHouseholdUpdateRequestInput(household)
}

// BuildFakeHouseholdCreationRequestInput builds a faked HouseholdCreationRequestInput.
func BuildFakeHouseholdCreationRequestInput() *types.HouseholdCreationRequestInput {
	household := BuildFakeHousehold()
	return converters.ConvertHouseholdToHouseholdCreationRequestInput(household)
}
