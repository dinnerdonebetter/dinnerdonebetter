package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/dinnerdonebetter/backend/internal/database"
	dbconfig "github.com/dinnerdonebetter/backend/internal/database/config"
	"github.com/dinnerdonebetter/backend/internal/database/postgres"
	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/logging/zap"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/pkg/pointers"
	"github.com/dinnerdonebetter/backend/internal/search"
	"github.com/dinnerdonebetter/backend/internal/search/algolia"
	searchcfg "github.com/dinnerdonebetter/backend/internal/search/config"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

func main() {
	indexPtr := flag.String("index", "", "index to initialize")
	wipePtr := flag.Bool("wipe", false, "whether to wipe the index or not")

	flag.Parse()

	index := *indexPtr
	if index == "" {
		log.Fatal("index is required")
	}

	ctx := context.Background()
	logger := zap.NewZapLogger(logging.DebugLevel)
	tracerProvider := tracing.NewNoopTracerProvider()

	cfg := &searchcfg.Config{
		Provider: searchcfg.AlgoliaProvider,
		Algolia: &algolia.Config{
			AppID:  os.Getenv("ALGOLIA_APP_ID"),
			APIKey: os.Getenv("ALGOLIA_API_KEY"),
		},
	}

	im, err := searchcfg.ProvideIndexManager[search.RecipeSearchSubset](ctx, logger, tracerProvider, cfg, index)
	if err != nil {
		log.Fatal(fmt.Errorf("initializing index manager: %w", err))
	}

	if *wipePtr {
		log.Println("wiping index")
		if err = im.Wipe(ctx); err != nil {
			log.Fatal(fmt.Errorf("wiping index: %w", err))
		}
		log.Println("wiped index")
	}

	dbConfig := &dbconfig.Config{
		ConnectionDetails: database.ConnectionDetails(os.Getenv("DATABASE_URL")),
	}

	dataManager, err := postgres.ProvideDatabaseClient(ctx, logger, dbConfig, tracing.NewNoopTracerProvider())
	if err != nil {
		log.Fatal(fmt.Errorf("initializing database client: %w", err))
	}

	switch index {
	case "recipes":
		filter := types.DefaultQueryFilter()
		filter.Limit = pointers.Pointer(uint8(50))

		var thresholdMet bool
		for !thresholdMet {
			var recipes *types.QueryFilteredResult[types.Recipe]
			recipes, err = dataManager.GetRecipes(ctx, filter)
			if err != nil {
				log.Fatal(fmt.Errorf("getting recipes: %w", err))
			}

			for _, x := range recipes.Data {
				var recipe *types.Recipe
				recipe, err = dataManager.GetRecipe(ctx, x.ID)
				if err != nil {
					log.Fatal(fmt.Errorf("getting recipe: %w", err))
				}

				toBeIndexed := search.SubsetFromRecipe(recipe)
				if err = im.Index(ctx, recipe.ID, toBeIndexed); err != nil {
					log.Fatal(fmt.Errorf("indexing recipe: %w", err))
				}
			}

			thresholdMet = len(recipes.Data) == 0
			*filter.Page++
		}
	}

	results, err := im.Search(ctx, "boil water")
	if err != nil {
		log.Fatal(fmt.Errorf("searching index: %w", err))
	}

	log.Println(results)
}
