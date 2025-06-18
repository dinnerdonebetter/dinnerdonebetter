package testing

import (
	"context"
	"fmt"
	"hash/fnv"
	"io"
	"log"
	"testing"
	"time"

	databasecfg "github.com/dinnerdonebetter/backend/internal/platform/database/config"

	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"gopkg.in/matryer/try.v1"
)

func hashStringToNumber(s string) uint64 {
	// Create a new FNV-1a 64-bit hash object
	h := fnv.New64a()

	// Write the bytes of the string into the hash object
	_, err := h.Write([]byte(s))
	if err != nil {
		// Handle error if necessary
		panic(err)
	}

	// Return the resulting hash value as a number (uint64)
	return h.Sum64()
}

func reverseString(input string) string {
	runes := []rune(input)
	length := len(runes)

	for i, j := 0, length-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}

	return string(runes)
}

func splitReverseConcat(input string) string {
	length := len(input)
	halfLength := length / 2

	firstHalf := input[:halfLength]
	secondHalf := input[halfLength:]

	reversedFirstHalf := reverseString(firstHalf)
	reversedSecondHalf := reverseString(secondHalf)

	return reversedSecondHalf + reversedFirstHalf
}

const (
	defaultPostgresImage = "postgres:17"
)

func BuildDatabaseClientForTest(t *testing.T, ctx context.Context) *postgres.PostgresContainer {
	t.Helper()

	dbUsername := fmt.Sprintf("%d", hashStringToNumber(t.Name()))
	testcontainers.Logger = log.New(io.Discard, "", log.LstdFlags)

	var container *postgres.PostgresContainer
	err := try.Do(func(attempt int) (bool, error) {
		var containerErr error
		container, containerErr = postgres.Run(
			ctx,
			defaultPostgresImage,
			postgres.WithDatabase(splitReverseConcat(dbUsername)),
			postgres.WithUsername(dbUsername),
			postgres.WithPassword(reverseString(dbUsername)),
			testcontainers.WithWaitStrategyAndDeadline(2*time.Minute, wait.ForLog("database system is ready to accept connections").WithOccurrence(2)),
		)

		return attempt < 5, containerErr
	})
	require.NoError(t, err)
	require.NotNil(t, container)

	connStr, err := container.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	dbConfig := &databasecfg.Config{
		RunMigrations:            true,
		OAuth2TokenEncryptionKey: "blahblahblahblahblahblahblahblah",
	}
	require.NoError(t, dbConfig.LoadConnectionDetailsFromURL(connStr))

	return container
}

/*

func createUserForTest(t *testing.T, ctx context.Context, exampleUser *types.User, dbc *Querier) *types.User {
	t.Helper()

	// create
	if exampleUser == nil {
		exampleUser = fakes.BuildFakeUser()
	}
	dbInput := converters.ConvertUserToUserDatabaseCreationInput(exampleUser)

	exampleUser.TwoFactorSecretVerifiedAt = nil
	created, err := dbc.CreateUser(ctx, dbInput)
	exampleUser.CreatedAt = created.CreatedAt
	exampleUser.TwoFactorSecretVerifiedAt = created.TwoFactorSecretVerifiedAt
	assert.NoError(t, err)
	assert.Equal(t, exampleUser, created)

	user, err := dbc.GetUser(ctx, created.ID)
	exampleUser.CreatedAt = user.CreatedAt
	exampleUser.Birthday = user.Birthday

	assert.NoError(t, err)
	assert.Equal(t, user, exampleUser)

	return created
}

func createAccountForTest(t *testing.T, ctx context.Context, exampleAccount *types.Account, dbc *Querier) *types.Account {
	t.Helper()

	// create
	if exampleAccount == nil {
		exampleUser := createUserForTest(t, ctx, nil, dbc)
		exampleAccount = fakes.BuildFakeAccount()
		exampleAccount.BelongsToUser = exampleUser.ID
	}
	exampleAccount.PaymentProcessorCustomerID = ""
	exampleAccount.Members = nil
	dbInput := converters.ConvertAccountToAccountDatabaseCreationInput(exampleAccount)

	created, err := dbc.CreateAccount(ctx, dbInput)
	assert.NoError(t, err)
	require.NotNil(t, created)
	exampleAccount.CreatedAt = created.CreatedAt
	exampleAccount.WebhookEncryptionKey = created.WebhookEncryptionKey
	assert.Equal(t, exampleAccount, created)

	account, err := dbc.GetAccount(ctx, created.ID)
	require.NoError(t, err)
	require.NotNil(t, account)

	exampleAccount.CreatedAt = account.CreatedAt
	exampleAccount.Members = account.Members
	exampleAccount.WebhookEncryptionKey = account.WebhookEncryptionKey

	assert.Equal(t, exampleAccount, account)

	return created
}

*/
