//go:build integration
// +build integration

package postgres

import (
	"context"
	"testing"

	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
)

func createValidInstrumentForTest(t *testing.T, ctx context.Context, exampleValidInstrument *types.ValidInstrument, dbc *Querier) *types.ValidInstrument {
	t.Helper()

	// create
	if exampleValidInstrument == nil {
		exampleValidInstrument = fakes.BuildFakeValidInstrument()
	}
	dbInput := converters.ConvertValidInstrumentToValidInstrumentDatabaseCreationInput(exampleValidInstrument)

	created, err := dbc.CreateValidInstrument(ctx, dbInput)
	exampleValidInstrument.CreatedAt = created.CreatedAt
	assert.NoError(t, err)
	assert.Equal(t, exampleValidInstrument, created)

	validInstrument, err := dbc.GetValidInstrument(ctx, created.ID)
	exampleValidInstrument.CreatedAt = validInstrument.CreatedAt

	assert.NoError(t, err)
	assert.Equal(t, validInstrument, exampleValidInstrument)

	return created
}

func TestQuerier_Integration_GetValidInstrument(T *testing.T) {
	T.Parallel()

	//nolint:paralleltest // this test uses a test container
	T.Run("integration", func(t *testing.T) {
		ctx := context.Background()

		c, container := buildDatabaseClientForTest(t, ctx)
		defer func(t *testing.T) {
			t.Helper()
			assert.NoError(t, container.Terminate(ctx))
		}(t)

		assert.NotNil(t, createValidInstrumentForTest(t, ctx, nil, c))
	})
}
