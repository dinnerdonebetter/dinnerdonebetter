package types

import (
	"context"
	"net/http"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const (
	// HouseholdInvitationDataType indicates an event is related to a household invitation.
	HouseholdInvitationDataType dataType = "household_invitation"

	// HouseholdInvitationCreatedCustomerEventType indicates a household invitation was created.
	HouseholdInvitationCreatedCustomerEventType CustomerEventType = "household_invitation_created"
	// HouseholdInvitationCanceledCustomerEventType indicates a household invitation was created.
	HouseholdInvitationCanceledCustomerEventType CustomerEventType = "household_invitation_canceled"
	// HouseholdInvitationAcceptedCustomerEventType indicates a household invitation was created.
	HouseholdInvitationAcceptedCustomerEventType CustomerEventType = "household_invitation_accepted"
	// HouseholdInvitationRejectedCustomerEventType indicates a household invitation was created.
	HouseholdInvitationRejectedCustomerEventType CustomerEventType = "household_invitation_rejected"

	// PendingHouseholdInvitationStatus indicates a household invitation is pending.
	PendingHouseholdInvitationStatus HouseholdInvitationStatus = "pending"
	// CancelledHouseholdInvitationStatus indicates a household invitation was accepted.
	CancelledHouseholdInvitationStatus HouseholdInvitationStatus = "cancelled"
	// AcceptedHouseholdInvitationStatus indicates a household invitation was accepted.
	AcceptedHouseholdInvitationStatus HouseholdInvitationStatus = "accepted"
	// RejectedHouseholdInvitationStatus indicates a household invitation was rejected.
	RejectedHouseholdInvitationStatus HouseholdInvitationStatus = "rejected"
)

type (
	// HouseholdInvitationStatus is the type to use/compare against when checking invitation status.
	HouseholdInvitationStatus string

	// HouseholdInvitationCreationRequestInput represents what a User could set as input for creating household invitations.
	HouseholdInvitationCreationRequestInput struct {
		_         struct{}
		ExpiresAt *time.Time `json:"expiresAt"`
		Note      string     `json:"note"`
		ToEmail   string     `json:"toEmail"`
	}

	// HouseholdInvitationDatabaseCreationInput represents what a User could set as input for creating household invitations.
	HouseholdInvitationDatabaseCreationInput struct {
		_ struct{}

		ID                     string
		FromUser               string
		ToUser                 *string
		Note                   string
		ToEmail                string
		Token                  string
		ExpiresAt              time.Time
		DestinationHouseholdID string
	}

	// HouseholdInvitation represents a household invitation.
	HouseholdInvitation struct {
		_                    struct{}
		CreatedAt            time.Time  `json:"createdAt"`
		LastUpdatedAt        *time.Time `json:"lastUpdatedAt"`
		ArchivedAt           *time.Time `json:"archivedAt"`
		ToUser               *string    `json:"toUser"`
		Status               string     `json:"status"`
		ToEmail              string     `json:"toEmail"`
		StatusNote           string     `json:"statusNote"`
		Token                string     `json:"token"`
		ID                   string     `json:"id"`
		Note                 string     `json:"note"`
		ExpiresAt            time.Time  `json:"expiresAt"`
		DestinationHousehold Household  `json:"destinationHousehold"`
		FromUser             User       `json:"fromUser"`
	}

	// HouseholdInvitationUpdateRequestInput is used by users to update the status of a given household invitation.
	HouseholdInvitationUpdateRequestInput struct {
		Token string `json:"token"`
		Note  string `json:"note"`
	}

	// HouseholdInvitationDataManager describes a structure capable of storing household invitations permanently.
	HouseholdInvitationDataManager interface {
		HouseholdInvitationExists(ctx context.Context, householdInvitationID string) (bool, error)
		GetHouseholdInvitationByHouseholdAndID(ctx context.Context, householdID, householdInvitationID string) (*HouseholdInvitation, error)
		GetHouseholdInvitationByTokenAndID(ctx context.Context, token, invitationID string) (*HouseholdInvitation, error)
		GetHouseholdInvitationByEmailAndToken(ctx context.Context, emailAddress, token string) (*HouseholdInvitation, error)
		GetPendingHouseholdInvitationsFromUser(ctx context.Context, userID string, filter *QueryFilter) (*QueryFilteredResult[HouseholdInvitation], error)
		GetPendingHouseholdInvitationsForUser(ctx context.Context, userID string, filter *QueryFilter) (*QueryFilteredResult[HouseholdInvitation], error)
		CreateHouseholdInvitation(ctx context.Context, input *HouseholdInvitationDatabaseCreationInput) (*HouseholdInvitation, error)
		CancelHouseholdInvitation(ctx context.Context, householdInvitationID, note string) error
		AcceptHouseholdInvitation(ctx context.Context, householdInvitationID, token, note string) error
		RejectHouseholdInvitation(ctx context.Context, householdInvitationID, note string) error
	}

	// HouseholdInvitationDataService describes a structure capable of serving traffic related to household invitations.
	HouseholdInvitationDataService interface {
		ReadHandler(res http.ResponseWriter, req *http.Request)
		InboundInvitesHandler(res http.ResponseWriter, req *http.Request)
		OutboundInvitesHandler(res http.ResponseWriter, req *http.Request)
		InviteMemberHandler(res http.ResponseWriter, req *http.Request)
		CancelInviteHandler(res http.ResponseWriter, req *http.Request)
		AcceptInviteHandler(res http.ResponseWriter, req *http.Request)
		RejectInviteHandler(res http.ResponseWriter, req *http.Request)
	}
)

var _ validation.ValidatableWithContext = (*HouseholdInvitationCreationRequestInput)(nil)

// ValidateWithContext validates a HouseholdCreationRequestInput.
func (x *HouseholdInvitationCreationRequestInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(ctx, x,
		validation.Field(&x.ToEmail, validation.Required),
	)
}

var _ validation.ValidatableWithContext = (*HouseholdInvitationUpdateRequestInput)(nil)

// ValidateWithContext validates a HouseholdInvitationUpdateRequestInput.
func (x *HouseholdInvitationUpdateRequestInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.Token, validation.Required),
	)
}
