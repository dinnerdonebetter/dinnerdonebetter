package fakes

import (
	fake "github.com/brianvoe/gofakeit/v5"
	"github.com/segmentio/ksuid"

	"github.com/prixfixeco/api_server/pkg/types"
)

// BuildFakeHousehold builds a faked household.
func BuildFakeHousehold() *types.Household {
	return &types.Household{
		ID:                         ksuid.New().String(),
		Name:                       fake.UUID(),
		BillingStatus:              types.PaidHouseholdBillingStatus,
		ContactEmail:               fake.Email(),
		ContactPhone:               fake.PhoneFormatted(),
		PaymentProcessorCustomerID: fake.UUID(),
		CreatedOn:                  uint64(uint32(fake.Date().Unix())),
		BelongsToUser:              fake.UUID(),
		Members:                    BuildFakeHouseholdUserMembershipList().HouseholdUserMemberships,
	}
}

// BuildFakeHouseholdForUser builds a faked household.
func BuildFakeHouseholdForUser(u *types.User) *types.Household {
	return &types.Household{
		ID:            ksuid.New().String(),
		Name:          u.Username,
		CreatedOn:     uint64(uint32(fake.Date().Unix())),
		BelongsToUser: u.ID,
	}
}

// BuildFakeHouseholdList builds a faked HouseholdList.
func BuildFakeHouseholdList() *types.HouseholdList {
	var examples []*types.Household
	for i := 0; i < exampleQuantity; i++ {
		examples = append(examples, BuildFakeHousehold())
	}

	return &types.HouseholdList{
		Pagination: types.Pagination{
			Page:          1,
			Limit:         20,
			FilteredCount: exampleQuantity / 2,
			TotalCount:    exampleQuantity,
		},
		Households: examples,
	}
}

// BuildFakeHouseholdUpdateInput builds a faked HouseholdUpdateInput from a household.
func BuildFakeHouseholdUpdateInput() *types.HouseholdUpdateInput {
	household := BuildFakeHousehold()
	return &types.HouseholdUpdateInput{
		Name:          household.Name,
		BelongsToUser: household.BelongsToUser,
	}
}

// BuildFakeHouseholdUpdateInputFromHousehold builds a faked HouseholdUpdateInput from a household.
func BuildFakeHouseholdUpdateInputFromHousehold(household *types.Household) *types.HouseholdUpdateInput {
	return &types.HouseholdUpdateInput{
		Name:          household.Name,
		BelongsToUser: household.BelongsToUser,
	}
}

// BuildFakeHouseholdCreationInput builds a faked HouseholdCreationRequestInput.
func BuildFakeHouseholdCreationInput() *types.HouseholdCreationRequestInput {
	household := BuildFakeHousehold()
	return BuildFakeHouseholdCreationRequestInputFromHousehold(household)
}

// BuildFakeHouseholdCreationRequestInputFromHousehold builds a faked HouseholdCreationRequestInput from a household.
func BuildFakeHouseholdCreationRequestInputFromHousehold(household *types.Household) *types.HouseholdCreationRequestInput {
	return &types.HouseholdCreationRequestInput{
		ID:            ksuid.New().String(),
		Name:          household.Name,
		ContactEmail:  household.ContactEmail,
		ContactPhone:  household.ContactPhone,
		BelongsToUser: household.BelongsToUser,
	}
}

// BuildFakeHouseholdDatabaseCreationInputFromHousehold builds a faked HouseholdCreationRequestInput.
func BuildFakeHouseholdDatabaseCreationInputFromHousehold(household *types.Household) *types.HouseholdDatabaseCreationInput {
	return &types.HouseholdDatabaseCreationInput{
		ID:            household.ID,
		Name:          household.Name,
		ContactEmail:  household.ContactEmail,
		ContactPhone:  household.ContactPhone,
		BelongsToUser: household.BelongsToUser,
	}
}
