package main

import (
	"context"

	"github.com/prixfixeco/api_server/internal/database"
	"github.com/prixfixeco/api_server/pkg/types"
)

var validPreparationCollection = struct {
	Dice,
	Slice,
	Plate,
	Sautee,
	Marinate,
	Boil,
	Grill,
	Whisk,
	Mix,
	Mince,
	Knead,
	Divide,
	Flatten,
	Rest,
	Griddle,
	Grind,
	_ *types.ValidPreparationDatabaseCreationInput
}{
	Dice: &types.ValidPreparationDatabaseCreationInput{
		ID:   padID("vprep_dice"),
		Name: "Dice",
	},
	Slice: &types.ValidPreparationDatabaseCreationInput{
		ID:   padID("vprep_slice"),
		Name: "Slice",
	},
	Plate: &types.ValidPreparationDatabaseCreationInput{
		ID:   padID("vprep_plate"),
		Name: "Plate",
	},
	Sautee: &types.ValidPreparationDatabaseCreationInput{
		ID:   padID("vprep_sautee"),
		Name: "Sautee",
	},
	Marinate: &types.ValidPreparationDatabaseCreationInput{
		ID:   padID("vprep_marinate"),
		Name: "Marinate",
	},
	Boil: &types.ValidPreparationDatabaseCreationInput{
		ID:   padID("vprep_boil"),
		Name: "Boil",
	},
	Grill: &types.ValidPreparationDatabaseCreationInput{
		ID:   padID("vprep_grill"),
		Name: "Grill",
	},
	Whisk: &types.ValidPreparationDatabaseCreationInput{
		ID:   padID("vprep_whisk"),
		Name: "Whisk",
	},
	Mix: &types.ValidPreparationDatabaseCreationInput{
		ID:   padID("vprep_mix"),
		Name: "Mix",
	},
	Mince: &types.ValidPreparationDatabaseCreationInput{
		ID:   padID("vprep_mince"),
		Name: "Mince",
	},
	Knead: &types.ValidPreparationDatabaseCreationInput{
		ID:   padID("vprep_knead"),
		Name: "Knead",
	},
	Divide: &types.ValidPreparationDatabaseCreationInput{
		ID:   padID("vprep_divide"),
		Name: "Divide",
	},
	Flatten: &types.ValidPreparationDatabaseCreationInput{
		ID:   padID("vprep_flatten"),
		Name: "Flatten",
	},
	Rest: &types.ValidPreparationDatabaseCreationInput{
		ID:   padID("vprep_rest"),
		Name: "Rest",
	},
	Griddle: &types.ValidPreparationDatabaseCreationInput{
		ID:   padID("vprep_griddle"),
		Name: "Griddle",
	},
	Grind: &types.ValidPreparationDatabaseCreationInput{
		ID:   padID("vprep_grind"),
		Name: "Grind",
	},
}

func scaffoldValidPreparations(ctx context.Context, db database.DataManager) error {
	validPreparations := []*types.ValidPreparationDatabaseCreationInput{
		validPreparationCollection.Dice,
		validPreparationCollection.Slice,
		validPreparationCollection.Plate,
		validPreparationCollection.Sautee,
		validPreparationCollection.Marinate,
		validPreparationCollection.Boil,
		validPreparationCollection.Grill,
		validPreparationCollection.Whisk,
		validPreparationCollection.Mix,
		validPreparationCollection.Mince,
		validPreparationCollection.Grind,
		validPreparationCollection.Knead,
		validPreparationCollection.Divide,
		validPreparationCollection.Flatten,
		validPreparationCollection.Rest,
		validPreparationCollection.Griddle,
	}

	for _, input := range validPreparations {
		if _, err := db.CreateValidPreparation(ctx, input); err != nil {
			return err
		}
	}

	return nil
}
