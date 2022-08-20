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
		ID:                 padID("vpi_dice_knife"),
		ValidPreparationID: padID("vprep_dice"),
		ValidInstrumentID:  padID("vinst_chefsknife"),
	},
	Slice: &types.ValidPreparationInstrumentDatabaseCreationInput{
		ID:                 padID("vpi_slice_knife"),
		ValidPreparationID: padID("vprep_slice"),
		ValidInstrumentID:  padID("vinst_chefsknife"),
	},
	Plate: &types.ValidPreparationInstrumentDatabaseCreationInput{
		ID:                 padID("vpi_plate_bh"),
		ValidPreparationID: padID("vprep_plate"),
		ValidInstrumentID:  padID("vinst_barehands"),
	},
	Sautee10: &types.ValidPreparationInstrumentDatabaseCreationInput{
		ID:                 padID("vpi_sautee10"),
		ValidPreparationID: padID("vprep_sautee"),
		ValidInstrumentID:  padID("vinst_10inchfp"),
	},
	Sautee12: &types.ValidPreparationInstrumentDatabaseCreationInput{
		ID:                 padID("vpi_sautee12"),
		ValidPreparationID: padID("vprep_sautee"),
		ValidInstrumentID:  padID("vinst_12inchfp"),
	},
	MarinateSmBowl: &types.ValidPreparationInstrumentDatabaseCreationInput{
		ID:                 padID("vpi_mary_sm"),
		ValidPreparationID: padID("vprep_marinate"),
		ValidInstrumentID:  padID("vinst_smmixbowl"),
	},
	MarinateMedBowl: &types.ValidPreparationInstrumentDatabaseCreationInput{
		ID:                 padID("vpi_mary_md"),
		ValidPreparationID: padID("vprep_marinate"),
		ValidInstrumentID:  padID("vinst_medmixbowl"),
	},
	MarinateLgBowl: &types.ValidPreparationInstrumentDatabaseCreationInput{
		ID:                 padID("vpi_mary_lg"),
		ValidPreparationID: padID("vprep_marinate"),
		ValidInstrumentID:  padID("vinst_lgmixbowl"),
	},
	Boil: &types.ValidPreparationInstrumentDatabaseCreationInput{
		ID:                 padID("vpi_boil_sp"),
		ValidPreparationID: padID("vprep_boil"),
		ValidInstrumentID:  padID("vinst_4qtsaucepan"),
	},
	Grill: &types.ValidPreparationInstrumentDatabaseCreationInput{
		ID:                 padID("vpi_grill2"),
		ValidPreparationID: padID("vprep_grill"),
		ValidInstrumentID:  padID("vinst_grill"),
	},
	Whisk: &types.ValidPreparationInstrumentDatabaseCreationInput{
		ID:                 padID("vpi_whisk2"),
		ValidPreparationID: padID("vprep_whisk"),
		ValidInstrumentID:  padID("vinst_whisk"),
	},
	MixWhisk: &types.ValidPreparationInstrumentDatabaseCreationInput{
		ID:                 padID("vpi_mix_whisk"),
		ValidPreparationID: padID("vprep_mix"),
		ValidInstrumentID:  padID("vinst_whisk"),
	},
	MixBareHands: &types.ValidPreparationInstrumentDatabaseCreationInput{
		ID:                 padID("vpi_mix_bh"),
		ValidPreparationID: padID("vprep_mix"),
		ValidInstrumentID:  padID("vinst_barehands"),
	},
	Knead: &types.ValidPreparationInstrumentDatabaseCreationInput{
		ID:                 padID("vpi_knead_bh"),
		ValidPreparationID: padID("vprep_knead"),
		ValidInstrumentID:  padID("vinst_barehands"),
	},
	Divide: &types.ValidPreparationInstrumentDatabaseCreationInput{
		ID:                 padID("vpi_divide_ps"),
		ValidPreparationID: padID("vprep_divide"),
		ValidInstrumentID:  padID("vinst_pastryscraper"),
	},
	Flatten: &types.ValidPreparationInstrumentDatabaseCreationInput{
		ID:                 padID("vpi_flatten_bh"),
		ValidPreparationID: padID("vprep_flatten"),
		ValidInstrumentID:  padID("vinst_barehands"),
	},
	RestSmBowl: &types.ValidPreparationInstrumentDatabaseCreationInput{
		ID:                 padID("vpi_rest_smb"),
		ValidPreparationID: padID("vprep_rest"),
		ValidInstrumentID:  padID("vinst_smmixbowl"),
	},
	RestMedBowl: &types.ValidPreparationInstrumentDatabaseCreationInput{
		ID:                 padID("vpi_rest_mb"),
		ValidPreparationID: padID("vprep_rest"),
		ValidInstrumentID:  padID("vinst_medmixbowl"),
	},
	RestLgBowl: &types.ValidPreparationInstrumentDatabaseCreationInput{
		ID:                 padID("vpi_rest_lb"),
		ValidPreparationID: padID("vprep_rest"),
		ValidInstrumentID:  padID("vinst_lgmixbowl"),
	},
	Griddle: &types.ValidPreparationInstrumentDatabaseCreationInput{
		ID:                 padID("vpi_griddle_comal"),
		ValidPreparationID: padID("vprep_griddle"),
		ValidInstrumentID:  padID("vinst_comal"),
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
