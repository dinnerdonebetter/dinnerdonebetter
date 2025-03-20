package grpc

import (
	"testing"

	"github.com/dinnerdonebetter/backend/internal/grpc/messages"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/lib/testutils"
	mockmanagers "github.com/dinnerdonebetter/backend/internal/services/eating/managers/mock"
	"github.com/dinnerdonebetter/backend/internal/services/eating/types"
	"github.com/dinnerdonebetter/backend/internal/services/eating/types/fakes"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func buildServiceImplForTest(t *testing.T) *serviceImpl {
	t.Helper()

	return &serviceImpl{
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

		exampleValidIngredientID := fakes.BuildFakeID()

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("ArchiveValidIngredient", testutils.ContextMatcher, exampleValidIngredientID).Return(nil)
		s.validEnumerationsManager = mvem

		res, err := s.ArchiveValidIngredient(ctx, &messages.ArchiveValidIngredientRequest{ValidIngredientID: exampleValidIngredientID})
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

		exampleValidIngredientGroupID := fakes.BuildFakeID()

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("ArchiveValidIngredientGroup", testutils.ContextMatcher, exampleValidIngredientGroupID).Return(nil)
		s.validEnumerationsManager = mvem

		res, err := s.ArchiveValidIngredientGroup(ctx, &messages.ArchiveValidIngredientGroupRequest{ValidIngredientGroupID: exampleValidIngredientGroupID})
		assert.NotNil(t, res)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mvem)
	})
}

func TestServiceImpl_ArchiveValidIngredientMeasurementUnit(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

	})
}

func TestServiceImpl_ArchiveValidIngredientPreparation(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		exampleValidIngredientPreparationID := fakes.BuildFakeID()

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("ArchiveValidIngredientPreparation", testutils.ContextMatcher, exampleValidIngredientPreparationID).Return(nil)
		s.validEnumerationsManager = mvem

		res, err := s.ArchiveValidIngredientPreparation(ctx, &messages.ArchiveValidIngredientPreparationRequest{ValidIngredientPreparationID: exampleValidIngredientPreparationID})
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

		exampleValidIngredientStateID := fakes.BuildFakeID()

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("ArchiveValidIngredientState", testutils.ContextMatcher, exampleValidIngredientStateID).Return(nil)
		s.validEnumerationsManager = mvem

		res, err := s.ArchiveValidIngredientState(ctx, &messages.ArchiveValidIngredientStateRequest{ValidIngredientStateID: exampleValidIngredientStateID})
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

		exampleValidIngredientStateIngredientID := fakes.BuildFakeID()

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("ArchiveValidIngredientStateIngredient", testutils.ContextMatcher, exampleValidIngredientStateIngredientID).Return(nil)
		s.validEnumerationsManager = mvem

		res, err := s.ArchiveValidIngredientStateIngredient(ctx, &messages.ArchiveValidIngredientStateIngredientRequest{ValidIngredientStateIngredientID: exampleValidIngredientStateIngredientID})
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

		exampleValidInstrumentID := fakes.BuildFakeID()

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("ArchiveValidInstrument", testutils.ContextMatcher, exampleValidInstrumentID).Return(nil)
		s.validEnumerationsManager = mvem

		res, err := s.ArchiveValidInstrument(ctx, &messages.ArchiveValidInstrumentRequest{ValidInstrumentID: exampleValidInstrumentID})
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

		exampleValidMeasurementUnitID := fakes.BuildFakeID()

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("ArchiveValidMeasurementUnit", testutils.ContextMatcher, exampleValidMeasurementUnitID).Return(nil)
		s.validEnumerationsManager = mvem

		res, err := s.ArchiveValidMeasurementUnit(ctx, &messages.ArchiveValidMeasurementUnitRequest{ValidMeasurementUnitID: exampleValidMeasurementUnitID})
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

		exampleValidMeasurementUnitConversionID := fakes.BuildFakeID()

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("ArchiveValidMeasurementUnitConversion", testutils.ContextMatcher, exampleValidMeasurementUnitConversionID).Return(nil)
		s.validEnumerationsManager = mvem

		res, err := s.ArchiveValidMeasurementUnitConversion(ctx, &messages.ArchiveValidMeasurementUnitConversionRequest{ValidMeasurementUnitConversionID: exampleValidMeasurementUnitConversionID})
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

		exampleValidPreparationID := fakes.BuildFakeID()

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("ArchiveValidPreparation", testutils.ContextMatcher, exampleValidPreparationID).Return(nil)
		s.validEnumerationsManager = mvem

		res, err := s.ArchiveValidPreparation(ctx, &messages.ArchiveValidPreparationRequest{ValidPreparationID: exampleValidPreparationID})
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

		exampleValidPreparationInstrumentID := fakes.BuildFakeID()

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("ArchiveValidPreparationInstrument", testutils.ContextMatcher, exampleValidPreparationInstrumentID).Return(nil)
		s.validEnumerationsManager = mvem

		res, err := s.ArchiveValidPreparationInstrument(ctx, &messages.ArchiveValidPreparationInstrumentRequest{ValidPreparationInstrumentID: exampleValidPreparationInstrumentID})
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

		exampleValidPreparationVesselID := fakes.BuildFakeID()

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("ArchiveValidPreparationVessel", testutils.ContextMatcher, exampleValidPreparationVesselID).Return(nil)
		s.validEnumerationsManager = mvem

		res, err := s.ArchiveValidPreparationVessel(ctx, &messages.ArchiveValidPreparationVesselRequest{ValidPreparationVesselID: exampleValidPreparationVesselID})
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

		exampleValidVesselID := fakes.BuildFakeID()

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("ArchiveValidVessel", testutils.ContextMatcher, exampleValidVesselID).Return(nil)
		s.validEnumerationsManager = mvem

		res, err := s.ArchiveValidVessel(ctx, &messages.ArchiveValidVesselRequest{ValidVesselID: exampleValidVesselID})
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

		exampleValidIngredient := fakes.BuildFakeValidIngredient()

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("CreateValidIngredient", testutils.ContextMatcher, testutils.MatchType[*types.ValidIngredientCreationRequestInput]()).Return(exampleValidIngredient, nil)
		s.validEnumerationsManager = mvem

		exampleInput := fakes.BuildFake[messages.CreateValidIngredientRequest]()

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

		exampleValidIngredientGroup := fakes.BuildFakeValidIngredientGroup()

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("CreateValidIngredientGroup", testutils.ContextMatcher, testutils.MatchType[*types.ValidIngredientGroupCreationRequestInput]()).Return(exampleValidIngredientGroup, nil)
		s.validEnumerationsManager = mvem

		exampleInput := fakes.BuildFake[messages.CreateValidIngredientGroupRequest]()

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

		exampleValidIngredientMeasurementUnit := fakes.BuildFakeValidIngredientMeasurementUnit()

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("CreateValidIngredientMeasurementUnit", testutils.ContextMatcher, testutils.MatchType[*types.ValidIngredientMeasurementUnitCreationRequestInput]()).Return(exampleValidIngredientMeasurementUnit, nil)
		s.validEnumerationsManager = mvem

		exampleInput := fakes.BuildFake[messages.CreateValidIngredientMeasurementUnitRequest]()

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

		exampleValidIngredientPreparation := fakes.BuildFakeValidIngredientPreparation()

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("CreateValidIngredientPreparation", testutils.ContextMatcher, testutils.MatchType[*types.ValidIngredientPreparationCreationRequestInput]()).Return(exampleValidIngredientPreparation, nil)
		s.validEnumerationsManager = mvem

		exampleInput := fakes.BuildFake[messages.CreateValidIngredientPreparationRequest]()

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

		exampleValidIngredientState := fakes.BuildFakeValidIngredientState()

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("CreateValidIngredientState", testutils.ContextMatcher, testutils.MatchType[*types.ValidIngredientStateCreationRequestInput]()).Return(exampleValidIngredientState, nil)
		s.validEnumerationsManager = mvem

		exampleInput := fakes.BuildFake[messages.CreateValidIngredientStateRequest]()

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

		exampleValidIngredientStateIngredient := fakes.BuildFakeValidIngredientStateIngredient()

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("CreateValidIngredientStateIngredient", testutils.ContextMatcher, testutils.MatchType[*types.ValidIngredientStateIngredientCreationRequestInput]()).Return(exampleValidIngredientStateIngredient, nil)
		s.validEnumerationsManager = mvem

		exampleInput := fakes.BuildFake[messages.CreateValidIngredientStateIngredientRequest]()

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

		exampleValidInstrument := fakes.BuildFakeValidInstrument()

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("CreateValidInstrument", testutils.ContextMatcher, testutils.MatchType[*types.ValidInstrumentCreationRequestInput]()).Return(exampleValidInstrument, nil)
		s.validEnumerationsManager = mvem

		exampleInput := fakes.BuildFake[messages.CreateValidInstrumentRequest]()

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

		exampleValidMeasurementUnit := fakes.BuildFakeValidMeasurementUnit()

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("CreateValidMeasurementUnit", testutils.ContextMatcher, testutils.MatchType[*types.ValidMeasurementUnitCreationRequestInput]()).Return(exampleValidMeasurementUnit, nil)
		s.validEnumerationsManager = mvem

		exampleInput := fakes.BuildFake[messages.CreateValidMeasurementUnitRequest]()

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

		exampleValidMeasurementUnitConversion := fakes.BuildFakeValidMeasurementUnitConversion()

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("CreateValidMeasurementUnitConversion", testutils.ContextMatcher, testutils.MatchType[*types.ValidMeasurementUnitConversionCreationRequestInput]()).Return(exampleValidMeasurementUnitConversion, nil)
		s.validEnumerationsManager = mvem

		exampleInput := fakes.BuildFake[messages.CreateValidMeasurementUnitConversionRequest]()

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

		exampleValidPreparation := fakes.BuildFakeValidPreparation()

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("CreateValidPreparation", testutils.ContextMatcher, testutils.MatchType[*types.ValidPreparationCreationRequestInput]()).Return(exampleValidPreparation, nil)
		s.validEnumerationsManager = mvem

		exampleInput := fakes.BuildFake[messages.CreateValidPreparationRequest]()

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

		exampleValidPreparationInstrument := fakes.BuildFakeValidPreparationInstrument()

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("CreateValidPreparationInstrument", testutils.ContextMatcher, testutils.MatchType[*types.ValidPreparationInstrumentCreationRequestInput]()).Return(exampleValidPreparationInstrument, nil)
		s.validEnumerationsManager = mvem

		exampleInput := fakes.BuildFake[messages.CreateValidPreparationInstrumentRequest]()

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

		exampleValidPreparationVessel := fakes.BuildFakeValidPreparationVessel()

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("CreateValidPreparationVessel", testutils.ContextMatcher, testutils.MatchType[*types.ValidPreparationVesselCreationRequestInput]()).Return(exampleValidPreparationVessel, nil)
		s.validEnumerationsManager = mvem

		exampleInput := fakes.BuildFake[messages.CreateValidPreparationVesselRequest]()

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

		exampleValidVessel := fakes.BuildFakeValidVessel()

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("CreateValidVessel", testutils.ContextMatcher, testutils.MatchType[*types.ValidVesselCreationRequestInput]()).Return(exampleValidVessel, nil)
		s.validEnumerationsManager = mvem

		exampleInput := fakes.BuildFake[messages.CreateValidVesselRequest]()

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

	})
}

func TestServiceImpl_GetRandomValidInstrument(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

	})
}

func TestServiceImpl_GetRandomValidPreparation(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

	})
}

func TestServiceImpl_GetRandomValidVessel(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

	})
}

func TestServiceImpl_GetValidIngredient(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

	})
}

func TestServiceImpl_GetValidIngredientGroup(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

	})
}

func TestServiceImpl_GetValidIngredientGroups(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

	})
}

func TestServiceImpl_GetValidIngredientMeasurementUnit(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

	})
}

func TestServiceImpl_GetValidIngredientMeasurementUnits(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

	})
}

func TestServiceImpl_GetValidIngredientMeasurementUnitsByIngredient(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

	})
}

func TestServiceImpl_GetValidIngredientMeasurementUnitsByMeasurementUnit(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

	})
}

func TestServiceImpl_GetValidIngredientPreparation(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

	})
}

func TestServiceImpl_GetValidIngredientPreparations(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

	})
}

func TestServiceImpl_GetValidIngredientPreparationsByIngredient(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

	})
}

func TestServiceImpl_GetValidIngredientPreparationsByPreparation(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

	})
}

func TestServiceImpl_GetValidIngredientState(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

	})
}

func TestServiceImpl_GetValidIngredientStateIngredient(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

	})
}

func TestServiceImpl_GetValidIngredientStateIngredients(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

	})
}

func TestServiceImpl_GetValidIngredientStateIngredientsByIngredient(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

	})
}

func TestServiceImpl_GetValidIngredientStateIngredientsByIngredientState(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

	})
}

func TestServiceImpl_GetValidIngredientStates(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

	})
}

func TestServiceImpl_GetValidIngredients(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

	})
}

func TestServiceImpl_GetValidInstrument(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

	})
}

func TestServiceImpl_GetValidInstruments(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

	})
}

func TestServiceImpl_GetValidMeasurementUnit(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

	})
}

func TestServiceImpl_GetValidMeasurementUnitConversion(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

	})
}

func TestServiceImpl_GetValidMeasurementUnitConversionsFromUnit(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

	})
}

func TestServiceImpl_GetValidMeasurementUnitConversionsToUnit(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

	})
}

func TestServiceImpl_GetValidMeasurementUnits(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

	})
}

func TestServiceImpl_GetValidPreparation(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

	})
}

func TestServiceImpl_GetValidPreparationInstrument(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

	})
}

func TestServiceImpl_GetValidPreparationInstruments(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

	})
}

func TestServiceImpl_GetValidPreparationInstrumentsByInstrument(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

	})
}

func TestServiceImpl_GetValidPreparationInstrumentsByPreparation(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

	})
}

func TestServiceImpl_GetValidPreparationVessel(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

	})
}

func TestServiceImpl_GetValidPreparationVessels(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

	})
}

func TestServiceImpl_GetValidPreparationVesselsByPreparation(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

	})
}

func TestServiceImpl_GetValidPreparationVesselsByVessel(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

	})
}

func TestServiceImpl_GetValidPreparations(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

	})
}

func TestServiceImpl_GetValidVessel(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

	})
}

func TestServiceImpl_GetValidVessels(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

	})
}

func TestServiceImpl_SearchForValidIngredientGroups(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

	})
}

func TestServiceImpl_SearchForValidIngredientStates(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

	})
}

func TestServiceImpl_SearchForValidIngredients(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

	})
}

func TestServiceImpl_SearchForValidInstruments(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

	})
}

func TestServiceImpl_SearchForValidMeasurementUnits(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

	})
}

func TestServiceImpl_SearchForValidPreparations(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

	})
}

func TestServiceImpl_SearchForValidVessels(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

	})
}

func TestServiceImpl_SearchValidIngredientsByPreparation(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

	})
}

func TestServiceImpl_SearchValidMeasurementUnitsByIngredient(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

	})
}

func TestServiceImpl_UpdateValidIngredient(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

	})
}

func TestServiceImpl_UpdateValidIngredientGroup(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

	})
}

func TestServiceImpl_UpdateValidIngredientMeasurementUnit(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

	})
}

func TestServiceImpl_UpdateValidIngredientPreparation(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

	})
}

func TestServiceImpl_UpdateValidIngredientState(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

	})
}

func TestServiceImpl_UpdateValidIngredientStateIngredient(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

	})
}

func TestServiceImpl_UpdateValidInstrument(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

	})
}

func TestServiceImpl_UpdateValidMeasurementUnit(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

	})
}

func TestServiceImpl_UpdateValidMeasurementUnitConversion(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

	})
}

func TestServiceImpl_UpdateValidPreparation(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

	})
}

func TestServiceImpl_UpdateValidPreparationInstrument(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

	})
}

func TestServiceImpl_UpdateValidPreparationVessel(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

	})
}

func TestServiceImpl_UpdateValidVessel(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

	})
}
