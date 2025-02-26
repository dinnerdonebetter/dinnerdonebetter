package managers

import (
	"testing"

	msgconfig "github.com/dinnerdonebetter/backend/internal/lib/messagequeue/config"
	mockpublishers "github.com/dinnerdonebetter/backend/internal/lib/messagequeue/mock"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/metrics"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/tracing"
	textsearchcfg "github.com/dinnerdonebetter/backend/internal/lib/search/text/config"
	"github.com/dinnerdonebetter/backend/internal/lib/testutils"
	"github.com/dinnerdonebetter/backend/internal/services/eating/database"
	"github.com/dinnerdonebetter/backend/internal/services/eating/events"
	"github.com/dinnerdonebetter/backend/internal/services/eating/types"
	"github.com/dinnerdonebetter/backend/internal/services/eating/types/fakes"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func buildValidEnumerationsManagerForTest(t *testing.T) *validEnumerationManager {
	t.Helper()

	queueCfg := &msgconfig.QueuesConfig{
		DataChangesTopicName: t.Name(),
	}

	mpp := &mockpublishers.PublisherProvider{}
	mpp.On("ProvidePublisher", queueCfg.DataChangesTopicName).Return(&mockpublishers.Publisher{}, nil)

	m, err := NewValidEnumerationsManager(
		t.Context(),
		logging.NewNoopLogger(),
		tracing.NewNoopTracerProvider(),
		database.NewMockDatabase(),
		queueCfg,
		mpp,
		&textsearchcfg.Config{},
		metrics.NewNoopMetricsProvider(),
	)
	require.NoError(t, err)

	mock.AssertExpectationsForObjects(t, mpp)

	return m.(*validEnumerationManager)
}

func TestValidEnumerationManager_SearchValidIngredientGroups(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		t.SkipNow()
	})
}

func TestValidEnumerationManager_ListValidIngredientGroups(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		vem := buildValidEnumerationsManagerForTest(t)

		expected := fakes.BuildFakeValidIngredientGroupsList()

		expectations := setupExpectationsForValidEnumerationManager(
			vem,
			func(db *database.MockDatabase) {
				db.ValidIngredientGroupDataManagerMock.On(testutils.GetMethodName(vem.db.GetValidIngredientGroups), testutils.ContextMatcher, testutils.QueryFilterMatcher).Return(expected, nil)
			},
		)

		actual, cursor, err := vem.ListValidIngredientGroups(ctx, nil)
		assert.NoError(t, err)
		assert.Empty(t, cursor)
		assert.Equal(t, expected.Data, actual)

		mock.AssertExpectationsForObjects(t, expectations...)
	})
}

func TestValidEnumerationManager_CreateValidIngredientGroup(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		vem := buildValidEnumerationsManagerForTest(t)

		expected := fakes.BuildFakeValidIngredientGroup()
		fakeInput := fakes.BuildFakeValidIngredientGroupCreationRequestInput()

		expectations := setupExpectationsForValidEnumerationManager(
			vem,
			func(db *database.MockDatabase) {
				db.ValidIngredientGroupDataManagerMock.On(testutils.GetMethodName(vem.db.CreateValidIngredientGroup), testutils.ContextMatcher, testutils.MatchType[*types.ValidIngredientGroupDatabaseCreationInput]()).Return(expected, nil)
			},
			map[string][]string{
				events.ValidIngredientGroupCreated: {keys.ValidIngredientGroupIDKey},
			},
		)

		actual, err := vem.CreateValidIngredientGroup(ctx, fakeInput)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mock.AssertExpectationsForObjects(t, expectations...)
	})
}

func TestValidEnumerationManager_ReadValidIngredientGroup(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		vem := buildValidEnumerationsManagerForTest(t)

		expected := fakes.BuildFakeValidIngredientGroup()

		expectations := setupExpectationsForValidEnumerationManager(
			vem,
			func(db *database.MockDatabase) {
				db.ValidIngredientGroupDataManagerMock.On(testutils.GetMethodName(vem.db.GetValidIngredientGroup), testutils.ContextMatcher, expected.ID).Return(expected, nil)
			},
		)

		actual, err := vem.ReadValidIngredientGroup(ctx, expected.ID)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mock.AssertExpectationsForObjects(t, expectations...)
	})
}

func TestValidEnumerationManager_UpdateValidIngredientGroup(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		mpm := buildValidEnumerationsManagerForTest(t)

		exampleValidIngredientGroup := fakes.BuildFakeValidIngredientGroup()
		exampleInput := fakes.BuildFakeValidIngredientGroupUpdateRequestInput()

		expectations := setupExpectationsForValidEnumerationManager(
			mpm,
			func(db *database.MockDatabase) {
				db.ValidIngredientGroupDataManagerMock.On(testutils.GetMethodName(mpm.db.GetValidIngredientGroup), testutils.ContextMatcher, exampleValidIngredientGroup.ID).Return(exampleValidIngredientGroup, nil)
				db.ValidIngredientGroupDataManagerMock.On(testutils.GetMethodName(mpm.db.UpdateValidIngredientGroup), testutils.ContextMatcher, testutils.MatchType[*types.ValidIngredientGroup]()).Return(nil)
			},
			map[string][]string{
				events.ValidIngredientGroupUpdated: {keys.ValidIngredientGroupIDKey},
			},
		)

		assert.NoError(t, mpm.UpdateValidIngredientGroup(ctx, exampleValidIngredientGroup.ID, exampleInput))

		mock.AssertExpectationsForObjects(t, expectations...)
	})
}

func TestValidEnumerationManager_ArchiveValidIngredientGroup(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		vem := buildValidEnumerationsManagerForTest(t)

		expected := fakes.BuildFakeValidIngredientGroup()

		expectations := setupExpectationsForValidEnumerationManager(
			vem,
			func(db *database.MockDatabase) {
				db.ValidIngredientGroupDataManagerMock.On(testutils.GetMethodName(vem.db.ArchiveValidIngredientGroup), testutils.ContextMatcher, expected.ID).Return(nil)
			},
			map[string][]string{
				events.ValidIngredientGroupArchived: {keys.ValidIngredientGroupIDKey},
			},
		)

		err := vem.ArchiveValidIngredientGroup(ctx, expected.ID)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, expectations...)
	})
}

func TestValidEnumerationManager_ListValidIngredientMeasurementUnits(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		vem := buildValidEnumerationsManagerForTest(t)

		expected := fakes.BuildFakeValidIngredientMeasurementUnitsList()

		expectations := setupExpectationsForValidEnumerationManager(
			vem,
			func(db *database.MockDatabase) {
				db.ValidIngredientMeasurementUnitDataManagerMock.On(testutils.GetMethodName(vem.db.GetValidIngredientMeasurementUnits), testutils.ContextMatcher, testutils.QueryFilterMatcher).Return(expected, nil)
			},
		)

		actual, cursor, err := vem.ListValidIngredientMeasurementUnits(ctx, nil)
		assert.NoError(t, err)
		assert.Empty(t, cursor)
		assert.Equal(t, expected.Data, actual)

		mock.AssertExpectationsForObjects(t, expectations...)
	})
}

func TestValidEnumerationManager_CreateValidIngredientMeasurementUnit(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		vem := buildValidEnumerationsManagerForTest(t)

		expected := fakes.BuildFakeValidIngredientMeasurementUnit()
		fakeInput := fakes.BuildFakeValidIngredientMeasurementUnitCreationRequestInput()

		expectations := setupExpectationsForValidEnumerationManager(
			vem,
			func(db *database.MockDatabase) {
				db.ValidIngredientMeasurementUnitDataManagerMock.On(testutils.GetMethodName(vem.db.CreateValidIngredientMeasurementUnit), testutils.ContextMatcher, testutils.MatchType[*types.ValidIngredientMeasurementUnitDatabaseCreationInput]()).Return(expected, nil)
			},
			map[string][]string{
				events.ValidIngredientMeasurementUnitCreated: {keys.ValidIngredientMeasurementUnitIDKey},
			},
		)

		actual, err := vem.CreateValidIngredientMeasurementUnit(ctx, fakeInput)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mock.AssertExpectationsForObjects(t, expectations...)
	})
}

func TestValidEnumerationManager_ReadValidIngredientMeasurementUnit(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		vem := buildValidEnumerationsManagerForTest(t)

		expected := fakes.BuildFakeValidIngredientMeasurementUnit()

		expectations := setupExpectationsForValidEnumerationManager(
			vem,
			func(db *database.MockDatabase) {
				db.ValidIngredientMeasurementUnitDataManagerMock.On(testutils.GetMethodName(vem.db.GetValidIngredientMeasurementUnit), testutils.ContextMatcher, expected.ID).Return(expected, nil)
			},
		)

		actual, err := vem.ReadValidIngredientMeasurementUnit(ctx, expected.ID)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mock.AssertExpectationsForObjects(t, expectations...)
	})
}

func TestValidEnumerationManager_UpdateValidIngredientMeasurementUnit(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		mpm := buildValidEnumerationsManagerForTest(t)

		exampleValidIngredientMeasurementUnit := fakes.BuildFakeValidIngredientMeasurementUnit()
		exampleInput := fakes.BuildFakeValidIngredientMeasurementUnitUpdateRequestInput()

		expectations := setupExpectationsForValidEnumerationManager(
			mpm,
			func(db *database.MockDatabase) {
				db.ValidIngredientMeasurementUnitDataManagerMock.On(testutils.GetMethodName(mpm.db.GetValidIngredientMeasurementUnit), testutils.ContextMatcher, exampleValidIngredientMeasurementUnit.ID).Return(exampleValidIngredientMeasurementUnit, nil)
				db.ValidIngredientMeasurementUnitDataManagerMock.On(testutils.GetMethodName(mpm.db.UpdateValidIngredientMeasurementUnit), testutils.ContextMatcher, testutils.MatchType[*types.ValidIngredientMeasurementUnit]()).Return(nil)
			},
			map[string][]string{
				events.ValidIngredientMeasurementUnitUpdated: {keys.ValidIngredientMeasurementUnitIDKey},
			},
		)

		assert.NoError(t, mpm.UpdateValidIngredientMeasurementUnit(ctx, exampleValidIngredientMeasurementUnit.ID, exampleInput))

		mock.AssertExpectationsForObjects(t, expectations...)
	})
}

func TestValidEnumerationManager_ArchiveValidIngredientMeasurementUnit(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		vem := buildValidEnumerationsManagerForTest(t)

		expected := fakes.BuildFakeValidIngredientMeasurementUnit()

		expectations := setupExpectationsForValidEnumerationManager(
			vem,
			func(db *database.MockDatabase) {
				db.ValidIngredientMeasurementUnitDataManagerMock.On(testutils.GetMethodName(vem.db.ArchiveValidIngredientMeasurementUnit), testutils.ContextMatcher, expected.ID).Return(nil)
			},
			map[string][]string{
				events.ValidIngredientMeasurementUnitArchived: {keys.ValidIngredientMeasurementUnitIDKey},
			},
		)

		err := vem.ArchiveValidIngredientMeasurementUnit(ctx, expected.ID)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, expectations...)
	})
}

func TestValidEnumerationManager_SearchValidIngredientMeasurementUnitsByIngredient(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		t.SkipNow()
	})
}

func TestValidEnumerationManager_SearchValidIngredientMeasurementUnitsByMeasurementUnit(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		t.SkipNow()
	})
}

func TestValidEnumerationManager_ListValidIngredientPreparations(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		vem := buildValidEnumerationsManagerForTest(t)

		expected := fakes.BuildFakeValidIngredientPreparationsList()

		expectations := setupExpectationsForValidEnumerationManager(
			vem,
			func(db *database.MockDatabase) {
				db.ValidIngredientPreparationDataManagerMock.On(testutils.GetMethodName(vem.db.GetValidIngredientPreparations), testutils.ContextMatcher, testutils.QueryFilterMatcher).Return(expected, nil)
			},
		)

		actual, cursor, err := vem.ListValidIngredientPreparations(ctx, nil)
		assert.NoError(t, err)
		assert.Empty(t, cursor)
		assert.Equal(t, expected.Data, actual)

		mock.AssertExpectationsForObjects(t, expectations...)
	})
}

func TestValidEnumerationManager_CreateValidIngredientPreparation(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		vem := buildValidEnumerationsManagerForTest(t)

		expected := fakes.BuildFakeValidIngredientPreparation()
		fakeInput := fakes.BuildFakeValidIngredientPreparationCreationRequestInput()

		expectations := setupExpectationsForValidEnumerationManager(
			vem,
			func(db *database.MockDatabase) {
				db.ValidIngredientPreparationDataManagerMock.On(testutils.GetMethodName(vem.db.CreateValidIngredientPreparation), testutils.ContextMatcher, testutils.MatchType[*types.ValidIngredientPreparationDatabaseCreationInput]()).Return(expected, nil)
			},
			map[string][]string{
				events.ValidIngredientPreparationCreated: {keys.ValidIngredientPreparationIDKey},
			},
		)

		actual, err := vem.CreateValidIngredientPreparation(ctx, fakeInput)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mock.AssertExpectationsForObjects(t, expectations...)
	})
}

func TestValidEnumerationManager_ReadValidIngredientPreparation(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		vem := buildValidEnumerationsManagerForTest(t)

		expected := fakes.BuildFakeValidIngredientPreparation()

		expectations := setupExpectationsForValidEnumerationManager(
			vem,
			func(db *database.MockDatabase) {
				db.ValidIngredientPreparationDataManagerMock.On(testutils.GetMethodName(vem.db.GetValidIngredientPreparation), testutils.ContextMatcher, expected.ID).Return(expected, nil)
			},
		)

		actual, err := vem.ReadValidIngredientPreparation(ctx, expected.ID)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mock.AssertExpectationsForObjects(t, expectations...)
	})
}

func TestValidEnumerationManager_UpdateValidIngredientPreparation(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		mpm := buildValidEnumerationsManagerForTest(t)

		exampleValidIngredientPreparation := fakes.BuildFakeValidIngredientPreparation()
		exampleInput := fakes.BuildFakeValidIngredientPreparationUpdateRequestInput()

		expectations := setupExpectationsForValidEnumerationManager(
			mpm,
			func(db *database.MockDatabase) {
				db.ValidIngredientPreparationDataManagerMock.On(testutils.GetMethodName(mpm.db.GetValidIngredientPreparation), testutils.ContextMatcher, exampleValidIngredientPreparation.ID).Return(exampleValidIngredientPreparation, nil)
				db.ValidIngredientPreparationDataManagerMock.On(testutils.GetMethodName(mpm.db.UpdateValidIngredientPreparation), testutils.ContextMatcher, testutils.MatchType[*types.ValidIngredientPreparation]()).Return(nil)
			},
			map[string][]string{
				events.ValidIngredientPreparationUpdated: {keys.ValidIngredientPreparationIDKey},
			},
		)

		assert.NoError(t, mpm.UpdateValidIngredientPreparation(ctx, exampleValidIngredientPreparation.ID, exampleInput))

		mock.AssertExpectationsForObjects(t, expectations...)
	})
}

func TestValidEnumerationManager_ArchiveValidIngredientPreparation(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		vem := buildValidEnumerationsManagerForTest(t)

		expected := fakes.BuildFakeValidIngredientPreparation()

		expectations := setupExpectationsForValidEnumerationManager(
			vem,
			func(db *database.MockDatabase) {
				db.ValidIngredientPreparationDataManagerMock.On(testutils.GetMethodName(vem.db.ArchiveValidIngredientPreparation), testutils.ContextMatcher, expected.ID).Return(nil)
			},
			map[string][]string{
				events.ValidIngredientPreparationArchived: {keys.ValidIngredientPreparationIDKey},
			},
		)

		err := vem.ArchiveValidIngredientPreparation(ctx, expected.ID)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, expectations...)
	})
}

func TestValidEnumerationManager_SearchValidIngredientPreparationsByIngredient(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		t.SkipNow()
	})
}

func TestValidEnumerationManager_SearchValidIngredientPreparationsByPreparation(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		t.SkipNow()
	})
}

func TestValidEnumerationManager_SearchValidIngredients(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		t.SkipNow()
	})
}

func TestValidEnumerationManager_ListValidIngredients(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		vem := buildValidEnumerationsManagerForTest(t)

		expected := fakes.BuildFakeValidIngredientsList()

		expectations := setupExpectationsForValidEnumerationManager(
			vem,
			func(db *database.MockDatabase) {
				db.ValidIngredientDataManagerMock.On(testutils.GetMethodName(vem.db.GetValidIngredients), testutils.ContextMatcher, testutils.QueryFilterMatcher).Return(expected, nil)
			},
		)

		actual, cursor, err := vem.ListValidIngredients(ctx, nil)
		assert.NoError(t, err)
		assert.Empty(t, cursor)
		assert.Equal(t, expected.Data, actual)

		mock.AssertExpectationsForObjects(t, expectations...)
	})
}

func TestValidEnumerationManager_CreateValidIngredient(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		vem := buildValidEnumerationsManagerForTest(t)

		expected := fakes.BuildFakeValidIngredient()
		fakeInput := fakes.BuildFakeValidIngredientCreationRequestInput()

		expectations := setupExpectationsForValidEnumerationManager(
			vem,
			func(db *database.MockDatabase) {
				db.ValidIngredientDataManagerMock.On(testutils.GetMethodName(vem.db.CreateValidIngredient), testutils.ContextMatcher, testutils.MatchType[*types.ValidIngredientDatabaseCreationInput]()).Return(expected, nil)
			},
			map[string][]string{
				events.ValidIngredientCreated: {keys.ValidIngredientIDKey},
			},
		)

		actual, err := vem.CreateValidIngredient(ctx, fakeInput)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mock.AssertExpectationsForObjects(t, expectations...)
	})
}

func TestValidEnumerationManager_ReadValidIngredient(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		vem := buildValidEnumerationsManagerForTest(t)

		expected := fakes.BuildFakeValidIngredient()

		expectations := setupExpectationsForValidEnumerationManager(
			vem,
			func(db *database.MockDatabase) {
				db.ValidIngredientDataManagerMock.On(testutils.GetMethodName(vem.db.GetValidIngredient), testutils.ContextMatcher, expected.ID).Return(expected, nil)
			},
		)

		actual, err := vem.ReadValidIngredient(ctx, expected.ID)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mock.AssertExpectationsForObjects(t, expectations...)
	})
}

func TestValidEnumerationManager_RandomValidIngredient(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		t.SkipNow()
	})
}

func TestValidEnumerationManager_UpdateValidIngredient(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		mpm := buildValidEnumerationsManagerForTest(t)

		exampleValidIngredient := fakes.BuildFakeValidIngredient()
		exampleInput := fakes.BuildFakeValidIngredientUpdateRequestInput()

		expectations := setupExpectationsForValidEnumerationManager(
			mpm,
			func(db *database.MockDatabase) {
				db.ValidIngredientDataManagerMock.On(testutils.GetMethodName(mpm.db.GetValidIngredient), testutils.ContextMatcher, exampleValidIngredient.ID).Return(exampleValidIngredient, nil)
				db.ValidIngredientDataManagerMock.On(testutils.GetMethodName(mpm.db.UpdateValidIngredient), testutils.ContextMatcher, testutils.MatchType[*types.ValidIngredient]()).Return(nil)
			},
			map[string][]string{
				events.ValidIngredientUpdated: {keys.ValidIngredientIDKey},
			},
		)

		assert.NoError(t, mpm.UpdateValidIngredient(ctx, exampleValidIngredient.ID, exampleInput))

		mock.AssertExpectationsForObjects(t, expectations...)
	})
}

func TestValidEnumerationManager_ArchiveValidIngredient(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		vem := buildValidEnumerationsManagerForTest(t)

		expected := fakes.BuildFakeValidIngredient()

		expectations := setupExpectationsForValidEnumerationManager(
			vem,
			func(db *database.MockDatabase) {
				db.ValidIngredientDataManagerMock.On(testutils.GetMethodName(vem.db.ArchiveValidIngredient), testutils.ContextMatcher, expected.ID).Return(nil)
			},
			map[string][]string{
				events.ValidIngredientArchived: {keys.ValidIngredientIDKey},
			},
		)

		err := vem.ArchiveValidIngredient(ctx, expected.ID)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, expectations...)
	})
}

func TestValidEnumerationManager_SearchValidIngredientsByPreparationAndIngredientName(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		t.SkipNow()
	})
}

func TestValidEnumerationManager_ListValidIngredientStateIngredients(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		vem := buildValidEnumerationsManagerForTest(t)

		expected := fakes.BuildFakeValidIngredientStateIngredientsList()

		expectations := setupExpectationsForValidEnumerationManager(
			vem,
			func(db *database.MockDatabase) {
				db.ValidIngredientStateIngredientDataManagerMock.On(testutils.GetMethodName(vem.db.GetValidIngredientStateIngredients), testutils.ContextMatcher, testutils.QueryFilterMatcher).Return(expected, nil)
			},
		)

		actual, cursor, err := vem.ListValidIngredientStateIngredients(ctx, nil)
		assert.NoError(t, err)
		assert.Empty(t, cursor)
		assert.Equal(t, expected.Data, actual)

		mock.AssertExpectationsForObjects(t, expectations...)
	})
}

func TestValidEnumerationManager_CreateValidIngredientStateIngredient(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		vem := buildValidEnumerationsManagerForTest(t)

		expected := fakes.BuildFakeValidIngredientStateIngredient()
		fakeInput := fakes.BuildFakeValidIngredientStateIngredientCreationRequestInput()

		expectations := setupExpectationsForValidEnumerationManager(
			vem,
			func(db *database.MockDatabase) {
				db.ValidIngredientStateIngredientDataManagerMock.On(testutils.GetMethodName(vem.db.CreateValidIngredientStateIngredient), testutils.ContextMatcher, testutils.MatchType[*types.ValidIngredientStateIngredientDatabaseCreationInput]()).Return(expected, nil)
			},
			map[string][]string{
				events.ValidIngredientStateIngredientCreated: {keys.ValidIngredientStateIngredientIDKey},
			},
		)

		actual, err := vem.CreateValidIngredientStateIngredient(ctx, fakeInput)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mock.AssertExpectationsForObjects(t, expectations...)
	})
}

func TestValidEnumerationManager_ReadValidIngredientStateIngredient(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		vem := buildValidEnumerationsManagerForTest(t)

		expected := fakes.BuildFakeValidIngredientStateIngredient()

		expectations := setupExpectationsForValidEnumerationManager(
			vem,
			func(db *database.MockDatabase) {
				db.ValidIngredientStateIngredientDataManagerMock.On(testutils.GetMethodName(vem.db.GetValidIngredientStateIngredient), testutils.ContextMatcher, expected.ID).Return(expected, nil)
			},
		)

		actual, err := vem.ReadValidIngredientStateIngredient(ctx, expected.ID)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mock.AssertExpectationsForObjects(t, expectations...)
	})
}

func TestValidEnumerationManager_UpdateValidIngredientStateIngredient(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		mpm := buildValidEnumerationsManagerForTest(t)

		exampleValidIngredientStateIngredient := fakes.BuildFakeValidIngredientStateIngredient()
		exampleInput := fakes.BuildFakeValidIngredientStateIngredientUpdateRequestInput()

		expectations := setupExpectationsForValidEnumerationManager(
			mpm,
			func(db *database.MockDatabase) {
				db.ValidIngredientStateIngredientDataManagerMock.On(testutils.GetMethodName(mpm.db.GetValidIngredientStateIngredient), testutils.ContextMatcher, exampleValidIngredientStateIngredient.ID).Return(exampleValidIngredientStateIngredient, nil)
				db.ValidIngredientStateIngredientDataManagerMock.On(testutils.GetMethodName(mpm.db.UpdateValidIngredientStateIngredient), testutils.ContextMatcher, testutils.MatchType[*types.ValidIngredientStateIngredient]()).Return(nil)
			},
			map[string][]string{
				events.ValidIngredientStateIngredientUpdated: {keys.ValidIngredientStateIngredientIDKey},
			},
		)

		assert.NoError(t, mpm.UpdateValidIngredientStateIngredient(ctx, exampleValidIngredientStateIngredient.ID, exampleInput))

		mock.AssertExpectationsForObjects(t, expectations...)
	})
}

func TestValidEnumerationManager_ArchiveValidIngredientStateIngredient(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		vem := buildValidEnumerationsManagerForTest(t)

		expected := fakes.BuildFakeValidIngredientStateIngredient()

		expectations := setupExpectationsForValidEnumerationManager(
			vem,
			func(db *database.MockDatabase) {
				db.ValidIngredientStateIngredientDataManagerMock.On(testutils.GetMethodName(vem.db.ArchiveValidIngredientStateIngredient), testutils.ContextMatcher, expected.ID).Return(nil)
			},
			map[string][]string{
				events.ValidIngredientStateIngredientArchived: {keys.ValidIngredientStateIngredientIDKey},
			},
		)

		err := vem.ArchiveValidIngredientStateIngredient(ctx, expected.ID)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, expectations...)
	})
}

func TestValidEnumerationManager_SearchValidIngredientStateIngredientsByIngredient(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		t.SkipNow()
	})
}

func TestValidEnumerationManager_SearchValidIngredientStateIngredientsByIngredientState(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		t.SkipNow()
	})
}

func TestValidEnumerationManager_SearchValidIngredientStates(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		t.SkipNow()
	})
}

func TestValidEnumerationManager_ListValidIngredientStates(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		vem := buildValidEnumerationsManagerForTest(t)

		expected := fakes.BuildFakeValidIngredientStatesList()

		expectations := setupExpectationsForValidEnumerationManager(
			vem,
			func(db *database.MockDatabase) {
				db.ValidIngredientStateDataManagerMock.On(testutils.GetMethodName(vem.db.GetValidIngredientStates), testutils.ContextMatcher, testutils.QueryFilterMatcher).Return(expected, nil)
			},
		)

		actual, cursor, err := vem.ListValidIngredientStates(ctx, nil)
		assert.NoError(t, err)
		assert.Empty(t, cursor)
		assert.Equal(t, expected.Data, actual)

		mock.AssertExpectationsForObjects(t, expectations...)
	})
}

func TestValidEnumerationManager_CreateValidIngredientState(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		vem := buildValidEnumerationsManagerForTest(t)

		expected := fakes.BuildFakeValidIngredientState()
		fakeInput := fakes.BuildFakeValidIngredientStateCreationRequestInput()

		expectations := setupExpectationsForValidEnumerationManager(
			vem,
			func(db *database.MockDatabase) {
				db.ValidIngredientStateDataManagerMock.On(testutils.GetMethodName(vem.db.CreateValidIngredientState), testutils.ContextMatcher, testutils.MatchType[*types.ValidIngredientStateDatabaseCreationInput]()).Return(expected, nil)
			},
			map[string][]string{
				events.ValidIngredientStateCreated: {keys.ValidIngredientStateIDKey},
			},
		)

		actual, err := vem.CreateValidIngredientState(ctx, fakeInput)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mock.AssertExpectationsForObjects(t, expectations...)
	})
}

func TestValidEnumerationManager_ReadValidIngredientState(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		vem := buildValidEnumerationsManagerForTest(t)

		expected := fakes.BuildFakeValidIngredientState()

		expectations := setupExpectationsForValidEnumerationManager(
			vem,
			func(db *database.MockDatabase) {
				db.ValidIngredientStateDataManagerMock.On(testutils.GetMethodName(vem.db.GetValidIngredientState), testutils.ContextMatcher, expected.ID).Return(expected, nil)
			},
		)

		actual, err := vem.ReadValidIngredientState(ctx, expected.ID)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mock.AssertExpectationsForObjects(t, expectations...)
	})
}

func TestValidEnumerationManager_UpdateValidIngredientState(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		mpm := buildValidEnumerationsManagerForTest(t)

		exampleValidIngredientState := fakes.BuildFakeValidIngredientState()
		exampleInput := fakes.BuildFakeValidIngredientStateUpdateRequestInput()

		expectations := setupExpectationsForValidEnumerationManager(
			mpm,
			func(db *database.MockDatabase) {
				db.ValidIngredientStateDataManagerMock.On(testutils.GetMethodName(mpm.db.GetValidIngredientState), testutils.ContextMatcher, exampleValidIngredientState.ID).Return(exampleValidIngredientState, nil)
				db.ValidIngredientStateDataManagerMock.On(testutils.GetMethodName(mpm.db.UpdateValidIngredientState), testutils.ContextMatcher, testutils.MatchType[*types.ValidIngredientState]()).Return(nil)
			},
			map[string][]string{
				events.ValidIngredientStateUpdated: {keys.ValidIngredientStateIDKey},
			},
		)

		assert.NoError(t, mpm.UpdateValidIngredientState(ctx, exampleValidIngredientState.ID, exampleInput))

		mock.AssertExpectationsForObjects(t, expectations...)
	})
}

func TestValidEnumerationManager_ArchiveValidIngredientState(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		vem := buildValidEnumerationsManagerForTest(t)

		expected := fakes.BuildFakeValidIngredientState()

		expectations := setupExpectationsForValidEnumerationManager(
			vem,
			func(db *database.MockDatabase) {
				db.ValidIngredientStateDataManagerMock.On(testutils.GetMethodName(vem.db.ArchiveValidIngredientState), testutils.ContextMatcher, expected.ID).Return(nil)
			},
			map[string][]string{
				events.ValidIngredientStateArchived: {keys.ValidIngredientStateIDKey},
			},
		)

		err := vem.ArchiveValidIngredientState(ctx, expected.ID)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, expectations...)
	})
}

func TestValidEnumerationManager_SearchValidMeasurementUnits(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		t.SkipNow()
	})
}

func TestValidEnumerationManager_SearchValidMeasurementUnitsByIngredientID(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		t.SkipNow()
	})
}

func TestValidEnumerationManager_ListValidMeasurementUnits(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		vem := buildValidEnumerationsManagerForTest(t)

		expected := fakes.BuildFakeValidMeasurementUnitsList()

		expectations := setupExpectationsForValidEnumerationManager(
			vem,
			func(db *database.MockDatabase) {
				db.ValidMeasurementUnitDataManagerMock.On(testutils.GetMethodName(vem.db.GetValidMeasurementUnits), testutils.ContextMatcher, testutils.QueryFilterMatcher).Return(expected, nil)
			},
		)

		actual, cursor, err := vem.ListValidMeasurementUnits(ctx, nil)
		assert.NoError(t, err)
		assert.Empty(t, cursor)
		assert.Equal(t, expected.Data, actual)

		mock.AssertExpectationsForObjects(t, expectations...)
	})
}

func TestValidEnumerationManager_CreateValidMeasurementUnit(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		vem := buildValidEnumerationsManagerForTest(t)

		expected := fakes.BuildFakeValidMeasurementUnit()
		fakeInput := fakes.BuildFakeValidMeasurementUnitCreationRequestInput()

		expectations := setupExpectationsForValidEnumerationManager(
			vem,
			func(db *database.MockDatabase) {
				db.ValidMeasurementUnitDataManagerMock.On(testutils.GetMethodName(vem.db.CreateValidMeasurementUnit), testutils.ContextMatcher, testutils.MatchType[*types.ValidMeasurementUnitDatabaseCreationInput]()).Return(expected, nil)
			},
			map[string][]string{
				events.ValidMeasurementUnitCreated: {keys.ValidMeasurementUnitIDKey},
			},
		)

		actual, err := vem.CreateValidMeasurementUnit(ctx, fakeInput)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mock.AssertExpectationsForObjects(t, expectations...)
	})
}

func TestValidEnumerationManager_ReadValidMeasurementUnit(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		vem := buildValidEnumerationsManagerForTest(t)

		expected := fakes.BuildFakeValidMeasurementUnit()

		expectations := setupExpectationsForValidEnumerationManager(
			vem,
			func(db *database.MockDatabase) {
				db.ValidMeasurementUnitDataManagerMock.On(testutils.GetMethodName(vem.db.GetValidMeasurementUnit), testutils.ContextMatcher, expected.ID).Return(expected, nil)
			},
		)

		actual, err := vem.ReadValidMeasurementUnit(ctx, expected.ID)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mock.AssertExpectationsForObjects(t, expectations...)
	})
}

func TestValidEnumerationManager_UpdateValidMeasurementUnit(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		mpm := buildValidEnumerationsManagerForTest(t)

		exampleValidMeasurementUnit := fakes.BuildFakeValidMeasurementUnit()
		exampleInput := fakes.BuildFakeValidMeasurementUnitUpdateRequestInput()

		expectations := setupExpectationsForValidEnumerationManager(
			mpm,
			func(db *database.MockDatabase) {
				db.ValidMeasurementUnitDataManagerMock.On(testutils.GetMethodName(mpm.db.GetValidMeasurementUnit), testutils.ContextMatcher, exampleValidMeasurementUnit.ID).Return(exampleValidMeasurementUnit, nil)
				db.ValidMeasurementUnitDataManagerMock.On(testutils.GetMethodName(mpm.db.UpdateValidMeasurementUnit), testutils.ContextMatcher, testutils.MatchType[*types.ValidMeasurementUnit]()).Return(nil)
			},
			map[string][]string{
				events.ValidMeasurementUnitUpdated: {keys.ValidMeasurementUnitIDKey},
			},
		)

		assert.NoError(t, mpm.UpdateValidMeasurementUnit(ctx, exampleValidMeasurementUnit.ID, exampleInput))

		mock.AssertExpectationsForObjects(t, expectations...)
	})
}

func TestValidEnumerationManager_ArchiveValidMeasurementUnit(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		vem := buildValidEnumerationsManagerForTest(t)

		expected := fakes.BuildFakeValidMeasurementUnit()

		expectations := setupExpectationsForValidEnumerationManager(
			vem,
			func(db *database.MockDatabase) {
				db.ValidMeasurementUnitDataManagerMock.On(testutils.GetMethodName(vem.db.ArchiveValidMeasurementUnit), testutils.ContextMatcher, expected.ID).Return(nil)
			},
			map[string][]string{
				events.ValidMeasurementUnitArchived: {keys.ValidMeasurementUnitIDKey},
			},
		)

		err := vem.ArchiveValidMeasurementUnit(ctx, expected.ID)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, expectations...)
	})
}

func TestValidEnumerationManager_SearchValidInstruments(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		t.SkipNow()
	})
}

func TestValidEnumerationManager_ListValidInstruments(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		vem := buildValidEnumerationsManagerForTest(t)

		expected := fakes.BuildFakeValidInstrumentsList()

		expectations := setupExpectationsForValidEnumerationManager(
			vem,
			func(db *database.MockDatabase) {
				db.ValidInstrumentDataManagerMock.On(testutils.GetMethodName(vem.db.GetValidInstruments), testutils.ContextMatcher, testutils.QueryFilterMatcher).Return(expected, nil)
			},
		)

		actual, cursor, err := vem.ListValidInstruments(ctx, nil)
		assert.NoError(t, err)
		assert.Empty(t, cursor)
		assert.Equal(t, expected.Data, actual)

		mock.AssertExpectationsForObjects(t, expectations...)
	})
}

func TestValidEnumerationManager_CreateValidInstrument(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		vem := buildValidEnumerationsManagerForTest(t)

		expected := fakes.BuildFakeValidInstrument()
		fakeInput := fakes.BuildFakeValidInstrumentCreationRequestInput()

		expectations := setupExpectationsForValidEnumerationManager(
			vem,
			func(db *database.MockDatabase) {
				db.ValidInstrumentDataManagerMock.On(testutils.GetMethodName(vem.db.CreateValidInstrument), testutils.ContextMatcher, testutils.MatchType[*types.ValidInstrumentDatabaseCreationInput]()).Return(expected, nil)
			},
			map[string][]string{
				events.ValidInstrumentCreated: {keys.ValidInstrumentIDKey},
			},
		)

		actual, err := vem.CreateValidInstrument(ctx, fakeInput)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mock.AssertExpectationsForObjects(t, expectations...)
	})
}

func TestValidEnumerationManager_ReadValidInstrument(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		vem := buildValidEnumerationsManagerForTest(t)

		expected := fakes.BuildFakeValidInstrument()

		expectations := setupExpectationsForValidEnumerationManager(
			vem,
			func(db *database.MockDatabase) {
				db.ValidInstrumentDataManagerMock.On(testutils.GetMethodName(vem.db.GetValidInstrument), testutils.ContextMatcher, expected.ID).Return(expected, nil)
			},
		)

		actual, err := vem.ReadValidInstrument(ctx, expected.ID)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mock.AssertExpectationsForObjects(t, expectations...)
	})
}

func TestValidEnumerationManager_RandomValidInstrument(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		t.SkipNow()
	})
}

func TestValidEnumerationManager_UpdateValidInstrument(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		mpm := buildValidEnumerationsManagerForTest(t)

		exampleValidInstrument := fakes.BuildFakeValidInstrument()
		exampleInput := fakes.BuildFakeValidInstrumentUpdateRequestInput()

		expectations := setupExpectationsForValidEnumerationManager(
			mpm,
			func(db *database.MockDatabase) {
				db.ValidInstrumentDataManagerMock.On(testutils.GetMethodName(mpm.db.GetValidInstrument), testutils.ContextMatcher, exampleValidInstrument.ID).Return(exampleValidInstrument, nil)
				db.ValidInstrumentDataManagerMock.On(testutils.GetMethodName(mpm.db.UpdateValidInstrument), testutils.ContextMatcher, testutils.MatchType[*types.ValidInstrument]()).Return(nil)
			},
			map[string][]string{
				events.ValidInstrumentUpdated: {keys.ValidInstrumentIDKey},
			},
		)

		assert.NoError(t, mpm.UpdateValidInstrument(ctx, exampleValidInstrument.ID, exampleInput))

		mock.AssertExpectationsForObjects(t, expectations...)
	})
}

func TestValidEnumerationManager_ArchiveValidInstrument(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		vem := buildValidEnumerationsManagerForTest(t)

		expected := fakes.BuildFakeValidInstrument()

		expectations := setupExpectationsForValidEnumerationManager(
			vem,
			func(db *database.MockDatabase) {
				db.ValidInstrumentDataManagerMock.On(testutils.GetMethodName(vem.db.ArchiveValidInstrument), testutils.ContextMatcher, expected.ID).Return(nil)
			},
			map[string][]string{
				events.ValidInstrumentArchived: {keys.ValidInstrumentIDKey},
			},
		)

		err := vem.ArchiveValidInstrument(ctx, expected.ID)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, expectations...)
	})
}

func TestValidEnumerationManager_ValidMeasurementUnitConversionsFromMeasurementUnit(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		t.SkipNow()
	})
}

func TestValidEnumerationManager_ValidMeasurementUnitConversionsToMeasurementUnit(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		t.SkipNow()
	})
}

func TestValidEnumerationManager_CreateValidMeasurementUnitConversion(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		vem := buildValidEnumerationsManagerForTest(t)

		expected := fakes.BuildFakeValidMeasurementUnitConversion()
		fakeInput := fakes.BuildFakeValidMeasurementUnitConversionCreationRequestInput()

		expectations := setupExpectationsForValidEnumerationManager(
			vem,
			func(db *database.MockDatabase) {
				db.ValidMeasurementUnitConversionDataManagerMock.On(testutils.GetMethodName(vem.db.CreateValidMeasurementUnitConversion), testutils.ContextMatcher, testutils.MatchType[*types.ValidMeasurementUnitConversionDatabaseCreationInput]()).Return(expected, nil)
			},
			map[string][]string{
				events.ValidMeasurementUnitConversionCreated: {keys.ValidMeasurementUnitConversionIDKey},
			},
		)

		actual, err := vem.CreateValidMeasurementUnitConversion(ctx, fakeInput)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mock.AssertExpectationsForObjects(t, expectations...)
	})
}

func TestValidEnumerationManager_ReadValidMeasurementUnitConversion(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		vem := buildValidEnumerationsManagerForTest(t)

		expected := fakes.BuildFakeValidMeasurementUnitConversion()

		expectations := setupExpectationsForValidEnumerationManager(
			vem,
			func(db *database.MockDatabase) {
				db.ValidMeasurementUnitConversionDataManagerMock.On(testutils.GetMethodName(vem.db.GetValidMeasurementUnitConversion), testutils.ContextMatcher, expected.ID).Return(expected, nil)
			},
		)

		actual, err := vem.ReadValidMeasurementUnitConversion(ctx, expected.ID)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mock.AssertExpectationsForObjects(t, expectations...)
	})
}

func TestValidEnumerationManager_UpdateValidMeasurementUnitConversion(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		mpm := buildValidEnumerationsManagerForTest(t)

		exampleValidMeasurementUnitConversion := fakes.BuildFakeValidMeasurementUnitConversion()
		exampleInput := fakes.BuildFakeValidMeasurementUnitConversionUpdateRequestInput()

		expectations := setupExpectationsForValidEnumerationManager(
			mpm,
			func(db *database.MockDatabase) {
				db.ValidMeasurementUnitConversionDataManagerMock.On(testutils.GetMethodName(mpm.db.GetValidMeasurementUnitConversion), testutils.ContextMatcher, exampleValidMeasurementUnitConversion.ID).Return(exampleValidMeasurementUnitConversion, nil)
				db.ValidMeasurementUnitConversionDataManagerMock.On(testutils.GetMethodName(mpm.db.UpdateValidMeasurementUnitConversion), testutils.ContextMatcher, testutils.MatchType[*types.ValidMeasurementUnitConversion]()).Return(nil)
			},
			map[string][]string{
				events.ValidMeasurementUnitConversionUpdated: {keys.ValidMeasurementUnitConversionIDKey},
			},
		)

		assert.NoError(t, mpm.UpdateValidMeasurementUnitConversion(ctx, exampleValidMeasurementUnitConversion.ID, exampleInput))

		mock.AssertExpectationsForObjects(t, expectations...)
	})
}

func TestValidEnumerationManager_ArchiveValidMeasurementUnitConversion(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		vem := buildValidEnumerationsManagerForTest(t)

		expected := fakes.BuildFakeValidMeasurementUnitConversion()

		expectations := setupExpectationsForValidEnumerationManager(
			vem,
			func(db *database.MockDatabase) {
				db.ValidMeasurementUnitConversionDataManagerMock.On(testutils.GetMethodName(vem.db.ArchiveValidMeasurementUnitConversion), testutils.ContextMatcher, expected.ID).Return(nil)
			},
			map[string][]string{
				events.ValidMeasurementUnitConversionArchived: {keys.ValidMeasurementUnitConversionIDKey},
			},
		)

		err := vem.ArchiveValidMeasurementUnitConversion(ctx, expected.ID)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, expectations...)
	})
}

func TestValidEnumerationManager_ListValidPreparationInstruments(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		vem := buildValidEnumerationsManagerForTest(t)

		expected := fakes.BuildFakeValidPreparationInstrumentsList()

		expectations := setupExpectationsForValidEnumerationManager(
			vem,
			func(db *database.MockDatabase) {
				db.ValidPreparationInstrumentDataManagerMock.On(testutils.GetMethodName(vem.db.GetValidPreparationInstruments), testutils.ContextMatcher, testutils.QueryFilterMatcher).Return(expected, nil)
			},
		)

		actual, cursor, err := vem.ListValidPreparationInstruments(ctx, nil)
		assert.NoError(t, err)
		assert.Empty(t, cursor)
		assert.Equal(t, expected.Data, actual)

		mock.AssertExpectationsForObjects(t, expectations...)
	})
}

func TestValidEnumerationManager_CreateValidPreparationInstrument(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		vem := buildValidEnumerationsManagerForTest(t)

		expected := fakes.BuildFakeValidPreparationInstrument()
		fakeInput := fakes.BuildFakeValidPreparationInstrumentCreationRequestInput()

		expectations := setupExpectationsForValidEnumerationManager(
			vem,
			func(db *database.MockDatabase) {
				db.ValidPreparationInstrumentDataManagerMock.On(testutils.GetMethodName(vem.db.CreateValidPreparationInstrument), testutils.ContextMatcher, testutils.MatchType[*types.ValidPreparationInstrumentDatabaseCreationInput]()).Return(expected, nil)
			},
			map[string][]string{
				events.ValidPreparationInstrumentCreated: {keys.ValidPreparationInstrumentIDKey},
			},
		)

		actual, err := vem.CreateValidPreparationInstrument(ctx, fakeInput)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mock.AssertExpectationsForObjects(t, expectations...)
	})
}

func TestValidEnumerationManager_ReadValidPreparationInstrument(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		vem := buildValidEnumerationsManagerForTest(t)

		expected := fakes.BuildFakeValidPreparationInstrument()

		expectations := setupExpectationsForValidEnumerationManager(
			vem,
			func(db *database.MockDatabase) {
				db.ValidPreparationInstrumentDataManagerMock.On(testutils.GetMethodName(vem.db.GetValidPreparationInstrument), testutils.ContextMatcher, expected.ID).Return(expected, nil)
			},
		)

		actual, err := vem.ReadValidPreparationInstrument(ctx, expected.ID)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mock.AssertExpectationsForObjects(t, expectations...)
	})
}

func TestValidEnumerationManager_UpdateValidPreparationInstrument(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		mpm := buildValidEnumerationsManagerForTest(t)

		exampleValidPreparationInstrument := fakes.BuildFakeValidPreparationInstrument()
		exampleInput := fakes.BuildFakeValidPreparationInstrumentUpdateRequestInput()

		expectations := setupExpectationsForValidEnumerationManager(
			mpm,
			func(db *database.MockDatabase) {
				db.ValidPreparationInstrumentDataManagerMock.On(testutils.GetMethodName(mpm.db.GetValidPreparationInstrument), testutils.ContextMatcher, exampleValidPreparationInstrument.ID).Return(exampleValidPreparationInstrument, nil)
				db.ValidPreparationInstrumentDataManagerMock.On(testutils.GetMethodName(mpm.db.UpdateValidPreparationInstrument), testutils.ContextMatcher, testutils.MatchType[*types.ValidPreparationInstrument]()).Return(nil)
			},
			map[string][]string{
				events.ValidPreparationInstrumentUpdated: {keys.ValidPreparationInstrumentIDKey},
			},
		)

		assert.NoError(t, mpm.UpdateValidPreparationInstrument(ctx, exampleValidPreparationInstrument.ID, exampleInput))

		mock.AssertExpectationsForObjects(t, expectations...)
	})
}

func TestValidEnumerationManager_ArchiveValidPreparationInstrument(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		vem := buildValidEnumerationsManagerForTest(t)

		expected := fakes.BuildFakeValidPreparationInstrument()

		expectations := setupExpectationsForValidEnumerationManager(
			vem,
			func(db *database.MockDatabase) {
				db.ValidPreparationInstrumentDataManagerMock.On(testutils.GetMethodName(vem.db.ArchiveValidPreparationInstrument), testutils.ContextMatcher, expected.ID).Return(nil)
			},
			map[string][]string{
				events.ValidPreparationInstrumentArchived: {keys.ValidPreparationInstrumentIDKey},
			},
		)

		err := vem.ArchiveValidPreparationInstrument(ctx, expected.ID)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, expectations...)
	})
}

func TestValidEnumerationManager_SearchValidPreparationInstrumentsByPreparation(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		t.SkipNow()
	})
}

func TestValidEnumerationManager_SearchValidPreparationInstrumentsByInstrument(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		t.SkipNow()
	})
}

func TestValidEnumerationManager_SearchValidPreparations(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		t.SkipNow()
	})
}

func TestValidEnumerationManager_ListValidPreparations(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		vem := buildValidEnumerationsManagerForTest(t)

		expected := fakes.BuildFakeValidPreparationsList()

		expectations := setupExpectationsForValidEnumerationManager(
			vem,
			func(db *database.MockDatabase) {
				db.ValidPreparationDataManagerMock.On(testutils.GetMethodName(vem.db.GetValidPreparations), testutils.ContextMatcher, testutils.QueryFilterMatcher).Return(expected, nil)
			},
		)

		actual, cursor, err := vem.ListValidPreparations(ctx, nil)
		assert.NoError(t, err)
		assert.Empty(t, cursor)
		assert.Equal(t, expected.Data, actual)

		mock.AssertExpectationsForObjects(t, expectations...)
	})
}

func TestValidEnumerationManager_CreateValidPreparation(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		vem := buildValidEnumerationsManagerForTest(t)

		expected := fakes.BuildFakeValidPreparation()
		fakeInput := fakes.BuildFakeValidPreparationCreationRequestInput()

		expectations := setupExpectationsForValidEnumerationManager(
			vem,
			func(db *database.MockDatabase) {
				db.ValidPreparationDataManagerMock.On(testutils.GetMethodName(vem.db.CreateValidPreparation), testutils.ContextMatcher, testutils.MatchType[*types.ValidPreparationDatabaseCreationInput]()).Return(expected, nil)
			},
			map[string][]string{
				events.ValidPreparationCreated: {keys.ValidPreparationIDKey},
			},
		)

		actual, err := vem.CreateValidPreparation(ctx, fakeInput)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mock.AssertExpectationsForObjects(t, expectations...)
	})
}

func TestValidEnumerationManager_ReadValidPreparation(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		vem := buildValidEnumerationsManagerForTest(t)

		expected := fakes.BuildFakeValidPreparation()

		expectations := setupExpectationsForValidEnumerationManager(
			vem,
			func(db *database.MockDatabase) {
				db.ValidPreparationDataManagerMock.On(testutils.GetMethodName(vem.db.GetValidPreparation), testutils.ContextMatcher, expected.ID).Return(expected, nil)
			},
		)

		actual, err := vem.ReadValidPreparation(ctx, expected.ID)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mock.AssertExpectationsForObjects(t, expectations...)
	})
}

func TestValidEnumerationManager_RandomValidPreparation(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		t.SkipNow()
	})
}

func TestValidEnumerationManager_UpdateValidPreparation(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		mpm := buildValidEnumerationsManagerForTest(t)

		exampleValidPreparation := fakes.BuildFakeValidPreparation()
		exampleInput := fakes.BuildFakeValidPreparationUpdateRequestInput()

		expectations := setupExpectationsForValidEnumerationManager(
			mpm,
			func(db *database.MockDatabase) {
				db.ValidPreparationDataManagerMock.On(testutils.GetMethodName(mpm.db.GetValidPreparation), testutils.ContextMatcher, exampleValidPreparation.ID).Return(exampleValidPreparation, nil)
				db.ValidPreparationDataManagerMock.On(testutils.GetMethodName(mpm.db.UpdateValidPreparation), testutils.ContextMatcher, testutils.MatchType[*types.ValidPreparation]()).Return(nil)
			},
			map[string][]string{
				events.ValidPreparationUpdated: {keys.ValidPreparationIDKey},
			},
		)

		assert.NoError(t, mpm.UpdateValidPreparation(ctx, exampleValidPreparation.ID, exampleInput))

		mock.AssertExpectationsForObjects(t, expectations...)
	})
}

func TestValidEnumerationManager_ArchiveValidPreparation(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		vem := buildValidEnumerationsManagerForTest(t)

		expected := fakes.BuildFakeValidPreparation()

		expectations := setupExpectationsForValidEnumerationManager(
			vem,
			func(db *database.MockDatabase) {
				db.ValidPreparationDataManagerMock.On(testutils.GetMethodName(vem.db.ArchiveValidPreparation), testutils.ContextMatcher, expected.ID).Return(nil)
			},
			map[string][]string{
				events.ValidPreparationArchived: {keys.ValidPreparationIDKey},
			},
		)

		err := vem.ArchiveValidPreparation(ctx, expected.ID)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, expectations...)
	})
}

func TestValidEnumerationManager_ListValidPreparationVessels(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		vem := buildValidEnumerationsManagerForTest(t)

		expected := fakes.BuildFakeValidPreparationVesselsList()

		expectations := setupExpectationsForValidEnumerationManager(
			vem,
			func(db *database.MockDatabase) {
				db.ValidPreparationVesselDataManagerMock.On(testutils.GetMethodName(vem.db.GetValidPreparationVessels), testutils.ContextMatcher, testutils.QueryFilterMatcher).Return(expected, nil)
			},
		)

		actual, cursor, err := vem.ListValidPreparationVessels(ctx, nil)
		assert.NoError(t, err)
		assert.Empty(t, cursor)
		assert.Equal(t, expected.Data, actual)

		mock.AssertExpectationsForObjects(t, expectations...)
	})
}

func TestValidEnumerationManager_CreateValidPreparationVessel(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		vem := buildValidEnumerationsManagerForTest(t)

		expected := fakes.BuildFakeValidPreparationVessel()
		fakeInput := fakes.BuildFakeValidPreparationVesselCreationRequestInput()

		expectations := setupExpectationsForValidEnumerationManager(
			vem,
			func(db *database.MockDatabase) {
				db.ValidPreparationVesselDataManagerMock.On(testutils.GetMethodName(vem.db.CreateValidPreparationVessel), testutils.ContextMatcher, testutils.MatchType[*types.ValidPreparationVesselDatabaseCreationInput]()).Return(expected, nil)
			},
			map[string][]string{
				events.ValidPreparationVesselCreated: {keys.ValidPreparationVesselIDKey},
			},
		)

		actual, err := vem.CreateValidPreparationVessel(ctx, fakeInput)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mock.AssertExpectationsForObjects(t, expectations...)
	})
}

func TestValidEnumerationManager_ReadValidPreparationVessel(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		vem := buildValidEnumerationsManagerForTest(t)

		expected := fakes.BuildFakeValidPreparationVessel()

		expectations := setupExpectationsForValidEnumerationManager(
			vem,
			func(db *database.MockDatabase) {
				db.ValidPreparationVesselDataManagerMock.On(testutils.GetMethodName(vem.db.GetValidPreparationVessel), testutils.ContextMatcher, expected.ID).Return(expected, nil)
			},
		)

		actual, err := vem.ReadValidPreparationVessel(ctx, expected.ID)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mock.AssertExpectationsForObjects(t, expectations...)
	})
}

func TestValidEnumerationManager_UpdateValidPreparationVessel(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		mpm := buildValidEnumerationsManagerForTest(t)

		exampleValidPreparationVessel := fakes.BuildFakeValidPreparationVessel()
		exampleInput := fakes.BuildFakeValidPreparationVesselUpdateRequestInput()

		expectations := setupExpectationsForValidEnumerationManager(
			mpm,
			func(db *database.MockDatabase) {
				db.ValidPreparationVesselDataManagerMock.On(testutils.GetMethodName(mpm.db.GetValidPreparationVessel), testutils.ContextMatcher, exampleValidPreparationVessel.ID).Return(exampleValidPreparationVessel, nil)
				db.ValidPreparationVesselDataManagerMock.On(testutils.GetMethodName(mpm.db.UpdateValidPreparationVessel), testutils.ContextMatcher, testutils.MatchType[*types.ValidPreparationVessel]()).Return(nil)
			},
			map[string][]string{
				events.ValidPreparationVesselUpdated: {keys.ValidPreparationVesselIDKey},
			},
		)

		assert.NoError(t, mpm.UpdateValidPreparationVessel(ctx, exampleValidPreparationVessel.ID, exampleInput))

		mock.AssertExpectationsForObjects(t, expectations...)
	})
}

func TestValidEnumerationManager_ArchiveValidPreparationVessel(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		vem := buildValidEnumerationsManagerForTest(t)

		expected := fakes.BuildFakeValidPreparationVessel()

		expectations := setupExpectationsForValidEnumerationManager(
			vem,
			func(db *database.MockDatabase) {
				db.ValidPreparationVesselDataManagerMock.On(testutils.GetMethodName(vem.db.ArchiveValidPreparationVessel), testutils.ContextMatcher, expected.ID).Return(nil)
			},
			map[string][]string{
				events.ValidPreparationVesselArchived: {keys.ValidPreparationVesselIDKey},
			},
		)

		err := vem.ArchiveValidPreparationVessel(ctx, expected.ID)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, expectations...)
	})
}

func TestValidEnumerationManager_SearchValidPreparationVesselsByPreparation(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		t.SkipNow()
	})
}

func TestValidEnumerationManager_SearchValidPreparationVesselsByVessel(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		t.SkipNow()
	})
}

func TestValidEnumerationManager_SearchValidVessels(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		t.SkipNow()
	})
}

func TestValidEnumerationManager_ListValidVessels(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		vem := buildValidEnumerationsManagerForTest(t)

		expected := fakes.BuildFakeValidVesselsList()

		expectations := setupExpectationsForValidEnumerationManager(
			vem,
			func(db *database.MockDatabase) {
				db.ValidVesselDataManagerMock.On(testutils.GetMethodName(vem.db.GetValidVessels), testutils.ContextMatcher, testutils.QueryFilterMatcher).Return(expected, nil)
			},
		)

		actual, cursor, err := vem.ListValidVessels(ctx, nil)
		assert.NoError(t, err)
		assert.Empty(t, cursor)
		assert.Equal(t, expected.Data, actual)

		mock.AssertExpectationsForObjects(t, expectations...)
	})
}

func TestValidEnumerationManager_CreateValidVessel(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		vem := buildValidEnumerationsManagerForTest(t)

		expected := fakes.BuildFakeValidVessel()
		fakeInput := fakes.BuildFakeValidVesselCreationRequestInput()

		expectations := setupExpectationsForValidEnumerationManager(
			vem,
			func(db *database.MockDatabase) {
				db.ValidVesselDataManagerMock.On(testutils.GetMethodName(vem.db.CreateValidVessel), testutils.ContextMatcher, testutils.MatchType[*types.ValidVesselDatabaseCreationInput]()).Return(expected, nil)
			},
			map[string][]string{
				events.ValidVesselCreated: {keys.ValidVesselIDKey},
			},
		)

		actual, err := vem.CreateValidVessel(ctx, fakeInput)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mock.AssertExpectationsForObjects(t, expectations...)
	})
}

func TestValidEnumerationManager_ReadValidVessel(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		vem := buildValidEnumerationsManagerForTest(t)

		expected := fakes.BuildFakeValidVessel()

		expectations := setupExpectationsForValidEnumerationManager(
			vem,
			func(db *database.MockDatabase) {
				db.ValidVesselDataManagerMock.On(testutils.GetMethodName(vem.db.GetValidVessel), testutils.ContextMatcher, expected.ID).Return(expected, nil)
			},
		)

		actual, err := vem.ReadValidVessel(ctx, expected.ID)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mock.AssertExpectationsForObjects(t, expectations...)
	})
}

func TestValidEnumerationManager_RandomValidVessel(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		t.SkipNow()
	})
}

func TestValidEnumerationManager_UpdateValidVessel(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		mpm := buildValidEnumerationsManagerForTest(t)

		exampleValidVessel := fakes.BuildFakeValidVessel()
		exampleInput := fakes.BuildFakeValidVesselUpdateRequestInput()

		expectations := setupExpectationsForValidEnumerationManager(
			mpm,
			func(db *database.MockDatabase) {
				db.ValidVesselDataManagerMock.On(testutils.GetMethodName(mpm.db.GetValidVessel), testutils.ContextMatcher, exampleValidVessel.ID).Return(exampleValidVessel, nil)
				db.ValidVesselDataManagerMock.On(testutils.GetMethodName(mpm.db.UpdateValidVessel), testutils.ContextMatcher, testutils.MatchType[*types.ValidVessel]()).Return(nil)
			},
			map[string][]string{
				events.ValidVesselUpdated: {keys.ValidVesselIDKey},
			},
		)

		assert.NoError(t, mpm.UpdateValidVessel(ctx, exampleValidVessel.ID, exampleInput))

		mock.AssertExpectationsForObjects(t, expectations...)
	})
}

func TestValidEnumerationManager_ArchiveValidVessel(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		vem := buildValidEnumerationsManagerForTest(t)

		expected := fakes.BuildFakeValidVessel()

		expectations := setupExpectationsForValidEnumerationManager(
			vem,
			func(db *database.MockDatabase) {
				db.ValidVesselDataManagerMock.On(testutils.GetMethodName(vem.db.ArchiveValidVessel), testutils.ContextMatcher, expected.ID).Return(nil)
			},
			map[string][]string{
				events.ValidVesselArchived: {keys.ValidVesselIDKey},
			},
		)

		err := vem.ArchiveValidVessel(ctx, expected.ID)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, expectations...)
	})
}
