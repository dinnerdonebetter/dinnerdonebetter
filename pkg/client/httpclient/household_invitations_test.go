package httpclient

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/suite"

	"github.com/prixfixeco/api_server/pkg/types"
	"github.com/prixfixeco/api_server/pkg/types/fakes"
)

func TestHouseholdInvitations(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(householdInvitationsTestSuite))
}

type householdInvitationsTestSuite struct {
	suite.Suite

	ctx                        context.Context
	exampleHousehold           *types.Household
	exampleHouseholdInvitation *types.HouseholdInvitation
	exampleUser                *types.User
	exampleHouseholdList       *types.HouseholdList
}

var _ suite.SetupTestSuite = (*householdInvitationsTestSuite)(nil)

func (s *householdInvitationsTestSuite) SetupTest() {
	s.ctx = context.Background()
	s.exampleUser = fakes.BuildFakeUser()
	s.exampleHousehold = fakes.BuildFakeHousehold()
	s.exampleHousehold.BelongsToUser = s.exampleUser.ID
	s.exampleHouseholdInvitation = fakes.BuildFakeHouseholdInvitation()
	s.exampleHouseholdInvitation.FromUser = s.exampleUser.ID
	s.exampleHouseholdInvitation.ToUser = func(s string) *string { return &s }(fakes.BuildFakeUser().ID)
	s.exampleHouseholdList = fakes.BuildFakeHouseholdList()
}

func (s *householdInvitationsTestSuite) TestClient_GetHouseholdInvitation() {
	const expectedPath = "/api/v1/households/%s/invitations/%s"

	s.Run("standard", func() {
		t := s.T()

		s.exampleHousehold.BelongsToUser = ""

		spec := newRequestSpec(true, http.MethodGet, "", expectedPath, s.exampleHousehold.ID, s.exampleHouseholdInvitation.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleHouseholdInvitation)

		actual, err := c.GetHouseholdInvitation(s.ctx, s.exampleHousehold.ID, s.exampleHouseholdInvitation.ID)
		assert.NoError(t, err)
		assert.NotNil(t, actual)
	})

	s.Run("with invalid household ID", func() {
		t := s.T()

		s.exampleHousehold.BelongsToUser = ""

		spec := newRequestSpec(true, http.MethodGet, "", expectedPath, s.exampleHousehold.ID, s.exampleHouseholdInvitation.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleHouseholdInvitation)

		actual, err := c.GetHouseholdInvitation(s.ctx, "", s.exampleHouseholdInvitation.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	s.Run("with invalid invitation ID", func() {
		t := s.T()

		s.exampleHousehold.BelongsToUser = ""

		spec := newRequestSpec(true, http.MethodGet, "", expectedPath, s.exampleHousehold.ID, s.exampleHouseholdInvitation.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleHouseholdInvitation)

		actual, err := c.GetHouseholdInvitation(s.ctx, s.exampleHousehold.ID, "")
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	s.Run("with error building request", func() {
		t := s.T()

		s.exampleHousehold.BelongsToUser = ""

		c := buildTestClientWithInvalidURL(t)

		actual, err := c.GetHouseholdInvitation(s.ctx, s.exampleHousehold.ID, s.exampleHouseholdInvitation.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		s.exampleHousehold.BelongsToUser = ""

		c, _ := buildTestClientThatWaitsTooLong(t)

		actual, err := c.GetHouseholdInvitation(s.ctx, s.exampleHousehold.ID, s.exampleHouseholdInvitation.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func (s *householdInvitationsTestSuite) TestClient_GetPendingHouseholdInvitationsForUser() {
	const expectedPath = "/api/v1/household_invitations/received"

	s.Run("standard", func() {
		t := s.T()

		filter := types.DefaultQueryFilter()
		s.exampleHousehold.BelongsToUser = ""

		spec := newRequestSpec(true, http.MethodGet, "includeArchived=false&limit=20&page=1&sortBy=asc", expectedPath)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleHouseholdInvitation)

		actual, err := c.GetPendingHouseholdInvitationsForUser(s.ctx, filter)
		assert.NoError(t, err)
		assert.NotNil(t, actual)
	})

	s.Run("with error building request", func() {
		t := s.T()

		filter := types.DefaultQueryFilter()

		c := buildTestClientWithInvalidURL(t)

		actual, err := c.GetPendingHouseholdInvitationsForUser(s.ctx, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		filter := types.DefaultQueryFilter()

		c, _ := buildTestClientThatWaitsTooLong(t)

		actual, err := c.GetPendingHouseholdInvitationsForUser(s.ctx, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func (s *householdInvitationsTestSuite) TestClient_GetPendingHouseholdInvitationsFromUser() {
	const expectedPath = "/api/v1/household_invitations/sent"

	s.Run("standard", func() {
		t := s.T()

		filter := types.DefaultQueryFilter()
		s.exampleHousehold.BelongsToUser = ""

		spec := newRequestSpec(true, http.MethodGet, "includeArchived=false&limit=20&page=1&sortBy=asc", expectedPath)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleHouseholdInvitation)

		actual, err := c.GetPendingHouseholdInvitationsFromUser(s.ctx, filter)
		assert.NoError(t, err)
		assert.NotNil(t, actual)
	})

	s.Run("with error building request", func() {
		t := s.T()

		filter := types.DefaultQueryFilter()

		c := buildTestClientWithInvalidURL(t)

		actual, err := c.GetPendingHouseholdInvitationsFromUser(s.ctx, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		filter := types.DefaultQueryFilter()

		c, _ := buildTestClientThatWaitsTooLong(t)

		actual, err := c.GetPendingHouseholdInvitationsFromUser(s.ctx, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func (s *householdInvitationsTestSuite) TestClient_CancelHouseholdInvitation() {
	const expectedPath = "/api/v1/households/%s/invitations/%s/cancel"

	s.Run("standard", func() {
		t := s.T()

		s.exampleHousehold.BelongsToUser = ""

		spec := newRequestSpec(false, http.MethodPut, "", expectedPath, s.exampleHousehold.ID, s.exampleHouseholdInvitation.ID)
		c, _ := buildTestClientWithStatusCodeResponse(t, spec, http.StatusAccepted)

		assert.NoError(t, c.CancelHouseholdInvitation(s.ctx, s.exampleHousehold.ID, s.exampleHouseholdInvitation.ID, t.Name()))
	})

	s.Run("with invalid household ID", func() {
		t := s.T()

		s.exampleHousehold.BelongsToUser = ""

		spec := newRequestSpec(false, http.MethodPut, "", expectedPath, s.exampleHousehold.ID, s.exampleHouseholdInvitation.ID)
		c, _ := buildTestClientWithStatusCodeResponse(t, spec, http.StatusAccepted)

		assert.Error(t, c.CancelHouseholdInvitation(s.ctx, "", s.exampleHouseholdInvitation.ID, t.Name()))
	})

	s.Run("with invalid household invitation ID", func() {
		t := s.T()

		s.exampleHousehold.BelongsToUser = ""

		spec := newRequestSpec(false, http.MethodPut, "", expectedPath, s.exampleHousehold.ID, s.exampleHouseholdInvitation.ID)
		c, _ := buildTestClientWithStatusCodeResponse(t, spec, http.StatusAccepted)

		assert.Error(t, c.CancelHouseholdInvitation(s.ctx, s.exampleHousehold.ID, "", t.Name()))
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)

		assert.Error(t, c.CancelHouseholdInvitation(s.ctx, s.exampleHousehold.ID, s.exampleHouseholdInvitation.ID, t.Name()))
	})

	s.Run("with error executing request", func() {
		t := s.T()

		s.exampleHousehold.BelongsToUser = ""

		c, _ := buildTestClientThatWaitsTooLong(t)

		assert.Error(t, c.CancelHouseholdInvitation(s.ctx, s.exampleHousehold.ID, s.exampleHouseholdInvitation.ID, t.Name()))
	})
}

func (s *householdInvitationsTestSuite) TestClient_AcceptHouseholdInvitation() {
	const expectedPath = "/api/v1/households/%s/invitations/%s/accept"

	s.Run("standard", func() {
		t := s.T()

		s.exampleHousehold.BelongsToUser = ""

		spec := newRequestSpec(false, http.MethodPut, "", expectedPath, s.exampleHousehold.ID, s.exampleHouseholdInvitation.ID)
		c, _ := buildTestClientWithStatusCodeResponse(t, spec, http.StatusAccepted)

		assert.NoError(t, c.AcceptHouseholdInvitation(s.ctx, s.exampleHousehold.ID, s.exampleHouseholdInvitation.ID, t.Name()))
	})

	s.Run("with invalid household ID", func() {
		t := s.T()

		s.exampleHousehold.BelongsToUser = ""

		spec := newRequestSpec(false, http.MethodPut, "", expectedPath, s.exampleHousehold.ID, s.exampleHouseholdInvitation.ID)
		c, _ := buildTestClientWithStatusCodeResponse(t, spec, http.StatusAccepted)

		assert.Error(t, c.AcceptHouseholdInvitation(s.ctx, "", s.exampleHouseholdInvitation.ID, t.Name()))
	})

	s.Run("with invalid household invitation ID", func() {
		t := s.T()

		s.exampleHousehold.BelongsToUser = ""

		spec := newRequestSpec(false, http.MethodPut, "", expectedPath, s.exampleHousehold.ID, s.exampleHouseholdInvitation.ID)
		c, _ := buildTestClientWithStatusCodeResponse(t, spec, http.StatusAccepted)

		assert.Error(t, c.AcceptHouseholdInvitation(s.ctx, s.exampleHousehold.ID, "", t.Name()))
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)

		assert.Error(t, c.AcceptHouseholdInvitation(s.ctx, s.exampleHousehold.ID, s.exampleHouseholdInvitation.ID, t.Name()))
	})

	s.Run("with error executing request", func() {
		t := s.T()

		c, _ := buildTestClientThatWaitsTooLong(t)

		assert.Error(t, c.AcceptHouseholdInvitation(s.ctx, s.exampleHousehold.ID, s.exampleHouseholdInvitation.ID, t.Name()))
	})
}

func (s *householdInvitationsTestSuite) TestClient_RejectHouseholdInvitation() {
	const expectedPath = "/api/v1/households/%s/invitations/%s/reject"

	s.Run("standard", func() {
		t := s.T()

		s.exampleHousehold.BelongsToUser = ""

		spec := newRequestSpec(false, http.MethodPut, "", expectedPath, s.exampleHousehold.ID, s.exampleHouseholdInvitation.ID)
		c, _ := buildTestClientWithStatusCodeResponse(t, spec, http.StatusAccepted)

		assert.NoError(t, c.RejectHouseholdInvitation(s.ctx, s.exampleHousehold.ID, s.exampleHouseholdInvitation.ID, t.Name()))
	})

	s.Run("with invalid household ID", func() {
		t := s.T()

		s.exampleHousehold.BelongsToUser = ""

		spec := newRequestSpec(false, http.MethodPut, "", expectedPath, s.exampleHousehold.ID, s.exampleHouseholdInvitation.ID)
		c, _ := buildTestClientWithStatusCodeResponse(t, spec, http.StatusAccepted)

		assert.Error(t, c.RejectHouseholdInvitation(s.ctx, "", s.exampleHouseholdInvitation.ID, t.Name()))
	})

	s.Run("with invalid household invitation ID", func() {
		t := s.T()

		s.exampleHousehold.BelongsToUser = ""

		spec := newRequestSpec(false, http.MethodPut, "", expectedPath, s.exampleHousehold.ID, s.exampleHouseholdInvitation.ID)
		c, _ := buildTestClientWithStatusCodeResponse(t, spec, http.StatusAccepted)

		assert.Error(t, c.RejectHouseholdInvitation(s.ctx, s.exampleHousehold.ID, "", t.Name()))
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)

		assert.Error(t, c.RejectHouseholdInvitation(s.ctx, s.exampleHousehold.ID, s.exampleHouseholdInvitation.ID, t.Name()))
	})

	s.Run("with error executing request", func() {
		t := s.T()

		c, _ := buildTestClientThatWaitsTooLong(t)

		assert.Error(t, c.RejectHouseholdInvitation(s.ctx, s.exampleHousehold.ID, s.exampleHouseholdInvitation.ID, t.Name()))
	})
}
