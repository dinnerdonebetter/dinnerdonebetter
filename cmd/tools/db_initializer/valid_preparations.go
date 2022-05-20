package main

import (
	"context"

	"github.com/prixfixeco/api_server/internal/database"
	"github.com/prixfixeco/api_server/pkg/types"

	"github.com/segmentio/ksuid"
)

var validPreparationCollection = struct {
	Dice,
	Slice,
	Plate,
	Sautee,
	Marinate,
	Boil,
	Grill,
	Drain *types.ValidPreparationDatabaseCreationInput
}{
	Dice: &types.ValidPreparationDatabaseCreationInput{
		ID:          ksuid.New().String(),
		Name:        "Dice",
		Description: "",
		IconPath:    "",
	},
	Slice: &types.ValidPreparationDatabaseCreationInput{
		ID:          ksuid.New().String(),
		Name:        "Slice",
		Description: "",
		IconPath:    "",
	},
	Plate: &types.ValidPreparationDatabaseCreationInput{
		ID:          ksuid.New().String(),
		Name:        "Plate",
		Description: "",
		IconPath:    "",
	},
	Sautee: &types.ValidPreparationDatabaseCreationInput{
		ID:          ksuid.New().String(),
		Name:        "Sautee",
		Description: "",
		IconPath:    "",
	},
	Marinate: &types.ValidPreparationDatabaseCreationInput{
		ID:          ksuid.New().String(),
		Name:        "Marinate",
		Description: "",
		IconPath:    "",
	},
	Boil: &types.ValidPreparationDatabaseCreationInput{
		ID:          ksuid.New().String(),
		Name:        "Boil",
		Description: "",
		IconPath:    "",
	},
	Grill: &types.ValidPreparationDatabaseCreationInput{
		ID:          ksuid.New().String(),
		Name:        "Grill",
		Description: "",
		IconPath:    "",
	},
	Drain: &types.ValidPreparationDatabaseCreationInput{
		ID:          ksuid.New().String(),
		Name:        "Drain",
		Description: "",
		IconPath:    "",
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
		validPreparationCollection.Drain,
	}

	for _, input := range validPreparations {
		if _, err := db.CreateValidPreparation(ctx, input); err != nil {
			return err
		}
	}

	return nil
}
