package types

import (
	"context"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const (
	// PendingHouseholdBillingStatus indicates a household invitation is pending.
	PendingHouseholdBillingStatus HouseholdInvitationStatus = "pending"
	// AcceptedHouseholdInvitationStatus indicates a household invitation was accepted.
	AcceptedHouseholdInvitationStatus HouseholdInvitationStatus = "accepted"
	// RejectedHouseholdInvitationStatus indicates a household invitation was rejected.
	RejectedHouseholdInvitationStatus HouseholdInvitationStatus = "rejected"
)

type (
	// HouseholdInvitationStatus is the type to use/compare against when checking billing status.
	HouseholdInvitationStatus string

	// HouseholdInvitation represents a household.
	HouseholdInvitation struct {
		_                    struct{}
		LastUpdatedOn        *uint64                   `json:"lastUpdatedOn"`
		ArchivedOn           *uint64                   `json:"archivedOn"`
		FromUser             string                    `json:"fromUser"`
		ToUser               string                    `json:"toUser"`
		DestinationHousehold string                    `json:"destinationHousehold"`
		ID                   string                    `json:"id"`
		Status               HouseholdInvitationStatus `json:"status"`
		CreatedOn            uint64                    `json:"createdOn"`
	}

	// HouseholdInvitationList represents a list of households.
	HouseholdInvitationList struct {
		_ struct{}

		HouseholdInvitations []*HouseholdInvitation `json:"households"`
		Pagination
	}

	// HouseholdInvitationCreationInput represents what a User could set as input for creating households.
	HouseholdInvitationCreationInput struct {
		_ struct{}

		ID                   string                    `json:"-"`
		Status               HouseholdInvitationStatus `json:"status"`
		FromUser             string                    `json:"fromUser"`
		ToUser               string                    `json:"toUser"`
		DestinationHousehold string                    `json:"destinationHousehold"`
		BelongsToUser        string                    `json:"-"`
	}

	// HouseholdInvitationDataManager describes a structure capable of storing households permanently.
	HouseholdInvitationDataManager interface {
		GetHouseholdInvitation(ctx context.Context, householdID, userID string) (*HouseholdInvitation, error)
		GetAllHouseholdInvitationsCount(ctx context.Context) (uint64, error)
		GetHouseholdInvitations(ctx context.Context, userID string, filter *QueryFilter) (*HouseholdInvitationList, error)
		CreateHouseholdInvitation(ctx context.Context, input *HouseholdInvitationCreationInput) (*HouseholdInvitation, error)
		ArchiveHouseholdInvitation(ctx context.Context, householdID string, userID string) error
	}

	// HouseholdInvitationDataService describes a structure capable of serving traffic related to households.
	HouseholdInvitationDataService interface {
		ListHandler(res http.ResponseWriter, req *http.Request)
		CreateHandler(res http.ResponseWriter, req *http.Request)
		ReadHandler(res http.ResponseWriter, req *http.Request)
		ArchiveHandler(res http.ResponseWriter, req *http.Request)
	}
)

var _ validation.ValidatableWithContext = (*HouseholdInvitationCreationInput)(nil)

// ValidateWithContext validates a HouseholdCreationInput.
func (x *HouseholdInvitationCreationInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(ctx, x,
		validation.Field(&x.ToUser, validation.Required),
		validation.Field(&x.FromUser, validation.Required),
		validation.Field(&x.DestinationHousehold, validation.Required),
	)
}
