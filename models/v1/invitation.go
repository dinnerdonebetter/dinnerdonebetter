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
		CreatedOn     uint64  `json:"createdOn"`
		LastUpdatedOn *uint64 `json:"lastUpdatedOn"`
		ArchivedOn    *uint64 `json:"archivedOn"`
		BelongsToUser uint64  `json:"belongsToUser"`
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
		GetAllInvitations(ctx context.Context, resultChannel chan []Invitation) error
		GetInvitations(ctx context.Context, filter *QueryFilter) (*InvitationList, error)
		GetInvitationsWithIDs(ctx context.Context, limit uint8, ids []uint64) ([]Invitation, error)
		CreateInvitation(ctx context.Context, input *InvitationCreationInput) (*Invitation, error)
		UpdateInvitation(ctx context.Context, updated *Invitation) error
		ArchiveInvitation(ctx context.Context, invitationID, userID uint64) error
	}

	// InvitationDataServer describes a structure capable of serving traffic related to invitations.
	InvitationDataServer interface {
		CreationInputMiddleware(next http.Handler) http.Handler
		UpdateInputMiddleware(next http.Handler) http.Handler

		ListHandler(res http.ResponseWriter, req *http.Request)
		CreateHandler(res http.ResponseWriter, req *http.Request)
		ExistenceHandler(res http.ResponseWriter, req *http.Request)
		ReadHandler(res http.ResponseWriter, req *http.Request)
		UpdateHandler(res http.ResponseWriter, req *http.Request)
		ArchiveHandler(res http.ResponseWriter, req *http.Request)
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
