package tableaccess

import (
	"context"
	"database/sql"
	"fmt"
	"hash/fnv"
	"io"
	"log"
	"net"
	"testing"
	"time"

	"github.com/dinnerdonebetter/backend/internal/lib/pointer"

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

	db, err := sql.Open("pgx", container.MustConnectionString(ctx, "sslmode=disable"))
	require.NoError(t, err)

	return db, container
}

func createDatabaseForTest(t *testing.T, db *sql.DB, database string) {
	t.Helper()

	_, err := db.Exec(fmt.Sprintf("CREATE DATABASE %s;", database))
	require.NoError(t, err)
}

func createDatabaseUserForTest(t *testing.T, db *sql.DB, username, password string) {
	t.Helper()

	_, err := db.Exec(fmt.Sprintf("CREATE USER %s WITH PASSWORD '%s';", username, password))
	require.NoError(t, err)
}

func grantRolesToUserForTest(t *testing.T, db *sql.DB, username string) {
	t.Helper()

	_, err := db.Exec(fmt.Sprintf("GRANT pg_read_all_data TO %s;", username))
	require.NoError(t, err)

	_, err = db.Exec(fmt.Sprintf("GRANT pg_write_all_data TO %s;", username))
	require.NoError(t, err)
}

func TestNewManager(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.SkipNow() // experimental
		t.Parallel()

		ctx := t.Context()

		adminDB, container := buildDatabaseConnectionForTest(t, ctx)
		defer container.Stop(ctx, pointer.To(10*time.Second))

		connStr, err := container.ConnectionString(ctx, "sslmode=disable")
		require.NoError(t, err)
		require.NotEmpty(t, connStr)

		const (
			databaseName = "exampledb"
			username     = "username"
			password     = "password123"
		)

		createDatabaseForTest(t, adminDB, databaseName)
		createDatabaseUserForTest(t, adminDB, username, password)
		grantRolesToUserForTest(t, adminDB, username)

		dbHost, err := container.Host(ctx)
		require.NoError(t, err)
		dbPort, err := container.MappedPort(ctx, "5432")
		require.NoError(t, err)

		createdUserDB, err := sql.Open("pgx", fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", username, password, net.JoinHostPort(dbHost, dbPort.Port()), databaseName))
		require.NoError(t, err)

		_, err = createdUserDB.ExecContext(ctx, `
CREATE TABLE IF NOT EXISTS example (
    id TEXT NOT NULL PRIMARY KEY,
    name TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    last_updated_at TIMESTAMP WITH TIME ZONE,
    archived_at TIMESTAMP WITH TIME ZONE
);
`)
		require.NoError(t, err)

		println("")
	})
}
