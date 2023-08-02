package postgres

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"io"
	"log"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/dinnerdonebetter/backend/internal/database"
	"github.com/dinnerdonebetter/backend/internal/database/config"
	"github.com/dinnerdonebetter/backend/internal/identifiers"
	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/pkg/cryptography/aes"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Masterminds/squirrel"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

var _ sqlmock.Argument = (*idMatcher)(nil)

type idMatcher struct{}

func (s *idMatcher) Match(v driver.Value) bool {
	x, ok := v.(string)
	if !ok {
		return false
	}

	if err := identifiers.Validate(x); err != nil {
		return false
	}

	return true
}

func assertArgCountMatchesQuery(t *testing.T, query string, args []any) {
	t.Helper()

	queryArgCount := len(regexp.MustCompile(`\$\d+`).FindAllString(query, -1))

	if len(args) > 0 {
		assert.Equal(t, queryArgCount, len(args))
	} else {
		assert.Zero(t, queryArgCount)
	}
}

func newArbitraryDatabaseResult() driver.Result {
	return sqlmock.NewResult(1, 1)
}

func formatQueryForSQLMock(query string) string {
	return strings.NewReplacer(
		"$", `\$`,
		"(", `\(`,
		")", `\)`,
		"=", `\=`,
		"*", `\*`,
		".", `\.`,
		"+", `\+`,
		"?", `\?`,
		",", `\,`,
		"-", `\-`,
		"[", `\[`,
		"]", `\]`,
	).Replace(query)
}

func interfaceToDriverValue(in []any) []driver.Value {
	out := []driver.Value{}

	for _, x := range in {
		out = append(out, driver.Value(x))
	}

	return out
}

func buildInvalidMockRowsFromListOfIDs(ids []string) *sqlmock.Rows {
	exampleRows := sqlmock.NewRows([]string{"id"})

	for range ids {
		exampleRows.AddRow(nil)
	}

	return exampleRows
}

type sqlmockExpecterWrapper struct {
	sqlmock.Sqlmock
}

func (e *sqlmockExpecterWrapper) AssertExpectations(t mock.TestingT) bool {
	return assert.NoError(t, e.Sqlmock.ExpectationsWereMet(), "not all database expectations were met")
}

func buildTestClient(t *testing.T) (*Querier, *sqlmockExpecterWrapper) {
	t.Helper()

	fakeDB, sqlMock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	require.NoError(t, err)

	encDec, err := aes.NewEncryptorDecryptor(tracing.NewNoopTracerProvider(), logging.NewNoopLogger(), []byte("blahblahblahblahblahblahblahblah"))
	require.NoError(t, err)

	c := &Querier{
		db:                      fakeDB,
		logQueries:              false,
		logger:                  logging.NewNoopLogger(),
		timeFunc:                defaultTimeFunc,
		tracer:                  tracing.NewTracerForTest("test"),
		sqlBuilder:              squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
		oauth2ClientTokenEncDec: encDec,
	}

	return c, &sqlmockExpecterWrapper{Sqlmock: sqlMock}
}

const (
	defaultImage            = "postgres:15"
	defaultDatabaseName     = "dinner-done-better"
	defaultDatabaseUsername = "dbuser"
	defaultDatabasePassword = "hunter2"
)

func buildDatabaseClientForTest(t *testing.T, ctx context.Context) (*Querier, *postgres.PostgresContainer) {
	t.Helper()

	testcontainers.Logger = log.New(io.Discard, "", log.LstdFlags)

	container, err := postgres.RunContainer(
		ctx,
		testcontainers.WithImage(defaultImage),
		postgres.WithDatabase(defaultDatabaseName),
		postgres.WithUsername(defaultDatabaseUsername),
		postgres.WithPassword(defaultDatabasePassword),
		testcontainers.WithWaitStrategyAndDeadline(
			time.Minute,
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2),
		),
	)

	require.NoError(t, err)
	require.NotNil(t, container)

	connStr, err := container.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	dbc, err := ProvideDatabaseClient(ctx, logging.NewNoopLogger(), &config.Config{ConnectionDetails: connStr, RunMigrations: true, OAuth2TokenEncryptionKey: "blahblahblahblahblahblahblahblah"}, tracing.NewNoopTracerProvider())
	require.NoError(t, err)
	require.NotNil(t, dbc)

	return dbc.(*Querier), container
}

func buildErroneousMockRow() *sqlmock.Rows {
	exampleRows := sqlmock.NewRows([]string{"columns", "don't", "match", "lol"}).AddRow(
		"doesn't",
		"matter",
		"what",
		"goes",
	)

	return exampleRows
}

// end helper funcs

func TestQuerier_IsReady(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, db := buildTestClient(t)

		db.ExpectPing().WillDelayFor(0)

		assert.True(t, c.IsReady(ctx, time.Second, 1))
	})

	T.Run("with error pinging database", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, db := buildTestClient(t)

		db.ExpectPing().WillReturnError(errors.New("blah"))

		assert.False(t, c.IsReady(ctx, time.Second, 1))
	})

	T.Run("exhausting all available queries", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		c, db := buildTestClient(t)

		c.IsReady(ctx, time.Second, 1)

		db.ExpectPing().WillReturnError(errors.New("blah"))

		assert.False(t, c.IsReady(ctx, time.Second, 1))
	})
}

func TestProvideDatabaseClient(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleConfig := &config.Config{
			Debug:                    true,
			OAuth2TokenEncryptionKey: "blahblahblahblahblahblahblahblah",
			RunMigrations:            false,
			MaxPingAttempts:          1,
		}

		actual, err := ProvideDatabaseClient(ctx, logging.NewNoopLogger(), exampleConfig, tracing.NewNoopTracerProvider())
		assert.NotNil(t, actual)
		assert.NoError(t, err)
	})
}

func TestDefaultTimeFunc(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		assert.NotZero(t, defaultTimeFunc())
	})
}

func TestQuerier_currentTime(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		c, _ := buildTestClient(t)

		assert.NotEmpty(t, c.currentTime())
	})

	T.Run("handles nil", func(t *testing.T) {
		t.Parallel()

		var c *Querier

		assert.NotEmpty(t, c.currentTime())
	})
}

func TestQuerier_rollbackTransaction(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, db := buildTestClient(t)

		db.ExpectBegin()
		db.ExpectRollback().WillReturnError(errors.New("blah"))

		tx, err := c.db.BeginTx(ctx, nil)
		require.NoError(t, err)

		c.rollbackTransaction(ctx, tx)
	})
}

func TestQuerier_handleRows(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		mockRows := &database.MockResultIterator{}
		mockRows.On("Err").Return(nil)
		mockRows.On("Close").Return(nil)

		c, _ := buildTestClient(t)

		err := c.checkRowsForErrorAndClose(ctx, mockRows)
		assert.NoError(t, err)
	})

	T.Run("with row error", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		expected := errors.New("blah")

		mockRows := &database.MockResultIterator{}
		mockRows.On("Err").Return(expected)

		c, _ := buildTestClient(t)

		err := c.checkRowsForErrorAndClose(ctx, mockRows)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, expected))
	})

	T.Run("with close error", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		expected := errors.New("blah")

		mockRows := &database.MockResultIterator{}
		mockRows.On("Err").Return(nil)
		mockRows.On("Close").Return(expected)

		c, _ := buildTestClient(t)

		err := c.checkRowsForErrorAndClose(ctx, mockRows)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, expected))
	})
}

func TestQuerier_performCreateQueryIgnoringReturn(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, db := buildTestClient(t)

		fakeQuery, fakeArgs := fakes.BuildFakeSQLQuery()

		db.ExpectExec(formatQueryForSQLMock(fakeQuery)).
			WithArgs(interfaceToDriverValue(fakeArgs)...).
			WillReturnResult(newArbitraryDatabaseResult())

		err := c.performWriteQuery(ctx, c.db, "example", fakeQuery, fakeArgs)

		assert.NoError(t, err)
	})
}

func TestQuerier_performCreateQuery(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, db := buildTestClient(t)

		fakeQuery, fakeArgs := fakes.BuildFakeSQLQuery()

		db.ExpectExec(formatQueryForSQLMock(fakeQuery)).
			WithArgs(interfaceToDriverValue(fakeArgs)...).
			WillReturnResult(newArbitraryDatabaseResult())

		err := c.performWriteQuery(ctx, c.db, "example", fakeQuery, fakeArgs)

		assert.NoError(t, err)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, db := buildTestClient(t)

		fakeQuery, fakeArgs := fakes.BuildFakeSQLQuery()

		db.ExpectExec(formatQueryForSQLMock(fakeQuery)).
			WithArgs(interfaceToDriverValue(fakeArgs)...).
			WillReturnError(errors.New("blah"))

		err := c.performWriteQuery(ctx, c.db, "example", fakeQuery, fakeArgs)

		assert.Error(t, err)
	})

	T.Run("with no rows returned", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, db := buildTestClient(t)

		fakeQuery, fakeArgs := fakes.BuildFakeSQLQuery()

		db.ExpectExec(formatQueryForSQLMock(fakeQuery)).
			WithArgs(interfaceToDriverValue(fakeArgs)...).
			WillReturnResult(sqlmock.NewResult(int64(1), 0))

		err := c.performWriteQuery(ctx, c.db, "example", fakeQuery, fakeArgs)

		assert.Error(t, err)
		assert.True(t, errors.Is(err, sql.ErrNoRows))
	})
}

func Test_timePointerFromNullTime(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		expected := time.Now()
		actual := timePointerFromNullTime(sql.NullTime{Time: expected, Valid: true})

		assert.NotNil(t, actual)
		assert.Equal(t, expected, *actual)
	})
}

func Test_stringPointerFromNullString(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		expected := t.Name()
		actual := stringPointerFromNullString(sql.NullString{String: expected, Valid: true})

		assert.NotNil(t, actual)
		assert.Equal(t, expected, *actual)
	})
}
