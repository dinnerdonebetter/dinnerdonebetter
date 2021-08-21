package types

import (
	"context"
	"fmt"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const (
	// PaidHouseholdBillingStatus indicates an household is fully paid.
	PaidHouseholdBillingStatus HouseholdBillingStatus = "paid"
	// UnpaidHouseholdBillingStatus indicates an household is not paid.
	UnpaidHouseholdBillingStatus HouseholdBillingStatus = "unpaid"
)

type (
	// HouseholdBillingStatus is the type to use/compare against when checking billing status.
	HouseholdBillingStatus string

	// Household represents an household.
	Household struct {
		ArchivedOn                 *uint64                    `json:"archivedOn"`
		SubscriptionPlanID         *uint64                    `json:"subscriptionPlanID"`
		LastUpdatedOn              *uint64                    `json:"lastUpdatedOn"`
		Name                       string                     `json:"name"`
		BillingStatus              HouseholdBillingStatus     `json:"billingStatus"`
		ContactEmail               string                     `json:"contactEmail"`
		ContactPhone               string                     `json:"contactPhone"`
		PaymentProcessorCustomerID string                     `json:"paymentProcessorCustomer"`
		ExternalID                 string                     `json:"externalID"`
		Members                    []*HouseholdUserMembership `json:"members"`
		CreatedOn                  uint64                     `json:"createdOn"`
		ID                         uint64                     `json:"id"`
		BelongsToUser              uint64                     `json:"belongsToUser"`
	}

	// HouseholdList represents a list of households.
	HouseholdList struct {
		Households []*Household `json:"households"`
		Pagination
	}

	// HouseholdCreationInput represents what a User could set as input for creating households.
	HouseholdCreationInput struct {
		Name          string `json:"name"`
		ContactEmail  string `json:"contactEmail"`
		ContactPhone  string `json:"contactPhone"`
		BelongsToUser uint64 `json:"-"`
	}

	// HouseholdUpdateInput represents what a User could set as input for updating households.
	HouseholdUpdateInput struct {
		Name          string `json:"name"`
		ContactEmail  string `json:"contactEmail"`
		ContactPhone  string `json:"contactPhone"`
		BelongsToUser uint64 `json:"-"`
	}

	// HouseholdDataManager describes a structure capable of storing households permanently.
	HouseholdDataManager interface {
		GetHousehold(ctx context.Context, householdID, userID uint64) (*Household, error)
		GetAllHouseholdsCount(ctx context.Context) (uint64, error)
		GetAllHouseholds(ctx context.Context, resultChannel chan []*Household, bucketSize uint16) error
		GetHouseholds(ctx context.Context, userID uint64, filter *QueryFilter) (*HouseholdList, error)
		GetHouseholdsForAdmin(ctx context.Context, filter *QueryFilter) (*HouseholdList, error)
		CreateHousehold(ctx context.Context, input *HouseholdCreationInput, createdByUser uint64) (*Household, error)
		UpdateHousehold(ctx context.Context, updated *Household, changedByUser uint64, changes []*FieldChangeSummary) error
		ArchiveHousehold(ctx context.Context, householdID, userID, archivedByUser uint64) error
		GetAuditLogEntriesForHousehold(ctx context.Context, householdID uint64) ([]*AuditLogEntry, error)
	}

	// HouseholdDataService describes a structure capable of serving traffic related to households.
	HouseholdDataService interface {
		ListHandler(res http.ResponseWriter, req *http.Request)
		CreateHandler(res http.ResponseWriter, req *http.Request)
		ReadHandler(res http.ResponseWriter, req *http.Request)
		UpdateHandler(res http.ResponseWriter, req *http.Request)
		ArchiveHandler(res http.ResponseWriter, req *http.Request)
		AddMemberHandler(res http.ResponseWriter, req *http.Request)
		RemoveMemberHandler(res http.ResponseWriter, req *http.Request)
		MarkAsDefaultHouseholdHandler(res http.ResponseWriter, req *http.Request)
		ModifyMemberPermissionsHandler(res http.ResponseWriter, req *http.Request)
		TransferHouseholdOwnershipHandler(res http.ResponseWriter, req *http.Request)
		AuditEntryHandler(res http.ResponseWriter, req *http.Request)
	}
)

// Update merges an HouseholdUpdateInput with an household.
func (x *Household) Update(input *HouseholdUpdateInput) []*FieldChangeSummary {
	var out []*FieldChangeSummary

	if input.Name != x.Name {
		out = append(out, &FieldChangeSummary{
			FieldName: "Name",
			OldValue:  x.Name,
			NewValue:  input.Name,
		})

		x.Name = input.Name
	}

	return out
}

var _ validation.ValidatableWithContext = (*HouseholdCreationInput)(nil)

// ValidateWithContext validates a HouseholdCreationInput.
func (x *HouseholdCreationInput) ValidateWithContext(ctx context.Context) error {
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
func HouseholdCreationInputForNewUser(u *User) *HouseholdCreationInput {
	return &HouseholdCreationInput{
		Name:          fmt.Sprintf("%s_default", u.Username),
		BelongsToUser: u.ID,
	}
}
