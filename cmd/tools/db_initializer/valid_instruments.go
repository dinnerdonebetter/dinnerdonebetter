package main

import (
	"context"
	"github.com/prixfixeco/api_server/internal/database"
	"github.com/prixfixeco/api_server/pkg/types"
)

var validInstrumentCollection = struct {
	ChefsKnife,
	Grill,
	BareHands,
	TenInchFryingPan,
	TwelveInchFryingPan,
	FourQuartSaucepan,
	Whisk,
	Comal,
	SmallMixingBowl,
	MediumMixingBowl,
	LargeMixingBowl,
	PastryScraper,
	RollingPin,

	_ *types.ValidInstrumentDatabaseCreationInput
}{
	ChefsKnife: &types.ValidInstrumentDatabaseCreationInput{
		ID:          padID("vinst_chefsknife"),
		Name:        "chef's knife",
		Description: "",
		IconPath:    "",
	},
	Grill: &types.ValidInstrumentDatabaseCreationInput{
		ID:          padID("vinst_grill"),
		Name:        "grill",
		Description: "",
		IconPath:    "",
	},
	BareHands: &types.ValidInstrumentDatabaseCreationInput{
		ID:          padID("vinst_barehands"),
		Name:        "bare hands",
		Description: "",
		IconPath:    "",
	},
	TenInchFryingPan: &types.ValidInstrumentDatabaseCreationInput{
		ID:          padID("vinst_10inchfp"),
		Name:        `10" frying pan`,
		Description: "",
		IconPath:    "",
	},
	TwelveInchFryingPan: &types.ValidInstrumentDatabaseCreationInput{
		ID:          padID("vinst_12inchfp"),
		Name:        `12" frying pan`,
		Description: "",
		IconPath:    "",
	},
	FourQuartSaucepan: &types.ValidInstrumentDatabaseCreationInput{
		ID:          padID("vinst_4qtsaucepan"),
		Name:        "4 quart saucepan",
		Description: "",
		IconPath:    "",
	},
	Whisk: &types.ValidInstrumentDatabaseCreationInput{
		ID:          padID("vinst_whisk"),
		Name:        "whisk",
		Description: "",
		IconPath:    "",
	},
	Comal: &types.ValidInstrumentDatabaseCreationInput{
		ID:          padID("vinst_comal"),
		Name:        "comal",
		Description: "",
		IconPath:    "",
	},
	SmallMixingBowl: &types.ValidInstrumentDatabaseCreationInput{
		ID:          padID("vinst_smmixbowl"),
		Name:        "small mixing bowl",
		Description: "",
		IconPath:    "",
	},
	MediumMixingBowl: &types.ValidInstrumentDatabaseCreationInput{
		ID:          padID("vinst_medmixbowl"),
		Name:        "medium mixing bowl",
		Description: "",
		IconPath:    "",
	},
	LargeMixingBowl: &types.ValidInstrumentDatabaseCreationInput{
		ID:          padID("vinst_lgmixbowl"),
		Name:        "large mixing bowl",
		Description: "",
		IconPath:    "",
	},
	PastryScraper: &types.ValidInstrumentDatabaseCreationInput{
		ID:          padID("vinst_pastryscraper"),
		Name:        "pastry scraper",
		Description: "",
		IconPath:    "",
	},
	RollingPin: &types.ValidInstrumentDatabaseCreationInput{
		ID:          padID("vinst_rollingpin"),
		Name:        "rolling pin",
		Description: "",
		IconPath:    "",
	},
}

func scaffoldValidInstruments(ctx context.Context, db database.DataManager) error {
	validInstruments := []*types.ValidInstrumentDatabaseCreationInput{
		validInstrumentCollection.ChefsKnife,
		validInstrumentCollection.Grill,
		validInstrumentCollection.BareHands,
		validInstrumentCollection.TenInchFryingPan,
		validInstrumentCollection.TwelveInchFryingPan,
		validInstrumentCollection.FourQuartSaucepan,
		validInstrumentCollection.Whisk,
		validInstrumentCollection.Comal,
		validInstrumentCollection.SmallMixingBowl,
		validInstrumentCollection.MediumMixingBowl,
		validInstrumentCollection.LargeMixingBowl,
		validInstrumentCollection.PastryScraper,
		validInstrumentCollection.RollingPin,
	}

	for _, input := range validInstruments {
		if _, err := db.CreateValidInstrument(ctx, input); err != nil {
			return err
		}
	}

	return nil
}
