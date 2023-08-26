package postgres

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestQuerier_ArchiveOAuth2ClientTokenByAccess(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, db := buildTestClient(t)

		assert.Error(t, c.ArchiveOAuth2ClientTokenByAccess(ctx, ""))

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_ArchiveOAuth2ClientTokenByCode(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, db := buildTestClient(t)

		assert.Error(t, c.ArchiveOAuth2ClientTokenByCode(ctx, ""))

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_ArchiveOAuth2ClientTokenByRefresh(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, db := buildTestClient(t)

		assert.Error(t, c.ArchiveOAuth2ClientTokenByRefresh(ctx, ""))

		mock.AssertExpectationsForObjects(t, db)
	})
}
