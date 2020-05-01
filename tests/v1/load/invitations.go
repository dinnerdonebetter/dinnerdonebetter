package main

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"

	client "gitlab.com/prixfixe/prixfixe/client/v1/http"
	models "gitlab.com/prixfixe/prixfixe/models/v1"
	fakemodels "gitlab.com/prixfixe/prixfixe/models/v1/fake"
)

// fetchRandomInvitation retrieves a random invitation from the list of available invitations.
func fetchRandomInvitation(ctx context.Context, c *client.V1Client) *models.Invitation {
	invitationsRes, err := c.GetInvitations(ctx, nil)
	if err != nil || invitationsRes == nil || len(invitationsRes.Invitations) == 0 {
		return nil
	}

	randIndex := rand.Intn(len(invitationsRes.Invitations))
	return &invitationsRes.Invitations[randIndex]
}

func buildInvitationActions(c *client.V1Client) map[string]*Action {
	return map[string]*Action{
		"CreateInvitation": {
			Name: "CreateInvitation",
			Action: func() (*http.Request, error) {
				ctx := context.Background()

				invitationInput := fakemodels.BuildFakeInvitationCreationInput()

				return c.BuildCreateInvitationRequest(ctx, invitationInput)
			},
			Weight: 100,
		},
		"GetInvitation": {
			Name: "GetInvitation",
			Action: func() (*http.Request, error) {
				ctx := context.Background()

				randomInvitation := fetchRandomInvitation(ctx, c)
				if randomInvitation == nil {
					return nil, fmt.Errorf("retrieving random invitation: %w", ErrUnavailableYet)
				}

				return c.BuildGetInvitationRequest(ctx, randomInvitation.ID)
			},
			Weight: 100,
		},
		"GetInvitations": {
			Name: "GetInvitations",
			Action: func() (*http.Request, error) {
				ctx := context.Background()

				return c.BuildGetInvitationsRequest(ctx, nil)
			},
			Weight: 100,
		},
		"UpdateInvitation": {
			Name: "UpdateInvitation",
			Action: func() (*http.Request, error) {
				ctx := context.Background()

				if randomInvitation := fetchRandomInvitation(ctx, c); randomInvitation != nil {
					newInvitation := fakemodels.BuildFakeInvitationCreationInput()
					randomInvitation.Code = newInvitation.Code
					randomInvitation.Consumed = newInvitation.Consumed
					return c.BuildUpdateInvitationRequest(ctx, randomInvitation)
				}

				return nil, ErrUnavailableYet
			},
			Weight: 100,
		},
		"ArchiveInvitation": {
			Name: "ArchiveInvitation",
			Action: func() (*http.Request, error) {
				ctx := context.Background()

				randomInvitation := fetchRandomInvitation(ctx, c)
				if randomInvitation == nil {
					return nil, fmt.Errorf("retrieving random invitation: %w", ErrUnavailableYet)
				}

				return c.BuildArchiveInvitationRequest(ctx, randomInvitation.ID)
			},
			Weight: 85,
		},
	}
}
