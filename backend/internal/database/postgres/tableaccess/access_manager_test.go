package tableaccess

import (
	"context"
	"database/sql"
	"fmt"
	"hash/fnv"
	"io"
	"log"
	"testing"
	"time"

	"github.com/dinnerdonebetter/backend/internal/pointer"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"gopkg.in/matryer/try.v1"
)

// TODO: lots of duplication with the upper postgres package

const (
	defaultPostgresImage = "postgres:17"
)

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

func buildDatabaseConnectionForTest(t *testing.T, ctx context.Context) (*sql.DB, *postgres.PostgresContainer) {
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

	db, err := sql.Open("pgx", connStr)
	require.NoError(t, err)

	return db, container
}

func createDatabaseForTest(t *testing.T, db *sql.DB, database string) {
	t.Helper()

	res, err := db.Exec(fmt.Sprintf("CREATE DATABASE %s", database))
	require.NoError(t, err)
	rowsAffected, err := res.RowsAffected()
	require.NoError(t, err)
	require.Equal(t, 1, rowsAffected)
}

func createDatabaseUserForTest(t *testing.T, db *sql.DB, username, password string) {
	t.Helper()

	res, err := db.Exec(fmt.Sprintf("CREATE USER %s WITH ENCRYPTED PASSWORD '%s'", username, password))
	require.NoError(t, err)
	rowsAffected, err := res.RowsAffected()
	require.NoError(t, err)
	require.Equal(t, 1, rowsAffected)
}

func grantAllPrivilegesToDatabaseForTest(t *testing.T, db *sql.DB, username, database string) {
	t.Helper()

	res, err := db.Exec(fmt.Sprintf("GRANT ALL PRIVILEGES ON DATABASE %s TO %s", database, username))
	require.NoError(t, err)
	rowsAffected, err := res.RowsAffected()
	require.NoError(t, err)
	require.Equal(t, 1, rowsAffected)
}

func revokeAllPrivilegesToDatabaseForTest(t *testing.T, db *sql.DB, username, database string) {
	t.Helper()

	res, err := db.Exec(fmt.Sprintf("REVOKE ALL PRIVILEGES ON DATABASE %s FROM %s", database, username))
	require.NoError(t, err)
	rowsAffected, err := res.RowsAffected()
	require.NoError(t, err)
	require.Equal(t, 1, rowsAffected)
}

func TestNewManager(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		db, container := buildDatabaseConnectionForTest(t, ctx)
		defer container.Stop(ctx, pointer.To(10*time.Second))

		connStr, err := container.ConnectionString(ctx, "sslmode=disable")
		require.NoError(t, err)
		require.NotEmpty(t, connStr)

		const (
			databaseName = "example"
			username     = "user"
			password     = "password123"
		)

		createDatabaseForTest(t, db, databaseName)
		createDatabaseUserForTest(t, db, username, password)
		grantAllPrivilegesToDatabaseForTest(t, db, username, databaseName)
		println("")
	})
}
