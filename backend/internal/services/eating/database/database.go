package database

import (
	"context"
	"database/sql"
	"errors"
	"io"

	"github.com/dinnerdonebetter/backend/internal/services/eating/types"
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
		Scan(dest ...any) error
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
		ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
		PrepareContext(context.Context, string) (*sql.Stmt, error)
		QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
		QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
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

	// DataManager describes anything that stores data for our services.
	DataManager interface {
		DB() *sql.DB
		Close()
		Migrate(ctx context.Context) error
		IsReady(ctx context.Context) (ready bool)

		types.MealPlanTaskDataManager
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
		types.ValidMeasurementUnitConversionDataManager
		types.RecipeMediaDataManager
		types.ValidIngredientStateDataManager
		types.RecipeStepCompletionConditionDataManager
		types.ValidIngredientStateIngredientDataManager
		types.RecipeStepVesselDataManager
		types.ValidIngredientGroupDataManager
		types.RecipeRatingDataManager
		types.ValidVesselDataManager
		types.ValidPreparationVesselDataManager
		types.IngredientPreferenceDataManager
		types.InstrumentOwnershipDataManager
	}
)
