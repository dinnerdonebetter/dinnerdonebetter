package types

import (
	"context"
	"net/http"
	"time"

	"github.com/dinnerdonebetter/backend/internal/lib/database/filtering"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const (
	// HouseholdInvitationCreatedServiceEventType indicates a household invitation was created.
	HouseholdInvitationCreatedServiceEventType = "household_invitation_created"
	// HouseholdInvitationCanceledServiceEventType indicates a household invitation was created.
	HouseholdInvitationCanceledServiceEventType = "household_invitation_canceled"
	// HouseholdInvitationAcceptedServiceEventType indicates a household invitation was created.
	HouseholdInvitationAcceptedServiceEventType = "household_invitation_accepted"
	// HouseholdInvitationRejectedServiceEventType indicates a household invitation was created.
	HouseholdInvitationRejectedServiceEventType = "household_invitation_rejected"

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
		_ struct{} `json:"-"`

		ExpiresAt *time.Time `json:"expiresAt"`
		Note      string     `json:"note"`
		ToEmail   string     `json:"toEmail"`
		ToName    string     `json:"toName"`
	}

	// HouseholdInvitationDatabaseCreationInput represents what a User could set as input for creating household invitations.
	HouseholdInvitationDatabaseCreationInput struct {
		_ struct{} `json:"-"`

		ID                     string    `json:"-"`
		FromUser               string    `json:"-"`
		ToUser                 *string   `json:"-"`
		Note                   string    `json:"-"`
		ToEmail                string    `json:"-"`
		Token                  string    `json:"-"`
		ToName                 string    `json:"-"`
		ExpiresAt              time.Time `json:"-"`
		DestinationHouseholdID string    `json:"-"`
	}

	// HouseholdInvitation represents a household invitation.
	HouseholdInvitation struct {
		_ struct{} `json:"-"`

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
		ToName               string     `json:"toName"`
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
		GetPendingHouseholdInvitationsFromUser(ctx context.Context, userID string, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[HouseholdInvitation], error)
		GetPendingHouseholdInvitationsForUser(ctx context.Context, userID string, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[HouseholdInvitation], error)
		CreateHouseholdInvitation(ctx context.Context, input *HouseholdInvitationDatabaseCreationInput) (*HouseholdInvitation, error)
		CancelHouseholdInvitation(ctx context.Context, householdInvitationID, note string) error
		AcceptHouseholdInvitation(ctx context.Context, householdInvitationID, token, note string) error
		RejectHouseholdInvitation(ctx context.Context, householdInvitationID, note string) error
	}

	// HouseholdInvitationDataService describes a structure capable of serving traffic related to household invitations.
	HouseholdInvitationDataService interface {
		ReadHouseholdInviteHandler(http.ResponseWriter, *http.Request)
		InboundInvitesHandler(http.ResponseWriter, *http.Request)
		OutboundInvitesHandler(http.ResponseWriter, *http.Request)
		InviteMemberHandler(http.ResponseWriter, *http.Request)
		CancelInviteHandler(http.ResponseWriter, *http.Request)
		AcceptInviteHandler(http.ResponseWriter, *http.Request)
		RejectInviteHandler(http.ResponseWriter, *http.Request)
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
