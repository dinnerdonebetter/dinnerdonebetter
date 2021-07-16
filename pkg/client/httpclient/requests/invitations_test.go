package requests

import (
	"net/http"
	"testing"

	"gitlab.com/prixfixe/prixfixe/pkg/types"
	"gitlab.com/prixfixe/prixfixe/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuilder_BuildInvitationExistsRequest(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/invitations/%d"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleInvitation := fakes.BuildFakeInvitation()

		actual, err := helper.builder.BuildInvitationExistsRequest(helper.ctx, exampleInvitation.ID)
		spec := newRequestSpec(true, http.MethodHead, "", expectedPathFormat, exampleInvitation.ID)

		assert.NoError(t, err)
		assertRequestQuality(t, actual, spec)
	})

	T.Run("with invalid invitation ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		actual, err := helper.builder.BuildInvitationExistsRequest(helper.ctx, 0)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		exampleInvitation := fakes.BuildFakeInvitation()

		actual, err := helper.builder.BuildInvitationExistsRequest(helper.ctx, exampleInvitation.ID)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildGetInvitationRequest(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/invitations/%d"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleInvitation := fakes.BuildFakeInvitation()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, exampleInvitation.ID)

		actual, err := helper.builder.BuildGetInvitationRequest(helper.ctx, exampleInvitation.ID)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with invalid invitation ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		actual, err := helper.builder.BuildGetInvitationRequest(helper.ctx, 0)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		exampleInvitation := fakes.BuildFakeInvitation()

		actual, err := helper.builder.BuildGetInvitationRequest(helper.ctx, exampleInvitation.ID)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildGetInvitationsRequest(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/invitations"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		filter := (*types.QueryFilter)(nil)
		spec := newRequestSpec(true, http.MethodGet, "includeArchived=false&limit=20&page=1&sortBy=asc", expectedPathFormat)

		actual, err := helper.builder.BuildGetInvitationsRequest(helper.ctx, filter)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		filter := (*types.QueryFilter)(nil)

		actual, err := helper.builder.BuildGetInvitationsRequest(helper.ctx, filter)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildCreateInvitationRequest(T *testing.T) {
	T.Parallel()

	const expectedPath = "/api/v1/invitations"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleInput := fakes.BuildFakeInvitationCreationInput()

		spec := newRequestSpec(false, http.MethodPost, "", expectedPath)

		actual, err := helper.builder.BuildCreateInvitationRequest(helper.ctx, exampleInput)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		actual, err := helper.builder.BuildCreateInvitationRequest(helper.ctx, nil)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		actual, err := helper.builder.BuildCreateInvitationRequest(helper.ctx, &types.InvitationCreationInput{})
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		exampleInput := fakes.BuildFakeInvitationCreationInput()

		actual, err := helper.builder.BuildCreateInvitationRequest(helper.ctx, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildUpdateInvitationRequest(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/invitations/%d"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleInvitation := fakes.BuildFakeInvitation()

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, exampleInvitation.ID)

		actual, err := helper.builder.BuildUpdateInvitationRequest(helper.ctx, exampleInvitation)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		actual, err := helper.builder.BuildUpdateInvitationRequest(helper.ctx, nil)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		exampleInvitation := fakes.BuildFakeInvitation()

		actual, err := helper.builder.BuildUpdateInvitationRequest(helper.ctx, exampleInvitation)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildArchiveInvitationRequest(T *testing.T) {
	T.Parallel()

	const expectedPathFormat = "/api/v1/invitations/%d"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleInvitation := fakes.BuildFakeInvitation()

		spec := newRequestSpec(true, http.MethodDelete, "", expectedPathFormat, exampleInvitation.ID)

		actual, err := helper.builder.BuildArchiveInvitationRequest(helper.ctx, exampleInvitation.ID)
		assert.NoError(t, err)

		assertRequestQuality(t, actual, spec)
	})

	T.Run("with invalid invitation ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		actual, err := helper.builder.BuildArchiveInvitationRequest(helper.ctx, 0)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		exampleInvitation := fakes.BuildFakeInvitation()

		actual, err := helper.builder.BuildArchiveInvitationRequest(helper.ctx, exampleInvitation.ID)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func TestBuilder_BuildGetAuditLogForInvitationRequest(T *testing.T) {
	T.Parallel()

	const expectedPath = "/api/v1/invitations/%d/audit"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		exampleInvitation := fakes.BuildFakeInvitation()

		actual, err := helper.builder.BuildGetAuditLogForInvitationRequest(helper.ctx, exampleInvitation.ID)
		require.NotNil(t, actual)
		assert.NoError(t, err)

		spec := newRequestSpec(true, http.MethodGet, "", expectedPath, exampleInvitation.ID)
		assertRequestQuality(t, actual, spec)
	})

	T.Run("with invalid invitation ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()

		actual, err := helper.builder.BuildGetAuditLogForInvitationRequest(helper.ctx, 0)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid request builder", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper()
		helper.builder = buildTestRequestBuilderWithInvalidURL()

		exampleInvitation := fakes.BuildFakeInvitation()

		actual, err := helper.builder.BuildGetAuditLogForInvitationRequest(helper.ctx, exampleInvitation.ID)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}
