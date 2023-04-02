package postgres

import (
	"context"
	"database/sql/driver"
	"errors"
	"testing"
	"time"

	"github.com/prixfixeco/backend/pkg/types"
	"github.com/prixfixeco/backend/pkg/types/converters"
	"github.com/prixfixeco/backend/pkg/types/fakes"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func buildMockRowsFromPasswordResetTokens(includeCounts bool, filteredCount uint64, tokens ...*types.PasswordResetToken) *sqlmock.Rows {
	columns := passwordResetTokensTableColumns

	if includeCounts {
		columns = append(columns, "filtered_count", "total_count")
	}

	exampleRows := sqlmock.NewRows(columns)

	for _, w := range tokens {
		rowValues := []driver.Value{
			w.ID,
			w.Token,
			w.ExpiresAt,
			w.CreatedAt,
			w.LastUpdatedAt,
			w.RedeemedAt,
			w.BelongsToUser,
		}

		if includeCounts {
			rowValues = append(rowValues, filteredCount, len(tokens))
		}

		exampleRows.AddRow(rowValues...)
	}

	return exampleRows
}

func TestSQLQuerier_GetPasswordResetTokenByToken(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleToken := fakes.BuildFakePasswordResetToken()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{
			exampleToken.Token,
		}

		db.ExpectQuery(formatQueryForSQLMock(getPasswordResetTokenQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromPasswordResetTokens(false, 0, exampleToken))

		actual, err := c.GetPasswordResetTokenByToken(ctx, exampleToken.Token)
		assert.NoError(t, err)
		assert.Equal(t, exampleToken, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with missing token", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetPasswordResetTokenByToken(ctx, "")
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		exampleToken := fakes.BuildFakePasswordResetToken()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{
			exampleToken.Token,
		}

		db.ExpectQuery(formatQueryForSQLMock(getPasswordResetTokenQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.GetPasswordResetTokenByToken(ctx, exampleToken.Token)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestSQLQuerier_CreatePasswordResetToken(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleToken := fakes.BuildFakePasswordResetToken()
		exampleInput := converters.ConvertPasswordResetTokenToPasswordResetTokenDatabaseCreationInput(exampleToken)

		ctx := context.Background()
		c, db := buildTestClient(t)

		c.timeFunc = func() time.Time {
			return exampleToken.CreatedAt
		}

		args := []any{
			exampleInput.ID,
			exampleInput.Token,
			exampleInput.BelongsToUser,
		}

		db.ExpectExec(formatQueryForSQLMock(passwordResetTokenCreationQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnResult(newArbitraryDatabaseResult())

		actual, err := c.CreatePasswordResetToken(ctx, exampleInput)
		assert.NoError(t, err)
		assert.Equal(t, exampleToken, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with missing input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.CreatePasswordResetToken(ctx, nil)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		exampleToken := fakes.BuildFakePasswordResetToken()
		exampleInput := converters.ConvertPasswordResetTokenToPasswordResetTokenDatabaseCreationInput(exampleToken)

		ctx := context.Background()
		c, db := buildTestClient(t)

		c.timeFunc = func() time.Time {
			return exampleToken.CreatedAt
		}

		args := []any{
			exampleInput.ID,
			exampleInput.Token,
			exampleInput.BelongsToUser,
		}

		db.ExpectExec(formatQueryForSQLMock(passwordResetTokenCreationQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.CreatePasswordResetToken(ctx, exampleInput)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestSQLQuerier_RedeemPasswordResetToken(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		examplePasswordResetTokenID := fakes.BuildFakeID()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{examplePasswordResetTokenID}

		db.ExpectExec(formatQueryForSQLMock(redeemPasswordResetTokenQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnResult(newArbitraryDatabaseResult())

		actual := c.RedeemPasswordResetToken(ctx, examplePasswordResetTokenID)
		assert.NoError(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with missing ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, db := buildTestClient(t)

		actual := c.RedeemPasswordResetToken(ctx, "")
		assert.Error(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		examplePasswordResetTokenID := fakes.BuildFakeID()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{examplePasswordResetTokenID}

		db.ExpectExec(formatQueryForSQLMock(redeemPasswordResetTokenQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual := c.RedeemPasswordResetToken(ctx, examplePasswordResetTokenID)
		assert.Error(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}
