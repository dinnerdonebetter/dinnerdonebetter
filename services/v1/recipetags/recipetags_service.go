package recipetags

import (
	"fmt"
	"net/http"

	"gitlab.com/prixfixe/prixfixe/internal/v1/encoding"
	"gitlab.com/prixfixe/prixfixe/internal/v1/metrics"
	models "gitlab.com/prixfixe/prixfixe/models/v1"

	"gitlab.com/verygoodsoftwarenotvirus/logging/v1"
	"gitlab.com/verygoodsoftwarenotvirus/newsman"
)

const (
	// CreateMiddlewareCtxKey is a string alias we can use for referring to recipe tag input data in contexts.
	CreateMiddlewareCtxKey models.ContextKey = "recipe_tag_create_input"
	// UpdateMiddlewareCtxKey is a string alias we can use for referring to recipe tag update data in contexts.
	UpdateMiddlewareCtxKey models.ContextKey = "recipe_tag_update_input"

	counterName        metrics.CounterName = "recipeTags"
	counterDescription string              = "the number of recipeTags managed by the recipeTags service"
	topicName          string              = "recipe_tags"
	serviceName        string              = "recipe_tags_service"
)

var (
	_ models.RecipeTagDataServer = (*Service)(nil)
)

type (
	// Service handles to-do list recipe tags
	Service struct {
		logger               logging.Logger
		recipeDataManager    models.RecipeDataManager
		recipeTagDataManager models.RecipeTagDataManager
		recipeIDFetcher      RecipeIDFetcher
		recipeTagIDFetcher   RecipeTagIDFetcher
		userIDFetcher        UserIDFetcher
		recipeTagCounter     metrics.UnitCounter
		encoderDecoder       encoding.EncoderDecoder
		reporter             newsman.Reporter
	}

	// UserIDFetcher is a function that fetches user IDs.
	UserIDFetcher func(*http.Request) uint64

	// RecipeIDFetcher is a function that fetches recipe IDs.
	RecipeIDFetcher func(*http.Request) uint64

	// RecipeTagIDFetcher is a function that fetches recipe tag IDs.
	RecipeTagIDFetcher func(*http.Request) uint64
)

// ProvideRecipeTagsService builds a new RecipeTagsService.
func ProvideRecipeTagsService(
	logger logging.Logger,
	recipeDataManager models.RecipeDataManager,
	recipeTagDataManager models.RecipeTagDataManager,
	recipeIDFetcher RecipeIDFetcher,
	recipeTagIDFetcher RecipeTagIDFetcher,
	userIDFetcher UserIDFetcher,
	encoder encoding.EncoderDecoder,
	recipeTagCounterProvider metrics.UnitCounterProvider,
	reporter newsman.Reporter,
) (*Service, error) {
	recipeTagCounter, err := recipeTagCounterProvider(counterName, counterDescription)
	if err != nil {
		return nil, fmt.Errorf("error initializing counter: %w", err)
	}

	svc := &Service{
		logger:               logger.WithName(serviceName),
		recipeIDFetcher:      recipeIDFetcher,
		recipeTagIDFetcher:   recipeTagIDFetcher,
		userIDFetcher:        userIDFetcher,
		recipeDataManager:    recipeDataManager,
		recipeTagDataManager: recipeTagDataManager,
		encoderDecoder:       encoder,
		recipeTagCounter:     recipeTagCounter,
		reporter:             reporter,
	}

	return svc, nil
}
