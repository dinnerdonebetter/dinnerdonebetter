package types

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const (
	// HouseholdDataType indicates an event is related to a household.
	HouseholdDataType dataType = "household"

	// HouseholdCreatedCustomerEventType indicates a household was created.
	HouseholdCreatedCustomerEventType CustomerEventType = "household_created"
	// HouseholdUpdatedCustomerEventType indicates a household was updated.
	HouseholdUpdatedCustomerEventType CustomerEventType = "household_updated"
	// HouseholdArchivedCustomerEventType indicates a household was archived.
	HouseholdArchivedCustomerEventType CustomerEventType = "household_archived"
	// HouseholdMemberRemovedCustomerEventType indicates a household member was removed.
	HouseholdMemberRemovedCustomerEventType CustomerEventType = "household_member_removed"
	// HouseholdMembershipPermissionsUpdatedCustomerEventType indicates a household member's permissions were modified.
	HouseholdMembershipPermissionsUpdatedCustomerEventType CustomerEventType = "household_membership_permissions_updated"
	// HouseholdOwnershipTransferredCustomerEventType indicates a household was transferred to another owner.
	HouseholdOwnershipTransferredCustomerEventType CustomerEventType = "household_ownership_transferred"

	// PaidHouseholdBillingStatus indicates a household is fully paid.
	PaidHouseholdBillingStatus HouseholdBillingStatus = "paid"
	// UnpaidHouseholdBillingStatus indicates a household is not paid.
	UnpaidHouseholdBillingStatus HouseholdBillingStatus = "unpaid"

	// DefaultHouseholdTimeZone is the default time zone we will assign to a household.
	DefaultHouseholdTimeZone = "US/Central"
)

type (
	// HouseholdBillingStatus is the type to use/compare against when checking billing status.
	HouseholdBillingStatus string

	// Household represents a household.
	Household struct {
		_                          struct{}
		ArchivedAt                 *uint64                            `json:"archivedAt"`
		SubscriptionPlanID         *uint64                            `json:"subscriptionPlanID"`
		LastUpdatedAt              *uint64                            `json:"lastUpdatedAt"`
		ContactPhone               string                             `json:"contactPhone"`
		BillingStatus              HouseholdBillingStatus             `json:"billingStatus"`
		ContactEmail               string                             `json:"contactEmail"`
		Name                       string                             `json:"name"`
		PaymentProcessorCustomerID string                             `json:"paymentProcessorCustomer"`
		BelongsToUser              string                             `json:"belongsToUser"`
		ID                         string                             `json:"id"`
		TimeZone                   string                             `json:"timeZone"`
		Members                    []*HouseholdUserMembershipWithUser `json:"members"`
		CreatedAt                  uint64                             `json:"createdAt"`
	}

	// HouseholdList represents a list of households.
	HouseholdList struct {
		_ struct{}

		Households []*Household `json:"data"`
		Pagination
	}

	// HouseholdCreationRequestInput represents what a User could set as input for creating households.
	HouseholdCreationRequestInput struct {
		_ struct{}

		ID            string `json:"-"`
		Name          string `json:"name"`
		ContactEmail  string `json:"contactEmail"`
		ContactPhone  string `json:"contactPhone"`
		TimeZone      string `json:"timeZone"`
		BelongsToUser string `json:"-"`
	}

	// HouseholdDatabaseCreationInput represents what a User could set as input for creating households.
	HouseholdDatabaseCreationInput struct {
		_ struct{}

		ID            string `json:"-"`
		Name          string `json:"name"`
		ContactEmail  string `json:"contactEmail"`
		ContactPhone  string `json:"contactPhone"`
		TimeZone      string `json:"timeZone"`
		BelongsToUser string `json:"-"`
	}

	// HouseholdUpdateRequestInput represents what a User could set as input for updating households.
	HouseholdUpdateRequestInput struct {
		_ struct{}

		Name          *string `json:"name"`
		ContactEmail  *string `json:"contactEmail"`
		ContactPhone  *string `json:"contactPhone"`
		TimeZone      *string `json:"timeZone"`
		BelongsToUser string  `json:"-"`
	}

	// HouseholdDataManager describes a structure capable of storing households permanently.
	HouseholdDataManager interface {
		GetHousehold(ctx context.Context, householdID, userID string) (*Household, error)
		GetHouseholdByID(ctx context.Context, householdID string) (*Household, error)
		GetAllHouseholdsCount(ctx context.Context) (uint64, error)
		GetHouseholds(ctx context.Context, userID string, filter *QueryFilter) (*HouseholdList, error)
		GetHouseholdsForAdmin(ctx context.Context, userID string, filter *QueryFilter) (*HouseholdList, error)
		CreateHousehold(ctx context.Context, input *HouseholdDatabaseCreationInput) (*Household, error)
		UpdateHousehold(ctx context.Context, updated *Household) error
		ArchiveHousehold(ctx context.Context, householdID string, userID string) error
	}

	// HouseholdDataService describes a structure capable of serving traffic related to households.
	HouseholdDataService interface {
		ListHandler(res http.ResponseWriter, req *http.Request)
		CreateHandler(res http.ResponseWriter, req *http.Request)
		CurrentInfoHandler(res http.ResponseWriter, req *http.Request)
		ReadHandler(res http.ResponseWriter, req *http.Request)
		UpdateHandler(res http.ResponseWriter, req *http.Request)
		ArchiveHandler(res http.ResponseWriter, req *http.Request)
		RemoveMemberHandler(res http.ResponseWriter, req *http.Request)
		MarkAsDefaultHouseholdHandler(res http.ResponseWriter, req *http.Request)
		ModifyMemberPermissionsHandler(res http.ResponseWriter, req *http.Request)
		TransferHouseholdOwnershipHandler(res http.ResponseWriter, req *http.Request)
	}
)

// Update merges a householdUpdateInput with a household.
func (x *Household) Update(input *HouseholdUpdateRequestInput) {
	if input.Name != nil && *input.Name != x.Name {
		x.Name = *input.Name
	}

	if input.ContactEmail != nil && *input.ContactEmail != x.ContactEmail {
		x.ContactEmail = *input.ContactEmail
	}

	if input.ContactPhone != nil && *input.ContactPhone != x.ContactPhone {
		x.ContactPhone = *input.ContactPhone
	}

	if input.TimeZone != nil && *input.TimeZone != x.TimeZone {
		if _, err := time.LoadLocation(*input.TimeZone); err != nil {
			// FIXME: we should return an error here, right?
			log.Println(err)
		}
		x.TimeZone = *input.TimeZone
	}
}

var _ validation.ValidatableWithContext = (*HouseholdCreationRequestInput)(nil)

// ValidateWithContext validates a HouseholdCreationRequestInput.
func (x *HouseholdCreationRequestInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(ctx, x,
		validation.Field(&x.Name, validation.Required),
	)
}

var _ validation.ValidatableWithContext = (*HouseholdUpdateRequestInput)(nil)

// ValidateWithContext validates a HouseholdUpdateRequestInput.
func (x *HouseholdUpdateRequestInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(ctx, x,
		validation.Field(&x.Name, validation.Required),
	)
}

// HouseholdCreationInputForNewUser creates a new HouseholdInputCreation struct for a given user.
func HouseholdCreationInputForNewUser(u *User) *HouseholdCreationRequestInput {
	return &HouseholdCreationRequestInput{
		Name:          fmt.Sprintf("%s_default", u.ID),
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
		TimeZone:      input.TimeZone,
	}

	return x
}

// HouseholdUpdateRequestInputFromHousehold creates a HouseholdCreationRequestInput from a CreationInput.
func HouseholdUpdateRequestInputFromHousehold(input *Household) *HouseholdUpdateRequestInput {
	x := &HouseholdUpdateRequestInput{
		Name:          &input.Name,
		ContactEmail:  &input.ContactEmail,
		ContactPhone:  &input.ContactPhone,
		BelongsToUser: input.BelongsToUser,
		TimeZone:      &input.TimeZone,
	}

	return x
}
