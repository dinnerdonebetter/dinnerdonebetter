package main

import (
	"context"
	"fmt"

	"github.com/prixfixeco/api_server/internal/database"
	"github.com/prixfixeco/api_server/pkg/types"
)

var validPreparationInstrumentCollection = struct {
	Dice,
	Slice,
	Plate,
	Sautee10,
	Sautee12,
	MarinateSmBowl,
	MarinateMedBowl,
	MarinateLgBowl,
	Boil,
	Grill,
	Whisk,
	MixWhisk,
	MixBareHands,
	Knead,
	Divide,
	Flatten,
	RestSmBowl,
	RestMedBowl,
	RestLgBowl,
	Griddle,
	_ *types.ValidPreparationInstrumentDatabaseCreationInput
}{
	Dice: &types.ValidPreparationInstrumentDatabaseCreationInput{
		ID:                 "vpi_dice_knife",
		ValidPreparationID: "vprep_dice",
		ValidInstrumentID:  "vinst_chefsknife",
	},
	Slice: &types.ValidPreparationInstrumentDatabaseCreationInput{
		ID:                 "vpi_slice_knife",
		ValidPreparationID: "vprep_slice",
		ValidInstrumentID:  "vinst_chefsknife",
	},
	Plate: &types.ValidPreparationInstrumentDatabaseCreationInput{
		ID:                 "vpi_plate_bh",
		ValidPreparationID: "vprep_plate",
		ValidInstrumentID:  "vinst_barehands",
	},
	Sautee10: &types.ValidPreparationInstrumentDatabaseCreationInput{
		ID:                 "vpi_sautee10",
		ValidPreparationID: "vprep_sautee",
		ValidInstrumentID:  "vinst_10inchfp",
	},
	Sautee12: &types.ValidPreparationInstrumentDatabaseCreationInput{
		ID:                 "vpi_sautee12",
		ValidPreparationID: "vprep_sautee",
		ValidInstrumentID:  "vinst_12inchfp",
	},
	MarinateSmBowl: &types.ValidPreparationInstrumentDatabaseCreationInput{
		ID:                 "vpi_mary_sm",
		ValidPreparationID: "vprep_marinate",
		ValidInstrumentID:  "vinst_smmixbowl",
	},
	MarinateMedBowl: &types.ValidPreparationInstrumentDatabaseCreationInput{
		ID:                 "vpi_mary_md",
		ValidPreparationID: "vprep_marinate",
		ValidInstrumentID:  "vinst_medmixbowl",
	},
	MarinateLgBowl: &types.ValidPreparationInstrumentDatabaseCreationInput{
		ID:                 "vpi_mary_lg",
		ValidPreparationID: "vprep_marinate",
		ValidInstrumentID:  "vinst_lgmixbowl",
	},
	Boil: &types.ValidPreparationInstrumentDatabaseCreationInput{
		ID:                 "vpi_boil_sp",
		ValidPreparationID: "vprep_boil",
		ValidInstrumentID:  "vinst_4qtsaucepan",
	},
	Grill: &types.ValidPreparationInstrumentDatabaseCreationInput{
		ID:                 "vpi_grill2",
		ValidPreparationID: "vprep_grill",
		ValidInstrumentID:  "vinst_grill",
	},
	Whisk: &types.ValidPreparationInstrumentDatabaseCreationInput{
		ID:                 "vpi_whisk2",
		ValidPreparationID: "vprep_whisk",
		ValidInstrumentID:  "vinst_whisk",
	},
	MixWhisk: &types.ValidPreparationInstrumentDatabaseCreationInput{
		ID:                 "vpi_mix_whisk",
		ValidPreparationID: "vprep_mix",
		ValidInstrumentID:  "vinst_whisk",
	},
	MixBareHands: &types.ValidPreparationInstrumentDatabaseCreationInput{
		ID:                 "vpi_mix_bh",
		ValidPreparationID: "vprep_mix",
		ValidInstrumentID:  "vinst_barehands",
	},
	Knead: &types.ValidPreparationInstrumentDatabaseCreationInput{
		ID:                 "vpi_knead_bh",
		ValidPreparationID: "vprep_knead",
		ValidInstrumentID:  "vinst_barehands",
	},
	Divide: &types.ValidPreparationInstrumentDatabaseCreationInput{
		ID:                 "vpi_divide_ps",
		ValidPreparationID: "vprep_divide",
		ValidInstrumentID:  "vinst_pastryscraper",
	},
	Flatten: &types.ValidPreparationInstrumentDatabaseCreationInput{
		ID:                 "vpi_flatten_bh",
		ValidPreparationID: "vprep_flatten",
		ValidInstrumentID:  "vinst_barehands",
	},
	RestSmBowl: &types.ValidPreparationInstrumentDatabaseCreationInput{
		ID:                 "vpi_rest_smb",
		ValidPreparationID: "vprep_rest",
		ValidInstrumentID:  "vinst_smmixbowl",
	},
	RestMedBowl: &types.ValidPreparationInstrumentDatabaseCreationInput{
		ID:                 "vpi_rest_mb",
		ValidPreparationID: "vprep_rest",
		ValidInstrumentID:  "vinst_medmixbowl",
	},
	RestLgBowl: &types.ValidPreparationInstrumentDatabaseCreationInput{
		ID:                 "vpi_rest_lb",
		ValidPreparationID: "vprep_rest",
		ValidInstrumentID:  "vinst_lgmixbowl",
	},
	Griddle: &types.ValidPreparationInstrumentDatabaseCreationInput{
		ID:                 "vpi_griddle_comal",
		ValidPreparationID: "vprep_griddle",
		ValidInstrumentID:  "vinst_comal",
	},
}

func scaffoldValidPreparationInstruments(ctx context.Context, db database.DataManager) error {
	validPreparationInstruments := []*types.ValidPreparationInstrumentDatabaseCreationInput{
		validPreparationInstrumentCollection.Dice,
		validPreparationInstrumentCollection.Slice,
		validPreparationInstrumentCollection.Plate,
		validPreparationInstrumentCollection.Sautee10,
		validPreparationInstrumentCollection.Sautee12,
		validPreparationInstrumentCollection.MarinateSmBowl,
		validPreparationInstrumentCollection.MarinateMedBowl,
		validPreparationInstrumentCollection.MarinateLgBowl,
		validPreparationInstrumentCollection.Boil,
		validPreparationInstrumentCollection.Grill,
		validPreparationInstrumentCollection.Whisk,
		validPreparationInstrumentCollection.MixWhisk,
		validPreparationInstrumentCollection.MixBareHands,
		validPreparationInstrumentCollection.Knead,
		validPreparationInstrumentCollection.Divide,
		validPreparationInstrumentCollection.Flatten,
		validPreparationInstrumentCollection.RestSmBowl,
		validPreparationInstrumentCollection.RestMedBowl,
		validPreparationInstrumentCollection.RestLgBowl,
		validPreparationInstrumentCollection.Griddle,
	}

	for i, input := range validPreparationInstruments {
		if _, err := db.CreateValidPreparationInstrument(ctx, input); err != nil {
			return fmt.Errorf("creating preparation instrument #%d: %w", i, err)
		}
	}

	return nil
}
