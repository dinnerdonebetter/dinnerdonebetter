package types

import (
	"context"
	"fmt"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const (
	// PaidHouseholdBillingStatus indicates a household is fully paid.
	PaidHouseholdBillingStatus HouseholdBillingStatus = "paid"
	// UnpaidHouseholdBillingStatus indicates a household is not paid.
	UnpaidHouseholdBillingStatus HouseholdBillingStatus = "unpaid"
)

type (
	// HouseholdBillingStatus is the type to use/compare against when checking billing status.
	HouseholdBillingStatus string

	// Household represents a household.
	Household struct {
		_ struct{}

		ArchivedOn                 *uint64                    `json:"archivedOn"`
		SubscriptionPlanID         *uint64                    `json:"subscriptionPlanID"`
		LastUpdatedOn              *uint64                    `json:"lastUpdatedOn"`
		Name                       string                     `json:"name"`
		BillingStatus              HouseholdBillingStatus     `json:"billingStatus"`
		ContactEmail               string                     `json:"contactEmail"`
		ContactPhone               string                     `json:"contactPhone"`
		PaymentProcessorCustomerID string                     `json:"paymentProcessorCustomer"`
		BelongsToUser              string                     `json:"belongsToUser"`
		ID                         string                     `json:"id"`
		Members                    []*HouseholdUserMembership `json:"members"`
		CreatedOn                  uint64                     `json:"createdOn"`
	}

	// HouseholdList represents a list of households.
	HouseholdList struct {
		_ struct{}

		Households []*Household `json:"households"`
		Pagination
	}

	// HouseholdCreationRequestInput represents what a User could set as input for creating households.
	HouseholdCreationRequestInput struct {
		_ struct{}

		ID            string `json:"-"`
		Name          string `json:"name"`
		ContactEmail  string `json:"contactEmail"`
		ContactPhone  string `json:"contactPhone"`
		BelongsToUser string `json:"-"`
	}

	// HouseholdDatabaseCreationInput represents what a User could set as input for creating households.
	HouseholdDatabaseCreationInput struct {
		_ struct{}

		ID            string `json:"-"`
		Name          string `json:"name"`
		ContactEmail  string `json:"contactEmail"`
		ContactPhone  string `json:"contactPhone"`
		BelongsToUser string `json:"-"`
	}

	// HouseholdUpdateInput represents what a User could set as input for updating households.
	HouseholdUpdateInput struct {
		_ struct{}

		Name          string `json:"name"`
		ContactEmail  string `json:"contactEmail"`
		ContactPhone  string `json:"contactPhone"`
		BelongsToUser string `json:"-"`
	}

	// HouseholdDataManager describes a structure capable of storing households permanently.
	HouseholdDataManager interface {
		GetHousehold(ctx context.Context, householdID, userID string) (*Household, error)
		GetAllHouseholdsCount(ctx context.Context) (uint64, error)
		GetHouseholds(ctx context.Context, userID string, filter *QueryFilter) (*HouseholdList, error)
		GetHouseholdsForAdmin(ctx context.Context, filter *QueryFilter) (*HouseholdList, error)
		CreateHousehold(ctx context.Context, input *HouseholdDatabaseCreationInput) (*Household, error)
		UpdateHousehold(ctx context.Context, updated *Household) error
		ArchiveHousehold(ctx context.Context, householdID string, userID string) error
	}

	// HouseholdDataService describes a structure capable of serving traffic related to households.
	HouseholdDataService interface {
		ListHandler(res http.ResponseWriter, req *http.Request)
		CreateHandler(res http.ResponseWriter, req *http.Request)
		ReadHandler(res http.ResponseWriter, req *http.Request)
		UpdateHandler(res http.ResponseWriter, req *http.Request)
		ArchiveHandler(res http.ResponseWriter, req *http.Request)
		InviteMemberHandler(res http.ResponseWriter, req *http.Request)
		RemoveMemberHandler(res http.ResponseWriter, req *http.Request)
		MarkAsDefaultHouseholdHandler(res http.ResponseWriter, req *http.Request)
		ModifyMemberPermissionsHandler(res http.ResponseWriter, req *http.Request)
		TransferHouseholdOwnershipHandler(res http.ResponseWriter, req *http.Request)
	}
)

// Update merges a householdUpdateInput with a household.
func (x *Household) Update(input *HouseholdUpdateInput) {
	if input.Name != "" && input.Name != x.Name {
		x.Name = input.Name
	}
}

var _ validation.ValidatableWithContext = (*HouseholdCreationRequestInput)(nil)

// ValidateWithContext validates a HouseholdCreationRequestInput.
func (x *HouseholdCreationRequestInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(ctx, x,
		validation.Field(&x.Name, validation.Required),
	)
}

var _ validation.ValidatableWithContext = (*HouseholdUpdateInput)(nil)

// ValidateWithContext validates a HouseholdUpdateInput.
func (x *HouseholdUpdateInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(ctx, x,
		validation.Field(&x.Name, validation.Required),
	)
}

// HouseholdCreationInputForNewUser creates a new HouseholdInputCreation struct for a given user.
func HouseholdCreationInputForNewUser(u *User) *HouseholdCreationRequestInput {
	return &HouseholdCreationRequestInput{
		Name:          fmt.Sprintf("%s_default", u.Username),
		BelongsToUser: u.ID,
	}
}

// HouseholdDatabaseCreationInputFromHouseholdCreationInput creates a DatabaseCreationInput from a CreationInput.
func HouseholdDatabaseCreationInputFromHouseholdCreationInput(input *HouseholdCreationRequestInput) *HouseholdDatabaseCreationInput {
	x := &HouseholdDatabaseCreationInput{
		ID:            input.ID,
		Name:          input.Name,
		ContactEmail:  input.ContactEmail,
		ContactPhone:  input.ContactPhone,
		BelongsToUser: input.BelongsToUser,
	}

	return x
}
