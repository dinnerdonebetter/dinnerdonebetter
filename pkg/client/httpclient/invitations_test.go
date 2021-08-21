package httpclient

import (
	"context"
	"net/http"
	"testing"

	"gitlab.com/prixfixe/prixfixe/pkg/types"
	"gitlab.com/prixfixe/prixfixe/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

func TestInvitations(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(invitationsTestSuite))
}

type invitationsBaseSuite struct {
	suite.Suite

	ctx               context.Context
	exampleInvitation *types.Invitation
}

var _ suite.SetupTestSuite = (*invitationsBaseSuite)(nil)

func (s *invitationsBaseSuite) SetupTest() {
	s.ctx = context.Background()
	s.exampleInvitation = fakes.BuildFakeInvitation()
}

type invitationsTestSuite struct {
	suite.Suite

	invitationsBaseSuite
}

func (s *invitationsTestSuite) TestClient_InvitationExists() {
	const expectedPathFormat = "/api/v1/invitations/%d"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodHead, "", expectedPathFormat, s.exampleInvitation.ID)

		c, _ := buildTestClientWithStatusCodeResponse(t, spec, http.StatusOK)
		actual, err := c.InvitationExists(s.ctx, s.exampleInvitation.ID)

		assert.NoError(t, err)
		assert.True(t, actual)
	})

	s.Run("with invalid invitation ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		actual, err := c.InvitationExists(s.ctx, 0)

		assert.Error(t, err)
		assert.False(t, actual)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.InvitationExists(s.ctx, s.exampleInvitation.ID)

		assert.Error(t, err)
		assert.False(t, actual)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		c, _ := buildTestClientThatWaitsTooLong(t)
		actual, err := c.InvitationExists(s.ctx, s.exampleInvitation.ID)

		assert.Error(t, err)
		assert.False(t, actual)
	})
}

func (s *invitationsTestSuite) TestClient_GetInvitation() {
	const expectedPathFormat = "/api/v1/invitations/%d"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, s.exampleInvitation.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleInvitation)
		actual, err := c.GetInvitation(s.ctx, s.exampleInvitation.ID)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, s.exampleInvitation, actual)
	})

	s.Run("with invalid invitation ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetInvitation(s.ctx, 0)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetInvitation(s.ctx, s.exampleInvitation.ID)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, s.exampleInvitation.ID)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetInvitation(s.ctx, s.exampleInvitation.ID)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *invitationsTestSuite) TestClient_GetInvitations() {
	const expectedPath = "/api/v1/invitations"

	s.Run("standard", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)

		exampleInvitationList := fakes.BuildFakeInvitationList()

		spec := newRequestSpec(true, http.MethodGet, "includeArchived=false&limit=20&page=1&sortBy=asc", expectedPath)
		c, _ := buildTestClientWithJSONResponse(t, spec, exampleInvitationList)
		actual, err := c.GetInvitations(s.ctx, filter)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, exampleInvitationList, actual)
	})

	s.Run("with error building request", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetInvitations(s.ctx, filter)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)

		spec := newRequestSpec(true, http.MethodGet, "includeArchived=false&limit=20&page=1&sortBy=asc", expectedPath)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetInvitations(s.ctx, filter)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *invitationsTestSuite) TestClient_CreateInvitation() {
	const expectedPath = "/api/v1/invitations"

	s.Run("standard", func() {
		t := s.T()

		exampleInput := fakes.BuildFakeInvitationCreationInput()
		exampleInput.BelongsToHousehold = 0

		spec := newRequestSpec(false, http.MethodPost, "", expectedPath)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleInvitation)

		actual, err := c.CreateInvitation(s.ctx, exampleInput)
		require.NotNil(t, actual)
		assert.NoError(t, err)

		assert.Equal(t, s.exampleInvitation, actual)
	})

	s.Run("with nil input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		actual, err := c.CreateInvitation(s.ctx, nil)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with invalid input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		exampleInput := &types.InvitationCreationInput{}

		actual, err := c.CreateInvitation(s.ctx, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		exampleInput := fakes.BuildFakeInvitationCreationInputFromInvitation(s.exampleInvitation)

		c := buildTestClientWithInvalidURL(t)

		actual, err := c.CreateInvitation(s.ctx, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		exampleInput := fakes.BuildFakeInvitationCreationInputFromInvitation(s.exampleInvitation)
		c, _ := buildTestClientThatWaitsTooLong(t)

		actual, err := c.CreateInvitation(s.ctx, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *invitationsTestSuite) TestClient_UpdateInvitation() {
	const expectedPathFormat = "/api/v1/invitations/%d"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, s.exampleInvitation.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleInvitation)

		err := c.UpdateInvitation(s.ctx, s.exampleInvitation)
		assert.NoError(t, err)
	})

	s.Run("with nil input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		err := c.UpdateInvitation(s.ctx, nil)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)

		err := c.UpdateInvitation(s.ctx, s.exampleInvitation)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		c, _ := buildTestClientThatWaitsTooLong(t)

		err := c.UpdateInvitation(s.ctx, s.exampleInvitation)
		assert.Error(t, err)
	})
}

func (s *invitationsTestSuite) TestClient_ArchiveInvitation() {
	const expectedPathFormat = "/api/v1/invitations/%d"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodDelete, "", expectedPathFormat, s.exampleInvitation.ID)
		c, _ := buildTestClientWithStatusCodeResponse(t, spec, http.StatusOK)

		err := c.ArchiveInvitation(s.ctx, s.exampleInvitation.ID)
		assert.NoError(t, err)
	})

	s.Run("with invalid invitation ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		err := c.ArchiveInvitation(s.ctx, 0)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)

		err := c.ArchiveInvitation(s.ctx, s.exampleInvitation.ID)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		c, _ := buildTestClientThatWaitsTooLong(t)

		err := c.ArchiveInvitation(s.ctx, s.exampleInvitation.ID)
		assert.Error(t, err)
	})
}

func (s *invitationsTestSuite) TestClient_GetAuditLogForInvitation() {
	const (
		expectedPath   = "/api/v1/invitations/%d/audit"
		expectedMethod = http.MethodGet
	)

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, expectedMethod, "", expectedPath, s.exampleInvitation.ID)
		exampleAuditLogEntryList := fakes.BuildFakeAuditLogEntryList().Entries

		c, _ := buildTestClientWithJSONResponse(t, spec, exampleAuditLogEntryList)

		actual, err := c.GetAuditLogForInvitation(s.ctx, s.exampleInvitation.ID)
		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, exampleAuditLogEntryList, actual)
	})

	s.Run("with invalid invitation ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		actual, err := c.GetAuditLogForInvitation(s.ctx, 0)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)

		actual, err := c.GetAuditLogForInvitation(s.ctx, s.exampleInvitation.ID)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		c, _ := buildTestClientThatWaitsTooLong(t)

		actual, err := c.GetAuditLogForInvitation(s.ctx, s.exampleInvitation.ID)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}
