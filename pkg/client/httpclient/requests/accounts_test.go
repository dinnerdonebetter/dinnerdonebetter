package requests

import (
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"gitlab.com/prixfixe/prixfixe/pkg/types"
	"gitlab.com/prixfixe/prixfixe/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuilder_BuildSwitchActiveHouseholdRequest(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/users/household/select"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		exampleHouseholdID := fakes.BuildFakeID()

		spec := newRequestSpec(false, http.MethodPost, "", expectedPathFormat)

		actual, err := helper.builder.BuildSwitchActiveHouseholdRequest(helper.ctx, exampleHouseholdID)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with invalid household ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		actual, err := helper.builder.BuildSwitchActiveHouseholdRequest(helper.ctx, 0)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()
		exampleHouseholdID := fakes.BuildFakeID()

		actual, err := helper.builder.BuildSwitchActiveHouseholdRequest(helper.ctx, exampleHouseholdID)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildGetHouseholdRequest(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/households/%d"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		exampleHouseholdID := fakes.BuildFakeID()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, exampleHouseholdID)

		actual, err := helper.builder.BuildGetHouseholdRequest(helper.ctx, exampleHouseholdID)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with invalid household ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		actual, err := helper.builder.BuildGetHouseholdRequest(helper.ctx, 0)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()
		exampleHouseholdID := fakes.BuildFakeID()

		actual, err := helper.builder.BuildGetHouseholdRequest(helper.ctx, exampleHouseholdID)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildGetHouseholdsRequest(T *testing.T) {
	T.Parallel()

	const expectedPath = "/api/v1/households"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		filter := (*types.QueryFilter)(nil)
		spec := newRequestSpec(true, http.MethodGet, "includeArchived=false&limit=20&page=1&sortBy=asc", expectedPath)

		actual, err := helper.builder.BuildGetHouseholdsRequest(helper.ctx, filter)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		filter := (*types.QueryFilter)(nil)

		actual, err := helper.builder.BuildGetHouseholdsRequest(helper.ctx, filter)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildCreateHouseholdRequest(T *testing.T) {
	T.Parallel()

	const expectedPath = "/api/v1/households"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		exampleHousehold := fakes.BuildFakeHousehold()
		exampleInput := fakes.BuildFakeHouseholdCreationInputFromHousehold(exampleHousehold)

		actual, err := helper.builder.BuildCreateHouseholdRequest(helper.ctx, exampleInput)
		assert.NoError(t, err)

		spec := newRequestSpec(false, http.MethodPost, "", expectedPath)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		actual, err := helper.builder.BuildCreateHouseholdRequest(helper.ctx, nil)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		actual, err := helper.builder.BuildCreateHouseholdRequest(helper.ctx, &types.HouseholdCreationInput{})
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()
		exampleHousehold := fakes.BuildFakeHousehold()
		exampleInput := fakes.BuildFakeHouseholdCreationInputFromHousehold(exampleHousehold)

		actual, err := helper.builder.BuildCreateHouseholdRequest(helper.ctx, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildUpdateHouseholdRequest(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/households/%d"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		exampleHousehold := fakes.BuildFakeHousehold()

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, exampleHousehold.ID)

		actual, err := helper.builder.BuildUpdateHouseholdRequest(helper.ctx, exampleHousehold)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		actual, err := helper.builder.BuildUpdateHouseholdRequest(helper.ctx, nil)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()
		exampleHousehold := fakes.BuildFakeHousehold()

		actual, err := helper.builder.BuildUpdateHouseholdRequest(helper.ctx, exampleHousehold)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildArchiveHouseholdRequest(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/households/%d"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		exampleHouseholdID := fakes.BuildFakeID()

		spec := newRequestSpec(true, http.MethodDelete, "", expectedPathFormat, exampleHouseholdID)

		actual, err := helper.builder.BuildArchiveHouseholdRequest(helper.ctx, exampleHouseholdID)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with invalid household ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		actual, err := helper.builder.BuildArchiveHouseholdRequest(helper.ctx, 0)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()
		exampleHouseholdID := fakes.BuildFakeID()

		actual, err := helper.builder.BuildArchiveHouseholdRequest(helper.ctx, exampleHouseholdID)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildAddUserRequest(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/households/%d/member"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleInput := fakes.BuildFakeAddUserToHouseholdInput()
		spec := newRequestSpec(false, http.MethodPost, "", expectedPathFormat, exampleInput.HouseholdID)

		actual, err := helper.builder.BuildAddUserRequest(helper.ctx, exampleInput)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		actual, err := helper.builder.BuildAddUserRequest(helper.ctx, nil)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		actual, err := helper.builder.BuildAddUserRequest(helper.ctx, &types.AddUserToHouseholdInput{})
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		exampleInput := fakes.BuildFakeAddUserToHouseholdInput()

		actual, err := helper.builder.BuildAddUserRequest(helper.ctx, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildMarkAsDefaultRequest(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/households/%d/default"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		exampleHouseholdID := fakes.BuildFakeID()

		spec := newRequestSpec(true, http.MethodPost, "", expectedPathFormat, exampleHouseholdID)

		actual, err := helper.builder.BuildMarkAsDefaultRequest(helper.ctx, exampleHouseholdID)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with invalid household ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		actual, err := helper.builder.BuildMarkAsDefaultRequest(helper.ctx, 0)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()
		exampleHouseholdID := fakes.BuildFakeID()

		actual, err := helper.builder.BuildMarkAsDefaultRequest(helper.ctx, exampleHouseholdID)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildRemoveUserRequest(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/households/%d/members/%d"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		exampleHouseholdID := fakes.BuildFakeID()

		reason := t.Name()
		expectedReason := url.QueryEscape(reason)
		spec := newRequestSpec(false, http.MethodDelete, fmt.Sprintf("reason=%s", expectedReason), expectedPathFormat, exampleHouseholdID, helper.exampleUser.ID)

		actual, err := helper.builder.BuildRemoveUserRequest(helper.ctx, exampleHouseholdID, helper.exampleUser.ID, reason)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with invalid household ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		reason := t.Name()

		actual, err := helper.builder.BuildRemoveUserRequest(helper.ctx, 0, helper.exampleUser.ID, reason)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid user ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		exampleHouseholdID := fakes.BuildFakeID()

		reason := t.Name()

		actual, err := helper.builder.BuildRemoveUserRequest(helper.ctx, exampleHouseholdID, 0, reason)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()
		exampleHouseholdID := fakes.BuildFakeID()

		reason := t.Name()

		actual, err := helper.builder.BuildRemoveUserRequest(helper.ctx, exampleHouseholdID, helper.exampleUser.ID, reason)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildModifyMemberPermissionsRequest(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/households/%d/members/%d/permissions"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		exampleHouseholdID := fakes.BuildFakeID()

		spec := newRequestSpec(false, http.MethodPatch, "", expectedPathFormat, exampleHouseholdID, helper.exampleUser.ID)
		exampleInput := fakes.BuildFakeUserPermissionModificationInput()

		actual, err := helper.builder.BuildModifyMemberPermissionsRequest(helper.ctx, exampleHouseholdID, helper.exampleUser.ID, exampleInput)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with invalid household ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleInput := fakes.BuildFakeUserPermissionModificationInput()

		actual, err := helper.builder.BuildModifyMemberPermissionsRequest(helper.ctx, 0, helper.exampleUser.ID, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid user ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		exampleHouseholdID := fakes.BuildFakeID()

		exampleInput := fakes.BuildFakeUserPermissionModificationInput()

		actual, err := helper.builder.BuildModifyMemberPermissionsRequest(helper.ctx, exampleHouseholdID, 0, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		exampleHouseholdID := fakes.BuildFakeID()

		actual, err := helper.builder.BuildModifyMemberPermissionsRequest(helper.ctx, exampleHouseholdID, helper.exampleUser.ID, nil)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		exampleHouseholdID := fakes.BuildFakeID()

		actual, err := helper.builder.BuildModifyMemberPermissionsRequest(helper.ctx, exampleHouseholdID, helper.exampleUser.ID, &types.ModifyUserPermissionsInput{})
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()
		exampleHouseholdID := fakes.BuildFakeID()

		exampleInput := fakes.BuildFakeUserPermissionModificationInput()

		actual, err := helper.builder.BuildModifyMemberPermissionsRequest(helper.ctx, exampleHouseholdID, helper.exampleUser.ID, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildTransferHouseholdOwnershipRequest(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/households/%d/transfer"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		exampleHouseholdID := fakes.BuildFakeID()

		spec := newRequestSpec(false, http.MethodPost, "", expectedPathFormat, exampleHouseholdID)
		exampleInput := fakes.BuildFakeTransferHouseholdOwnershipInput()

		actual, err := helper.builder.BuildTransferHouseholdOwnershipRequest(helper.ctx, exampleHouseholdID, exampleInput)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with invalid household ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleInput := fakes.BuildFakeTransferHouseholdOwnershipInput()

		actual, err := helper.builder.BuildTransferHouseholdOwnershipRequest(helper.ctx, 0, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		exampleHouseholdID := fakes.BuildFakeID()

		actual, err := helper.builder.BuildTransferHouseholdOwnershipRequest(helper.ctx, exampleHouseholdID, nil)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		exampleHouseholdID := fakes.BuildFakeID()

		actual, err := helper.builder.BuildTransferHouseholdOwnershipRequest(helper.ctx, exampleHouseholdID, &types.HouseholdOwnershipTransferInput{})
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()
		exampleHouseholdID := fakes.BuildFakeID()

		exampleInput := fakes.BuildFakeTransferHouseholdOwnershipInput()

		actual, err := helper.builder.BuildTransferHouseholdOwnershipRequest(helper.ctx, exampleHouseholdID, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildGetAuditLogForHouseholdRequest(T *testing.T) {
	T.Parallel()

	const expectedPath = "/api/v1/households/%d/audit"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		exampleHouseholdID := fakes.BuildFakeID()

		actual, err := helper.builder.BuildGetAuditLogForHouseholdRequest(helper.ctx, exampleHouseholdID)
		require.NotNil(t, actual)
		assert.NoError(t, err)

		spec := newRequestSpec(true, http.MethodGet, "", expectedPath, exampleHouseholdID)
		assertRequestQuality(t, actual, spec)
	})

	T.Run("with invalid household ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		actual, err := helper.builder.BuildGetAuditLogForHouseholdRequest(helper.ctx, 0)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()
		exampleHouseholdID := fakes.BuildFakeID()

		actual, err := helper.builder.BuildGetAuditLogForHouseholdRequest(helper.ctx, exampleHouseholdID)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}
