package main

import (
	"context"

	"github.com/prixfixeco/api_server/internal/observability/logging"
	"github.com/prixfixeco/api_server/internal/search"
	"github.com/prixfixeco/api_server/internal/search/realasticsearch"
)

func main() {
	logger := logging.NewZerologLogger()
	//httpClient := observability.HTTPClient()

	cfg := &search.Config{
		Provider: search.ElasticsearchProvider,
		Address:  "https://e43f4c67a36b44db8b37a111cb28ae69.us-east-1.aws.found.io:9243",
		Username: "elastic",
		Password: "lRLCRwZ600vntA6i7zeS6tRA",
	}

	imp, err := elasticsearch.NewIndexManagerProvider(logger, cfg)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	im, err := imp.ProvideIndexManager(ctx, logger, "things", "name", "description")
	if err != nil {
		panic(err)
	}

	im.Search(ctx, "", "", "")
}
