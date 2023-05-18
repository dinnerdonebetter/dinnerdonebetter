package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/dinnerdonebetter/backend/internal/database"
	dbconfig "github.com/dinnerdonebetter/backend/internal/database/config"
	"github.com/dinnerdonebetter/backend/internal/database/postgres"
	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/logging/zerolog"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
)

func main() {
	ctx := context.Background()
	logger := zerolog.NewZerologLogger(logging.DebugLevel)

	dbConfig := &dbconfig.Config{
		ConnectionDetails: database.ConnectionDetails(os.Getenv("DB_CONNECTION_STRING")),
	}

	tracerProvider := tracing.NewNoopTracerProvider()

	dataManager, err := postgres.ProvideDatabaseClient(ctx, logger, dbConfig, tracerProvider)
	if err != nil {
		panic(err)
	}

	defer dataManager.Close()

	var prompt strings.Builder
	prompt.WriteString(`You are a recipe state machine. Your job is to produce JSON representations of recipes. You are to adhere to a strict schema:

1. Every possible preparation you could describe in a recipe has a unique identifier.
2. Every possible ingredient you could describe in a recipe has a unique identifier.
3. Every possible instrument you could describe in a recipe has a unique identifier.
4. Every possible cooking or preparation vessel you could describe in a recipe has a unique identifier.
5. Every possible measurement unit you could describe in a recipe has a unique identifier.
6. Every step is comprised of exactly one preparation, any number of ingredients, any number of instruments, any number of vessels.


`)

	fmt.Println(prompt.String())
}
