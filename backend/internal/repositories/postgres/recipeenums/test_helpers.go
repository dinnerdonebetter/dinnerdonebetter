package recipeenums

import (
	"context"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/domain/recipeenums"
	"github.com/dinnerdonebetter/backend/internal/domain/recipeenums/converters"
	"github.com/dinnerdonebetter/backend/internal/domain/recipeenums/fakes"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func CreateValidInstrumentForTest(t *testing.T, ctx context.Context, exampleValidInstrument *recipeenums.ValidInstrument, dbc recipeenums.Repository) *recipeenums.ValidInstrument {
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

	return validInstrument
}

func CreateValidIngredientForTest(t *testing.T, ctx context.Context, exampleValidIngredient *recipeenums.ValidIngredient, dbc recipeenums.Repository) *recipeenums.ValidIngredient {
	t.Helper()

	// create
	if exampleValidIngredient == nil {
		exampleValidIngredient = fakes.BuildFakeValidIngredient()
	}
	dbInput := converters.ConvertValidIngredientToValidIngredientDatabaseCreationInput(exampleValidIngredient)

	created, err := dbc.CreateValidIngredient(ctx, dbInput)
	exampleValidIngredient.CreatedAt = created.CreatedAt
	assert.NoError(t, err)
	assert.Equal(t, exampleValidIngredient, created)

	validIngredient, err := dbc.GetValidIngredient(ctx, created.ID)
	exampleValidIngredient.CreatedAt = validIngredient.CreatedAt

	assert.NoError(t, err)
	assert.Equal(t, validIngredient, exampleValidIngredient)

	return validIngredient
}

func CreateValidVesselForTest(t *testing.T, ctx context.Context, exampleValidVessel *recipeenums.ValidVessel, dbc recipeenums.Repository) *recipeenums.ValidVessel {
	t.Helper()

	// create
	if exampleValidVessel == nil {
		createdValidMeasurementUnit := CreateValidMeasurementUnitForTest(t, ctx, nil, dbc)
		exampleValidVessel = fakes.BuildFakeValidVessel()
		exampleValidVessel.CapacityUnit = createdValidMeasurementUnit
	}
	dbInput := converters.ConvertValidVesselToValidVesselDatabaseCreationInput(exampleValidVessel)

	created, err := dbc.CreateValidVessel(ctx, dbInput)
	assert.NoError(t, err)
	require.NotNil(t, created)

	exampleValidVessel.CreatedAt = created.CreatedAt
	exampleValidVessel.CapacityUnit = &recipeenums.ValidMeasurementUnit{ID: exampleValidVessel.CapacityUnit.ID}
	assert.Equal(t, exampleValidVessel, created)

	validVessel, err := dbc.GetValidVessel(ctx, created.ID)
	exampleValidVessel.CreatedAt = validVessel.CreatedAt
	exampleValidVessel.CapacityUnit = validVessel.CapacityUnit

	assert.NoError(t, err)
	assert.Equal(t, validVessel, exampleValidVessel)

	return validVessel
}

func CreateValidPreparationForTest(t *testing.T, ctx context.Context, exampleValidPreparation *recipeenums.ValidPreparation, dbc recipeenums.Repository) *recipeenums.ValidPreparation {
	t.Helper()

	// create
	if exampleValidPreparation == nil {
		exampleValidPreparation = fakes.BuildFakeValidPreparation()
	}
	dbInput := converters.ConvertValidPreparationToValidPreparationDatabaseCreationInput(exampleValidPreparation)

	created, err := dbc.CreateValidPreparation(ctx, dbInput)
	exampleValidPreparation.CreatedAt = created.CreatedAt
	assert.NoError(t, err)
	assert.Equal(t, exampleValidPreparation, created)

	validPreparation, err := dbc.GetValidPreparation(ctx, created.ID)
	exampleValidPreparation.CreatedAt = validPreparation.CreatedAt

	assert.NoError(t, err)
	assert.Equal(t, validPreparation, exampleValidPreparation)

	return validPreparation
}

func CreateValidPreparationVesselForTest(t *testing.T, ctx context.Context, exampleValidPreparationVessel *recipeenums.ValidPreparationVessel, dbc recipeenums.Repository) *recipeenums.ValidPreparationVessel {
	t.Helper()

	// create
	if exampleValidPreparationVessel == nil {
		exampleValidVessel := CreateValidVesselForTest(t, ctx, nil, dbc)
		exampleValidPreparation := CreateValidPreparationForTest(t, ctx, nil, dbc)
		exampleValidPreparationVessel = fakes.BuildFakeValidPreparationVessel()
		exampleValidPreparationVessel.Vessel = *exampleValidVessel
		exampleValidPreparationVessel.Preparation = *exampleValidPreparation
	}

	dbInput := converters.ConvertValidPreparationVesselToValidPreparationVesselDatabaseCreationInput(exampleValidPreparationVessel)

	created, err := dbc.CreateValidPreparationVessel(ctx, dbInput)
	assert.NoError(t, err)
	require.NotNil(t, created)
	exampleValidPreparationVessel.CreatedAt = created.CreatedAt
	assert.Equal(t, exampleValidPreparationVessel, created)

	validPreparationVessel, err := dbc.GetValidPreparationVessel(ctx, created.ID)
	exampleValidPreparationVessel.CreatedAt = validPreparationVessel.CreatedAt
	exampleValidPreparationVessel.Preparation = validPreparationVessel.Preparation
	exampleValidPreparationVessel.Vessel = validPreparationVessel.Vessel

	assert.NoError(t, err)
	assert.Equal(t, validPreparationVessel, exampleValidPreparationVessel)

	return created
}

func CreateValidPreparationInstrumentForTest(t *testing.T, ctx context.Context, exampleValidPreparationInstrument *recipeenums.ValidPreparationInstrument, dbc recipeenums.Repository) *recipeenums.ValidPreparationInstrument {
	t.Helper()

	// create
	if exampleValidPreparationInstrument == nil {
		exampleValidInstrument := CreateValidInstrumentForTest(t, ctx, nil, dbc)
		exampleValidPreparation := CreateValidPreparationForTest(t, ctx, nil, dbc)
		exampleValidPreparationInstrument = fakes.BuildFakeValidPreparationInstrument()
		exampleValidPreparationInstrument.Instrument = *exampleValidInstrument
		exampleValidPreparationInstrument.Preparation = *exampleValidPreparation
	}

	dbInput := converters.ConvertValidPreparationInstrumentToValidPreparationInstrumentDatabaseCreationInput(exampleValidPreparationInstrument)

	created, err := dbc.CreateValidPreparationInstrument(ctx, dbInput)
	assert.NoError(t, err)
	require.NotNil(t, created)
	exampleValidPreparationInstrument.CreatedAt = created.CreatedAt
	assert.Equal(t, exampleValidPreparationInstrument, created)

	validPreparationInstrument, err := dbc.GetValidPreparationInstrument(ctx, created.ID)
	exampleValidPreparationInstrument.CreatedAt = validPreparationInstrument.CreatedAt
	exampleValidPreparationInstrument.Preparation = validPreparationInstrument.Preparation
	exampleValidPreparationInstrument.Instrument = validPreparationInstrument.Instrument

	assert.NoError(t, err)
	assert.Equal(t, validPreparationInstrument, exampleValidPreparationInstrument)

	return created
}

func CreateValidMeasurementUnitForTest(t *testing.T, ctx context.Context, exampleValidMeasurementUnit *recipeenums.ValidMeasurementUnit, dbc recipeenums.Repository) *recipeenums.ValidMeasurementUnit {
	t.Helper()

	// create
	if exampleValidMeasurementUnit == nil {
		exampleValidMeasurementUnit = fakes.BuildFakeValidMeasurementUnit()
	}
	dbInput := converters.ConvertValidMeasurementUnitToValidMeasurementUnitDatabaseCreationInput(exampleValidMeasurementUnit)

	created, err := dbc.CreateValidMeasurementUnit(ctx, dbInput)
	exampleValidMeasurementUnit.CreatedAt = created.CreatedAt
	assert.NoError(t, err)
	assert.Equal(t, exampleValidMeasurementUnit, created)

	validMeasurementUnit, err := dbc.GetValidMeasurementUnit(ctx, created.ID)
	exampleValidMeasurementUnit.CreatedAt = validMeasurementUnit.CreatedAt

	assert.NoError(t, err)
	assert.Equal(t, validMeasurementUnit, exampleValidMeasurementUnit)

	return validMeasurementUnit
}

func CreateValidMeasurementUnitConversionForTest(t *testing.T, ctx context.Context, exampleValidMeasurementUnitConversion *recipeenums.ValidMeasurementUnitConversion, dbc recipeenums.Repository) *recipeenums.ValidMeasurementUnitConversion {
	t.Helper()

	// create
	if exampleValidMeasurementUnitConversion == nil {
		exampleValidMeasurementUnitConversion = fakes.BuildFakeValidMeasurementUnitConversion()
	}
	dbInput := converters.ConvertValidMeasurementUnitConversionToValidMeasurementUnitConversionDatabaseCreationInput(exampleValidMeasurementUnitConversion)

	created, err := dbc.CreateValidMeasurementUnitConversion(ctx, dbInput)
	require.NoError(t, err)
	require.NotNil(t, created)

	validMeasurementUnitConversion, err := dbc.GetValidMeasurementUnitConversion(ctx, created.ID)
	require.NoError(t, err)
	require.NotNil(t, validMeasurementUnitConversion)

	exampleValidMeasurementUnitConversion.CreatedAt = validMeasurementUnitConversion.CreatedAt
	assert.Equal(t, exampleValidMeasurementUnitConversion.From.ID, validMeasurementUnitConversion.From.ID)
	exampleValidMeasurementUnitConversion.From.CreatedAt = validMeasurementUnitConversion.From.CreatedAt
	assert.Equal(t, exampleValidMeasurementUnitConversion.To.ID, validMeasurementUnitConversion.To.ID)
	exampleValidMeasurementUnitConversion.To.CreatedAt = validMeasurementUnitConversion.To.CreatedAt
	exampleValidMeasurementUnitConversion.OnlyForIngredient = validMeasurementUnitConversion.OnlyForIngredient

	assert.Equal(t, validMeasurementUnitConversion, exampleValidMeasurementUnitConversion)

	return validMeasurementUnitConversion
}

func CreateValidIngredientStateForTest(t *testing.T, ctx context.Context, exampleValidIngredientState *recipeenums.ValidIngredientState, dbc recipeenums.Repository) *recipeenums.ValidIngredientState {
	t.Helper()

	// create
	if exampleValidIngredientState == nil {
		exampleValidIngredientState = fakes.BuildFakeValidIngredientState()
	}
	dbInput := converters.ConvertValidIngredientStateToValidIngredientStateDatabaseCreationInput(exampleValidIngredientState)

	created, err := dbc.CreateValidIngredientState(ctx, dbInput)
	exampleValidIngredientState.CreatedAt = created.CreatedAt
	assert.NoError(t, err)
	assert.Equal(t, exampleValidIngredientState, created)

	validIngredientState, err := dbc.GetValidIngredientState(ctx, created.ID)
	exampleValidIngredientState.CreatedAt = validIngredientState.CreatedAt

	assert.NoError(t, err)
	assert.Equal(t, validIngredientState, exampleValidIngredientState)

	return validIngredientState
}

func CreateValidIngredientStateIngredientForTest(t *testing.T, ctx context.Context, exampleValidIngredientStateIngredient *recipeenums.ValidIngredientStateIngredient, dbc recipeenums.Repository) *recipeenums.ValidIngredientStateIngredient {
	t.Helper()

	// create
	if exampleValidIngredientStateIngredient == nil {
		exampleValidIngredient := CreateValidIngredientForTest(t, ctx, nil, dbc)
		exampleValidIngredientState := CreateValidIngredientStateForTest(t, ctx, nil, dbc)
		exampleValidIngredientStateIngredient = fakes.BuildFakeValidIngredientStateIngredient()
		exampleValidIngredientStateIngredient.Ingredient = *exampleValidIngredient
		exampleValidIngredientStateIngredient.IngredientState = *exampleValidIngredientState
	}

	dbInput := converters.ConvertValidIngredientStateIngredientToValidIngredientStateIngredientDatabaseCreationInput(exampleValidIngredientStateIngredient)

	created, err := dbc.CreateValidIngredientStateIngredient(ctx, dbInput)
	assert.NoError(t, err)
	require.NotNil(t, created)
	exampleValidIngredientStateIngredient.CreatedAt = created.CreatedAt
	assert.Equal(t, exampleValidIngredientStateIngredient, created)

	validIngredientStateIngredient, err := dbc.GetValidIngredientStateIngredient(ctx, created.ID)
	exampleValidIngredientStateIngredient.CreatedAt = validIngredientStateIngredient.CreatedAt
	exampleValidIngredientStateIngredient.IngredientState = validIngredientStateIngredient.IngredientState
	exampleValidIngredientStateIngredient.Ingredient = validIngredientStateIngredient.Ingredient

	assert.NoError(t, err)
	assert.Equal(t, validIngredientStateIngredient, exampleValidIngredientStateIngredient)

	return validIngredientStateIngredient
}

func CreateValidIngredientPreparationForTest(t *testing.T, ctx context.Context, exampleValidIngredientPreparation *recipeenums.ValidIngredientPreparation, dbc recipeenums.Repository) *recipeenums.ValidIngredientPreparation {
	t.Helper()

	// create
	if exampleValidIngredientPreparation == nil {
		exampleValidIngredient := CreateValidIngredientForTest(t, ctx, nil, dbc)
		exampleValidPreparation := CreateValidPreparationForTest(t, ctx, nil, dbc)
		exampleValidIngredientPreparation = fakes.BuildFakeValidIngredientPreparation()
		exampleValidIngredientPreparation.Ingredient = *exampleValidIngredient
		exampleValidIngredientPreparation.Preparation = *exampleValidPreparation
	}

	dbInput := converters.ConvertValidIngredientPreparationToValidIngredientPreparationDatabaseCreationInput(exampleValidIngredientPreparation)

	created, err := dbc.CreateValidIngredientPreparation(ctx, dbInput)
	assert.NoError(t, err)
	require.NotNil(t, created)
	exampleValidIngredientPreparation.CreatedAt = created.CreatedAt
	assert.Equal(t, exampleValidIngredientPreparation, created)

	validIngredientPreparation, err := dbc.GetValidIngredientPreparation(ctx, created.ID)
	require.NotNil(t, validIngredientPreparation)
	exampleValidIngredientPreparation.CreatedAt = validIngredientPreparation.CreatedAt
	exampleValidIngredientPreparation.Preparation = validIngredientPreparation.Preparation
	exampleValidIngredientPreparation.Ingredient = validIngredientPreparation.Ingredient

	assert.NoError(t, err)
	assert.Equal(t, validIngredientPreparation, exampleValidIngredientPreparation)

	return created
}

func CreateValidIngredientMeasurementUnitForTest(t *testing.T, ctx context.Context, exampleValidIngredientMeasurementUnit *recipeenums.ValidIngredientMeasurementUnit, dbc recipeenums.Repository) *recipeenums.ValidIngredientMeasurementUnit {
	t.Helper()

	// create
	if exampleValidIngredientMeasurementUnit == nil {
		exampleValidIngredient := CreateValidIngredientForTest(t, ctx, nil, dbc)
		exampleValidMeasurementUnit := CreateValidMeasurementUnitForTest(t, ctx, nil, dbc)
		exampleValidIngredientMeasurementUnit = fakes.BuildFakeValidIngredientMeasurementUnit()
		exampleValidIngredientMeasurementUnit.Ingredient = *exampleValidIngredient
		exampleValidIngredientMeasurementUnit.MeasurementUnit = *exampleValidMeasurementUnit
	}

	dbInput := converters.ConvertValidIngredientMeasurementUnitToValidIngredientMeasurementUnitDatabaseCreationInput(exampleValidIngredientMeasurementUnit)

	created, err := dbc.CreateValidIngredientMeasurementUnit(ctx, dbInput)
	assert.NoError(t, err)
	require.NotNil(t, created)
	exampleValidIngredientMeasurementUnit.CreatedAt = created.CreatedAt
	assert.Equal(t, exampleValidIngredientMeasurementUnit, created)

	validIngredientMeasurementUnit, err := dbc.GetValidIngredientMeasurementUnit(ctx, created.ID)
	require.NoError(t, err)
	require.NotNil(t, validIngredientMeasurementUnit)
	exampleValidIngredientMeasurementUnit.CreatedAt = validIngredientMeasurementUnit.CreatedAt
	exampleValidIngredientMeasurementUnit.MeasurementUnit = validIngredientMeasurementUnit.MeasurementUnit
	exampleValidIngredientMeasurementUnit.Ingredient = validIngredientMeasurementUnit.Ingredient

	assert.NoError(t, err)
	assert.Equal(t, validIngredientMeasurementUnit, exampleValidIngredientMeasurementUnit)

	return created
}

func CreateValidIngredientGroupForTest(t *testing.T, ctx context.Context, exampleValidIngredientGroup *recipeenums.ValidIngredientGroup, dbc recipeenums.Repository) *recipeenums.ValidIngredientGroup {
	t.Helper()

	// create
	if exampleValidIngredientGroup == nil {
		exampleValidIngredientGroup = fakes.BuildFakeValidIngredientGroup()
	}
	dbInput := converters.ConvertValidIngredientGroupToValidIngredientGroupDatabaseCreationInput(exampleValidIngredientGroup)

	created, err := dbc.CreateValidIngredientGroup(ctx, dbInput)
	require.NoError(t, err)
	require.NotNil(t, created)
	exampleValidIngredientGroup.CreatedAt = created.CreatedAt
	for i := range exampleValidIngredientGroup.Members {
		exampleValidIngredientGroup.Members[i].CreatedAt = created.Members[i].CreatedAt
		exampleValidIngredientGroup.Members[i].ValidIngredient = created.Members[i].ValidIngredient
	}
	assert.Equal(t, exampleValidIngredientGroup, created)

	validIngredientGroup, err := dbc.GetValidIngredientGroup(ctx, created.ID)
	require.NoError(t, err)
	require.NotNil(t, validIngredientGroup)

	exampleValidIngredientGroup.CreatedAt = validIngredientGroup.CreatedAt
	for i := range exampleValidIngredientGroup.Members {
		exampleValidIngredientGroup.Members[i].CreatedAt = validIngredientGroup.Members[i].CreatedAt
		exampleValidIngredientGroup.Members[i].ValidIngredient = validIngredientGroup.Members[i].ValidIngredient
	}

	assert.NoError(t, err)
	assert.Equal(t, validIngredientGroup, exampleValidIngredientGroup)

	return created
}
