package fakes

import (
	"gitlab.com/prixfixe/prixfixe/pkg/types"

	fake "github.com/brianvoe/gofakeit/v5"
)

// BuildFakeHousehold builds a faked household.
func BuildFakeHousehold() *types.Household {
	return &types.Household{
		ID:                         uint64(fake.Uint32()),
		ExternalID:                 fake.UUID(),
		Name:                       fake.Password(true, true, true, false, false, 32),
		BillingStatus:              types.PaidHouseholdBillingStatus,
		ContactEmail:               fake.Email(),
		ContactPhone:               fake.PhoneFormatted(),
		PaymentProcessorCustomerID: fake.UUID(),
		CreatedOn:                  uint64(uint32(fake.Date().Unix())),
		BelongsToUser:              fake.Uint64(),
		Members:                    BuildFakeHouseholdUserMembershipList().HouseholdUserMemberships,
	}
}

// BuildFakeHouseholdForUser builds a faked household.
func BuildFakeHouseholdForUser(u *types.User) *types.Household {
	return &types.Household{
		ID:            uint64(fake.Uint32()),
		ExternalID:    fake.UUID(),
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

// BuildFakeHouseholdUpdateInput builds a faked HouseholdUpdateInput from an household.
func BuildFakeHouseholdUpdateInput() *types.HouseholdUpdateInput {
	household := BuildFakeHousehold()
	return &types.HouseholdUpdateInput{
		Name:          household.Name,
		BelongsToUser: household.BelongsToUser,
	}
}

// BuildFakeHouseholdUpdateInputFromHousehold builds a faked HouseholdUpdateInput from an household.
func BuildFakeHouseholdUpdateInputFromHousehold(household *types.Household) *types.HouseholdUpdateInput {
	return &types.HouseholdUpdateInput{
		Name:          household.Name,
		BelongsToUser: household.BelongsToUser,
	}
}

// BuildFakeHouseholdCreationInput builds a faked HouseholdCreationInput.
func BuildFakeHouseholdCreationInput() *types.HouseholdCreationInput {
	household := BuildFakeHousehold()
	return BuildFakeHouseholdCreationInputFromHousehold(household)
}

// BuildFakeHouseholdCreationInputFromHousehold builds a faked HouseholdCreationInput from an household.
func BuildFakeHouseholdCreationInputFromHousehold(household *types.Household) *types.HouseholdCreationInput {
	return &types.HouseholdCreationInput{
		Name:          household.Name,
		ContactEmail:  household.ContactEmail,
		ContactPhone:  household.ContactPhone,
		BelongsToUser: household.BelongsToUser,
	}
}
