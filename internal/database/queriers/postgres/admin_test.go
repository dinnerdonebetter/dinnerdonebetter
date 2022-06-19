package postgres

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/prixfixeco/api_server/pkg/types"
	"github.com/prixfixeco/api_server/pkg/types/fakes"
)

func TestQuerier_UpdateUserAccountStatus(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleUser := fakes.BuildFakeUser()
		exampleInput := &types.UserAccountStatusUpdateInput{
			TargetUserID: exampleUser.ID,
			NewStatus:    "new",
			Reason:       "because",
		}

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []interface{}{
			exampleInput.NewStatus,
			exampleInput.Reason,
			exampleInput.TargetUserID,
		}

		db.ExpectExec(formatQueryForSQLMock(setUserAccountStatusQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnResult(newArbitraryDatabaseResult())

		assert.NoError(t, c.UpdateUserAccountStatus(ctx, exampleUser.ID, exampleInput))

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		exampleUser := fakes.BuildFakeUser()
		exampleInput := &types.UserAccountStatusUpdateInput{
			TargetUserID: exampleUser.ID,
			NewStatus:    "new",
			Reason:       "because",
		}

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []interface{}{
			exampleInput.NewStatus,
			exampleInput.Reason,
			exampleInput.TargetUserID,
		}

		db.ExpectExec(formatQueryForSQLMock(setUserAccountStatusQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		assert.Error(t, c.UpdateUserAccountStatus(ctx, exampleUser.ID, exampleInput))

		mock.AssertExpectationsForObjects(t, db)
	})
}
