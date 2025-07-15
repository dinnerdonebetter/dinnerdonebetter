package grpc

import (
	"testing"

	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning"
	mealplanningfakes "github.com/dinnerdonebetter/backend/internal/domain/mealplanning/fakes"
	mockmanagers "github.com/dinnerdonebetter/backend/internal/domain/mealplanning/managers/mock"
	mealplanninggrpc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/platform/fake"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/platform/testutils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func buildServiceImplForTest(t *testing.T) *ServiceImpl {
	t.Helper()

	return &ServiceImpl{
		tracer: tracing.NewTracerForTest(t.Name()),
		logger: logging.NewNoopLogger(),
	}
}

func TestServiceImpl_ArchiveValidIngredient(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		exampleValidIngredientID := mealplanningfakes.BuildFakeID()

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("ArchiveValidIngredient", testutils.ContextMatcher, exampleValidIngredientID).Return(nil)
		s.validEnumerationsManager = mvem

		res, err := s.ArchiveValidIngredient(ctx, &mealplanninggrpc.ArchiveValidIngredientRequest{ValidIngredientID: exampleValidIngredientID})
		assert.NotNil(t, res)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mvem)
	})
}

func TestServiceImpl_ArchiveValidIngredientGroup(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		exampleValidIngredientGroupID := mealplanningfakes.BuildFakeID()

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("ArchiveValidIngredientGroup", testutils.ContextMatcher, exampleValidIngredientGroupID).Return(nil)
		s.validEnumerationsManager = mvem

		res, err := s.ArchiveValidIngredientGroup(ctx, &mealplanninggrpc.ArchiveValidIngredientGroupRequest{ValidIngredientGroupID: exampleValidIngredientGroupID})
		assert.NotNil(t, res)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mvem)
	})
}

func TestServiceImpl_ArchiveValidIngredientMeasurementUnit(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		exampleValidIngredientMeasurementUnitID := mealplanningfakes.BuildFakeID()

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("ArchiveValidIngredientMeasurementUnit", testutils.ContextMatcher, exampleValidIngredientMeasurementUnitID).Return(nil)
		s.validEnumerationsManager = mvem

		res, err := s.ArchiveValidIngredientMeasurementUnit(ctx, &mealplanninggrpc.ArchiveValidIngredientMeasurementUnitRequest{ValidIngredientMeasurementUnitID: exampleValidIngredientMeasurementUnitID})
		assert.NotNil(t, res)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mvem)
	})
}

func TestServiceImpl_ArchiveValidIngredientPreparation(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		exampleValidIngredientPreparationID := mealplanningfakes.BuildFakeID()

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("ArchiveValidIngredientPreparation", testutils.ContextMatcher, exampleValidIngredientPreparationID).Return(nil)
		s.validEnumerationsManager = mvem

		res, err := s.ArchiveValidIngredientPreparation(ctx, &mealplanninggrpc.ArchiveValidIngredientPreparationRequest{ValidIngredientPreparationID: exampleValidIngredientPreparationID})
		assert.NotNil(t, res)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mvem)
	})
}

func TestServiceImpl_ArchiveValidIngredientState(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		exampleValidIngredientStateID := mealplanningfakes.BuildFakeID()

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("ArchiveValidIngredientState", testutils.ContextMatcher, exampleValidIngredientStateID).Return(nil)
		s.validEnumerationsManager = mvem

		res, err := s.ArchiveValidIngredientState(ctx, &mealplanninggrpc.ArchiveValidIngredientStateRequest{ValidIngredientStateID: exampleValidIngredientStateID})
		assert.NotNil(t, res)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mvem)
	})
}

func TestServiceImpl_ArchiveValidIngredientStateIngredient(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		exampleValidIngredientStateIngredientID := mealplanningfakes.BuildFakeID()

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("ArchiveValidIngredientStateIngredient", testutils.ContextMatcher, exampleValidIngredientStateIngredientID).Return(nil)
		s.validEnumerationsManager = mvem

		res, err := s.ArchiveValidIngredientStateIngredient(ctx, &mealplanninggrpc.ArchiveValidIngredientStateIngredientRequest{ValidIngredientStateIngredientID: exampleValidIngredientStateIngredientID})
		assert.NotNil(t, res)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mvem)
	})
}

func TestServiceImpl_ArchiveValidInstrument(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		exampleValidInstrumentID := mealplanningfakes.BuildFakeID()

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("ArchiveValidInstrument", testutils.ContextMatcher, exampleValidInstrumentID).Return(nil)
		s.validEnumerationsManager = mvem

		res, err := s.ArchiveValidInstrument(ctx, &mealplanninggrpc.ArchiveValidInstrumentRequest{ValidInstrumentID: exampleValidInstrumentID})
		assert.NotNil(t, res)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mvem)
	})
}

func TestServiceImpl_ArchiveValidMeasurementUnit(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		exampleValidMeasurementUnitID := mealplanningfakes.BuildFakeID()

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("ArchiveValidMeasurementUnit", testutils.ContextMatcher, exampleValidMeasurementUnitID).Return(nil)
		s.validEnumerationsManager = mvem

		res, err := s.ArchiveValidMeasurementUnit(ctx, &mealplanninggrpc.ArchiveValidMeasurementUnitRequest{ValidMeasurementUnitID: exampleValidMeasurementUnitID})
		assert.NotNil(t, res)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mvem)
	})
}

func TestServiceImpl_ArchiveValidMeasurementUnitConversion(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		exampleValidMeasurementUnitConversionID := mealplanningfakes.BuildFakeID()

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("ArchiveValidMeasurementUnitConversion", testutils.ContextMatcher, exampleValidMeasurementUnitConversionID).Return(nil)
		s.validEnumerationsManager = mvem

		res, err := s.ArchiveValidMeasurementUnitConversion(ctx, &mealplanninggrpc.ArchiveValidMeasurementUnitConversionRequest{ValidMeasurementUnitConversionID: exampleValidMeasurementUnitConversionID})
		assert.NotNil(t, res)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mvem)
	})
}

func TestServiceImpl_ArchiveValidPreparation(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		exampleValidPreparationID := mealplanningfakes.BuildFakeID()

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("ArchiveValidPreparation", testutils.ContextMatcher, exampleValidPreparationID).Return(nil)
		s.validEnumerationsManager = mvem

		res, err := s.ArchiveValidPreparation(ctx, &mealplanninggrpc.ArchiveValidPreparationRequest{ValidPreparationID: exampleValidPreparationID})
		assert.NotNil(t, res)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mvem)
	})
}

func TestServiceImpl_ArchiveValidPreparationInstrument(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		exampleValidPreparationInstrumentID := mealplanningfakes.BuildFakeID()

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("ArchiveValidPreparationInstrument", testutils.ContextMatcher, exampleValidPreparationInstrumentID).Return(nil)
		s.validEnumerationsManager = mvem

		res, err := s.ArchiveValidPreparationInstrument(ctx, &mealplanninggrpc.ArchiveValidPreparationInstrumentRequest{ValidPreparationInstrumentID: exampleValidPreparationInstrumentID})
		assert.NotNil(t, res)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mvem)
	})
}

func TestServiceImpl_ArchiveValidPreparationVessel(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		exampleValidPreparationVesselID := mealplanningfakes.BuildFakeID()

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("ArchiveValidPreparationVessel", testutils.ContextMatcher, exampleValidPreparationVesselID).Return(nil)
		s.validEnumerationsManager = mvem

		res, err := s.ArchiveValidPreparationVessel(ctx, &mealplanninggrpc.ArchiveValidPreparationVesselRequest{ValidPreparationVesselID: exampleValidPreparationVesselID})
		assert.NotNil(t, res)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mvem)
	})
}

func TestServiceImpl_ArchiveValidVessel(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		exampleValidVesselID := mealplanningfakes.BuildFakeID()

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("ArchiveValidVessel", testutils.ContextMatcher, exampleValidVesselID).Return(nil)
		s.validEnumerationsManager = mvem

		res, err := s.ArchiveValidVessel(ctx, &mealplanninggrpc.ArchiveValidVesselRequest{ValidVesselID: exampleValidVesselID})
		assert.NotNil(t, res)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mvem)
	})
}

func TestServiceImpl_CreateValidIngredient(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		exampleValidIngredient := mealplanningfakes.BuildFakeValidIngredient()

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("CreateValidIngredient", testutils.ContextMatcher, testutils.MatchType[*mealplanning.ValidIngredientCreationRequestInput]()).Return(exampleValidIngredient, nil)
		s.validEnumerationsManager = mvem

		exampleInput := fake.BuildFakeForTest[mealplanninggrpc.CreateValidIngredientRequest](t)

		actual, err := s.CreateValidIngredient(ctx, exampleInput)
		assert.NotNil(t, actual)
		assert.NoError(t, err)
	})
}

func TestServiceImpl_CreateValidIngredientGroup(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		exampleValidIngredientGroup := mealplanningfakes.BuildFakeValidIngredientGroup()

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("CreateValidIngredientGroup", testutils.ContextMatcher, testutils.MatchType[*mealplanning.ValidIngredientGroupCreationRequestInput]()).Return(exampleValidIngredientGroup, nil)
		s.validEnumerationsManager = mvem

		exampleInput := fake.BuildFakeForTest[mealplanninggrpc.CreateValidIngredientGroupRequest](t)

		actual, err := s.CreateValidIngredientGroup(ctx, exampleInput)
		assert.NotNil(t, actual)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mvem)
	})
}

func TestServiceImpl_CreateValidIngredientMeasurementUnit(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		exampleValidIngredientMeasurementUnit := mealplanningfakes.BuildFakeValidIngredientMeasurementUnit()

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("CreateValidIngredientMeasurementUnit", testutils.ContextMatcher, testutils.MatchType[*mealplanning.ValidIngredientMeasurementUnitCreationRequestInput]()).Return(exampleValidIngredientMeasurementUnit, nil)
		s.validEnumerationsManager = mvem

		exampleInput := fake.BuildFakeForTest[mealplanninggrpc.CreateValidIngredientMeasurementUnitRequest](t)

		actual, err := s.CreateValidIngredientMeasurementUnit(ctx, exampleInput)
		assert.NotNil(t, actual)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mvem)
	})
}

func TestServiceImpl_CreateValidIngredientPreparation(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		exampleValidIngredientPreparation := mealplanningfakes.BuildFakeValidIngredientPreparation()

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("CreateValidIngredientPreparation", testutils.ContextMatcher, testutils.MatchType[*mealplanning.ValidIngredientPreparationCreationRequestInput]()).Return(exampleValidIngredientPreparation, nil)
		s.validEnumerationsManager = mvem

		exampleInput := fake.BuildFakeForTest[mealplanninggrpc.CreateValidIngredientPreparationRequest](t)

		actual, err := s.CreateValidIngredientPreparation(ctx, exampleInput)
		assert.NotNil(t, actual)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mvem)
	})
}

func TestServiceImpl_CreateValidIngredientState(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		exampleValidIngredientState := mealplanningfakes.BuildFakeValidIngredientState()

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("CreateValidIngredientState", testutils.ContextMatcher, testutils.MatchType[*mealplanning.ValidIngredientStateCreationRequestInput]()).Return(exampleValidIngredientState, nil)
		s.validEnumerationsManager = mvem

		exampleInput := fake.BuildFakeForTest[mealplanninggrpc.CreateValidIngredientStateRequest](t)

		actual, err := s.CreateValidIngredientState(ctx, exampleInput)
		assert.NotNil(t, actual)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mvem)
	})
}
func TestServiceImpl_CreateValidIngredientStateIngredient(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		exampleValidIngredientStateIngredient := mealplanningfakes.BuildFakeValidIngredientStateIngredient()

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("CreateValidIngredientStateIngredient", testutils.ContextMatcher, testutils.MatchType[*mealplanning.ValidIngredientStateIngredientCreationRequestInput]()).Return(exampleValidIngredientStateIngredient, nil)
		s.validEnumerationsManager = mvem

		exampleInput := fake.BuildFakeForTest[mealplanninggrpc.CreateValidIngredientStateIngredientRequest](t)

		actual, err := s.CreateValidIngredientStateIngredient(ctx, exampleInput)
		assert.NotNil(t, actual)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mvem)
	})
}

func TestServiceImpl_CreateValidInstrument(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		exampleValidInstrument := mealplanningfakes.BuildFakeValidInstrument()

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("CreateValidInstrument", testutils.ContextMatcher, testutils.MatchType[*mealplanning.ValidInstrumentCreationRequestInput]()).Return(exampleValidInstrument, nil)
		s.validEnumerationsManager = mvem

		exampleInput := fake.BuildFakeForTest[mealplanninggrpc.CreateValidInstrumentRequest](t)

		actual, err := s.CreateValidInstrument(ctx, exampleInput)
		assert.NotNil(t, actual)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mvem)
	})
}

func TestServiceImpl_CreateValidMeasurementUnit(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		exampleValidMeasurementUnit := mealplanningfakes.BuildFakeValidMeasurementUnit()

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("CreateValidMeasurementUnit", testutils.ContextMatcher, testutils.MatchType[*mealplanning.ValidMeasurementUnitCreationRequestInput]()).Return(exampleValidMeasurementUnit, nil)
		s.validEnumerationsManager = mvem

		exampleInput := fake.BuildFakeForTest[mealplanninggrpc.CreateValidMeasurementUnitRequest](t)

		actual, err := s.CreateValidMeasurementUnit(ctx, exampleInput)
		assert.NotNil(t, actual)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mvem)
	})
}

func TestServiceImpl_CreateValidMeasurementUnitConversion(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		exampleValidMeasurementUnitConversion := mealplanningfakes.BuildFakeValidMeasurementUnitConversion()

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("CreateValidMeasurementUnitConversion", testutils.ContextMatcher, testutils.MatchType[*mealplanning.ValidMeasurementUnitConversionCreationRequestInput]()).Return(exampleValidMeasurementUnitConversion, nil)
		s.validEnumerationsManager = mvem

		exampleInput := fake.BuildFakeForTest[mealplanninggrpc.CreateValidMeasurementUnitConversionRequest](t)

		actual, err := s.CreateValidMeasurementUnitConversion(ctx, exampleInput)
		assert.NotNil(t, actual)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mvem)
	})
}

func TestServiceImpl_CreateValidPreparation(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		exampleValidPreparation := mealplanningfakes.BuildFakeValidPreparation()

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("CreateValidPreparation", testutils.ContextMatcher, testutils.MatchType[*mealplanning.ValidPreparationCreationRequestInput]()).Return(exampleValidPreparation, nil)
		s.validEnumerationsManager = mvem

		exampleInput := fake.BuildFakeForTest[mealplanninggrpc.CreateValidPreparationRequest](t)

		actual, err := s.CreateValidPreparation(ctx, exampleInput)
		assert.NotNil(t, actual)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mvem)
	})
}

func TestServiceImpl_CreateValidPreparationInstrument(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		exampleValidPreparationInstrument := mealplanningfakes.BuildFakeValidPreparationInstrument()

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("CreateValidPreparationInstrument", testutils.ContextMatcher, testutils.MatchType[*mealplanning.ValidPreparationInstrumentCreationRequestInput]()).Return(exampleValidPreparationInstrument, nil)
		s.validEnumerationsManager = mvem

		exampleInput := fake.BuildFakeForTest[mealplanninggrpc.CreateValidPreparationInstrumentRequest](t)

		actual, err := s.CreateValidPreparationInstrument(ctx, exampleInput)
		assert.NotNil(t, actual)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mvem)
	})
}

func TestServiceImpl_CreateValidPreparationVessel(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		exampleValidPreparationVessel := mealplanningfakes.BuildFakeValidPreparationVessel()

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("CreateValidPreparationVessel", testutils.ContextMatcher, testutils.MatchType[*mealplanning.ValidPreparationVesselCreationRequestInput]()).Return(exampleValidPreparationVessel, nil)
		s.validEnumerationsManager = mvem

		exampleInput := fake.BuildFakeForTest[mealplanninggrpc.CreateValidPreparationVesselRequest](t)

		actual, err := s.CreateValidPreparationVessel(ctx, exampleInput)
		assert.NotNil(t, actual)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mvem)
	})
}

func TestServiceImpl_CreateValidVessel(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		exampleValidVessel := mealplanningfakes.BuildFakeValidVessel()

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("CreateValidVessel", testutils.ContextMatcher, testutils.MatchType[*mealplanning.ValidVesselCreationRequestInput]()).Return(exampleValidVessel, nil)
		s.validEnumerationsManager = mvem

		exampleInput := fake.BuildFakeForTest[mealplanninggrpc.CreateValidVesselRequest](t)

		actual, err := s.CreateValidVessel(ctx, exampleInput)
		assert.NotNil(t, actual)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mvem)
	})
}

func TestServiceImpl_GetRandomValidIngredient(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleResult := mealplanningfakes.BuildFakeValidIngredient()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("RandomValidIngredient", testutils.ContextMatcher).Return(exampleResult, nil)
		s.validEnumerationsManager = mvem

		result, err := s.GetRandomValidIngredient(ctx, &mealplanninggrpc.GetRandomValidIngredientRequest{})
		assert.Equal(t, exampleResult.ID, result.Result.ID)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mvem)
	})
}

func TestServiceImpl_GetRandomValidInstrument(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleResult := mealplanningfakes.BuildFakeValidInstrument()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("RandomValidInstrument", testutils.ContextMatcher).Return(exampleResult, nil)
		s.validEnumerationsManager = mvem

		result, err := s.GetRandomValidInstrument(ctx, &mealplanninggrpc.GetRandomValidInstrumentRequest{})
		assert.Equal(t, exampleResult.ID, result.Result.ID)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mvem)
	})
}

func TestServiceImpl_GetRandomValidPreparation(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleResult := mealplanningfakes.BuildFakeValidPreparation()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("RandomValidPreparation", testutils.ContextMatcher).Return(exampleResult, nil)
		s.validEnumerationsManager = mvem

		result, err := s.GetRandomValidPreparation(ctx, &mealplanninggrpc.GetRandomValidPreparationRequest{})
		assert.Equal(t, exampleResult.ID, result.Result.ID)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mvem)
	})
}

func TestServiceImpl_GetRandomValidVessel(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleResult := mealplanningfakes.BuildFakeValidVessel()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("RandomValidVessel", testutils.ContextMatcher).Return(exampleResult, nil)
		s.validEnumerationsManager = mvem

		result, err := s.GetRandomValidVessel(ctx, &mealplanninggrpc.GetRandomValidVesselRequest{})
		assert.Equal(t, exampleResult.ID, result.Result.ID)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mvem)
	})
}

func TestServiceImpl_GetValidIngredient(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleResult := mealplanningfakes.BuildFakeValidIngredient()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("ReadValidIngredient", testutils.ContextMatcher, exampleResult.ID).Return(exampleResult, nil)
		s.validEnumerationsManager = mvem

		result, err := s.GetValidIngredient(ctx, &mealplanninggrpc.GetValidIngredientRequest{ValidIngredientID: exampleResult.ID})
		assert.Equal(t, exampleResult.ID, result.Result.ID)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mvem)
	})
}

func TestServiceImpl_GetValidIngredientGroup(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleResult := mealplanningfakes.BuildFakeValidIngredientGroup()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("ReadValidIngredientGroup", testutils.ContextMatcher, exampleResult.ID).Return(exampleResult, nil)
		s.validEnumerationsManager = mvem

		result, err := s.GetValidIngredientGroup(ctx, &mealplanninggrpc.GetValidIngredientGroupRequest{ValidIngredientGroupID: exampleResult.ID})
		assert.Equal(t, exampleResult.ID, result.Result.ID)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mvem)
	})
}

func TestServiceImpl_GetValidIngredientGroups(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleResult := mealplanningfakes.BuildFakeValidIngredientGroupsList()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("ListValidIngredientGroups", testutils.ContextMatcher, testutils.QueryFilterMatcher).Return(exampleResult.Data, "", nil)
		s.validEnumerationsManager = mvem

		result, err := s.GetValidIngredientGroups(ctx, &mealplanninggrpc.GetValidIngredientGroupsRequest{})
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result.Results, len(exampleResult.Data))

		mock.AssertExpectationsForObjects(t, mvem)
	})
}

func TestServiceImpl_GetValidIngredientMeasurementUnit(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleResult := mealplanningfakes.BuildFakeValidIngredientMeasurementUnit()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("ReadValidIngredientMeasurementUnit", testutils.ContextMatcher, exampleResult.ID).Return(exampleResult, nil)
		s.validEnumerationsManager = mvem

		result, err := s.GetValidIngredientMeasurementUnit(ctx, &mealplanninggrpc.GetValidIngredientMeasurementUnitRequest{ValidIngredientMeasurementUnitID: exampleResult.ID})
		assert.Equal(t, exampleResult.ID, result.Result.ID)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mvem)
	})
}

func TestServiceImpl_GetValidIngredientMeasurementUnits(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleResult := mealplanningfakes.BuildFakeValidMeasurementUnitsList()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("ListValidMeasurementUnits", testutils.ContextMatcher, testutils.QueryFilterMatcher).Return(exampleResult.Data, "", nil)
		s.validEnumerationsManager = mvem

		result, err := s.GetValidMeasurementUnits(ctx, &mealplanninggrpc.GetValidMeasurementUnitsRequest{})
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result.Results, len(exampleResult.Data))

		mock.AssertExpectationsForObjects(t, mvem)
	})
}

func TestServiceImpl_GetValidIngredientMeasurementUnitsByIngredient(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleID := mealplanningfakes.BuildFakeID()
		exampleResult := mealplanningfakes.BuildFakeValidIngredientMeasurementUnitsList()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("SearchValidIngredientMeasurementUnitsByIngredient", testutils.ContextMatcher, exampleID, testutils.QueryFilterMatcher).Return(exampleResult.Data, nil)
		s.validEnumerationsManager = mvem

		result, err := s.GetValidIngredientMeasurementUnitsByIngredient(ctx, &mealplanninggrpc.GetValidIngredientMeasurementUnitsByIngredientRequest{
			ValidIngredientID: exampleID,
		})
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result.Results, len(exampleResult.Data))

		mock.AssertExpectationsForObjects(t, mvem)
	})
}

func TestServiceImpl_GetValidIngredientMeasurementUnitsByMeasurementUnit(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleID := mealplanningfakes.BuildFakeID()
		exampleResult := mealplanningfakes.BuildFakeValidIngredientMeasurementUnitsList()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("SearchValidIngredientMeasurementUnitsByMeasurementUnit", testutils.ContextMatcher, exampleID, testutils.QueryFilterMatcher).Return(exampleResult.Data, nil)
		s.validEnumerationsManager = mvem

		result, err := s.GetValidIngredientMeasurementUnitsByMeasurementUnit(ctx, &mealplanninggrpc.GetValidIngredientMeasurementUnitsByMeasurementUnitRequest{
			ValidMeasurementUnitID: exampleID,
		})
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result.Results, len(exampleResult.Data))

		mock.AssertExpectationsForObjects(t, mvem)
	})
}

func TestServiceImpl_GetValidIngredientPreparation(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleResult := mealplanningfakes.BuildFakeValidIngredientPreparation()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("ReadValidIngredientPreparation", testutils.ContextMatcher, exampleResult.ID).Return(exampleResult, nil)
		s.validEnumerationsManager = mvem

		result, err := s.GetValidIngredientPreparation(ctx, &mealplanninggrpc.GetValidIngredientPreparationRequest{ValidIngredientPreparationID: exampleResult.ID})
		assert.Equal(t, exampleResult.ID, result.Result.ID)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mvem)
	})
}

func TestServiceImpl_GetValidIngredientPreparations(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleResult := mealplanningfakes.BuildFakeValidIngredientPreparationsList()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("ListValidIngredientPreparations", testutils.ContextMatcher, testutils.QueryFilterMatcher).Return(exampleResult.Data, "", nil)
		s.validEnumerationsManager = mvem

		result, err := s.GetValidIngredientPreparations(ctx, &mealplanninggrpc.GetValidIngredientPreparationsRequest{})
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result.Results, len(exampleResult.Data))

		mock.AssertExpectationsForObjects(t, mvem)
	})
}

func TestServiceImpl_GetValidIngredientPreparationsByIngredient(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleID := mealplanningfakes.BuildFakeID()
		exampleResult := mealplanningfakes.BuildFakeValidIngredientPreparationsList()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("SearchValidIngredientPreparationsByIngredient", testutils.ContextMatcher, exampleID, testutils.QueryFilterMatcher).Return(exampleResult.Data, nil)
		s.validEnumerationsManager = mvem

		result, err := s.GetValidIngredientPreparationsByIngredient(ctx, &mealplanninggrpc.GetValidIngredientPreparationsByIngredientRequest{
			ValidIngredientID: exampleID,
		})
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result.Results, len(exampleResult.Data))

		mock.AssertExpectationsForObjects(t, mvem)
	})
}

func TestServiceImpl_GetValidIngredientPreparationsByPreparation(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleID := mealplanningfakes.BuildFakeID()
		exampleResult := mealplanningfakes.BuildFakeValidIngredientPreparationsList()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("SearchValidIngredientPreparationsByPreparation", testutils.ContextMatcher, exampleID, testutils.QueryFilterMatcher).Return(exampleResult.Data, nil)
		s.validEnumerationsManager = mvem

		result, err := s.GetValidIngredientPreparationsByPreparation(ctx, &mealplanninggrpc.GetValidIngredientPreparationsByPreparationRequest{
			ValidPreparationID: exampleID,
		})
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result.Results, len(exampleResult.Data))

		mock.AssertExpectationsForObjects(t, mvem)
	})
}

func TestServiceImpl_GetValidIngredientState(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleResult := mealplanningfakes.BuildFakeValidIngredientState()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("ReadValidIngredientState", testutils.ContextMatcher, exampleResult.ID).Return(exampleResult, nil)
		s.validEnumerationsManager = mvem

		result, err := s.GetValidIngredientState(ctx, &mealplanninggrpc.GetValidIngredientStateRequest{ValidIngredientStateID: exampleResult.ID})
		assert.Equal(t, exampleResult.ID, result.Result.ID)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mvem)
	})
}

func TestServiceImpl_GetValidIngredientStateIngredient(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleResult := mealplanningfakes.BuildFakeValidIngredientStateIngredient()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("ReadValidIngredientStateIngredient", testutils.ContextMatcher, exampleResult.ID).Return(exampleResult, nil)
		s.validEnumerationsManager = mvem

		result, err := s.GetValidIngredientStateIngredient(ctx, &mealplanninggrpc.GetValidIngredientStateIngredientRequest{ValidIngredientStateIngredientID: exampleResult.ID})
		assert.Equal(t, exampleResult.ID, result.Result.ID)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mvem)
	})
}

func TestServiceImpl_GetValidIngredientStateIngredients(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleResult := mealplanningfakes.BuildFakeValidIngredientStateIngredientsList()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("ListValidIngredientStateIngredients", testutils.ContextMatcher, testutils.QueryFilterMatcher).Return(exampleResult.Data, "", nil)
		s.validEnumerationsManager = mvem

		result, err := s.GetValidIngredientStateIngredients(ctx, &mealplanninggrpc.GetValidIngredientStateIngredientsRequest{})
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result.Results, len(exampleResult.Data))

		mock.AssertExpectationsForObjects(t, mvem)
	})
}

func TestServiceImpl_GetValidIngredientStateIngredientsByIngredient(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleID := mealplanningfakes.BuildFakeID()
		exampleResult := mealplanningfakes.BuildFakeValidIngredientStateIngredientsList()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("SearchValidIngredientStateIngredientsByIngredient", testutils.ContextMatcher, exampleID, testutils.QueryFilterMatcher).Return(exampleResult.Data, nil)
		s.validEnumerationsManager = mvem

		result, err := s.GetValidIngredientStateIngredientsByIngredient(ctx, &mealplanninggrpc.GetValidIngredientStateIngredientsByIngredientRequest{
			ValidIngredientID: exampleID,
		})
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result.Results, len(exampleResult.Data))

		mock.AssertExpectationsForObjects(t, mvem)
	})
}

func TestServiceImpl_GetValidIngredientStateIngredientsByIngredientState(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleID := mealplanningfakes.BuildFakeID()
		exampleResult := mealplanningfakes.BuildFakeValidIngredientStateIngredientsList()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("SearchValidIngredientStateIngredientsByIngredientState", testutils.ContextMatcher, exampleID, testutils.QueryFilterMatcher).Return(exampleResult.Data, nil)
		s.validEnumerationsManager = mvem

		result, err := s.GetValidIngredientStateIngredientsByIngredientState(ctx, &mealplanninggrpc.GetValidIngredientStateIngredientsByIngredientStateRequest{
			ValidIngredientStateID: exampleID,
		})
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result.Results, len(exampleResult.Data))

		mock.AssertExpectationsForObjects(t, mvem)
	})
}

func TestServiceImpl_GetValidIngredientStates(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleResult := mealplanningfakes.BuildFakeValidIngredientStatesList()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("ListValidIngredientStates", testutils.ContextMatcher, testutils.QueryFilterMatcher).Return(exampleResult.Data, "", nil)
		s.validEnumerationsManager = mvem

		result, err := s.GetValidIngredientStates(ctx, &mealplanninggrpc.GetValidIngredientStatesRequest{})
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result.Results, len(exampleResult.Data))

		mock.AssertExpectationsForObjects(t, mvem)
	})
}

func TestServiceImpl_GetValidIngredients(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleResult := mealplanningfakes.BuildFakeValidIngredientsList()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("ListValidIngredients", testutils.ContextMatcher, testutils.QueryFilterMatcher).Return(exampleResult.Data, "", nil)
		s.validEnumerationsManager = mvem

		result, err := s.GetValidIngredients(ctx, &mealplanninggrpc.GetValidIngredientsRequest{})
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result.Results, len(exampleResult.Data))

		mock.AssertExpectationsForObjects(t, mvem)
	})
}

func TestServiceImpl_GetValidInstrument(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleResult := mealplanningfakes.BuildFakeValidInstrument()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("ReadValidInstrument", testutils.ContextMatcher, exampleResult.ID).Return(exampleResult, nil)
		s.validEnumerationsManager = mvem

		result, err := s.GetValidInstrument(ctx, &mealplanninggrpc.GetValidInstrumentRequest{ValidInstrumentID: exampleResult.ID})
		assert.Equal(t, exampleResult.ID, result.Result.ID)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mvem)
	})
}

func TestServiceImpl_GetValidInstruments(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleResult := mealplanningfakes.BuildFakeValidInstrumentsList()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("ListValidInstruments", testutils.ContextMatcher, testutils.QueryFilterMatcher).Return(exampleResult.Data, "", nil)
		s.validEnumerationsManager = mvem

		result, err := s.GetValidInstruments(ctx, &mealplanninggrpc.GetValidInstrumentsRequest{})
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result.Results, len(exampleResult.Data))

		mock.AssertExpectationsForObjects(t, mvem)
	})
}

func TestServiceImpl_GetValidMeasurementUnit(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleResult := mealplanningfakes.BuildFakeValidMeasurementUnit()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("ReadValidMeasurementUnit", testutils.ContextMatcher, exampleResult.ID).Return(exampleResult, nil)
		s.validEnumerationsManager = mvem

		result, err := s.GetValidMeasurementUnit(ctx, &mealplanninggrpc.GetValidMeasurementUnitRequest{ValidMeasurementUnitID: exampleResult.ID})
		assert.Equal(t, exampleResult.ID, result.Result.ID)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mvem)
	})
}

func TestServiceImpl_GetValidMeasurementUnitConversion(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleResult := mealplanningfakes.BuildFakeValidMeasurementUnitConversion()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("ReadValidMeasurementUnitConversion", testutils.ContextMatcher, exampleResult.ID).Return(exampleResult, nil)
		s.validEnumerationsManager = mvem

		result, err := s.GetValidMeasurementUnitConversion(ctx, &mealplanninggrpc.GetValidMeasurementUnitConversionRequest{ValidMeasurementUnitConversionID: exampleResult.ID})
		assert.Equal(t, exampleResult.ID, result.Result.ID)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mvem)
	})
}

func TestServiceImpl_GetValidMeasurementUnitConversionsFromUnit(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleID := mealplanningfakes.BuildFakeID()
		exampleResult := mealplanningfakes.BuildFakeValidMeasurementUnitConversionsList()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("ValidMeasurementUnitConversionsFromMeasurementUnit", testutils.ContextMatcher, exampleID).Return(exampleResult.Data, nil)
		s.validEnumerationsManager = mvem

		result, err := s.GetValidMeasurementUnitConversionsFromUnit(ctx, &mealplanninggrpc.GetValidMeasurementUnitConversionsFromUnitRequest{
			ValidMeasurementUnitID: exampleID,
		})
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result.Results, len(exampleResult.Data))

		mock.AssertExpectationsForObjects(t, mvem)
	})
}

func TestServiceImpl_GetValidMeasurementUnitConversionsToUnit(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleID := mealplanningfakes.BuildFakeID()
		exampleResult := mealplanningfakes.BuildFakeValidMeasurementUnitConversionsList()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("ValidMeasurementUnitConversionsToMeasurementUnit", testutils.ContextMatcher, exampleID).Return(exampleResult.Data, nil)
		s.validEnumerationsManager = mvem

		result, err := s.GetValidMeasurementUnitConversionsToUnit(ctx, &mealplanninggrpc.GetValidMeasurementUnitConversionsToUnitRequest{
			ValidMeasurementUnitID: exampleID,
		})
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result.Results, len(exampleResult.Data))

		mock.AssertExpectationsForObjects(t, mvem)
	})
}

func TestServiceImpl_GetValidMeasurementUnits(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleResult := mealplanningfakes.BuildFakeValidMeasurementUnitsList()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("ListValidMeasurementUnits", testutils.ContextMatcher, testutils.QueryFilterMatcher).Return(exampleResult.Data, "", nil)
		s.validEnumerationsManager = mvem

		result, err := s.GetValidMeasurementUnits(ctx, &mealplanninggrpc.GetValidMeasurementUnitsRequest{})
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result.Results, len(exampleResult.Data))

		mock.AssertExpectationsForObjects(t, mvem)
	})
}

func TestServiceImpl_GetValidPreparation(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleResult := mealplanningfakes.BuildFakeValidPreparation()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("ReadValidPreparation", testutils.ContextMatcher, exampleResult.ID).Return(exampleResult, nil)
		s.validEnumerationsManager = mvem

		result, err := s.GetValidPreparation(ctx, &mealplanninggrpc.GetValidPreparationRequest{ValidPreparationID: exampleResult.ID})
		assert.Equal(t, exampleResult.ID, result.Result.ID)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mvem)
	})
}

func TestServiceImpl_GetValidPreparationInstrument(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleResult := mealplanningfakes.BuildFakeValidPreparationInstrument()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("ReadValidPreparationInstrument", testutils.ContextMatcher, exampleResult.ID).Return(exampleResult, nil)
		s.validEnumerationsManager = mvem

		result, err := s.GetValidPreparationInstrument(ctx, &mealplanninggrpc.GetValidPreparationInstrumentRequest{ValidPreparationInstrumentID: exampleResult.ID})
		assert.Equal(t, exampleResult.ID, result.Result.ID)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mvem)
	})
}

func TestServiceImpl_GetValidPreparationInstruments(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleResult := mealplanningfakes.BuildFakeValidPreparationInstrumentsList()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("ListValidPreparationInstruments", testutils.ContextMatcher, testutils.QueryFilterMatcher).Return(exampleResult.Data, "", nil)
		s.validEnumerationsManager = mvem

		result, err := s.GetValidPreparationInstruments(ctx, &mealplanninggrpc.GetValidPreparationInstrumentsRequest{})
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result.Results, len(exampleResult.Data))

		mock.AssertExpectationsForObjects(t, mvem)
	})
}

func TestServiceImpl_GetValidPreparationInstrumentsByInstrument(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleID := mealplanningfakes.BuildFakeID()
		exampleResult := mealplanningfakes.BuildFakeValidPreparationInstrumentsList()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("SearchValidPreparationInstrumentsByInstrument", testutils.ContextMatcher, exampleID, testutils.QueryFilterMatcher).Return(exampleResult.Data, nil)
		s.validEnumerationsManager = mvem

		result, err := s.GetValidPreparationInstrumentsByInstrument(ctx, &mealplanninggrpc.GetValidPreparationInstrumentsByInstrumentRequest{
			ValidInstrumentID: exampleID,
		})
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result.Results, len(exampleResult.Data))

		mock.AssertExpectationsForObjects(t, mvem)
	})
}

func TestServiceImpl_GetValidPreparationInstrumentsByPreparation(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleID := mealplanningfakes.BuildFakeID()
		exampleResult := mealplanningfakes.BuildFakeValidPreparationInstrumentsList()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("SearchValidPreparationInstrumentsByPreparation", testutils.ContextMatcher, exampleID, testutils.QueryFilterMatcher).Return(exampleResult.Data, nil)
		s.validEnumerationsManager = mvem

		result, err := s.GetValidPreparationInstrumentsByPreparation(ctx, &mealplanninggrpc.GetValidPreparationInstrumentsByPreparationRequest{
			ValidPreparationID: exampleID,
		})
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result.Results, len(exampleResult.Data))

		mock.AssertExpectationsForObjects(t, mvem)
	})
}

func TestServiceImpl_GetValidPreparationVessel(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleResult := mealplanningfakes.BuildFakeValidPreparationVessel()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("ReadValidPreparationVessel", testutils.ContextMatcher, exampleResult.ID).Return(exampleResult, nil)
		s.validEnumerationsManager = mvem

		result, err := s.GetValidPreparationVessel(ctx, &mealplanninggrpc.GetValidPreparationVesselRequest{ValidPreparationVesselID: exampleResult.ID})
		assert.Equal(t, exampleResult.ID, result.Result.ID)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mvem)
	})
}

func TestServiceImpl_GetValidPreparationVessels(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleResult := mealplanningfakes.BuildFakeValidPreparationVesselsList()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("ListValidPreparationVessels", testutils.ContextMatcher, testutils.QueryFilterMatcher).Return(exampleResult.Data, "", nil)
		s.validEnumerationsManager = mvem

		result, err := s.GetValidPreparationVessels(ctx, &mealplanninggrpc.GetValidPreparationVesselsRequest{})
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result.Results, len(exampleResult.Data))

		mock.AssertExpectationsForObjects(t, mvem)
	})
}

func TestServiceImpl_GetValidPreparationVesselsByPreparation(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleID := mealplanningfakes.BuildFakeID()
		exampleResult := mealplanningfakes.BuildFakeValidPreparationVesselsList()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("SearchValidPreparationVesselsByPreparation", testutils.ContextMatcher, exampleID, testutils.QueryFilterMatcher).Return(exampleResult.Data, nil)
		s.validEnumerationsManager = mvem

		result, err := s.GetValidPreparationVesselsByPreparation(ctx, &mealplanninggrpc.GetValidPreparationVesselsByPreparationRequest{
			ValidPreparationID: exampleID,
		})
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result.Results, len(exampleResult.Data))

		mock.AssertExpectationsForObjects(t, mvem)
	})
}

func TestServiceImpl_GetValidPreparationVesselsByVessel(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleID := mealplanningfakes.BuildFakeID()
		exampleResult := mealplanningfakes.BuildFakeValidPreparationVesselsList()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("SearchValidPreparationVesselsByVessel", testutils.ContextMatcher, exampleID, testutils.QueryFilterMatcher).Return(exampleResult.Data, nil)
		s.validEnumerationsManager = mvem

		result, err := s.GetValidPreparationVesselsByVessel(ctx, &mealplanninggrpc.GetValidPreparationVesselsByVesselRequest{
			ValidVesselID: exampleID,
		})
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result.Results, len(exampleResult.Data))

		mock.AssertExpectationsForObjects(t, mvem)
	})
}

func TestServiceImpl_GetValidPreparations(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleResult := mealplanningfakes.BuildFakeValidPreparationsList()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("ListValidPreparations", testutils.ContextMatcher, testutils.QueryFilterMatcher).Return(exampleResult.Data, "", nil)
		s.validEnumerationsManager = mvem

		result, err := s.GetValidPreparations(ctx, &mealplanninggrpc.GetValidPreparationsRequest{})
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result.Results, len(exampleResult.Data))

		mock.AssertExpectationsForObjects(t, mvem)
	})
}

func TestServiceImpl_GetValidVessel(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleResult := mealplanningfakes.BuildFakeValidVessel()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("ReadValidVessel", testutils.ContextMatcher, exampleResult.ID).Return(exampleResult, nil)
		s.validEnumerationsManager = mvem

		result, err := s.GetValidVessel(ctx, &mealplanninggrpc.GetValidVesselRequest{ValidVesselID: exampleResult.ID})
		assert.Equal(t, exampleResult.ID, result.Result.ID)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mvem)
	})
}

func TestServiceImpl_GetValidVessels(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleResult := mealplanningfakes.BuildFakeValidVesselsList()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("ListValidVessels", testutils.ContextMatcher, testutils.QueryFilterMatcher).Return(exampleResult.Data, "", nil)
		s.validEnumerationsManager = mvem

		result, err := s.GetValidVessels(ctx, &mealplanninggrpc.GetValidVesselsRequest{})
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result.Results, len(exampleResult.Data))

		mock.AssertExpectationsForObjects(t, mvem)
	})
}

func TestServiceImpl_SearchForValidIngredientGroups(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleResult := mealplanningfakes.BuildFakeValidIngredientGroupsList()
		exampleRequest := fake.BuildFakeForTest[mealplanninggrpc.SearchForValidIngredientGroupsRequest](t)

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("SearchValidIngredientGroups", testutils.ContextMatcher, exampleRequest.Query, exampleRequest.UseDatabase, testutils.QueryFilterMatcher).Return(exampleResult.Data, nil)
		s.validEnumerationsManager = mvem

		result, err := s.SearchForValidIngredientGroups(ctx, exampleRequest)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result.Results, len(exampleResult.Data))

		mock.AssertExpectationsForObjects(t, mvem)
	})
}

func TestServiceImpl_SearchForValidIngredientStates(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleResult := mealplanningfakes.BuildFakeValidIngredientStatesList()
		exampleRequest := fake.BuildFakeForTest[mealplanninggrpc.SearchForValidIngredientStatesRequest](t)

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("SearchValidIngredientStates", testutils.ContextMatcher, exampleRequest.Query, exampleRequest.UseDatabase, testutils.QueryFilterMatcher).Return(exampleResult.Data, nil)
		s.validEnumerationsManager = mvem

		result, err := s.SearchForValidIngredientStates(ctx, exampleRequest)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result.Results, len(exampleResult.Data))

		mock.AssertExpectationsForObjects(t, mvem)
	})
}

func TestServiceImpl_SearchForValidIngredients(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleResult := mealplanningfakes.BuildFakeValidIngredientsList()
		exampleRequest := fake.BuildFakeForTest[mealplanninggrpc.SearchForValidIngredientsRequest](t)

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("SearchValidIngredients", testutils.ContextMatcher, exampleRequest.Query, exampleRequest.UseDatabase, testutils.QueryFilterMatcher).Return(exampleResult.Data, nil)
		s.validEnumerationsManager = mvem

		result, err := s.SearchForValidIngredients(ctx, exampleRequest)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result.Results, len(exampleResult.Data))

		mock.AssertExpectationsForObjects(t, mvem)
	})
}

func TestServiceImpl_SearchForValidInstruments(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleResult := mealplanningfakes.BuildFakeValidInstrumentsList()
		exampleRequest := fake.BuildFakeForTest[mealplanninggrpc.SearchForValidInstrumentsRequest](t)

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("SearchValidInstruments", testutils.ContextMatcher, exampleRequest.Query, exampleRequest.UseDatabase, testutils.QueryFilterMatcher).Return(exampleResult.Data, nil)
		s.validEnumerationsManager = mvem

		result, err := s.SearchForValidInstruments(ctx, exampleRequest)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result.Results, len(exampleResult.Data))

		mock.AssertExpectationsForObjects(t, mvem)
	})
}

func TestServiceImpl_SearchForValidMeasurementUnits(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleResult := mealplanningfakes.BuildFakeValidMeasurementUnitsList()
		exampleRequest := fake.BuildFakeForTest[mealplanninggrpc.SearchForValidMeasurementUnitsRequest](t)

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("SearchValidMeasurementUnits", testutils.ContextMatcher, exampleRequest.Query, exampleRequest.UseDatabase, testutils.QueryFilterMatcher).Return(exampleResult.Data, nil)
		s.validEnumerationsManager = mvem

		result, err := s.SearchForValidMeasurementUnits(ctx, exampleRequest)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result.Results, len(exampleResult.Data))

		mock.AssertExpectationsForObjects(t, mvem)
	})
}

func TestServiceImpl_SearchForValidPreparations(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleResult := mealplanningfakes.BuildFakeValidPreparationsList()
		exampleRequest := fake.BuildFakeForTest[mealplanninggrpc.SearchForValidPreparationsRequest](t)

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("SearchValidPreparations", testutils.ContextMatcher, exampleRequest.Query, exampleRequest.UseDatabase, testutils.QueryFilterMatcher).Return(exampleResult.Data, nil)
		s.validEnumerationsManager = mvem

		result, err := s.SearchForValidPreparations(ctx, exampleRequest)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result.Results, len(exampleResult.Data))

		mock.AssertExpectationsForObjects(t, mvem)
	})
}

func TestServiceImpl_SearchForValidVessels(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleResult := mealplanningfakes.BuildFakeValidVesselsList()
		exampleRequest := fake.BuildFakeForTest[mealplanninggrpc.SearchForValidVesselsRequest](t)

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("SearchValidVessels", testutils.ContextMatcher, exampleRequest.Query, exampleRequest.UseDatabase, testutils.QueryFilterMatcher).Return(exampleResult.Data, nil)
		s.validEnumerationsManager = mvem

		result, err := s.SearchForValidVessels(ctx, exampleRequest)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result.Results, len(exampleResult.Data))

		mock.AssertExpectationsForObjects(t, mvem)
	})
}

func TestServiceImpl_SearchValidIngredientsByPreparation(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleResult := mealplanningfakes.BuildFakeValidIngredientsList()
		exampleRequest := fake.BuildFakeForTest[mealplanninggrpc.SearchValidIngredientsByPreparationRequest](t)

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("SearchValidIngredientsByPreparationAndIngredientName", testutils.ContextMatcher, exampleRequest.ValidPreparationID, exampleRequest.Query, testutils.QueryFilterMatcher).Return(exampleResult.Data, nil)
		s.validEnumerationsManager = mvem

		result, err := s.SearchValidIngredientsByPreparation(ctx, exampleRequest)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result.Results, len(exampleResult.Data))

		mock.AssertExpectationsForObjects(t, mvem)
	})
}

func TestServiceImpl_SearchValidMeasurementUnitsByIngredient(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleResult := mealplanningfakes.BuildFakeValidMeasurementUnitsList()
		exampleRequest := fake.BuildFakeForTest[mealplanninggrpc.SearchValidMeasurementUnitsByIngredientRequest](t)

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("SearchValidMeasurementUnitsByIngredientID", testutils.ContextMatcher, exampleRequest.ValidIngredientID, testutils.QueryFilterMatcher).Return(exampleResult.Data, nil)
		s.validEnumerationsManager = mvem

		result, err := s.SearchValidMeasurementUnitsByIngredient(ctx, exampleRequest)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result.Results, len(exampleResult.Data))

		mock.AssertExpectationsForObjects(t, mvem)
	})
}

func TestServiceImpl_UpdateValidIngredient(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		exampleRequest := fake.BuildFakeForTest[mealplanninggrpc.UpdateValidIngredientRequest](t)
		exampleResponse := mealplanningfakes.BuildFakeValidIngredient()

		s := buildServiceImplForTest(t)

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("UpdateValidIngredient", testutils.ContextMatcher, exampleRequest.ValidIngredientID, testutils.MatchType[*mealplanning.ValidIngredientUpdateRequestInput]()).Return(exampleResponse, nil)
		s.validEnumerationsManager = mvem

		res, err := s.UpdateValidIngredient(ctx, exampleRequest)
		assert.NoError(t, err)
		assert.Equal(t, exampleResponse.ID, res.Result.ID)

		mock.AssertExpectationsForObjects(t, mvem)
	})
}

func TestServiceImpl_UpdateValidIngredientGroup(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		exampleRequest := fake.BuildFakeForTest[mealplanninggrpc.UpdateValidIngredientGroupRequest](t)
		exampleResponse := mealplanningfakes.BuildFakeValidIngredientGroup()

		s := buildServiceImplForTest(t)

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("UpdateValidIngredientGroup", testutils.ContextMatcher, exampleRequest.ValidIngredientGroupID, testutils.MatchType[*mealplanning.ValidIngredientGroupUpdateRequestInput]()).Return(exampleResponse, nil)
		s.validEnumerationsManager = mvem

		res, err := s.UpdateValidIngredientGroup(ctx, exampleRequest)
		assert.NoError(t, err)
		assert.Equal(t, exampleResponse.ID, res.Result.ID)

		mock.AssertExpectationsForObjects(t, mvem)
	})
}

func TestServiceImpl_UpdateValidIngredientMeasurementUnit(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		exampleRequest := fake.BuildFakeForTest[mealplanninggrpc.UpdateValidIngredientMeasurementUnitRequest](t)
		exampleResponse := mealplanningfakes.BuildFakeValidIngredientMeasurementUnit()

		s := buildServiceImplForTest(t)

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("UpdateValidIngredientMeasurementUnit", testutils.ContextMatcher, exampleRequest.ValidIngredientMeasurementUnitID, testutils.MatchType[*mealplanning.ValidIngredientMeasurementUnitUpdateRequestInput]()).Return(exampleResponse, nil)
		s.validEnumerationsManager = mvem

		res, err := s.UpdateValidIngredientMeasurementUnit(ctx, exampleRequest)
		assert.NoError(t, err)
		assert.Equal(t, exampleResponse.ID, res.Result.ID)

		mock.AssertExpectationsForObjects(t, mvem)
	})
}

func TestServiceImpl_UpdateValidIngredientPreparation(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		exampleRequest := fake.BuildFakeForTest[mealplanninggrpc.UpdateValidIngredientPreparationRequest](t)
		exampleResponse := mealplanningfakes.BuildFakeValidIngredientPreparation()

		s := buildServiceImplForTest(t)

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("UpdateValidIngredientPreparation", testutils.ContextMatcher, exampleRequest.ValidIngredientPreparationID, testutils.MatchType[*mealplanning.ValidIngredientPreparationUpdateRequestInput]()).Return(exampleResponse, nil)
		s.validEnumerationsManager = mvem

		res, err := s.UpdateValidIngredientPreparation(ctx, exampleRequest)
		assert.NoError(t, err)
		assert.Equal(t, exampleResponse.ID, res.Result.ID)

		mock.AssertExpectationsForObjects(t, mvem)
	})
}

func TestServiceImpl_UpdateValidIngredientState(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		exampleRequest := fake.BuildFakeForTest[mealplanninggrpc.UpdateValidIngredientStateRequest](t)
		exampleResponse := mealplanningfakes.BuildFakeValidIngredientState()

		s := buildServiceImplForTest(t)

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("UpdateValidIngredientState", testutils.ContextMatcher, exampleRequest.ValidIngredientStateID, testutils.MatchType[*mealplanning.ValidIngredientStateUpdateRequestInput]()).Return(exampleResponse, nil)
		s.validEnumerationsManager = mvem

		res, err := s.UpdateValidIngredientState(ctx, exampleRequest)
		assert.NoError(t, err)
		assert.Equal(t, exampleResponse.ID, res.Result.ID)

		mock.AssertExpectationsForObjects(t, mvem)
	})
}

func TestServiceImpl_UpdateValidIngredientStateIngredient(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		exampleRequest := fake.BuildFakeForTest[mealplanninggrpc.UpdateValidIngredientStateIngredientRequest](t)
		exampleResponse := mealplanningfakes.BuildFakeValidIngredientStateIngredient()

		s := buildServiceImplForTest(t)

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("UpdateValidIngredientStateIngredient", testutils.ContextMatcher, exampleRequest.ValidIngredientStateIngredientID, testutils.MatchType[*mealplanning.ValidIngredientStateIngredientUpdateRequestInput]()).Return(exampleResponse, nil)
		s.validEnumerationsManager = mvem

		res, err := s.UpdateValidIngredientStateIngredient(ctx, exampleRequest)
		assert.NoError(t, err)
		assert.Equal(t, exampleResponse.ID, res.Result.ID)

		mock.AssertExpectationsForObjects(t, mvem)
	})
}

func TestServiceImpl_UpdateValidInstrument(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		exampleRequest := fake.BuildFakeForTest[mealplanninggrpc.UpdateValidInstrumentRequest](t)
		exampleResponse := mealplanningfakes.BuildFakeValidInstrument()

		s := buildServiceImplForTest(t)

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("UpdateValidInstrument", testutils.ContextMatcher, exampleRequest.ValidInstrumentID, testutils.MatchType[*mealplanning.ValidInstrumentUpdateRequestInput]()).Return(exampleResponse, nil)
		s.validEnumerationsManager = mvem

		res, err := s.UpdateValidInstrument(ctx, exampleRequest)
		assert.NoError(t, err)
		assert.Equal(t, exampleResponse.ID, res.Result.ID)

		mock.AssertExpectationsForObjects(t, mvem)
	})
}

func TestServiceImpl_UpdateValidMeasurementUnit(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		exampleRequest := fake.BuildFakeForTest[mealplanninggrpc.UpdateValidMeasurementUnitRequest](t)
		exampleResponse := mealplanningfakes.BuildFakeValidMeasurementUnit()

		s := buildServiceImplForTest(t)

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("UpdateValidMeasurementUnit", testutils.ContextMatcher, exampleRequest.ValidMeasurementUnitID, testutils.MatchType[*mealplanning.ValidMeasurementUnitUpdateRequestInput]()).Return(exampleResponse, nil)
		s.validEnumerationsManager = mvem

		res, err := s.UpdateValidMeasurementUnit(ctx, exampleRequest)
		assert.NoError(t, err)
		assert.Equal(t, exampleResponse.ID, res.Result.ID)

		mock.AssertExpectationsForObjects(t, mvem)
	})
}

func TestServiceImpl_UpdateValidMeasurementUnitConversion(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		exampleRequest := fake.BuildFakeForTest[mealplanninggrpc.UpdateValidMeasurementUnitConversionRequest](t)
		exampleResponse := mealplanningfakes.BuildFakeValidMeasurementUnitConversion()

		s := buildServiceImplForTest(t)

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("UpdateValidMeasurementUnitConversion", testutils.ContextMatcher, exampleRequest.ValidMeasurementUnitConversionID, testutils.MatchType[*mealplanning.ValidMeasurementUnitConversionUpdateRequestInput]()).Return(exampleResponse, nil)
		s.validEnumerationsManager = mvem

		res, err := s.UpdateValidMeasurementUnitConversion(ctx, exampleRequest)
		assert.NoError(t, err)
		assert.Equal(t, exampleResponse.ID, res.Result.ID)

		mock.AssertExpectationsForObjects(t, mvem)
	})
}

func TestServiceImpl_UpdateValidPreparation(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		exampleRequest := fake.BuildFakeForTest[mealplanninggrpc.UpdateValidPreparationRequest](t)
		exampleResponse := mealplanningfakes.BuildFakeValidPreparation()

		s := buildServiceImplForTest(t)

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("UpdateValidPreparation", testutils.ContextMatcher, exampleRequest.ValidPreparationID, testutils.MatchType[*mealplanning.ValidPreparationUpdateRequestInput]()).Return(exampleResponse, nil)
		s.validEnumerationsManager = mvem

		res, err := s.UpdateValidPreparation(ctx, exampleRequest)
		assert.NoError(t, err)
		assert.Equal(t, exampleResponse.ID, res.Result.ID)

		mock.AssertExpectationsForObjects(t, mvem)
	})
}

func TestServiceImpl_UpdateValidPreparationInstrument(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		exampleRequest := fake.BuildFakeForTest[mealplanninggrpc.UpdateValidPreparationInstrumentRequest](t)
		exampleResponse := mealplanningfakes.BuildFakeValidPreparationInstrument()

		s := buildServiceImplForTest(t)

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("UpdateValidPreparationInstrument", testutils.ContextMatcher, exampleRequest.ValidPreparationInstrumentID, testutils.MatchType[*mealplanning.ValidPreparationInstrumentUpdateRequestInput]()).Return(exampleResponse, nil)
		s.validEnumerationsManager = mvem

		res, err := s.UpdateValidPreparationInstrument(ctx, exampleRequest)
		assert.NoError(t, err)
		assert.Equal(t, exampleResponse.ID, res.Result.ID)

		mock.AssertExpectationsForObjects(t, mvem)
	})
}

func TestServiceImpl_UpdateValidPreparationVessel(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		exampleRequest := fake.BuildFakeForTest[mealplanninggrpc.UpdateValidPreparationVesselRequest](t)
		exampleResponse := mealplanningfakes.BuildFakeValidPreparationVessel()

		s := buildServiceImplForTest(t)

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("UpdateValidPreparationVessel", testutils.ContextMatcher, exampleRequest.ValidPreparationVesselID, testutils.MatchType[*mealplanning.ValidPreparationVesselUpdateRequestInput]()).Return(exampleResponse, nil)
		s.validEnumerationsManager = mvem

		res, err := s.UpdateValidPreparationVessel(ctx, exampleRequest)
		assert.NoError(t, err)
		assert.Equal(t, exampleResponse.ID, res.Result.ID)

		mock.AssertExpectationsForObjects(t, mvem)
	})
}

func TestServiceImpl_UpdateValidVessel(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		exampleRequest := fake.BuildFakeForTest[mealplanninggrpc.UpdateValidVesselRequest](t)
		exampleResponse := mealplanningfakes.BuildFakeValidVessel()

		s := buildServiceImplForTest(t)

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("UpdateValidVessel", testutils.ContextMatcher, exampleRequest.ValidVesselID, testutils.MatchType[*mealplanning.ValidVesselUpdateRequestInput]()).Return(exampleResponse, nil)
		s.validEnumerationsManager = mvem

		res, err := s.UpdateValidVessel(ctx, exampleRequest)
		assert.NoError(t, err)
		assert.Equal(t, exampleResponse.ID, res.Result.ID)

		mock.AssertExpectationsForObjects(t, mvem)
	})
}
