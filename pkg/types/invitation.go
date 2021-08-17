package types

import (
	"context"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// Invitation represents an invitation.
	Invitation struct {
		ArchivedOn       *uint64 `json:"archivedOn"`
		LastUpdatedOn    *uint64 `json:"lastUpdatedOn"`
		ExternalID       string  `json:"externalID"`
		Code             string  `json:"code"`
		CreatedOn        uint64  `json:"createdOn"`
		ID               uint64  `json:"id"`
		BelongsToAccount uint64  `json:"belongsToAccount"`
		Consumed         bool    `json:"consumed"`
	}

	// InvitationList represents a list of invitations.
	InvitationList struct {
		Invitations []*Invitation `json:"invitations"`
		Pagination
	}

	// InvitationCreationInput represents what a user could set as input for creating invitations.
	InvitationCreationInput struct {
		Code             string `json:"code"`
		Consumed         bool   `json:"consumed"`
		BelongsToAccount uint64 `json:"-"`
	}

	// InvitationUpdateInput represents what a user could set as input for updating invitations.
	InvitationUpdateInput struct {
		Code             string `json:"code"`
		Consumed         bool   `json:"consumed"`
		BelongsToAccount uint64 `json:"-"`
	}

	// InvitationDataManager describes a structure capable of storing invitations permanently.
	InvitationDataManager interface {
		InvitationExists(ctx context.Context, invitationID uint64) (bool, error)
		GetInvitation(ctx context.Context, invitationID uint64) (*Invitation, error)
		GetAllInvitationsCount(ctx context.Context) (uint64, error)
		GetAllInvitations(ctx context.Context, resultChannel chan []*Invitation, bucketSize uint16) error
		GetInvitations(ctx context.Context, filter *QueryFilter) (*InvitationList, error)
		GetInvitationsWithIDs(ctx context.Context, accountID uint64, limit uint8, ids []uint64) ([]*Invitation, error)
		CreateInvitation(ctx context.Context, input *InvitationCreationInput, createdByUser uint64) (*Invitation, error)
		UpdateInvitation(ctx context.Context, updated *Invitation, changedByUser uint64, changes []*FieldChangeSummary) error
		ArchiveInvitation(ctx context.Context, invitationID, accountID, archivedBy uint64) error
		GetAuditLogEntriesForInvitation(ctx context.Context, invitationID uint64) ([]*AuditLogEntry, error)
	}

	// InvitationDataService describes a structure capable of serving traffic related to invitations.
	InvitationDataService interface {
		AuditEntryHandler(res http.ResponseWriter, req *http.Request)
		ListHandler(res http.ResponseWriter, req *http.Request)
		CreateHandler(res http.ResponseWriter, req *http.Request)
		ExistenceHandler(res http.ResponseWriter, req *http.Request)
		ReadHandler(res http.ResponseWriter, req *http.Request)
		UpdateHandler(res http.ResponseWriter, req *http.Request)
		ArchiveHandler(res http.ResponseWriter, req *http.Request)
	}
)

// Update merges an InvitationUpdateInput with an invitation.
func (x *Invitation) Update(input *InvitationUpdateInput) []*FieldChangeSummary {
	var out []*FieldChangeSummary

	if input.Code != x.Code {
		out = append(out, &FieldChangeSummary{
			FieldName: "Code",
			OldValue:  x.Code,
			NewValue:  input.Code,
		})

		x.Code = input.Code
	}

	if input.Consumed != x.Consumed {
		out = append(out, &FieldChangeSummary{
			FieldName: "Consumed",
			OldValue:  x.Consumed,
			NewValue:  input.Consumed,
		})

		x.Consumed = input.Consumed
	}

	return out
}

var _ validation.ValidatableWithContext = (*InvitationCreationInput)(nil)

// ValidateWithContext validates a InvitationCreationInput.
func (x *InvitationCreationInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.Code, validation.Required),
	)
}

var _ validation.ValidatableWithContext = (*InvitationUpdateInput)(nil)

// ValidateWithContext validates a InvitationUpdateInput.
func (x *InvitationUpdateInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.Code, validation.Required),
	)
}
