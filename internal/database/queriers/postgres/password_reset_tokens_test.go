package postgres

import (
	"context"
	"testing"
)

func TestSQLQuerier_scanPasswordResetToken(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		q, _ := buildTestClient(t)

		_, _ = ctx, q
	})
}

func TestSQLQuerier_GetPasswordResetToken(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		q, _ := buildTestClient(t)

		_, _ = ctx, q
	})
}

func TestSQLQuerier_GetTotalPasswordResetTokenCount(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		q, _ := buildTestClient(t)

		_, _ = ctx, q
	})
}

func TestSQLQuerier_CreatePasswordResetToken(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		q, _ := buildTestClient(t)

		_, _ = ctx, q
	})
}

func TestSQLQuerier_ArchivePasswordResetToken(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		q, _ := buildTestClient(t)

		_, _ = ctx, q
	})
}
