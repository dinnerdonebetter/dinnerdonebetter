package dbclient

import (
	"context"
	"testing"

	models "gitlab.com/prixfixe/prixfixe/models/v1"
	fakemodels "gitlab.com/prixfixe/prixfixe/models/v1/fake"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestClient_InvitationExists(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fakemodels.BuildFakeUser()
		exampleInvitation := fakemodels.BuildFakeInvitation()
		exampleInvitation.BelongsToUser = exampleUser.ID

		c, mockDB := buildTestClient()
		mockDB.InvitationDataManager.On("InvitationExists", mock.Anything, exampleInvitation.ID).Return(true, nil)

		actual, err := c.InvitationExists(ctx, exampleInvitation.ID)
		assert.NoError(t, err)
		assert.True(t, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_GetInvitation(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fakemodels.BuildFakeUser()
		exampleInvitation := fakemodels.BuildFakeInvitation()
		exampleInvitation.BelongsToUser = exampleUser.ID

		c, mockDB := buildTestClient()
		mockDB.InvitationDataManager.On("GetInvitation", mock.Anything, exampleInvitation.ID).Return(exampleInvitation, nil)

		actual, err := c.GetInvitation(ctx, exampleInvitation.ID)
		assert.NoError(t, err)
		assert.Equal(t, exampleInvitation, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_GetAllInvitationsCount(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		exampleCount := uint64(123)

		c, mockDB := buildTestClient()
		mockDB.InvitationDataManager.On("GetAllInvitationsCount", mock.Anything).Return(exampleCount, nil)

		actual, err := c.GetAllInvitationsCount(ctx)
		assert.NoError(t, err)
		assert.Equal(t, exampleCount, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_GetInvitations(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		filter := models.DefaultQueryFilter()
		exampleInvitationList := fakemodels.BuildFakeInvitationList()

		c, mockDB := buildTestClient()
		mockDB.InvitationDataManager.On("GetInvitations", mock.Anything, filter).Return(exampleInvitationList, nil)

		actual, err := c.GetInvitations(ctx, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleInvitationList, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with nil filter", func(t *testing.T) {
		ctx := context.Background()

		filter := (*models.QueryFilter)(nil)
		exampleInvitationList := fakemodels.BuildFakeInvitationList()

		c, mockDB := buildTestClient()
		mockDB.InvitationDataManager.On("GetInvitations", mock.Anything, filter).Return(exampleInvitationList, nil)

		actual, err := c.GetInvitations(ctx, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleInvitationList, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_CreateInvitation(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fakemodels.BuildFakeUser()
		exampleInvitation := fakemodels.BuildFakeInvitation()
		exampleInvitation.BelongsToUser = exampleUser.ID
		exampleInput := fakemodels.BuildFakeInvitationCreationInputFromInvitation(exampleInvitation)

		c, mockDB := buildTestClient()
		mockDB.InvitationDataManager.On("CreateInvitation", mock.Anything, exampleInput).Return(exampleInvitation, nil)

		actual, err := c.CreateInvitation(ctx, exampleInput)
		assert.NoError(t, err)
		assert.Equal(t, exampleInvitation, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_UpdateInvitation(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()
		var expected error

		exampleUser := fakemodels.BuildFakeUser()
		exampleInvitation := fakemodels.BuildFakeInvitation()
		exampleInvitation.BelongsToUser = exampleUser.ID

		c, mockDB := buildTestClient()

		mockDB.InvitationDataManager.On("UpdateInvitation", mock.Anything, exampleInvitation).Return(expected)

		err := c.UpdateInvitation(ctx, exampleInvitation)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_ArchiveInvitation(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		var expected error

		exampleUser := fakemodels.BuildFakeUser()
		exampleInvitation := fakemodels.BuildFakeInvitation()
		exampleInvitation.BelongsToUser = exampleUser.ID

		c, mockDB := buildTestClient()
		mockDB.InvitationDataManager.On("ArchiveInvitation", mock.Anything, exampleInvitation.ID, exampleInvitation.BelongsToUser).Return(expected)

		err := c.ArchiveInvitation(ctx, exampleInvitation.ID, exampleInvitation.BelongsToUser)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}
