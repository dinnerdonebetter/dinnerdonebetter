package models

import (
	"context"
	"net/http"
)

type (
	// Invitation represents an invitation.
	Invitation struct {
		ID            uint64  `json:"id"`
		Code          string  `json:"code"`
		Consumed      bool    `json:"consumed"`
		CreatedOn     uint64  `json:"created_on"`
		UpdatedOn     *uint64 `json:"updated_on"`
		ArchivedOn    *uint64 `json:"archived_on"`
		BelongsToUser uint64  `json:"belongs_to_user"`
	}

	// InvitationList represents a list of invitations.
	InvitationList struct {
		Pagination
		Invitations []Invitation `json:"invitations"`
	}

	// InvitationCreationInput represents what a user could set as input for creating invitations.
	InvitationCreationInput struct {
		Code          string `json:"code"`
		Consumed      bool   `json:"consumed"`
		BelongsToUser uint64 `json:"-"`
	}

	// InvitationUpdateInput represents what a user could set as input for updating invitations.
	InvitationUpdateInput struct {
		Code          string `json:"code"`
		Consumed      bool   `json:"consumed"`
		BelongsToUser uint64 `json:"-"`
	}

	// InvitationDataManager describes a structure capable of storing invitations permanently.
	InvitationDataManager interface {
		InvitationExists(ctx context.Context, invitationID uint64) (bool, error)
		GetInvitation(ctx context.Context, invitationID uint64) (*Invitation, error)
		GetAllInvitationsCount(ctx context.Context) (uint64, error)
		GetInvitations(ctx context.Context, filter *QueryFilter) (*InvitationList, error)
		CreateInvitation(ctx context.Context, input *InvitationCreationInput) (*Invitation, error)
		UpdateInvitation(ctx context.Context, updated *Invitation) error
		ArchiveInvitation(ctx context.Context, invitationID, userID uint64) error
	}

	// InvitationDataServer describes a structure capable of serving traffic related to invitations.
	InvitationDataServer interface {
		CreationInputMiddleware(next http.Handler) http.Handler
		UpdateInputMiddleware(next http.Handler) http.Handler

		ListHandler() http.HandlerFunc
		CreateHandler() http.HandlerFunc
		ExistenceHandler() http.HandlerFunc
		ReadHandler() http.HandlerFunc
		UpdateHandler() http.HandlerFunc
		ArchiveHandler() http.HandlerFunc
	}
)

// Update merges an InvitationInput with an invitation.
func (x *Invitation) Update(input *InvitationUpdateInput) {
	if input.Code != "" && input.Code != x.Code {
		x.Code = input.Code
	}

	if input.Consumed != x.Consumed {
		x.Consumed = input.Consumed
	}
}

// ToUpdateInput creates a InvitationUpdateInput struct for an invitation.
func (x *Invitation) ToUpdateInput() *InvitationUpdateInput {
	return &InvitationUpdateInput{
		Code:     x.Code,
		Consumed: x.Consumed,
	}
}
