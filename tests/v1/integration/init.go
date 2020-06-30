package integration

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	client "gitlab.com/prixfixe/prixfixe/client/v1/http"
	"gitlab.com/prixfixe/prixfixe/database/v1"
	"gitlab.com/prixfixe/prixfixe/database/v1/queriers/postgres"
	models "gitlab.com/prixfixe/prixfixe/models/v1"
	"gitlab.com/prixfixe/prixfixe/tests/v1/testutil"

	"github.com/Masterminds/squirrel"
	"gitlab.com/verygoodsoftwarenotvirus/logging/v1/zerolog"
)

const (
	debug         = true
	nonexistentID = 999999999
)

var (
	urlToUse       string
	prixfixeClient *client.V1Client
)

func init() {
	urlToUse = testutil.DetermineServiceURL()
	logger := zerolog.NewZeroLogger()

	logger.WithValue("url", urlToUse).Info("checking server")
	testutil.EnsureServerIsUp(urlToUse)

	ogUser, userCreationErr := testutil.CreateObligatoryUser(urlToUse, debug)
	if userCreationErr != nil {
		logger.Fatal(userCreationErr)
	}

	// make the user an admin
	dbURL := testutil.DetermineDatabaseURL()
	db, dbConnectionErr := postgres.ProvidePostgresDB(logger, database.ConnectionDetails(dbURL))
	if dbConnectionErr != nil {
		logger.Fatal(dbConnectionErr)
	}

	makeAdminQuery, makeAdminArgs, queryCreationErr := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).
		Update("users").
		Set("is_admin", true).
		Where(squirrel.Eq{
			"id": ogUser.ID,
		}).
		Suffix("RETURNING updated_on").
		ToSql()
	if queryCreationErr != nil {
		logger.Fatal(queryCreationErr)
	}

	_, userModificationErr := db.Exec(makeAdminQuery, makeAdminArgs...)
	if userModificationErr != nil {
		logger.Fatal(userModificationErr)
	}

	if dbCloseErr := db.Close(); dbCloseErr != nil {
		logger.Fatal(dbCloseErr)
	}

	oa2Client, clientCreationErr := testutil.CreateObligatoryClient(urlToUse, ogUser)
	if clientCreationErr != nil {
		logger.Fatal(clientCreationErr)
	}

	prixfixeClient = initializeClient(oa2Client)

	fiftySpaces := strings.Repeat("\n", 50)
	fmt.Printf("%s\tRunning tests%s", fiftySpaces, fiftySpaces)
}

func buildHTTPClient() *http.Client {
	httpc := &http.Client{
		Transport: http.DefaultTransport,
		Timeout:   5 * time.Second,
	}

	return httpc
}

func initializeClient(oa2Client *models.OAuth2Client) *client.V1Client {
	uri, err := url.Parse(urlToUse)
	if err != nil {
		panic(err)
	}

	c, err := client.NewClient(
		context.Background(),
		oa2Client.ClientID,
		oa2Client.ClientSecret,
		uri,
		zerolog.NewZeroLogger(),
		buildHTTPClient(),
		oa2Client.Scopes,
		debug,
	)
	if err != nil {
		panic(err)
	}
	return c
}
