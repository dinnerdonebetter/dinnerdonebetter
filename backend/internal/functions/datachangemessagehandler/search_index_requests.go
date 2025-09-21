package datachangemessagehandler

import (
	"context"
	"fmt"
	"time"

	textsearch "github.com/dinnerdonebetter/backend/internal/platform/search/text"
	coreindexing "github.com/dinnerdonebetter/backend/internal/services/identity/indexing"
	eatingindexing "github.com/dinnerdonebetter/backend/internal/services/mealplanning/indexing"
)

func (a *AsyncDataChangeMessageHandler) SearchIndexRequestsEventHandler(ctx context.Context, rawMsg []byte) error {
	ctx, span := a.tracer.StartSpan(ctx)
	defer span.End()

	start := time.Now()

	var searchIndexRequest textsearch.IndexRequest
	if err := a.decoder.DecodeBytes(ctx, rawMsg, &searchIndexRequest); err != nil {
		return fmt.Errorf("decoding JSON body: %w", err)
	}

	switch searchIndexRequest.IndexType {
	case eatingindexing.IndexTypeRecipes,
		eatingindexing.IndexTypeMeals,
		eatingindexing.IndexTypeValidIngredients,
		eatingindexing.IndexTypeValidInstruments,
		eatingindexing.IndexTypeValidMeasurementUnits,
		eatingindexing.IndexTypeValidPreparations,
		eatingindexing.IndexTypeValidIngredientStates,
		eatingindexing.IndexTypeValidVessels:
		// we don't want to retry indexing perpetually in the event of a fundamental error, so we just log it and move on
		if err := a.mealPlanningDataIndexer.HandleIndexRequest(ctx, &searchIndexRequest); err != nil {
			return fmt.Errorf("handling search indexing request: %w", err)
		}

	case coreindexing.IndexTypeUsers:
		// we don't want to retry indexing perpetually in the event of a fundamental error, so we just log it and move on
		if err := a.userDataIndexer.HandleIndexRequest(ctx, &searchIndexRequest); err != nil {
			return fmt.Errorf("handling search indexing request: %w", err)
		}
	}

	a.searchIndexRequestsExecutionTimeHistogram.Record(ctx, float64(time.Since(start).Milliseconds()))

	return nil
}
