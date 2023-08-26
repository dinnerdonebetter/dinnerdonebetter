package postgres

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSQLQuerier_GetPasswordResetTokenByToken(T *testing.T) {
	T.Parallel()

	T.Run("with missing token", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetPasswordResetTokenByToken(ctx, "")
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestSQLQuerier_CreatePasswordResetToken(T *testing.T) {
	T.Parallel()

	T.Run("with missing input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.CreatePasswordResetToken(ctx, nil)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestSQLQuerier_RedeemPasswordResetToken(T *testing.T) {
	T.Parallel()

	T.Run("with missing ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, db := buildTestClient(t)

		actual := c.RedeemPasswordResetToken(ctx, "")
		assert.Error(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}
