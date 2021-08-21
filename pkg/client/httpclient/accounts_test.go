package httpclient

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	"gitlab.com/prixfixe/prixfixe/internal/observability/keys"
	"gitlab.com/prixfixe/prixfixe/pkg/types"
	"gitlab.com/prixfixe/prixfixe/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
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

		s.exampleHousehold.BelongsToUser = 0

		spec := newRequestSpec(false, http.MethodPost, "", expectedPath)
		c, _ := buildTestClientWithStatusCodeResponse(t, spec, http.StatusAccepted)
		c.authMethod = cookieAuthMethod

		assert.NoError(t, c.SwitchActiveHousehold(s.ctx, s.exampleHousehold.ID))
	})

	s.Run("with invalid household ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		c.authMethod = cookieAuthMethod

		assert.Error(t, c.SwitchActiveHousehold(s.ctx, 0))
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

func (s *householdsTestSuite) TestClient_GetHousehold() {
	const expectedPathFormat = "/api/v1/households/%d"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, s.exampleHousehold.ID)

		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleHousehold)

		actual, err := c.GetHousehold(s.ctx, s.exampleHousehold.ID)
		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, s.exampleHousehold, actual)
	})

	s.Run("with invalid household ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		actual, err := c.GetHousehold(s.ctx, 0)
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

	spec := newRequestSpec(true, http.MethodGet, "includeArchived=false&limit=20&page=1&sortBy=asc", expectedPath)
	filter := (*types.QueryFilter)(nil)

	s.Run("standard", func() {
		t := s.T()

		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleHouseholdList)
		actual, err := c.GetHouseholds(s.ctx, filter)

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

		s.exampleHousehold.BelongsToUser = 0
		exampleInput := fakes.BuildFakeHouseholdCreationInputFromHousehold(s.exampleHousehold)

		c := buildTestClientWithRequestBodyValidation(t, spec, exampleInput, exampleInput, s.exampleHousehold)
		actual, err := c.CreateHousehold(s.ctx, exampleInput)

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

		exampleInput := &types.HouseholdCreationInput{}
		c, _ := buildSimpleTestClient(t)

		actual, err := c.CreateHousehold(s.ctx, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		exampleInput := fakes.BuildFakeHouseholdCreationInputFromHousehold(s.exampleHousehold)
		c := buildTestClientWithInvalidURL(t)

		actual, err := c.CreateHousehold(s.ctx, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		s.exampleHousehold.BelongsToUser = 0
		exampleInput := fakes.BuildFakeHouseholdCreationInputFromHousehold(s.exampleHousehold)

		c, _ := buildTestClientThatWaitsTooLong(t)

		actual, err := c.CreateHousehold(s.ctx, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *householdsTestSuite) TestClient_UpdateHousehold() {
	const expectedPathFormat = "/api/v1/households/%d"

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
	const expectedPathFormat = "/api/v1/households/%d"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodDelete, "", expectedPathFormat, s.exampleHousehold.ID)
		c, _ := buildTestClientWithStatusCodeResponse(t, spec, http.StatusOK)

		assert.NoError(t, c.ArchiveHousehold(s.ctx, s.exampleHousehold.ID), "no error should be returned")
	})

	s.Run("with invalid household ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		assert.Error(t, c.ArchiveHousehold(s.ctx, 0), "no error should be returned")
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

func (s *householdsTestSuite) TestClient_AddUserToHousehold() {
	const expectedPathFormat = "/api/v1/households/%d/member"

	s.Run("standard", func() {
		t := s.T()

		exampleInput := fakes.BuildFakeAddUserToHouseholdInput()
		spec := newRequestSpec(false, http.MethodPost, "", expectedPathFormat, exampleInput.HouseholdID)
		c, _ := buildTestClientWithStatusCodeResponse(t, spec, http.StatusOK)

		assert.NoError(t, c.AddUserToHousehold(s.ctx, exampleInput))
	})

	s.Run("with nil input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		assert.Error(t, c.AddUserToHousehold(s.ctx, nil))
	})

	s.Run("with invalid input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		assert.Error(t, c.AddUserToHousehold(s.ctx, &types.AddUserToHouseholdInput{}))
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)
		exampleInput := fakes.BuildFakeAddUserToHouseholdInput()

		assert.Error(t, c.AddUserToHousehold(s.ctx, exampleInput))
	})

	s.Run("with error executing request", func() {
		t := s.T()

		exampleInput := fakes.BuildFakeAddUserToHouseholdInput()
		c, _ := buildTestClientThatWaitsTooLong(t)

		assert.Error(t, c.AddUserToHousehold(s.ctx, exampleInput))
	})
}

func (s *householdsTestSuite) TestClient_MarkAsDefault() {
	const expectedPathFormat = "/api/v1/households/%d/default"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodPost, "", expectedPathFormat, s.exampleHousehold.ID)
		c, _ := buildTestClientWithStatusCodeResponse(t, spec, http.StatusOK)

		assert.NoError(t, c.MarkAsDefault(s.ctx, s.exampleHousehold.ID))
	})

	s.Run("with invalid household ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		assert.Error(t, c.MarkAsDefault(s.ctx, 0))
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
	const expectedPathFormat = "/api/v1/households/%d/members/%d"

	s.Run("standard", func() {
		t := s.T()

		query := url.Values{keys.ReasonKey: []string{t.Name()}}.Encode()
		spec := newRequestSpec(true, http.MethodDelete, query, expectedPathFormat, s.exampleHousehold.ID, s.exampleUser.ID)
		c, _ := buildTestClientWithStatusCodeResponse(t, spec, http.StatusOK)

		assert.NoError(t, c.RemoveUserFromHousehold(s.ctx, s.exampleHousehold.ID, s.exampleUser.ID, t.Name()))
	})

	s.Run("with invalid household ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		assert.Error(t, c.RemoveUserFromHousehold(s.ctx, 0, s.exampleUser.ID, t.Name()))
	})

	s.Run("with invalid user ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		assert.Error(t, c.RemoveUserFromHousehold(s.ctx, s.exampleHousehold.ID, 0, t.Name()))
	})

	s.Run("with invalid reason", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		assert.Error(t, c.RemoveUserFromHousehold(s.ctx, s.exampleHousehold.ID, s.exampleUser.ID, ""))
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)

		assert.Error(t, c.RemoveUserFromHousehold(s.ctx, s.exampleHousehold.ID, s.exampleUser.ID, t.Name()))
	})

	s.Run("with error executing request", func() {
		t := s.T()

		c, _ := buildTestClientThatWaitsTooLong(t)

		assert.Error(t, c.RemoveUserFromHousehold(s.ctx, s.exampleHousehold.ID, s.exampleUser.ID, t.Name()))
	})
}

func (s *householdsTestSuite) TestClient_ModifyMemberPermissions() {
	const expectedPathFormat = "/api/v1/households/%d/members/%d/permissions"

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

		assert.Error(t, c.ModifyMemberPermissions(s.ctx, 0, s.exampleUser.ID, exampleInput))
	})

	s.Run("with invalid user ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		exampleInput := fakes.BuildFakeUserPermissionModificationInput()

		assert.Error(t, c.ModifyMemberPermissions(s.ctx, s.exampleHousehold.ID, 0, exampleInput))
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
	const expectedPathFormat = "/api/v1/households/%d/transfer"

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

		assert.Error(t, c.TransferHouseholdOwnership(s.ctx, 0, exampleInput))
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

func (s *householdsTestSuite) TestClient_GetAuditLogForHousehold() {
	const (
		expectedPath   = "/api/v1/households/%d/audit"
		expectedMethod = http.MethodGet
	)

	s.Run("standard", func() {
		t := s.T()

		exampleAuditLogEntryList := fakes.BuildFakeAuditLogEntryList().Entries
		spec := newRequestSpec(true, expectedMethod, "", expectedPath, s.exampleHousehold.ID)

		c, _ := buildTestClientWithJSONResponse(t, spec, exampleAuditLogEntryList)

		actual, err := c.GetAuditLogForHousehold(s.ctx, s.exampleHousehold.ID)
		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, exampleAuditLogEntryList, actual)
	})

	s.Run("with invalid household ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		actual, err := c.GetAuditLogForHousehold(s.ctx, 0)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)

		actual, err := c.GetAuditLogForHousehold(s.ctx, s.exampleHousehold.ID)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		spec := newRequestSpec(true, expectedMethod, "", expectedPath, s.exampleHousehold.ID)
		c := buildTestClientWithInvalidResponse(t, spec)

		actual, err := c.GetAuditLogForHousehold(s.ctx, s.exampleHousehold.ID)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}
