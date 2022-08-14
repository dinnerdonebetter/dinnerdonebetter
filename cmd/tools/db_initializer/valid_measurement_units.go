package main

import (
	"context"

	"github.com/prixfixeco/api_server/internal/database"
	"github.com/prixfixeco/api_server/pkg/types"
)

var validMeasurementUnitCollection = struct {
	Gram,
	Milliliter,
	Unit,
	Clove,
	Teaspoon,
	Tablespoon,
	Can,
	Cup,
	Percent *types.ValidMeasurementUnitDatabaseCreationInput
}{
	Gram: &types.ValidMeasurementUnitDatabaseCreationInput{
		ID:         padID("vmu_gram"),
		Name:       "gram",
		Volumetric: false,
	},
	Milliliter: &types.ValidMeasurementUnitDatabaseCreationInput{
		ID:         padID("vmu_milliliter"),
		Name:       "milliliter",
		Volumetric: true,
	},
	Unit: &types.ValidMeasurementUnitDatabaseCreationInput{
		ID:         padID("vmu_unit"),
		Name:       "unit",
		Volumetric: false,
	},
	Clove: &types.ValidMeasurementUnitDatabaseCreationInput{
		ID:         padID("vmu_clove"),
		Name:       "clove",
		Volumetric: false,
	},
	Teaspoon: &types.ValidMeasurementUnitDatabaseCreationInput{
		ID:         padID("vmu_teaspoon"),
		Name:       "teaspoon",
		Volumetric: true,
	},
	Tablespoon: &types.ValidMeasurementUnitDatabaseCreationInput{
		ID:         padID("vmu_tablespoon"),
		Name:       "tablespoon",
		Volumetric: true,
	},
	Can: &types.ValidMeasurementUnitDatabaseCreationInput{
		ID:         padID("vmu_can"),
		Name:       "can",
		Volumetric: false,
	},
	Cup: &types.ValidMeasurementUnitDatabaseCreationInput{
		ID:         padID("vmu_cup"),
		Name:       "cup",
		Volumetric: true,
	},
	Percent: &types.ValidMeasurementUnitDatabaseCreationInput{
		ID:         padID("vmu_percent"),
		Name:       "percent",
		Volumetric: false,
	},
}

func scaffoldValidMeasurementUnits(ctx context.Context, db database.DataManager) error {
	validMeasurementUnits := []*types.ValidMeasurementUnitDatabaseCreationInput{
		validMeasurementUnitCollection.Gram,
		validMeasurementUnitCollection.Milliliter,
		validMeasurementUnitCollection.Unit,
		validMeasurementUnitCollection.Clove,
		validMeasurementUnitCollection.Teaspoon,
		validMeasurementUnitCollection.Tablespoon,
		validMeasurementUnitCollection.Can,
		validMeasurementUnitCollection.Cup,
		validMeasurementUnitCollection.Percent,
	}

	for _, input := range validMeasurementUnits {
		if _, err := db.CreateValidMeasurementUnit(ctx, input); err != nil {
			return err
		}
	}

	return nil
}
