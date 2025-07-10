package datachangemessagehandler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/metrics"
	textsearch "github.com/dinnerdonebetter/backend/internal/platform/search/text"
	textsearchcfg "github.com/dinnerdonebetter/backend/internal/platform/search/text/config"
	coreindexing "github.com/dinnerdonebetter/backend/internal/services/core/indexing"
	eatingindexing "github.com/dinnerdonebetter/backend/internal/services/eating/indexing"
	"time"
)

func (a *AsyncDataChangeMessageHandler) buildSearchIndexRequestsEventHandler(
	metricsProvider metrics.Provider,
	searchCfg *textsearchcfg.Config,
) func(context.Context, []byte) error {
	return func(ctx context.Context, rawMsg []byte) error {
		ctx, span := a.tracer.StartSpan(ctx)
		defer span.End()

		start := time.Now()

		var searchIndexRequest textsearch.IndexRequest
		if err := json.NewDecoder(bytes.NewReader(rawMsg)).Decode(&searchIndexRequest); err != nil {
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
			if err := eatingindexing.HandleIndexRequest(ctx, a.logger, tracerProvider, metricsProvider, searchCfg, a.mealPlanningRepo, &searchIndexRequest); err != nil {
				return fmt.Errorf("handling search indexing request: %w", err)
			}

		case coreindexing.IndexTypeUsers:
			// we don't want to retry indexing perpetually in the event of a fundamental error, so we just log it and move on
			if err := coreindexing.HandleIndexRequest(ctx, a.logger, tracerProvider, metricsProvider, searchCfg, a.identityRepo, &searchIndexRequest); err != nil {
				return fmt.Errorf("handling search indexing request: %w", err)
			}
		}

		a.searchIndexRequestsExecutionTimeHistogram.Record(ctx, float64(time.Since(start).Milliseconds()))

		return nil
	}
}
