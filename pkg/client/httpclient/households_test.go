package httpclient

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/prixfixeco/api_server/pkg/types"
	"github.com/prixfixeco/api_server/pkg/types/fakes"
)

func TestHouseholds(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(householdsTestSuite))
}

type householdsTestSuite struct {
	suite.Suite

	ctx                  context.Context
	exampleHousehold     *types.Household
	exampleUser          *types.User
	exampleHouseholdList *types.HouseholdList
}

var _ suite.SetupTestSuite = (*householdsTestSuite)(nil)

func (s *householdsTestSuite) SetupTest() {
	s.ctx = context.Background()
	s.exampleUser = fakes.BuildFakeUser()
	s.exampleHousehold = fakes.BuildFakeHousehold()
	s.exampleHousehold.BelongsToUser = s.exampleUser.ID
	s.exampleHouseholdList = fakes.BuildFakeHouseholdList()
}

func (s *householdsTestSuite) TestClient_SwitchActiveHousehold() {
	const expectedPath = "/users/household/select"

	s.Run("standard", func() {
		t := s.T()

		s.exampleHousehold.BelongsToUser = ""

		spec := newRequestSpec(false, http.MethodPost, "", expectedPath)
		c, _ := buildTestClientWithStatusCodeResponse(t, spec, http.StatusAccepted)
		c.authMethod = cookieAuthMethod

		assert.NoError(t, c.SwitchActiveHousehold(s.ctx, s.exampleHousehold.ID))
	})

	s.Run("with invalid household ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		c.authMethod = cookieAuthMethod

		assert.Error(t, c.SwitchActiveHousehold(s.ctx, ""))
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)
		c.authMethod = cookieAuthMethod

		assert.Error(t, c.SwitchActiveHousehold(s.ctx, s.exampleHousehold.ID))
	})

	s.Run("with error executing request", func() {
		t := s.T()

		c, _ := buildTestClientThatWaitsTooLong(t)
		c.authMethod = cookieAuthMethod

		assert.Error(t, c.SwitchActiveHousehold(s.ctx, s.exampleHousehold.ID))
	})
}

func (s *householdsTestSuite) TestClient_GetCurrentHousehold() {
	const expectedPathFormat = "/api/v1/households/current"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat)

		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleHousehold)

		actual, err := c.GetCurrentHousehold(s.ctx)

		for i := range actual.Members {
			actual.Members[i].BelongsToUser.TwoFactorSecretVerifiedAt = s.exampleHousehold.Members[i].BelongsToUser.TwoFactorSecretVerifiedAt
		}

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, s.exampleHousehold, actual)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetCurrentHousehold(s.ctx)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat)

		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetCurrentHousehold(s.ctx)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *householdsTestSuite) TestClient_GetHousehold() {
	const expectedPathFormat = "/api/v1/households/%s"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, s.exampleHousehold.ID)

		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleHousehold)

		actual, err := c.GetHousehold(s.ctx, s.exampleHousehold.ID)

		for i := range actual.Members {
			actual.Members[i].BelongsToUser.TwoFactorSecretVerifiedAt = s.exampleHousehold.Members[i].BelongsToUser.TwoFactorSecretVerifiedAt
		}

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, s.exampleHousehold, actual)
	})

	s.Run("with invalid household ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		actual, err := c.GetHousehold(s.ctx, "")
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetHousehold(s.ctx, s.exampleHousehold.ID)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, s.exampleHousehold.ID)

		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetHousehold(s.ctx, s.exampleHousehold.ID)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *householdsTestSuite) TestClient_GetHouseholds() {
	const expectedPath = "/api/v1/households"

	spec := newRequestSpec(true, http.MethodGet, "limit=20&page=1&sortBy=asc", expectedPath)
	filter := (*types.QueryFilter)(nil)

	s.Run("standard", func() {
		t := s.T()

		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleHouseholdList)
		actual, err := c.GetHouseholds(s.ctx, filter)

		for i, household := range actual.Households {
			for j := range household.Members {
				actual.Households[i].Members[j].BelongsToUser.TwoFactorSecretVerifiedAt = s.exampleHouseholdList.Households[i].Members[j].BelongsToUser.TwoFactorSecretVerifiedAt
			}
		}

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, s.exampleHouseholdList, actual)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetHouseholds(s.ctx, filter)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetHouseholds(s.ctx, filter)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *householdsTestSuite) TestClient_CreateHousehold() {
	const expectedPath = "/api/v1/households"

	spec := newRequestSpec(false, http.MethodPost, "", expectedPath)

	s.Run("standard", func() {
		t := s.T()

		s.exampleHousehold.BelongsToUser = ""
		exampleInput := fakes.BuildFakeHouseholdCreationRequestInputFromHousehold(s.exampleHousehold)

		c := buildTestClientWithRequestBodyValidation(t, spec, exampleInput, exampleInput, s.exampleHousehold)
		actual, err := c.CreateHousehold(s.ctx, exampleInput)

		for i := range actual.Members {
			actual.Members[i].BelongsToUser.TwoFactorSecretVerifiedAt = s.exampleHousehold.Members[i].BelongsToUser.TwoFactorSecretVerifiedAt
		}

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, s.exampleHousehold, actual)
	})

	s.Run("with nil input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		actual, err := c.CreateHousehold(s.ctx, nil)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with invalid input", func() {
		t := s.T()

		exampleInput := &types.HouseholdCreationRequestInput{}
		c, _ := buildSimpleTestClient(t)

		actual, err := c.CreateHousehold(s.ctx, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		exampleInput := fakes.BuildFakeHouseholdCreationRequestInputFromHousehold(s.exampleHousehold)
		c := buildTestClientWithInvalidURL(t)

		actual, err := c.CreateHousehold(s.ctx, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		s.exampleHousehold.BelongsToUser = ""
		exampleInput := fakes.BuildFakeHouseholdCreationRequestInputFromHousehold(s.exampleHousehold)

		c, _ := buildTestClientThatWaitsTooLong(t)

		actual, err := c.CreateHousehold(s.ctx, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *householdsTestSuite) TestClient_UpdateHousehold() {
	const expectedPathFormat = "/api/v1/households/%s"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, s.exampleHousehold.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleHousehold)

		assert.NoError(t, c.UpdateHousehold(s.ctx, s.exampleHousehold), "no error should be returned")
	})

	s.Run("with nil input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		assert.Error(t, c.UpdateHousehold(s.ctx, nil), "error should be returned")
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)

		err := c.UpdateHousehold(s.ctx, s.exampleHousehold)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		c, _ := buildTestClientThatWaitsTooLong(t)

		assert.Error(t, c.UpdateHousehold(s.ctx, s.exampleHousehold), "error should be returned")
	})
}

func (s *householdsTestSuite) TestClient_ArchiveHousehold() {
	const expectedPathFormat = "/api/v1/households/%s"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodDelete, "", expectedPathFormat, s.exampleHousehold.ID)
		c, _ := buildTestClientWithStatusCodeResponse(t, spec, http.StatusOK)

		assert.NoError(t, c.ArchiveHousehold(s.ctx, s.exampleHousehold.ID), "no error should be returned")
	})

	s.Run("with invalid household ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		assert.Error(t, c.ArchiveHousehold(s.ctx, ""), "no error should be returned")
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)
		assert.Error(t, c.ArchiveHousehold(s.ctx, s.exampleHousehold.ID), "error should be returned")
	})

	s.Run("with error executing request", func() {
		t := s.T()

		c, _ := buildTestClientThatWaitsTooLong(t)

		assert.Error(t, c.ArchiveHousehold(s.ctx, s.exampleHousehold.ID), "no error should be returned")
	})
}

func (s *householdsTestSuite) TestClient_InviteUserToHousehold() {
	const expectedPathFormat = "/api/v1/households/%s/invite"

	s.Run("standard", func() {
		t := s.T()

		invitation := fakes.BuildFakeHouseholdInvitation()
		invitation.FromUser.TwoFactorSecret = ""

		exampleInput := fakes.BuildFakeHouseholdInvitationCreationInputFromHouseholdInvitation(invitation)
		spec := newRequestSpec(false, http.MethodPost, "", expectedPathFormat, exampleInput.DestinationHouseholdID)
		c, _ := buildTestClientWithJSONResponse(t, spec, invitation)

		householdInvite, err := c.InviteUserToHousehold(s.ctx, exampleInput)
		assert.Equal(t, invitation, householdInvite)
		assert.NoError(t, err)
	})

	s.Run("with nil input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		householdInvite, err := c.InviteUserToHousehold(s.ctx, nil)
		assert.Nil(t, householdInvite)
		assert.Error(t, err)
	})

	s.Run("with invalid input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		householdInvite, err := c.InviteUserToHousehold(s.ctx, &types.HouseholdInvitationCreationRequestInput{})
		assert.Nil(t, householdInvite)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)
		exampleInput := fakes.BuildFakeHouseholdInvitationCreationRequestInput()

		householdInvite, err := c.InviteUserToHousehold(s.ctx, exampleInput)
		assert.Nil(t, householdInvite)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		exampleInput := fakes.BuildFakeHouseholdInvitationCreationRequestInput()
		c, _ := buildTestClientThatWaitsTooLong(t)

		householdInvite, err := c.InviteUserToHousehold(s.ctx, exampleInput)
		assert.Nil(t, householdInvite)
		assert.Error(t, err)
	})
}

func (s *householdsTestSuite) TestClient_MarkAsDefault() {
	const expectedPathFormat = "/api/v1/households/%s/default"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodPost, "", expectedPathFormat, s.exampleHousehold.ID)
		c, _ := buildTestClientWithStatusCodeResponse(t, spec, http.StatusOK)

		assert.NoError(t, c.MarkAsDefault(s.ctx, s.exampleHousehold.ID))
	})

	s.Run("with invalid household ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		assert.Error(t, c.MarkAsDefault(s.ctx, ""))
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)

		assert.Error(t, c.MarkAsDefault(s.ctx, s.exampleHousehold.ID))
	})

	s.Run("with error executing request", func() {
		t := s.T()

		c, _ := buildTestClientThatWaitsTooLong(t)

		assert.Error(t, c.MarkAsDefault(s.ctx, s.exampleHousehold.ID))
	})
}

func (s *householdsTestSuite) TestClient_RemoveUserFromHousehold() {
	const expectedPathFormat = "/api/v1/households/%s/members/%s"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodDelete, "", expectedPathFormat, s.exampleHousehold.ID, s.exampleUser.ID)
		c, _ := buildTestClientWithStatusCodeResponse(t, spec, http.StatusOK)

		assert.NoError(t, c.RemoveUserFromHousehold(s.ctx, s.exampleHousehold.ID, s.exampleUser.ID))
	})

	s.Run("with invalid household ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		assert.Error(t, c.RemoveUserFromHousehold(s.ctx, "", s.exampleUser.ID))
	})

	s.Run("with invalid user ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		assert.Error(t, c.RemoveUserFromHousehold(s.ctx, s.exampleHousehold.ID, ""))
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)

		assert.Error(t, c.RemoveUserFromHousehold(s.ctx, s.exampleHousehold.ID, s.exampleUser.ID))
	})

	s.Run("with error executing request", func() {
		t := s.T()

		c, _ := buildTestClientThatWaitsTooLong(t)

		assert.Error(t, c.RemoveUserFromHousehold(s.ctx, s.exampleHousehold.ID, s.exampleUser.ID))
	})
}

func (s *householdsTestSuite) TestClient_ModifyMemberPermissions() {
	const expectedPathFormat = "/api/v1/households/%s/members/%s/permissions"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(false, http.MethodPatch, "", expectedPathFormat, s.exampleHousehold.ID, s.exampleUser.ID)
		c, _ := buildTestClientWithStatusCodeResponse(t, spec, http.StatusOK)
		exampleInput := fakes.BuildFakeUserPermissionModificationInput()

		assert.NoError(t, c.ModifyMemberPermissions(s.ctx, s.exampleHousehold.ID, s.exampleUser.ID, exampleInput))
	})

	s.Run("with invalid household ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		exampleInput := fakes.BuildFakeUserPermissionModificationInput()

		assert.Error(t, c.ModifyMemberPermissions(s.ctx, "", s.exampleUser.ID, exampleInput))
	})

	s.Run("with invalid user ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		exampleInput := fakes.BuildFakeUserPermissionModificationInput()

		assert.Error(t, c.ModifyMemberPermissions(s.ctx, s.exampleHousehold.ID, "", exampleInput))
	})

	s.Run("with nil input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		assert.Error(t, c.ModifyMemberPermissions(s.ctx, s.exampleHousehold.ID, s.exampleUser.ID, nil))
	})

	s.Run("with invalid input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		exampleInput := &types.ModifyUserPermissionsInput{}

		assert.Error(t, c.ModifyMemberPermissions(s.ctx, s.exampleHousehold.ID, s.exampleUser.ID, exampleInput))
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)
		exampleInput := fakes.BuildFakeUserPermissionModificationInput()

		assert.Error(t, c.ModifyMemberPermissions(s.ctx, s.exampleHousehold.ID, s.exampleUser.ID, exampleInput))
	})

	s.Run("with error executing request", func() {
		t := s.T()

		c, _ := buildTestClientThatWaitsTooLong(t)
		exampleInput := fakes.BuildFakeUserPermissionModificationInput()

		assert.Error(t, c.ModifyMemberPermissions(s.ctx, s.exampleHousehold.ID, s.exampleUser.ID, exampleInput))
	})
}

func (s *householdsTestSuite) TestClient_TransferHouseholdOwnership() {
	const expectedPathFormat = "/api/v1/households/%s/transfer"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(false, http.MethodPost, "", expectedPathFormat, s.exampleHousehold.ID)
		c, _ := buildTestClientWithStatusCodeResponse(t, spec, http.StatusOK)
		exampleInput := fakes.BuildFakeTransferHouseholdOwnershipInput()

		assert.NoError(t, c.TransferHouseholdOwnership(s.ctx, s.exampleHousehold.ID, exampleInput))
	})

	s.Run("with invalid household ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		exampleInput := fakes.BuildFakeTransferHouseholdOwnershipInput()

		assert.Error(t, c.TransferHouseholdOwnership(s.ctx, "", exampleInput))
	})

	s.Run("with nil input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		assert.Error(t, c.TransferHouseholdOwnership(s.ctx, s.exampleHousehold.ID, nil))
	})

	s.Run("with invalid input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		exampleInput := &types.HouseholdOwnershipTransferInput{}

		assert.Error(t, c.TransferHouseholdOwnership(s.ctx, s.exampleHousehold.ID, exampleInput))
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)
		exampleInput := fakes.BuildFakeTransferHouseholdOwnershipInput()

		assert.Error(t, c.TransferHouseholdOwnership(s.ctx, s.exampleHousehold.ID, exampleInput))
	})

	s.Run("with error executing request", func() {
		t := s.T()

		c, _ := buildTestClientThatWaitsTooLong(t)
		exampleInput := fakes.BuildFakeTransferHouseholdOwnershipInput()

		assert.Error(t, c.TransferHouseholdOwnership(s.ctx, s.exampleHousehold.ID, exampleInput))
	})
}
