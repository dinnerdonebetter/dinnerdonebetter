package postgres

import (
	"context"
	"testing"

	"github.com/dinnerdonebetter/backend/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
)

func TestQuerier_HouseholdInstrumentOwnershipExists(T *testing.T) {
	T.Parallel()

	T.Run("with invalid household instrument ownership ID", func(t *testing.T) {
		t.Parallel()

		exampleHouseholdID := fakes.BuildFakeID()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.HouseholdInstrumentOwnershipExists(ctx, "", exampleHouseholdID)
		assert.Error(t, err)
		assert.False(t, actual)
	})
}

func TestQuerier_GetHouseholdInstrumentOwnership(T *testing.T) {
	T.Parallel()

	T.Run("with invalid household instrument ownership ID", func(t *testing.T) {
		t.Parallel()

		exampleHouseholdID := fakes.BuildFakeID()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetHouseholdInstrumentOwnership(ctx, "", exampleHouseholdID)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestQuerier_CreateHouseholdInstrumentOwnership(T *testing.T) {
	T.Parallel()

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.CreateHouseholdInstrumentOwnership(ctx, nil)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestQuerier_UpdateHouseholdInstrumentOwnership(T *testing.T) {
	T.Parallel()

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.UpdateHouseholdInstrumentOwnership(ctx, nil))
	})
}

func TestQuerier_ArchiveHouseholdInstrumentOwnership(T *testing.T) {
	T.Parallel()

	T.Run("with invalid household instrument ownership ID", func(t *testing.T) {
		t.Parallel()

		exampleHouseholdID := fakes.BuildFakeID()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.ArchiveHouseholdInstrumentOwnership(ctx, "", exampleHouseholdID))
	})
}
