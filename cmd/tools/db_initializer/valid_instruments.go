package main

import (
	"context"

	"github.com/prixfixeco/api_server/internal/database"
	"github.com/prixfixeco/api_server/pkg/types"

	"github.com/segmentio/ksuid"
)

var validInstrumentCollection = struct {
	Spoon *types.ValidInstrumentDatabaseCreationInput
}{
	Spoon: &types.ValidInstrumentDatabaseCreationInput{
		ID:          ksuid.New().String(),
		Name:        "spoon",
		Variant:     "",
		Description: "",
		IconPath:    "",
	},
}

func scaffoldValidInstruments(ctx context.Context, db database.DataManager) error {
	validInstruments := []*types.ValidInstrumentDatabaseCreationInput{
		validInstrumentCollection.Spoon,
	}

	for _, input := range validInstruments {
		if _, err := db.CreateValidInstrument(ctx, input); err != nil {
			return err
		}
	}

	return nil
}
