package dbclient

import (
	"context"
	"testing"

	models "gitlab.com/prixfixe/prixfixe/models/v1"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestClient_GetIterationMedia(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		exampleIterationMediaID := uint64(123)
		exampleUserID := uint64(123)
		expected := &models.IterationMedia{}

		c, mockDB := buildTestClient()
		mockDB.IterationMediaDataManager.On("GetIterationMedia", mock.Anything, exampleIterationMediaID, exampleUserID).Return(expected, nil)

		actual, err := c.GetIterationMedia(context.Background(), exampleIterationMediaID, exampleUserID)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mockDB.AssertExpectations(t)
	})
}

func TestClient_GetIterationMediaCount(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		expected := uint64(321)
		exampleUserID := uint64(123)

		c, mockDB := buildTestClient()
		mockDB.IterationMediaDataManager.On("GetIterationMediaCount", mock.Anything, models.DefaultQueryFilter(), exampleUserID).Return(expected, nil)

		actual, err := c.GetIterationMediaCount(context.Background(), models.DefaultQueryFilter(), exampleUserID)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mockDB.AssertExpectations(t)
	})

	T.Run("with nil filter", func(t *testing.T) {
		expected := uint64(321)
		exampleUserID := uint64(123)

		c, mockDB := buildTestClient()
		mockDB.IterationMediaDataManager.On("GetIterationMediaCount", mock.Anything, (*models.QueryFilter)(nil), exampleUserID).Return(expected, nil)

		actual, err := c.GetIterationMediaCount(context.Background(), nil, exampleUserID)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mockDB.AssertExpectations(t)
	})
}

func TestClient_GetAllIterationMediasCount(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		expected := uint64(321)
		c, mockDB := buildTestClient()
		mockDB.IterationMediaDataManager.On("GetAllIterationMediasCount", mock.Anything).Return(expected, nil)

		actual, err := c.GetAllIterationMediasCount(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mockDB.AssertExpectations(t)
	})
}

func TestClient_GetIterationMedias(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		exampleUserID := uint64(123)
		c, mockDB := buildTestClient()
		expected := &models.IterationMediaList{}

		mockDB.IterationMediaDataManager.On("GetIterationMedias", mock.Anything, models.DefaultQueryFilter(), exampleUserID).Return(expected, nil)

		actual, err := c.GetIterationMedias(context.Background(), models.DefaultQueryFilter(), exampleUserID)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mockDB.AssertExpectations(t)
	})

	T.Run("with nil filter", func(t *testing.T) {
		exampleUserID := uint64(123)
		c, mockDB := buildTestClient()
		expected := &models.IterationMediaList{}

		mockDB.IterationMediaDataManager.On("GetIterationMedias", mock.Anything, (*models.QueryFilter)(nil), exampleUserID).Return(expected, nil)

		actual, err := c.GetIterationMedias(context.Background(), nil, exampleUserID)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mockDB.AssertExpectations(t)
	})
}

func TestClient_CreateIterationMedia(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		exampleInput := &models.IterationMediaCreationInput{}
		c, mockDB := buildTestClient()
		expected := &models.IterationMedia{}

		mockDB.IterationMediaDataManager.On("CreateIterationMedia", mock.Anything, exampleInput).Return(expected, nil)

		actual, err := c.CreateIterationMedia(context.Background(), exampleInput)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mockDB.AssertExpectations(t)
	})
}

func TestClient_UpdateIterationMedia(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		exampleInput := &models.IterationMedia{}
		c, mockDB := buildTestClient()
		var expected error

		mockDB.IterationMediaDataManager.On("UpdateIterationMedia", mock.Anything, exampleInput).Return(expected)

		err := c.UpdateIterationMedia(context.Background(), exampleInput)
		assert.NoError(t, err)
	})
}

func TestClient_ArchiveIterationMedia(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		exampleUserID := uint64(123)
		exampleIterationMediaID := uint64(123)
		var expected error

		c, mockDB := buildTestClient()
		mockDB.IterationMediaDataManager.On("ArchiveIterationMedia", mock.Anything, exampleIterationMediaID, exampleUserID).Return(expected)

		err := c.ArchiveIterationMedia(context.Background(), exampleUserID, exampleIterationMediaID)
		assert.NoError(t, err)
	})
}
