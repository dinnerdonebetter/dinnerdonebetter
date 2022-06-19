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

func TestAdmin(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(adminTestSuite))
}

type adminTestSuite struct {
	suite.Suite

	ctx                  context.Context
	exampleHousehold     *types.Household
	exampleHouseholdList *types.HouseholdList
}

var _ suite.SetupTestSuite = (*adminTestSuite)(nil)

func (s *adminTestSuite) SetupTest() {
	s.ctx = context.Background()
	s.exampleHousehold = fakes.BuildFakeHousehold()
	s.exampleHouseholdList = fakes.BuildFakeHouseholdList()
}

func (s *adminTestSuite) TestClient_UpdateAccountStatus() {
	const expectedPath = "/api/v1/admin/users/status"

	s.Run("standard", func() {
		t := s.T()

		exampleInput := fakes.BuildFakeUserAccountStatusUpdateInput()
		spec := newRequestSpec(false, http.MethodPost, "", expectedPath)
		c, _ := buildTestClientWithStatusCodeResponse(t, spec, http.StatusAccepted)

		assert.NoError(t, c.UpdateUserAccountStatus(s.ctx, exampleInput))
	})

	s.Run("with nil input", func() {
		t := s.T()

		spec := newRequestSpec(false, http.MethodPost, "", expectedPath)
		c, _ := buildTestClientWithStatusCodeResponse(t, spec, http.StatusAccepted)

		assert.Error(t, c.UpdateUserAccountStatus(s.ctx, nil))
	})

	s.Run("with invalid input", func() {
		t := s.T()

		exampleInput := &types.UserAccountStatusUpdateInput{}
		spec := newRequestSpec(false, http.MethodPost, "", expectedPath)
		c, _ := buildTestClientWithStatusCodeResponse(t, spec, http.StatusAccepted)

		assert.Error(t, c.UpdateUserAccountStatus(s.ctx, exampleInput))
	})

	s.Run("with bad request response", func() {
		t := s.T()

		exampleInput := fakes.BuildFakeUserAccountStatusUpdateInput()
		spec := newRequestSpec(false, http.MethodPost, "", expectedPath)
		c, _ := buildTestClientWithStatusCodeResponse(t, spec, http.StatusBadRequest)

		assert.Error(t, c.UpdateUserAccountStatus(s.ctx, exampleInput))
	})

	s.Run("with otherwise invalid status code response", func() {
		t := s.T()

		exampleInput := fakes.BuildFakeUserAccountStatusUpdateInput()
		spec := newRequestSpec(false, http.MethodPost, "", expectedPath)
		c, _ := buildTestClientWithStatusCodeResponse(t, spec, http.StatusInternalServerError)

		assert.Error(t, c.UpdateUserAccountStatus(s.ctx, exampleInput))
	})

	s.Run("with error building request", func() {
		t := s.T()

		exampleInput := fakes.BuildFakeUserAccountStatusUpdateInput()
		c := buildTestClientWithInvalidURL(t)

		assert.Error(t, c.UpdateUserAccountStatus(s.ctx, exampleInput))
	})

	s.Run("with timeout", func() {
		t := s.T()

		exampleInput := fakes.BuildFakeUserAccountStatusUpdateInput()
		c, _ := buildTestClientThatWaitsTooLong(t)

		assert.Error(t, c.UpdateUserAccountStatus(s.ctx, exampleInput))
	})
}
