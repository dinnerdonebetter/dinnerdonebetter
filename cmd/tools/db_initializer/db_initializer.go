package main

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"sync"

	flag "github.com/spf13/pflag"
	"go.opentelemetry.io/otel/trace"

	logcfg "github.com/prixfixeco/api_server/internal/observability/logging/config"
	"github.com/prixfixeco/api_server/pkg/client/httpclient"
	"github.com/prixfixeco/api_server/pkg/types"
	testutils "github.com/prixfixeco/api_server/tests/utils"
)

var (
	address    string
	username   string
	password   string
	totpSecret string
	debug      bool
)

func init() {
	flag.StringVarP(&address, "address", "a", "", "where the target instance is hosted")
	flag.StringVarP(&username, "username", "u", "", "admin username")
	flag.StringVarP(&password, "password", "p", "", "admin password")
	flag.StringVarP(&totpSecret, "two-factor-secret", "t", "", "admin 2FA secret")
	flag.BoolVarP(&debug, "debug", "d", false, "whether debug mode is enabled")
}

func main() {
	flag.Parse()

	ctx := context.Background()
	logger := (&logcfg.Config{Provider: logcfg.ProviderZerolog}).ProvideLogger()

	if address == "" {
		logger.Fatal(errors.New("uri must be valid"))
	}

	if username == "" || password == "" || address == "" {
		logger.Fatal(errors.New("all credentials must be provided"))
	}

	parsedURI, uriParseErr := url.Parse(address)
	if uriParseErr != nil {
		logger.Fatal(fmt.Errorf("parsing provided url: %w", uriParseErr))
	}
	if parsedURI.Scheme == "" {
		logger.Fatal(errors.New("provided URI missing scheme"))
	}

	user := &types.User{
		Username:        username,
		TwoFactorSecret: totpSecret,
		HashedPassword:  password,
	}

	cookie, cookieErr := testutils.GetLoginCookie(ctx, address, user)
	if cookieErr != nil {
		logger.Fatal(fmt.Errorf("getting cookie: %w", cookieErr))
	}

	client, err := httpclient.NewClient(parsedURI, trace.NewNoopTracerProvider(), httpclient.UsingLogger(logger), httpclient.UsingCookie(cookie))
	if err != nil {
		logger.Fatal(fmt.Errorf("initializing client: %w", err))
	}

	logger.Debug("initialized API client")

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		for _, instrument := range validInstruments {
			if _, instrumentCreationErr := client.CreateValidInstrument(ctx, instrument); instrumentCreationErr != nil {
				logger.Error(instrumentCreationErr, "creating valid instrument")
			}
		}
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		for _, preparation := range validPreparations {
			if _, preparationCreationErr := client.CreateValidPreparation(ctx, preparation); preparationCreationErr != nil {
				logger.Error(preparationCreationErr, "creating valid preparation")
			}
		}
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		for _, ingredient := range validIngredients {
			if _, ingredientCreationErr := client.CreateValidIngredient(ctx, ingredient); ingredientCreationErr != nil {
				logger.Error(ingredientCreationErr, "creating valid ingredient")
			}
		}
		wg.Done()
	}()

	wg.Wait()

	for _, recipe := range recipes {
		if _, recipeCreationErr := client.CreateRecipe(ctx, recipe); recipeCreationErr != nil {
			logger.Error(recipeCreationErr, "creating valid recipe")
		}
	}
}
