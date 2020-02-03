package dbclient

import (
	"context"
	"testing"

	models "gitlab.com/prixfixe/prixfixe/models/v1"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestClient_GetInvitation(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		exampleInvitationID := uint64(123)
		exampleUserID := uint64(123)
		expected := &models.Invitation{}

		c, mockDB := buildTestClient()
		mockDB.InvitationDataManager.On("GetInvitation", mock.Anything, exampleInvitationID, exampleUserID).Return(expected, nil)

		actual, err := c.GetInvitation(context.Background(), exampleInvitationID, exampleUserID)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mockDB.AssertExpectations(t)
	})
}

func TestClient_GetInvitationCount(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		expected := uint64(321)
		exampleUserID := uint64(123)

		c, mockDB := buildTestClient()
		mockDB.InvitationDataManager.On("GetInvitationCount", mock.Anything, models.DefaultQueryFilter(), exampleUserID).Return(expected, nil)

		actual, err := c.GetInvitationCount(context.Background(), models.DefaultQueryFilter(), exampleUserID)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mockDB.AssertExpectations(t)
	})

	T.Run("with nil filter", func(t *testing.T) {
		expected := uint64(321)
		exampleUserID := uint64(123)

		c, mockDB := buildTestClient()
		mockDB.InvitationDataManager.On("GetInvitationCount", mock.Anything, (*models.QueryFilter)(nil), exampleUserID).Return(expected, nil)

		actual, err := c.GetInvitationCount(context.Background(), nil, exampleUserID)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mockDB.AssertExpectations(t)
	})
}

func TestClient_GetAllInvitationsCount(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		expected := uint64(321)
		c, mockDB := buildTestClient()
		mockDB.InvitationDataManager.On("GetAllInvitationsCount", mock.Anything).Return(expected, nil)

		actual, err := c.GetAllInvitationsCount(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mockDB.AssertExpectations(t)
	})
}

func TestClient_GetInvitations(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		exampleUserID := uint64(123)
		c, mockDB := buildTestClient()
		expected := &models.InvitationList{}

		mockDB.InvitationDataManager.On("GetInvitations", mock.Anything, models.DefaultQueryFilter(), exampleUserID).Return(expected, nil)

		actual, err := c.GetInvitations(context.Background(), models.DefaultQueryFilter(), exampleUserID)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mockDB.AssertExpectations(t)
	})

	T.Run("with nil filter", func(t *testing.T) {
		exampleUserID := uint64(123)
		c, mockDB := buildTestClient()
		expected := &models.InvitationList{}

		mockDB.InvitationDataManager.On("GetInvitations", mock.Anything, (*models.QueryFilter)(nil), exampleUserID).Return(expected, nil)

		actual, err := c.GetInvitations(context.Background(), nil, exampleUserID)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mockDB.AssertExpectations(t)
	})
}

func TestClient_CreateInvitation(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		exampleInput := &models.InvitationCreationInput{}
		c, mockDB := buildTestClient()
		expected := &models.Invitation{}

		mockDB.InvitationDataManager.On("CreateInvitation", mock.Anything, exampleInput).Return(expected, nil)

		actual, err := c.CreateInvitation(context.Background(), exampleInput)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mockDB.AssertExpectations(t)
	})
}

func TestClient_UpdateInvitation(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		exampleInput := &models.Invitation{}
		c, mockDB := buildTestClient()
		var expected error

		mockDB.InvitationDataManager.On("UpdateInvitation", mock.Anything, exampleInput).Return(expected)

		err := c.UpdateInvitation(context.Background(), exampleInput)
		assert.NoError(t, err)
	})
}

func TestClient_ArchiveInvitation(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		exampleUserID := uint64(123)
		exampleInvitationID := uint64(123)
		var expected error

		c, mockDB := buildTestClient()
		mockDB.InvitationDataManager.On("ArchiveInvitation", mock.Anything, exampleInvitationID, exampleUserID).Return(expected)

		err := c.ArchiveInvitation(context.Background(), exampleUserID, exampleInvitationID)
		assert.NoError(t, err)
	})
}
