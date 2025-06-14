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

	"github.com/dinnerdonebetter/backend/internal/platform/pointer"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/assert"
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

func buildConnectionString(t *testing.T, container *postgres.PostgresContainer, dbName, username, password string) string {
	t.Helper()
	ctx := t.Context()

	containerPort, err := container.MappedPort(ctx, "5432/tcp")
	require.NoError(t, err)

	host, err := container.Host(ctx)
	require.NoError(t, err)

	return fmt.Sprintf("postgres://%s:%s@%s/%s", username, password, net.JoinHostPort(host, containerPort.Port()), dbName)
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

func TestNewManager(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		adminDB, container := buildDatabaseConnectionForTest(t, ctx)
		defer container.Stop(ctx, pointer.To(10*time.Second))

		mgr := NewManager(adminDB)

		username := "example"
		password := "hunter2"
		databaseName := "records"

		assert.NoError(t, mgr.CreateUser(ctx, username, password))
		assert.NoError(t, mgr.CreateDatabase(ctx, databaseName, username))

		canAccess, err := mgr.UserCanAccessDatabase(ctx, username, databaseName)
		assert.NoError(t, err)
		assert.True(t, canAccess)

		db2, err := sql.Open("pgx", buildConnectionString(t, container, databaseName, username, password))
		require.NoError(t, err)

		var dbName string
		db2.QueryRowContext(ctx, `SELECT current_database()`).Scan(&dbName)
		assert.Equal(t, databaseName, dbName)
	})
}
