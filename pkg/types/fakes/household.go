package fakes

import (
	"github.com/prixfixeco/backend/pkg/types"
	"github.com/prixfixeco/backend/pkg/types/converters"

	fake "github.com/brianvoe/gofakeit/v5"
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

	return &types.Household{
		ID:                         householdID,
		Name:                       fake.UUID(),
		BillingStatus:              string(types.PaidHouseholdBillingStatus),
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
		Latitude:                   &fakeAddress.Latitude,
		Longitude:                  &fakeAddress.Longitude,
	}
}

// BuildFakeHouseholdForUser builds a faked household.
func BuildFakeHouseholdForUser(u *types.User) *types.Household {
	h := BuildFakeHousehold()
	h.BelongsToUser = u.ID
	return h
}

// BuildFakeHouseholdList builds a faked HouseholdList.
func BuildFakeHouseholdList() *types.QueryFilteredResult[types.Household] {
	var examples []*types.Household
	for i := 0; i < exampleQuantity; i++ {
		examples = append(examples, BuildFakeHousehold())
	}

	return &types.QueryFilteredResult[types.Household]{
		Pagination: types.Pagination{
			Page:          1,
			Limit:         20,
			FilteredCount: exampleQuantity / 2,
			TotalCount:    exampleQuantity,
		},
		Data: examples,
	}
}

// BuildFakeHouseholdUpdateInput builds a faked HouseholdUpdateRequestInput from a household.
func BuildFakeHouseholdUpdateInput() *types.HouseholdUpdateRequestInput {
	household := BuildFakeHousehold()
	return &types.HouseholdUpdateRequestInput{
		Name:          &household.Name,
		ContactPhone:  &household.ContactPhone,
		AddressLine1:  &household.AddressLine1,
		AddressLine2:  &household.AddressLine2,
		City:          &household.City,
		State:         &household.State,
		ZipCode:       &household.ZipCode,
		Country:       &household.Country,
		Latitude:      household.Latitude,
		Longitude:     household.Longitude,
		BelongsToUser: household.BelongsToUser,
	}
}

// BuildFakeHouseholdCreationInput builds a faked HouseholdCreationRequestInput.
func BuildFakeHouseholdCreationInput() *types.HouseholdCreationRequestInput {
	household := BuildFakeHousehold()
	return converters.ConvertHouseholdToHouseholdCreationRequestInput(household)
}
