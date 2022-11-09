package database

import (
	"context"
	"database/sql"
	"errors"
	"io"
	"time"

	"github.com/alexedwards/scs/v2"

	"github.com/prixfixeco/backend/pkg/types"
)

var (
	// ErrDatabaseNotReady indicates the given database is not ready.
	ErrDatabaseNotReady = errors.New("database is not ready yet")

	// ErrUserAlreadyExists indicates that a user with that username has already been created.
	ErrUserAlreadyExists = errors.New("user already exists")
)

type (
	// Scanner represents any database response (i.e. sql.Row[s]).
	Scanner interface {
		Scan(dest ...interface{}) error
	}

	// ResultIterator represents any iterable database response (i.e. sql.Rows).
	ResultIterator interface {
		Next() bool
		Err() error
		Scanner
		io.Closer
	}

	// SQLQueryExecutor is a subset interface for sql.{DB|Tx} objects.
	SQLQueryExecutor interface {
		ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
		PrepareContext(context.Context, string) (*sql.Stmt, error)
		QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
		QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	}

	// SQLTransactionManager is a subset interface for sql.{DB|Tx} objects.
	SQLTransactionManager interface {
		Rollback() error
	}

	// SQLQueryExecutorAndTransactionManager is a subset interface for sql.{DB|Tx} objects.
	SQLQueryExecutorAndTransactionManager interface {
		SQLQueryExecutor
		SQLTransactionManager
	}

	// MetricsCollectionInterval defines the interval at which we collect database metrics.
	MetricsCollectionInterval time.Duration

	// ConnectionDetails is a string alias for dependency injection.
	ConnectionDetails string

	// DataManager describes anything that stores data for our services.
	DataManager interface {
		DB() *sql.DB
		Migrate(ctx context.Context, maxAttempts uint8) error
		IsReady(ctx context.Context, maxAttempts uint8) (ready bool)
		ProvideSessionStore() scs.Store

		types.MealPlanTaskDataManager
		types.AdminUserDataManager
		types.HouseholdDataManager
		types.HouseholdInvitationDataManager
		types.HouseholdUserMembershipDataManager
		types.UserDataManager
		types.APIClientDataManager
		types.PasswordResetTokenDataManager
		types.WebhookDataManager
		types.ValidInstrumentDataManager
		types.ValidIngredientDataManager
		types.ValidPreparationDataManager
		types.ValidIngredientPreparationDataManager
		types.MealDataManager
		types.RecipeDataManager
		types.RecipeStepDataManager
		types.RecipeStepProductDataManager
		types.RecipeStepInstrumentDataManager
		types.RecipeStepIngredientDataManager
		types.MealPlanDataManager
		types.MealPlanEventDataManager
		types.MealPlanOptionDataManager
		types.MealPlanOptionVoteDataManager
		types.ValidMeasurementUnitDataManager
		types.ValidPreparationInstrumentDataManager
		types.ValidIngredientMeasurementUnitDataManager
		types.RecipePrepTaskDataManager
		types.MealPlanGroceryListItemDataManager
		types.ValidMeasurementConversionDataManager
		types.RecipeMediaDataManager
	}
)
