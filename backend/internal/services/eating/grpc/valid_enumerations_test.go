package grpc

import (
	"testing"

	"github.com/dinnerdonebetter/backend/internal/lib/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/tracing"

	"github.com/stretchr/testify/assert"
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

		assert.NotNil(t, ctx)
		assert.NotNil(t, s)

		// s.ArchiveValidIngredient(ctx)
	})
}

func TestServiceImpl_ArchiveValidIngredientGroup(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		assert.NotNil(t, ctx)
		assert.NotNil(t, s)

		// s.ArchiveValidIngredientGroup(ctx)
	})
}

func TestServiceImpl_ArchiveValidIngredientMeasurementUnit(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		assert.NotNil(t, ctx)
		assert.NotNil(t, s)

		// s.ArchiveValidIngredientMeasurementUnit(ctx)
	})
}

func TestServiceImpl_ArchiveValidIngredientPreparation(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		assert.NotNil(t, ctx)
		assert.NotNil(t, s)

		// s.ArchiveValidIngredientPreparation(ctx)
	})
}

func TestServiceImpl_ArchiveValidIngredientState(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		assert.NotNil(t, ctx)
		assert.NotNil(t, s)

		// s.ArchiveValidIngredientState(ctx)
	})
}

func TestServiceImpl_ArchiveValidIngredientStateIngredient(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		assert.NotNil(t, ctx)
		assert.NotNil(t, s)

		// s.ArchiveValidIngredientStateIngredient(ctx)
	})
}

func TestServiceImpl_ArchiveValidInstrument(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		assert.NotNil(t, ctx)
		assert.NotNil(t, s)

		// s.ArchiveValidInstrument(ctx)
	})
}

func TestServiceImpl_ArchiveValidMeasurementUnit(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		assert.NotNil(t, ctx)
		assert.NotNil(t, s)

		// s.ArchiveValidMeasurementUnit(ctx)
	})
}

func TestServiceImpl_ArchiveValidMeasurementUnitConversion(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		assert.NotNil(t, ctx)
		assert.NotNil(t, s)

		// s.ArchiveValidMeasurementUnitConversion(ctx)
	})
}

func TestServiceImpl_ArchiveValidPreparation(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		assert.NotNil(t, ctx)
		assert.NotNil(t, s)

		// s.ArchiveValidPreparation(ctx)
	})
}

func TestServiceImpl_ArchiveValidPreparationInstrument(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		assert.NotNil(t, ctx)
		assert.NotNil(t, s)

		// s.ArchiveValidPreparationInstrument(ctx)
	})
}

func TestServiceImpl_ArchiveValidPreparationVessel(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		assert.NotNil(t, ctx)
		assert.NotNil(t, s)

		// s.ArchiveValidPreparationVessel(ctx)
	})
}

func TestServiceImpl_ArchiveValidVessel(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		assert.NotNil(t, ctx)
		assert.NotNil(t, s)

		// s.ArchiveValidVessel(ctx)
	})
}

func TestServiceImpl_CreateValidIngredient(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		assert.NotNil(t, ctx)
		assert.NotNil(t, s)

		// s.CreateValidIngredient(ctx)
	})
}

func TestServiceImpl_CreateValidIngredientGroup(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		assert.NotNil(t, ctx)
		assert.NotNil(t, s)

		// s.CreateValidIngredientGroup(ctx)
	})
}

func TestServiceImpl_CreateValidIngredientMeasurementUnit(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		assert.NotNil(t, ctx)
		assert.NotNil(t, s)

		// s.CreateValidIngredientMeasurementUnit(ctx)
	})
}

func TestServiceImpl_CreateValidIngredientPreparation(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		assert.NotNil(t, ctx)
		assert.NotNil(t, s)

		// s.CreateValidIngredientPreparation(ctx)
	})
}

func TestServiceImpl_CreateValidIngredientState(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		assert.NotNil(t, ctx)
		assert.NotNil(t, s)

		// s.CreateValidIngredientState(ctx)
	})
}

func TestServiceImpl_CreateValidIngredientStateIngredient(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		assert.NotNil(t, ctx)
		assert.NotNil(t, s)

		// s.CreateValidIngredientStateIngredient(ctx)
	})
}

func TestServiceImpl_CreateValidInstrument(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		assert.NotNil(t, ctx)
		assert.NotNil(t, s)

		// s.CreateValidInstrument(ctx)
	})
}

func TestServiceImpl_CreateValidMeasurementUnit(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		assert.NotNil(t, ctx)
		assert.NotNil(t, s)

		// s.CreateValidMeasurementUnit(ctx)
	})
}

func TestServiceImpl_CreateValidMeasurementUnitConversion(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		assert.NotNil(t, ctx)
		assert.NotNil(t, s)

		// s.CreateValidMeasurementUnitConversion(ctx)
	})
}

func TestServiceImpl_CreateValidPreparation(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		assert.NotNil(t, ctx)
		assert.NotNil(t, s)

		// s.CreateValidPreparation(ctx)
	})
}

func TestServiceImpl_CreateValidPreparationInstrument(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		assert.NotNil(t, ctx)
		assert.NotNil(t, s)

		// s.CreateValidPreparationInstrument(ctx)
	})
}

func TestServiceImpl_CreateValidPreparationVessel(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		assert.NotNil(t, ctx)
		assert.NotNil(t, s)

		// s.CreateValidPreparationVessel(ctx)
	})
}

func TestServiceImpl_CreateValidVessel(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		assert.NotNil(t, ctx)
		assert.NotNil(t, s)

		// s.CreateValidVessel(ctx)
	})
}

func TestServiceImpl_GetRandomValidIngredient(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		assert.NotNil(t, ctx)
		assert.NotNil(t, s)

		// s.GetRandomValidIngredient(ctx)
	})
}

func TestServiceImpl_GetRandomValidInstrument(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		assert.NotNil(t, ctx)
		assert.NotNil(t, s)

		// s.GetRandomValidInstrument(ctx)
	})
}

func TestServiceImpl_GetRandomValidPreparation(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		assert.NotNil(t, ctx)
		assert.NotNil(t, s)

		// s.GetRandomValidPreparation(ctx)
	})
}

func TestServiceImpl_GetRandomValidVessel(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		assert.NotNil(t, ctx)
		assert.NotNil(t, s)

		// s.GetRandomValidVessel(ctx)
	})
}

func TestServiceImpl_GetValidIngredient(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		assert.NotNil(t, ctx)
		assert.NotNil(t, s)

		// s.GetValidIngredient(ctx)
	})
}

func TestServiceImpl_GetValidIngredientGroup(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		assert.NotNil(t, ctx)
		assert.NotNil(t, s)

		// s.GetValidIngredientGroup(ctx)
	})
}

func TestServiceImpl_GetValidIngredientGroups(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		assert.NotNil(t, ctx)
		assert.NotNil(t, s)

		// s.GetValidIngredientGroups(ctx)
	})
}

func TestServiceImpl_GetValidIngredientMeasurementUnit(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		assert.NotNil(t, ctx)
		assert.NotNil(t, s)

		// s.GetValidIngredientMeasurementUnit(ctx)
	})
}

func TestServiceImpl_GetValidIngredientMeasurementUnits(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		assert.NotNil(t, ctx)
		assert.NotNil(t, s)

		// s.GetValidIngredientMeasurementUnits(ctx)
	})
}

func TestServiceImpl_GetValidIngredientMeasurementUnitsByIngredient(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		assert.NotNil(t, ctx)
		assert.NotNil(t, s)

		// s.GetValidIngredientMeasurementUnitsByIngredient(ctx)
	})
}

func TestServiceImpl_GetValidIngredientMeasurementUnitsByMeasurementUnit(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		assert.NotNil(t, ctx)
		assert.NotNil(t, s)

		// s.GetValidIngredientMeasurementUnitsByMeasurementUnit(ctx)
	})
}

func TestServiceImpl_GetValidIngredientPreparation(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		assert.NotNil(t, ctx)
		assert.NotNil(t, s)

		// s.GetValidIngredientPreparation(ctx)
	})
}

func TestServiceImpl_GetValidIngredientPreparations(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		assert.NotNil(t, ctx)
		assert.NotNil(t, s)

		// s.GetValidIngredientPreparations(ctx)
	})
}

func TestServiceImpl_GetValidIngredientPreparationsByIngredient(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		assert.NotNil(t, ctx)
		assert.NotNil(t, s)

		// s.GetValidIngredientPreparationsByIngredient(ctx)
	})
}

func TestServiceImpl_GetValidIngredientPreparationsByPreparation(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		assert.NotNil(t, ctx)
		assert.NotNil(t, s)

		// s.GetValidIngredientPreparationsByPreparation(ctx)
	})
}

func TestServiceImpl_GetValidIngredientState(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		assert.NotNil(t, ctx)
		assert.NotNil(t, s)

		// s.GetValidIngredientState(ctx)
	})
}

func TestServiceImpl_GetValidIngredientStateIngredient(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		assert.NotNil(t, ctx)
		assert.NotNil(t, s)

		// s.GetValidIngredientStateIngredient(ctx)
	})
}

func TestServiceImpl_GetValidIngredientStateIngredients(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		assert.NotNil(t, ctx)
		assert.NotNil(t, s)

		// s.GetValidIngredientStateIngredients(ctx)
	})
}

func TestServiceImpl_GetValidIngredientStateIngredientsByIngredient(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		assert.NotNil(t, ctx)
		assert.NotNil(t, s)

		// s.GetValidIngredientStateIngredientsByIngredient(ctx)
	})
}

func TestServiceImpl_GetValidIngredientStateIngredientsByIngredientState(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		assert.NotNil(t, ctx)
		assert.NotNil(t, s)

		// s.GetValidIngredientStateIngredientsByIngredientState(ctx)
	})
}

func TestServiceImpl_GetValidIngredientStates(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		assert.NotNil(t, ctx)
		assert.NotNil(t, s)

		// s.GetValidIngredientStates(ctx)
	})
}

func TestServiceImpl_GetValidIngredients(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		assert.NotNil(t, ctx)
		assert.NotNil(t, s)

		// s.GetValidIngredients(ctx)
	})
}

func TestServiceImpl_GetValidInstrument(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		assert.NotNil(t, ctx)
		assert.NotNil(t, s)

		// s.GetValidInstrument(ctx)
	})
}

func TestServiceImpl_GetValidInstruments(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		assert.NotNil(t, ctx)
		assert.NotNil(t, s)

		// s.GetValidInstruments(ctx)
	})
}

func TestServiceImpl_GetValidMeasurementUnit(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		assert.NotNil(t, ctx)
		assert.NotNil(t, s)

		// s.GetValidMeasurementUnit(ctx)
	})
}

func TestServiceImpl_GetValidMeasurementUnitConversion(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		assert.NotNil(t, ctx)
		assert.NotNil(t, s)

		// s.GetValidMeasurementUnitConversion(ctx)
	})
}

func TestServiceImpl_GetValidMeasurementUnitConversionsFromUnit(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		assert.NotNil(t, ctx)
		assert.NotNil(t, s)

		// s.GetValidMeasurementUnitConversionsFromUnit(ctx)
	})
}

func TestServiceImpl_GetValidMeasurementUnitConversionsToUnit(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		assert.NotNil(t, ctx)
		assert.NotNil(t, s)

		// s.GetValidMeasurementUnitConversionsToUnit(ctx)
	})
}

func TestServiceImpl_GetValidMeasurementUnits(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		assert.NotNil(t, ctx)
		assert.NotNil(t, s)

		// s.GetValidMeasurementUnits(ctx)
	})
}

func TestServiceImpl_GetValidPreparation(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		assert.NotNil(t, ctx)
		assert.NotNil(t, s)

		// s.GetValidPreparation(ctx)
	})
}

func TestServiceImpl_GetValidPreparationInstrument(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		assert.NotNil(t, ctx)
		assert.NotNil(t, s)

		// s.GetValidPreparationInstrument(ctx)
	})
}

func TestServiceImpl_GetValidPreparationInstruments(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		assert.NotNil(t, ctx)
		assert.NotNil(t, s)

		// s.GetValidPreparationInstruments(ctx)
	})
}

func TestServiceImpl_GetValidPreparationInstrumentsByInstrument(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		assert.NotNil(t, ctx)
		assert.NotNil(t, s)

		// s.GetValidPreparationInstrumentsByInstrument(ctx)
	})
}

func TestServiceImpl_GetValidPreparationInstrumentsByPreparation(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		assert.NotNil(t, ctx)
		assert.NotNil(t, s)

		// s.GetValidPreparationInstrumentsByPreparation(ctx)
	})
}

func TestServiceImpl_GetValidPreparationVessel(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		assert.NotNil(t, ctx)
		assert.NotNil(t, s)

		// s.GetValidPreparationVessel(ctx)
	})
}

func TestServiceImpl_GetValidPreparationVessels(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		assert.NotNil(t, ctx)
		assert.NotNil(t, s)

		// s.GetValidPreparationVessels(ctx)
	})
}

func TestServiceImpl_GetValidPreparationVesselsByPreparation(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		assert.NotNil(t, ctx)
		assert.NotNil(t, s)

		// s.GetValidPreparationVesselsByPreparation(ctx)
	})
}

func TestServiceImpl_GetValidPreparationVesselsByVessel(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		assert.NotNil(t, ctx)
		assert.NotNil(t, s)

		// s.GetValidPreparationVesselsByVessel(ctx)
	})
}

func TestServiceImpl_GetValidPreparations(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		assert.NotNil(t, ctx)
		assert.NotNil(t, s)

		// s.GetValidPreparations(ctx)
	})
}

func TestServiceImpl_GetValidVessel(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		assert.NotNil(t, ctx)
		assert.NotNil(t, s)

		// s.GetValidVessel(ctx)
	})
}

func TestServiceImpl_GetValidVessels(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		assert.NotNil(t, ctx)
		assert.NotNil(t, s)

		// s.GetValidVessels(ctx)
	})
}

func TestServiceImpl_SearchForValidIngredientGroups(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		assert.NotNil(t, ctx)
		assert.NotNil(t, s)

		// s.SearchForValidIngredientGroups(ctx)
	})
}

func TestServiceImpl_SearchForValidIngredientStates(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		assert.NotNil(t, ctx)
		assert.NotNil(t, s)

		// s.SearchForValidIngredientStates(ctx)
	})
}

func TestServiceImpl_SearchForValidIngredients(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		assert.NotNil(t, ctx)
		assert.NotNil(t, s)

		// s.SearchForValidIngredients(ctx)
	})
}

func TestServiceImpl_SearchForValidInstruments(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		assert.NotNil(t, ctx)
		assert.NotNil(t, s)

		// s.SearchForValidInstruments(ctx)
	})
}

func TestServiceImpl_SearchForValidMeasurementUnits(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		assert.NotNil(t, ctx)
		assert.NotNil(t, s)

		// s.SearchForValidMeasurementUnits(ctx)
	})
}

func TestServiceImpl_SearchForValidPreparations(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		assert.NotNil(t, ctx)
		assert.NotNil(t, s)

		// s.SearchForValidPreparations(ctx)
	})
}

func TestServiceImpl_SearchForValidVessels(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		assert.NotNil(t, ctx)
		assert.NotNil(t, s)

		// s.SearchForValidVessels(ctx)
	})
}

func TestServiceImpl_SearchValidIngredientsByPreparation(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		assert.NotNil(t, ctx)
		assert.NotNil(t, s)

		// s.SearchValidIngredientsByPreparation(ctx)
	})
}

func TestServiceImpl_SearchValidMeasurementUnitsByIngredient(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		assert.NotNil(t, ctx)
		assert.NotNil(t, s)

		// s.SearchValidMeasurementUnitsByIngredient(ctx)
	})
}

func TestServiceImpl_UpdateValidIngredient(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		assert.NotNil(t, ctx)
		assert.NotNil(t, s)

		// s.UpdateValidIngredient(ctx)
	})
}

func TestServiceImpl_UpdateValidIngredientGroup(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		assert.NotNil(t, ctx)
		assert.NotNil(t, s)

		// s.UpdateValidIngredientGroup(ctx)
	})
}

func TestServiceImpl_UpdateValidIngredientMeasurementUnit(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		assert.NotNil(t, ctx)
		assert.NotNil(t, s)

		// s.UpdateValidIngredientMeasurementUnit(ctx)
	})
}

func TestServiceImpl_UpdateValidIngredientPreparation(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		assert.NotNil(t, ctx)
		assert.NotNil(t, s)

		// s.UpdateValidIngredientPreparation(ctx)
	})
}

func TestServiceImpl_UpdateValidIngredientState(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		assert.NotNil(t, ctx)
		assert.NotNil(t, s)

		// s.UpdateValidIngredientState(ctx)
	})
}

func TestServiceImpl_UpdateValidIngredientStateIngredient(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		assert.NotNil(t, ctx)
		assert.NotNil(t, s)

		// s.UpdateValidIngredientStateIngredient(ctx)
	})
}

func TestServiceImpl_UpdateValidInstrument(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		assert.NotNil(t, ctx)
		assert.NotNil(t, s)

		// s.UpdateValidInstrument(ctx)
	})
}

func TestServiceImpl_UpdateValidMeasurementUnit(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		assert.NotNil(t, ctx)
		assert.NotNil(t, s)

		// s.UpdateValidMeasurementUnit(ctx)
	})
}

func TestServiceImpl_UpdateValidMeasurementUnitConversion(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		assert.NotNil(t, ctx)
		assert.NotNil(t, s)

		// s.UpdateValidMeasurementUnitConversion(ctx)
	})
}

func TestServiceImpl_UpdateValidPreparation(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		assert.NotNil(t, ctx)
		assert.NotNil(t, s)

		// s.UpdateValidPreparation(ctx)
	})
}

func TestServiceImpl_UpdateValidPreparationInstrument(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		assert.NotNil(t, ctx)
		assert.NotNil(t, s)

		// s.UpdateValidPreparationInstrument(ctx)
	})
}

func TestServiceImpl_UpdateValidPreparationVessel(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		assert.NotNil(t, ctx)
		assert.NotNil(t, s)

		// s.UpdateValidPreparationVessel(ctx)
	})
}

func TestServiceImpl_UpdateValidVessel(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForTest(t)

		assert.NotNil(t, ctx)
		assert.NotNil(t, s)

		// s.UpdateValidVessel(ctx)
	})
}
