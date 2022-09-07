package fakes

import (
	fake "github.com/brianvoe/gofakeit/v5"
	"github.com/segmentio/ksuid"

	"github.com/prixfixeco/api_server/pkg/types"
)

// BuildFakeHousehold builds a faked household.
func BuildFakeHousehold() *types.Household {
	householdID := ksuid.New().String()

	var memberships []*types.HouseholdUserMembershipWithUser
	for i := 0; i < exampleQuantity; i++ {
		membership := BuildFakeHouseholdUserMembershipWithUser()
		membership.BelongsToHousehold = householdID
		memberships = append(memberships, membership)
	}

	return &types.Household{
		ID:                         householdID,
		Name:                       fake.UUID(),
		BillingStatus:              types.PaidHouseholdBillingStatus,
		ContactEmail:               fake.Email(),
		ContactPhone:               fake.PhoneFormatted(),
		PaymentProcessorCustomerID: fake.UUID(),
		CreatedAt:                  fake.Date(),
		BelongsToUser:              fake.UUID(),
		Members:                    memberships,
		TimeZone:                   types.DefaultHouseholdTimeZone,
	}
}

// BuildFakeHouseholdForUser builds a faked household.
func BuildFakeHouseholdForUser(u *types.User) *types.Household {
	h := BuildFakeHousehold()
	h.BelongsToUser = u.ID
	return h
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

// BuildFakeHouseholdUpdateInput builds a faked HouseholdUpdateRequestInput from a household.
func BuildFakeHouseholdUpdateInput() *types.HouseholdUpdateRequestInput {
	household := BuildFakeHousehold()
	return &types.HouseholdUpdateRequestInput{
		Name:          &household.Name,
		BelongsToUser: household.BelongsToUser,
		TimeZone:      &household.TimeZone,
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
		TimeZone:      household.TimeZone,
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
		TimeZone:      household.TimeZone,
	}
}
