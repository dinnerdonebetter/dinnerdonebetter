package main

import (
	"context"
	"math/rand"
	"net/http"

	client "gitlab.com/prixfixe/prixfixe/client/v1/http"
	models "gitlab.com/prixfixe/prixfixe/models/v1"
	randmodel "gitlab.com/prixfixe/prixfixe/tests/v1/testutil/rand/model"
)

// fetchRandomInvitation retrieves a random invitation from the list of available invitations
func fetchRandomInvitation(c *client.V1Client) *models.Invitation {
	invitationsRes, err := c.GetInvitations(context.Background(), nil)
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
				return c.BuildCreateInvitationRequest(context.Background(), randmodel.RandomInvitationCreationInput())
			},
			Weight: 100,
		},
		"GetInvitation": {
			Name: "GetInvitation",
			Action: func() (*http.Request, error) {
				if randomInvitation := fetchRandomInvitation(c); randomInvitation != nil {
					return c.BuildGetInvitationRequest(context.Background(), randomInvitation.ID)
				}
				return nil, ErrUnavailableYet
			},
			Weight: 100,
		},
		"GetInvitations": {
			Name: "GetInvitations",
			Action: func() (*http.Request, error) {
				return c.BuildGetInvitationsRequest(context.Background(), nil)
			},
			Weight: 100,
		},
		"UpdateInvitation": {
			Name: "UpdateInvitation",
			Action: func() (*http.Request, error) {
				if randomInvitation := fetchRandomInvitation(c); randomInvitation != nil {
					randomInvitation.Code = randmodel.RandomInvitationCreationInput().Code
					randomInvitation.Consumed = randmodel.RandomInvitationCreationInput().Consumed
					return c.BuildUpdateInvitationRequest(context.Background(), randomInvitation)
				}
				return nil, ErrUnavailableYet
			},
			Weight: 100,
		},
		"ArchiveInvitation": {
			Name: "ArchiveInvitation",
			Action: func() (*http.Request, error) {
				if randomInvitation := fetchRandomInvitation(c); randomInvitation != nil {
					return c.BuildArchiveInvitationRequest(context.Background(), randomInvitation.ID)
				}
				return nil, ErrUnavailableYet
			},
			Weight: 85,
		},
	}
}
