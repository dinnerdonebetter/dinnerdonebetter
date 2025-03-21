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

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		exampleValidIngredientMeasurementUnitID := fakes.BuildFakeID()

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("ArchiveValidIngredientMeasurementUnit", testutils.ContextMatcher, exampleValidIngredientMeasurementUnitID).Return(nil)
		s.validEnumerationsManager = mvem

		res, err := s.ArchiveValidIngredientMeasurementUnit(ctx, &messages.ArchiveValidIngredientMeasurementUnitRequest{ValidIngredientMeasurementUnitID: exampleValidIngredientMeasurementUnitID})
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

		exampleResult := fakes.BuildFakeValidIngredient()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("RandomValidIngredient", testutils.ContextMatcher).Return(exampleResult, nil)
		s.validEnumerationsManager = mvem

		result, err := s.GetRandomValidIngredient(ctx, &messages.GetRandomValidIngredientRequest{})
		assert.Equal(t, exampleResult.ID, result.Result.ID)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mvem)
	})
}

func TestServiceImpl_GetRandomValidInstrument(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleResult := fakes.BuildFakeValidInstrument()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("RandomValidInstrument", testutils.ContextMatcher).Return(exampleResult, nil)
		s.validEnumerationsManager = mvem

		result, err := s.GetRandomValidInstrument(ctx, &messages.GetRandomValidInstrumentRequest{})
		assert.Equal(t, exampleResult.ID, result.Result.ID)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mvem)
	})
}

func TestServiceImpl_GetRandomValidPreparation(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleResult := fakes.BuildFakeValidPreparation()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("RandomValidPreparation", testutils.ContextMatcher).Return(exampleResult, nil)
		s.validEnumerationsManager = mvem

		result, err := s.GetRandomValidPreparation(ctx, &messages.GetRandomValidPreparationRequest{})
		assert.Equal(t, exampleResult.ID, result.Result.ID)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mvem)
	})
}

func TestServiceImpl_GetRandomValidVessel(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleResult := fakes.BuildFakeValidVessel()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("RandomValidVessel", testutils.ContextMatcher).Return(exampleResult, nil)
		s.validEnumerationsManager = mvem

		result, err := s.GetRandomValidVessel(ctx, &messages.GetRandomValidVesselRequest{})
		assert.Equal(t, exampleResult.ID, result.Result.ID)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mvem)
	})
}

func TestServiceImpl_GetValidIngredient(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleResult := fakes.BuildFakeValidIngredient()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("ReadValidIngredient", testutils.ContextMatcher, exampleResult.ID).Return(exampleResult, nil)
		s.validEnumerationsManager = mvem

		result, err := s.GetValidIngredient(ctx, &messages.GetValidIngredientRequest{ValidIngredientID: exampleResult.ID})
		assert.Equal(t, exampleResult.ID, result.Result.ID)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mvem)
	})
}

func TestServiceImpl_GetValidIngredientGroup(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleResult := fakes.BuildFakeValidIngredientGroup()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("ReadValidIngredientGroup", testutils.ContextMatcher, exampleResult.ID).Return(exampleResult, nil)
		s.validEnumerationsManager = mvem

		result, err := s.GetValidIngredientGroup(ctx, &messages.GetValidIngredientGroupRequest{ValidIngredientGroupID: exampleResult.ID})
		assert.Equal(t, exampleResult.ID, result.Result.ID)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mvem)
	})
}

func TestServiceImpl_GetValidIngredientGroups(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleResult := fakes.BuildFakeValidIngredientGroupsList()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("ListValidIngredientGroups", testutils.ContextMatcher, testutils.QueryFilterMatcher).Return(exampleResult.Data, "", nil)
		s.validEnumerationsManager = mvem

		result, err := s.GetValidIngredientGroups(ctx, &messages.GetValidIngredientGroupsRequest{})
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

		exampleResult := fakes.BuildFakeValidIngredientMeasurementUnit()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("ReadValidIngredientMeasurementUnit", testutils.ContextMatcher, exampleResult.ID).Return(exampleResult, nil)
		s.validEnumerationsManager = mvem

		result, err := s.GetValidIngredientMeasurementUnit(ctx, &messages.GetValidIngredientMeasurementUnitRequest{ValidIngredientMeasurementUnitID: exampleResult.ID})
		assert.Equal(t, exampleResult.ID, result.Result.ID)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mvem)
	})
}

func TestServiceImpl_GetValidIngredientMeasurementUnits(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleResult := fakes.BuildFakeValidMeasurementUnitsList()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("ListValidMeasurementUnits", testutils.ContextMatcher, testutils.QueryFilterMatcher).Return(exampleResult.Data, "", nil)
		s.validEnumerationsManager = mvem

		result, err := s.GetValidMeasurementUnits(ctx, &messages.GetValidMeasurementUnitsRequest{})
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

		exampleID := fakes.BuildFakeID()
		exampleResult := fakes.BuildFakeValidIngredientMeasurementUnitsList()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("SearchValidIngredientMeasurementUnitsByIngredient", testutils.ContextMatcher, exampleID, testutils.QueryFilterMatcher).Return(exampleResult.Data, nil)
		s.validEnumerationsManager = mvem

		result, err := s.GetValidIngredientMeasurementUnitsByIngredient(ctx, &messages.GetValidIngredientMeasurementUnitsByIngredientRequest{
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

		exampleID := fakes.BuildFakeID()
		exampleResult := fakes.BuildFakeValidIngredientMeasurementUnitsList()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("SearchValidIngredientMeasurementUnitsByMeasurementUnit", testutils.ContextMatcher, exampleID, testutils.QueryFilterMatcher).Return(exampleResult.Data, nil)
		s.validEnumerationsManager = mvem

		result, err := s.GetValidIngredientMeasurementUnitsByMeasurementUnit(ctx, &messages.GetValidIngredientMeasurementUnitsByMeasurementUnitRequest{
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

		exampleResult := fakes.BuildFakeValidIngredientPreparation()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("ReadValidIngredientPreparation", testutils.ContextMatcher, exampleResult.ID).Return(exampleResult, nil)
		s.validEnumerationsManager = mvem

		result, err := s.GetValidIngredientPreparation(ctx, &messages.GetValidIngredientPreparationRequest{ValidIngredientPreparationID: exampleResult.ID})
		assert.Equal(t, exampleResult.ID, result.Result.ID)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mvem)
	})
}

func TestServiceImpl_GetValidIngredientPreparations(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleResult := fakes.BuildFakeValidIngredientPreparationsList()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("ListValidIngredientPreparations", testutils.ContextMatcher, testutils.QueryFilterMatcher).Return(exampleResult.Data, "", nil)
		s.validEnumerationsManager = mvem

		result, err := s.GetValidIngredientPreparations(ctx, &messages.GetValidIngredientPreparationsRequest{})
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

		exampleID := fakes.BuildFakeID()
		exampleResult := fakes.BuildFakeValidIngredientPreparationsList()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("SearchValidIngredientPreparationsByIngredient", testutils.ContextMatcher, exampleID, testutils.QueryFilterMatcher).Return(exampleResult.Data, nil)
		s.validEnumerationsManager = mvem

		result, err := s.GetValidIngredientPreparationsByIngredient(ctx, &messages.GetValidIngredientPreparationsByIngredientRequest{
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

		exampleID := fakes.BuildFakeID()
		exampleResult := fakes.BuildFakeValidIngredientPreparationsList()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("SearchValidIngredientPreparationsByPreparation", testutils.ContextMatcher, exampleID, testutils.QueryFilterMatcher).Return(exampleResult.Data, nil)
		s.validEnumerationsManager = mvem

		result, err := s.GetValidIngredientPreparationsByPreparation(ctx, &messages.GetValidIngredientPreparationsByPreparationRequest{
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

		exampleResult := fakes.BuildFakeValidIngredientState()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("ReadValidIngredientState", testutils.ContextMatcher, exampleResult.ID).Return(exampleResult, nil)
		s.validEnumerationsManager = mvem

		result, err := s.GetValidIngredientState(ctx, &messages.GetValidIngredientStateRequest{ValidIngredientStateID: exampleResult.ID})
		assert.Equal(t, exampleResult.ID, result.Result.ID)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mvem)
	})
}

func TestServiceImpl_GetValidIngredientStateIngredient(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleResult := fakes.BuildFakeValidIngredientStateIngredient()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("ReadValidIngredientStateIngredient", testutils.ContextMatcher, exampleResult.ID).Return(exampleResult, nil)
		s.validEnumerationsManager = mvem

		result, err := s.GetValidIngredientStateIngredient(ctx, &messages.GetValidIngredientStateIngredientRequest{ValidIngredientStateIngredientID: exampleResult.ID})
		assert.Equal(t, exampleResult.ID, result.Result.ID)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mvem)
	})
}

func TestServiceImpl_GetValidIngredientStateIngredients(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleResult := fakes.BuildFakeValidIngredientStateIngredientsList()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("ListValidIngredientStateIngredients", testutils.ContextMatcher, testutils.QueryFilterMatcher).Return(exampleResult.Data, "", nil)
		s.validEnumerationsManager = mvem

		result, err := s.GetValidIngredientStateIngredients(ctx, &messages.GetValidIngredientStateIngredientsRequest{})
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

		exampleID := fakes.BuildFakeID()
		exampleResult := fakes.BuildFakeValidIngredientStateIngredientsList()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("SearchValidIngredientStateIngredientsByIngredient", testutils.ContextMatcher, exampleID, testutils.QueryFilterMatcher).Return(exampleResult.Data, nil)
		s.validEnumerationsManager = mvem

		result, err := s.GetValidIngredientStateIngredientsByIngredient(ctx, &messages.GetValidIngredientStateIngredientsByIngredientRequest{
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

		exampleID := fakes.BuildFakeID()
		exampleResult := fakes.BuildFakeValidIngredientStateIngredientsList()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("SearchValidIngredientStateIngredientsByIngredientState", testutils.ContextMatcher, exampleID, testutils.QueryFilterMatcher).Return(exampleResult.Data, nil)
		s.validEnumerationsManager = mvem

		result, err := s.GetValidIngredientStateIngredientsByIngredientState(ctx, &messages.GetValidIngredientStateIngredientsByIngredientStateRequest{
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

		exampleResult := fakes.BuildFakeValidIngredientStatesList()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("ListValidIngredientStates", testutils.ContextMatcher, testutils.QueryFilterMatcher).Return(exampleResult.Data, "", nil)
		s.validEnumerationsManager = mvem

		result, err := s.GetValidIngredientStates(ctx, &messages.GetValidIngredientStatesRequest{})
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

		exampleResult := fakes.BuildFakeValidIngredientsList()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("ListValidIngredients", testutils.ContextMatcher, testutils.QueryFilterMatcher).Return(exampleResult.Data, "", nil)
		s.validEnumerationsManager = mvem

		result, err := s.GetValidIngredients(ctx, &messages.GetValidIngredientsRequest{})
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

		exampleResult := fakes.BuildFakeValidInstrument()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("ReadValidInstrument", testutils.ContextMatcher, exampleResult.ID).Return(exampleResult, nil)
		s.validEnumerationsManager = mvem

		result, err := s.GetValidInstrument(ctx, &messages.GetValidInstrumentRequest{ValidInstrumentID: exampleResult.ID})
		assert.Equal(t, exampleResult.ID, result.Result.ID)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mvem)
	})
}

func TestServiceImpl_GetValidInstruments(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleResult := fakes.BuildFakeValidInstrumentsList()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("ListValidInstruments", testutils.ContextMatcher, testutils.QueryFilterMatcher).Return(exampleResult.Data, "", nil)
		s.validEnumerationsManager = mvem

		result, err := s.GetValidInstruments(ctx, &messages.GetValidInstrumentsRequest{})
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

		exampleResult := fakes.BuildFakeValidMeasurementUnit()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("ReadValidMeasurementUnit", testutils.ContextMatcher, exampleResult.ID).Return(exampleResult, nil)
		s.validEnumerationsManager = mvem

		result, err := s.GetValidMeasurementUnit(ctx, &messages.GetValidMeasurementUnitRequest{ValidMeasurementUnitID: exampleResult.ID})
		assert.Equal(t, exampleResult.ID, result.Result.ID)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mvem)
	})
}

func TestServiceImpl_GetValidMeasurementUnitConversion(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleResult := fakes.BuildFakeValidMeasurementUnitConversion()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("ReadValidMeasurementUnitConversion", testutils.ContextMatcher, exampleResult.ID).Return(exampleResult, nil)
		s.validEnumerationsManager = mvem

		result, err := s.GetValidMeasurementUnitConversion(ctx, &messages.GetValidMeasurementUnitConversionRequest{ValidMeasurementUnitConversionID: exampleResult.ID})
		assert.Equal(t, exampleResult.ID, result.Result.ID)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mvem)
	})
}

func TestServiceImpl_GetValidMeasurementUnitConversionsFromUnit(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleID := fakes.BuildFakeID()
		exampleResult := fakes.BuildFakeValidMeasurementUnitConversionsList()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("ValidMeasurementUnitConversionsFromMeasurementUnit", testutils.ContextMatcher, exampleID).Return(exampleResult.Data, nil)
		s.validEnumerationsManager = mvem

		result, err := s.GetValidMeasurementUnitConversionsFromUnit(ctx, &messages.GetValidMeasurementUnitConversionsFromUnitRequest{
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

		exampleID := fakes.BuildFakeID()
		exampleResult := fakes.BuildFakeValidMeasurementUnitConversionsList()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("ValidMeasurementUnitConversionsToMeasurementUnit", testutils.ContextMatcher, exampleID).Return(exampleResult.Data, nil)
		s.validEnumerationsManager = mvem

		result, err := s.GetValidMeasurementUnitConversionsToUnit(ctx, &messages.GetValidMeasurementUnitConversionsToUnitRequest{
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

		exampleResult := fakes.BuildFakeValidMeasurementUnitsList()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("ListValidMeasurementUnits", testutils.ContextMatcher, testutils.QueryFilterMatcher).Return(exampleResult.Data, "", nil)
		s.validEnumerationsManager = mvem

		result, err := s.GetValidMeasurementUnits(ctx, &messages.GetValidMeasurementUnitsRequest{})
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

		exampleResult := fakes.BuildFakeValidPreparation()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("ReadValidPreparation", testutils.ContextMatcher, exampleResult.ID).Return(exampleResult, nil)
		s.validEnumerationsManager = mvem

		result, err := s.GetValidPreparation(ctx, &messages.GetValidPreparationRequest{ValidPreparationID: exampleResult.ID})
		assert.Equal(t, exampleResult.ID, result.Result.ID)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mvem)
	})
}

func TestServiceImpl_GetValidPreparationInstrument(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleResult := fakes.BuildFakeValidPreparationInstrument()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("ReadValidPreparationInstrument", testutils.ContextMatcher, exampleResult.ID).Return(exampleResult, nil)
		s.validEnumerationsManager = mvem

		result, err := s.GetValidPreparationInstrument(ctx, &messages.GetValidPreparationInstrumentRequest{ValidPreparationInstrumentID: exampleResult.ID})
		assert.Equal(t, exampleResult.ID, result.Result.ID)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mvem)
	})
}

func TestServiceImpl_GetValidPreparationInstruments(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleResult := fakes.BuildFakeValidPreparationInstrumentsList()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("ListValidPreparationInstruments", testutils.ContextMatcher, testutils.QueryFilterMatcher).Return(exampleResult.Data, "", nil)
		s.validEnumerationsManager = mvem

		result, err := s.GetValidPreparationInstruments(ctx, &messages.GetValidPreparationInstrumentsRequest{})
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

		exampleID := fakes.BuildFakeID()
		exampleResult := fakes.BuildFakeValidPreparationInstrumentsList()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("SearchValidPreparationInstrumentsByInstrument", testutils.ContextMatcher, exampleID, testutils.QueryFilterMatcher).Return(exampleResult.Data, nil)
		s.validEnumerationsManager = mvem

		result, err := s.GetValidPreparationInstrumentsByInstrument(ctx, &messages.GetValidPreparationInstrumentsByInstrumentRequest{
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

		exampleID := fakes.BuildFakeID()
		exampleResult := fakes.BuildFakeValidPreparationInstrumentsList()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("SearchValidPreparationInstrumentsByPreparation", testutils.ContextMatcher, exampleID, testutils.QueryFilterMatcher).Return(exampleResult.Data, nil)
		s.validEnumerationsManager = mvem

		result, err := s.GetValidPreparationInstrumentsByPreparation(ctx, &messages.GetValidPreparationInstrumentsByPreparationRequest{
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

		exampleResult := fakes.BuildFakeValidPreparationVessel()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("ReadValidPreparationVessel", testutils.ContextMatcher, exampleResult.ID).Return(exampleResult, nil)
		s.validEnumerationsManager = mvem

		result, err := s.GetValidPreparationVessel(ctx, &messages.GetValidPreparationVesselRequest{ValidPreparationVesselID: exampleResult.ID})
		assert.Equal(t, exampleResult.ID, result.Result.ID)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mvem)
	})
}

func TestServiceImpl_GetValidPreparationVessels(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleResult := fakes.BuildFakeValidPreparationVesselsList()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("ListValidPreparationVessels", testutils.ContextMatcher, testutils.QueryFilterMatcher).Return(exampleResult.Data, "", nil)
		s.validEnumerationsManager = mvem

		result, err := s.GetValidPreparationVessels(ctx, &messages.GetValidPreparationVesselsRequest{})
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

		exampleID := fakes.BuildFakeID()
		exampleResult := fakes.BuildFakeValidPreparationVesselsList()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("SearchValidPreparationVesselsByPreparation", testutils.ContextMatcher, exampleID, testutils.QueryFilterMatcher).Return(exampleResult.Data, nil)
		s.validEnumerationsManager = mvem

		result, err := s.GetValidPreparationVesselsByPreparation(ctx, &messages.GetValidPreparationVesselsByPreparationRequest{
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

		exampleID := fakes.BuildFakeID()
		exampleResult := fakes.BuildFakeValidPreparationVesselsList()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("SearchValidPreparationVesselsByVessel", testutils.ContextMatcher, exampleID, testutils.QueryFilterMatcher).Return(exampleResult.Data, nil)
		s.validEnumerationsManager = mvem

		result, err := s.GetValidPreparationVesselsByVessel(ctx, &messages.GetValidPreparationVesselsByVesselRequest{
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

		exampleResult := fakes.BuildFakeValidPreparationsList()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("ListValidPreparations", testutils.ContextMatcher, testutils.QueryFilterMatcher).Return(exampleResult.Data, "", nil)
		s.validEnumerationsManager = mvem

		result, err := s.GetValidPreparations(ctx, &messages.GetValidPreparationsRequest{})
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

		exampleResult := fakes.BuildFakeValidVessel()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("ReadValidVessel", testutils.ContextMatcher, exampleResult.ID).Return(exampleResult, nil)
		s.validEnumerationsManager = mvem

		result, err := s.GetValidVessel(ctx, &messages.GetValidVesselRequest{ValidVesselID: exampleResult.ID})
		assert.Equal(t, exampleResult.ID, result.Result.ID)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mvem)
	})
}

func TestServiceImpl_GetValidVessels(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleResult := fakes.BuildFakeValidVesselsList()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		mvem := &mockmanagers.MockValidEnumerationsManager{}
		mvem.On("ListValidVessels", testutils.ContextMatcher, testutils.QueryFilterMatcher).Return(exampleResult.Data, "", nil)
		s.validEnumerationsManager = mvem

		result, err := s.GetValidVessels(ctx, &messages.GetValidVesselsRequest{})
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

		exampleResult := fakes.BuildFakeValidIngredientGroupsList()
		exampleRequest := fakes.BuildFake[messages.SearchForValidIngredientGroupsRequest]()

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

		exampleResult := fakes.BuildFakeValidIngredientStatesList()
		exampleRequest := fakes.BuildFake[messages.SearchForValidIngredientStatesRequest]()

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

		exampleResult := fakes.BuildFakeValidIngredientsList()
		exampleRequest := fakes.BuildFake[messages.SearchForValidIngredientsRequest]()

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

		exampleResult := fakes.BuildFakeValidInstrumentsList()
		exampleRequest := fakes.BuildFake[messages.SearchForValidInstrumentsRequest]()

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

		exampleResult := fakes.BuildFakeValidMeasurementUnitsList()
		exampleRequest := fakes.BuildFake[messages.SearchForValidMeasurementUnitsRequest]()

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

		exampleResult := fakes.BuildFakeValidPreparationsList()
		exampleRequest := fakes.BuildFake[messages.SearchForValidPreparationsRequest]()

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

		exampleResult := fakes.BuildFakeValidVesselsList()
		exampleRequest := fakes.BuildFake[messages.SearchForValidVesselsRequest]()

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

		exampleResult := fakes.BuildFakeValidIngredientsList()
		exampleRequest := fakes.BuildFake[messages.SearchValidIngredientsByPreparationRequest]()

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

		exampleResult := fakes.BuildFakeValidMeasurementUnitsList()
		exampleRequest := fakes.BuildFake[messages.SearchValidMeasurementUnitsByIngredientRequest]()

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
